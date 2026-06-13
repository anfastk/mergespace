package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/anfastk/mergespace/profile/internal/profile/application/dto"
	"github.com/anfastk/mergespace/profile/internal/profile/application/port/outbound"
	"github.com/anfastk/mergespace/profile/internal/profile/domain/entity"
)

type ProfileUseCase struct {
	repo outbound.ProfileRepository
}

func NewProfileUseCase(repo outbound.ProfileRepository) *ProfileUseCase {
	
	return &ProfileUseCase{
		repo: repo,
	}
}

func (u *ProfileUseCase) CreateProfile(ctx context.Context, event *dto.UserCreatedEvent) error {

	profile := &entity.Profile{
		ID:        uuid.NewString(),
		UserID:    event.UserID,
		Email:     event.Email,
		Username:  event.Username,
		FirstName: event.FirstName,
		LastName:  event.LastName,
		Bio:       "",
		AvatarURL: "",
		CreatedAt: time.Now(),
	}

	return u.repo.Create(
		ctx,
		profile,
	)
}
