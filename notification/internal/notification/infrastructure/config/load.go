package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() *Config {

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env not found:", err)
	}

	return &Config{
		Kafka: KafkaConfig{
			Brokers: []string{
				os.Getenv("KAFKA_BROKER"),
			},
			Topic:   os.Getenv("KAFKA_TOPIC"),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
		},

		SendGrid: SendGridConfig{
			APIKey: os.Getenv("SENDGRID_API_KEY"),
			From:   os.Getenv("EMAIL_FROM"),
		},
	}
}
