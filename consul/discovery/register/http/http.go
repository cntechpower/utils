package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cntechpower/utils/consul"
	"github.com/hashicorp/consul/api"
)

const (
	healthCheckInterval = 5
	healthCheckTimeout  = 2
	deregisterTimeout   = 10
)

func serviceName(appName string) string {
	return fmt.Sprintf("http-%s", appName)
}

func Register(appName, address string, port int) error {
	regHttp := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v:%v", serviceName(appName), address, strconv.Itoa(port)),
		Name:    serviceName(appName),
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("health-http-%v-%v:%v", appName, address, strconv.Itoa(port)),
			Name:                           serviceName(appName),
			Interval:                       (healthCheckInterval * time.Second).String(),
			Timeout:                        (healthCheckTimeout * time.Second).String(),
			HTTP:                           fmt.Sprintf("http://%v:%v/ping", address, port),
			DeregisterCriticalServiceAfter: (deregisterTimeout * time.Second).String(),
		},
	}
	return consul.Client.Agent().ServiceRegister(regHttp)
}
