package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GinMiddlewareOption interface {
	apply(*ginMiddlewareOptions)
}

type ginMiddlewareOptions struct {
	skipMonitorPaths map[string]struct{}
	logStart         bool
	logEnd           bool
	traceEnable      bool
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
	o := &ginMiddlewareOptions{
		skipMonitorPaths: map[string]struct{}{},
		logStart:         false,
		logEnd:           false,
	}
	for _, f := range opts {
		f.apply(o)
	}
	return func(ctx *gin.Context) {
		_, skip := o.skipMonitorPaths[ctx.Request.URL.Path]
		if skip {
			ctx.Next()
			return
		}
		if o.traceEnable {
			inject(ctx)
		}
		labels := []string{
			ctx.Request.RequestURI,
			strconv.Itoa(ctx.Writer.Status())}

		//doing request
		start := time.Now()
		o.doStartLog(skip, ctx)
		ctx.Next()
		ts := float64(time.Now().Sub(start).Microseconds())
		o.doEndLog(skip, ctx, ts)
		httpDurationTime.WithLabelValues(labels...).Set(ts)
		httpDurationTimeHist.Observe(ts)
		httpQueriesTotal.WithLabelValues(labels...).Inc()
	}
}
