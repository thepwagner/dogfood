package dogfood

import "gopkg.in/alexcesaro/statsd.v2"

type Metric interface {
	Name() string
	HasTags

	apply(client *statsd.Client)
}

type MetricOpt func(s *countMetric)

func NewCountMetric(name string, opts ...MetricOpt) Metric {
	m := &countMetric{
		name:  name,
		count: 1,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func WithTags(tags HasTags) MetricOpt {
	return func(m *countMetric) {
		m.tags = tags
	}
}

func WithCount(count int) MetricOpt {
	return func(m *countMetric) {
		m.count = count
	}
}

type countMetric struct {
	name  string
	count int
	tags  HasTags
}

func (m *countMetric) Name() string { return m.name }
func (m *countMetric) Tags() Tags {
	if m.tags != nil {
		return m.tags.Tags()
	}
	return nil
}

func (m *countMetric) apply(c *statsd.Client) {
	c.Count(m.name, m.count)
}
