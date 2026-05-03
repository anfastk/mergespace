package outbound

import "time"

type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, time.Time, error)
	GenerateRefreshToken(userID string) (string, error)
}