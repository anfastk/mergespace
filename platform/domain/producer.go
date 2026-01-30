package domain

import "context"

type Producer interface {
	Publish(ctx context.Context, eventName string, key []byte, event any) error
	Close() error
}
