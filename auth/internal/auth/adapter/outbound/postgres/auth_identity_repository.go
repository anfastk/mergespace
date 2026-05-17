package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type AuthIdentityRepository struct {
	db *pgxpool.Pool
}

func NewAuthIdentityRepository(db *pgxpool.Pool) *AuthIdentityRepository {
	return &AuthIdentityRepository{db: db}
}

func (r *AuthIdentityRepository) FindByProviderAndProviderUserID(ctx context.Context, provider entity.AuthProvider, providerUserID string) (*entity.AuthIdentity, error) {

	query := `
	SELECT id, user_id, provider, provider_user_id
	FROM auth_identities
	WHERE provider = $1 AND provider_user_id = $2
	LIMIT 1
	`

	row := r.db.QueryRow(ctx, query, provider, providerUserID)

	var identity entity.AuthIdentity

	err := row.Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderUserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &identity, nil
}

func (r *AuthIdentityRepository) Create(ctx context.Context, identity *entity.AuthIdentity) error {

	query := `
	INSERT INTO auth_identities (
		id,
		user_id,
		provider,
		provider_user_id
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.ProviderUserID,
	)

	return err
}
