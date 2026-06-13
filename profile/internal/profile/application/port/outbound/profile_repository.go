package outbound

import (
	"context"

	"github.com/anfastk/mergespace/profile/internal/profile/domain/entity"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *entity.Profile) error
}
