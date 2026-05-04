package http

import (
	"auth/pkg/api/auth"
	"auth/pkg/logger"
	"context"
	"fmt"
	"net/http"

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
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux,
			marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
		}),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			if id, ok := ctx.Value(logger.RequestID).(string); ok {
				return metadata.Pairs("x-request-id", id)
			}
			return nil
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	handler := HTTPMiddleware(ctx, mux)

	if err := auth.RegisterAuthHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
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
