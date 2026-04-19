package grpc

import (
	"context"
	"fmt"
	"log"

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

func (h *AuthHandler) CheckUsernameAvailability(ctx context.Context, req *connect.Request[authv1.CheckUsernameRequest]) (*connect.Response[authv1.CheckUsernameResponse], error) {
	result, err := h.usecase.CheckUsernameAvailability(ctx, mapper.ToCheckUsernameDTO(req.Msg))
	if err != nil {
		fmt.Println(err)
		return nil, mapper.MapDomainError(err)
	}

	return connect.NewResponse(&authv1.CheckUsernameResponse{
		Available:   result.Available,
		Message:     result.Message,
		Suggestions: result.Suggestions,
	}), nil
}

func (h *AuthHandler) InitiateSignup(ctx context.Context, req *connect.Request[authv1.InitiateSignupRequest]) (*connect.Response[authv1.InitiateSignupResponse], error) {
	if h.usecase == nil {
		log.Fatal("usecase is nil")
	}

	res, err := h.usecase.InitiateSignup(ctx, mapper.ToInitiateDTO(req.Msg))

	if err != nil {
		return nil, mapper.MapDomainError(err)
	}

	return connect.NewResponse(&authv1.InitiateSignupResponse{
		TempId:  res.TempID,
		Message: res.Message,
	}), nil
}
