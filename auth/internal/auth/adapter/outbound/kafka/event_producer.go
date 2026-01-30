package kafka

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/kafka/mapper"
	authEvent "github.com/anfastk/mergespace/auth/internal/auth/application/event"
	"github.com/anfastk/mergespace/platform/domain"
)

type EventProducer struct {
	producer domain.Producer
}

func NewEventProducer(p domain.Producer) *EventProducer {
	return &EventProducer{producer: p}
}

func (e *EventProducer) PublishSendOTP(ctx context.Context, ev authEvent.SendOTP) error {

	return e.producer.Publish(ctx, "auth.send_otp", []byte(ev.Email), mapper.ToSendOTPAvro(ev))
}
