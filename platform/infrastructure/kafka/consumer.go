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

func NewConsumer(brokers []string, group string, topics []string, codec domain.Codec, handler domain.Handler) (*Consumer, error) {

	client, err := kgo.NewClient(
		ConsumerOpts(
			brokers,
			group,
			topics,
		)...,
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

		select {

		case <-ctx.Done():
			log.Println("consumer shutting down")
			return ctx.Err()

		default:
		}

		fetches := c.client.PollFetches(ctx)

		if fetches.IsClientClosed() {
			return nil
		}

		if errs := fetches.Errors(); len(errs) > 0 {

			for _, err := range errs {
				log.Printf(
					"kafka fetch error: %v",
					err,
				)
			}

			continue
		}

		fetches.EachRecord(func(r *kgo.Record) {

			var (
				eventName string
				eventID   string
			)

			for _, h := range r.Headers {

				switch h.Key {

				case "event_name":
					eventName = string(h.Value)

				case "event_id":
					eventID = string(h.Value)
				}
			}

			data, err := c.codec.Decode(
				r.Value,
			)
			if err != nil {

				log.Printf(
					"event=%s event_id=%s decode_error=%v",
					eventName,
					eventID,
					err,
				)

				return
			}

			envelope := domain.Envelope{
				ID:      eventID,
				Name:    eventName,
				Payload: data,
			}

			err = retry(
				ctx,
				3,
				func() error {

					return c.handle(
						ctx,
						envelope,
					)
				},
			)

			if err != nil {

				log.Printf(
					"event=%s event_id=%s handler_error=%v",
					eventName,
					eventID,
					err,
				)

				return
			}

			if err := c.client.CommitRecords(
				ctx,
				r,
			); err != nil {

				log.Printf(
					"event=%s event_id=%s commit_error=%v",
					eventName,
					eventID,
					err,
				)
			}

			log.Printf(
				"event=%s event_id=%s processed=true",
				eventName,
				eventID,
			)
		})
	}
}

func (c *Consumer) Close() error {
	c.client.Close()
	return nil
}
