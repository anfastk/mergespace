package kafka

import (
	"context"
	"time"

	"github.com/anfastk/mergespace/platform/domain"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	codec  domain.Codec
	topic  string
}

func NewProducer(
	brokers []string,
	topic string,
	codec domain.Codec,
) (*Producer, error) {

	client, err := kgo.NewClient(ProducerOpts(brokers)...)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: client,
		codec:  codec,
		topic:  topic,
	}, nil
}

func (p *Producer) Publish(
	ctx context.Context,
	eventName string,
	key []byte,
	event any,
) error {

	value, err := p.codec.Encode(eventName, event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return p.client.ProduceSync(ctx, &kgo.Record{
		Topic: p.topic,
		Key:   key,
		Value: value,
		Headers: []kgo.RecordHeader{
			{Key: "event_name", Value: []byte(eventName)},
		},
	}).FirstErr()
}

func (p *Producer) Close() error {
	p.client.Close()
	return nil
}
