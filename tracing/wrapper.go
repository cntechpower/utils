package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go/ext"
)

func Do(ctx context.Context, operationName string, f func() error) (err error) {
	span, ctx := New(ctx, operationName)
	err = f()
	if err != nil {
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func DoCtx(ctx context.Context, operationName string, f func(context.Context) error) (err error) {
	span, ctx := New(ctx, operationName)
	err = f(ctx)
	if err != nil {
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}
