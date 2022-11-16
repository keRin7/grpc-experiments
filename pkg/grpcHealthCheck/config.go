package grpcHealthCheck

import "sync"

type Config struct {
	Host          string `yaml:"GRPC_HOST"`
	Port          string `yaml:"GRPC_PORT"`
	TLS           bool   `yaml:"TLS"`
	FlConnTimeout int    `yaml:"CONN_TIMEOUT"`
	FlRPCTimeout  int    `yaml:"RPC_TIMEOUT"`
	FlService     string `yaml:"SERVICE"`
	mutex         sync.Mutex
}

func NewConfig() *Config {
	return &Config{}
}
