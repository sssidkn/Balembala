package consumer

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

const (
	sessionTimeout = 10000
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewComsumer(cfg config.Config, consumerGroup string) (*Consumer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":        cfg.Brokers,
		"security.protocol":        cfg.SecurityProtocol,
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.commit":       "false",
		"enable.auto.offset.store": "false",
	}
	if strings.EqualFold(cfg.SecurityProtocol, "ssl") {
		_ = conf.SetKey("ssl.ca.location", cfg.CaCertPath)
		_ = conf.SetKey("ssl.certificate.location", cfg.ClientCertPath)
		_ = conf.SetKey("ssl.key.location", cfg.KeyFilePath)
		_ = conf.SetKey("ssl.key.password", "123456")
	}
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, err
	}
	err = c.Subscribe(cfg.InputTopic, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: c}, nil

}

func (c *Consumer) Consume(ctx context.Context, channel chan dto.Message) {
	logg := logger.GetLoggerFromContext(ctx)

	defer func() {
		if err := c.Stop(); err != nil {
			logg.Error(ctx, fmt.Sprintf("stop consumer error: %v", err))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				logg.Error(ctx, fmt.Sprintf("read error: %v", err))
				continue
			}

			var message dto.Message
			if err = json.Unmarshal(msg.Value, &message); err != nil {
				logg.Error(ctx, fmt.Sprintf("unmarshal error : %v", err))
				continue
			}

			select {
			case channel <- message:
				logg.Info(ctx, "get message successfully")
				if _, err = c.consumer.StoreMessage(msg); err != nil {
					logg.Error(ctx, fmt.Sprintf("failed to store message: %v\n", err))
				}
				if _, err = c.consumer.CommitMessage(msg); err != nil {
					logg.Error(ctx, fmt.Sprintf("failed to commit message: %v\n", err))
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	_, err := c.consumer.Commit()
	if err != nil {
		return err
	}
	return c.consumer.Close()
}
