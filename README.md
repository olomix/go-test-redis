# go-test-redis
Helper tool to test Go programs against Redis

`go-test-redis` expects the redis server address in REDISADDR
environment variable.

Example:

```go
import go_test_redis "github.com/olomix/go-test-redis"

func TestCache(t *testing.T) {
	rdb := go_test_redis.WithRedis(t)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
}
```

`WithRedis` looking for the empty database and locks it to prevent other
parallel tests to use the same database. If all databases are busy,
we are waiting for free one.

When we need to wait for redis to be available in CI if it starting
in paralle container, we can use `WaitForRedis` helper. Example:

```go
func TestMain(m *testing.M) {
	err :=
		go_test_redis.WaitForRedis(go_test_redis.WithTimeout(30 * time.Second))
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
```

If we skip `WithTimeout` option, 5 seconds is the default one.