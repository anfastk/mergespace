package entity

import (
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderGithub AuthProvider = "github"
)

type User struct {
	UserID              valueobject.UserID
	Username            valueobject.Username
	Email               valueobject.Email
	Password            *valueobject.Password
	Status              UserStatus
	DeletionScheduledAt *time.Time
}

func NewLocalUser(id valueobject.UserID, email valueobject.Email, username valueobject.Username, password valueobject.Password) *User {
	return &User{
		UserID:   id,
		Email:    email,
		Username: username,
		Password: &password,
		Status:   UserStatusActive,
	}
}
