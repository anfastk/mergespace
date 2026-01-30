package event

type SendOTP struct {
	Email string
	OTP   string
}
