package scenarios

import (
	"fmt"
	"time"

	"github.com/thepwagner/dogfood/dogfood"
	"gopkg.in/yaml.v3"
)

type scenarioFile struct {
	Name        string                 `yaml:"name"`
	Concurrency int                    `yaml:"concurrency"`
	Tags        map[string]interface{} `yaml:"tags"`
	Phases      []phaseFile            `yaml:"phases"`
}

type phaseFile struct {
	Name     string                 `yaml:"name"`
	Duration string                 `yaml:"duration"`
	Tags     map[string]interface{} `yaml:"tags"`
	Metrics  []metricFile           `yaml:"metrics"`
}

type metricFile struct {
	Name string                 `yaml:"name"`
	Tags map[string]interface{} `yaml:"tags"`

	Count  int `yaml:"count"`
	Timing *struct {
		Min string `yaml:"min"`
		Max string `yaml:"max"`
	} `yaml:"timing"`
}

func LoadScenario(b []byte) (dogfood.Scenario, error) {
	var sf scenarioFile
	if err := yaml.Unmarshal(b, &sf); err != nil {
		return nil, fmt.Errorf("parsing sf: %w", err)
	}

	var name string
	if sf.Name != "" {
		name = sf.Name
	} else {
		name = "(no name)"
	}

	var opts []dogfood.ScenarioOpt
	if sf.Concurrency != 0 {
		opts = append(opts, dogfood.WithConcurrency(sf.Concurrency))
	}
	if sf.Tags != nil {
		st, err := loadTags(sf.Tags)
		if err != nil {
			return nil, fmt.Errorf("loading scenario tags: %w", err)
		}
		opts = append(opts, dogfood.WithScenarioTags(st))
	}

	phases := make([]dogfood.ScenarioPhase, 0, len(sf.Phases))
	for phaseIndex, phase := range sf.Phases {
		var phaseName string
		if phase.Name != "" {
			phaseName = phase.Name
		} else {
			phaseName = fmt.Sprintf("phase%d", phaseIndex)
		}

		var phaseOpts []dogfood.ScenarioPhaseOpt
		if phase.Duration != "" {
			dur, err := time.ParseDuration(phase.Duration)
			if err != nil {
				return nil, fmt.Errorf("parsing phase duration %q: %w", phaseName, err)
			}
			phaseOpts = append(phaseOpts, dogfood.WithPhaseDuration(dur))
		}
		if phase.Tags != nil {
			pt, err := loadTags(phase.Tags)
			if err != nil {
				return nil, fmt.Errorf("loading phase tags %q %q: %w", name, phaseName, err)
			}
			phaseOpts = append(phaseOpts, dogfood.WithPhaseTags(pt))
		}

		phaseMetrics := make([]dogfood.Metric, 0, len(phase.Metrics))
		for metricIndex, m := range phase.Metrics {
			if m.Name == "" {
				return nil, fmt.Errorf("metric missing name phase tags %q %q %d", name, phaseName, metricIndex)
			}
			mt, err := loadTags(m.Tags)
			if err != nil {
				return nil, fmt.Errorf("loading metric tags %q %q %q: %w", name, phaseName, m.Name, err)
			}

			if m.Count > 0 {
				phaseMetrics = append(phaseMetrics, dogfood.NewCountMetric(m.Name, dogfood.WithCount(m.Count), dogfood.WithTags(mt)))
			} else if m.Timing != nil {
				timeMin, err := time.ParseDuration(m.Timing.Min)
				if err != nil {
					return nil, fmt.Errorf("parsing min time: %w", err)
				}
				timeMax, err := time.ParseDuration(m.Timing.Max)
				if err != nil {
					return nil, fmt.Errorf("parsing max time: %w", err)
				}

				phaseMetrics = append(phaseMetrics, dogfood.NewTimingMetric(m.Name, dogfood.RandomTiming(timeMin, timeMax), dogfood.WithTags(mt)))
			} else {
				return nil, fmt.Errorf("unknown metric type %q %q %q", name, phaseName, m.Name)
			}
		}
		phaseOpts = append(phaseOpts, dogfood.WithMetrics(phaseMetrics...))

		phases = append(phases, dogfood.NewScenarioPhase(phaseName, phaseOpts...))
	}
	opts = append(opts, dogfood.WithPhases(phases...))

	return dogfood.NewScenario(name, opts...), nil
}

func loadTags(t map[string]interface{}) (dogfood.HasTags, error) {
	tags := make([]dogfood.HasTags, 0, len(t))
	for name, config := range t {
		switch v := config.(type) {
		case string:
			tags = append(tags, dogfood.Tags{name: v})
		case map[string]interface{}:
			if weights, ok := v["weights"].(map[string]interface{}); ok {
				w := make(map[string]int, len(weights))
				for k, v := range weights {
					if weight, ok := v.(int); ok {
						w[k] = weight
					} else {
						return nil, fmt.Errorf("invalid weighted tag %q value %q", name, k)
					}
				}
				tags = append(tags, dogfood.NewWeightedTag(name, w))
			} else {
				return nil, fmt.Errorf("invalid tag %q", name)
			}
		default:
			return nil, fmt.Errorf("invalid tag configuration %q", name)
		}

	}
	return dogfood.NewMergedTags(tags...), nil
}
