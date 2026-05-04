package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sender-service/internal/config"
	"sender-service/internal/consumer"
	"sender-service/internal/dto"
	"sender-service/internal/producer"
	"sender-service/internal/sender"
	"sender-service/pkg/logger"
	"sync"
	"syscall"
)

const numPartitions = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, err := logger.New(ctx)
	if err != nil {
		log.Fatal("failed to initialize logger")
	}

	logg := logger.GetLoggerFromContext(ctx)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to load config: %v", err))
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	mailChannels := make([]chan dto.Message, numPartitions)
	for i := range mailChannels {
		mailChannels[i] = make(chan dto.Message, 100)
	}

	mailSender := sender.New(*cfg)

	prod, err := producer.NewProducer(cfg)
	if err != nil {
		logg.Fatal(ctx, fmt.Sprintf("failed to initialize producer: %v", err))
	}

	var wg sync.WaitGroup
	for i := 0; i < numPartitions; i++ {
		wg.Add(1)
		go func(partition int) {
			defer wg.Done()

			cont, err := logger.New(context.Background())
			if err != nil {
				logg.Error(ctx, "failed to add logger in context to consumer")
				cancel()
			}
			cons, err := consumer.NewComsumer(*cfg, "sender-service")
			if err != nil {
				logg.Error(ctx, "failed to create consumer")
				return
			}

			go func() {
				cons.Consume(cont, mailChannels[partition])
			}()

			for msg := range mailChannels[partition] {
				if err, retryMessage := mailSender.Send(msg); err != nil {
					logg.Error(ctx, fmt.Sprintf("error: %v\n", err))
					prod.Produce(cont, *cfg, retryMessage)
				} else {
					if len(retryMessage.ToList) > 0 {
						prod.Produce(cont, *cfg, retryMessage)
					}
					logg.Info(ctx, "send successfully")
				}
			}
		}(i)
	}

	select {
	case <-sigChan:
		logg.Info(ctx, "Received shutdown signal")
	case <-ctx.Done():
		logg.Info(ctx, "Shutdown initiated due to error")
	}

	logg.Info(ctx, "starting graceful shutdown...")
	cancel()

	for _, ch := range mailChannels {
		close(ch)
	}
	wg.Wait()

	prod.Close()

	logg.Info(ctx, "shutdown completed")
}
