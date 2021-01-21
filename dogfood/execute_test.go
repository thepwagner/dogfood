package dogfood_test

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/dogfood/dogfood"
	"gopkg.in/alexcesaro/statsd.v2"
)

func TestExecute(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Simulate an HTTP request every second:
	s := dogfood.NewScenario(
		"http test",
		1.*time.Second,
		// Common tags for every metric:
		dogfood.WithScenarioTags(
			dogfood.NewMergedTags(
				dogfood.Tags{"host": "test.test.com"},
				// Keep the `method` tag consistent
				dogfood.NewWeightedTag("method", map[string]int{
					"GET":    100,
					"POST":   10,
					"DELETE": 1,
				}),
			),
		),
		dogfood.WithMetrics(
			// Count the request:
			dogfood.NewCountMetric("http.request"),
			// Count the response:
			dogfood.NewCountMetric("http.response", dogfood.WithTags(
				dogfood.NewWeightedTag("status", map[string]int{
					"200": 9,
					"500": 1,
				})),
			),
			dogfood.NewTimingMetric("http.duration", dogfood.NewRandomTiming(200, 500)),
		),
	)

	c, err := statsd.New(statsd.TagsFormat(statsd.Datadog))
	require.NoError(t, err)
	defer c.Close()

	e := dogfood.Executor{Client: c}
	for {
		err := e.Execute(s)
		require.NoError(t, err)
		time.Sleep(200 * time.Millisecond)
	}
}
