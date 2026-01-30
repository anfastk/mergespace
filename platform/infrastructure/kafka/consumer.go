package kafka

import (
	"context"
	"log"

	"github.com/anfastk/mergespace/platform/domain"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
	codec  domain.Codec
	handle domain.Handler
}

func NewConsumer(
	brokers []string,
	group string,
	topics []string,
	codec domain.Codec,
	handler domain.Handler,
) (*Consumer, error) {

	client, err := kgo.NewClient(
		ConsumerOpts(brokers, group, topics)...,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client: client,
		codec:  codec,
		handle: handler,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		fetches.EachRecord(func(r *kgo.Record) {

			var eventName string
			for _, h := range r.Headers {
				if h.Key == "event_name" {
					eventName = string(h.Value)
					break
				}
			}

			data, err := c.codec.Decode(r.Value)
			if err != nil {
				log.Println("decode error:", err)
				return
			}

			if err := c.handle(ctx, domain.Envelope{
				Name:    eventName,
				Payload: data,
			}); err != nil {
				log.Println("handler error:", err)
			}
		})

		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			log.Println("commit error:", err)
		}
	}
}

func (c *Consumer) Close() error {
	c.client.Close()
	return nil
}
