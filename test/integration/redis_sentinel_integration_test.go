package integration

import (
	"testing"

	"github.com/Trendyol/go-dcp-redis/config"
)

func TestRedisSentinel(t *testing.T) {
	testRedisConnection(t, "config-sentinel.yml", config.Redis{
		DB: 1, // Use different database to avoid conflicts
		Sentinel: &config.RedisSentinel{
			MasterName: "mymaster",
			SentinelAddrs: []string{
				"localhost:26379",
				"localhost:26380",
				"localhost:26381",
			},
		},
	})
}
