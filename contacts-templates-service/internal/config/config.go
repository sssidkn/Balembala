package config

import (
	"fmt"
	"report/pkg/cache"
	"report/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres postgres.Config
	GrpcPort int `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"9080"`
	RestPort int `yaml:"REST_PORT" env:"REST_PORT" env-default:"8082"`
	Cache    cache.Config
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read enviroment variables %w: ", err)
	}

	return &cfg, nil
}
