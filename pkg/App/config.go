package app

import (
	"io/ioutil"

	"github.com/keRin7/grpc-experiments/pkg/grpcHealthCheck"
	"github.com/keRin7/grpc-experiments/pkg/webServer"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GrpcConfigs  []*grpcHealthCheck.Config `yaml:"GRPC_HOSTS"`
	QueryTimeout int                       `yaml:"QUERY_TIMEOUT"`
	LogLevel     string                    `default:"warn" yaml:"LOG_LEVEL"`
	WebService   *webServer.Config         `yaml:"WEB_SERVER"`
}

func NewConfig() *Config {
	GrpcConfigs := make([]*grpcHealthCheck.Config, 0, 2)
	return &Config{
		GrpcConfigs: GrpcConfigs,
		WebService:  webServer.NewConfig(),
	}
}

func (k *Config) ReadFromFile(fileName string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, k)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
	}
}
