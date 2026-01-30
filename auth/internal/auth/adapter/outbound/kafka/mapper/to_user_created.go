package mapper

import (
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/kafka/avro"
	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
)

func ToUserCreatedAvro(e event.UserCreated) avro.UserCreatedAvro {
	return avro.UserCreatedAvro{
		UserID:    e.UserID,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Email:     e.Email,
		Avatar:    e.Avatar,
	}
}
