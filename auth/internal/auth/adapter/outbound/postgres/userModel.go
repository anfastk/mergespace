package postgres

import "time"

type UserModel struct {
	UserID       string    `db:"user_id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash *string   `db:"password_hash"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
