package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/redis/go-redis/v9"
)

const signupContextKeyPrefix = "signup_context:pending:"
const signupContextTTL = 10 * time.Minute

type SignupContextRedisStore struct {
	client *redis.Client
}

var _ outbound.SignupContextStore = (*SignupContextRedisStore)(nil)

func NewSignupContextRedisStore(client *redis.Client) outbound.SignupContextStore {
	return &SignupContextRedisStore{client: client}
}

func (r *SignupContextRedisStore) redisKey(id entity.SignupContextID) string {
	return fmt.Sprintf("%s%s", signupContextKeyPrefix, id)
}

func (r *SignupContextRedisStore) Set(ctx context.Context, signup *entity.SignupContext) error {
	model := SignupContexModel{
		ID:           string(signup.ID),
		FirstName:    signup.FirstName,
		LastName:     signup.LastName,
		Email:        signup.Email,
		Username:     signup.Username,
		PasswordHash: signup.PasswordHash,
		OTP:          signup.OTP,
	}

	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.redisKey(signup.ID), data, 10*time.Minute).Err()
}
