package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/samber/lo"
	"os"
	"strconv"

	"github.com/liluoliluoli/gnboot/internal/server"

	"github.com/go-cinch/common/log"
	"github.com/go-cinch/common/log/caller"
	_ "github.com/go-cinch/common/plugins/gorm/filter"
	"github.com/go-cinch/common/plugins/k8s/pod"
	"github.com/go-cinch/common/plugins/kratos/config/env"
	_ "github.com/go-cinch/common/plugins/kratos/encoding/yml"
	"github.com/go-cinch/common/utils"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liluoliluoli/gnboot/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "gnboot"
	// EnvPrefix is the prefix of the env params
	EnvPrefix = "SERVICE"
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	flagConf string
	// beforeReadConfigLogLevel is log level before read config.
	beforeReadConfigLogLevel = log.InfoLevel

	id, _ = os.Hostname()
)

func init() {
	wd, _ := os.Getwd()
	env := os.Getenv("APP_ENV")
	configPath := lo.Ternary(env == "", "/configs/dev", "/configs/"+env)
	flag.StringVar(&flagConf, "c", wd+configPath, "")
}

func newApp(gs *grpc.Server, hs *http.Server, jb *server.Job) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(log.DefaultWrapper.Options().Logger()),
		kratos.Server(
			gs,
			hs,
			jb,
		),
	)
}

func main() {
	flag.Parse()
	// set default log before read config
	logOps := []func(*log.Options){
		log.WithJSON(false),
		log.WithLevel(beforeReadConfigLogLevel),
		log.WithValuer("service.id", id),
		log.WithValuer("service.name", Name),
		log.WithValuer("service.version", Version),
		log.WithValuer("trace.id", tracing.TraceID()),
		log.WithValuer("span.id", tracing.SpanID()),
		log.WithCallerOptions(
			caller.WithSource(false),
			caller.WithLevel(2),
			caller.WithVersion(true),
		),
	}
	log.DefaultWrapper = log.NewWrapper(logOps...)

	//sc := []constant.ServerConfig{
	//	*constant.NewServerConfig("Nacos", 8848),
	//}
	//
	//cc := &constant.ClientConfig{
	//	NamespaceId:         "8e3d53a1-e2b4-45fa-ab89-8c9c7cc0d7cc",
	//	TimeoutMs:           5000,
	//	NotLoadCacheAtStart: true,
	//	LogDir:              "./docker-compose/log",
	//	CacheDir:            "./docker-compose/cache",
	//	LogLevel:            "debug",
	//}
	//
	//client, err := clients.NewConfigClient(
	//	vo.NacosClientParam{
	//		ClientConfig:  cc,
	//		ServerConfigs: sc,
	//	},
	//)
	//if err != nil {
	//	log.Error(err)
	//}
	c := config.New(
		//config.WithSource(file.NewSource(flagConf), knacos.NewConfigSource(
		//	client,
		//	knacos.WithGroup("DEFAULT_GROUP"),
		//	knacos.WithDataID("dynamic.yml"),
		//)),
		config.WithSource(file.NewSource(flagConf)),
		config.WithResolver(
			env.NewRevolver(
				env.WithPrefix(EnvPrefix),
				env.WithLoaded(func(k string, v interface{}) {
					log.Info("env loaded: %s=%v", k, v)
				}),
			),
		),
	)
	defer c.Close()

	fields := log.Fields{
		"conf": flagConf,
	}
	if err := c.Load(); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("load conf failed")
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("scan conf failed")
	}
	bc.Name = Name
	bc.Version = Version
	// override log level after read config
	logOps = append(logOps,
		[]func(*log.Options){
			log.WithLevel(log.NewLevel(bc.Log.Level)),
			log.WithJSON(bc.Log.JSON),
		}...,
	)
	log.DefaultWrapper = log.NewWrapper(logOps...)
	if bc.Server.MachineId == "" {
		// if machine id not set, gen from pod ip
		machineId, err := pod.MachineId()
		if err == nil {
			bc.Server.MachineId = strconv.FormatUint(uint64(machineId), 10)
		} else {
			bc.Server.MachineId = "0"
		}
	}
	// os.Setenv("COPIERX_UTC", "true")

	app, cleanup, err := wireApp(&bc)
	if err != nil {
		str := utils.Struct2Json(&bc)
		log.
			WithError(err).
			Error("wire app failed")
		// env str maybe very long, log with another line
		log.
			WithFields(fields).
			Fatal(str)
	}
	defer cleanup()

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("run app failed")
	}
}
