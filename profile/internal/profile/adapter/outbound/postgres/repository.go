package postgres

import (
	"context"

	"github.com/anfastk/mergespace/profile/internal/profile/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {

	return &Repository{
		db: db,
	}

}

func (r *Repository) Create(ctx context.Context, profile *entity.Profile) error {

	query := `
INSERT INTO profiles (
	id,
	user_id,
	email,
	username,
	first_name,
	last_name,
	bio,
	avatar_url,
	created_at
)
VALUES (
	$1,$2,$3,$4,$5,$6,$7,$8,$9
)
`

	_, err := r.db.Exec(
		ctx,
		query,
		profile.ID,
		profile.UserID,
		profile.Email,
		profile.Username,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.AvatarURL,
		profile.CreatedAt,
	)
	return err
}
