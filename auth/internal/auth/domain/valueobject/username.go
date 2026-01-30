package valueobject

import (
	"strings"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/errs"
)

type Username struct {
	value string
}

func NewUsername(v string) (Username, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return Username{}, errs.ErrEmptyUsername
	}
	if len(v) < 3 {
		return Username{}, errs.ErrUsernameTooShort
	}
	return Username{value: v}, nil
}

func (u Username) String() string {
	return u.value
}
