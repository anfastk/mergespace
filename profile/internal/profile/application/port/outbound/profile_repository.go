package outbound

import (
	"context"

	"github.com/anfastk/mergespace/profile/internal/profile/application/dto"
	"github.com/anfastk/mergespace/profile/internal/profile/domain/entity"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *entity.Profile) error
	UpdateProfile(ctx context.Context, req *dto.UpdateProfileRequest) error
}
