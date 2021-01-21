package dogfood_test

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestExecute(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Simulate an HTTP request every second:

	//c, err := statsd.New(statsd.TagsFormat(statsd.Datadog))
	//require.NoError(t, err)
	//defer c.Close()
	//
	//e := dogfood.Executor{statsd: c}
	//for {
	//	err := e.Execute(s)
	//	require.NoError(t, err)
	//	time.Sleep(200 * time.Millisecond)
	//}
}
