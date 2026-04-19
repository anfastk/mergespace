package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
	"github.com/redis/go-redis/v9"
)

const signupContextKeyPrefix = "signup_context:pending:id:"
const signupEmailIndexPrefix = "signup_context:pending:email:"
const signupContextTTL = 10 * time.Minute

type SignupContextModel struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	OTP          string `json:"otp"`
}

type SignupContextRedisStore struct {
	client *redis.Client
}

var _ outbound.SignupContextStore = (*SignupContextRedisStore)(nil)

func NewSignupContextRedisStore(client *redis.Client) outbound.SignupContextStore {
	return &SignupContextRedisStore{client: client}
}

func (r *SignupContextRedisStore) redisIDKey(id entity.SignupContextID) string {
	return fmt.Sprintf("%s%s", signupContextKeyPrefix, id)
}

func (r *SignupContextRedisStore) redisEmailKey(email string) string {
	return fmt.Sprintf("%s%s", signupEmailIndexPrefix, email)
}

func (r *SignupContextRedisStore) Save(ctx context.Context, signup *entity.SignupContext) error {
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
		return fmt.Errorf("redis save - failed to marshal model: %w", err) // FIXED: Wrapped error
	}

	pipe := r.client.TxPipeline()

	pipe.Set(ctx, r.redisEmailKey(signup.Email), string(signup.ID), signupContextTTL)
	pipe.Set(ctx, r.redisIDKey(signup.ID), data, signupContextTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis save - pipeline execution failed: %w", err)
	}

	return nil
}

func (r *SignupContextRedisStore) FindByID(ctx context.Context, id entity.SignupContextID) (*entity.SignupContext, error) {
	data, err := r.client.Get(ctx, r.redisIDKey(id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("redis find - failed to get key: %w", err)
	}

	var model SignupContextModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &entity.SignupContext{
		ID:           entity.SignupContextID(model.ID),
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Email:        model.Email,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		OTP:          model.OTP,
	}, nil
}

func (r *SignupContextRedisStore) Delete(ctx context.Context, id entity.SignupContextID) error {
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("redis delete - failed to fetch existing context: %w", err)
	}
	if existing == nil {
		return nil
	}

	pipe := r.client.TxPipeline()

	pipe.Del(ctx, r.redisIDKey(id))
	pipe.Del(ctx, r.redisEmailKey(fmt.Sprint(existing.Email)))

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis delete - pipeline execution failed: %w", err)
	}
	return nil
}

func (r *SignupContextRedisStore) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.SignupContext, error) {
	emailStr := email.String()
	idStr, err := r.client.Get(ctx, r.redisEmailKey(emailStr)).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("redis find by email - failed to fetch index: %w", err)
	}

	return r.FindByID(ctx, entity.SignupContextID(idStr))
}

func (r *SignupContextRedisStore) AcquireSignupSlot(ctx context.Context, email valueobject.Email) (bool, error) {

	key := "signup:lock:" + email.String()

	ok, err := r.client.SetNX(ctx, key, "1", 10*time.Minute).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *SignupContextRedisStore) ReleaseSignupSlot(ctx context.Context, email valueobject.Email) error {

	key := "signup:lock:" + email.String()
	return r.client.Del(ctx, key).Err()
}
