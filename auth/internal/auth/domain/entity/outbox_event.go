package entity

import "time"

type OutboxEvent struct {
	ID          string
	EventType   string
	Payload     []byte
	Status      string
	RetryCount  int
	NextRetryAt time.Time
	LastError   *string
}
