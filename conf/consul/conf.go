package consul

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"time"

	"github.com/cntechpower/utils/log"

	"github.com/hashicorp/consul/api"
)

const kvPrefix = "config-center"

var consulClient *api.Client

type IConf interface {
	GetAppName() string
	GetConfKey() string
}

func Init(consulAddr string) {
	var err error
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err = api.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}
}

func Get(c IConf) (err error) {
	kv, _, err := consulClient.KV().Get(path.Join(kvPrefix, c.GetAppName(), c.GetConfKey()), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kv.Value, c)
	return err
}

func Save(c IConf) (err error) {
	content, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	_, err = consulClient.KV().Put(&api.KVPair{
		Key:   path.Join(kvPrefix, c.GetAppName(), c.GetConfKey()),
		Value: content,
	}, nil)
	return
}

var lastIndex uint64

func GetAndWatch(c IConf, interval time.Duration, onChange func(c interface{}) error) (err error) {
	kv, meta, err := consulClient.KV().Get(path.Join(kvPrefix, c.GetAppName(), c.GetConfKey()), nil)
	if err != nil {
		return err
	}
	lastIndex = meta.LastIndex
	h := log.NewHeader(fmt.Sprintf("GetAndWatch-%v-%v", c.GetAppName(), c.GetConfKey()))
	go func() {
		for range time.NewTicker(interval).C {
			kv, meta, err := consulClient.KV().Get(path.Join(kvPrefix, c.GetAppName(), c.GetConfKey()), nil)
			if err != nil {
				h.Errorf("get consul kv error: %v", err)
			}
			if lastIndex == meta.LastIndex {
				continue
			}
			nc := reflect.New(reflect.TypeOf(c)).Interface()
			err = json.Unmarshal(kv.Value, nc)
			if err != nil {
				h.Errorf("json.Unmarshal error: %v", err)
			}
			err = onChange(nc)
			if err != nil {
				h.Errorf("call onChange error: %v", err)
			}
			lastIndex = meta.LastIndex
		}
	}()
	err = json.Unmarshal(kv.Value, c)
	return
}
