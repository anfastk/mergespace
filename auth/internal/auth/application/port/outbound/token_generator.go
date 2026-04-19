package outbound

import "github.com/anfastk/mergespace/auth/internal/auth/domain/entity"

type TokenGenerator interface {
	Generate(user *entity.User) (string, string, int64, error)
}
