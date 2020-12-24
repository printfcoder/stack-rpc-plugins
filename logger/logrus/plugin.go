package logrus

import (
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugin"
)

var options struct {
	Stack struct {
		Logger struct {
			Logrus struct {
				SplitLevel   bool `sc:"split-level"`
				ReportCaller bool `sc:"report-caller"`
			} `sc:"slogrus"`
		} `sc:"logger"`
	} `sc:"stack"`
}

type logrusLogPlugin struct {
}

func (l *logrusLogPlugin) Name() string {
	return "slogrus"
}

func (l *logrusLogPlugin) Options() []logger.Option {
	var opts []logger.Option

	opts = append(opts, SplitLevel(options.Stack.Logger.Logrus.SplitLevel))
	opts = append(opts, ReportCaller(options.Stack.Logger.Logrus.ReportCaller))

	return opts
}

func (l *logrusLogPlugin) New(opts ...logger.Option) logger.Logger {
	return NewLogger(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.LoggerPlugins["slogrus"] = &logrusLogPlugin{}
}
