package valueobject

import (
	"strings"
	"unicode"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

type Username struct {
	value string
}

const (
	minLength = 3
	maxLength = 30
)

func NewUsername(username string) (Username, error) {
	username = normalize(username)

	if err := validateNotEmpty(username); err != nil {
		return Username{}, err
	}

	if err := validateLength(username); err != nil {
		return Username{}, err
	}

	if err := validateEdges(username); err != nil {
		return Username{}, err
	}

	if err := validateCharacters(username); err != nil {
		return Username{}, err
	}

	if err := validateSequence(username); err != nil {
		return Username{}, err
	}

	if err := validateContainsAlphabet(username); err != nil {
		return Username{}, err
	}

	return Username{value: username}, nil
}

func (u Username) String() string {
	return u.value
}

func normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func validateNotEmpty(u string) error {
	if u == "" {
		return errs.ErrEmptyUsername
	}
	return nil
}

func validateLength(u string) error {
	if len(u) < minLength {
		return errs.ErrUsernameTooShort
	}
	if len(u) > maxLength {
		return errs.ErrUsernameTooLong
	}
	return nil
}

func validateEdges(u string) error {
	if isSpecialChar(rune(u[0])) {
		return errs.ErrStartsWithSpecialChar
	}
	if isSpecialChar(rune(u[len(u)-1])) {
		return errs.ErrEndsWithSpecialChar
	}
	return nil
}

func validateCharacters(u string) error {
	for _, r := range u {
		if !isAllowedChar(r) {
			return errs.ErrInvalidCharacter
		}
	}
	return nil
}

func validateSequence(u string) error {
	for i := 1; i < len(u); i++ {
		if isSpecialChar(rune(u[i])) && isSpecialChar(rune(u[i-1])) {
			return errs.ErrConsecutiveSpecialChars
		}
	}
	return nil
}

func validateContainsAlphabet(u string) error {
	for _, r := range u {
		if unicode.IsLetter(r) {
			return nil
		}
	}
	return errs.ErrUsernameNoAlphabet
}

func isSpecialChar(r rune) bool {
	return r == '.' || r == '_' || r == '-'
}

func isAllowedChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || isSpecialChar(r)
}