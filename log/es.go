package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	appId string
	addr  string
	http  http.Client
}

func newEsWriter(appId, addr string) *esWriter {
	w := &esWriter{
		appId: appId,
		addr:  addr,
	}
	w.http = http.Client{
		Timeout: time.Second,
	}
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
		resp, err = w.http.Post(fmt.Sprintf("http://%v/%v/_doc", w.addr, w.appId), "application/json", strings.NewReader(s))
		if err != nil || resp.StatusCode != http.StatusCreated {
			fmt.Printf("esWriter Println Post fail, err: %v, resp: %v\n", err, resp)
			continue
		}
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			fmt.Printf("esWriter Println Decode fail, err: %v, resp: %v\n", err, resp)
			continue
		}
		break
	}
	if r.Result != "created" {
		fmt.Printf("esWriter Println Post fail, Result: %v\n", r.Result)
		return
	}
}
