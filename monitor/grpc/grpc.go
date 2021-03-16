package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type InterceptorOption interface {
	apply(*interceptorOptions)
}

type interceptorOptions struct {
	skipMonitorPaths map[string]struct{}
	logEnable        bool
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
		labels := []string{cc.Target(), method}
		start := time.Now()
		o.doStartLog(skip, method, req)
		err = invoker(ctx, method, req, reply, cc, opts...)
		o.doEndLog(skip, method, reply)
		ts := float64(time.Now().Sub(start).Microseconds())
		if !skip {
			clientGrpcDurationTimeHist.Observe(ts)
			clientGrpcDurationTime.WithLabelValues(labels...).Set(ts)
			clientGrpcQueriesTotal.WithLabelValues(labels...).Inc()
			if err != nil {
				clientGrpcErrorsTotal.WithLabelValues(labels...).Inc()
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
		labels := []string{info.FullMethod}
		start := time.Now()
		o.doStartLog(skip, info.FullMethod, req)
		resp, err = handler(ctx, req)
		o.doEndLog(skip, info.FullMethod, resp)
		ts := float64(time.Now().Sub(start).Microseconds())
		if !skip {
			serverGrpcDurationTimeHist.Observe(ts)
			serverGrpcDurationTime.WithLabelValues(labels...).Set(ts)
			serverGrpcQueriesTotal.WithLabelValues(labels...).Inc()
			if err != nil {
				serverGrpcErrorsTotal.WithLabelValues(labels...).Inc()
			}
		}
		return
	}
}
