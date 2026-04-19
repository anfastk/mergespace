package dto

type AuthResponse struct {
	User            UserRes
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt int64
}
