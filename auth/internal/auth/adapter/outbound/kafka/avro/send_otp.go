package avro

type SendOTP struct {
	Email string `avro:"email"`
	OTP   string `avro:"otp"`
}
