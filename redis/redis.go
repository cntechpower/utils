package redis

import (
	"context"

	"github.com/cntechpower/utils/trans"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/cntechpower/utils/tracing"

	"github.com/go-redis/redis/v8"
)

const (
	dbTypeRedis = "redis"
	maxStmtLen  = 500
)

var skipCmd = []string{
	"ping",
}

func New(options *redis.Options) (cli *redis.Client) {
	cli = redis.NewClient(options)
	cli.AddHook(&hook{})
	return cli
}

type hook struct {
}

func (h *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if trans.StringInSlice(cmd.Name(), skipCmd) {
		return ctx, nil
	}
	span, ctx := tracing.New(ctx, "redis."+cmd.Name())
	ext.DBType.Set(span, dbTypeRedis)
	stmt := cmd.String()
	if len(stmt) > maxStmtLen {
		stmt = stmt[:maxStmtLen]
	}
	ext.DBStatement.Set(span, stmt)
	return ctx, nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if trans.StringInSlice(cmd.Name(), skipCmd) {
		return nil
	}
	span := tracing.SpanFromContext(ctx)
	if span != nil {
		if cmd.Err() != nil {
			ext.LogError(span, cmd.Err())
			ext.Error.Set(span, true)
		}
		span.Finish()
	}
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
