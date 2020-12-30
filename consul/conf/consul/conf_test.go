package consul

import (
	"testing"
	"time"

	"github.com/cntechpower/utils/log"

	"github.com/stretchr/testify/assert"
)

type testConf struct {
	Id   int64
	Addr string
	IpS  []string
}

func (c *testConf) GetAppName() string {
	return "test"
}

func (c *testConf) GetConfKey() string {
	return "test.json"
}

func Test_Conf(t *testing.T) {
	c1 := &testConf{
		Id:   1,
		Addr: "127.0.0.1:2233",
		IpS:  []string{"a", "b"},
	}

	log.InitLogger("")
	Init("10.0.0.2:8500")
	assert.Equal(t, nil, Save(c1))
	c2 := &testConf{}
	assert.Equal(t, nil, Get(c2))
	t.Logf("%+v", c2)
	c3 := &testConf{}
	err := GetAndWatch(c3, time.Second*3, func(c interface{}) error {
		tc := c.(**testConf)
		t.Logf("Changed: pointer %v, values: %+v\n", &c, **tc)
		return nil
	})
	assert.Equal(t, nil, err)
	t.Logf("c3: %v ,%v", &c3, c3)

	time.Sleep(30 * time.Second)
	t.Logf("c3: %v ,%v", &c3, c3)

}
