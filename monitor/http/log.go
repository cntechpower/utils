package http

import (
	"strings"

	"github.com/cntechpower/utils/log"
	"github.com/gin-gonic/gin"
)

const (
	fieldNameHTTPMethod   = "http.method"
	fieldNameHTTPPath     = "http.path"
	fieldNameHTTPParams   = "http.params"
	fieldNameRequestHost  = "http.host"
	fieldNameResponseCode = "http.code"
	fieldNameHTTPDuration = "http.duration"
)

func WithLog(logStart, logEnd bool) GinMiddlewareOption {
	return newFuncServerOption(func(options *ginMiddlewareOptions) {
		options.logStart = logStart
		options.logEnd = logEnd
	})
}

func (o *ginMiddlewareOptions) doStartLog(skip bool, ctx *gin.Context) {
	if !o.logStart || skip {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHTTPMethod:  ctx.Request.Method,
		fieldNameHTTPPath:    ctx.Request.URL.Path,
		fieldNameHTTPParams:  strings.TrimLeft(ctx.Request.RequestURI, ctx.Request.URL.Path+"?"),
		fieldNameRequestHost: strings.Split(ctx.Request.RemoteAddr, ":")[0],
	}).WithReportFileLine(false).Infoc(ctx, "request received")
}

func (o *ginMiddlewareOptions) doEndLog(skip bool, ctx *gin.Context, dur float64) {
	if !o.logEnd || skip {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHTTPMethod:   ctx.Request.Method,
		fieldNameHTTPPath:     ctx.Request.URL.Path,
		fieldNameRequestHost:  strings.Split(ctx.Request.RemoteAddr, ":")[0],
		fieldNameResponseCode: ctx.Writer.Status(),
		fieldNameHTTPDuration: dur,
	}).WithReportFileLine(false).Infoc(ctx, "request finish")
}
