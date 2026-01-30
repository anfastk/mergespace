package config

type KafkaConfig struct {
	Brokers           []string
	SchemaRegistryURL string
	UserSignupTopic   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	DatabaseURL string
	Kafka       KafkaConfig
	Redis       RedisConfig
}
