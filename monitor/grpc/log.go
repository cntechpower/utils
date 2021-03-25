package grpc

import (
	"context"

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

func (o *interceptorOptions) doStartLog(ctx context.Context, skip bool, method string, req interface{}) {
	if o.logStart == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod: method,
		fieldNameRpcReq:    req,
	}).WithReportFileLine(false).Infoc(ctx, "request received")
}

func (o *interceptorOptions) doEndLog(ctx context.Context, skip bool, method string, reply interface{}, ts float64) {
	if o.logEnd == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod:   method,
		fieldNameRpcReply:    reply,
		fieldNameRpcDuration: ts,
	}).WithReportFileLine(false).Infoc(ctx, "request finish")
}
