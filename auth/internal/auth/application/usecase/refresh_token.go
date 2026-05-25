package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {

	if req.RefreshToken == "" {
		return nil, errors.New("refresh token required")
	}

	session, err := s.refreshSessionRepo.FindByToken(
		ctx,
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	if session.Revoked {
		return nil, errors.New("session revoked")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	userID, err := s.tokenGenerator.ValidateRefreshToken(
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	accessToken, accessExpiry, err :=
		s.tokenGenerator.GenerateAccessToken(
			userID,
		)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:     accessToken,
		AccessExpiresAt: accessExpiry,
	}, nil
}
