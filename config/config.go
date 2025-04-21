package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr string `envconfig:"ADDR" default:":8080"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
