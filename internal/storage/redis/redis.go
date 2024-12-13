package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)


// TODO: сейчас это просто обертка, возможно нужно добавить другие методы, связанные с нашей логикой

type Redis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Ping(ctx context.Context) error
	MSet(ctx context.Context, values map[string]interface{}) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port, password string, databaseNumber int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     host + ":" + port,
			Password: password,
			DB:       databaseNumber,
		}),
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *RedisClient) MSet(ctx context.Context, values map[string]interface{}) error {
	return r.client.MSet(ctx, values).Err()
}

func (r *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}
