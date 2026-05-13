package dto

type VerifyForgotPasswordOTPRequest struct {
	ResetID string
	OTP     string
}