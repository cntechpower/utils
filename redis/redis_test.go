package redis

import (
	"context"
	"testing"

	"github.com/cntechpower/utils/tracing"

	"github.com/go-redis/redis/v8"
)

func TestRedisTracing(t *testing.T) {
	tracing.Init("unit-test", "10.0.0.2:6831")
	defer tracing.Close()

	cli := New(&redis.Options{
		Addr:     "10.0.0.2:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer func() {
		_ = cli.Close()
	}()

	cli.Del(context.Background(), "abc")
}
