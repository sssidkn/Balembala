package config

import (
	"auth/pkg/db/postgres"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres       postgres.Config
	Secret         string `env:"SECRET" env-default:"secret"`
	GRPCServerPort int    `env:"GRPC_PORT" env-default:"9090"`
	HTTPServerPort int    `env:"HTTP_PORT" env-default:"8081"`
	Host           string `env:"HOST" env-default:"localhost"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config.New: failed to read config: %v", err)
	}
	return &cfg, nil
}
