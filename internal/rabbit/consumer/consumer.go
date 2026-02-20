package consumer

import (
	"context"
	"delayed/internal/config"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type service interface {
	Process(ctx context.Context, id string) error
}

type consumer struct {
	consumer *rabbitmq.Consumer
}

type Message struct {
	ID string `json:"id"`
}

func New(cfgApp *config.AppConfig, service service) (*consumer, error) {
	cfg := rabbitmq.ClientConfig{
		URL: cfgApp.RabbitConfig.Url,
		ReconnectStrat: retry.Strategy{
			Attempts: 3,
			Delay:    5,
			Backoff:  1,
		},
		ConsumingStrat: retry.Strategy{
			Attempts: 3,
			Delay:    5,
			Backoff:  1,
		},
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
	}

	cfgConsumer := rabbitmq.ConsumerConfig{
		Queue:         cfgApp.RabbitConfig.Queue,
		ConsumerTag:   "notification_consumer",
		AutoAck:       false,
		Workers:       1,
		PrefetchCount: 1,
		Nack:          rabbitmq.NackConfig{Multiple: false, Requeue: true},
		Ask:           rabbitmq.AskConfig{Multiple: false},
	}

	handler := func(ctx context.Context, msg amqp091.Delivery) error {
		return handleMessage(ctx, msg, service)
	}

	cons := rabbitmq.NewConsumer(client, cfgConsumer, handler)

	return &consumer{
		consumer: cons,
	}, nil
}

func (c *consumer) Start(ctx context.Context) error {
	return c.consumer.Start(ctx)
}

func handleMessage(ctx context.Context, msg amqp091.Delivery, service service) error {
	var message Message
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if message.ID == "" {
		return nil // возвращаем nil, чтобы подтвердить сообщение (или не возвращать в очередь)
	}

	// Вызываем сервис для обработки
	err := service.Process(ctx, message.ID)
	if err != nil {
		return err // ошибка будет обработана оберткой (ретранслируется)
	}

	return nil
}
