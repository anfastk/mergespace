package config

type Config struct {
	DatabaseURL string
	Kafka       KafkaConfig
}

type KafkaConfig struct {
	Brokers           []string
	Topic             string
	GroupID           string
	SchemaRegistryURL string
}
