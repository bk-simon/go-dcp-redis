package client

import (
	"context"
	"fmt"

	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Redis) (*redis.Client, error) {
	var client *redis.Client

	// Check if Sentinel configuration is provided
	if cfg.Sentinel != nil && len(cfg.Sentinel.SentinelAddrs) > 0 {
		client = newSentinelClient(cfg)
	} else {
		client = newStandaloneClient(cfg)
	}

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func newStandaloneClient(cfg config.Redis) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.Username != "" {
		options.Username = cfg.Username
	}

	return redis.NewClient(options)
}

func newSentinelClient(cfg config.Redis) *redis.Client {
	options := &redis.FailoverOptions{
		MasterName:    cfg.Sentinel.MasterName,
		SentinelAddrs: cfg.Sentinel.SentinelAddrs,
		DB:            cfg.DB,
	}

	// Use Redis config credentials if Sentinel specific ones are not provided
	if cfg.Sentinel.Username != "" {
		options.Username = cfg.Sentinel.Username
	} else if cfg.Username != "" {
		options.Username = cfg.Username
	}

	if cfg.Sentinel.Password != "" {
		options.Password = cfg.Sentinel.Password
	} else if cfg.Password != "" {
		options.Password = cfg.Password
	}

	return redis.NewFailoverClient(options)
}
