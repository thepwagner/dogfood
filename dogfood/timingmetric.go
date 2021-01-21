package dogfood

import (
	"math/rand"

	"gopkg.in/alexcesaro/statsd.v2"
)

func NewTimingMetric(name string, value func() int, opts ...MetricOpt) Metric {
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
	value func() int
}

func (m *timingMetric) apply(c *statsd.Client) {
	value := m.value()
	c.Timing(m.name, value)
}

func NewRandomTiming(min, max int) func() int {
	return func() int {
		return rand.Intn(max-min) + min
	}
}
