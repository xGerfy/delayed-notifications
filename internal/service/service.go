package service

import (
	"context"
	"delayed/internal/entities"
	"errors"
	"fmt"
	"time"
)

type repository interface {
	Save(ctx context.Context, n *entities.Notification) error
	Get(ctx context.Context, id string) (*entities.Notification, error)
}

type publisher interface {
	PublishWithDelay(ctx context.Context, id string, delay time.Duration) error
}

type service struct {
	repo      repository
	publisher publisher
}

func New(repo repository, publisher publisher) *service {
	return &service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *service) Create(ctx context.Context, message string, sendAt time.Time) (*entities.Notification, error) {
	if sendAt.Before(time.Now()) {
		return nil, errors.New("send time must be in the future")
	}

	delay := time.Until(sendAt)
	n := entities.New(message, delay)

	if err := s.repo.Save(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	if err := s.publisher.PublishWithDelay(ctx, n.ID, delay); err != nil {
		return nil, fmt.Errorf("failed to publish notification: %w", err)
	}

	return n, nil
}

func (s *service) Get(ctx context.Context, id string) (*entities.Notification, error) {
	n, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return n, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	n, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}
	if n.Status == entities.StatusSent {
		return errors.New("cannot delete sent notification")
	}

	n.Status = entities.StatusCancelled
	n.UpdatedAt = time.Now()

	if err = s.repo.Save(ctx, n); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

func (s *service) Process(ctx context.Context, id string) error {
	// Получаем уведомление
	n, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	// Проверяем статус
	if n.Status != entities.StatusPending {
		return fmt.Errorf("notification is not pending, current status: %s", n.Status)
	}

	// Пытаемся отправить уведомление
	if err := s.deliverNotification(n); err != nil {
		// Обработка ошибки с ретраями
		return s.handleDeliveryError(ctx, n, err)
	}

	// Успешная отправка
	n.Status = entities.StatusSent
	n.UpdatedAt = time.Now()

	return s.repo.Save(ctx, n)
}

func (s *service) deliverNotification(n *entities.Notification) error {
	// TODO: Реализовать реальную отправку (email, push, etc.)
	fmt.Printf("Sending notification %s: %s\n", n.ID, n.Message)

	// Имитация ошибки для тестирования ретраев
	// return errors.New("delivery failed")

	return nil
}

func (s *service) handleDeliveryError(ctx context.Context, n *entities.Notification, err error) error {
	const maxRetries = 5

	n.RetryCount++
	n.UpdatedAt = time.Now()

	if n.RetryCount >= maxRetries {
		// Превышено максимальное количество попыток
		n.Status = entities.StatusFailed
		updateErr := s.repo.Save(ctx, n)
		if updateErr != nil {
			return fmt.Errorf("delivery failed after %d retries and status update failed: %v, original error: %w",
				maxRetries, updateErr, err)
		}
		return fmt.Errorf("delivery failed after %d retries: %w", maxRetries, err)
	}

	// Экспоненциальная задержка: 2^retry * baseDelay
	baseDelay := time.Second * 5
	delay := baseDelay * time.Duration(1<<uint(n.RetryCount-1))

	fmt.Printf("Retry %d for notification %s after %v\n", n.RetryCount, n.ID, delay)

	// Обновляем статус в репозитории
	if updateErr := s.repo.Save(ctx, n); updateErr != nil {
		return fmt.Errorf("failed to update retry count: %w", updateErr)
	}

	// Возвращаем в очередь с задержкой
	return s.publisher.PublishWithDelay(ctx, n.ID, delay)
}
