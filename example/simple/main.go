package main

import (
	"fmt"
	"time"

	dcpredis "github.com/Trendyol/go-dcp-redis"
	"github.com/Trendyol/go-dcp-redis/couchbase"
	"github.com/Trendyol/go-dcp-redis/redis"
)

func mapper(event couchbase.Event) []redis.Model {
	var set = redis.Set{
		Key:   fmt.Sprintf("doc:%s", string(event.Key)),
		Value: string(event.Value),
		TTL:   time.Hour * 24, // 24 hours TTL
	}

	return []redis.Model{&set}
}

func main() {
	connector, err := dcpredis.NewConnectorBuilder("config.yml").
		SetMapper(mapper).
		Build()
	if err != nil {
		panic(err)
	}

	defer connector.Close()
	connector.Start()
}
