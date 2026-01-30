package dto

type AuthResponse struct {
	AccessToken     string
	AccessExpiresAt int64
	User            UserRes
}
