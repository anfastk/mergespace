package event

type ForgotPasswordOTP struct {
	Email string `json:"email" avro:"email"`
	OTP   string `json:"otp" avro:"otp"`
}
