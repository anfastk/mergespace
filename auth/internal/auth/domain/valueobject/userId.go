package valueobject

import "github.com/google/uuid"

type UserID struct {
	value string
}

func NewUserID(value string) (UserID, error) {
	if _, err := uuid.Parse(value); err != nil {
		return UserID{}, err
	}
	return UserID{value: value}, nil
}

func (id UserID) String() string {
	return id.value
}
