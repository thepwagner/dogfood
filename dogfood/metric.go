package dogfood

import "gopkg.in/alexcesaro/statsd.v2"

type Metric interface {
	Name() string
	HasTags

	apply(client *statsd.Client)
}

type MetricOpt func(s Metric)

func WithTags(tags HasTags) MetricOpt {
	return func(m Metric) {
		if m, ok := m.(setTags); ok {
			m.setTags(tags)
		}
	}
}

type metric struct {
	name string
	tags HasTags
}

func (m *metric) Name() string         { return m.name }
func (m *metric) apply(*statsd.Client) {}
func (m *metric) Tags() Tags {
	if m.tags != nil {
		return m.tags.Tags()
	}
	return nil
}

type setTags interface{ setTags(t HasTags) }

func (m *metric) setTags(t HasTags) {
	m.tags = t
}
