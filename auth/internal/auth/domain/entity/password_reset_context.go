package entity

type PasswordResetContextID string

type PasswordResetContext struct {
	ID          PasswordResetContextID
	UserID      string
	Email       string
	OTP         string
	OTPVerified bool
}
