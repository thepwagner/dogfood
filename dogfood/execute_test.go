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

	s := dogfood.NewScenario(
		// Simulate an HTTP request every second:
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
			dogfood.NewCountMetric("http.request"),
			dogfood.NewCountMetric("http.response", dogfood.WithTags(
				dogfood.NewWeightedTag("status", map[string]int{
					"200": 9,
					"500": 1,
				})),
			),
		),
	)

	c, err := statsd.New(statsd.TagsFormat(statsd.Datadog))
	require.NoError(t, err)
	defer c.Close()

	e := dogfood.Executor{Client: c}
	for {
		err := e.Execute(s)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
	}
}
