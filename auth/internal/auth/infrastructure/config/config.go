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

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Config struct {
	DatabaseURL string
	Kafka       KafkaConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Google      GoogleConfig
}
