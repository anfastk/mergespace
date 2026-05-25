package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	accessSecret  string
	refreshSecret string
}

func NewJWTGenerator(accessSecret, refreshSecret string) *JWTGenerator {

	return &JWTGenerator{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (j *JWTGenerator) GenerateAccessToken(userID string) (string, time.Time, error) {

	expiry := time.Now().Add(
		15 * time.Minute,
	)

	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  expiry.Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signed, err := token.SignedString(
		[]byte(j.accessSecret),
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiry, nil
}

func (j *JWTGenerator) GenerateRefreshToken(userID string) (string, time.Time, error) {

	expiry := time.Now().Add(
		7 * 24 * time.Hour,
	)

	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  expiry.Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signed, err := token.SignedString(
		[]byte(j.refreshSecret),
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiry, nil
}

func (j *JWTGenerator) ValidateAccessToken(tokenString string) (string, error) {

	if tokenString == "" {
		return "", errors.New("empty access token")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.accessSecret), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid access token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject")
	}

	return userID, nil
}

func (j *JWTGenerator) ValidateRefreshToken(tokenString string) (string, error) {

	// 🔥 IMPORTANT SECURITY FIX
	if tokenString == "" {
		return "", errors.New("empty refresh token")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.refreshSecret), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject")
	}

	return userID, nil
}
