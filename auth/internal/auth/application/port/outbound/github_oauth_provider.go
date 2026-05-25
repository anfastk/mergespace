package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

type GitHubOAuthProvider interface {
	GetGitHubUser(ctx context.Context, code string) (*dto.GitHubUser, error)
}
