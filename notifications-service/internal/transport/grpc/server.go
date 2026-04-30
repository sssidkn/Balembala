package server

import (
	"context"
	"fmt"
	"net"
	"notifications/internal/service"
	"notifications/pkg/api/notifications"
	"notifications/pkg/logger"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context, port int, service service.NotificationsService) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(ContextWithLogger(*logger.GetLoggerFromContext(ctx))))
	grpcServer := grpc.NewServer(opts...)
	notifications.RegisterNotificationsServer(grpcServer, &service)
	return &Server{grpcServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Info(ctx, "starting grpc server")
	eg := errgroup.Group{}
	eg.Go(func() error {
		return s.grpcServer.Serve(s.listener)
	})
	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Info(ctx, "stopping grpc server")
	s.grpcServer.GracefulStop()
	return nil
}
