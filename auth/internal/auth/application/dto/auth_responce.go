package dto

import "time"

type AuthResponse struct {
	User            UserRes
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
}
