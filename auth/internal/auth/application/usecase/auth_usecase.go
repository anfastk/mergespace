package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
	appErr "github.com/anfastk/mergespace/auth/internal/auth/application/errors"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/inbound"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	domainErr "github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ inbound.AuthUseCase = (*AuthService)(nil)

type AuthService struct {
	db                 *pgxpool.Pool
	userRepo           outbound.UserRepository
	otpGen             outbound.OTPGenerator
	idGen              outbound.IDGenerator
	signupCtxStore     outbound.SignupContextStore
	passwordHasher     outbound.PasswordHasher
	eventProducer      outbound.EventProducer
	outboxRepo         outbound.OutboxRepository
	tokenGenerator     outbound.TokenGenerator
	passwordResetStore outbound.PasswordResetStore
	oauthProvider      outbound.OAuthProvider
	authIdentityRepo   outbound.AuthIdentityRepository
}

func NewAuthService(db *pgxpool.Pool, user outbound.UserRepository, otpGen outbound.OTPGenerator, idGen outbound.IDGenerator, signupCtxStore outbound.SignupContextStore, passwordHasher outbound.PasswordHasher, producer outbound.EventProducer, outboxRepo outbound.OutboxRepository, tokenGen outbound.TokenGenerator, passwordResetStore outbound.PasswordResetStore, oauthProvider outbound.OAuthProvider, authIdentityRepo outbound.AuthIdentityRepository) *AuthService {
	return &AuthService{
		db:                 db,
		userRepo:           user,
		otpGen:             otpGen,
		idGen:              idGen,
		signupCtxStore:     signupCtxStore,
		passwordHasher:     passwordHasher,
		eventProducer:      producer,
		outboxRepo:         outboxRepo,
		tokenGenerator:     tokenGen,
		passwordResetStore: passwordResetStore,
		oauthProvider:      oauthProvider,
		authIdentityRepo:   authIdentityRepo,
	}
}

