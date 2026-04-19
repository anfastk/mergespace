package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
)

type OutboxWorker struct {
	repo          outbound.OutboxRepository
	eventProducer outbound.EventProducer
}

func NewOutboxWorker(repo outbound.OutboxRepository, producer outbound.EventProducer) *OutboxWorker {
	return &OutboxWorker{
		repo:          repo,
		eventProducer: producer,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			w.process(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *OutboxWorker) process(ctx context.Context) {
	events, err := w.repo.FetchPending(ctx, 10)
	if err != nil {
		log.Println("fetch error:", err)
		return
	}

	for _, e := range events {

		if e.RetryCount >= 5 {
			log.Printf("event moved to dead (id=%s)", e.ID)

			if err := w.repo.MarkDead(ctx, e.ID); err != nil {
				log.Printf("failed to mark dead: %v", err)
			}
			continue
		}

		if err := w.handle(ctx, e); err != nil {

			log.Printf("event failed (id=%s): %v", e.ID, err)

			if err := w.repo.MarkFailed(ctx, e.ID, err.Error(), e.RetryCount+1); err != nil {
				log.Printf("failed to update retry: %v", err)
			}

			continue
		}

		if err := w.repo.MarkCompleted(ctx, e.ID); err != nil {
			log.Printf("failed to mark completed: %v", err)
		}
	}
}

func (w *OutboxWorker) handle(ctx context.Context, e *entity.OutboxEvent) error {

	if e.EventType != "SendOTP" {
		return nil
	}

	var payload map[string]string
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return err
	}

	return w.eventProducer.PublishSendOTP(ctx, &event.SendOTP{
		Email: payload["email"],
		OTP:   payload["otp"],
	})
}
