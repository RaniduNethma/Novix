package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/novix/services/streaming-service/internal/config"
	"github.com/novix/services/streaming-service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected",
		zap.String("addr", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
	)

	return &Client{rdb: rdb}, nil
}

// Set stores a value with TTL
func (c *Client) Set(ctx context.Context, key string,
	value interface{}, ttl time.Duration) error {

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

// Get retrieves and unmarshals a value
func (c *Client) Get(ctx context.Context, key string,
	dest interface{}) error {

	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key: %w", err)
	}
	return count > 0, nil
}

// Keys returns all keys matching a pattern
func (c *Client) Keys(ctx context.Context,
	pattern string) ([]string, error) {

	keys, err := c.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	return keys, nil
}

// Increment increments a counter key
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key: %w", err)
	}
	return val, nil
}

// Expire sets TTL on existing key
func (c *Client) Expire(ctx context.Context,
	key string, ttl time.Duration) error {

	if err := c.rdb.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expiry: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}