func (s *AuthService) InitiateSignup(ctx context.Context, req *dto.InitiateSignUpRequest) (*dto.InitiateSignUpResponce, error) {

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "email",
			Err:   err,
		}
	}

	firstname, err := valueobject.NewName(req.FirstName)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "first_name",
			Err:   err,
		}
	}

	lastname, err := valueobject.NewName(req.LastName)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "last_name",
			Err:   err,
		}
	}

	userName, err := valueobject.NewUsername(req.UserName)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "username",
			Err:   err,
		}
	}

	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "password",
			Err:   err,
		}
	}

	existing, err := s.signupCtxStore.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing != nil {

		lastSentAt, err := s.signupCtxStore.GetLastOTPSentAt(ctx, existing.ID)
		if err == nil {

			remaining := 60*time.Second - time.Since(lastSentAt)

			if remaining > 0 {
				return &dto.InitiateSignUpResponce{
					TempID: string(existing.ID),
					Message: fmt.Sprintf(
						"OTP already sent. Please wait %d seconds before requesting again.",
						int(remaining.Seconds()),
					),
				}, nil
			}
		}

		return &dto.InitiateSignUpResponce{
			TempID:  string(existing.ID),
			Message: "OTP already sent. You can use resend OTP.",
		}, nil
	}

	acquired, err := s.signupCtxStore.AcquireSignupSlot(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire signup slot: %w", err)
	}
	if !acquired {
		return nil, errs.ErrOTPAlreadySent
	}

	release := func() {
		if err := s.signupCtxStore.ReleaseSignupSlot(ctx, email); err != nil {
			log.Printf("failed to release signup slot: %v", err)
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		release()
		return nil, err
	}
	defer tx.Rollback(ctx)

	if exists, err := s.userRepo.ExistsByEmail(ctx, email.String()); err != nil {
		release()
		return nil, err
	} else if exists {
		release()
		return nil, errs.ErrEmailAlreadyExists
	}

	if exists, err := s.userRepo.ExistsByUsername(ctx, userName.String()); err != nil {
		release()
		return nil, err
	} else if exists {
		release()
		return nil, domainErr.ErrUsernameExists
	}

	passwordHash, err := s.passwordHasher.Hash(password.String())
	if err != nil {
		release()
		return nil, err
	}

	tempID := entity.SignupContextID(s.idGen.NewID(ctx))

	otp, err := s.otpGen.Generate(ctx)
	if err != nil {
		release()
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

	if err := s.signupCtxStore.Save(ctx, pending); err != nil {
		release()
		return nil, err
	}

	if err := s.signupCtxStore.SetLastOTPSentAt(ctx, tempID, time.Now()); err != nil {
		log.Printf("failed to save otp sent time: %v", err)
	}

	payload, _ := json.Marshal(map[string]string{
		"email": email.String(),
		"otp":   otp,
	})

	outbox := &entity.OutboxEvent{
		ID:        s.idGen.NewID(ctx),
		EventType: "SendOTP",
		Payload:   payload,
		Status:    "pending",
	}

	if err := s.outboxRepo.Save(ctx, tx, outbox); err != nil {
		release()
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		release()
		return nil, err
	}

	return &dto.InitiateSignUpResponce{
		TempID:  string(tempID),
		Message: "Otp Send Success",
	}, nil
}

func (s *AuthService) VerifySignup(ctx context.Context, req *dto.VerifySignupRequest) (*dto.AuthResponse, error) {

	signupCtx, err := s.signupCtxStore.FindByID(ctx, entity.SignupContextID(req.TempID))
	if err != nil {
		return nil, err
	}

	if signupCtx == nil {
		return nil, errs.ErrSignupContextNotFound
	}

	attempts, err := s.signupCtxStore.GetAttempts(ctx, entity.SignupContextID(req.TempID))
	if err != nil {
		return nil, err
	}

	if attempts >= 5 {
		return nil, errs.ErrTooManyAttempts
	}

	if signupCtx.OTP != req.OTP {

		if err := s.signupCtxStore.IncrementAttempts(ctx, entity.SignupContextID(req.TempID), 10*time.Minute); err != nil {
			log.Println("failed to increment attempts:", err)
		}

		return nil, errs.ErrOTPInvalid
	}

	if err = s.signupCtxStore.DeleteAttempts(ctx, entity.SignupContextID(req.TempID)); err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(s.idGen.NewID(ctx))
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(signupCtx.Email)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "email",
			Err:   err,
		}
	}

	username, err := valueobject.NewUsername(signupCtx.Username)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "username",
			Err:   err,
		}
	}

	firstname, err := valueobject.NewName(signupCtx.FirstName)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "first_name",
			Err:   err,
		}
	}

	lastname, err := valueobject.NewName(signupCtx.LastName)
	if err != nil {
		return nil, appErr.FieldError{
			Field: "last_name",
			Err:   err,
		}
	}

	user := entity.NewLocalUser(
		userID,
		email,
		username,
		&signupCtx.PasswordHash,
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(map[string]string{
		"user_id":   userID.String(),
		"firstName": firstname.String(),
		"lastName":  lastname.String(),
		"email":     email.String(),
		"username":  username.String(),
	})
	if err != nil {
		return nil, err
	}

	outbox := &entity.OutboxEvent{
		ID:        s.idGen.NewID(ctx),
		EventType: "UserCreated",
		Payload:   payload,
		Status:    "pending",
	}

	if err := s.outboxRepo.Save(ctx, tx, outbox); err != nil {
		return nil, err
	}

	if err = s.signupCtxStore.Delete(ctx, entity.SignupContextID(req.TempID)); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	accessToken, accessExpiry, err := s.tokenGenerator.GenerateAccessToken(userID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(userID.String())
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserRes{
			ID:       user.UserID.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Status:   string(user.Status),
		},
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiry,
	}, nil
}

