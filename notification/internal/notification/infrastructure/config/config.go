package config

type Config struct {
	Kafka    KafkaConfig
	SendGrid SendGridConfig
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type SendGridConfig struct {
	APIKey string
	From   string
}
