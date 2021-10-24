package log_v2

import (
	"os"
	"runtime"

	"github.com/cntechpower/utils/net"
	log "github.com/sirupsen/logrus"
)

const (
	fieldNameHostname    = "os.hostname"
	fieldNameIpAddr      = "os.ip"
	fieldNameRuntimeArch = "runtime.arch"
	fieldNameRuntimeOs   = "runtime.os"
	fieldNameRuntimeGo   = "runtime.go-version"
	fieldNameFileName    = "file_name"
	fieldNameFileLine    = "file_line"
	fieldNameTraceId     = "trace_id"
	fieldNameTraceName   = "trace_name"
	FieldNameBizName     = "biz"
)

var hostName string
var ipAddr string
var defaultFields log.Fields
var hostIpFields log.Fields
var defaultLogger *log.Entry

func SetDefaultFields(fs ...log.Fields) {
	if defaultFields == nil {
		defaultFields = make(map[string]interface{}, 0)
	}
	for _, f := range fs {
		for k, v := range f {
			defaultFields[k] = v
		}
	}
	defaultLogger = log.WithFields(defaultFields)
}

func init() {
	hostName, _ = os.Hostname()
	ipAddr, _ = net.GetFirstLocalIp()
	hostIpFields = log.Fields{
		fieldNameHostname:    hostName,
		fieldNameIpAddr:      ipAddr,
		fieldNameRuntimeArch: runtime.GOARCH,
		fieldNameRuntimeOs:   runtime.GOOS,
		fieldNameRuntimeGo:   runtime.Version(),
	}
	SetDefaultFields(hostIpFields)
}
