package config

import (
	"testing"
	"time"
)

func TestConnector_ApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    Connector
		expected Connector
	}{
		{
			name:  "empty config gets defaults",
			input: Connector{Redis: Redis{}},
			expected: Connector{
				Redis: Redis{
					Port:                6379,
					BatchTickerDuration: 10 * time.Second,
				},
			},
		},
		{
			name:  "custom port preserved",
			input: Connector{Redis: Redis{Port: 6380}},
			expected: Connector{
				Redis: Redis{
					Port:                6380,
					BatchTickerDuration: 10 * time.Second,
				},
			},
		},
		{
			name: "collection key mapping storage type defaults to string",
			input: Connector{
				Redis: Redis{
					CollectionKeyMapping: []CollectionKeyMapping{{Collection: "test"}},
				},
			},
			expected: Connector{
				Redis: Redis{
					Port:                6379,
					BatchTickerDuration: 10 * time.Second,
					CollectionKeyMapping: []CollectionKeyMapping{
						{Collection: "test", StorageType: "string"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.ApplyDefaults()

			if tt.input.Redis.Port != tt.expected.Redis.Port {
				t.Errorf("Port = %d, want %d", tt.input.Redis.Port, tt.expected.Redis.Port)
			}

			if tt.input.Redis.BatchTickerDuration != tt.expected.Redis.BatchTickerDuration {
				t.Errorf("BatchTickerDuration = %v, want %v", tt.input.Redis.BatchTickerDuration, tt.expected.Redis.BatchTickerDuration)
			}

			if len(tt.input.Redis.CollectionKeyMapping) > 0 {
				actual := tt.input.Redis.CollectionKeyMapping[0].StorageType
				expected := tt.expected.Redis.CollectionKeyMapping[0].StorageType
				if actual != expected {
					t.Errorf("StorageType = %s, want %s", actual, expected)
				}
			}
		})
	}
}

func TestRedisClusterConfig(t *testing.T) {
	t.Run("cluster config with addrs", func(t *testing.T) {
		cfg := Redis{
			Cluster: &RedisCluster{
				Addrs: []string{
					"localhost:7000",
					"localhost:7001",
					"localhost:7002",
				},
				RouteByLatency: true,
				RouteRandomly:  false,
				ReadOnly:       true,
			},
		}

		if len(cfg.Cluster.Addrs) != 3 {
			t.Errorf("Addrs count = %d, want 3", len(cfg.Cluster.Addrs))
		}

		if !cfg.Cluster.RouteByLatency {
			t.Error("RouteByLatency should be true")
		}

		if cfg.Cluster.RouteRandomly {
			t.Error("RouteRandomly should be false")
		}

		if !cfg.Cluster.ReadOnly {
			t.Error("ReadOnly should be true")
		}
	})

	t.Run("cluster config with credentials", func(t *testing.T) {
		cfg := Redis{
			Username: "default_user",
			Password: "default_pass",
			Cluster: &RedisCluster{
				Addrs:    []string{"localhost:7000"},
				Username: "cluster_user",
				Password: "cluster_pass",
			},
		}

		// Cluster credentials should take priority
		if cfg.Cluster.Username != "cluster_user" {
			t.Errorf("Cluster Username = %s, want cluster_user", cfg.Cluster.Username)
		}

		if cfg.Cluster.Password != "cluster_pass" {
			t.Errorf("Cluster Password = %s, want cluster_pass", cfg.Cluster.Password)
		}
	})
}

func TestRedisSentinelConfig(t *testing.T) {
	t.Run("sentinel config structure", func(t *testing.T) {
		cfg := Redis{
			Sentinel: &RedisSentinel{
				MasterName: "mymaster",
				SentinelAddrs: []string{
					"localhost:26379",
					"localhost:26380",
				},
				Username: "sentinel_user",
				Password: "sentinel_pass",
			},
		}

		if cfg.Sentinel.MasterName != "mymaster" {
			t.Errorf("MasterName = %s, want mymaster", cfg.Sentinel.MasterName)
		}

		if len(cfg.Sentinel.SentinelAddrs) != 2 {
			t.Errorf("SentinelAddrs count = %d, want 2", len(cfg.Sentinel.SentinelAddrs))
		}
	})
}
