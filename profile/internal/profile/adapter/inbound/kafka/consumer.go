package kafka

import (
	"context"
	"fmt"
	"log"

	platform "github.com/anfastk/mergespace/platform/domain"
	platformEvents "github.com/anfastk/mergespace/platform/events"

	"github.com/anfastk/mergespace/profile/internal/profile/application/dto"
	"github.com/anfastk/mergespace/profile/internal/profile/application/usecase"
)

type ConsumerHandler struct {
	usecase *usecase.ProfileUseCase
}

func NewConsumerHandler(usecase *usecase.ProfileUseCase) *ConsumerHandler {

	return &ConsumerHandler{
		usecase: usecase,
	}

}

func (h *ConsumerHandler) Handle(ctx context.Context, event platform.Envelope) error {

	switch event.Name {

	case platformEvents.EventUserCreated:

		log.Println(
			"PROCESSING USER CREATED EVENT",
		)

		payloadMap := event.Payload

		payload := dto.UserCreatedEvent{
			UserID: fmt.Sprintf(
				"%v",
				payloadMap["user_id"],
			),

			Email: fmt.Sprintf(
				"%v",
				payloadMap["email"],
			),

			Username: fmt.Sprintf(
				"%v",
				payloadMap["username"],
			),

			FirstName: fmt.Sprintf(
				"%v",
				payloadMap["first_name"],
			),

			LastName: fmt.Sprintf(
				"%v",
				payloadMap["last_name"],
			),
		}

		log.Printf(
			"PAYLOAD DEBUG: %+v",
			payload,
		)

		return h.usecase.CreateProfile(
			ctx,
			&payload,
		)
	}

	log.Println(
		"UNKNOWN EVENT:",
		event.Name,
	)

	return nil
}
