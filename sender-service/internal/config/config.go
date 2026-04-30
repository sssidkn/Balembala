package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Brokers          string `env:"KAFKA_BROKERS" env-default:"localhost:9092,localhost:9094,localhost:9096" env-separator:","`
	SecurityProtocol string `env:"KAFKA_SECURITY_PROTOCOL" env-default:"plaintext"`
	InputTopic       string `env:"KAFKA_INPUT_TOPIC" env-default:"notifications"`
	RetryTopic       string `env:"KAFKA_RETRY_TOPIC" env-default:"retry"`
	CaCertPath       string `env:"CA_CERT_PATH"`
	ClientCertPath   string `env:"CLIENT_CERT_PATH"`
	KeyFilePath      string `env:"CLIENT_KEY_PATH"`

	Username string `env:"USERNAME" env-default:"notificat2025@mail.ru"`
	Password string `env:"PASSWORD" env-default:"Lm1gNygA98dCedEfx1Fc"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read enviroment variables: %w", err)
	}
	return &cfg, nil
}
