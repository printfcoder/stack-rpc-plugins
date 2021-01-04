package plugins

import (
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/registry"
)

type Option func(o *Options)

type Options struct {
	// todo move all of below to runtime
	Registry registry.Registry
	Selector selector.Selector
	Client   client.Client
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}
