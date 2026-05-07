package inbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

type AuthUseCase interface {
	CheckUsernameAvailability(ctx context.Context, req *dto.CheckUsernameReq) (*dto.CheckUsernameRes, error)
	InitiateSignup(ctx context.Context, req *dto.InitiateSignUpRequest) (*dto.InitiateSignUpResponce, error)
	VerifySignup(ctx context.Context, req *dto.VerifySignupRequest) (*dto.AuthResponse, error)
	ResendOTP(ctx context.Context, req *dto.ResendOTPRequest) (*dto.InitiateSignUpResponce, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
}
