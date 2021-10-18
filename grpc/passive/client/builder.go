package client

import (
	"errors"

	"google.golang.org/grpc/resolver"
)

const NAME = "passive"

var (
	errSchemeNotSupport = errors.New("passive resolver: scheme not support")
)

var builder = &passiveBuilder{}

type passiveBuilder struct {
}

// Build creates and starts a DNS resolver that watches the name resolution of the target.
func (b *passiveBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.Scheme != NAME {
		return nil, errSchemeNotSupport
	}
	r := &passiveResolver{
		target: target,
		cc:     cc,
		opts:   opts,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "passive".
func (b *passiveBuilder) Scheme() string {
	return NAME
}

type passiveResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	opts   resolver.BuildOptions
}

// ResolveNow invoke an immediate resolution of the target that this consulResolver watches.
func (r *passiveResolver) ResolveNow(resolver.ResolveNowOptions) {
	r.cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{{
			Addr:       r.target.Endpoint,
			ServerName: "",
			Attributes: nil,
		}},
		ServiceConfig: nil,
		Attributes:    nil,
	})
}

// Close closes the consulResolver.
func (r *passiveResolver) Close() {
}
