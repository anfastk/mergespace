package valueobject

import (
	"errors"
)

type UserID struct {
	value string
}

func NewUserID(value string) (UserID, error) {
	if value == "" {
		return UserID{}, errors.New("empty id")
	}
	return UserID{value: value}, nil
}

func (id UserID) String() string {
	return id.value
}
