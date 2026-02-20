package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         string    `json:"id"`
	Message    string    `json:"message"`
	SendAt     time.Time `json:"send_at"`
	Status     Status    `json:"status"`
	RetryCount int       `json:"retry_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

func New(message string, delay time.Duration) *Notification {
	now := time.Now()
	return &Notification{
		ID:         uuid.New().String(),
		Message:    message,
		SendAt:     now.Add(delay),
		Status:     StatusPending,
		RetryCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
