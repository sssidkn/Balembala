package http

import (
	"context"
	"fmt"
	"net/http"
	"notifications/pkg/api/notifications"
	"notifications/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	httpServer *http.Server
}

func New(ctx context.Context, grpcAddr string, httpPort int) (*Server, error) {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			if id, ok := ctx.Value(logger.RequestID).(string); ok {
				return metadata.Pairs("x-request-id", id,
					"x-user-id", r.Header.Get("x-user-id"))
			}
			return nil
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	handler := HTTPMiddleware(ctx, mux)

	if err := notifications.RegisterNotificationsHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, err
	}
	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", httpPort),
			Handler: handler,
		},
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Info(ctx, "starting http server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Info(ctx, "stopping http server")
	return s.httpServer.Shutdown(ctx)
}
