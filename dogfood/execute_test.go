package dogfood_test

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/dogfood/dogfood"
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
		// Every metric should have some common tags:
		dogfood.WithTags(
			dogfood.NewMergedTags(
				dogfood.Tags{"host": "test.test.com"},
				dogfood.NewWeightedTag("method", map[string]int{
					"GET":    100,
					"POST":   10,
					"DELETE": 1,
				}),
			),
		),
	)
	//s := &dogfood.Scenario{
	//	Name:      "http request",
	//	Frequency: 1 * time.Second,
	//	//Tags: dogfood.TagDistribution{
	//	//	"method": {
	//	//		"GET":    100,
	//	//		"POST":   10,
	//	//		"DELETE": 1,
	//	//	},
	//	//},
	//	Sequences: []dogfood.MetricSequence{
	//		{
	//			Metric: "http.request",
	//		},
	//	},
	//}
	for i := 0; i < 10; i++ {
		logrus.WithField("loop", i).Info("execute")
		err := dogfood.Execute(s)
		require.NoError(t, err)
	}

}
