package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type AuthIdentityRepository interface {
	FindByProviderAndProviderUserID(ctx context.Context, provider entity.AuthProvider, providerUserID string) (*entity.AuthIdentity, error)
	Create(ctx context.Context, identity *entity.AuthIdentity) error
}
