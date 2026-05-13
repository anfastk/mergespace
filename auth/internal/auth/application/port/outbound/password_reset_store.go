package outbound

import (
	"context"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type PasswordResetStore interface {
	Save(ctx context.Context, reset *entity.PasswordResetContext) error
	FindByID(ctx context.Context, id entity.PasswordResetContextID) (*entity.PasswordResetContext, error)
	Delete(ctx context.Context, id entity.PasswordResetContextID) error
	Update(ctx context.Context, reset *entity.PasswordResetContext) error

	GetAttempts(ctx context.Context, id entity.PasswordResetContextID) (int, error)
	IncrementAttempts(ctx context.Context, id entity.PasswordResetContextID, ttl time.Duration) error
	DeleteAttempts(ctx context.Context, id entity.PasswordResetContextID) error

	SetLastOTPSentAt(ctx context.Context, id entity.PasswordResetContextID, t time.Time) error
	GetLastOTPSentAt(ctx context.Context, id entity.PasswordResetContextID) (time.Time, error)
}
