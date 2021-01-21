package dogfood

import "github.com/sirupsen/logrus"

func Execute(s Scenario) error {
	log := logrus.WithField("scenario", s.Name())
	log.Debug("executing scenario...")

	tags := s.Tags()
	logrus.WithField("tags", tags).Debug("collected common tags")

	return nil
}
