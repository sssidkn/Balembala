package main

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"
	"auth/internal/transport/grpc"
	"auth/internal/transport/http"
	_ "auth/pkg/api/auth"
	"auth/pkg/db/postgres"
	"auth/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	)
	defer stop()

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	repo := repository.NewUsersRepository(db)
	authService := service.NewAuthService(repo)

	grpcSrv, err := grpc.New(ctx, cfg.GRPCServerPort, *authService)
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

	<-ctx.Done()

	logger.GetLoggerFromContext(ctx).Info(ctx, "shutting down")
	grpcSrv.Stop(ctx)
	httpSrv.Stop(ctx)
	logger.GetLoggerFromContext(ctx).Info(ctx, "service stopped")
}
