package cmd

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/backend/internal/proxy"
	z "github.com/stack-labs/stack-rpc-plugins/service/stackweb/internal/zap"
	"github.com/stack-labs/stack-rpc-plugins/service/stackweb/plugins"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/config/cmd"
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
	namespace        = "go.micro.web"
	registerTTL      = 30 * time.Second
	registerInterval = 10 * time.Second
	logger           = z.GetLogger()
)

// Init app
func Init(ops ...plugins.Option) {
	app := cmd.App()
	app.Flags = append(app.Flags,
		&cli.StringFlag{
			Name:    "root_path",
			Usage:   "Set the root path of micro web",
			EnvVars: []string{"MICRO_WEB_NAMESPACE"},
		},
		&cli.StringFlag{
			Name:    "static_dir",
			Usage:   "Set the static dir of micro web",
			EnvVars: []string{"MICRO_WEB_STATIC_DIR"},
		},
	)

	app.Action = func(c *cli.Context) error {
		loadConfig(c)
		run(c)
		return nil
	}

	if err := cmd.Init(cmd.Name(name)); err != nil {
		panic(err)
	}
}

func run(ctx *cli.Context, srvOpts ...plugins.Option) {
	parseFlags(ctx)

	s := web.NewService(
		web.Name(name),
		web.Version(version),
		web.Address(address),
		web.RegisterTTL(registerTTL),
		web.RegisterInterval(registerInterval),
		web.Id(name+"-"+uuid.New().String()),
	)

	// favicon.ico
	s.HandleFunc("/favicon.ico", faviconHandler)

	// static dir
	s.Handle(rootPath+"/", http.StripPrefix(rootPath+"/", http.FileServer(http.Dir(StaticDir))))

	webProxyPath := rootPath + "/web/"
	s.Handle(webProxyPath, webProxy())

	logger.Info("handler web at ：", zap.Any("path", webProxyPath))

	loadModules(ctx, s)

	if err := s.Init(
		web.Action(
			func(c *cli.Context) {
				// do something
			}),
	); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

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

func webProxy() http.Handler {
	sel := selector.NewSelector(
		selector.Registry(*cmd.DefaultOptions().Registry),
	)

	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		parts := strings.Split(r.URL.Path, "/web/")
		if len(parts) < 2 {
			kill()
			return
		}

		if !re.MatchString(parts[1]) {
			kill()
			return
		}

		next, err := sel.Select(namespace + "." + parts[1])
		if err != nil {
			kill()
			return
		}

		s, err := next()
		if err != nil {
			kill()
			return
		}

		path := "/" + strings.Join(parts[2:], "/")

		logger.Debug("proxy to", zap.String("path", path))

		r.Header.Set(proxy.BasePathHeader, "/"+parts[1])
		r.URL.Host = fmt.Sprintf("%s", s.Address)
		r.URL.Path = path
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &proxy.Proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}
