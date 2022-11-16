package webServer

import (
	"net/http"
)

type Config struct {
	mux  *http.ServeMux
	Port string `yaml:"LISTEN_PORT"`
}

func NewConfig() *Config {
	return &Config{
		mux: http.NewServeMux(),
	}
}
