package es

import (
	"context"

	"github.com/cntechpower/utils/log"

	"github.com/olivere/elastic/v7"
)

// SearchBoolMust 多字段AND查询
func SearchBoolMust(ctx context.Context, index string, m interface{}, from, size int, extraConfig ...*SearchConfig) (res *elastic.SearchResult, err error) {
	cli := MustGetCli(ctx)
	service := cli.Search(index)
	res, err = service.Query(QueryBuilderBoolMustSearchQuery(m, extraConfig...)).From(from).Size(size).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.SearchBoolMust Query DB error %v", err)
		return
	}
	return
}

// SearchMultiMatch 多字段模糊查询
func SearchMultiMatch(ctx context.Context, index, keyword string, m interface{}, from, size int,
	extraConfig ...*SearchConfig) (res *elastic.SearchResult, err error) {
	cli := MustGetCli(ctx)
	service := cli.Search(index)
	res, err = service.Query(QueryBuilderMultiMatchSearch(keyword, m, extraConfig...)).
		From(from).Size(size).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.SearchMultiMatch Query DB error %v", err)
		return
	}
	return
}

// SearchMultiMatchWithFilter 多字段模糊查询+指定字段过滤
func SearchMultiMatchWithFilter(ctx context.Context, index, keyword string, m interface{}, from, size int,
	queryExtraConfig []*SearchConfig, filterExtraConfig []*SearchConfig) (res *elastic.SearchResult, err error) {
	cli := MustGetCli(ctx)
	service := cli.Search(index)
	res, err = service.Query(QueryBuilderMultiMatchSearch(keyword, m, queryExtraConfig...)).
		PostFilter(QueryBuilderBoolMustSearchQuery(m, filterExtraConfig...)).
		From(from).Size(size).Do(ctx)
	if err != nil {
		log.ErrorC(ctx, "olivere.SearchMultiMatchWithFilter Query DB error %v", err)
		return
	}
	return
}
