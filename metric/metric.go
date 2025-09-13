package metric

import (
	"github.com/Trendyol/go-dcp-redis/redis/bulk"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	RedisConnectorLatency            = "redis_connector_latency_ms"
	RedisConnectorBulkRequestLatency = "redis_connector_bulk_request_process_latency_ms"
)

type MetricCollector struct {
	bulk                      *bulk.Bulk
	redisConnectorLatency     prometheus.Gauge
	redisConnectorBulkLatency prometheus.Gauge
}

func NewMetricCollector(bulk *bulk.Bulk) *MetricCollector {
	return &MetricCollector{
		bulk: bulk,
		redisConnectorLatency: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: RedisConnectorLatency,
			Help: "Time to adding to the batch.",
		}),
		redisConnectorBulkLatency: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: RedisConnectorBulkRequestLatency,
			Help: "Time to process bulk request.",
		}),
	}
}

func (s *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	metric := s.bulk.GetMetric()

	s.redisConnectorLatency.Set(float64(metric.ProcessLatencyMs))
	s.redisConnectorBulkLatency.Set(float64(metric.BulkRequestProcessLatencyMs))

	s.redisConnectorLatency.Collect(ch)
	s.redisConnectorBulkLatency.Collect(ch)
}

func (s *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	s.redisConnectorLatency.Describe(ch)
	s.redisConnectorBulkLatency.Describe(ch)
}
