package http

import (
	"context"
	"net/http"
	"notifications/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func HTTPMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		md := metadata.New(map[string]string{
			"x-request-id": requestID,
			"x-user-id":    r.Header.Get("x-user-id"),
		})

		ctx = metadata.NewOutgoingContext(ctx, md)
		ctx = context.WithValue(ctx, logger.RequestID, requestID)

		logger.GetLoggerFromContext(ctx).Info(
			ctx,
			"request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("user-id", r.Header.Get("x-user-id")),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
