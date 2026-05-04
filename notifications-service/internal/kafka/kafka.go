package kafkaSettings

type Config struct {
	Brokers          string `env:"BROKERS" env-default:"kafka1:9092,kafka2:9094,kafka3:9096"`
	SecurityProtocol string `env:"KAFKA_SECURITY_PROTOCOL" env-default:"plaintext"`
	CaCertPath       string `env:"CA_CERT_PATH"`
	ClientCertPath   string `env:"CLIENT_CERT_PATH"`
	KeyFilePath      string `env:"CLIENT_KEY_PATH"`
	ChannelSize      int    `env:"CHANNEL_SIZE" env-default:"1000"`
	ChannelNumber    int    `env:"CHANNEL_NUMBER" env-default:"3"`
	BatchSize        int    `env:"BATCH_SIZE" env-default:"20"`
}
