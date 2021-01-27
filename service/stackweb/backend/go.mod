module github.com/stack-labs/stack-rpc-plugins/service/stackweb

go 1.14

replace (
	github.com/stack-labs/stack-rpc v1.0.1 => ../../../../stack-rpc
	github.com/stack-labs/stack-rpc-plugins/logger/logrus v1.0.0 => ../../../../stack-rpc-plugins/logger/logrus
)

require (
	github.com/stack-labs/stack-rpc v1.0.1
	github.com/stack-labs/stack-rpc-plugins/logger/logrus v1.0.0
)