func (s *AuthService) CheckUsernameAvailability(ctx context.Context, req *dto.CheckUsernameReq) (*dto.CheckUsernameRes, error) {

	username, err := valueobject.NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	exists, err := s.userRepo.ExistsByUsername(ctx, username.String())
	if err != nil {
		return nil, err
	}

	if !exists {
		return &dto.CheckUsernameRes{
			Available:   true,
			Message:     "Username is available",
			Suggestions: []string{},
		}, nil
	}

	suggestions := s.generateSuggestions(username.String())

	return &dto.CheckUsernameRes{
		Available:   false,
		Message:     "Username is already taken",
		Suggestions: suggestions,
	}, nil
}

func (s *AuthService) generateSuggestions(base string) []string {
	var suggestions []string
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	s1 := fmt.Sprintf("%s%d", base, randSource.Intn(900)+100)
	suggestions = append(suggestions, s1)

	s2 := fmt.Sprintf("%s%d", base, time.Now().Year())
	suggestions = append(suggestions, s2)

	s3 := fmt.Sprintf("official.%s", base)
	suggestions = append(suggestions, s3)

	suffixes := []string{"_dev", "_pro", "_code", "x"}
	s4 := fmt.Sprintf("%s%s", base, suffixes[randSource.Intn(len(suffixes))])
	suggestions = append(suggestions, s4)

	return suggestions
}

