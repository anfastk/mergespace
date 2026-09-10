package main

import (
	"context"
	"log"
	"net/http"

	profilev1connect "github.com/anfastk/mergespace/contracts/gen/go/proto/profile/v1/profilev1connect"

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

	go func() {

		log.Println(
			"profile kafka consumer started",
		)

		if err := consumer.Run(
			context.Background(),
		); err != nil {

			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()

	path, handler := profilev1connect.NewProfileServiceHandler(
		app.GrpcHandler,
	)

	mux.Handle(
		path,
		handler,
	)

	log.Println(
		"profile grpc server started on :8082",
	)

	if err := http.ListenAndServe(
		":8082",
		mux,
	); err != nil {

		log.Fatal(err)
	}

}
