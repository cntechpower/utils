package log

import (
	"encoding/json"
	"os"

	"github.com/cntechpower/utils/net"
)

const (
	fieldNameHostname = "hostname"
	fieldNameIpAddr   = "ip"
	fieldNameTime     = "time"
	fieldNameFileName = "file_name"
	fieldNameFileLine = "file_line"
	fieldNameHeader   = "sub_module"
	fieldNameLevel    = "level"
	fieldNameMessage  = "message"
)

var defaultFields Fields
var hostName string
var ipAddr string

var HostIpFields = Fields{
	fieldNameHostname: hostName,
	fieldNameIpAddr:   ipAddr,
}

func init() {
	hostName, _ = os.Hostname()
	ipAddr, _ = net.GetFirstLocalIp()
}

func SetDefaultFields(fs ...Fields) {
	for _, f := range fs {
		for k, v := range f {
			defaultFields[k] = v
		}
	}
}

type Fields map[string]interface{}

func (f Fields) String() string {
	bs, _ := json.Marshal(f)
	return string(bs)
}

func (f Fields) DeepCopy() Fields {
	nf := make(map[string]interface{}, len(f))
	for k, v := range f {
		nf[k] = v
	}
	return nf
}

type Header struct {
	name   string
	fields Fields
}

func NewHeader(n string) *Header {
	return &Header{
		name:   n,
		fields: map[string]interface{}{},
	}
}

func (h *Header) WithField(key string, value interface{}) *Header {
	return h.WithFields(Fields{key: value})
}

func (h *Header) WithFields(f Fields) *Header {
	for k, v := range f {
		h.fields[k] = v
	}
	return h
}

func (h *Header) String() string {
	return h.name
}

func (h *Header) Info(format string, a ...interface{}) {
	logOutput(3, h, levelInfo, format, a...)
}
func (h *Header) Infof(format string, a ...interface{}) {
	logOutput(3, h, levelInfo, format, a...)
}

func (h *Header) Errorf(format string, a ...interface{}) {
	logOutput(3, h, levelError, format, a...)
}

func (h *Header) Error(err error, format string, a ...interface{}) {
	logOutput(3, h, levelError, "%v", err)
	logOutput(3, h, levelError, format, a...)
}

func (h *Header) Warnf(format string, a ...interface{}) {
	logOutput(3, h, levelWarn, format, a...)
}

func (h *Header) Fatalf(format string, a ...interface{}) {
	logOutput(3, h, levelFatal, format, a...)
	panic(nil)
}

func Infof(h *Header, format string, a ...interface{}) {
	logOutput(3, h, levelInfo, format, a...)
}

func Errorf(h *Header, format string, a ...interface{}) {
	logOutput(3, h, levelError, format, a...)
}

func Warnf(h *Header, format string, a ...interface{}) {
	logOutput(3, h, levelWarn, format, a...)
}

func Fatalf(h *Header, format string, a ...interface{}) {
	logOutput(3, h, levelError, format, a...)
}
