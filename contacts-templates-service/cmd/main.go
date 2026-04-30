package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"report/internal/config"
	"report/internal/repository"
	"report/internal/service"
	"report/pkg/api/report"
	"report/pkg/cache"
	"report/pkg/logger"
	"report/pkg/postgres"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, err := logger.New(ctx)
	if err != nil {
		panic(err)
	}
	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to read config", zap.Error(err))
	}

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to db", zap.Error(err))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen: %v", zap.Error(err))
	}

	redisCache := cache.New(cfg.Cache)
	repo := repository.NewRepository(db)
	srv := service.NewService(repo, redisCache)
	server := grpc.NewServer(grpc.UnaryInterceptor(logger.LogMiddleware(*logger.GetLoggerFromCtx(ctx))))
	report.RegisterReportServiceServer(server, srv)
	go func() {
		if err := server.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "server start grpc")
		}
	}()

	mux := runtime.NewServeMux(runtime.WithMetadata(service.HTTPToGRPCMiddleware))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = report.RegisterReportServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", cfg.GrpcPort), opts)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to register test server", zap.Error(err))
	}
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.RestPort), mux); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "server start rest")
		}
	}()

	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	server.GracefulStop()
	redisCache.Close()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
}
