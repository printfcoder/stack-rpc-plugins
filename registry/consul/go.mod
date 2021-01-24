module github.com/stack-labs/stack-rpc-plugins/registry/consul

go 1.14

replace (
	github.com/stack-labs/stack-rpc v1.0.0 => ../../../stack-rpc
)

require (
	github.com/hashicorp/consul/api v1.3.0
	github.com/stack-labs/stack-rpc v1.0.0
	github.com/mitchellh/hashstructure v1.0.0
)
