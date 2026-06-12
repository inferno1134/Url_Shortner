package store

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr string) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis ping failed: %w", err)
    }

    return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(key string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(key, value string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.client.Set(ctx, key, value, 24*time.Hour).Err()
}

func (c *RedisCache) Close() error {
    return c.client.Close()
}
