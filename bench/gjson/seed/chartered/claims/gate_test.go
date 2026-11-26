package claims_test

import (
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/tidwall/gjson/claims/gate"
)

// TestGate holds the tree to the charter: every active decision realized
// and every scenario verified, every surface accounted for.
func TestGate(t *testing.T) {
	chartertest.Check(t, gate.Options)
}
