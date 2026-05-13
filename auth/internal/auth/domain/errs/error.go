package errs

import "errors"

var (
	// User
	ErrUserNotFound      = errors.New("USER_NOT_FOUND")
	ErrUserAlreadyExists = errors.New("USER_ALREADY_EXISTS")
	ErrUserDeleted       = errors.New("USER_DELETED")

	// Name
	ErrEmptyName    = errors.New("NAME_EMPTY")
	ErrInvalidName  = errors.New("NAME_INVALID")
	ErrNameTooShort = errors.New("NAME_TOO_SHORT")
	ErrNameTooLong  = errors.New("NAME_TOO_LONG")

	// Username
	ErrEmptyUsername           = errors.New("USERNAME_EMPTY")
	ErrUsernameTooShort        = errors.New("USERNAME_TOO_SHORT")
	ErrUsernameTooLong         = errors.New("USERNAME_TOO_LONG")
	ErrUsernameExists          = errors.New("USERNAME_ALREADY_EXISTS")
	ErrInvalidUsername         = errors.New("USERNAME_INVALID_FORMAT")
	ErrInvalidCharacter        = errors.New("USERNAME_INVALID_CHARACTER")
	ErrStartsWithSpecialChar   = errors.New("USERNAME_STARTS_WITH_SPECIAL_CHAR")
	ErrEndsWithSpecialChar     = errors.New("USERNAME_ENDS_WITH_SPECIAL_CHAR")
	ErrConsecutiveSpecialChars = errors.New("USERNAME_CONSECUTIVE_SPECIAL_CHARS")
	ErrUsernameNoAlphabet      = errors.New("USERNAME_NO_ALPHABET")

	// Email
	ErrEmptyEmail         = errors.New("EMAIL_EMPTY")
	ErrInvalidEmail       = errors.New("EMAIL_INVALID")
	ErrEmailAlreadyExists = errors.New("EMAIL_ALREADY_EXISTS")

	// Auth
	ErrInvalidPassword       = errors.New("PASSWORD_INVALID")
	ErrInvalidCredentials    = errors.New("INVALID_CREDENTIALS")
	ErrAccountSuspended      = errors.New("ACCOUNT_SUSPENDED")
	ErrSignupContextNotFound = errors.New("SIGNUP_CONTEXT_NOT_FOUND")
	ErrUserInactive          = errors.New("USER_INACTIVE")

	// OTP
	ErrOTPInvalid     = errors.New("OTP_INVALID")
	ErrOTPExpired     = errors.New("OTP_EXPIRED")
	ErrOTPAlreadySent = errors.New("OTP_ALREADY_SENT")
	ErrOTPNotVerified = errors.New("OTP_NOT_VERIFIED")

	// System
	ErrInternalServer  = errors.New("INTERNAL_SERVER_ERROR")
	ErrTooManyAttempts = errors.New("TOO_MANY_ATTEMPTS")
	ErrTooManyRequests = errors.New("TOO_MANY_REQUESTS")
)
