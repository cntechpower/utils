package log_v2

import (
	"context"
	"os"
	"runtime"

	"github.com/cntechpower/utils/log.v2/output"
	"github.com/cntechpower/utils/tracing"

	log "github.com/sirupsen/logrus"
)

var closer func()

func Init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	SetDefaultFields(hostIpFields)
}

func InitWithES(appId, addr string) {
	log.SetFormatter(&log.JSONFormatter{})
	o, c := output.NewES(appId, addr)
	log.SetOutput(o)
	closer = c
	SetDefaultFields(hostIpFields)
}

func Close() {
	if closer != nil {
		closer()
	}
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		return file, line
	}
	return "unknown.go", 0
}

func getRuntimeFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})
	fields[fieldNameTraceId], fields[fieldNameSpanId] = tracing.TraceSpanIdFromContext(ctx)
	fields[fieldNameTraceName] = tracing.OperationNameFromContext(ctx)
	file, line := getCaller(3)
	fields[fieldNameFileName] = file
	fields[fieldNameFileLine] = line
	return fields
}

func InfoC(ctx context.Context, fields log.Fields, format string, args ...interface{}) {
	fs := getRuntimeFields(ctx)
	defaultLogger.WithFields(fs).WithFields(fields).Infof(format, args...)
}

func WarnC(ctx context.Context, fields log.Fields, format string, args ...interface{}) {
	fs := getRuntimeFields(ctx)
	defaultLogger.WithFields(fs).WithFields(fields).Warnf(format, args...)
}

func ErrorC(ctx context.Context, fields log.Fields, format string, args ...interface{}) {
	fs := getRuntimeFields(ctx)
	defaultLogger.WithFields(fs).WithFields(fields).Errorf(format, args...)
}

func FatalC(ctx context.Context, fields log.Fields, format string, args ...interface{}) {
	fs := getRuntimeFields(ctx)
	defaultLogger.WithFields(fs).WithFields(fields).Fatalf(format, args...)
}

func Infof(fields log.Fields, format string, args ...interface{}) {
	defaultLogger.WithFields(fields).Infof(format, args...)
}

func Warnf(fields log.Fields, format string, args ...interface{}) {
	defaultLogger.WithFields(fields).Warnf(format, args...)
}

func Errorf(fields log.Fields, format string, args ...interface{}) {
	defaultLogger.WithFields(fields).Errorf(format, args...)
}

func Fatalf(fields log.Fields, format string, args ...interface{}) {
	defaultLogger.WithFields(fields).Fatalf(format, args...)
}
