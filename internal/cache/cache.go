package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/amirhossein-karimi/todo/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}

type cache struct {
	rdb *redis.Client
}

func NewRedisCache(config *config.Config) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
	return &cache{
		rdb: rdb,
	}
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("set cache key %q: %w", key, err)
	}

	return nil
}

func (c *cache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get cache key %q: %w", key, err)
	}

	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("unmarshal cache value: %w", err)
	}

	return value, nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete cache key %q: %w", key, err)
	}
	return nil
}
