package grpc

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/cntechpower/utils/tracing"
)

func WithTrace() InterceptorOption {
	return newFuncServerOption(func(options *interceptorOptions) {
		options.traceEnable = true
	})
}

func inject(method string, ctx context.Context) (ctxNew context.Context) {
	ctxNew = ctx
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span, ctxNew = tracing.New(ctx, method)
	}
	return
}
