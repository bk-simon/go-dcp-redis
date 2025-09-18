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

func testRedisConnection(t *testing.T, configFile string, redisConfig config.Redis) {
	time.Sleep(time.Second * 30)

	connector, err := dcpredis.NewConnectorBuilder(configFile).Build()
	if err != nil {
		t.Fatal(err)
		return
	}

	// Start connector in a goroutine
	go connector.Start()

	// Wait a bit for connector to initialize
	time.Sleep(2 * time.Second)

	// Create Redis client for testing
	redisClient, err := client.NewRedisClient(redisConfig)
	if err != nil {
		t.Fatalf("could not open connection to redis %s", err)
	}

	// Ensure cleanup happens even if test fails
	defer func() {
		logger.Log.Info("starting cleanup process")
		connector.Close()
		time.Sleep(3 * time.Second) // Wait for connector to fully stop
		redisClient.Close()
		logger.Log.Info("cleanup completed")
	}()

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
				return // defer will handle cleanup
			}
			time.Sleep(2 * time.Second)
		}
	}
}
