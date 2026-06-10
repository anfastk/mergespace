package outbound

import (
	"context"
)

type EventProducer interface {
	Publish(ctx context.Context, eventName string, key []byte, payload any) error
}
