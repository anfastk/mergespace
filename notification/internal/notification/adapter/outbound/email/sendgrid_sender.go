package email

import (
	"context"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridSender struct {
	apiKey string
	from   string
}

func NewSendGridSender(apiKey string, from string) *SendGridSender {

	return &SendGridSender{
		apiKey: apiKey,
		from:   from,
	}
}

func (s *SendGridSender) send(to string, subject string, body string) error {

	from := mail.NewEmail(
		"MergeSpace",
		s.from,
	)

	recipient := mail.NewEmail(
		"",
		to,
	)

	message := mail.NewSingleEmail(
		from,
		subject,
		recipient,
		"",
		body,
	)

	client := sendgrid.NewSendClient(
		s.apiKey,
	)

	response, err := client.Send(message)

	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {

		log.Printf(
			"sendgrid failed status=%d body=%s",
			response.StatusCode,
			response.Body,
		)

		return fmt.Errorf(
			"sendgrid error: status=%d",
			response.StatusCode,
		)
	}

	log.Printf(
		"email sent status=%d to=%s",
		response.StatusCode,
		to,
	)

	return nil
}

func (s *SendGridSender) SendOTP(ctx context.Context, email string, otp string) error {

	body := fmt.Sprintf(`
		<h2>Your OTP</h2>

		<p>Your verification code is:</p>

		<h1>%s</h1>

		<p>Expires in 10 minutes.</p>
	`, otp)

	return s.send(
		email,
		"Your OTP Code",
		body,
	)
}

func (s *SendGridSender) SendWelcome(ctx context.Context, email string, firstName string) error {

	body := fmt.Sprintf(`
		<h1>Welcome to MergeSpace 🚀</h1>

		<p>Hello %s,</p>

		<p>Your account was created successfully.</p>

		<p>We're excited to have you onboard.</p>
	`, firstName)

	return s.send(
		email,
		"Welcome to MergeSpace",
		body,
	)
}

func (s *SendGridSender) SendForgotPasswordOTP(ctx context.Context, to string, otp string) error {

	body := fmt.Sprintf(`
	<h1>Reset Your Password</h1>

	<p>Your password reset OTP is:</p>

	<h2>%s</h2>

	<p>This OTP expires in 10 minutes.</p>
	`, otp)

	return s.send(
		to,
		"Reset Password OTP",
		body,
	)
}
