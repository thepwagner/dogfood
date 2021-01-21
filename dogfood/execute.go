package dogfood

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

type Executor struct {
	statsd *statsd.Client
	wg     sync.WaitGroup
}

func NewExecutor(c *statsd.Client) *Executor {
	return &Executor{statsd: c}
}

type DelayFunc func() time.Duration

func FixedDelay(dur time.Duration) DelayFunc {
	return func() time.Duration { return dur }
}

func RandomDelay(min, max time.Duration) DelayFunc {
	r := int64(max - min)
	return func() time.Duration {
		return time.Duration(rand.Int63n(r)) + min
	}
}

type startScenarioOptions struct {
	concurrency int
	delay       DelayFunc
}

type StartScenarioOpt func(o *startScenarioOptions)

func WithDelayFunc(f DelayFunc) StartScenarioOpt {
	return func(o *startScenarioOptions) {
		o.delay = f
	}
}

func WithConcurrency(concurrency int) StartScenarioOpt {
	return func(o *startScenarioOptions) {
		o.concurrency = concurrency
	}
}

func (e *Executor) Start(ctx context.Context, s Scenario, opts ...StartScenarioOpt) {
	o := &startScenarioOptions{
		delay:       FixedDelay(1 * time.Second),
		concurrency: 1,
	}
	for _, opt := range opts {
		opt(o)
	}
	for i := 0; i < o.concurrency; i++ {
		e.wg.Add(1)
		go e.executorLoop(ctx, s, i, o.delay)
	}
}

func (e *Executor) executorLoop(ctx context.Context, s Scenario, id int, dur DelayFunc) {
	defer e.wg.Done()

	log := logrus.WithFields(logrus.Fields{
		"scenario": s.Name(),
		"id":       id,
	})
	log.Debug("starting scenario...")
	defer func() {
		log.Debug("scenario complete")
	}()

	tags := s.Tags()
	log.WithField("tags", tags).Debug("computed scenario tags")
	client := taggedClient(e.statsd, tags)

	delay := time.Duration(0)
	for {
		if err := delayedExecute(ctx, log, s, client, delay); err != nil {
			if ctx.Err() == nil {
				log.WithError(err).Error("scenario terminated")
			}
			return
		}
		delay = dur()
		log.WithField("delay", delay.Truncate(time.Millisecond)).Info("processed scenario")
	}
}

func delayedExecute(ctx context.Context, log logrus.FieldLogger, s Scenario, client *statsd.Client, delay time.Duration) error {
	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		for i, m := range s.Metrics() {
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

func (e *Executor) Wait(ctx context.Context) error {
	c := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(c)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return nil
	}
}
