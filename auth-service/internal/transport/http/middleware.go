package http

import (
	"auth/pkg/logger"
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func HTTPMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, logger.RequestID, requestID)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			"x-request-id", requestID,
		))
		logger.GetLoggerFromContext(ctx).Info(
			ctx,
			"request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
