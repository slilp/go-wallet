package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=./redis_adapter.go -destination=./mocks/mock_redis_adapter.go -package=mock_adapters
type RedisAdapter interface {
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
}

type redisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) RedisAdapter {
	return &redisAdapter{client: client}
}

func (r *redisAdapter) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *redisAdapter) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
