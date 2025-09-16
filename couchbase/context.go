package couchbase

import (
	"github.com/Trendyol/go-dcp/tracing"
	"github.com/redis/go-redis/v9"
)

type Context struct {
	Tracer    tracing.ListenerTrace
	RedisConn *redis.Conn
	Event     Event
}
