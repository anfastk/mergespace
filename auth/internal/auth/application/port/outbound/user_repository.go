package outbound

import (
	"context"

/* 	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
 */)

type UserRepository interface {
	/* CreateUser(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error) */
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
