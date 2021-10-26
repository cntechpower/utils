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

var BackupActiveSpanKey = "BAS"

// Init returns an instance of Jaeger Tracer that samples probability% of traces and logs all spans to stdout.
func Init(appName, reporterAddr string) {
	cfg := &config.Configuration{
		ServiceName: appName,
		Sampler: &config.SamplerConfig{
			Type:  "probabilistic",
			Param: 0.1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: reporterAddr,
			LogSpans:           true,
		},
		Headers: &jaeger.HeadersConfig{
			TraceContextHeaderName: "trace-id",
		},
	}
	var err error
	tracer, tracerCloser, err = cfg.NewTracer(config.Logger(jaeger.NullLogger))
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
	parent := SpanFromContext(ctx)
	var opts []opentracing.StartSpanOption
	if parent != nil {
		opts = []opentracing.StartSpanOption{opentracing.ChildOf(parent.Context())}
	}
	span, ctxNew = opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName, opts...)
	ctxNew = context.WithValue(ctx, BackupActiveSpanKey, span)
	return
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func SpanFromContext(ctx context.Context) (span opentracing.Span) {
	span = opentracing.SpanFromContext(ctx)
	if span == nil {
		val := ctx.Value(BackupActiveSpanKey)
		if sp, ok := val.(opentracing.Span); ok {
			span = sp
		}
	}
	return span
}

// TraceIdFromContext returns the `traceId` previously associated with `ctx`, or
// `""` if not found.
func TraceIdFromContext(ctx context.Context) (traceId string) {
	if ctx == nil {
		return
	}
	return TraceIdFromSpan(SpanFromContext(ctx))
}

// TraceIdFromSpan returns the `traceId` previously associated with `span`, or
// `""` if not found.
func TraceIdFromSpan(span opentracing.Span) (traceId string) {
	if span == nil {
		return
	}
	sc, ok := span.Context().(jaeger.SpanContext)
	if ok {
		traceId = sc.TraceID().String()
	}
	return
}

// OperationNameFromContext returns the `operationName` previously associated with `ctx`, or
// `""` if not found.
func OperationNameFromContext(ctx context.Context) (traceId string) {
	if ctx == nil {
		return
	}
	return OperationNameFromSpan(SpanFromContext(ctx))
}

// OperationNameFromSpan returns the `operationName` previously associated with `span`, or
// `""` if not found.
func OperationNameFromSpan(span opentracing.Span) (operationName string) {
	if span == nil {
		return
	}
	s, ok := span.(*jaeger.Span)
	if ok && s != nil {
		operationName = s.OperationName()
	}
	return
}
