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
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, pattern string) error
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

func (c *cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("get cache key %q: %w", key, err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("unmarshal cache key %q: %w", key, err)
	}

	return nil
}

func (c *cache) Delete(ctx context.Context, pattern string) error {
	var cursor uint64

	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("scan cache keys %q: %w", pattern, err)
		}

		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("delete cache keys %q: %w", pattern, err)
			}
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}
