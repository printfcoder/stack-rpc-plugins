package main

import (
	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins"
	cfg "github.com/stack-labs/stack-rpc/config"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/web"

	_ "github.com/stack-labs/stack-rpc-plugins/logger/logrus"
	_ "github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins/basic"
)

type config struct {
	Stack struct {
		Stackweb struct {
			Name       string `sc:"name"`
			Address    string `sc:"address"`
			ApiPath    string `sc:"api-path"`
			RootPath   string `sc:"root-path"`
			StaticDir  string `sc:"static-dir"`
			FaviconIco string `sc:"favicon-ico"`
		} `sc:"stackweb"`
	} `sc:"stack"`
}

var (
	stackwebConfig config
)

func init() {
	cfg.RegisterOptions(&stackwebConfig)
}

func main() {
	s := stack.NewWebService(
		stack.Name("stack.stackweb"),
		stack.WebHandleFuncs(),
	)
	if err := s.Init(stack.AfterStart(func() error {
		loadPlugins(s.Options())
		return nil
	})); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

func loadPlugins(s service.Options) {
	for _, m := range plugins.Plugins() {
		err := m.Init(plugins.Service(s))
		if err != nil {
			log.Errorf("plugin [%s] init err: %s", m.Name(), err)
			continue
		}

		for k, h := range m.Handlers() {
			if h.IsFunc() {
				web.HandleFuncs(web.HandlerFunc{
					Route: k,
					Func:  h.Func,
				})
			}
		}
	}
}
