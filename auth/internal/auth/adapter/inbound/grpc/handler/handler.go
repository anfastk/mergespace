package grpc

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/inbound/grpc/mapper"
	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/inbound"
	authv1 "github.com/anfastk/mergespace/contracts/gen/go/proto/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (h *AuthHandler) VerifySignup(ctx context.Context, req *connect.Request[authv1.VerifySignupRequest]) (*connect.Response[authv1.AuthResponse], error) {
	res, err := h.usecase.VerifySignup(ctx, mapper.ToVerifySignupDTO(req.Msg))
	if err != nil {
		log.Println(err)
		return nil, mapper.MapDomainError(err)
	}

	response := connect.NewResponse(&authv1.AuthResponse{
		User: &authv1.UserRes{
			Id:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
			Avatar:   res.User.Avatar,
			Status:   res.User.Status,
		},
		AccessToken:     res.AccessToken,
		AccessExpiresAt: timestamppb.New(res.AccessExpiresAt),
	})

	response.Header().Add("Set-Cookie", buildRefreshCookie(res.RefreshToken))

	return response, nil
}

func buildRefreshCookie(token string) string {
	return fmt.Sprintf(
		"refresh_token=%s; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=%d",
		token,
		7*24*60*60,
	)
}

func (h *AuthHandler) ResendOTP(ctx context.Context, req *connect.Request[authv1.ResendOTPRequest]) (*connect.Response[authv1.InitiateSignupResponse], error) {

	res, err := h.usecase.ResendOTP(ctx, &dto.ResendOTPRequest{
		TempID: req.Msg.TempId,
	})
	if err != nil {
		return nil, mapper.MapDomainError(err)
	}

	return connect.NewResponse(&authv1.InitiateSignupResponse{
		TempId:  res.TempID,
		Message: res.Message,
	}), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.AuthResponse], error) {

	res, err := h.usecase.Login(ctx, &dto.LoginRequest{
		EmailOrUsername: req.Msg.EmailOrUsername,
		Password:        req.Msg.Password,
	})
	if err != nil {
		return nil, mapper.MapDomainError(err)
	}

	response := connect.NewResponse(&authv1.AuthResponse{
		User: &authv1.UserRes{
			Id:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
			Status:   res.User.Status,
		},
		AccessToken:     res.AccessToken,
		AccessExpiresAt: timestamppb.New(res.AccessExpiresAt),
	})

	response.Header().Add(
		"Set-Cookie",
		buildRefreshCookie(res.RefreshToken),
	)

	return response, nil
}
