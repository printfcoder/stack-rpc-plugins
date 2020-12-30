package cmd

import (
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"time"

	z "github.com/stack-labs/stack-rpc-plugins/service/stackweb/internal/zap"
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/config/source/file"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/web"
)

var (
	re               = regexp.MustCompile("^[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9]*)?$")
	address          = ":9082"
	name             = "go.micro.web.platform"
	version          = "1.0.1-beta"
	rootPath         = "/platform"
	apiPath          = "/api/v1"
	configFile       = "./conf/micro.yml"
	StaticDir        = "webapp"
	registerTTL      = 30 * time.Second
	registerInterval = 10 * time.Second
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webapp/favicon.ico")
	return
}

func loadConfig(ctx *cli.Context) {
	if len(ctx.String("config_file")) > 0 {
		configFile = ctx.String("config_file")
	}
	log.Infof("[loadConfig] load config file: %s", configFile)

	if err := config.Load(file.NewSource(file.WithPath(configFile))); err != nil {
		panic(err)
	}
}

func loadModules(ctx *cli.Context, s web.Service) {
	// init modules
	for _, m := range plugins.Plugins() {
		logger.Info("loading module：", zap.Any("module", m.Name()))

		m.Init(ctx)
		r := m.Path()

		for k, h := range m.Handlers() {
			route := rootPath + apiPath + r + k

			if h.IsFunc() {
				logger.Info("handler Func", zap.Any("route", route))
				s.HandleFunc(route, h.Func)
			} else {
				logger.Info("handler Handle", zap.Any("route", route))
				s.Handle(route, h.Hld)
			}
		}
	}
}

func parseFlags(ctx *cli.Context) {
	if len(ctx.String("server_name")) > 0 {
		name = ctx.String("server_name")
	}

	if len(ctx.String("server_version")) > 0 {
		version = ctx.String("server_version")
	}

	if len(ctx.String("namespace")) > 0 {
		namespace = ctx.String("namespace")
	}

	if len(ctx.String("address")) > 0 {
		version = ctx.String("address")
	}

	if len(ctx.String("root_path")) > 0 {
		rootPath = ctx.String("root_path")
	}

	if len(ctx.String("address")) > 0 {
		address = ctx.String("address")
	}

	if len(ctx.String("register_ttl")) > 0 {
		registerTTL = ctx.Duration("register_ttl")
	}

	if len(ctx.String("register_interval")) > 0 {
		registerInterval = ctx.Duration("register_interval")
	}
}
