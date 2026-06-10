package main

import (
	"context"
	"log"

	"github.com/anfastk/mergespace/notification/internal/notification/infrastructure/di"

	platformAvro "github.com/anfastk/mergespace/platform/infrastructure/avro"

	platformKafka "github.com/anfastk/mergespace/platform/infrastructure/kafka"
)

func main() {

	app := di.BuildApp()

	registry := platformAvro.NewRegistry(
		"http://localhost:8081",
	)

	codec := platformAvro.NewCodec(
		registry,
	)

	log.Println(app.Config.Kafka.Topic)

	consumer, err := platformKafka.NewConsumer(
		app.Config.Kafka.Brokers,
		app.Config.Kafka.GroupID,
		[]string{
			app.Config.Kafka.Topic,
		},
		codec,
		app.ConsumerHandler.Handle,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		"notification service started",
	)

	if err := consumer.Run(
		context.Background(),
	); err != nil {
		log.Fatal(err)
	}
}