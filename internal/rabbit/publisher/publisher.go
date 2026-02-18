package publisher

import (
	"context"
	"delayed/internal/entities"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

type publisher struct {
	client *rabbitmq.RabbitClient
	pub    *rabbitmq.Publisher
}

var (
	exchange = "delayed_exchange"
	queue    = "notifications_queue"
)

func New(url string) (*publisher, error) {
	cfg := rabbitmq.ClientConfig{
		URL:            url,
		ConnectTimeout: 0,
		Heartbeat:      0,
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
		ConsumingStrat: retry.Strategy{
			Attempts: 3,
			Delay:    5,
			Backoff:  1,
		},
	}

	client, err := rabbitmq.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = client.DeclareExchange(
		exchange,
		"x-delayed-message",
		true,
		true,
		false,
		amqp091.Table{"x-delayed-type": "direct"},
	)
	if err != nil {
		client.Close()
		return nil, err
	}

	err = client.DeclareQueue(
		queue,
		exchange,
		queue,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		client.Close()
		return nil, err
	}

	pub := rabbitmq.NewPublisher(client, exchange, "application/json")

	return &publisher{
		client: client,
		pub:    pub,
	}, nil
}

type msg struct {
	ID string
}

func (p *publisher) PublishWithDelay(ctx context.Context, msg entities.NotificationQueueMsg, delay time.Duration) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	opts := []rabbitmq.PublishOption{
		rabbitmq.WithHeaders(amqp091.Table{
			"x-delay": int(delay.Milliseconds()),
		}),
	}

	return p.pub.Publish(ctx, body, queue, opts...)
}

func (p *publisher) Close() error {
	return p.client.Close()
}
