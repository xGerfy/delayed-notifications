package service

import (
	"context"
	"delayed/internal/entities"
	"time"
)

type Repository interface {
	Save(ctx context.Context, n *entities.Notification) error
	Get(ctx context.Context, id string) (*entities.Notification, error)
}

type Queue interface {
	PublishWithDelay(ctx context.Context, msg entities.NotificationQueueMsg, delay time.Duration) error
}

type service struct {
	repo  Repository
	queue Queue
}

func New(repo Repository, queue Queue) *service {
	return &service{
		repo:  repo,
		queue: queue,
	}
}

func (s *service) Create(ctx context.Context, message string) (*entities.Notification, error) {
	n := entities.New(message)

	if err := s.repo.Save(ctx, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (s *service) Get(ctx context.Context, id string) (*entities.Notification, error) {
	return s.repo.Get(ctx, id)
}
