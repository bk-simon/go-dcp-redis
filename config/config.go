package config

import (
	"time"

	"github.com/Trendyol/go-dcp/config"
)

type Redis struct {
	Host                 string                 `yaml:"host"`
	Port                 uint16                 `yaml:"port"`
	Username             string                 `yaml:"username"`
	Password             string                 `yaml:"password"`
	DB                   int                    `yaml:"db"`
	CollectionKeyMapping []CollectionKeyMapping `yaml:"collectionKeyMapping,omitempty"`
	BatchTickerDuration  time.Duration          `yaml:"batchTickerDuration"`
	DefaultTTL           time.Duration          `yaml:"defaultTTL"`
	Sentinel             *RedisSentinel         `yaml:"sentinel,omitempty"`
}

type RedisSentinel struct {
	MasterName    string   `yaml:"masterName"`
	SentinelAddrs []string `yaml:"sentinelAddrs"`
	Username      string   `yaml:"username,omitempty"`
	Password      string   `yaml:"password,omitempty"`
}

type CollectionKeyMapping struct {
	Collection  string        `yaml:"collection"`
	KeyPrefix   string        `yaml:"keyPrefix"`
	KeySuffix   string        `yaml:"keySuffix"`
	StorageType string        `yaml:"storageType"` // "string", "hash", "json"
	TTL         time.Duration `yaml:"ttl"`
	HashField   string        `yaml:"hashField,omitempty"`
}

type Connector struct {
	Redis Redis      `yaml:"redis" mapstructure:"redis"`
	Dcp   config.Dcp `yaml:",inline" mapstructure:",squash"`
}

func (c *Connector) ApplyDefaults() {
	if c.Redis.Port == 0 {
		c.Redis.Port = 6379
	}

	if c.Redis.BatchTickerDuration == 0 {
		c.Redis.BatchTickerDuration = 10 * time.Second
	}

	if c.Redis.DB == 0 {
		c.Redis.DB = 0
	}

	// Set default storage type for mappings
	for i := range c.Redis.CollectionKeyMapping {
		if c.Redis.CollectionKeyMapping[i].StorageType == "" {
			c.Redis.CollectionKeyMapping[i].StorageType = "string"
		}
	}
}
