package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ErrorCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
}
type RedisClient struct {
	client *redis.Client
}

func NewClient(Addr string, UserName string, PassWord string) Client {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     Addr,
			Username: UserName,
			Password: PassWord,
			DB:       0,
		}),
	}

}
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()

}
func (c *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorCacheMiss
		}
		return nil, err
	}
	return data, nil
}
