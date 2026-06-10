package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/anfastk/mergespace/auth/internal/auth/application/event"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/outbound"
	"github.com/anfastk/mergespace/auth/internal/auth/domain/entity"
	"github.com/anfastk/mergespace/platform/events"
	platformEvents "github.com/anfastk/mergespace/platform/events"
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

func (w *OutboxWorker) handle(
	ctx context.Context,
	e *entity.OutboxEvent,
) error {

	log.Println(
		"PUBLISHING EVENT:",
		e.EventType,
	)

	switch e.EventType {

	case platformEvents.EventSendOTP:

		var payload event.SendOTP

		if err := json.Unmarshal(
			e.Payload,
			&payload,
		); err != nil {
			return err
		}

		return w.eventProducer.Publish(
			ctx,
			platformEvents.EventSendOTP,
			[]byte(payload.Email),
			payload,
		)

	case platformEvents.EventUserCreated:

		var payload event.UserCreated

		if err := json.Unmarshal(
			e.Payload,
			&payload,
		); err != nil {
			return err
		}

		return w.eventProducer.Publish(
			ctx,
			platformEvents.EventUserCreated,
			[]byte(payload.UserID),
			payload,
		)

	case events.EventForgotPasswordOTP:

		var payload event.ForgotPasswordOTP

		if err := json.Unmarshal(
			e.Payload,
			&payload,
		); err != nil {
			return err
		}

		return w.eventProducer.Publish(
			ctx,
			events.EventForgotPasswordOTP,
			[]byte(payload.Email),
			payload,
		)
	}

	log.Println(
		"UNKNOWN EVENT:",
		e.EventType,
	)

	return nil
}
