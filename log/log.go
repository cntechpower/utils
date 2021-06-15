package log

import (
	"path"
	"runtime"
	"sync"
)

var loggers []*loggerWithConfig
var options *logOptions
var wg sync.WaitGroup
var closing bool

type Logger interface {
	Println(v ...interface{})
}

type outputType string

const (
	OutputTypeText outputType = "TEXT"
	OutputTypeJson outputType = "JSON"
)

type loggerWithConfig struct {
	typ    outputType
	buffer chan string
	Logger
}

func (l *loggerWithConfig) run() {
	defer wg.Done()
	for {
		select {
		case s, ok := <-l.buffer:
			if !ok {
				return
			}
			l.Println(s)
		}
	}
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
	wg.Add(len(loggers))
}

func Close() {
	closing = true
	for _, l := range loggers {
		close(l.buffer)
	}
	wg.Wait()
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
