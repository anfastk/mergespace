package entity

import "time"

type Profile struct {
	ID        string
	UserID    string
	Email     string
	Username  string
	FirstName string
	LastName  string
	Bio       string
	AvatarURL string
	CreatedAt time.Time
}
