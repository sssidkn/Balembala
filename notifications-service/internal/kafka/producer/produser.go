package producer

import (
	kafkaSettings "notifications/internal/kafka"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const timeout = 10000

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg kafkaSettings.Config, id string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.Brokers,
		"security.protocol":  cfg.SecurityProtocol,
		"client.id":          id,
		"acks":               "all",
		"enable.idempotence": "true",
	}
	if strings.EqualFold(cfg.SecurityProtocol, "ssl") {
		_ = conf.SetKey("ssl.ca.location", cfg.CaCertPath)
		_ = conf.SetKey("ssl.certificate.location", cfg.ClientCertPath)
		_ = conf.SetKey("ssl.key.location", cfg.KeyFilePath)
		_ = conf.SetKey("ssl.key.password", "123456")
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Start(msg []byte, topic string, partition int32) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          msg,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}
	e := <-kafkaChan
	switch e.(type) {
	case *kafka.Message:
		return nil
	default:
		return e.(error)
	}
}

func (p *Producer) Close() {
	p.producer.Flush(timeout)
	p.producer.Close()
}
