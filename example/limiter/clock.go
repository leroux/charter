package limiter

import "time"

// Clock tells the limiter the time. The limiter reads time through it
// alone, so a test drives time by hand.
//
//charter:realizes SW-261126-2XG8
type Clock interface {
	Now() time.Time
}

// SystemClock is the wall clock.
type SystemClock struct{}

// Now returns the current time.
func (SystemClock) Now() time.Time { return time.Now() }
