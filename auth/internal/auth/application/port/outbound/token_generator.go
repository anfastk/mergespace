package outbound

import "time"

type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, time.Time, error)
	GenerateRefreshToken(userID string) (string, time.Time, error)
	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, error)
}
