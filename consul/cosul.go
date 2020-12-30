package consul

import "github.com/hashicorp/consul/api"

var Client *api.Client

func Init(config *api.Config) {
	c, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	api.DefaultConfig()
	Client = c
}
