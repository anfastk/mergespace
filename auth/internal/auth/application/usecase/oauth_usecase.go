package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
)

func (s *AuthService) GoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {

	googleUser, err := s.googleProvider.GetGoogleUser(
		ctx,
		code,
	)
	if err != nil {
		return nil, err
	}

	identity, err := s.authIdentityRepo.FindByProviderAndProviderUserID(
		ctx,
		entity.AuthProviderGoogle,
		googleUser.ID,
	)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	if identity != nil {

		user, err = s.userRepo.FindByID(
			ctx,
			identity.UserID,
		)
		if err != nil {
			return nil, err
		}

	} else {

		user, err = s.userRepo.FindByEmail(
			ctx,
			googleUser.Email,
		)

		if err != nil {

			if !errors.Is(err, errs.ErrUserNotFound) {
				return nil, err
			}

			userID, err := valueobject.NewUserID(
				s.idGen.NewID(ctx),
			)
			if err != nil {
				return nil, err
			}

			email, err := valueobject.NewEmail(
				googleUser.Email,
			)
			if err != nil {
				return nil, err
			}

			username, err := valueobject.NewUsername(
				strings.ToLower(
					strings.ReplaceAll(
						googleUser.Name,
						" ",
						"",
					),
				) + "_" + userID.String()[len(userID.String())-5:],
			)
			if err != nil {
				return nil, err
			}

			user = entity.NewOAuthUser(
				userID,
				email,
				username,
			)

			if err := s.userRepo.CreateUser(
				ctx,
				user,
			); err != nil {
				return nil, err
			}
		}

		identity = &entity.AuthIdentity{
			ID:             s.idGen.NewID(ctx),
			UserID:         user.UserID.String(),
			Provider:       entity.AuthProviderGoogle,
			ProviderUserID: googleUser.ID,
		}

		if err := s.authIdentityRepo.Create(
			ctx,
			identity,
		); err != nil {
			return nil, err
		}
	}

	accessToken, accessExpiry, err := s.tokenGenerator.GenerateAccessToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiry, err := s.tokenGenerator.GenerateRefreshToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	session := &entity.RefreshSession{
		ID: s.idGen.NewID(ctx),

		UserID: user.UserID.String(),

		RefreshToken: refreshToken,

		ExpiresAt: refreshExpiry,
	}

	if err := s.refreshSessionRepo.Create(
		ctx,
		session,
	); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserRes{
			ID:       user.UserID.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Status:   string(user.Status),
		},
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiry,
	}, nil
}

func (s *AuthService) GitHubLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {

	githubUser, err := s.githubProvider.GetGitHubUser(
		ctx,
		code,
	)
	if err != nil {
		return nil, err
	}

	identity, err := s.authIdentityRepo.FindByProviderAndProviderUserID(
		ctx,
		entity.AuthProviderGithub,
		githubUser.Login,
	)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	if identity != nil {

		user, err = s.userRepo.FindByID(
			ctx,
			identity.UserID,
		)
		if err != nil {
			return nil, err
		}

	} else {

		user, err = s.userRepo.FindByEmail(
			ctx,
			githubUser.Email,
		)

		if err != nil {

			if !errors.Is(err, errs.ErrUserNotFound) {
				return nil, err
			}

			userID, err := valueobject.NewUserID(
				s.idGen.NewID(ctx),
			)
			if err != nil {
				return nil, err
			}

			email, err := valueobject.NewEmail(
				githubUser.Email,
			)
			if err != nil {
				return nil, err
			}

			username, err := valueobject.NewUsername(
				strings.ToLower(
					githubUser.Login,
				),
			)
			if err != nil {
				return nil, err
			}

			user = entity.NewOAuthUser(
				userID,
				email,
				username,
			)

			if err := s.userRepo.CreateUser(
				ctx,
				user,
			); err != nil {
				return nil, err
			}
		}

		identity = &entity.AuthIdentity{
			ID:             s.idGen.NewID(ctx),
			UserID:         user.UserID.String(),
			Provider:       entity.AuthProviderGithub,
			ProviderUserID: githubUser.Login,
		}

		if err := s.authIdentityRepo.Create(
			ctx,
			identity,
		); err != nil {
			return nil, err
		}
	}

	accessToken, accessExpiry, err := s.tokenGenerator.GenerateAccessToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiry, err := s.tokenGenerator.GenerateRefreshToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	session := &entity.RefreshSession{
		ID: s.idGen.NewID(ctx),

		UserID: user.UserID.String(),

		RefreshToken: refreshToken,

		ExpiresAt: refreshExpiry,
	}

	if err := s.refreshSessionRepo.Create(
		ctx,
		session,
	); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserRes{
			ID:       user.UserID.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Status:   string(user.Status),
		},
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiry,
	}, nil
}
