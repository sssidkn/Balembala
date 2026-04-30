package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigDefaults(t *testing.T) {
	os.Clearenv()

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, "localhost:9092,localhost:9094,localhost:9096", cfg.Brokers)
	assert.Equal(t, "plaintext", cfg.SecurityProtocol)
	assert.Equal(t, "notifications", cfg.NotificationsTopic)
	assert.Equal(t, "retry", cfg.RetryTopic)
	assert.Equal(t, "dlq", cfg.DLQTopic)
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Equal(t, 3, cfg.PartitionsCount)
	assert.Equal(t, 100, cfg.ChannelSize)
}

func TestNewConfigReadsEnvironment(t *testing.T) {
	os.Clearenv()
	t.Setenv("KAFKA_BROKERS", "kafka1:8092,kafka2:8094,kafka3:8096")
	t.Setenv("KAFKA_SECURITY_PROTOCOL", "ssl")
	t.Setenv("KAFKA_NOTIFICATIONS_TOPIC", "notifications-custom")
	t.Setenv("KAFKA_RETRY_TOPIC", "retry-custom")
	t.Setenv("KAFKA_DLQ_TOPIC", "dlq-custom")
	t.Setenv("CA_CERT_PATH", "/etc/secrets/ca.crt")
	t.Setenv("CLIENT_CERT_PATH", "/etc/secrets/client.crt")
	t.Setenv("CLIENT_KEY_PATH", "/etc/secrets/client.key")
	t.Setenv("MAX_RETRIES", "5")
	t.Setenv("PARTITIONS_COUNT", "7")
	t.Setenv("CHANNEL_SIZE", "250")

	cfg, err := NewConfig()

	require.NoError(t, err)
	assert.Equal(t, "kafka1:8092,kafka2:8094,kafka3:8096", cfg.Brokers)
	assert.Equal(t, "ssl", cfg.SecurityProtocol)
	assert.Equal(t, "notifications-custom", cfg.NotificationsTopic)
	assert.Equal(t, "retry-custom", cfg.RetryTopic)
	assert.Equal(t, "dlq-custom", cfg.DLQTopic)
	assert.Equal(t, "/etc/secrets/ca.crt", cfg.CaCertPath)
	assert.Equal(t, "/etc/secrets/client.crt", cfg.ClientCertPath)
	assert.Equal(t, "/etc/secrets/client.key", cfg.KeyFilePath)
	assert.Equal(t, 5, cfg.MaxRetries)
	assert.Equal(t, 7, cfg.PartitionsCount)
	assert.Equal(t, 250, cfg.ChannelSize)
}
