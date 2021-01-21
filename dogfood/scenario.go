package dogfood

import "time"

type Scenario interface {
	Name() string
	HasTags
}

type ScenarioOpt func(s *scenario)

func NewScenario(name string, frequency time.Duration, opts ...ScenarioOpt) Scenario {
	s := &scenario{
		name:      name,
		frequency: frequency,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithTags(tags HasTags) ScenarioOpt {
	return func(s *scenario) {
		s.tags = tags
	}
}

type scenario struct {
	name      string
	frequency time.Duration
	tags      HasTags
}

func (s *scenario) Name() string {
	return s.name
}

func (s *scenario) Tags() Tags {
	return s.tags.Tags()
}
