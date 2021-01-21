package dogfood

type Scenario interface {
	Name() string
	Concurrency() int
	Phases() []ScenarioPhase
	HasTags
}

type ScenarioOpt func(s *scenario)

func NewScenario(name string, opts ...ScenarioOpt) Scenario {
	s := &scenario{
		name:        name,
		concurrency: 1,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithConcurrency(concurrency int) ScenarioOpt {
	return func(s *scenario) {
		s.concurrency = concurrency
	}
}

func WithPhases(phases ...ScenarioPhase) ScenarioOpt {
	return func(s *scenario) {
		s.phases = append(s.phases, phases...)
	}
}

func WithScenarioTags(tags HasTags) ScenarioOpt {
	return func(s *scenario) {
		s.tags = tags
	}
}

type scenario struct {
	name        string
	concurrency int
	phases      []ScenarioPhase
	tags        HasTags
}

func (s *scenario) Name() string            { return s.name }
func (s *scenario) Concurrency() int        { return s.concurrency }
func (s *scenario) Phases() []ScenarioPhase { return s.phases }
func (s *scenario) Tags() Tags {
	if s.tags != nil {
		return s.tags.Tags()
	}
	return nil
}
