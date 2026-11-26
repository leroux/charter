package chartertest

import (
	"reflect"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestCollectRefsFindsCitationsAndReferences(t *testing.T) {
	Verifies(t, claims.CH_023, "go-citation-read", "reference-is-not-citation")
	root := copyFixture(t)
	r, tr, pkg, _ := parseFixture(t, root)
	refs, cites, problems := collectRefs(tr, r, pkg)
	if len(problems) != 0 {
		t.Fatalf("clean fixture reported %q", problems)
	}
	wantCites := map[string][]Ref{
		"SW_001":  {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 11}},
		"CON_002": {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 12}},
		"STK_001": {{File: "pkg/a/a_test.go", Func: "TestRestart", Line: 13}},
		"CON_001": {{File: "pkg/a/a_test.go", Func: "TestLock", Line: 17}},
		"SW_002":  {{File: "pkg/a/a_test.go", Func: "TestLock", Line: 18}},
	}
	if !reflect.DeepEqual(cites, wantCites) {
		t.Errorf("cites = %v, want %v", cites, wantCites)
	}
	if got := refs["HealthIsJSON"]; len(got) != 1 || got[0].File != "pkg/a/a_probe_test.go" || got[0].Func != "TestProbeHealth" {
		t.Errorf("the probe's reference = %v", got)
	}
	if got := refs["SW_001"]; len(got) != 1 || got[0].Line != 11 {
		t.Errorf("a citation is a reference too: %v", got)
	}
}

// A citation opening a subtest names it as go test does below the test,
// a title that is no literal standing for any name.
func TestCollectRefsNamesTheSubtestACitationOpens(t *testing.T) {
	Verifies(t, claims.CH_023, "go-citation-read")
	root := copyFixture(t)
	writeFile(t, root, "pkg/a/sub_test.go", "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n"+
		"func TestSub(t *testing.T) {\n\tt.Run(\"clamp holds\", func(t *testing.T) {\n\t\tchartertest.Verifies(t, claims.SW_001)\n"+
		"\t\tt.Run(\"a/b\", func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.SW_002)\n\t\t})\n\t})\n"+
		"\tfor _, name := range []string{\"x\"} {\n\t\tt.Run(name, func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.CON_001)\n\t\t})\n\t}\n}\n")
	r, tr, pkg, _ := parseFixture(t, root)
	_, cites, problems := collectRefs(tr, r, pkg)
	if len(problems) != 0 {
		t.Fatalf("problems %q", problems)
	}
	sub := map[string][]string{}
	for name, refs := range cites {
		for _, ref := range refs {
			if ref.File == "pkg/a/sub_test.go" {
				sub[name] = ref.Subtest
			} else if ref.Subtest != nil {
				t.Errorf("%s at %s:%d names subtest %q, want none", name, ref.File, ref.Line, ref.Subtest)
			}
		}
	}
	want := map[string][]string{"SW_001": {"clamp_holds"}, "SW_002": {"clamp_holds", "a", "b"}, "CON_001": {""}}
	if !reflect.DeepEqual(sub, want) {
		t.Errorf("subtests %q, want %q", sub, want)
	}
}

func TestCollectRefsRejectsDotImportsAndAssignments(t *testing.T) {
	Verifies(t, claims.CH_025, "assignment-refused", "dot-import-refused")
	root := copyFixture(t)
	writeFile(t, root, "pkg/b/b_test.go", "package b\n\nimport . \"example.com/fixture/claims\"\n\nvar _ = SW_001\n")
	writeFile(t, root, "pkg/b/c_test.go", "package b\n\nimport \"example.com/fixture/claims\"\n\nfunc init() {\n\tclaims.HealthIsJSON.Text = \"x\"\n\tclaims.CON_002 = nil\n\tclaims.CON_002.Accounts[0] = nil\n}\n")
	r, tr, pkg, _ := parseFixture(t, root)
	_, _, problems := collectRefs(tr, r, pkg)
	for _, want := range []string{
		"pkg/b/b_test.go:1: rule 10: the entries package must not be dot-imported",
		"pkg/b/c_test.go:6: rule 1: HealthIsJSON is assigned outside its declaration",
		"pkg/b/c_test.go:7: rule 1: CON_002 is assigned outside its declaration",
		"pkg/b/c_test.go:8: rule 1: CON_002 is assigned outside its declaration",
	} {
		if !anyContains(problems, want) {
			t.Errorf("want %q in %q", want, problems)
		}
	}
}
