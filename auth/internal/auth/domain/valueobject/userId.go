package valueobject

import "github.com/google/uuid"

type UserID struct {
	value uuid.UUID
}

func NewUserID() UserID {
	return UserID{value: uuid.New()}
}

func (id UserID) String() uuid.UUID {
	return id.value
}
