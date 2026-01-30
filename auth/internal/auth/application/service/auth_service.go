package service

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	domainErr "github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
)

type AuthService struct {
	user              outbound.UserRepository
	usernameAllocator outbound.UsernameAllocator
	otpGen            outbound.OTPGenerator
	idGen             outbound.IDGenerator
	signupCtxStore    outbound.SignupContextStore
	passwordHasher    outbound.PasswordHasher
	producer          outbound.EventProducer
}

func NewAuthService(user outbound.UserRepository, usernameAllocator outbound.UsernameAllocator, otpGen outbound.OTPGenerator, idGen outbound.IDGenerator, signupCtxStore outbound.SignupContextStore, passwordHasher outbound.PasswordHasher, producer outbound.EventProducer) *AuthService {
	return &AuthService{
		user:              user,
		usernameAllocator: usernameAllocator,
		otpGen:            otpGen,
		idGen:             idGen,
		signupCtxStore:    signupCtxStore,
		passwordHasher:    passwordHasher,
		producer:          producer,
	}
}

func (s *AuthService) InitiateSignup(ctx context.Context, req *dto.InitiateSignUpRequest) (*dto.InitiateSignUpResponce, error) {
	if u, err := s.user.ExistsByEmail(ctx, req.Email); err != nil {
		return nil, err
	} else if u {
		return nil, errs.ErrEmailAlreadyExists
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	firstname, err := valueobject.NewName(req.FirstName)
	if err != nil {
		return nil, err
	}

	lastname, err := valueobject.NewName(req.LastName)
	if err != nil {
		return nil, err
	}

	username, err := s.usernameAllocator.Allocate(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	userName, err := valueobject.NewUsername(username)
	if err != nil {
		return nil, err
	}

	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if u, err := s.user.ExistsByUsername(ctx, userName.String()); err != nil {
		return nil, err
	} else if u {
		return nil, domainErr.ErrUsernameExists
	}

	passwordHash, err := s.passwordHasher.Hash(password.String())
	if err != nil {
		return nil, err
	}

	tempID := entity.SignupContextID(s.idGen.NewID(ctx))
	otp, err := s.otpGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	pending := &entity.SignupContext{
		ID:           tempID,
		Email:        email.String(),
		FirstName:    firstname.String(),
		LastName:     lastname.String(),
		Username:     userName.String(),
		PasswordHash: string(passwordHash),
		OTP:          otp,
	}

	if err := s.signupCtxStore.Set(ctx, pending); err != nil {
		return nil, err
	}

	event := event.SendOTP{
		Email: email.String(),
		OTP:   otp,
	}

	if err = s.producer.PublishSendOTP(ctx, event); err != nil {
		return nil, err
	}

	return &dto.InitiateSignUpResponce{
		TempID:  string(tempID),
		Message: "Otp Send Success",
	}, nil
}
