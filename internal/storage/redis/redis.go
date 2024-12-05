package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)



type Redis interface {
	
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

func (r *RedisClient) SetSession(ctx context.Context, sessionID string, data string, expiration time.Duration) error {
    return r.client.Set(ctx, sessionID, data, expiration).Err()
}

func (r *RedisClient) GetSession(ctx context.Context, sessionID string) (string, error) {
    return r.client.Get(ctx, sessionID).Result()
}

func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
    return r.client.Del(ctx, sessionID).Err()
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}