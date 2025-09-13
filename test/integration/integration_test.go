package integration

import (
	"context"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Start()
	}()

	time.Sleep(1 * time.Second)

	go func() {
		redisClient, err := client.NewRedisClient(config.Redis{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		})
		if err != nil {
			t.Fatalf("could not open connection to redis %s", err)
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)

	CountCheckLoop:
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("deadline exceed")
			default:
				// Check if keys exist in Redis
				keys, err := redisClient.Keys(ctx, "doc:*").Result()
				if err != nil {
					t.Fatalf("redis query error %s", err)
				}
				if len(keys) >= 100 { // Check for at least 100 keys
					logger.Log.Info("done")
					connector.Close()
					break CountCheckLoop
				}
				time.Sleep(2 * time.Second)
			}
		}

	}()

	wg.Wait()
}
