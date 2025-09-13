# Go Dcp Redis [![Go Reference](https://pkg.go.dev/badge/github.com/Trendyol/go-dcp-redis.svg)](https://pkg.go.dev/github.com/Trendyol/go-dcp-redis) [![Go Report Card](https://goreportcard.com/badge/github.com/Trendyol/go-dcp-redis)](https://goreportcard.com/report/github.com/Trendyol/go-dcp-redis) [![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Trendyol/go-dcp-redis/badge)](https://scorecard.dev/viewer/?uri=github.com/Trendyol/go-dcp-redis)

**Go Dcp Redis** streams documents from Couchbase Database Change Protocol (DCP) and writes to
Redis in near real-time.

## Features

* Custom Redis commands **per** DCP event.
* **Multiple Redis operations** for a DCP event(see [Example](#example)).
* Handling different DCP events such as **expiration, deletion and mutation**(see [Example](#example)).
* **Managing batch configurations** such as batch ticker durations.
* **Multiple storage types** (string, hash, json) with TTL support.
* **Scale up and down** by custom membership algorithms(Couchbase, KubernetesHa, Kubernetes StatefulSet or
  Static, see [examples](https://github.com/Trendyol/go-dcp#examples)).
* **Easily manageable configurations**.

## Example

**Note:** If you prefer to use the default mapper by entering the configuration instead of creating a custom mapper, please refer to [this](#collection-key-mapping-configuration) topic.
Otherwise, you can refer to the example provided below:

```go
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
    SetMapper(mapper). // NOT NEEDED IF YOU'RE USING DEFAULT MAPPER. JUST CALL Build() FUNCTION
    Build()
  if err != nil {
    panic(err)
  }

  defer connector.Close()
  connector.Start()
}
```

### Advanced Redis Usage Examples

**Hash Storage Example:**
```go
func hashMapper(event couchbase.Event) []redis.Model {
  return []redis.Model{
    &redis.HSet{
      Key:   fmt.Sprintf("user:%s", string(event.Key)),
      Field: "profile",
      Value: string(event.Value),
      TTL:   time.Hour * 48,
    },
  }
}
```

**Multiple Operations Example:**
```go
func multiMapper(event couchbase.Event) []redis.Model {
  key := string(event.Key)
  value := string(event.Value)
  
  return []redis.Model{
    // Store as string
    &redis.Set{
      Key:   fmt.Sprintf("doc:%s", key),
      Value: value,
      TTL:   time.Hour * 24,
    },
    // Store in hash for indexing
    &redis.HSet{
      Key:   "docs:index",
      Field: key,
      Value: time.Now().Unix(),
    },
  }
}
```

**Raw Command Example:**
```go
func rawMapper(event couchbase.Event) []redis.Model {
  return []redis.Model{
    &redis.Raw{
      Operation: "ZADD",
      Key:       "docs:timeline",
      Args:      []interface{}{time.Now().Unix(), string(event.Key)},
    },
  }
}
```

## Configuration

### Dcp Configuration

Check out on [go-dcp](https://github.com/Trendyol/go-dcp#configuration)

### Redis Specific Configuration

| Variable                                                | Type                     | Required | Default | Description                                                                                        |                                                           
|---------------------------------------------------------|--------------------------|----------|---------|----------------------------------------------------------------------------------------------------|
| `redis.host`                                            | string                   | yes      |         | Redis host                                                                                         |
| `redis.port`                                            | int                      | no       | 6379    | Redis port                                                                                         |
| `redis.username`                                        | string                   | no       |         | Redis username (Redis 6.0+)                                                                       |
| `redis.password`                                        | string                   | no       |         | Redis password                                                                                     |
| `redis.db`                                              | int                      | no       | 0       | Redis database number                                                                              |
| `redis.defaultTTL`                                      | time.Duration            | no       |         | Default TTL for keys                                                                               |
| `redis.batchTickerDuration`                             | time.Duration            | no       | 10s     | Batch is being flushed automatically at specific time intervals for long waiting messages in batch |
| `redis.collectionKeyMapping`                            | []CollectionKeyMapping   | no       |         | Will be used for default mapper. Please read the next topic.                                       |

### Collection Key Mapping Configuration

Collection key mapping configuration is optional. This configuration should only be provided if you are using the default mapper. If you are implementing your own custom mapper function, this configuration is not needed.

| Variable                                                 | Type          | Required | Default | Description                                                                  |                                                           
|----------------------------------------------------------|---------------|----------|---------|------------------------------------------------------------------------------|
| `redis.collectionKeyMapping[].collection`               | string        | yes      |         | Couchbase collection name                                                    |
| `redis.collectionKeyMapping[].keyPrefix`                 | string        | no       |         | Prefix to add to Redis keys                                                  |
| `redis.collectionKeyMapping[].keySuffix`                 | string        | no       |         | Suffix to add to Redis keys                                                  |
| `redis.collectionKeyMapping[].storageType`               | string        | no       | string  | Storage type: "string", "hash", "json"                                       |
| `redis.collectionKeyMapping[].ttl`                       | time.Duration | no       |         | TTL for this collection's keys                                               |
| `redis.collectionKeyMapping[].hashField`                 | string        | no       | value   | Field name when using hash storage type                                      |

## Exposed metrics

| Metric Name                                     | Description                   | Labels | Value Type |
|-------------------------------------------------|-------------------------------|--------|------------|
| redis_connector_latency_ms                      | Time to adding to the batch.  | N/A    | Gauge      |
| redis_connector_bulk_request_process_latency_ms | Time to process bulk request. | N/A    | Gauge      |

You can also use all DCP-related metrics explained [here](https://github.com/Trendyol/go-dcp#exposed-metrics).
All DCP-related metrics are automatically injected. It means you don't need to do anything. 

## Contributing

Go Dcp Redis is always open for direct contributions. For more information please check
our [Contribution Guideline document](./CONTRIBUTING.md).

## License

Released under the [MIT License](LICENSE).
