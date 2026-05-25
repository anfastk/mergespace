package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		user.UserID,
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {

	query := `
		SELECT
			id,
			email,
			username,
			password_hash,
			status
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var (
		id           string
		emailStr     string
		usernameStr  string
		status       string
		passwordHash sql.NullString
	)

	err := row.Scan(
		&id,
		&emailStr,
		&usernameStr,
		&passwordHash,
		&status,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf(
			"failed to find user by email: %w",
			err,
		)
	}

	return toDomainUser(
		id,
		emailStr,
		usernameStr,
		passwordHash.String,
		status,
	)
}

func (r *UserRepository) FindByEmailOrUsername(ctx context.Context, value string) (*entity.User, error) {

	query := `
	SELECT id, email, username, password_hash, status
	FROM users
	WHERE email = $1 OR username = $1
	LIMIT 1
	`

	row := r.db.QueryRow(ctx, query, value)

	var (
		id           string
		email        string
		username     string
		passwordHash *string
		status       string
	)

	err := row.Scan(
		&id,
		&email,
		&username,
		&passwordHash,
		&status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	userID, err := valueobject.NewUserID(id)
	if err != nil {
		return nil, err
	}

	userEmail, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	userUsername, err := valueobject.NewUsername(username)
	if err != nil {
		return nil, err
	}

	user := entity.NewLocalUser(
		userID,
		userEmail,
		userUsername,
		passwordHash,
	)

	user.Status = entity.UserStatus(status)

	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, hash string) error {

	query := `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, hash, userID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*entity.User, error) {

	query := `
		SELECT
			id,
			email,
			username,
			password_hash,
			status
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, userID)

	var (
		id           string
		emailStr     string
		usernameStr  string
		status       string
		passwordHash sql.NullString
	)

	err := row.Scan(
		&id,
		&emailStr,
		&usernameStr,
		&passwordHash,
		&status,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf(
			"failed to find user by id: %w",
			err,
		)
	}

	return toDomainUser(
		id,
		emailStr,
		usernameStr,
		passwordHash.String,
		status,
	)
}
