package claims_test

import (
	"testing"

	"example.com/toml/claims/gate"
	"github.com/leroux/charter/chartertest"
)

// TestGate holds the tree to the charter: every active decision realized
// and every scenario verified.
func TestGate(t *testing.T) {
	chartertest.Check(t, gate.Options)
}
