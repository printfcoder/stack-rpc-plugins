package plugins

import (
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/registry"
)

type Options struct {
	// todo move all of below to runtime
	Registry registry.Registry
	Selector selector.Selector
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
