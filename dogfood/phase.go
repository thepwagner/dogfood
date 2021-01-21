package dogfood

import "time"

type ScenarioPhase interface {
	Name() string
	Duration() time.Duration
	Metrics() []Metric
	Delay() TimeFunc
	HasTags
}

func NewScenarioPhase(name string, opts ...ScenarioPhaseOpt) ScenarioPhase {
	p := &phase{
		name:  name,
		delay: RandomTiming(1*time.Second, 2*time.Second),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type ScenarioPhaseOpt func(p *phase)

type phase struct {
	name    string
	dur     time.Duration
	metrics []Metric
	delay   TimeFunc
	tags    HasTags
}

func (p phase) Name() string            { return p.name }
func (p phase) Duration() time.Duration { return p.dur }
func (p phase) Metrics() []Metric       { return p.metrics }
func (p phase) Delay() TimeFunc         { return p.delay }
func (p phase) Tags() Tags {
	if p.tags != nil {
		return p.tags.Tags()
	}
	return nil
}

func WithMetrics(metrics ...Metric) ScenarioPhaseOpt {
	return func(p *phase) {
		p.metrics = append(p.metrics, metrics...)
	}
}

func WithPhaseDuration(dur time.Duration) ScenarioPhaseOpt {
	return func(p *phase) {
		p.dur = dur
	}
}

func WithPhaseTags(tags HasTags) ScenarioPhaseOpt {
	return func(p *phase) {
		p.tags = tags
	}
}

func WithDelayFunc(delay TimeFunc) ScenarioPhaseOpt {
	return func(p *phase) {
		p.delay = delay
	}
}
