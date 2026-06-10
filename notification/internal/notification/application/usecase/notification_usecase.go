package usecase

import (
	"context"

	"github.com/anfastk/mergespace/notification/internal/notification/application/dto"
	"github.com/anfastk/mergespace/notification/internal/notification/application/port/outbound"
)

type NotificationUseCase struct {
	emailSender outbound.EmailSender
}

func NewNotificationUseCase(emailSender outbound.EmailSender) *NotificationUseCase {

	return &NotificationUseCase{
		emailSender: emailSender,
	}
}

func (u *NotificationUseCase) HandleSendOTP(ctx context.Context, event *dto.SendOTPEvent) error {

	return u.emailSender.SendOTP(
		ctx,
		event.Email,
		event.OTP,
	)
}

func (u *NotificationUseCase) HandleUserCreated(ctx context.Context, event *dto.UserCreatedEvent) error {

	return u.emailSender.SendWelcome(
		ctx,
		event.Email,
		event.FirstName,
	)
}

func (u *NotificationUseCase) HandleForgotPasswordOTP(ctx context.Context, event *dto.ForgotPasswordOTPEvent) error {

	return u.emailSender.SendForgotPasswordOTP(
		ctx,
		event.Email,
		event.OTP,
	)
}
