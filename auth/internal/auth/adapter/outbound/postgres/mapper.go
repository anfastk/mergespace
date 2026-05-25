package postgres

import (
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
)

func toDomainUser(id string, email string, username string, passwordHash string, status string) (*entity.User, error) {

	userID, err := valueobject.NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	usernameVO, err := valueobject.NewUsername(username)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		UserID:   userID,
		Email:    emailVO,
		Username: usernameVO,
		Password: &passwordHash,
		Status:   entity.UserStatus(status),
	}, nil
}
