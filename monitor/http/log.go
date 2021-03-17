package http

import (
	"strings"

	"github.com/cntechpower/utils/log"
	"github.com/gin-gonic/gin"
)

const (
	fieldNameHttpMethod   = "http.method"
	fieldNameHttpPath     = "http.path"
	fieldNameRequestHost  = "http.host"
	fieldNameResponseCode = "http.code"
	fieldNameHttpDuration = "http.duration"
)

func WithLog(logStart, logEnd bool) GinMiddlewareOption {
	return newFuncServerOption(func(options *ginMiddlewareOptions) {
		options.logStart = logStart
		options.logEnd = logEnd
	})
}

func (o *ginMiddlewareOptions) doStartLog(skip bool, ctx *gin.Context) {
	if o.logStart == false || skip {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHttpMethod:  ctx.Request.Method,
		fieldNameHttpPath:    ctx.Request.RequestURI,
		fieldNameRequestHost: strings.Split(ctx.Request.RemoteAddr, ":")[0],
	}).WithReportFileLine(false).Infof("request received")
}

func (o *ginMiddlewareOptions) doEndLog(skip bool, ctx *gin.Context, dur float64) {
	if o.logEnd == false || skip {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHttpMethod:   ctx.Request.Method,
		fieldNameHttpPath:     ctx.Request.RequestURI,
		fieldNameRequestHost:  strings.Split(ctx.Request.RemoteAddr, ":")[0],
		fieldNameResponseCode: ctx.Writer.Status(),
		fieldNameHttpDuration: dur,
	}).WithReportFileLine(false).Infof("request finish")
}
