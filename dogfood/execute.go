package dogfood

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

type Executor struct {
	Client *statsd.Client
}

func (e *Executor) Execute(s Scenario) error {
	log := logrus.WithField("scenario", s.Name())
	log.Debug("executing scenario...")

	tags := s.Tags()
	log.WithField("tags", tags).Debug("computed scenario tags")
	client := taggedClient(e.Client, tags)

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

func taggedClient(client *statsd.Client, tags Tags) *statsd.Client {
	flat := make([]string, 0, len(tags)*2)
	for k, v := range tags {
		flat = append(flat, k, v)
	}
	return client.Clone(statsd.Tags(flat...))
}
