package token

import (
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	accessSecret  string
	refreshSecret string
}

var _ outbound.TokenGenerator = (*JWTGenerator)(nil)

func NewJWTGenerator(accessSecret, refreshSecret string) *JWTGenerator {
	return &JWTGenerator{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (j *JWTGenerator) GenerateAccessToken(userID string) (string, time.Time, error) {

	expiry := time.Now().Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiry.Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.accessSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiry, nil
}

func (j *JWTGenerator) GenerateRefreshToken(userID string) (string, error) {

	expiry := time.Now().Add(7 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiry.Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.refreshSecret))
}
