package chartertest

import (
	"testing"

	"github.com/leroux/charter"
)

// Observes records that the calling probe checks the given observations
// against their live subject. Each becomes a test attribute, which
// `go test -json` reports as an attr event on the test. Call it after the
// probe's preconditions have been checked, so a skipped probe emits nothing.
func Observes(tb testing.TB, observations ...*charter.Observation) {
	tb.Helper()
	for _, o := range observations {
		tb.Attr("observation", o.ID)
	}
}
