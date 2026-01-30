package inbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

type AuthUseCase interface {
	InitiateSignup(ctx context.Context, req *dto.InitiateSignUpRequest) (*dto.InitiateSignUpResponce, error)
}
