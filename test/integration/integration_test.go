package integration

import (
	"context"
	"testing"
	"time"

	dcpredis "github.com/Trendyol/go-dcp-redis"
	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/Trendyol/go-dcp-redis/redis/client"
	"github.com/Trendyol/go-dcp/logger"
)

type AirlineEvent struct {
	name string
}

func TestRedis(t *testing.T) {
	time.Sleep(time.Second * 30)

	connector, err := dcpredis.NewConnectorBuilder("config.yml").Build()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Start connector in a goroutine
	go connector.Start()

	// Wait a bit for connector to initialize
	time.Sleep(2 * time.Second)

	// Create Redis client for testing
	redisClient, err := client.NewRedisClient(config.Redis{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	})
	if err != nil {
		t.Fatalf("could not open connection to redis %s", err)
	}
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Check for keys periodically
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("deadline exceeded - not enough keys found")
		default:
			// Check if keys exist in Redis
			keys, err := redisClient.Keys(ctx, "doc:*").Result()
			if err != nil {
				t.Fatalf("redis query error %s", err)
			}

			logger.Log.Info("found %d keys", len(keys))

			if len(keys) >= 100 { // Check for at least 100 keys
				logger.Log.Info("test completed successfully - found %d keys", len(keys))
				connector.Close()
				return
			}
			time.Sleep(2 * time.Second)
		}
	}
}
