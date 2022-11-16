package app

import (
	"context"
	"time"

	"github.com/keRin7/grpc-experiments/pkg/grpcHealthCheck"
	"github.com/keRin7/grpc-experiments/pkg/vPrometheus"
	"github.com/keRin7/grpc-experiments/pkg/webServer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	Config    *Config
	GrpcHC    []*grpcHealthCheck.IntGRPC
	Prom      *vPrometheus.AppPrometheus
	WebServer *webServer.WebServer
}

func New(config *Config) *App {
	grpcHC := make([]*grpcHealthCheck.IntGRPC, 0, 2)

	switch config.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Warnf("Home: invalid log level supplied: '%s'", config.LogLevel)
	}

	for _, el := range config.GrpcConfigs {
		grpcHost := grpcHealthCheck.New(el)
		grpcHC = append(grpcHC, grpcHost)
	}
	return &App{
		Config:    config,
		GrpcHC:    grpcHC,
		Prom:      vPrometheus.CreatePrometheus(),
		WebServer: webServer.CreateWebServer(config.WebService),
	}
}

func (c *App) WebStart() {
	//c.prom.
	logrus.Printf("Web started on port: %s", c.Config.WebService.Port)
	c.Prom.InitPrometheus()
	c.WebServer.AddHeandler("/metrics", promhttp.Handler())
	c.WebServer.Start()
}

func (c *App) GrpcTest(ctx context.Context) {
	for _, v := range c.GrpcHC {
		go c.runTest(ctx, v)
	}
}

func (c *App) runTest(ctx context.Context, grpcHost *grpcHealthCheck.IntGRPC) {
	//logrus.Println("Start: " + grpc.Config.Host + "\n")
	var GRPCopts []grpc.DialOption
	if grpcHost.Config.TLS {
		GRPCopts = grpcHost.InitOptsTLS()
	} else {
		GRPCopts = grpcHost.InitOptsNoTLS()
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(c.Config.QueryTimeout) * time.Second)
			if ok, time, err := grpcHost.Check(GRPCopts); ok {
				logrus.Debug(grpcHost.Config.Host + " test OK" + time.String())
				c.Prom.WriteMetric(grpcHost.Config.Host, float64(time.Milliseconds()))
			} else {
				logrus.Debug(grpcHost.Config.Host + " test FAIL" + err.Error())
				c.Prom.WriteMetric(grpcHost.Config.Host, float64(0))
			}
		}
	}
}
