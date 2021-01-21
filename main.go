package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/dogfood/dogfood"
	"github.com/thepwagner/dogfood/scenarios"
	"gopkg.in/alexcesaro/statsd.v2"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

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
	exec := dogfood.NewExecutor(c)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exec.Start(ctx,
		scenarios.HttpTraffic,
		dogfood.WithDelayFunc(dogfood.RandomDelay(200*time.Millisecond, 300*time.Millisecond)),
		dogfood.WithConcurrency(5),
	)
	_ = exec.Wait(ctx)
	return nil
}
