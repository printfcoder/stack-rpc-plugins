package main

import (
	"fmt"

	"github.com/stack-labs/stack-rpc"
	_ "github.com/stack-labs/stack-rpc-plugins/logger/logrus"
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins"
	_ "github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins/basic"
	cfg "github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/web"
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
		stack.Name("stack.rpc.stackweb"),
		stack.Address(":8090"),
		stack.WebHandleFuncs(handlers()...),
	)
	if err := s.Init(stack.AfterStart(func() error {
		return loadPlugins(s.Options())
	})); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

func handlers() []web.HandlerFunc {
	handlers := make([]web.HandlerFunc, 0)
	for _, m := range plugins.Plugins() {
		for k, h := range m.Handlers() {
			if h.IsFunc() {
				handlers = append(handlers, web.HandlerFunc{
					Route: k,
					Func:  h.Func,
				})
			}
		}
	}

	return handlers
}

func loadPlugins(s service.Options) (err error) {
	for _, m := range plugins.Plugins() {
		err := m.Init(plugins.Service(s))
		if err != nil {
			return fmt.Errorf("plugin [%s] init err: %s", m.Name(), err)
		}
	}

	return nil
}
