package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("successful logger creation", func(t *testing.T) {
		ctx := context.Background()

		newCtx, err := New(ctx)
		require.NoError(t, err)
		require.NotNil(t, newCtx)

		logger := GetLoggerFromContext(newCtx)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.Logger)
	})

	t.Run("failed logger creation", func(t *testing.T) {
		t.Skip("This test requires mocking zap.NewProduction")
	})
}

func TestGetLoggerFromContext(t *testing.T) {
	t.Run("logger exists in context", func(t *testing.T) {
		ctx := context.Background()
		core, _ := observer.New(zap.InfoLevel)
		logger := zap.New(core)

		ctx = context.WithValue(ctx, Key, &Logger{logger})

		retrieved := GetLoggerFromContext(ctx)
		assert.NotNil(t, retrieved)
		assert.Equal(t, logger, retrieved.Logger)
	})

	t.Run("logger not in context", func(t *testing.T) {
		ctx := context.Background()

		// Это вызовет панику, так как происходит type assertion к *Logger
		assert.Panics(t, func() {
			GetLoggerFromContext(ctx)
		})
	})

	t.Run("wrong type in context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, Key, "not a logger")

		assert.Panics(t, func() {
			GetLoggerFromContext(ctx)
		})
	})
}

func TestLoggerMethods(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		method   func(*Logger, context.Context, string, ...zap.Field)
		msg      string
		fields   []zap.Field
		ctx      context.Context
		expected []zap.Field
	}{
		{
			name:   "Info without partition",
			method: (*Logger).Info,
			msg:    "test info",
			fields: []zap.Field{zap.String("key", "value")},
			ctx:    ctx,
			expected: []zap.Field{
				zap.String("key", "value"),
			},
		},
		{
			name:   "Info with partition",
			method: (*Logger).Info,
			msg:    "test info with partition",
			fields: []zap.Field{zap.String("key", "value")},
			ctx:    context.WithValue(ctx, Partition, 42),
			expected: []zap.Field{
				zap.String("key", "value"),
				zap.Int(Partition, 42),
			},
		},
		{
			name:   "Error without partition",
			method: (*Logger).Error,
			msg:    "test error",
			fields: []zap.Field{zap.String("error", "something went wrong")},
			ctx:    ctx,
			expected: []zap.Field{
				zap.String("error", "something went wrong"),
			},
		},
		{
			name:   "Error with partition",
			method: (*Logger).Error,
			msg:    "test error with partition",
			fields: []zap.Field{zap.String("error", "something went wrong")},
			ctx:    context.WithValue(ctx, Partition, 99),
			expected: []zap.Field{
				zap.String("error", "something went wrong"),
				zap.Int(Partition, 99),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := &Logger{zap.New(core)}

			tt.method(logger, tt.ctx, tt.msg, tt.fields...)

			logs := recorded.All()
			require.Len(t, logs, 1, "Expected exactly one log entry")
			assert.Equal(t, tt.msg, logs[0].Message)

			for _, expectedField := range tt.expected {
				found := false
				for _, actualField := range logs[0].Context {
					if actualField.Key == expectedField.Key {
						found = true
						assert.Equal(t, expectedField.Interface, actualField.Interface)
						break
					}
				}
				assert.True(t, found, "Expected field %q not found in log entry", expectedField.Key)
			}

			assert.Len(t, logs[0].Context, len(tt.expected),
				"Unexpected number of fields in log entry")
		})
	}
}

func TestFatal(t *testing.T) {
	logger := &Logger{zap.NewNop()}
	assert.NotNil(t, logger.Fatal)
}
