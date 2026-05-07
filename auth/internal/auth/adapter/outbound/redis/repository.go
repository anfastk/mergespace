package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"
	"github.com/redis/go-redis/v9"
)

const signupContextKeyPrefix = "signup_context:pending:id:"
const signupEmailIndexPrefix = "signup_context:pending:email:"
const signupContextOTPAttemptsKeyPrefix = "otp:attempts:id:"
const signupOTPSentKeyPrefix = "signup:otp_sent:id:"
const signupContextTTL = 10 * time.Minute

type SignupContextRedisStore struct {
	client *redis.Client
}

var _ outbound.SignupContextStore = (*SignupContextRedisStore)(nil)

func NewSignupContextRedisStore(client *redis.Client) outbound.SignupContextStore {
	return &SignupContextRedisStore{client: client}
}

func (r *SignupContextRedisStore) redisSignupContextKey(id entity.SignupContextID) string {
	return fmt.Sprintf("%s%s", signupContextKeyPrefix, id)
}

func (r *SignupContextRedisStore) redisEmailKey(email string) string {
	return fmt.Sprintf("%s%s", signupEmailIndexPrefix, email)
}

func (r *SignupContextRedisStore) redisOTPAttemptsKey(id entity.SignupContextID) string {
	return fmt.Sprintf("%s%s", signupContextOTPAttemptsKeyPrefix, id)
}

func (r *SignupContextRedisStore) redisOTPSentKey(id entity.SignupContextID) string {
	return fmt.Sprintf("%s%s", signupOTPSentKeyPrefix, id)
}

func (r *SignupContextRedisStore) Save(ctx context.Context, signup *entity.SignupContext) error {
	model := SignupContextModel{
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
		return fmt.Errorf("redis save - failed to marshal model: %w", err)
	}

	pipe := r.client.TxPipeline()

	pipe.Set(ctx, r.redisEmailKey(signup.Email), string(signup.ID), signupContextTTL)
	pipe.Set(ctx, r.redisSignupContextKey(signup.ID), data, signupContextTTL)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis save - pipeline execution failed: %w", err)
	}

	return nil
}

func (r *SignupContextRedisStore) FindByID(ctx context.Context, id entity.SignupContextID) (*entity.SignupContext, error) {
	data, err := r.client.Get(ctx, r.redisSignupContextKey(id)).Bytes()
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

	pipe.Del(ctx, r.redisSignupContextKey(id))
	pipe.Del(ctx, r.redisEmailKey(fmt.Sprint(existing.Email)))
	pipe.Del(ctx, r.redisOTPAttemptsKey(id))
	pipe.Del(ctx, r.redisOTPSentKey(id))
	pipe.Del(ctx, "signup:lock:"+existing.Email)

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

func (r *SignupContextRedisStore) GetAttempts(ctx context.Context, id entity.SignupContextID) (int, error) {
	val, err := r.client.Get(ctx, r.redisOTPAttemptsKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	attempts, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return attempts, nil
}

func (r *SignupContextRedisStore) IncrementAttempts(ctx context.Context, id entity.SignupContextID, ttl time.Duration) error {
	key := r.redisOTPAttemptsKey(id)

	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	log.Printf("OTP attempts incremented: %d", incr.Val())

	return nil
}

func (r *SignupContextRedisStore) DeleteAttempts(ctx context.Context, id entity.SignupContextID) error {
	return r.client.Del(ctx, r.redisOTPAttemptsKey(id)).Err()
}

func (r *SignupContextRedisStore) GetLastOTPSentAt(ctx context.Context, id entity.SignupContextID) (time.Time, error) {

	val, err := r.client.Get(ctx, r.redisOTPSentKey(id)).Result()

	if err == redis.Nil {
		return time.Time{}, nil
	}

	if err != nil {
		return time.Time{}, err
	}

	parsed, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func (r *SignupContextRedisStore) SetLastOTPSentAt(ctx context.Context, id entity.SignupContextID, t time.Time) error {

	return r.client.Set(
		ctx,
		r.redisOTPSentKey(id),
		t.Format(time.RFC3339),
		signupContextTTL,
	).Err()
}

func (r *SignupContextRedisStore) Update(ctx context.Context, signup *entity.SignupContext) error {

	model := SignupContextModel{
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

	key := r.redisSignupContextKey(signup.ID)

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	if ttl <= 0 {
		ttl = signupContextTTL
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}
