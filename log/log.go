package log

import (
	"path"
	"runtime"
)

var loggers []*loggerWithConfig
var options *logOptions

type Logger interface {
	Println(v ...interface{})
}

type outputType string

const (
	OutputTypeText outputType = "TEXT"
	OutputTypeJson outputType = "JSON"
)

type loggerWithConfig struct {
	typ outputType
	Logger
}

func Init(opts ...Option) {
	if loggers != nil {
		panic("Logger already init")
	}
	loggers = make([]*loggerWithConfig, 0)
	options = &logOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
}

type Level string

const (
	levelInfo  Level = "INFO"
	levelWarn  Level = "WARN"
	levelError Level = "ERROR"
	levelFatal Level = "FATAL"
)

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		return path.Base(file), line
	}
	return "unknown.go", 0
}
