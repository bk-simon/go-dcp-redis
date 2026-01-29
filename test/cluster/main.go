package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/Trendyol/go-dcp-redis/redis/client"
)

func main() {
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

	redisClient, err := client.NewRedisClient(cfg)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer redisClient.Close()

	ctx := context.Background()

	// Test SET
	err = redisClient.Set(ctx, "test:cluster:key1", "Hello Redis Cluster!", time.Hour).Err()
	if err != nil {
		fmt.Printf("SET failed: %v\n", err)
		return
	}
	fmt.Println("SET: test:cluster:key1 = 'Hello Redis Cluster!'")

	// Test GET
	val, err := redisClient.Get(ctx, "test:cluster:key1").Result()
	if err != nil {
		fmt.Printf("GET failed: %v\n", err)
		return
	}
	fmt.Printf("GET: test:cluster:key1 = '%s'\n", val)

	// Test multiple keys (will be distributed across nodes)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:cluster:key%d", i)
		err = redisClient.Set(ctx, key, fmt.Sprintf("value%d", i), time.Hour).Err()
		if err != nil {
			fmt.Printf("SET failed for %s: %v\n", key, err)
			return
		}
	}
	fmt.Println("SET: 10 keys successfully stored across cluster nodes")

	// Cleanup
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:cluster:key%d", i)
		redisClient.Del(ctx, key)
	}
	fmt.Println("Cleanup: All test keys deleted")

	fmt.Println("\n✅ Redis Cluster connection test PASSED!")
}
