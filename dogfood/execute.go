package dogfood

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

type Executor struct {
	statsd *statsd.Client
}

func NewExecutor(c *statsd.Client) *Executor {
	return &Executor{statsd: c}
}

func (e *Executor) Run(ctx context.Context, s Scenario) error {
	var wg sync.WaitGroup
	for i := 0; i < s.Concurrency(); i++ {
		wg.Add(1)
		go e.executeScenario(ctx, &wg, s, i)
	}
	return wait(ctx, &wg)
}

func (e *Executor) executeScenario(ctx context.Context, wg *sync.WaitGroup, s Scenario, id int) {
	defer wg.Done()

	phases := s.Phases()
	log := logrus.WithField("id", id)
	log.WithFields(logrus.Fields{
		"scenario": s.Name(),
		"phases":   len(phases),
	}).Info("starting scenario...")
	defer func() {
		log.Debug("scenario complete")
	}()

	scenarioTags := s.Tags()
	log.WithField("tags", scenarioTags).Debug("computed scenario tags")
	scenarioClient := taggedClient(e.statsd, scenarioTags)

phaseLoop:
	for _, phase := range phases {
		duration := phase.Duration()
		log.WithFields(logrus.Fields{
			"phase": phase.Name(),
			"dur":   duration,
		}).Info("starting phase...")

		phaseEnd := time.Now().Add(duration)

		delayFunc := phase.Delay()
		delay := time.Duration(0)
		for {
			if time.Now().After(phaseEnd) {
				continue phaseLoop
			}

			phaseTags := phase.Tags()
			log.WithField("tags", phaseTags).Debug("computed phase tags")
			phaseClient := taggedClient(scenarioClient, phaseTags)

			// TODO: gate on phase duration
			if err := executePhaseLoop(ctx, log, phase, phaseClient, delay); err != nil {
				if ctx.Err() == nil {
					log.WithError(err).Error("scenario terminated")
				}
				return
			}
			delay = delayFunc()
			log.WithField("delay", delay.Truncate(time.Millisecond)).Debug("processed phase loop")
		}
	}
}

func executePhaseLoop(ctx context.Context, log logrus.FieldLogger, sp ScenarioPhase, client *statsd.Client, delay time.Duration) error {
	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		for i, m := range sp.Metrics() {
			metricLog := log.WithFields(logrus.Fields{
				"metric_index": i,
				"metric":       m.Name(),
			})
			metricLog.Debug("processing metric...")

			var metricClient *statsd.Client
			if metricTags := m.Tags(); metricTags != nil {
				metricLog.WithField("tags", metricTags).Debug("computed metric tags")
				metricClient = taggedClient(client, metricTags)
			} else {
				metricClient = client
			}

			m.apply(metricClient)
		}
		return nil
	}
}

func taggedClient(client *statsd.Client, tags Tags) *statsd.Client {
	flat := make([]string, 0, len(tags)*2)
	for k, v := range tags {
		flat = append(flat, k, v)
	}
	return client.Clone(statsd.Tags(flat...))
}

func wait(ctx context.Context, wg *sync.WaitGroup) error {
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return nil
	}
}
