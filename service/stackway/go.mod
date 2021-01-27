module github.com/stack-labs/stack-rpc-plugins/service/stackway

go 1.14

replace (
	github.com/stack-labs/stack-rpc v1.0.1 => ../../../stack-rpc
)

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/stack-labs/stack-rpc v1.0.1
	github.com/stretchr/testify v1.4.0
)
