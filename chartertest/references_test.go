package chartertest

import (
	"reflect"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestReferencesListCitingTests(t *testing.T) {
	root := copyFixture(t)
	refs := refsOf(t, fixtureOptions(root))
	want := map[string][]Ref{
		"SW_001":       {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 11}, {File: "web/a.test.ts", Func: "survives a restart", Line: 7}},
		"STK_001":      {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 13}},
		"CON_002":      {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 12}, {File: "web/a.test.ts", Func: "folds %d", Line: 10}},
		"CON_001":      {{File: "pkg/a/a_test.go", Func: "TestLock", Line: 17}, {File: "web/a.test.ts", Func: `locks ${"maybe"}`, Line: 13}},
		"SW_002":       {{File: "pkg/a/a_test.go", Func: "TestLock", Line: 18}, {File: "web/b.spec.ts", Func: "windows is not supported yet", Line: 6}},
		"HealthIsJSON": {{File: "pkg/a/a_probe_test.go", Func: "TestProbeHealth", Line: 13}},
	}
	if !reflect.DeepEqual(refs, want) {
		t.Errorf("refs = %v, want %v", refs, want)
	}
}

func TestReferencesFollowAnAlias(t *testing.T) {
	Verifies(t, claims.CH_023, "alias-followed")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", `"example.com/fixture/claims"`, `c "example.com/fixture/claims"`)
	replaceIn(t, root, "pkg/a/a_test.go", "chartertest.Verifies(t, claims.SW_001)", "chartertest.Verifies(t, c.SW_001)")
	refs := refsOf(t, fixtureOptions(root))
	if got := refs["SW_001"]; len(got) != 2 || got[0].Func != "TestRestart" {
		t.Errorf("aliased reference not found: %v", refs)
	}
}

// refsOf reads the charter and returns, by entry name, the tests that cite
// each entry, and for an observation the probes that reference it.
func refsOf(t *testing.T, opts Options) map[string][]Ref {
	t.Helper()
	reg, err := Read(opts)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]Ref{}
	for _, e := range reg.Entries {
		if len(e.Refs) > 0 {
			out[e.Name] = e.Refs
		}
	}
	return out
}
