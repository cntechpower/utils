package es

import (
	"context"
	"fmt"
	"sync"

	"github.com/olivere/elastic/v7"
)

var once sync.Once

var cli *elastic.Client

type logger struct {
}

func (logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Init(addr, userName, password string, traceLog bool) {
	once.Do(func() {
		opt := make([]elastic.ClientOptionFunc, 0)
		opt = append(opt, elastic.SetURL(addr))
		if userName != "" && password != "" {
			opt = append(opt, elastic.SetBasicAuth(userName, password))
		}
		if traceLog {
			opt = append(opt, elastic.SetTraceLog(logger{}))
		}
		var err error
		cli, err = elastic.NewClient(opt...)
		if err != nil {
			panic(err)
		}
	})
}

func MustGetCli(ctx context.Context) *elastic.Client {
	if cli == nil {
		panic("cli is nil")
	}
	return cli
}
