package redis

import (
	// "context"
	// "time"

	"context"

	"github.com/go-redis/redis/v8"
)



type Redis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}


type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port, password string, database_number int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     host + ":" + port,
			Password: password,
			DB: database_number,
		}),
	}
}


func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
