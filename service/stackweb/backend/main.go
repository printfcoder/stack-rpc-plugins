package main

import (
	"net/http"

	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins"
	cfg "github.com/stack-labs/stack-rpc/config"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/web"

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
	s := web.NewService(
		web.Name("stack.stackweb"),
	)
	if err := s.Init(web.AfterStart(func() error {
		loadPlugins(s)
		return nil
	})); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

func loadPlugins(s web.Service) {
	rootPath := stackwebConfig.Stack.Stackweb.RootPath
	staticDir := stackwebConfig.Stack.Stackweb.StaticDir
	apiPath := stackwebConfig.Stack.Stackweb.ApiPath
	log.Infof("stackweb runs rootPath at %s", rootPath)
	log.Infof("stackweb deploys staticDir at %s", staticDir)
	log.Infof("stackweb applies rootPath at %s", apiPath)

	// favicon.ico todo
	// s.HandleFunc("/favicon.ico", faviconHandler)
	// static dir
	s.Handle(rootPath+"/", http.StripPrefix(rootPath+"/", http.FileServer(http.Dir(staticDir))))

	for _, m := range plugins.Plugins() {
		err := m.Init(plugins.Client(s.Options().Service.Client()))
		if err != nil {
			log.Errorf("plugin [%s] init err: %s", m.Name(), err)
			continue
		}

		r := m.Path()
		for k, h := range m.Handlers() {
			route := rootPath + apiPath + r + k

			if h.IsFunc() {
				s.HandleFunc(route, h.Func)
			} else {
				log.Infof("handler Handle: %s", route)
				s.Handle(route, h.Hld)
			}
		}
	}
}
