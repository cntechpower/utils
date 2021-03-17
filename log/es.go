package log

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

/*
{
    "_index": "main.anywhere",
    "_type": "_doc",
    "_id": "uSXHO3gBDZnrgKvXqiV3",
    "_version": 1,
    "result": "created",
    "_shards": {
        "total": 2,
        "successful": 1,
        "failed": 0
    },
    "_seq_no": 1,
    "_primary_term": 1
}
*/

type esResp struct {
	Result string `json:"result"`
}

type esWriter struct {
	appId    string
	addr     string
	http     http.Client
	esClient *elasticsearch.Client
}

func newEsWriter(appId, addr string) *esWriter {
	w := &esWriter{
		appId: appId,
		addr:  addr,
	}
	w.http = http.Client{
		Timeout: time.Second,
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
	var resp *http.Response
	var r *esResp
	for i := 0; i < 3; i++ {
		resp, err := w.esClient.Index(w.appId, strings.NewReader(s), nil)
		if err == nil && !resp.IsError() {
			break
		}
		//resp, err = w.http.Post(fmt.Sprintf("http://%v/%v/_doc", w.addr, w.appId), "application/json", strings.NewReader(s))
		//if err == nil && resp.StatusCode == http.StatusCreated {
		//	err = json.NewDecoder(resp.Body).Decode(&r)
		//	if err == nil && r.Result == "created" {
		//		break
		//	}
		//}
	}
	if err != nil {
		fmt.Printf("esWriter Println Post fail, err: %v\n", err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("esWriter Println Post got code: %v\n", resp.StatusCode)
		return
	}
	if r.Result != "created" {
		fmt.Printf("esWriter Println Post fail, Result: %v\n", r.Result)
		return
	}
}
