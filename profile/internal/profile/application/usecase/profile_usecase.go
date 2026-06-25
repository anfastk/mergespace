package usecase

import (
	"context"
	"errors"
	"strings"
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

func (u *ProfileUseCase) UpdateProfile(ctx context.Context, req *dto.UpdateProfileRequest) error {

	if req.FirstName != nil {

		trimmed := strings.TrimSpace(
			*req.FirstName,
		)

		if trimmed == "" {
			return errors.New(
				"first name cannot be empty",
			)
		}

		*req.FirstName = trimmed
	}

	if req.LastName != nil {

		trimmed := strings.TrimSpace(
			*req.LastName,
		)

		if trimmed == "" {
			return errors.New(
				"last name cannot be empty",
			)
		}

		*req.LastName = trimmed
	}

	if req.Bio != nil {

		trimmed := strings.TrimSpace(
			*req.Bio,
		)

		if len(trimmed) > 500 {
			return errors.New(
				"bio too long",
			)
		}

		*req.Bio = trimmed
	}

	return u.repo.UpdateProfile(
		ctx,
		req,
	)

}
