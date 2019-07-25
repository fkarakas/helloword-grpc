package resolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"strconv"
	"time"
)

func NewBuilderWithScheme(scheme string) *Resolver {
	return &Resolver{
		scheme: scheme,
	}
}

type Resolver struct {
	scheme string
	cc     resolver.ClientConn
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.cc = cc
	if r.bootstrapState != nil {
		r.UpdateState(*r.bootstrapState)
	}
	return r, nil
}

// Scheme returns the test scheme.
func (r *Resolver) Scheme() string {
	return r.scheme
}

func (*Resolver) ResolveNow(o resolver.ResolveNowOption) {
	fmt.Println("@@@@@@@@ CALLING ResolveNow with", o)
}

// Close is a noop for Resolver.
func (*Resolver) Close() {
	fmt.Println("@@@@@@@@ CALLING Close")
}

// UpdateState calls cc.UpdateState.
func (r *Resolver) UpdateState(s resolver.State) {
	r.cc.UpdateState(s)
}

// GenerateAndRegisterManualResolver generates a random scheme and a Resolver
// with it. It also registers this Resolver.
// It returns the Resolver and a cleanup function to unregister it.
func GenerateAndRegisterManualResolver() (*Resolver, func()) {
	scheme := strconv.FormatInt(time.Now().UnixNano(), 36)
	r := NewBuilderWithScheme(scheme)
	resolver.Register(r)
	return r, func() { resolver.UnregisterForTesting(scheme) }
}
