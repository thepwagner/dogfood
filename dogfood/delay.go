package dogfood

import (
	"math/rand"
	"time"
)

type TimeFunc func() time.Duration

func FixedTiming(dur time.Duration) TimeFunc {
	return func() time.Duration { return dur }
}

func RandomTiming(min, max time.Duration) TimeFunc {
	r := int64(max - min)
	return func() time.Duration {
		return time.Duration(rand.Int63n(r)) + min
	}
}
