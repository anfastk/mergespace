package entity

import "time"

type RefreshSession struct {
	ID           string
	UserID       string
	RefreshToken string
	UserAgent    string
	IPAddress    string
	Revoked      bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
