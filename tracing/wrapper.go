package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"

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

func SetSpanWithFields(span opentracing.Span, fields map[string]interface{}) {
	for k, v := range fields {
		tk := k
		tv := v
		span.SetTag(tk, tv)
	}
}

func DoF(ctx context.Context, operationName string, f func() error, fields map[string]interface{}) (err error) {
	span, ctx := New(ctx, operationName)
	SetSpanWithFields(span, fields)
	err = f()
	if err != nil {
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func DoCtxF(ctx context.Context, operationName string, f func(context.Context) error, fields map[string]interface{}) (err error) {
	span, ctx := New(ctx, operationName)
	SetSpanWithFields(span, fields)
	err = f(ctx)
	if err != nil {
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}
