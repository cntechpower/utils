package http

import (
	"context"

	"github.com/cntechpower/utils/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func WithTrace() GinMiddlewareOption {
	return newFuncServerOption(func(options *ginMiddlewareOptions) {
		options.traceEnable = true
	})
}

func inject(ctx *gin.Context) (span opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
	if err == nil {
		span = opentracing.GlobalTracer().StartSpan(ctx.Request.URL.Path, ext.RPCServerOption(spanCtx))
	} else {
		span, _ = tracing.New(context.Background(), ctx.Request.URL.Path)
	}
	ext.HTTPMethod.Set(span, ctx.Request.Method)
	ext.HTTPUrl.Set(span, ctx.Request.URL.Path)
	ctx.Set(tracing.BackupActiveSpanKey, span)
	ctx.Header(tracing.TraceID, tracing.TraceIdFromSpan(span))
	return
}
