// Package chartertest is the gate and the test marker of a charter, built
// on the charter package: Verifies records which specs a test verifies,
// Generate writes the Go of spec files, Check enforces the registry's rules against the tree, and
// Read lists the entries with the tests that cite them. Observes records
// which observations a probe checks.
package chartertest

import (
	"testing"

	"github.com/leroux/charter"
)

// Verifies records that the calling test verifies a spec, and which of its
// scenarios, by handle. It is the first statement of a test function; the
// gate reads the call statically and refuses one anywhere else, or one
// naming a handle the spec does not have. A test that verifies several
// specs calls it once per spec.
//
// Each citation becomes a test attribute, "SW-118" or "SW-118/handle",
// which `go test -json` reports as an attr event on the test, so one run
// yields pass, fail or skip per scenario; the report reads them.
//
//charter:realizes CH-024
func Verifies(tb testing.TB, spec *charter.Spec, scenarios ...string) {
	tb.Helper()
	if spec == nil {
		tb.Fatal("chartertest: Verifies needs a spec")
		return
	}
	if len(scenarios) == 0 {
		tb.Attr("claim", spec.ID)
		return
	}
	for _, s := range scenarios {
		tb.Attr("claim", spec.ID+"/"+s)
	}
}
