package scenarios

import (
	"github.com/thepwagner/dogfood/dogfood"
)

var HttpTraffic = dogfood.NewScenario(
	"http test",
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
