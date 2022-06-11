package es

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testIndexName = "test.item"
)

var initData = []*Item{
	{1, "华为笔记本电脑MateBook D 14 SE版 14英寸 11代酷睿 i5 锐炬显卡 8G+512G 轻薄本/高清护眼防眩光屏 银", nil},

	{2, "惠普(HP)战66五代 锐龙版 14英寸轻薄笔记本电脑(全新2022款锐龙 R5-5625U 16G 512G 高色域低功耗屏 长续航)", nil},

	{3, "戴尔DELL 笔记本电脑 成就3400 14英寸性能商务办公网课手提轻薄本(11代i5-1135G7 16G 512G)一年上门+7x24", nil},

	{4, "联想ThinkPad E14 英特尔酷睿i5 14英寸轻薄笔记本电脑(i5-1135G7 16G 512G 100%sRGB)银",
		[]*InnerRow{
			{11, "内部行", []*InnerInnerRow{
				{111, "内部内部行"}}},
		},
	},

	{5, "Apple MacBook Pro 14英寸 M1 Pro芯片(10核中央处理器 16核图形处理器) 16G 1T深空灰笔记本电脑 MKGQ3CH/A", nil},

	{6, "华硕无双 英特尔Evo平台 12代酷睿i5标压 14.0英寸2.8K 90Hz OLED护眼轻薄笔记本电脑(i5-12500H 16G 512G)银", nil},
}

type Item struct {
	Id   int64       `json:"id" search.type:"terms"`                         //主键id
	Name string      `json:"name" search.type:"terms" search.keyword:"true"` //名称
	Rows []*InnerRow `json:"rows" search.type:"terms"`
}

type InnerRow struct {
	Id   int64            `json:"id" search.type:"terms"`   //主键id
	Name string           `json:"name" search.type:"terms"` //名称
	Rows []*InnerInnerRow `json:"rows" search.type:"terms"`
}
type InnerInnerRow struct {
	Id   int64  `json:"id" search.type:"terms"`   //主键id
	Name string `json:"name" search.type:"terms"` //名称
}

func (m Item) GetID() string {
	return strconv.FormatInt(m.Id, 10)
}

// UNITTEST_ES_ADDR=http://10.x.x.x:9200 go test -v -count 1 .
func TestSearch(t *testing.T) {
	Init("http://10.0.0.2:9200", "", "", true)
	ctx := context.Background()

	vs := make([]Model, 0, len(initData)-1)
	for i := 0; i < len(initData)-1; i++ {
		vs = append(vs, initData[i])
	}
	err := BatchUpsert(ctx, testIndexName, 2, vs)
	if !assert.Equal(t, nil, err) {
		t.FailNow()
	}
	err = Upsert(ctx, testIndexName, initData[len(initData)-1])
	assert.Equal(t, nil, err)

	err = RefreshIndex(ctx, testIndexName)
	assert.Equal(t, nil, err)

	resp, err := SearchMultiMatch(ctx, testIndexName, "电脑", Item{}, 0, 10)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(resp.Hits.Hits))

	resp, err = SearchMultiMatch(ctx, testIndexName, "无双", Item{}, 0, 10)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(resp.Hits.Hits))

	resp, err = SearchBoolMust(ctx, testIndexName, Item{Id: 1, Name: "华硕无双"}, 0, 10)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(resp.Hits.Hits)) //id实际为6

	//包含物料关键字, 且filter id=5,理论上只出现 Apple MacBook Pro
	resp, err = SearchMultiMatchWithFilter(ctx, testIndexName, "电脑", Item{Id: 5}, 0, 10, nil, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(resp.Hits.Hits))

	err = Delete(ctx, testIndexName, Item{Id: 5})
	assert.Equal(t, nil, err)

	err = RefreshIndex(ctx, testIndexName)
	assert.Equal(t, nil, err)

	resp, err = SearchMultiMatch(ctx, testIndexName, "电脑", Item{}, 0, 10)
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(resp.Hits.Hits))

	resp, err = SearchMultiMatch(ctx, testIndexName, "内部内部行", Item{}, 0, 10, &SearchConfig{
		Json:       "rows.rows.name",
		Type:       "terms",
		ForKeyWord: false,
		Boost:      1,
		Val:        "内部内部行",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(resp.Hits.Hits))

	resp, err = SearchMultiMatch(ctx, testIndexName, "内部行", Item{}, 0, 10, &SearchConfig{
		Json:       "rows.name",
		Type:       "terms",
		ForKeyWord: false,
		Boost:      1,
		Val:        "内部行",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(resp.Hits.Hits))

	err = DeleteIndex(ctx, testIndexName)
	assert.Equal(t, nil, err)

}
