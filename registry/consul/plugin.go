package consul

import "github.com/stack-labs/stack-rpc/registry"

type consulRegistryPlugin struct{}

func (c *consulRegistryPlugin) Name() string {
	return "consul"
}

func (c *consulRegistryPlugin) Options() []registry.Option {
	return nil
}

func (c *consulRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}
