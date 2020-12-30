package main

import (
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/cmd"

	_ "github.com/stack-labs/stack-rpc-plugins/logger/logrus"
)

func main() {
	cmd.Init()
}
