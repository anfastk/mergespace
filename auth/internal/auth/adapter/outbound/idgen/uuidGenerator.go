package idgen

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/google/uuid"
)

type UUIDGenerator struct{}

var _ outbound.IDGenerator = (*UUIDGenerator)(nil)

func NewUUIDGenerator() outbound.IDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) NewID(ctx context.Context) string {
	return uuid.NewString()
}
