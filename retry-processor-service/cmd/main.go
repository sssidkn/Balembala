package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rebalancer/internal/config"
	"rebalancer/internal/kafka/consumer"
	"rebalancer/internal/kafka/producer"
	"rebalancer/internal/models"
	"rebalancer/pkg/logger"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	ctx, err := logger.New(ctx)
	if err != nil {
		panic(err)
	}
	l := logger.GetLoggerFromContext(ctx)

	cfg, err := config.NewConfig()
	if err != nil {
		l.Fatal(ctx, "failed to load config: %v", zap.Error(err))
	}

	dlqTopic := cfg.DLQTopic
	retryTopic := cfg.RetryTopic
	notificationsTopic := cfg.NotificationsTopic
	dlqProducers := make([]*producer.Producer, cfg.PartitionsCount)
	notificationProducers := make([]*producer.Producer, cfg.PartitionsCount)

	for i := 0; i < cfg.PartitionsCount; i++ {
		dlqProducers[i], err = producer.NewProducer(cfg)
		if err != nil {
			l.Fatal(ctx, fmt.Sprintf("failed to create DLQ producer %d: %v", i, zap.Error(err)))
		}
		defer dlqProducers[i].Close()

		notificationProducers[i], err = producer.NewProducer(cfg)
		if err != nil {
			l.Fatal(ctx, fmt.Sprintf("failed to create Notification producer %d: %v", i, zap.Error(err)))
		}
		defer notificationProducers[i].Close()
	}

	retryChans := make([]chan models.Message, cfg.PartitionsCount)
	for i := range retryChans {
		retryChans[i] = make(chan models.Message, cfg.ChannelSize)
	}

	retryConsumers := make([]*consumer.Consumer, cfg.PartitionsCount)
	for i := 0; i < cfg.PartitionsCount; i++ {
		retryConsumers[i], err = consumer.NewConsumer(*cfg, retryTopic, "retry-group")
		if err != nil {
			l.Fatal(ctx, fmt.Sprintf("failed to create retry consumer %d: %v", i, zap.Error(err)))
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < cfg.PartitionsCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			retryConsumers[idx].Consume(ctx, retryChans[idx])
		}(i)
	}

	for i := 0; i < cfg.PartitionsCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-retryChans[idx]:
					producerIdx := idx % cfg.PartitionsCount

					if msg.RetryCount >= cfg.MaxRetries {
						l.Info(ctx, fmt.Sprintf("message to %v sent to DLQ by dlqProducer-%d", msg.ToList, idx))
						if err := dlqProducers[producerIdx].Produce(ctx, dlqTopic, int32(producerIdx), msg); err != nil {
							l.Error(ctx, fmt.Sprintf("producer %d failed to produce to DLQ topic: %v",
								idx, zap.Error(err)))
						}
					} else {
						l.Info(ctx, fmt.Sprintf("message to %v sent to notifications by producer %d (retry: %d)",
							msg.ToList, idx, msg.RetryCount))
						if err := notificationProducers[producerIdx].Produce(ctx, notificationsTopic, int32(producerIdx),
							msg); err != nil {
							l.Error(ctx, fmt.Sprintf("producer %d failed to produce to notification topic: %v",
								idx, zap.Error(err)))
						}
					}
				}
			}
		}(i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	l.Info(ctx, "shutting down retry processor...")
	cancel()
	wg.Wait()
	l.Info(ctx, "shutdown complete")
}
