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

type User struct {
	UserID              valueobject.UserID
	Username            valueobject.Username
	Email               valueobject.Email
	Password            *string
	Status              UserStatus
	DeletionScheduledAt *time.Time
}

func NewLocalUser(id valueobject.UserID, email valueobject.Email, username valueobject.Username, password *string) *User {
	return &User{
		UserID:   id,
		Email:    email,
		Username: username,
		Password: password,
		Status:   UserStatusActive,
	}
}

func NewOAuthUser(id valueobject.UserID, email valueobject.Email, username valueobject.Username) *User {
	return &User{
		UserID:   id,
		Email:    email,
		Username: username,
		Password: nil,
		Status:   UserStatusActive,
	}
}
