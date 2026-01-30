package otp

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
)

type CryptoOTPGenerator struct{}

var _ outbound.OTPGenerator = (*CryptoOTPGenerator)(nil)

const length = 6
const digits = "0123456789"

var ErrOTPRandFailed = errors.New("failed to generate secure OTP")

func NewCryptoOTPGenerator() outbound.OTPGenerator {
	return &CryptoOTPGenerator{}
}

func (CryptoOTPGenerator) Generate(ctx context.Context) (string, error) {
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", ErrOTPRandFailed
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}
