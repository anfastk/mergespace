package outbound

import "context"

type IDGenerator interface {
	NewID(ctx context.Context) string
}
