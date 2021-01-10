package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "http_queries_count",
			Help:        "http_queries_count",
			ConstLabels: nil,
		}, []string{"path", "code"})
	httpDurationTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "http_time_duration",
			Help:        "",
			ConstLabels: nil,
		}, []string{"path", "code"})
	httpDurationTimeHist = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        "http_request_duration_us",
		Help:        "",
		ConstLabels: nil,
		Buckets:     []float64{10, 20, 50, 100, 1000},
	})
)

func init() {
	prometheus.MustRegister(
		httpQueriesTotal,
		httpDurationTime,
		httpDurationTimeHist,
	)
}

func GinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}
		labels := []string{
			ctx.Request.RequestURI,
			strconv.Itoa(ctx.Writer.Status())}

		//doing request
		start := time.Now()
		ctx.Next()
		ts := float64(time.Now().Sub(start).Microseconds())
		httpDurationTime.WithLabelValues(labels...).Set(ts)
		httpDurationTimeHist.Observe(ts)
		httpQueriesTotal.WithLabelValues(labels...).Inc()
	}
}
