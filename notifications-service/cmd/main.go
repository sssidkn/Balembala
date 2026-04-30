package main

import (
	"context"
	"encoding/json"
	"fmt"
	report "notifications/internal/client"
	"notifications/internal/config"
	"notifications/internal/kafka/producer"
	"notifications/internal/models"
	"notifications/internal/service"
	server "notifications/internal/transport/grpc"
	"notifications/internal/transport/http"
	cache "notifications/pkg/db"
	"notifications/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

const (
	topic = "notifications"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	ctx, err := logger.New(context.Background())
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	reportClient, err := report.NewClient(cfg.ReportServer)
	if err != nil {
		panic(err)
	}
	msgChannels := make([]chan models.KafkaMsg, cfg.Kafka.ChannelNumber)
	for i := range msgChannels {
		msgChannels[i] = make(chan models.KafkaMsg, cfg.Kafka.ChannelSize)
	}
	redis := cache.New(cfg.Redis)
	notificationsService := service.NewNotificationsService(reportClient, redis, msgChannels)
	notificationsService.SetupBatchSize(cfg.Kafka.BatchSize)

	grpcSrv, err := server.New(ctx, cfg.GRPCServerPort, *notificationsService)
	if err != nil {
		panic(err)
	}
	go func() {
		err = grpcSrv.Start(ctx)
		if err != nil {
			stop()
		}
	}()

	httpSrv, err := http.New(ctx, fmt.Sprintf("localhost:%d", cfg.GRPCServerPort), cfg.HTTPServerPort)
	if err != nil {
		panic(err)
	}
	go func() {
		err = httpSrv.Start(ctx)
		if err != nil {
			stop()
		}
	}()
	var wg sync.WaitGroup
	for i := range cfg.Kafka.ChannelNumber {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := producer.NewProducer(cfg.Kafka, fmt.Sprintf("notifications-producer-%d", i))
			if err != nil {
				logger.GetLoggerFromContext(ctx).Error(ctx, "failed to create producer", zap.Error(err))
				return
			}
			defer p.Close()

			for msg := range msgChannels[i] {
				msgString, err := json.Marshal(msg)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Error(ctx, "failed to marshal message", zap.Error(err))
				}
				err = p.Start(msgString, topic, int32(i))
				if err != nil {
					logger.GetLoggerFromContext(ctx).Error(ctx, "failed to start producer", zap.Error(err))
				}
				logger.GetLoggerFromContext(ctx).Info(ctx, "message sent", zap.String("message", string(msgString)))
			}
		}()
	}
	<-ctx.Done()

	logger.GetLoggerFromContext(ctx).Info(ctx, "shutting down")
	grpcSrv.Stop(ctx)
	httpSrv.Stop(ctx)
	for _, ch := range msgChannels {
		close(ch)
	}

	logger.GetLoggerFromContext(ctx).Info(ctx, "service stopped")
}
