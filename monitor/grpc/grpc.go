package grpc

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type InterceptorOption interface {
	apply(*interceptorOptions)
}

type interceptorOptions struct {
	skipMonitorPaths map[string]struct{}
	logStart         bool
	logEnd           bool
	traceEnable      bool
}

// funcInterceptorOption wraps a function that modifies interceptorOptions into an
// implementation of the InterceptorOption interface.
type funcInterceptorOption struct {
	f func(*interceptorOptions)
}

func (fdo *funcInterceptorOption) apply(do *interceptorOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*interceptorOptions)) *funcInterceptorOption {
	return &funcInterceptorOption{
		f: f,
	}
}

func WithBlackList(l []string) InterceptorOption {
	return newFuncServerOption(func(options *interceptorOptions) {
		for _, p := range l {
			options.skipMonitorPaths[p] = struct{}{}
		}
	})
}

func GetUnaryClientInterceptor(opts ...InterceptorOption) grpc.UnaryClientInterceptor {
	o := &interceptorOptions{skipMonitorPaths: map[string]struct{}{}}
	for _, f := range opts {
		f.apply(o)
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		_, skip := o.skipMonitorPaths[method]
		var span opentracing.Span
		if !skip && o.traceEnable {
			ctx, span = clientTryInject(method, ctx)
		}
		labels := []string{cc.Target(), method}
		start := time.Now()
		o.doStartLog(ctx, skip, method, req)
		err = invoker(ctx, method, req, reply, cc, opts...)
		ts := float64(time.Now().Sub(start).Microseconds())
		o.doEndLog(ctx, skip, method, reply, ts)
		if !skip {
			clientGrpcDurationTimeHist.Observe(ts)
			clientGrpcDurationTime.WithLabelValues(labels...).Set(ts)
			clientGrpcQueriesTotal.WithLabelValues(labels...).Inc()
			if err != nil {
				clientGrpcErrorsTotal.WithLabelValues(labels...).Inc()
				if span != nil {
					SetSpanTags(span, err, true)
				}
			}
			if span != nil {
				span.Finish()
			}
		}
		return err
	}
}

func GetUnaryServerInterceptor(opts ...InterceptorOption) grpc.UnaryServerInterceptor {
	o := &interceptorOptions{skipMonitorPaths: map[string]struct{}{}}
	for _, f := range opts {
		f.apply(o)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		_, skip := o.skipMonitorPaths[info.FullMethod]
		var span opentracing.Span
		if !skip && o.traceEnable {
			ctx, span = serverTryExtract(info.FullMethod, ctx)
		}
		labels := []string{info.FullMethod}
		start := time.Now()
		o.doStartLog(ctx, skip, info.FullMethod, req)
		resp, err = handler(ctx, req)
		ts := float64(time.Now().Sub(start).Microseconds())
		o.doEndLog(ctx, skip, info.FullMethod, resp, ts)
		if !skip {
			serverGrpcDurationTimeHist.Observe(ts)
			serverGrpcDurationTime.WithLabelValues(labels...).Set(ts)
			serverGrpcQueriesTotal.WithLabelValues(labels...).Inc()
			if err != nil {
				serverGrpcErrorsTotal.WithLabelValues(labels...).Inc()
				if span != nil {
					SetSpanTags(span, err, false)
				}
			}
			if span != nil {
				span.Finish()
			}
		}
		return
	}
}
