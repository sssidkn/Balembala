package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"sender-service/internal/config"
	"sender-service/internal/dto"
	"sender-service/pkg/logger"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const timeout = 10000

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.Brokers,
		"security.protocol":  cfg.SecurityProtocol,
		"client.id":          "sender-service",
		"acks":               "all",
		"enable.idempotence": "true",
	}
	if strings.EqualFold(cfg.SecurityProtocol, "ssl") {
		_ = conf.SetKey("ssl.ca.location", cfg.CaCertPath)
		_ = conf.SetKey("ssl.certificate.location", cfg.ClientCertPath)
		_ = conf.SetKey("ssl.key.location", cfg.KeyFilePath)
		_ = conf.SetKey("ssl.key.password", "123456")
		_ = conf.SetKey("debug", "security,broker")
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p}, nil

}

func (p *Producer) Produce(ctx context.Context, cfg config.Config, msg dto.Message) error {
	logg := logger.GetLoggerFromContext(ctx)
	val, err := json.Marshal(msg)
	if err != nil {
		logg.Error(ctx, fmt.Sprintf("marshal error: %v", err))
		return err
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &cfg.RetryTopic, Partition: kafka.PartitionAny},
		Value:          val,
	}

	kafkaChan := make(chan kafka.Event)
	if err = p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}
	e := <-kafkaChan
	switch e.(type) {
	case *kafka.Message:
		logg.Info(ctx, "sent message successfully")
		return nil
	default:
		return e.(error)
	}

}

func (p *Producer) Close() {
	p.producer.Flush(timeout)
	p.producer.Close()
}
