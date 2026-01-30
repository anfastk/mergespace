package outbound

import "context"

type UsernameAllocator  interface {
	Allocate(ctx context.Context, firstName string, lastName string) (string, error)
}
