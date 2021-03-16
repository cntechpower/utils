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
)

func WithLog() GinMiddlewareOption {
	return newFuncServerOption(func(options *ginMiddlewareOptions) {
		options.logEnable = true
	})
}

func (o *ginMiddlewareOptions) doStartLog(ctx *gin.Context) {
	if o.logEnable == false {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHttpMethod:  ctx.Request.Method,
		fieldNameHttpPath:    ctx.Request.RequestURI,
		fieldNameRequestHost: strings.Split(ctx.Request.RemoteAddr, ":")[0],
	}).Infof("request received")
}

func (o *ginMiddlewareOptions) doEndLog(ctx *gin.Context) {
	if o.logEnable == false {
		return
	}
	log.NewHeader("http-access").WithFields(log.Fields{
		fieldNameHttpMethod:   ctx.Request.Method,
		fieldNameHttpPath:     ctx.Request.RequestURI,
		fieldNameRequestHost:  strings.Split(ctx.Request.RemoteAddr, ":")[0],
		fieldNameResponseCode: ctx.Writer.Status(),
	}).Infof("request end")
}
