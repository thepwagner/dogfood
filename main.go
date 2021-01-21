package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

func main() {
	if err := run(); err != nil {
		logrus.WithError(err).Fatal("failed")
	}
}

func run() error {
	c, err := statsd.New(statsd.TagsFormat(statsd.Datadog))
	if err != nil {
		return err
	}
	defer c.Close()

	func() {
		tagged := c.Clone(statsd.Tags("foo", "1", "bar", "1"))
		tagged.Increment("pwagner.test")
	}()

	func() {
		tagged := c.Clone(statsd.Tags("foo", "2", "bar", "1"))
		tagged.Increment("pwagner.test")
	}()

	return nil
}
