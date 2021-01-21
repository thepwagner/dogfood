package dogfood

import (
	"time"

	"gopkg.in/alexcesaro/statsd.v2"
)

func NewTimingMetric(name string, value TimeFunc, opts ...MetricOpt) Metric {
	m := &timingMetric{
		metric: metric{name: name},
		value:  value,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

type timingMetric struct {
	metric
	value func() time.Duration
}

func (m *timingMetric) apply(c *statsd.Client) {
	value := m.value()
	c.Timing(m.name, int64(value/time.Millisecond))
}
