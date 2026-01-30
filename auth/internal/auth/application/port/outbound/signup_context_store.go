package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type SignupContextStore interface {
	Set(ctx context.Context, signup *entity.SignupContext) error
/* 	Get(ctx context.Context, id entity.SignupContextID) (*entity.SignupContext, error)
	Delete(ctx context.Context, id entity.SignupContextID) error */
}
