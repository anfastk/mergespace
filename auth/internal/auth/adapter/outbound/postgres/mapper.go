package postgres
/* 
import (
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

func toModel(u *entity.User) *UserModel {
	var passwordHash *string
	if u.Password != nil {
		hash := u.Password.String()
		passwordHash = &hash
	}

	return &UserModel{
		UserID:       u.UserID.String().String(),
		Username:     u.Username.String(),
		Email:        u.Email.String(),
		PasswordHash: passwordHash,
		Status:       string(u.Status),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
} */

/* func toEntity(m *UserModel) (*entity.User, error) {
	userID, err := valueobject.NewUserID(m.ID)
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	username, err := valueobject.NewUsername(m.Username)
	if err != nil {
		return nil, err
	}

	firstName, _ := valueobject.NewName(m.FirstName)
	lastName, _ := valueobject.NewName(m.LastName)

	var password *valueobject.Password
	if m.PasswordHash != nil {
		password = valueobject.NewPassword(*m.PasswordHash)
	}

	return &entity.User{
		UserID:       userID,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Username:     username,
		AuthProvider: entity.AuthProvider(m.AuthProvider),
		Password:     password,
		ProviderID:   m.providerID,
		Status:       entity.UserStatus(m.Status),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}
*/
