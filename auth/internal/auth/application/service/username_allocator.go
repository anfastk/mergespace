package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

const DefaultMaxUsernameAttempts = 10

var usernameNormalizeRegex = regexp.MustCompile(`[^a-z0-9]`)

var _ outbound.UsernameAllocator = (*UsernameAllocatorService)(nil)

type UsernameAllocatorService  struct {
	userRepo    outbound.UserRepository
	maxAttempts int
}


func NewUsernameAllocator(userRepo outbound.UserRepository, maxAttempts int) *UsernameAllocatorService  {
	if maxAttempts <= 0 {
		maxAttempts = DefaultMaxUsernameAttempts
	}

	return &UsernameAllocatorService {
		userRepo:    userRepo,
		maxAttempts: maxAttempts,
	}
}

func (a *UsernameAllocatorService) Allocate(ctx context.Context, firstName string, lastName string) (string, error) {

	base := normalizeUsername(firstName) + "." + normalizeUsername(lastName)

	for attempt := 0; attempt < a.maxAttempts; attempt++ {
		username := base
		if attempt > 0 {
			username = fmt.Sprintf("%s%d", base, attempt)
		}

		exists, err := a.userRepo.ExistsByUsername(ctx, username)
		if err != nil {
			return "", fmt.Errorf("username availability check failed: %w", err)
		}
		if !exists {
			return username, nil
		}
	}

	return "", errs.ErrUsernameGenerationFailed
}

func normalizeUsername(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	return usernameNormalizeRegex.ReplaceAllString(value, "")
}
