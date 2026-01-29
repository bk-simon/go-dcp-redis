package integration

import (
	"testing"

	"github.com/Trendyol/go-dcp-redis/config"
)

func TestRedisCluster(t *testing.T) {
	testRedisConnection(t, "config-cluster.yml", config.Redis{
		DB: 0,
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
	})
}
