package claims_test

import (
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/claims/gate"
)

// TestGate is the charter's own gate: every rule of chartertest.Check over
// claims and the tooling that realizes it, and nothing else.
func TestGate(t *testing.T) {
	chartertest.Verifies(t, claims.CH_006, "self-gate-walks-the-tooling-only")
	chartertest.Check(t, gate.Options)
}
