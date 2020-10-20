package go_test_redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type waitOptions struct {
	timeout time.Duration
}

type waitOptionFn func(o *waitOptions)

// WithTimeout overwrite default timeout to wait for redis availability in
// WaitForRedis function
func WithTimeout(timeout time.Duration) waitOptionFn {
	return func(o *waitOptions) {
		o.timeout = timeout
	}
}

// WaitForRedis is useful to use in TestMain function to wait until
// redis would be available. It may be used if redis is not running
// all the time and is starting in parallel with tests. Default timeout to
// wait for redis is 5 seconds. It may be overwritten using WithTimeout option.
func WaitForRedis(ops ...waitOptionFn) error {
	var options = waitOptions{timeout: 5 * time.Second}
	for _, fn := range ops {
		fn(&options)
	}

	redisAddr := os.Getenv("REDISADDR")
	if redisAddr == "" {
		return errors.New("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.timeout)
	defer cancel()

	err := waitForSocket(ctx, redisAddr)
	if err != nil {
		return err
	}

	return waitRedisLoaded(ctx)
}

func waitRedisLoaded(ctx context.Context) (err error) {
	cli := redis.NewClient(newRedisOpts(0))
	defer func() {
		err2 := cli.Close()
		if err2 != nil && err == nil {
			err = err2
		}
	}()

	tmMin := 50 * time.Millisecond
	tmMax := time.Second

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for redis loading failed: %w", ctx.Err())
		default:
		}

		r, err := cli.Info(ctx, "persistence").Result()
		if err != nil {
			return err
		}

		if loading := parseInfoResponse(r)["loading"]; loading == "0" {
			return nil
		}

		time.Sleep(tmMin)
		tmMin *= 2
		if tmMin > tmMax {
			tmMin = tmMax
		}
	}
}

func parseInfoResponse(in string) map[string]string {
	lines := strings.Split(in, "\n")
	result := make(map[string]string, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) < 2 || l[0] == '#' {
			continue
		}
		idx := strings.IndexRune(l, ':')
		if idx <= 0 {
			continue
		}
		result[strings.TrimSpace(l[:idx])] = strings.TrimSpace(l[idx+1:])
	}
	return result
}

func waitForSocket(ctx context.Context, addr string) error {
	var (
		err   error
		conn  net.Conn
		tmMin = 50 * time.Millisecond
		tmMax = time.Second
	)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for %v failed: %w", addr, err)
		default:
		}

		conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		if err != nil {
			time.Sleep(tmMin)
			tmMin *= 2
			if tmMin > tmMax {
				tmMin = tmMax
			}
			continue
		}

		return conn.Close()
	}
}
