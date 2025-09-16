package main

import (
	dcpredis "github.com/Trendyol/go-dcp-redis"
	"github.com/Trendyol/go-dcp-redis/couchbase"
	"github.com/Trendyol/go-dcp-redis/redis"
)

func mapper(event couchbase.Context) []redis.Model {
	if event.Event.IsMutated {
		return []redis.Model{
			&redis.Set{
				Key:   string(event.Event.Key),
				Value: event.Event.Value,
			},
		}
	}

	return []redis.Model{
		&redis.Del{
			Key: string(event.Event.Key),
		},
	}
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
