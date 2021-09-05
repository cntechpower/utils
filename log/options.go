package log

import (
	_log "log"
	"os"
)

type logOptions struct {
}

type Option interface {
	apply(option *logOptions)
}

type funcLogOptions struct {
	f func(option *logOptions)
}

func (fdo *funcLogOptions) apply(option *logOptions) {
	fdo.f(option)
}

func newLogOption(f func(*logOptions)) *funcLogOptions {
	return &funcLogOptions{
		f: f,
	}
}

func WithStd(typ outputType) Option {
	return newLogOption(func(option *logOptions) {
		l := &_log.Logger{}
		l.SetOutput(os.Stdout)
		lc := &loggerWithConfig{
			typ:    typ,
			buffer: make(chan string, 1000),
			Logger: l,
		}
		loggers = append(loggers, lc)
		go lc.run()
	})
}

func WithFile(typ outputType, fileName string) Option {
	return newLogOption(func(option *logOptions) {
		l := &_log.Logger{}
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		l.SetOutput(file)
		lc := &loggerWithConfig{
			typ:    typ,
			buffer: make(chan string, 1000),
			Logger: l,
		}
		loggers = append(loggers, lc)
		go lc.run()
	})
}

func WithEs(appID, esAddr string) Option {
	return newLogOption(func(option *logOptions) {
		l := newEsWriter(appID, esAddr)
		lc := &loggerWithConfig{
			typ:    OutputTypeJson,
			buffer: make(chan string, 1000),
			Logger: l,
		}
		loggers = append(loggers, lc)
		go lc.run()
	})
}
