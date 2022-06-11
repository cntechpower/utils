package es

import (
	"context"

	"github.com/cntechpower/utils/log"

	"github.com/olivere/elastic/v7"
)

func Upsert(ctx context.Context, index string, value Model) (err error) {
	cli := MustGetCli(ctx)
	resp, err := cli.Index().Index(index).Id(value.GetID()).BodyJson(value).Refresh("true").Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.Upsert Query DB error %v, resp: %v", err, resp)
		return
	}
	return
}

func BatchUpsert(ctx context.Context, index string, batchSize int, values []Model) (err error) {
	cli := MustGetCli(ctx)
	req := make([]elastic.BulkableRequest, 0)
	var resp *elastic.BulkResponse
	for idx, v := range values {
		tv := v
		req = append(req, elastic.NewBulkIndexRequest().Id(tv.GetID()).Doc(tv))
		if (((idx+1)%batchSize == 0) || idx == len(values)-1) && len(req) != 0 {
			resp, err = cli.Bulk().Index(index).Type(DocType).Add(req...).Do(ctx)
			if err != nil || (resp != nil && resp.Errors == true) {
				log.ErrorC(ctx, "olivere.BatchUpsert error %v,resp %v", err, resp)
				return
			}
			req = make([]elastic.BulkableRequest, 0)
		}
	}
	err = RefreshIndex(ctx, index)
	return
}

func Delete(ctx context.Context, index string, m Model) (err error) {
	cli := MustGetCli(ctx)
	resp, err := cli.Delete().Index(index).Id(m.GetID()).Refresh("true").Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.Delete Query error %v, resp: %v", err, resp)
		return
	}
	return
}

func RefreshIndex(ctx context.Context, index ...string) (err error) {
	cli := MustGetCli(ctx)
	resp, err := cli.Refresh(index...).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.RefreshIndex Query error %v, resp: %v", err, resp)
		return
	}
	return
}

func DeleteIndex(ctx context.Context, index string) (err error) {
	cli := MustGetCli(ctx)
	exist, err := cli.IndexExists(index).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.DeleteIndex Query error %v", err)
		return
	}
	if !exist {
		return
	}
	resp, err := cli.DeleteIndex(index).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.DeleteIndex Query error %v, resp: %v", err, resp)
		return
	}
	return
}

func CreateIndex(ctx context.Context, index, mappingString string) (err error) {
	resp, err := MustGetCli(ctx).CreateIndex(index).BodyString(mappingString).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.CreateIndex Query error %v, resp: %v", err, resp)
		return
	}
	return
}
