//go:build fxprobe

package a

import (
	"testing"

	"example.com/fixture/claims"
	"github.com/leroux/charter/chartertest"
)

func TestProbeHealth(t *testing.T) {
	chartertest.Observes(t, claims.HealthIsJSON)
}
