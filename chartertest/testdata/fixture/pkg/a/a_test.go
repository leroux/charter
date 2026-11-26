package a

import (
	"testing"

	"example.com/fixture/claims"
	"github.com/leroux/charter/chartertest"
)

func TestRestart(t *testing.T) {
	chartertest.Verifies(t, claims.SW_001)
	chartertest.Verifies(t, claims.CON_002)
	chartertest.Verifies(t, claims.STK_001)
}

func TestLock(t *testing.T) {
	chartertest.Verifies(t, claims.CON_001)
	chartertest.Verifies(t, claims.SW_002)
}
