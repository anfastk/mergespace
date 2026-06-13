package main

import (
	"context"
	"log"

	"github.com/anfastk/mergespace/profile/internal/profile/infrastructure/di"

	platformAvro "github.com/anfastk/mergespace/platform/infrastructure/avro"

	platformKafka "github.com/anfastk/mergespace/platform/infrastructure/kafka"
)

func main() {

	app := di.BuildApp()

	registry := platformAvro.NewRegistry(
		app.Config.Kafka.SchemaRegistryURL,
	)

	codec := platformAvro.NewCodec(
		registry,
	)

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
		"profile service started",
	)

	if err := consumer.Run(
		context.Background(),
	); err != nil {

		log.Fatal(err)
	}

}
