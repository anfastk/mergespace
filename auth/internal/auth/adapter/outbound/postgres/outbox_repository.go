package postgres

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepo struct {
	db *pgxpool.Pool
}

func NewOutboxRepo(db *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{db: db}
}

var _ outbound.OutboxRepository = (*OutboxRepo)(nil)

func (r *OutboxRepo) Save(ctx context.Context, tx pgx.Tx, e *entity.OutboxEvent) error {

	query := `
		INSERT INTO outbox_events (id, event_type, payload, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query,
		e.ID,
		e.EventType,
		e.Payload,
		e.Status,
	)

	return err
}

func (r *OutboxRepo) FetchPending(ctx context.Context, limit int) ([]*entity.OutboxEvent, error) {
	query := `
		SELECT id, event_type, payload, status, retry_count, next_retry_at
		FROM outbox_events
		WHERE status = 'pending'
  			AND next_retry_at <= NOW()
		ORDER BY created_at
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.OutboxEvent

	for rows.Next() {
		e := &entity.OutboxEvent{}
		if err := rows.Scan(
			&e.ID,
			&e.EventType,
			&e.Payload,
			&e.Status,
			&e.RetryCount,
			&e.NextRetryAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *OutboxRepo) MarkCompleted(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE outbox_events SET status = 'completed' WHERE id=$1`, id,
	)
	return err
}

func (r *OutboxRepo) MarkFailed(ctx context.Context, id string, errMsg string, retryCount int) error {
	base := 2 * time.Second
	max := 1 * time.Minute

	exp := time.Duration(math.Pow(2, float64(retryCount))) * base
	if exp > max {
		exp = max
	}

	delay := time.Duration(rand.Int63n(int64(exp)))

	nextRetry := time.Now().Add(delay)

	query := `
		UPDATE outbox_events
		SET 
			retry_count = $2,
			last_error = $3,
			next_retry_at = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, retryCount, errMsg, nextRetry)
	return err
}

func (r *OutboxRepo) MarkDead(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE outbox_events SET status='failed' WHERE id=$1`,
		id,
	)
	return err
}
