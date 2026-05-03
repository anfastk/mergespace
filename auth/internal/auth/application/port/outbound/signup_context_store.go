package outbound

import (
	"context"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
)

type SignupContextStore interface {
	Save(ctx context.Context, signup *entity.SignupContext) error
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.SignupContext, error)
	FindByID(ctx context.Context, id entity.SignupContextID) (*entity.SignupContext, error)
	Delete(ctx context.Context, id entity.SignupContextID) error

	AcquireSignupSlot(ctx context.Context, email valueobject.Email) (bool, error)
	ReleaseSignupSlot(ctx context.Context, email valueobject.Email) error

	GetAttempts(ctx context.Context, id entity.SignupContextID) (int, error)
	IncrementAttempts(ctx context.Context, id entity.SignupContextID, ttl time.Duration) error
	DeleteAttempts(ctx context.Context, id entity.SignupContextID) error
}
