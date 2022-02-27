package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cntechpower/utils/tracing"

	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/ext"
)

type Interface interface {
	Key() string
	GetFromDB(ctx context.Context) error
}

var mem gcache.Cache
var cache *redis.Client

func Init(r *redis.Client) {
	mem = gcache.New(100).LFU().
		Build()
	cache = r
}

func Get(ctx context.Context, arg Interface) (err error) {
	span, ctx := tracing.New(ctx, "cache.Get")
	var bs []byte
	var shouldUpdateMem, shouldUpdateRedis bool
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			ext.LogError(span, err)
			return
		}
		if !shouldUpdateRedis && !shouldUpdateMem {
			span.Finish()
			return
		}
		bs, _ = json.Marshal(arg)
		if shouldUpdateMem {
			span, _ := tracing.New(ctx, "set_mem")
			err1 := mem.Set(arg.Key(), bs)
			if err1 != nil {
				ext.Error.Set(span, true)
				ext.LogError(span, err1)
			}
			span.Finish()

		}
		if shouldUpdateRedis {
			span, _ := tracing.New(ctx, "set_redis")
			err1 := cache.Set(ctx, arg.Key(), bs, time.Second*60).Err()
			if err1 != nil {
				ext.Error.Set(span, true)
				ext.LogError(span, err1)
			}
			span.Finish()
		}
		span.Finish()
	}()
	if mem != nil {
		span, _ := tracing.New(ctx, "get_from_mem")
		var res interface{}
		var ok bool
		res, err = mem.Get(arg.Key())
		if err != nil {
			ext.Error.Set(span, true)
			ext.LogError(span, err)
		} else {
			bs, ok = res.([]byte)
			if ok {
				err = json.Unmarshal(bs, arg)
				if err != nil {
					ext.Error.Set(span, true)
					ext.LogError(span, err)
				}
			}
		}

		span.Finish()
		if err == nil {
			return
		}
	}
	shouldUpdateMem = mem != nil

	if cache != nil {
		span, _ := tracing.New(ctx, "get_from_redis")
		bs, err = cache.Get(ctx, arg.Key()).Bytes()
		if err != nil {
			ext.Error.Set(span, true)
			ext.LogError(span, err)
		} else {
			err = json.Unmarshal(bs, arg)
			if err != nil {
				ext.Error.Set(span, true)
				ext.LogError(span, err)
			}
		}
		span.Finish()
		if err == nil {
			return
		}
	}
	shouldUpdateRedis = cache != nil

	err = arg.GetFromDB(ctx)
	return
}
