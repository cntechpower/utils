package grpc

import (
	"context"

	"github.com/cntechpower/utils/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

func WithTrace() InterceptorOption {
	return newFuncServerOption(func(options *interceptorOptions) {
		options.traceEnable = true
	})
}

func serverTryExtract(method string, ctx context.Context) (ctxNew context.Context, span opentracing.Span) {
	ctxNew = ctx
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	clientContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
	if err == nil {
		span = opentracing.GlobalTracer().StartSpan(
			method, ext.RPCServerOption(clientContext),
			gRPCComponentTag,
		)
	} else {
		span, ctxNew = tracing.New(ctx, method)
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, metadataReaderWriter{md})
		if err == nil {
			ctxNew = metadata.NewIncomingContext(ctxNew, md)
		}
	}
	ctxNew = context.WithValue(ctxNew, tracing.BackupActiveSpanKey, span)
	return
}

func clientTryInject(method string, ctx context.Context) (ctxNew context.Context, span opentracing.Span) {
	ctxNew = ctx
	span = tracing.SpanFromContext(ctx)
	if span == nil {
		span, ctxNew = tracing.New(ctx, method)
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, metadataReaderWriter{md})
	if err == nil {
		ctxNew = metadata.NewOutgoingContext(ctx, md)
	}
	return
}
