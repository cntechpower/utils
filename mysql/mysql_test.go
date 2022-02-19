package mysql

import (
	"context"
	"testing"

	"github.com/cntechpower/utils/tracing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLTracing(t *testing.T) {
	tracing.Init("unit-test", "10.0.0.2:6831")
	defer tracing.Close()

	db, err := New("anywhere:anywhere@tcp(10.0.0.2:3306)/anywhere?charset=utf8")
	if !assert.Equal(t, nil, err) {
		assert.FailNow(t, "connect to db error: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()
	_, err = db.Query(context.Background(), "SELECT COUNT(1) FROM rewards_cdkey_v3_1")

	assert.Equal(t, nil, err)

}
