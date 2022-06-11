package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/cntechpower/utils/net"
)

const (
	fieldNameHostname    = "os.hostname"
	fieldNameIpAddr      = "os.ip"
	fieldNameRuntimeArch = "runtime.arch"
	fieldNameRuntimeOs   = "runtime.os"
	fieldNameRuntimeGo   = "runtime.go-version"
	fieldNameTime        = "time"
	fieldNameFileName    = "file_name"
	fieldNameFileLine    = "file_line"
	fieldNameHeader      = "module"
	fieldNameLevel       = "level"
	fieldNameMessage     = "message"
	fieldNameTracing     = "trace_id"
)

var defaultFields Fields
var hostName string
var ipAddr string

var HostIpFields Fields

func init() {
	hostName, _ = os.Hostname()
	ipAddr, _ = net.GetFirstLocalIp()
	HostIpFields = Fields{
		fieldNameHostname:    hostName,
		fieldNameIpAddr:      ipAddr,
		fieldNameRuntimeArch: runtime.GOARCH,
		fieldNameRuntimeOs:   runtime.GOOS,
		fieldNameRuntimeGo:   runtime.Version(),
	}
	SetDefaultFields(HostIpFields)
}

func SetDefaultFields(fs ...Fields) {
	if defaultFields == nil {
		defaultFields = make(map[string]interface{}, 0)
	}
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
	name           string
	fields         Fields
	skipCallers    int
	reportFileLine bool
}

func NewHeader(n string) *Header {
	return &Header{
		name:           n,
		fields:         map[string]interface{}{},
		skipCallers:    3,
		reportFileLine: true,
	}
}

func (h *Header) WithSkipCallers(n int) *Header {
	h.skipCallers = n
	return h
}

func (h *Header) WithReportFileLine(b bool) *Header {
	h.reportFileLine = b
	return h
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

func (h *Header) Printf(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelInfo, format, a...)
}

func (h *Header) Info(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelInfo, format, a...)
}
func (h *Header) Infof(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelInfo, format, a...)
}

func (h *Header) Errorf(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelError, format, a...)
}

func (h *Header) Error(err error, format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelError, "%v", err)
	logOutput(context.Background(), h.skipCallers, h, levelError, format, a...)
}

func (h *Header) Warnf(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelWarn, format, a...)
}

func (h *Header) Fatalf(format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelFatal, format, a...)
	panic(fmt.Sprintf(format, a...))
}

func (h *Header) Infoc(ctx context.Context, format string, a ...interface{}) {
	logOutput(ctx, h.skipCallers, h, levelInfo, format, a...)
}

func (h *Header) Errorc(ctx context.Context, format string, a ...interface{}) {
	logOutput(ctx, h.skipCallers, h, levelError, format, a...)
}

func (h *Header) Warnc(ctx context.Context, format string, a ...interface{}) {
	logOutput(ctx, h.skipCallers, h, levelWarn, format, a...)
}

func (h *Header) Fatalc(ctx context.Context, format string, a ...interface{}) {
	logOutput(ctx, h.skipCallers, h, levelFatal, format, a...)
	panic(nil)
}

func Infof(h *Header, format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelInfo, format, a...)
}

func Errorf(h *Header, format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelError, format, a...)
}

func Warnf(h *Header, format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelWarn, format, a...)
}

func Fatalf(h *Header, format string, a ...interface{}) {
	logOutput(context.Background(), h.skipCallers, h, levelError, format, a...)
}

type headerKey struct {
}

func Inject(ctx context.Context, name string) {
	ctx = context.WithValue(ctx, headerKey{}, NewHeader(name))
}

func Extract(ctx context.Context) (h *Header) {
	h = ctx.Value(headerKey{}).(*Header)
	return
}

func TryExtract(ctx context.Context) (h *Header) {
	h = ctx.Value(headerKey{}).(*Header)
	if h == nil {
		h = NewHeader("default-header")
	}
	return
}

func FatalC(ctx context.Context, format string, a ...interface{}) {
	h := TryExtract(ctx)
	logOutput(context.Background(), h.skipCallers, h, levelFatal, format, a...)
	panic(fmt.Sprintf(format, a...))
}

func InfoC(ctx context.Context, format string, a ...interface{}) {
	h := TryExtract(ctx)
	logOutput(ctx, h.skipCallers, h, levelInfo, format, a...)
}

func ErrorC(ctx context.Context, format string, a ...interface{}) {
	h := TryExtract(ctx)
	logOutput(ctx, h.skipCallers, h, levelError, format, a...)
}

func WarnC(ctx context.Context, format string, a ...interface{}) {
	h := TryExtract(ctx)
	logOutput(ctx, h.skipCallers, h, levelWarn, format, a...)
}
