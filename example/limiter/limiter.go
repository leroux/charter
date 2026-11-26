// Package limiter is a token-bucket rate limiter keyed by caller.
package limiter

import (
	"errors"
	"sync"
	"time"
)

// Limiter lets each key through at a sustained rate with a burst on top.
// It is safe for concurrent use.
type Limiter struct {
	rate  float64 // tokens per second
	burst float64 // the most tokens a bucket holds
	clock Clock

	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	tokens float64
	last   time.Time
}

// New returns a limiter that earns rate tokens per second per key, up to
// burst. A rate that is not above zero, or a burst below one, is refused.
func New(rate float64, burst int, clock Clock) (*Limiter, error) {
	if rate <= 0 {
		return nil, errors.New("limiter: the rate must be above zero")
	}
	if burst < 1 {
		return nil, errors.New("limiter: the burst must be at least one")
	}
	if clock == nil {
		clock = SystemClock{}
	}
	return &Limiter{rate: rate, burst: float64(burst), clock: clock, buckets: map[string]*bucket{}}, nil
}

// Allow reports whether key may make one request now, and spends a token
// when it may.
//
//charter:realizes SW-261126-VX8Z, SW-261126-39F6
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock.Now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	if elapsed := now.Sub(b.last).Seconds(); elapsed > 0 {
		b.tokens = min(l.burst, b.tokens+elapsed*l.rate)
		b.last = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
