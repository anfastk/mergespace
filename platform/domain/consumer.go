package domain

import "context"

type Handler func(context.Context, Envelope) error

type Consumer interface {
	Run(ctx context.Context) error
	Close() error
}
