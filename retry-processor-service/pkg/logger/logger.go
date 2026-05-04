package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	Key       = "logger"
	Partition = "partition"
)

type Logger struct {
	Logger *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromContext(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(Partition) != nil {
		fields = append(fields, zap.Int(Partition, ctx.Value(Partition).(int)))
	}
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(Partition) != nil {
		fields = append(fields, zap.Int(Partition, ctx.Value(Partition).(int)))
	}
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(Partition) != nil {
		fields = append(fields, zap.Int(Partition, ctx.Value(Partition).(int)))
	}
	l.Logger.Fatal(msg, fields...)
}
