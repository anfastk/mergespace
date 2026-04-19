package outbound

import (
	"context"

	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/jackc/pgx/v5"
)

type OutboxRepository interface {
	Save(ctx context.Context, tx pgx.Tx, e *entity.OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)
	MarkCompleted(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, errMsg string, retryCount int) error
	MarkDead(ctx context.Context, id string) error
}
