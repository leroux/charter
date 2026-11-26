package claims_test

import (
	"testing"

	"example.com/ratelimit/claims/gate"
	"github.com/leroux/charter/chartertest"
)

// TestGate holds the tree to the charter: every active decision realized
// and every scenario verified, every surface accounted for.
func TestGate(t *testing.T) {
	chartertest.Check(t, gate.Options)
}
