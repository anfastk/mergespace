package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env not found:", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	return &Config{
		DatabaseURL: dsn,
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers:           strings.Split(getEnv("KAFKA_BROKERS", "localhost:29092"), ","),
			SchemaRegistryURL: getEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081"),
			UserSignupTopic:   getEnv("TOPIC_USER_AUTH", "auth.events"),
		},
	}, nil
}
