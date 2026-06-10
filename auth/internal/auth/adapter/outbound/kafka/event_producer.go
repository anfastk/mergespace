package kafka

import (
	"context"

	platformKafka "github.com/anfastk/mergespace/platform/infrastructure/kafka"
)

type EventProducer struct {
	producer *platformKafka.Producer
}

func NewEventProducer(producer *platformKafka.Producer) *EventProducer {

	return &EventProducer{
		producer: producer,
	}
}

func (p *EventProducer) Publish(ctx context.Context, eventName string, key []byte, payload any) error {

	return p.producer.Publish(
		ctx,
		eventName,
		key,
		payload,
	)
}
