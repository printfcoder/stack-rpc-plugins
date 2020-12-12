package hook

import (
	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/util/log"

	"github.com/stack-labs/stack-rpc-plugins/service/stackway/api"
	"github.com/stack-labs/stack-rpc-plugins/service/stackway/plugin"
)

func Hook(svc stack.Service) {
	apiServer := api.NewServer(svc)

	// stackway options
	_ = svc.Init(api.Options()...)

	// stackway hook
	_ = svc.Init(
		stack.AfterStart(apiServer.Start),
		stack.AfterStop(apiServer.Stop),
	)

	// plugin tags
	plugins := plugin.Plugins()
	for _, p := range plugins {
		log.Debugf("plugin: %s", p.String())
		if flags := p.Flags(); len(flags) > 0 {
			log.Debugf("flags: %+#s", flags)
			_ = svc.Init(stack.Flags(flags...))
		}
	}

	return
}
