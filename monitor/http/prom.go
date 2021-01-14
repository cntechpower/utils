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

type GinMiddlewareOption interface {
	apply(*ginMiddlewareOptions)
}

type ginMiddlewareOptions struct {
	skipMonitorPaths map[string]struct{}
}

// funcMiddlewareOption wraps a function that modifies ginMiddlewareOptions into an
// implementation of the GinMiddlewareOption interface.
type funcMiddlewareOption struct {
	f func(*ginMiddlewareOptions)
}

func (fdo *funcMiddlewareOption) apply(do *ginMiddlewareOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*ginMiddlewareOptions)) *funcMiddlewareOption {
	return &funcMiddlewareOption{
		f: f,
	}
}

func WithBlackList(l []string) GinMiddlewareOption {
	return newFuncServerOption(func(options *ginMiddlewareOptions) {
		for _, p := range l {
			options.skipMonitorPaths[p] = struct{}{}
		}
	})
}

func GinMiddleware(opts ...GinMiddlewareOption) gin.HandlerFunc {
	o := &ginMiddlewareOptions{skipMonitorPaths: map[string]struct{}{}}
	for _, f := range opts {
		f.apply(o)
	}
	return func(ctx *gin.Context) {
		if _, ok := o.skipMonitorPaths[ctx.Request.URL.Path]; ok {
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
