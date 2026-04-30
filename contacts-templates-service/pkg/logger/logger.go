package logger

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Key       = "logger"
	RequestID = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	l := &Logger{l: logger}
	ctx = context.WithValue(ctx, Key, l)
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if l, oke := ctx.Value(Key).(*Logger); oke {
		return l
	}
	return nil
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if requestID, ok := ctx.Value(RequestID).(string); ok {
		fields = append(fields, zap.String(RequestID, requestID))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if requestID, ok := ctx.Value(RequestID).(string); ok {
		fields = append(fields, zap.String(RequestID, requestID))
	}
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if requestID, ok := ctx.Value(RequestID).(string); ok {
		fields = append(fields, zap.String(RequestID, requestID))
	}
	l.l.Error(msg, fields...)
}

func LogMiddleware(l Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l.Info(ctx, "request started",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		rsp, err := handler(ctx, req)

		l.Info(ctx, "request completed",
			zap.String("method", info.FullMethod),
			zap.Any("response", rsp),
			zap.Error(err),
		)
		return rsp, err
	}
}
