package valueobject

import (
	"regexp"
	"strings"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

var nameRegex = regexp.MustCompile(`^[\p{L}]+$`)

type Name struct {
	value string
}

func 	NewName(raw string) (Name, error) {
	raw = strings.TrimSpace(raw)

	if len(raw) < 1 {
		return Name{}, errs.ErrNameTooShort
	}
	if len(raw) > 50 {
		return Name{}, errs.ErrNameTooLong
	}

	if !nameRegex.MatchString(raw) {
		return Name{}, errs.ErrInvalidName
	}

	return Name{value: raw}, nil
}

func (n Name) String() string {
	return n.value
}
