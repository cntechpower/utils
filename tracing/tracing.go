package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	tracer       opentracing.Tracer
	tracerCloser io.Closer
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(appName, reporterAddr string) {
	cfg := &config.Configuration{
		ServiceName: appName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: reporterAddr,
			LogSpans:           true,
		},
	}
	var err error
	tracer, tracerCloser, err = cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
}

func Close() {
	_ = tracerCloser.Close()
}

// New trace instance with given operationName.
func New(ctx context.Context, operationName string) (span opentracing.Span, ctxNew context.Context) {
	span, ctxNew = opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)
	return
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// TraceIdFromContext returns the `traceId` previously associated with `ctx`, or
// `""` if not found.
func TraceIdFromContext(ctx context.Context) (traceId string) {
	span := SpanFromContext(ctx)
	if span == nil {
		return
	}
	sc, ok := span.Context().(jaeger.SpanContext)
	if ok {
		traceId = sc.TraceID().String()
	}
	return
}

// TraceIdFromSpan returns the `traceId` previously associated with `span`, or
// `""` if not found.
func TraceIdFromSpan(span opentracing.Span) (traceId string) {
	sc, ok := span.Context().(jaeger.SpanContext)
	if ok {
		traceId = sc.TraceID().String()
	}
	return
}
