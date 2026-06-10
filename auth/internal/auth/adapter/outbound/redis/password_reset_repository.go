package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/valueobject"

	"github.com/redis/go-redis/v9"
)

type PasswordResetRedisStore struct {
	client *redis.Client
}

var _ outbound.PasswordResetStore = (*PasswordResetRedisStore)(nil)

func NewPasswordResetRedisStore(client *redis.Client) *PasswordResetRedisStore {
	return &PasswordResetRedisStore{client: client}
}

func (r *PasswordResetRedisStore) Save(ctx context.Context, reset *entity.PasswordResetContext) error {

	key := fmt.Sprintf(
		"password_reset:%s",
		reset.ID,
	)

	emailKey := fmt.Sprintf(
		"password_reset_email:%s",
		reset.Email,
	)

	data, err := json.Marshal(reset)
	if err != nil {
		return err
	}

	pipe := r.client.TxPipeline()

	pipe.Set(
		ctx,
		key,
		data,
		10*time.Minute,
	)

	pipe.Set(
		ctx,
		emailKey,
		string(reset.ID),
		10*time.Minute,
	)

	_, err = pipe.Exec(ctx)

	return err

}

func (r *PasswordResetRedisStore) FindByID(ctx context.Context, id entity.PasswordResetContextID) (*entity.PasswordResetContext, error) {

	key := fmt.Sprintf("password_reset:%s", id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var reset entity.PasswordResetContext

	if err := json.Unmarshal(data, &reset); err != nil {
		return nil, err
	}

	return &reset, nil
}

func (r *PasswordResetRedisStore) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.PasswordResetContext, error) {
	emailKey := fmt.Sprintf(
		"password_reset_email:%s",
		email.String(),
	)

	resetID, err := r.client.Get(
		ctx,
		emailKey,
	).Result()

	if err != nil {

		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return r.FindByID(
		ctx,
		entity.PasswordResetContextID(resetID),
	)

}

func (r *PasswordResetRedisStore) Update(ctx context.Context, reset *entity.PasswordResetContext) error {
	return r.Save(ctx, reset)
}

func (r *PasswordResetRedisStore) Delete(ctx context.Context, id entity.PasswordResetContextID) error {

	key := fmt.Sprintf("password_reset:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *PasswordResetRedisStore) GetAttempts(ctx context.Context, id entity.PasswordResetContextID) (int, error) {

	key := fmt.Sprintf("password_reset_attempts:%s", id)

	count, err := r.client.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (r *PasswordResetRedisStore) IncrementAttempts(ctx context.Context, id entity.PasswordResetContextID, ttl time.Duration) error {

	key := fmt.Sprintf("password_reset_attempts:%s", id)

	pipe := r.client.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *PasswordResetRedisStore) DeleteAttempts(ctx context.Context, id entity.PasswordResetContextID) error {

	key := fmt.Sprintf("password_reset_attempts:%s", id)
	return r.client.Del(ctx, key).Err()
}

func (r *PasswordResetRedisStore) SetLastOTPSentAt(ctx context.Context, id entity.PasswordResetContextID, t time.Time) error {

	key := fmt.Sprintf("password_reset_last_sent:%s", id)

	return r.client.Set(ctx, key, t.Unix(), 10*time.Minute).Err()
}

func (r *PasswordResetRedisStore) GetLastOTPSentAt(ctx context.Context, id entity.PasswordResetContextID) (time.Time, error) {

	key := fmt.Sprintf("password_reset_last_sent:%s", id)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return time.Time{}, err
	}

	unixTime, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixTime, 0), nil
}
