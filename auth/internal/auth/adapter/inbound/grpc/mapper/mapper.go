package mapper

import (
	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
	authv1 "github.com/anfastk/mergespace/contracts/gen/go/proto/auth/v1"
)

func ToInitiateDTO(req *authv1.InitiateSignupRequest) *dto.InitiateSignUpRequest {
	return &dto.InitiateSignUpRequest{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Password: req.Password}
}
