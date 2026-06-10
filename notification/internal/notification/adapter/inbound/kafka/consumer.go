package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/anfastk/mergespace/notification/internal/notification/application/dto"
	"github.com/anfastk/mergespace/notification/internal/notification/application/usecase"

	platform "github.com/anfastk/mergespace/platform/domain"
	"github.com/anfastk/mergespace/platform/events"
	platformEvents "github.com/anfastk/mergespace/platform/events"
)

type ConsumerHandler struct {
	usecase *usecase.NotificationUseCase
}

func NewConsumerHandler(usecase *usecase.NotificationUseCase) *ConsumerHandler {

	return &ConsumerHandler{
		usecase: usecase,
	}
}

func (h *ConsumerHandler) Handle(ctx context.Context, event platform.Envelope) error {

	log.Println(
		"EVENT RECEIVED:",
		event.Name,
	)

	switch event.Name {

	case platformEvents.EventSendOTP:

		log.Println(
			"PROCESSING SEND OTP EVENT",
		)

		var payload dto.SendOTPEvent

		bytes, err := json.Marshal(
			event.Payload,
		)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(
			bytes,
			&payload,
		); err != nil {
			return err
		}

		log.Println(
			"SENDING OTP EMAIL TO:",
			payload.Email,
		)

		return h.usecase.HandleSendOTP(
			ctx,
			&payload,
		)

	case platformEvents.EventUserCreated:

		log.Println(
			"PROCESSING USER CREATED EVENT",
		)

		var payload dto.UserCreatedEvent

		bytes, err := json.Marshal(
			event.Payload,
		)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(
			bytes,
			&payload,
		); err != nil {
			return err
		}

		log.Println(
			"SENDING WELCOME EMAIL TO:",
			payload.Email,
		)

		return h.usecase.HandleUserCreated(
			ctx,
			&payload,
		)

	case events.EventForgotPasswordOTP:

		log.Println(
			"PROCESSING FORGOT PASSWORD OTP EVENT",
		)

		var payload dto.ForgotPasswordOTPEvent

		bytes, err := json.Marshal(
			event.Payload,
		)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(
			bytes,
			&payload,
		); err != nil {
			return err
		}

		log.Println(
			"SENDING FORGOT PASSWORD OTP EMAIL TO:",
			payload.Email,
		)

		return h.usecase.HandleForgotPasswordOTP(
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
