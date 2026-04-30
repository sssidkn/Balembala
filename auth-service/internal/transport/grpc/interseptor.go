package grpc

import (
	"auth/pkg/logger"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var requestID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				requestID = ids[0]
				ctx = context.WithValue(ctx, logger.RequestID, requestID)
			}
		}

		if requestID == "" {
			requestID = uuid.New().String()
			ctx = context.WithValue(ctx, logger.RequestID, requestID)
			l.Info(ctx, "Direct gRPC request started",
				zap.String("method", info.FullMethod),
			)
		} else {
			l.Info(ctx, "gRPC request (via HTTP)",
				zap.String("method", info.FullMethod),
			)
		}

		return handler(ctx, req)
	}
}
