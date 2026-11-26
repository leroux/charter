package chartertest

import (
	"testing"

	"github.com/leroux/charter/claims"
)

// Lint prints what the gate fails on before its advice.
func TestLintRunsTheWholeGateFirst(t *testing.T) {
	Verifies(t, claims.CH_041, "gate-problems-first")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	t.Setenv(UpdateSpecsEnv, "1")
	t.Setenv(UpdateTypeScriptEnv, "1")
	if _, err := Problems(fixtureOptions(root)); err != nil {
		t.Fatal(err)
	}
	lines, err := Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(lines, "rule 25: scenario waited-request-passes of SW-010 is verified by no test") {
		t.Errorf("lint lacks the gate's problem: %q", lines)
	}
}
