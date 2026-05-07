package mapper

import (
	"errors"
	"strings"

	"connectrpc.com/connect"
	appErr "github.com/anfastk/mergespace/auth/internal/auth/application/errors"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type errorMeta struct {
	code    connect.Code
	message string
}

var errorMap = map[error]errorMeta{

	// ---------------- USER ----------------
	errs.ErrUserNotFound: {
		code:    connect.CodeNotFound,
		message: "User not found.",
	},
	errs.ErrUserAlreadyExists: {
		code:    connect.CodeAlreadyExists,
		message: "User already exists.",
	},
	errs.ErrUserDeleted: {
		code:    connect.CodeFailedPrecondition,
		message: "User account is deleted.",
	},

	// ---------------- NAME ----------------
	errs.ErrEmptyName: {
		code:    connect.CodeInvalidArgument,
		message: "Name cannot be empty.",
	},
	errs.ErrInvalidName: {
		code:    connect.CodeInvalidArgument,
		message: "Name must contain letters only.",
	},
	errs.ErrNameTooShort: {
		code:    connect.CodeInvalidArgument,
		message: "Name is too short.",
	},
	errs.ErrNameTooLong: {
		code:    connect.CodeInvalidArgument,
		message: "Name is too long.",
	},

	// ---------------- USERNAME ----------------
	errs.ErrEmptyUsername: {
		code:    connect.CodeInvalidArgument,
		message: "Username cannot be empty.",
	},
	errs.ErrUsernameTooShort: {
		code:    connect.CodeInvalidArgument,
		message: "Username must be at least 3 characters long.",
	},
	errs.ErrUsernameTooLong: {
		code:    connect.CodeInvalidArgument,
		message: "Username must not exceed 30 characters.",
	},
	errs.ErrInvalidCharacter: {
		code:    connect.CodeInvalidArgument,
		message: "Username can only contain letters, numbers, '.', '_' and '-'.",
	},
	errs.ErrStartsWithSpecialChar: {
		code:    connect.CodeInvalidArgument,
		message: "Username cannot start with '.', '_' or '-'.",
	},
	errs.ErrEndsWithSpecialChar: {
		code:    connect.CodeInvalidArgument,
		message: "Username cannot end with '.', '_' or '-'.",
	},
	errs.ErrConsecutiveSpecialChars: {
		code:    connect.CodeInvalidArgument,
		message: "Username cannot contain consecutive special characters.",
	},
	errs.ErrUsernameNoAlphabet: {
		code:    connect.CodeInvalidArgument,
		message: "Username must contain at least one letter.",
	},
	errs.ErrUsernameExists: {
		code:    connect.CodeAlreadyExists,
		message: "Username is already taken.",
	},

	// ---------------- EMAIL ----------------
	errs.ErrEmptyEmail: {
		code:    connect.CodeInvalidArgument,
		message: "Email cannot be empty.",
	},
	errs.ErrInvalidEmail: {
		code:    connect.CodeInvalidArgument,
		message: "Invalid email address.",
	},
	errs.ErrEmailAlreadyExists: {
		code:    connect.CodeAlreadyExists,
		message: "Email address already exists.",
	},

	// ---------------- AUTH ----------------
	errs.ErrInvalidPassword: {
		code:    connect.CodeInvalidArgument,
		message: "Invalid password.",
	},
	errs.ErrInvalidCredentials: {
		code:    connect.CodeUnauthenticated,
		message: "Invalid credentials.",
	},
	errs.ErrAccountSuspended: {
		code:    connect.CodePermissionDenied,
		message: "Account is suspended.",
	},
	errs.ErrSignupContextNotFound: {
		code:    connect.CodeNotFound,
		message: "Signup context not found.",
	},
	errs.ErrUserInactive: {
		code:    connect.CodePermissionDenied,
		message: "Your account is inactive. Please contact support.",
	},

	// ---------------- OTP ----------------
	errs.ErrOTPInvalid: {
		code:    connect.CodeInvalidArgument,
		message: "Invalid OTP.",
	},
	errs.ErrOTPExpired: {
		code:    connect.CodeFailedPrecondition,
		message: "OTP has expired.",
	},
	errs.ErrOTPAlreadySent: {
		code:    connect.CodeAlreadyExists,
		message: "OTP already sent.",
	},

	// ---------------- SYSTEM ----------------
	errs.ErrTooManyRequests: {
		code:    connect.CodeResourceExhausted,
		message: "Too many requests. Please try again later.",
	},

	errs.ErrTooManyAttempts: {
		code:    connect.CodeResourceExhausted,
		message: "Too many OTP attempts. Please try again later.",
	},
}

func MapDomainError(err error) error {

	var fe appErr.FieldError
	if errors.As(err, &fe) {

		for domainErr, meta := range errorMap {
			if errors.Is(fe.Err, domainErr) {

				c := cases.Title(language.English)

				fieldName := strings.ReplaceAll(fe.Field, "_", " ")
				fieldName = c.String(fieldName)

				msg := strings.Replace(meta.message, "Name", fieldName, 1)

				return connect.NewError(meta.code, errors.New(msg))
			}
		}
	}

	if meta, ok := errorMap[err]; ok {
		return connect.NewError(meta.code, errors.New(meta.message))
	}

	for domainErr, meta := range errorMap {
		if errors.Is(err, domainErr) {
			return connect.NewError(meta.code, errors.New(meta.message))
		}
	}

	return connect.NewError(connect.CodeInternal, errors.New("Internal server error"))
}
