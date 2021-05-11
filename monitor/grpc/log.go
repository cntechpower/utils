package grpc

import (
	"context"
	"encoding/json"

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
	bs, _ := json.Marshal(req)
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod: method,
		fieldNameRpcReq:    string(bs),
	}).WithReportFileLine(false).Infoc(ctx, "request received")
}

func (o *interceptorOptions) doEndLog(ctx context.Context, skip bool, method string, reply interface{}, ts float64) {
	if o.logEnd == false || skip {
		return
	}
	bs, _ := json.Marshal(reply)
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod:   method,
		fieldNameRpcReply:    string(bs),
		fieldNameRpcDuration: ts,
	}).WithReportFileLine(false).Infoc(ctx, "request finish")
}
