package client

import (
	"context"
	"fmt"

	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Redis) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.Username != "" {
		options.Username = cfg.Username
	}

	client := redis.NewClient(options)

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
