package consul

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cntechpower/utils/log"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const NAME = "consul"

var (
	errSchemeNotSupport = errors.New("consul resolver: scheme not support")
)

// NewBuilder creates a dnsBuilder which is used to factory DNS resolvers.
func NewBuilder(consulAddr string, refreshInterval time.Duration) resolver.Builder {
	return &consulBuilder{
		consulAddr:      consulAddr,
		refreshInterval: refreshInterval,
	}
}

type consulBuilder struct {
	consulAddr      string
	refreshInterval time.Duration
}

// Build creates and starts a DNS resolver that watches the name resolution of the target.
func (b *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.Scheme != NAME {
		return nil, errSchemeNotSupport
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = b.consulAddr
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	d := &consulResolver{
		h:                    log.NewHeader("consul_resolver"),
		consulClient:         consulClient,
		name:                 target.Endpoint,
		refreshInterval:      b.refreshInterval,
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
		rn:                   make(chan struct{}, 1),
		disableServiceConfig: opts.DisableServiceConfig,
	}

	d.wg.Add(1)
	go d.watcher()
	d.ResolveNow(resolver.ResolveNowOptions{})
	return d, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "consul".
func (b *consulBuilder) Scheme() string {
	return NAME
}

type consulResolver struct {
	h                    *log.Header
	consulClient         *api.Client
	name                 string
	refreshInterval      time.Duration
	wg                   sync.WaitGroup
	ctx                  context.Context
	cancel               context.CancelFunc
	rn                   chan struct{}
	cc                   resolver.ClientConn
	disableServiceConfig bool
	lastIndex            uint64
	lastAddresses        map[string]struct{}
}

// ResolveNow invoke an immediate resolution of the target that this consulResolver watches.
func (r *consulResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

// Close closes the consulResolver.
func (r *consulResolver) Close() {
	r.cancel()
	r.wg.Wait()
}

func getAddressesByMap(m map[string]struct{}) []string {
	addrS := make([]string, len(m))
	for addr := range m {
		addrS = append(addrS, addr)
	}
	return addrS
}

func (r *consulResolver) watcher() {
	defer r.wg.Done()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-r.rn:
		}
		services, meta, err := r.consulClient.Health().Service(r.name, "", true, &api.QueryOptions{WaitIndex: r.lastIndex})
		if err != nil {
			r.h.Errorf("get consul services error: %v, waiting next retry", err)
			continue
		}
		r.lastIndex = meta.LastIndex
		newState := resolver.State{
			Addresses:     make([]resolver.Address, 0),
			ServiceConfig: nil,
			Attributes:    nil,
		}
		addresses := make(map[string]struct{}, len(services))
		changed := false
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			if _, ok := r.lastAddresses[addr]; !ok {
				r.h.Infof("discovering new service: %v address: %v", service.Service.Service, addr)
				changed = true
			}
			addresses[addr] = struct{}{}
			newState.Addresses = append(newState.Addresses, resolver.Address{Addr: addr})
		}
		if changed || len(addresses) != len(r.lastAddresses) {
			r.h.Infof("service changed, old: %v, new: %v", getAddressesByMap(r.lastAddresses), getAddressesByMap(addresses))
		}
		r.lastAddresses = addresses
		r.cc.UpdateState(newState)
		t := time.NewTimer(r.refreshInterval)
		select {
		case <-t.C:
		case <-r.ctx.Done():
			t.Stop()
			return
		}
	}
}
