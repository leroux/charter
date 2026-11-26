package limiter

import (
	"sync"
	"testing"
	"time"

	"example.com/ratelimit/claims"
)

// fakeClock is a clock a test moves by hand.
type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newLimiter(t *testing.T, rate float64, burst int) (*Limiter, *fakeClock) {
	t.Helper()
	c := &fakeClock{t: time.Unix(1_700_000_000, 0)}
	l, err := New(rate, burst, c)
	if err != nil {
		t.Fatal(err)
	}
	return l, c
}

// allowed counts how many of n requests at this instant go through.
func allowed(l *Limiter, key string, n int) int {
	got := 0
	for range n {
		if l.Allow(key) {
			got++
		}
	}
	return got
}

func TestANewKeyBurstsThenIsDenied(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_VX8Z, "starts-full")
	l, _ := newLimiter(t, 1, 3)
	if got := allowed(l, "a", 4); got != 3 {
		t.Errorf("%d of 4 requests at once went through, want the burst of 3", got)
	}
}

func TestAnEmptyBucketEarnsTheRate(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_VX8Z, "earns-the-rate")
	l, c := newLimiter(t, 2, 3)
	allowed(l, "a", 3)
	c.advance(time.Second)
	if got := allowed(l, "a", 3); got != 2 {
		t.Errorf("after a second at rate 2, %d of 3 went through, want 2", got)
	}
}

func TestAQuietKeyIsCappedAtTheBurst(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_VX8Z, "capped-at-burst")
	l, c := newLimiter(t, 10, 4)
	allowed(l, "a", 4)
	c.advance(time.Hour)
	if got := allowed(l, "a", 100); got != 4 {
		t.Errorf("after an hour quiet, %d went through, want the burst of 4", got)
	}
}

func TestADenialSpendsNothing(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_VX8Z, "denial-spends-nothing")
	l, c := newLimiter(t, 1, 1)
	allowed(l, "a", 1)
	c.advance(500 * time.Millisecond)
	if l.Allow("a") {
		t.Fatal("half a token let a request through")
	}
	c.advance(500 * time.Millisecond)
	if !l.Allow("a") {
		t.Error("the denial delayed the next token: a full second earned nothing")
	}
}

func TestAnEmptyKeyLeavesANewKeyItsBurst(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_39F6, "keys-independent")
	l, _ := newLimiter(t, 1, 2)
	allowed(l, "noisy", 10)
	if got := allowed(l, "quiet", 3); got != 2 {
		t.Errorf("the quiet key got %d, want its full burst of 2", got)
	}
}

func TestConcurrentKeysEachGetTheirBurst(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_39F6, "concurrent-keys")
	l, _ := newLimiter(t, 1, 5)
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	counts := make([]int, len(keys))
	var wg sync.WaitGroup
	for i, k := range keys {
		for range 4 {
			wg.Go(func() {
				n := allowed(l, k, 5)
				l.mu.Lock()
				counts[i] += n
				l.mu.Unlock()
			})
		}
	}
	wg.Wait()
	for i, k := range keys {
		if counts[i] != 5 {
			t.Errorf("key %s got %d through, want exactly the burst of 5", k, counts[i])
		}
	}
}

func TestAStillClockEarnsNothing(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_2XG8, "time-from-clock")
	l, _ := newLimiter(t, 1000, 1)
	allowed(l, "a", 1)
	time.Sleep(20 * time.Millisecond) // twenty tokens of wall time
	if l.Allow("a") {
		t.Error("the limiter earned from the wall clock; it must read time from its Clock alone")
	}
}
