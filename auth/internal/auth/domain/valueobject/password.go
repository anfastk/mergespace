package valueobject

import (
	"strings"
	"unicode"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
)

type Password struct {
	value string
}

func NewPassword(v string, forbidden ...string) (Password, error) {
	v = strings.TrimSpace(v)

	if v == "" {
		return Password{}, errs.ErrInvalidPassword
	}

	if len(v) < minPasswordLength || len(v) > maxPasswordLength {
		return Password{}, errs.ErrInvalidPassword
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasDigit  bool
		hasSymbol bool

		repeatCount int
		lastRune    rune
	)

	for i, r := range v {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}

		if i > 0 {
			if r == lastRune {
				repeatCount++
				if repeatCount >= 2 {
					return Password{}, errs.ErrInvalidPassword
				}
			} else {
				repeatCount = 0
			}
		}
		lastRune = r
	}

	if !(hasUpper && hasLower && hasDigit && hasSymbol) {
		return Password{}, errs.ErrInvalidPassword
	}

	if containsSequentialPatterns(v) {
		return Password{}, errs.ErrInvalidPassword
	}

	if hasWeakSuffix(v) {
		return Password{}, errs.ErrInvalidPassword
	}

	if containsForbiddenSubstring(v, forbidden) {
		return Password{}, errs.ErrInvalidPassword
	}

	if isCommonPassword(v) {
		return Password{}, errs.ErrInvalidPassword
	}

	return Password{value: v}, nil
}

func (p Password) String() string {
	return p.value
}
