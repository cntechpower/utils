package grpc

import "github.com/prometheus/client_golang/prometheus"

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
