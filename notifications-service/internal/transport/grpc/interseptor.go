package server

import (
	"context"
	"notifications/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var requestID string
		var userID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				requestID = ids[0]
			}

			if uids := md.Get("x-user-id"); len(uids) > 0 {
				userID = uids[0]
				ctx = context.WithValue(ctx, "x-user-id", userID)
			}
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, logger.RequestID, requestID)

		logFields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("request-id", requestID),
		}

		if userID != "" {
			logFields = append(logFields, zap.String("user-id", userID))
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok && len(md.Get("x-request-id")) > 0 {
			l.Info(ctx, "gRPC request (via HTTP)", logFields...)
		} else {
			l.Info(ctx, "Direct gRPC request started", logFields...)
		}

		return handler(ctx, req)
	}
}
