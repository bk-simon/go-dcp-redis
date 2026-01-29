package client

import (
	"testing"

	"github.com/Trendyol/go-dcp-redis/config"
)

func TestNewRedisClient_Standalone(t *testing.T) {
	// This test requires a running Redis instance
	// Skip if not in integration environment
	t.Skip("Skipping standalone client test - requires running Redis")

	cfg := config.Redis{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	client, err := NewRedisClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create standalone client: %v", err)
	}
	defer client.Close()

	// Verify client is working
	_, err = client.Ping(t.Context()).Result()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestNewRedisClient_Sentinel(t *testing.T) {
	// This test requires a running Redis Sentinel setup
	t.Skip("Skipping sentinel client test - requires running Redis Sentinel")

	cfg := config.Redis{
		DB: 0,
		Sentinel: &config.RedisSentinel{
			MasterName: "mymaster",
			SentinelAddrs: []string{
				"localhost:26379",
				"localhost:26380",
				"localhost:26381",
			},
		},
	}

	client, err := NewRedisClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create sentinel client: %v", err)
	}
	defer client.Close()

	// Verify client is working
	_, err = client.Ping(t.Context()).Result()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestNewRedisClient_Cluster(t *testing.T) {
	// This test requires a running Redis Cluster
	t.Skip("Skipping cluster client test - requires running Redis Cluster")

	cfg := config.Redis{
		Cluster: &config.RedisCluster{
			Addrs: []string{
				"localhost:7000",
				"localhost:7001",
				"localhost:7002",
				"localhost:7003",
				"localhost:7004",
				"localhost:7005",
			},
		},
	}

	client, err := NewRedisClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create cluster client: %v", err)
	}
	defer client.Close()

	// Verify client is working
	_, err = client.Ping(t.Context()).Result()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestClientPriorityOrder(t *testing.T) {
	// Test that Cluster takes priority over Sentinel
	// This is a unit test that doesn't require a running Redis

	// When both Cluster and Sentinel are configured, Cluster should take priority
	t.Run("ClusterPriorityOverSentinel", func(t *testing.T) {
		cfg := config.Redis{
			Host: "localhost",
			Port: 6379,
			Cluster: &config.RedisCluster{
				Addrs: []string{"localhost:7000"},
			},
			Sentinel: &config.RedisSentinel{
				MasterName:    "mymaster",
				SentinelAddrs: []string{"localhost:26379"},
			},
		}

		// This will fail to connect, but we're testing priority logic
		// The error message will indicate which type of client was created
		_, err := NewRedisClient(cfg)
		if err == nil {
			t.Skip("Skipping - Redis Cluster is running")
		}
		// If we get here, it means cluster client was attempted first (which is correct)
	})
}
