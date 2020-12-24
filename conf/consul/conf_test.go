package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type c struct {
	Id   int64
	Addr string
	IpS  []string
}

func (c *c) GetAppName() string {
	return "test"
}

func (c *c) GetConfKey() string {
	return "test.json"
}

func Test_Conf(t *testing.T) {
	c1 := &c{
		Id:   1,
		Addr: "127.0.0.1:2233",
		IpS:  []string{"a", "b"},
	}

	Init("10.0.0.2:8500")
	assert.Equal(t, nil, Save(c1))
	c2 := &c{}
	assert.Equal(t, nil, Get(c2))
	t.Logf("%+v", c2)

}
