package dogfood

type Scenario interface {
	Name() string
	Metrics() []Metric
	HasTags
}

type ScenarioOpt func(s *scenario)

func NewScenario(name string, opts ...ScenarioOpt) Scenario {
	s := &scenario{
		name: name,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithScenarioTags(tags HasTags) ScenarioOpt {
	return func(s *scenario) {
		s.tags = tags
	}
}

func WithMetric(metric Metric) ScenarioOpt {
	return func(s *scenario) {
		s.metrics = append(s.metrics, metric)
	}
}

func WithMetrics(metrics ...Metric) ScenarioOpt {
	return func(s *scenario) {
		s.metrics = append(s.metrics, metrics...)
	}
}

type scenario struct {
	name    string
	metrics []Metric
	tags    HasTags
}

func (s *scenario) Name() string      { return s.name }
func (s *scenario) Metrics() []Metric { return s.metrics }
func (s *scenario) Tags() Tags        { return s.tags.Tags() }
