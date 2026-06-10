package outbound

import (
	"context"
)

type EmailSender interface {
	SendOTP(ctx context.Context, email string, otp string) error
	SendWelcome(ctx context.Context, email string, firstName string) error
	SendForgotPasswordOTP(ctx context.Context, email string, otp string) error
}
