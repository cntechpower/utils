package es

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/olivere/elastic/v7"
)

const (
	TypeTerm  = "terms"
	TypeGt    = "gt"
	TypeGet   = "gte"
	TypeLt    = "lt"
	TypeLte   = "lte"
	TypeMatch = "match"

	//标志该字段是否参与关键字MultiMatchSearch
	searchTypeKeyword = "keyword"
)

type SearchConfig struct {
	Json       string      `json:"json"`
	Type       string      `json:"type"`
	ForKeyWord bool        `json:"for_key_word"`
	Boost      float64     `json:"boost"`
	Val        interface{} `json:"val"`
}

// parseSearchTag 解析struct tag获取搜索配置
func parseSearchTag(field reflect.StructField) (c *SearchConfig, ok bool) {
	c = &SearchConfig{}
	c.Json = field.Tag.Get("json")
	if c.Json == "" {
		// 并非json字段
		return
	}
	ok = true
	c.Type = field.Tag.Get("search.type")
	if c.Type == "" {
		c.Type = TypeTerm //默认走terms查询
	}
	c.ForKeyWord = strings.ToLower(field.Tag.Get("search.keyword")) == "true"
	c.Boost = float64(1)
	b, err1 := strconv.ParseFloat(field.Tag.Get("search.boost"), 64)
	if err1 == nil {
		c.Boost = b
	}
	return
}

func addMustToBoolQuery(c *SearchConfig, query *elastic.BoolQuery) {
	switch strings.ToLower(c.Type) {
	case TypeTerm:
		query.Must(elastic.NewTermsQuery(c.Json, c.Val).Boost(c.Boost))
	case TypeGt:
		query.Must(elastic.NewRangeQuery(c.Json).Gt(c.Val).Boost(c.Boost))
	case TypeGet:
		query.Must(elastic.NewRangeQuery(c.Json).Gte(c.Val).Boost(c.Boost))
	case TypeLt:
		query.Must(elastic.NewRangeQuery(c.Json).Lt(c.Val).Boost(c.Boost))
	case TypeLte:
		query.Must(elastic.NewRangeQuery(c.Json).Lte(c.Val).Boost(c.Boost))
	case TypeMatch:
		query.Must(elastic.NewMatchQuery(c.Json, c.Val).Boost(c.Boost))
	default:
	}
}

// QueryBuilderBoolMustSearchQuery 多字段AND查询构建
func QueryBuilderBoolMustSearchQuery(m interface{}, configs ...*SearchConfig) elastic.Query {
	query := elastic.NewBoolQuery()
	t := reflect.TypeOf(m)
	val := reflect.ValueOf(m) //获取reflect.Type类型
	for i := 0; i < t.NumField(); i++ {
		isZero := val.Field(i).IsZero()
		if isZero {
			continue
		}
		c, ok := parseSearchTag(t.Field(i))
		if !ok {
			continue
		}
		c.Val = val.Field(i).Interface()
		addMustToBoolQuery(c, query)
	}

	for _, c := range configs {
		addMustToBoolQuery(c, query)
	}
	return query
}

// QueryBuilderMultiMatchSearch 多字段模糊查询构建
func QueryBuilderMultiMatchSearch(keyword string, m interface{}, extraConfig ...*SearchConfig) (query elastic.Query) {
	if keyword == "" {
		query = elastic.NewMatchAllQuery()
		return
	}
	fields := make([]string, 0)
	t := reflect.TypeOf(m)
	for i := 0; i < t.NumField(); i++ {
		c, ok := parseSearchTag(t.Field(i))
		if !ok || !c.ForKeyWord {
			continue
		}

		fields = append(fields, fmt.Sprintf("%s^%v", c.Json, c.Boost)) //support boost
	}
	for _, c := range extraConfig {
		fields = append(fields, fmt.Sprintf("%s^%v", c.Json, c.Boost)) //support boost
	}
	query = elastic.NewMultiMatchQuery(keyword, fields...)
	return query
}
