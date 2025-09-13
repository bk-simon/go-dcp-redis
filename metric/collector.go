package metric

import (
	"github.com/Trendyol/go-dcp-redis/redis/bulk"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	bulk *bulk.Bulk

	processLatency            *prometheus.Desc
	bulkRequestProcessLatency *prometheus.Desc
}

func (s *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(s, ch)
}

func (s *Collector) Collect(ch chan<- prometheus.Metric) {
	bulkMetric := s.bulk.GetMetric()

	ch <- prometheus.MustNewConstMetric(
		s.processLatency,
		prometheus.GaugeValue,
		float64(bulkMetric.ProcessLatencyMs),
		[]string{}...,
	)

	ch <- prometheus.MustNewConstMetric(
		s.bulkRequestProcessLatency,
		prometheus.GaugeValue,
		float64(bulkMetric.BulkRequestProcessLatencyMs),
		[]string{}...,
	)
}
