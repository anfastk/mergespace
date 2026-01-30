package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
)

type EventProducer interface {
	PublishSendOTP(ctx context.Context, event event.SendOTP) error
}
