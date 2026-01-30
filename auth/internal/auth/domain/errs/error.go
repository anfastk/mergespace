package errs

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserDeleted              = errors.New("user is deleted")
	ErrEmptyName                = errors.New("name is empty")
	ErrInvalidName              = errors.New("name must contain letters only")
	ErrNameTooShort             = errors.New("name is too short")
	ErrNameTooLong              = errors.New("name is too long")
	ErrEmptyUsername            = errors.New("username must not be empty")
	ErrUsernameTooShort         = errors.New("username too short")
	ErrUsernameExists           = errors.New("username already taken")
	ErrEmptyEmail               = errors.New("email must not be empty")
	ErrInvalidEmail             = errors.New("invalid email")
	ErrEmailAlreadyExists       = errors.New("email address already exists")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrAccountSuspended         = errors.New("account suspended")
	ErrOTPInvalid               = errors.New("invalid otp")
	ErrOTPExpired               = errors.New("otp expired")
	ErrUsernameGenerationFailed = errors.New("unable to generate username")
	ErrInternalServer           = errors.New("internal server error")
	ErrTooManyRequests          = errors.New("too many requests, please try again later")
)
