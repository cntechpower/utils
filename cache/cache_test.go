package cache

import (
	"context"
	"fmt"
	"testing"

	tracingMySQL "github.com/cntechpower/utils/mysql"
	"github.com/stretchr/testify/assert"

	tracingRedis "github.com/cntechpower/utils/redis"
	"github.com/cntechpower/utils/tracing"
	"github.com/go-redis/redis/v8"
)

var db *tracingMySQL.DB

func TestCache(t *testing.T) {
	tracing.Init("unit-test", "10.0.0.2:6831")
	defer tracing.Close()

	cli := tracingRedis.New(&redis.Options{
		Addr:     "10.0.0.2:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer func() {
		_ = cli.Close()
	}()
	var err error

	db, err = tracingMySQL.New("anywhere:anywhere@tcp(10.0.0.2:3306)/anywhere?charset=utf8")
	if !assert.Equal(t, nil, err) {
		assert.FailNow(t, "connect to db error: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	Init(cli)

	t1 := &TestStruct{Id: 97}
	err = Get(context.Background(), t1)
	assert.Equal(t, nil, err)
	fmt.Println(t1.Name)

	err = Get(context.Background(), t1)
	assert.Equal(t, nil, err)
	fmt.Println(t1.Name)
}

type TestStruct struct {
	Id   int64
	Name string
}

func (t *TestStruct) Key() string {
	return fmt.Sprintf("test-struct-%d", t.Id)
}

func (t *TestStruct) GetFromDB(ctx context.Context) error {
	if t.Id == 0 {
		return fmt.Errorf("id not exists")
	}
	return db.QueryRow(ctx, "select address_en from whitelist_deny_history where id= ? ", t.Id).Scan(&t.Name)
}
