package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

type GoogleOAuthProvider interface {
	GetGoogleUser(ctx context.Context, code string) (*dto.GoogleUser, error)
}
