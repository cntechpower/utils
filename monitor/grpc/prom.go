package grpc

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

//client side metrics
var (
	clientGrpcQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "client_grpc_queries_count",
			Help:        "client_grpc_queries_count",
			ConstLabels: nil,
		}, []string{"target", "method"})
	clientGrpcErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "client_grpc_errors_count",
			Help:        "client_grpc_errors_count",
			ConstLabels: nil,
		}, []string{"target", "method"})
	clientGrpcDurationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "client_grpc_time_duration",
			Help:        "",
			ConstLabels: nil,
		}, []string{"target", "method"})
	clientGrpcDurationTimeHist = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "client_grpc_request_duration_us",
		Help:        "",
		ConstLabels: nil,
		Buckets:     []float64{10, 20, 50, 100, 1000},
	})
)

//server side metrics
var (
	serverGrpcQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "server_grpc_queries_count",
			Help:        "server_grpc_queries_count",
			ConstLabels: nil,
		}, []string{"method"})
	serverGrpcErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "server_grpc_errors_count",
			Help:        "server_grpc_errors_count",
			ConstLabels: nil,
		}, []string{"method"})
	serverGrpcDurationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "server_grpc_time_duration",
			Help:        "",
			ConstLabels: nil,
		}, []string{"method"})
	serverGrpcDurationTimeHist = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "server_grpc_request_duration_us",
		Help:        "",
		ConstLabels: nil,
		Buckets:     []float64{10, 20, 50, 100, 1000},
	})
)

func init() {
	prometheus.MustRegister(
		//client side
		clientGrpcQueriesTotal,
		clientGrpcErrorsTotal,
		clientGrpcDurationTime,
		clientGrpcDurationTimeHist,
		//server side
		serverGrpcQueriesTotal,
		serverGrpcErrorsTotal,
		serverGrpcDurationTime,
		serverGrpcDurationTimeHist,
	)
}

func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		labels := []string{cc.Target(), method}
		start := time.Now()
		err = invoker(ctx, method, req, reply, cc, opts...)
		ts := float64(time.Now().Sub(start).Microseconds())
		clientGrpcDurationTimeHist.Observe(ts)
		clientGrpcDurationTime.WithLabelValues(labels...).Set(ts)
		clientGrpcQueriesTotal.WithLabelValues(labels...).Inc()
		if err != nil {
			clientGrpcErrorsTotal.WithLabelValues(labels...).Inc()
		}
		return err
	}
}

func GetUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		labels := []string{info.FullMethod}
		start := time.Now()
		resp, err = handler(ctx, req)
		ts := float64(time.Now().Sub(start).Microseconds())
		serverGrpcDurationTimeHist.Observe(ts)
		serverGrpcDurationTime.WithLabelValues(labels...).Set(ts)
		serverGrpcQueriesTotal.WithLabelValues(labels...).Inc()
		if err != nil {
			serverGrpcErrorsTotal.WithLabelValues(labels...).Inc()
		}
		return
	}
}
