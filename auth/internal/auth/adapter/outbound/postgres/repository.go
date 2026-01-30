package postgres

import (
	"context"
	"fmt"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

var _ outbound.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) outbound.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}
