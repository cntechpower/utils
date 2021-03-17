package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type esWriter struct {
	appId    string
	addr     string
	esClient *elasticsearch.Client
}

func newEsWriter(appId, addr string) *esWriter {
	w := &esWriter{
		appId: appId,
		addr:  addr,
	}
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:  []string{addr},
		MaxRetries: 3,
	})
	if err != nil {
		panic(err)
	}
	w.esClient = c

	return w
}
func (w *esWriter) Println(v ...interface{}) {
	if len(v) != 1 {
		fmt.Println("esWriter Println got multi v")
		return
	}
	s, ok := v[0].(string)
	if !ok {
		fmt.Println("esWriter Println got non string")
		return
	}
	var err error
	var resp *esapi.Response
	for i := 0; i < 3; i++ {
		resp, err = w.esClient.Index(w.appId, strings.NewReader(s),
			w.esClient.Index.WithTimeout(time.Second))
		if err == nil && !resp.IsError() {
			break
		}
	}
}
