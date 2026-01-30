package grpc

import (
	"context"

	"connectrpc.com/connect"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/inbound/grpc/mapper"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/inbound"
	authv1 "github.com/anfastk/mergespace/contracts/gen/go/proto/auth/v1"
)

type AuthHandler struct {
	usecase inbound.AuthUseCase
}

func NewAuthHandler(usecase inbound.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (h *AuthHandler) InitiateSignup(ctx context.Context, req *connect.Request[authv1.InitiateSignupRequest]) (*connect.Response[authv1.InitiateSignupResponse], error) {
	res, err := h.usecase.InitiateSignup(ctx, mapper.ToInitiateDTO(req.Msg))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&authv1.InitiateSignupResponse{
		TempId:  res.TempID,
		Message: res.Message,
	}), nil
}
