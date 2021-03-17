package grpc

import (
	"github.com/cntechpower/utils/log"
)

const (
	fieldNameRpcMethod   = "grpc.method"
	fieldNameRpcReq      = "grpc.req"
	fieldNameRpcReply    = "grpc.reply"
	fieldNameRpcDuration = "grpc.duration"
)

func WithLog(logStart, logEnd bool) InterceptorOption {
	return newFuncServerOption(func(options *interceptorOptions) {
		options.logStart = logStart
		options.logEnd = logEnd
	})
}

func (o *interceptorOptions) doStartLog(skip bool, method string, req interface{}) {
	if o.logStart == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod: method,
		fieldNameRpcReq:    req,
	}).WithReportFileLine(false).Infof("request received")
}

func (o *interceptorOptions) doEndLog(skip bool, method string, reply interface{}, ts float64) {
	if o.logEnd == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod:   method,
		fieldNameRpcReply:    reply,
		fieldNameRpcDuration: ts,
	}).WithReportFileLine(false).Infof("request finish")
}
