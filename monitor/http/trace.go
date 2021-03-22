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

func t(ctx *gin.Context) {
	var span opentracing.Span
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err == nil {
		span = opentracing.GlobalTracer().StartSpan(ctx.Request.RequestURI, ext.RPCServerOption(spanCtx))
	} else {
		span, _ = tracing.New(context.Background(), ctx.Request.RequestURI)
	}
	ctx.Header(tracing.TraceID, tracing.TraceIdFromSpan(span))
}
