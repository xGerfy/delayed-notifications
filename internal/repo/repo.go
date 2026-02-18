package repo

import (
	"context"
	"delayed/internal/entities"
	"encoding/json"

	"github.com/wb-go/wbf/redis"
)

type repository struct {
	client *redis.Client
}

func New(addr, password string, db int) *repository {
	return &repository{
		client: redis.New(addr, password, db),
	}
}

func (r *repository) Save(ctx context.Context, n *entities.Notification) error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, n.ID, data)
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Notification, error) {
	data, err := r.client.Get(ctx, id)
	if err == redis.NoMatches {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var n entities.Notification
	if err := json.Unmarshal([]byte(data), &n); err != nil {
		return nil, err
	}

	return &n, nil
}
