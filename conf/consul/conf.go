package consul

import (
	"encoding/json"
	"path"

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
