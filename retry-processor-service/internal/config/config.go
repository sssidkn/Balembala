package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Brokers            string `env:"KAFKA_BROKERS" env-default:"localhost:9092,localhost:9094,localhost:9096" env-separator:","`
	SecurityProtocol   string `env:"KAFKA_SECURITY_PROTOCOL" env-default:"plaintext"`
	NotificationsTopic string `env:"KAFKA_NOTIFICATIONS_TOPIC" env-default:"notifications"`
	RetryTopic         string `env:"KAFKA_RETRY_TOPIC" env-default:"retry"`
	DLQTopic           string `env:"KAFKA_DLQ_TOPIC" env-default:"dlq"`
	CaCertPath         string `env:"CA_CERT_PATH"`
	ClientCertPath     string `env:"CLIENT_CERT_PATH"`
	KeyFilePath        string `env:"CLIENT_KEY_PATH"`
	MaxRetries         int    `env:"MAX_RETRIES" env-default:"3"`
	PartitionsCount    int    `env:"PARTITIONS_COUNT" env-default:"3"`
	ChannelSize        int    `env:"CHANNEL_SIZE" env-default:"100"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read enviroment variables: %w", err)
	}
	return &cfg, nil
}
