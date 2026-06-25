package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/anfastk/mergespace/profile/internal/profile/application/dto"
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

func (r *Repository) UpdateProfile(
	ctx context.Context,
	req *dto.UpdateProfileRequest,
) error {
	sets := []string{}
	args := []any{}
	argPos := 1

	add := func(column string, value *string) {
		if value == nil {
			return
		}

		sets = append(
			sets,
			fmt.Sprintf("%s = $%d", column, argPos),
		)

		args = append(args, *value)
		argPos++
	}

	add("first_name", req.FirstName)
	add("last_name", req.LastName)
	add("bio", req.Bio)
	add("avatar_url", req.AvatarURL)

	sets = append(sets, "updated_at = NOW()")

	query := fmt.Sprintf(`
UPDATE profiles
SET %s
WHERE user_id = $%d
`,
		strings.Join(sets, ", "),
		argPos,
	)

	args = append(args, req.UserID)

	result, err := r.db.Exec(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"profile not found",
		)
	}

	return nil
}	
