package publisher

import (
	"context"
	"delayed/internal/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type publisher struct {
	client *rabbitmq.RabbitClient
	pub    *rabbitmq.Publisher
	cfgApp *config.AppConfig
}

type message struct {
	ID string
}

func New(cfgApp *config.AppConfig) (*publisher, error) {
	cfg := rabbitmq.ClientConfig{
		URL: cfgApp.RabbitConfig.Url,
		ReconnectStrat: retry.Strategy{
			Attempts: 3,
			Delay:    5,
			Backoff:  1,
		},
		ProducingStrat: retry.Strategy{
			Attempts: 3,
			Delay:    5,
			Backoff:  1,
		},
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
	}

	err = client.DeclareExchange(
		cfgApp.RabbitConfig.Exchange,
		"x-delayed-message",
		true,
		true,
		false,
		amqp091.Table{"x-delayed-type": "direct"},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchenge: %w", err)
	}

	err = client.DeclareQueue(
		cfgApp.RabbitConfig.Queue,
		cfgApp.RabbitConfig.Exchange,
		cfgApp.RabbitConfig.Queue,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	pub := rabbitmq.NewPublisher(client, cfgApp.RabbitConfig.Exchange, "application/json")

	return &publisher{
		client: client,
		pub:    pub,
		cfgApp: cfgApp,
	}, nil
}

func (p *publisher) PublishWithDelay(ctx context.Context, id string, delay time.Duration) error {
	msg := message{ID: id}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	delayMs := delay.Milliseconds()
	headers := amqp091.Table{"x-delay": delayMs}

	return p.pub.Publish(ctx, body, p.cfgApp.RabbitConfig.Queue, rabbitmq.WithHeaders(headers))
}

func (p *publisher) Close() error {
	return p.client.Close()
}
