package dogfood

import "gopkg.in/alexcesaro/statsd.v2"

func NewCountMetric(name string, opts ...MetricOpt) Metric {
	m := &countMetric{
		metric: metric{name: name},
		count:  1,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func WithCount(count int) MetricOpt {
	return func(m Metric) {
		if m, ok := m.(*countMetric); ok {
			m.count = count
		}
	}
}

type countMetric struct {
	metric
	count int
}

func (m *countMetric) apply(c *statsd.Client) {
	c.Count(m.name, m.count)
}
