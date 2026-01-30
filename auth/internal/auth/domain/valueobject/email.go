package valueobject

import (
	"regexp"
	"strings"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return Email{}, errs.ErrEmptyEmail
	}
	if !emailRegex.MatchString(raw) {
		return Email{}, errs.ErrInvalidEmail
	}
	return Email{value: raw}, nil
}

func (e Email) String() string {
	return e.value
}
