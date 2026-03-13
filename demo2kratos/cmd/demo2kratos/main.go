package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/yylego/done"
	"github.com/yylego/kratos-config/configkratos"
	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
	"github.com/yylego/must"
	"github.com/yylego/rese"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(done.VCE(os.Hostname()).Omit()),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", kratos.ID(done.VCE(os.Hostname()).Omit()),
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// demo2 uses DataSource to load dynamic config with watch support
	// demo2 使用 DataSource 加载动态配置，支持 Watch 机制
	var sources []config.Source
	for _, item := range rese.A1(os.ReadDir(flagconf)) {
		configPath := filepath.Join(flagconf, item.Name())
		configData := rese.A1(os.ReadFile(configPath))
		sources = append(sources, configkratos.NewDataSource(configData, "yaml"))
	}
	c := config.New(
		config.WithSource(sources...),
	)
	defer rese.F0(c.Close)

	must.Done(c.Load())

	var cfg conf.Bootstrap
	must.Done(c.Scan(&cfg))

	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, logger))
	defer cleanup()

	// start and wait for stop signal
	must.Done(app.Run())
}
