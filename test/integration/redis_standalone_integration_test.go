package integration

import (
	"testing"

	"github.com/Trendyol/go-dcp-redis/config"
)

func TestRedisStandalone(t *testing.T) {
	testRedisConnection(t, "config.yml", config.Redis{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	})
}
