package config

import (
	"fmt"
	kafkaSettings "notifications/internal/kafka"
	cache "notifications/pkg/db"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServerPort int    `env:"GRPC_PORT" env-default:"8095"`
	HTTPServerPort int    `env:"HTTP_PORT" env-default:"8090"`
	Host           string `env:"HOST" env-default:""`
	ReportServer   string `env:"REPORT_SERVER" env-default:"localhost:9080"`
	Kafka          kafkaSettings.Config
	Redis          cache.Config
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config.New: failed to read config: %v", err)
	}
	return &cfg, nil
}