func (s *AuthService) ResendOTP(ctx context.Context, req *dto.ResendOTPRequest) (*dto.InitiateSignUpResponce, error) {

	signupCtx, err := s.signupCtxStore.FindByID(ctx, entity.SignupContextID(req.TempID))
	if err != nil {
		return nil, err
	}

	if signupCtx == nil {
		return nil, errs.ErrSignupContextNotFound
	}

	lastSentAt, err := s.signupCtxStore.GetLastOTPSentAt(ctx, signupCtx.ID)
	if err == nil {
		if time.Since(lastSentAt) < 60*time.Second {
			return nil, errs.ErrTooManyRequests
		}
	}

	otp, err := s.otpGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	signupCtx.OTP = otp

	if err := s.signupCtxStore.Update(ctx, signupCtx); err != nil {
		return nil, err
	}

	if err := s.signupCtxStore.SetLastOTPSentAt(ctx, signupCtx.ID, time.Now()); err != nil {
		log.Println("failed to set OTP sent time:", err)
	}

	payload, _ := json.Marshal(map[string]string{
		"email": signupCtx.Email,
		"otp":   otp,
	})

	outbox := &entity.OutboxEvent{
		ID:        s.idGen.NewID(ctx),
		EventType: "SendOTP",
		Payload:   payload,
		Status:    "pending",
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.outboxRepo.Save(ctx, tx, outbox); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.InitiateSignUpResponce{
		TempID:  string(signupCtx.ID),
		Message: "OTP resent successfully",
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {

	user, err := s.userRepo.FindByEmailOrUsername(
		ctx,
		req.EmailOrUsername,
	)
	if err != nil {
		return nil, err
	}

	if user.Password == nil {
		return nil, errs.ErrInvalidCredentials
	}

	if err := s.passwordHasher.Compare(
		*user.Password,
		req.Password,
	); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	if user.Status != entity.UserStatusActive {
		return nil, errs.ErrUserInactive
	}

	accessToken, accessExpiry, err := s.tokenGenerator.GenerateAccessToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(
		user.UserID.String(),
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserRes{
			ID:       user.UserID.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Status:   string(user.Status),
		},
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessExpiresAt: accessExpiry,
	}, nil
}

func (s *AuthService) ForgotPasswordInitiate(ctx context.Context, req *dto.ForgotPasswordRequest) (*dto.InitiateSignUpResponce, error) {

	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.EmailOrUsername)

	// silent success
	if err != nil {
		return &dto.InitiateSignUpResponce{
			Message: "If account exists, OTP has been sent.",
		}, nil
	}

	otp, err := s.otpGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	resetID := entity.PasswordResetContextID(s.idGen.NewID(ctx))

	reset := &entity.PasswordResetContext{
		ID:          resetID,
		UserID:      user.UserID.String(),
		Email:       user.Email.String(),
		OTP:         otp,
		OTPVerified: false,
	}

	if err := s.passwordResetStore.Save(ctx, reset); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(map[string]string{
		"email": user.Email.String(),
		"otp":   otp,
	})

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	outbox := &entity.OutboxEvent{
		ID:        s.idGen.NewID(ctx),
		EventType: "SendOTP",
		Payload:   payload,
		Status:    "pending",
	}

	if err := s.outboxRepo.Save(ctx, tx, outbox); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.InitiateSignUpResponce{
		TempID:  string(resetID),
		Message: "If account exists, OTP has been sent.",
	}, nil
}

func (s *AuthService) VerifyForgotPasswordOTP(ctx context.Context, req *dto.VerifyForgotPasswordOTPRequest) error {

	reset, err := s.passwordResetStore.FindByID(
		ctx,
		entity.PasswordResetContextID(req.ResetID),
	)
	if err != nil {
		return err
	}

	if reset == nil {
		return errs.ErrSignupContextNotFound
	}

	attempts, err := s.passwordResetStore.GetAttempts(ctx, reset.ID)
	if err != nil {
		return err
	}

	if attempts >= 5 {
		return errs.ErrTooManyAttempts
	}

	if reset.OTP != req.OTP {

		if err := s.passwordResetStore.IncrementAttempts(
			ctx,
			reset.ID,
			10*time.Minute,
		); err != nil {
			log.Println(err)
		}

		return errs.ErrOTPInvalid
	}

	reset.OTPVerified = true

	if err := s.passwordResetStore.Update(ctx, reset); err != nil {
		return err
	}

	if err := s.passwordResetStore.DeleteAttempts(ctx, reset.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {

	reset, err := s.passwordResetStore.FindByID(
		ctx,
		entity.PasswordResetContextID(req.ResetID),
	)
	if err != nil {
		return err
	}

	if reset == nil {
		return errs.ErrSignupContextNotFound
	}

	if !reset.OTPVerified {
		return errs.ErrOTPInvalid
	}

	password, err := valueobject.NewPassword(req.NewPassword)
	if err != nil {
		return errs.ErrInvalidPassword
	}

	hash, err := s.passwordHasher.Hash(password.String())
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(
		ctx,
		reset.UserID,
		string(hash),
	); err != nil {
		return err
	}

	if err := s.passwordResetStore.Delete(ctx, reset.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendForgotPasswordOTP(ctx context.Context, req *dto.ResendForgotPasswordOTPRequest) (*dto.InitiateSignUpResponce, error) {

	reset, err := s.passwordResetStore.FindByID(
		ctx,
		entity.PasswordResetContextID(req.ResetID),
	)
	if err != nil {
		return nil, err
	}

	if reset == nil {
		return nil, errs.ErrSignupContextNotFound
	}

	// cooldown check
	lastSentAt, err := s.passwordResetStore.GetLastOTPSentAt(
		ctx,
		reset.ID,
	)

	if err == nil {
		if time.Since(lastSentAt) < 60*time.Second {
			return nil, errs.ErrTooManyRequests
		}
	}

	otp, err := s.otpGen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	reset.OTP = otp

	if err := s.passwordResetStore.Update(ctx, reset); err != nil {
		return nil, err
	}

	if err := s.passwordResetStore.SetLastOTPSentAt(
		ctx,
		reset.ID,
		time.Now(),
	); err != nil {
		log.Println("failed to update resend timestamp:", err)
	}

	payload, err := json.Marshal(map[string]string{
		"email": reset.Email,
		"otp":   otp,
	})
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	outbox := &entity.OutboxEvent{
		ID:        s.idGen.NewID(ctx),
		EventType: "SendOTP",
		Payload:   payload,
		Status:    "pending",
	}

	if err := s.outboxRepo.Save(ctx, tx, outbox); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.InitiateSignUpResponce{
		TempID:  string(reset.ID),
		Message: "OTP resent successfully",
	}, nil
}
