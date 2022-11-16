package grpcHealthCheck

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type IntGRPC struct {
	Config *Config
}

func New(config *Config) *IntGRPC {
	//logrus.Println("Config", config)
	return &IntGRPC{
		Config: config,
	}
}

func (c *IntGRPC) ShowHost() {
	logrus.Println("Host:", c.Config.Host)
}

func (c *IntGRPC) InitOptsTLS() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUserAgent("grpc_health_probe"),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}
}

func (c *IntGRPC) InitOptsNoTLS() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUserAgent("grpc_health_probe"),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
}

func (c *IntGRPC) Check(GRPCopts []grpc.DialOption) (bool, time.Duration, error) {
	var err error

	start := time.Now()

	GRPCctx, Cancel := context.WithCancel(context.Background())

	defer Cancel()

	GRPCDialctx, DialCancel := context.WithTimeout(GRPCctx, time.Duration(c.Config.FlConnTimeout)*time.Second)

	defer DialCancel()

	Conn, err := grpc.DialContext(GRPCDialctx, c.Config.Host+":"+c.Config.Port, GRPCopts...)
	if err != nil {
		return false, time.Since(start), err
	}

	GRPCRPCctx, RpcCancel := context.WithTimeout(GRPCctx, time.Duration(c.Config.FlRPCTimeout)*time.Second)

	defer RpcCancel()

	resp, err := healthpb.NewHealthClient(Conn).Check(GRPCRPCctx,
		&healthpb.HealthCheckRequest{
			Service: c.Config.FlService})

	c.Config.mutex.Lock()
	Conn.Close()
	c.Config.mutex.Unlock()

	elapsed := time.Since(start)

	if err != nil {
		return false, elapsed, err
	}

	if resp.GetStatus() == healthpb.HealthCheckResponse_SERVICE_UNKNOWN {
		return true, elapsed, nil
	}

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return false, elapsed, nil
	} else {
		return true, elapsed, nil
	}

}
