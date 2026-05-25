package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type RefreshSessionRepo struct {
	db *pgxpool.Pool
}

func NewRefreshSessionRepo(db *pgxpool.Pool) *RefreshSessionRepo {

	return &RefreshSessionRepo{
		db: db,
	}
}

func (r *RefreshSessionRepo) Create(ctx context.Context, session *entity.RefreshSession) error {

	query := `
		INSERT INTO refresh_sessions (
			id,
			user_id,
			refresh_token,
			user_agent,
			ip_address,
			expires_at
		)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
	)

	return err
}

func (r *RefreshSessionRepo) FindByToken(ctx context.Context, token string) (*entity.RefreshSession, error) {

	query := `
		SELECT
			id,
			user_id,
			refresh_token,
			user_agent,
			ip_address,
			expires_at,
			revoked,
			created_at
		FROM refresh_sessions
		WHERE refresh_token = $1
	`

	row := r.db.QueryRow(ctx, query, token)

	var session entity.RefreshSession

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.Revoked,
		&session.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("session not found")
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *RefreshSessionRepo) Revoke(ctx context.Context, token string) error {

	query := `
		UPDATE refresh_sessions
		SET revoked = true
		WHERE refresh_token = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		token,
	)

	return err
}
