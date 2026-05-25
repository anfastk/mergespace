package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type RefreshSessionRepository interface {
	Create(ctx context.Context, session *entity.RefreshSession) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshSession, error)
	Revoke(ctx context.Context, token string) error
}
