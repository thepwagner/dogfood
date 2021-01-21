package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/dogfood/dogfood"
	"github.com/thepwagner/dogfood/scenarios"
	"gopkg.in/alexcesaro/statsd.v2"
)

func main() {
	//logrus.SetLevel(logrus.DebugLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func() {
		for range sigC {
			cancel()
		}
	}()

	if err := run(ctx, os.Args...); err != nil && ctx.Err() == nil {
		logrus.WithError(err).Fatal("failed")
	}
}

func run(ctx context.Context, args ...string) error {
	if len(args) != 2 {
		return errors.New("specify scenario path")
	}

	scenario, err := readScenario(args[1])
	if err != nil {
		return err
	}

	c, err := statsd.New(statsd.TagsFormat(statsd.Datadog))
	if err != nil {
		return err
	}
	defer c.Close()
	return dogfood.NewExecutor(c).Run(ctx, scenario)
}

func readScenario(path string) (dogfood.Scenario, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading script: %w", err)
	}
	return scenarios.LoadScenario(b)
}
