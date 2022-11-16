package main

import (
	"context"

	app "github.com/keRin7/grpc-experiments/pkg/App"
	"github.com/sirupsen/logrus"
)

const ConfigFile = "config/config.yaml"

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	ctx, finish := context.WithCancel(context.Background())
	defer finish()
	config := app.NewConfig()
	config.ReadFromFile(ConfigFile)

	MyApp := app.New(config)
	MyApp.GrpcTest(ctx)
	MyApp.WebStart()
	//fmt.Scanln()

}
