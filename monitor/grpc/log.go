package grpc

import (
	"github.com/cntechpower/utils/log"
)

const (
	fieldNameRpcMethod = "grpc.method"
	fieldNameRpcReq    = "grpc.req"
	fieldNameRpcReply  = "grpc.reply"
)

func WithLog() InterceptorOption {
	return newFuncServerOption(func(options *interceptorOptions) {
		options.logEnable = true
	})
}

func (o *interceptorOptions) doStartLog(skip bool, method string, req interface{}) {
	if o.logEnable == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod: method,
		fieldNameRpcReq:    req,
	}).Infof("request received")
}

func (o *interceptorOptions) doEndLog(skip bool, method string, reply interface{}) {
	if o.logEnable == false || skip {
		return
	}
	log.NewHeader("grpc-access").WithFields(log.Fields{
		fieldNameRpcMethod: method,
		fieldNameRpcReply:  reply,
	}).Infof("request end")
}
