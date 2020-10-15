package go_test_redis

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

const broadcastChName = "redis-test-broadcast"
const lockDBKeyTmpl = "redis-test-%d"
const lockTimeout = time.Minute * 20
const waitForDBTimeout = time.Minute

func closeOrFatal(t testing.TB, c io.Closer) {
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
}

func lockKeyFmt(n int) string {
	return fmt.Sprintf(lockDBKeyTmpl, n)
}

type testRedisOptions struct {
	debug            bool
	waitForDBTimeout time.Duration
}

type Option func(*testRedisOptions)

func newRedisOpts(db int) *redis.Options {
	opts := &redis.Options{DB: db}
	addr, ok := os.LookupEnv("REDISADDR")
	if ok {
		opts.Addr = addr
	}
	return opts

}

func WithRedis(t testing.TB, opts ...Option) *redis.Client {
	var op = testRedisOptions{
		waitForDBTimeout: waitForDBTimeout,
	}
	for _, setup := range opts {
		setup(&op)
	}

	cli := redis.NewClient(newRedisOpts(0))
	defer closeOrFatal(t, cli)

	n := databasesNum(t, cli)
	if n < 2 {
		t.Fatalf(
			"Minimal acceptable number of databases on redis should be 2, "+
				"currently: %v",
			n,
		)
	}
	ctx := context.Background()
	chosenDB := getOrWaitFreeDB(ctx, t, cli, n, op.waitForDBTimeout)
	if op.debug {
		t.Logf("Number of databases: %v, chosen: %v", n, chosenDB)
	}

	chosenCli := redis.NewClient(newRedisOpts(chosenDB))
	if op.debug {
		t.Logf("Return redis cli with DB = %v", currentDB(ctx, t, chosenCli))
	}

	t.Cleanup(func() {
		if op.debug {
			t.Logf("Release redis cli %v", chosenDB)
		}
		ctx := context.Background()

		if err := chosenCli.FlushDB(ctx).Err(); err != nil {
			// we should not stop here and delete lock key
			t.Errorf("can't flush db: %+v", err)
		}

		conn := chosenCli.Conn(ctx)
		defer closeOrFatal(t, conn)
		if err := conn.Select(ctx, 0).Err(); err != nil {
			t.Fatal(err)
		}
		if err := conn.Del(ctx, lockKeyFmt(chosenDB)).Err(); err != nil {
			t.Fatal(err)
		}

		err := conn.Publish(ctx, broadcastChName, strconv.Itoa(chosenDB)).Err()
		if err != nil {
			t.Fatal(err)
		}

		closeOrFatal(t, chosenCli)
	})

	return chosenCli
}

func getOrWaitFreeDB(
	ctx context.Context, t testing.TB, cli *redis.Client, dbsNum int,
	timeout time.Duration,
) int {
	pubsub := cli.Subscribe(ctx, broadcastChName)
	defer closeOrFatal(t, pubsub)

	ch := pubsub.Channel()
	ticker := time.NewTicker(5 * time.Second)
	timer := time.NewTimer(timeout)
	var chosenDB int

	// find free database
	chosenDB = lockFreeDB(ctx, t, cli, dbsNum)
	if chosenDB > 0 {
		return chosenDB
	}

	for {
		select {
		case msg := <-ch:
			// check if freed database is actually free and can be locked
			n, err := strconv.Atoi(msg.Payload)
			if err != nil {
				t.Fatal(err)
			}
			if tryLockDB(ctx, t, cli, n) {
				return n
			}
		case <-ticker.C:
			// rescan all databases, may be we will find empty one
			chosenDB = lockFreeDB(ctx, t, cli, dbsNum)
			if chosenDB > 0 {
				return chosenDB
			}
		case <-timer.C:
			t.Fatal("wait for free database timeout")
		}
	}
}

func tryLockDB(ctx context.Context, t testing.TB, cli *redis.Client, db int) bool {
	conn := cli.Conn(ctx)
	defer closeOrFatal(t, conn)

	now := time.Now().Format(time.RFC3339)
	ok, err := conn.SetNX(ctx, lockKeyFmt(db), now, lockTimeout).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		r, err := conn.Select(ctx, db).Result()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("select result [3]: %v", r)
		r, err = conn.RandomKey(ctx).Result()
		if err == redis.Nil {
			return true
		} else if err != nil {
			t.Fatal(err)
		}
		t.Fatalf(
			"expected clean database after lock released: %v, but found key %v",
			db, r,
		)
	}
	return false
}

// return:
//  -1 if no db chosen
func lockFreeDB(
	ctx context.Context, t testing.TB, cli *redis.Client, dbsNum int,
) int {
	conn := cli.Conn(ctx)
	defer closeOrFatal(t, conn)

	var foundLockedDatabases = false
	for i := 1; i < dbsNum; i++ {
		now := time.Now().Format(time.RFC3339)
		ok, err := conn.SetNX(ctx, lockKeyFmt(i), now, lockTimeout).Result()
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			if err = conn.Select(ctx, i).Err(); err != nil {
				t.Fatal(err)
			}
			if err = conn.RandomKey(ctx).Err(); err == redis.Nil {
				return i
			} else if err != nil {
				t.Fatal(err)
			}
			if err = conn.Select(ctx, 0).Err(); err != nil {
				t.Fatal(err)
			}
			if err = conn.Del(ctx, lockKeyFmt(i)).Err(); err != nil {
				t.Fatalf("can't release lock of dirty database: %v", err)
			}
		} else {
			foundLockedDatabases = true
		}
	}

	if !foundLockedDatabases {
		t.Fatal("clean databases not found, try to flush few databases")
	}

	return -1
}

var clientLineIDRegex = regexp.MustCompile(`^id=(\d+) .* db=(\d+) .*$`)

// Get currently connected database
func currentDB(ctx context.Context, t testing.TB, cli *redis.Client) int {
	var clientIDCmd *redis.IntCmd
	var clientListCmd *redis.StringCmd
	_, err := cli.Pipelined(ctx, func(p redis.Pipeliner) error {
		clientIDCmd = p.ClientID(ctx)
		clientListCmd = p.ClientList(ctx)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	clientID, err := clientIDCmd.Result()
	if err != nil {
		t.Fatal(err)
	}
	clientList, err := clientListCmd.Result()
	if err != nil {
		t.Fatal(err)
	}
	for _, statusLn := range strings.Split(clientList, "\n") {
		res := clientLineIDRegex.FindStringSubmatch(statusLn)
		if len(res) != 3 {
			continue
		}
		currentID, err := strconv.Atoi(res[1])
		if err != nil {
			t.Fatal(err)
		}
		if int64(currentID) != clientID {
			continue
		}

		dbNum, err := strconv.Atoi(res[2])
		if err != nil {
			t.Fatal(err)
		}
		return dbNum
	}

	t.Fatal("[assertion] can't find self client line")
	return 0
}

func databasesNum(t testing.TB, cli *redis.Client) int {
	paramDatabases := "databases"
	res, err := cli.ConfigGet(context.Background(), paramDatabases).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("unexpected number of returned arguments: %v", len(res))
	}
	paramName, ok := res[0].(string)
	if !ok || paramName != paramDatabases {
		t.Fatalf(
			"unexpected parameter name: %v, expected %v",
			res[0], paramDatabases,
		)
	}
	paramVal, ok := res[1].(string)
	if !ok {
		t.Fatalf("expected param value to be string(%[1]T %[1]v)", res[1])
	}
	i, err := strconv.ParseUint(paramVal, 10, 16)
	if err != nil {
		t.Fatal(err)
	}
	return int(i)
}

