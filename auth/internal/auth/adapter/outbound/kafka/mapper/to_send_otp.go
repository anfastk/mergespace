package mapper

import (
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/kafka/avro"
	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
)

func ToSendOTPAvro(e event.SendOTP) avro.SendOTP {
	return avro.SendOTP{
		Email: e.Email,
		OTP:   e.OTP,
	}
}
