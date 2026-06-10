package dto

type ForgotPasswordOTPEvent struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
