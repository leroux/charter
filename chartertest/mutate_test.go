package chartertest

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/leroux/charter/claims"
)

// The mutate tests break internal/mutateprobe for real. Its anchors are
// described here, in a registry of their own, since the gate ignores the
// probe.

const probeDir = "chartertest/internal/mutateprobe"

// repoRoot is the module root, one level above this package.
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// probeAnchors are the realizes directives in probe.go, by what they
// document.
func probeAnchors(t *testing.T, root string) map[string]Anchor {
	t.Helper()
	src, err := os.ReadFile(filepath.Join(root, probeDir, "probe.go"))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]Anchor{}
	lines := strings.Split(string(src), "\n")
	for i, l := range lines {
		if l != "//charter:realizes SW-001" {
			continue
		}
		decl := strings.Fields(lines[i+1])[1]
		decl, _, _ = strings.Cut(decl, "(")
		out[decl] = Anchor{Kind: "realizes", File: probeDir + "/probe.go", Line: i + 1, Decl: decl}
	}
	return out
}

func probeRef(test string, handles ...string) Ref {
	return Ref{File: probeDir + "/probe_test.go", Func: test, Handles: handles}
}

// probeRegistry: SW-001 is realized by Clamp, Label and Limit, its bounds
// scenario cited by a test that notices a broken Clamp and its label
// scenario by one that notices nothing. SW-002 and SW-003 are realized by
// Clamp and cited by the unsteady and the littering tests.
func probeRegistry(t *testing.T, root string) *Registry {
	t.Helper()
	a := probeAnchors(t, root)
	spec := func(id string, handles ...string) Entry {
		e := Entry{ID: id, Name: strings.ReplaceAll(id, "-", "_"), Text: id + " title", Level: levelSpec, Status: "active"}
		for _, h := range handles {
			e.Scenarios = append(e.Scenarios, Scenario{Handle: h})
		}
		return e
	}
	sw1 := spec("SW-001", "bounds", "label")
	sw1.Anchors = []Anchor{a["Clamp"], a["Label"], a["Limit"]}
	sw1.Refs = []Ref{probeRef("TestClampHoldsTheBounds", "bounds"), probeRef("TestLabelIsNamed", "label")}
	sw2 := spec("SW-002", "steady")
	sw2.Anchors = []Anchor{a["Clamp"]}
	sw2.Refs = []Ref{probeRef("TestUnsteady", "steady")}
	sw3 := spec("SW-003", "tidy")
	sw3.Anchors = []Anchor{a["Clamp"]}
	sw3.Refs = []Ref{probeRef("TestLitter", "tidy")}
	return &Registry{Entries: []Entry{sw1, sw2, sw3}}
}

// One real run over SW-001, with a given mutant that kills and one that
// does not build, shared by the tests that read it.
var (
	probeRunOnce sync.Once
	probeRun     *MutationReport
	probeRunErr  error
)

func runProbe(t *testing.T) EntryMutation {
	t.Helper()
	root := repoRoot(t)
	probeRunOnce.Do(func() {
		probeRun, probeRunErr = MutateRegistry(root, probeRegistry(t, root), MutateOptions{
			IDs: []string{"SW-001"}, Workers: 4, Cache: "shared",
			Given: []GivenMutant{
				{Entry: "SW-001", File: probeDir + "/probe.go", Old: "return v\n", New: "return v + 1\n", Note: "off by one inside"},
				{Entry: "SW-001", File: probeDir + "/probe.go", Old: `return "many"`, New: `return "many" +`, Note: "does not build"},
				{Entry: "SW-001", File: probeDir + "/probe.go", Old: "\treturn v\n}", New: "\tgo func() { panic(\"crash\") }()\n\tselect {}\n}", Note: "crashes the test binary"},
			},
		})
	})
	if probeRunErr != nil {
		t.Fatal(probeRunErr)
	}
	if len(probeRun.Entries) != 1 {
		t.Fatalf("entries %+v, want SW-001 alone", probeRun.Entries)
	}
	return probeRun.Entries[0]
}

func TestMutantsStayInTheAnchoredFunctions(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "mutants-in-anchored-code")
	root := repoRoot(t)
	reg := probeRegistry(t, root)
	a := probeAnchors(t, root)

	plans, err := planMutants(root, reg, MutateOptions{IDs: []string{"SW-001"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].entry.ID != "SW-001" {
		t.Fatalf("plans %+v, want SW-001 alone", plans)
	}
	funcs := map[string]int{}
	for _, m := range plans[0].mutants {
		funcs[m.Func]++
		if m.File != probeDir+"/probe.go" {
			t.Errorf("mutant in %s, outside the anchored file", m.File)
		}
	}
	if funcs["Clamp"] == 0 || funcs["Label"] == 0 || len(funcs) != 2 {
		t.Errorf("mutants per function %v, want Clamp and Label only", funcs)
	}

	// A diff that changed only Label's comparison line: mutants on that
	// line only, and only for the entries Label realizes.
	labelLine := a["Label"].Line + 2
	diff := "diff --git a/" + probeDir + "/probe.go b/" + probeDir + "/probe.go\n" +
		"--- a/" + probeDir + "/probe.go\n+++ b/" + probeDir + "/probe.go\n" +
		"@@ -20 +" + strconv.Itoa(labelLine) + " @@\n-\tif n == 2 {\n+\tif n == 1 {\n"
	changed := parseDiff([]byte(diff))
	plans, err = planMutants(root, reg, MutateOptions{Since: "HEAD"}, changed)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].entry.ID != "SW-001" {
		t.Fatalf("plans for a diff in Label: %d, want SW-001 alone (SW-002 and SW-003 are Clamp's)", len(plans))
	}
	if len(plans[0].mutants) == 0 {
		t.Fatal("no mutant on the changed line")
	}
	for _, m := range plans[0].mutants {
		if m.Line != labelLine || m.Func != "Label" {
			t.Errorf("mutant at %s:%d in %s, want only line %d of Label", m.File, m.Line, m.Func, labelLine)
		}
	}
}

func TestParseDiffReadsTheNewSide(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "mutants-in-anchored-code")
	got := parseDiff([]byte("diff --git a/x.go b/x.go\n--- a/x.go\n+++ b/x.go\n@@ -3,2 +3,4 @@ func f\n+a\n@@ -9 +11,0 @@\n-b\n@@ -20 +22 @@\n" +
		"diff --git a/gone.go b/gone.go\n--- a/gone.go\n+++ /dev/null\n@@ -1 +0,0 @@\n"))
	want := lineSet{{3, 6}, {22, 22}}
	if !slices.Equal(got["x.go"], want) || len(got) != 1 {
		t.Errorf("parseDiff = %v, want x.go %v alone: a pure deletion adds no line, a deleted file none", got, want)
	}
}

func TestAKilledMutantCountsForTheScenariosOfItsTests(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "killed-by-citing-test")
	e := runProbe(t)
	var bounds, survivors int
	for _, m := range e.Mutants {
		switch m.Outcome { //nolint:exhaustive // an unbuilt mutant counts neither way
		case Killed:
			if !slices.Equal(m.Handles, []string{"bounds"}) || !slices.Equal(m.Failed, []string{probeDir + "/TestClampHoldsTheBounds"}) {
				t.Errorf("killed %s:%d %s: failed %v handles %v, want TestClampHoldsTheBounds and bounds", m.File, m.Line, m.Op, m.Failed, m.Handles)
			}
			bounds++
		case Survived:
			survivors++
			if m.File == "" || m.Line == 0 || m.Func == "" || m.Op == "" || m.Orig == "" {
				t.Errorf("survivor %+v lacks its place or change", m.Mutant)
			}
		}
	}
	if bounds == 0 || survivors == 0 || e.Killed != bounds || e.Survived != survivors {
		t.Errorf("killed %d (tally %d), survived %d (tally %d): want both", bounds, e.Killed, survivors, e.Survived)
	}
	for _, s := range e.Scenarios {
		if s.Handle == "bounds" && s.Killed != bounds {
			t.Errorf("bounds killed %d, want %d", s.Killed, bounds)
		}
	}
}

// A kill is credited to the scenarios of the subtests that failed, not to
// those of a sibling subtest that passed.
func TestAKillCountsForTheSubtestsThatFailed(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "killed-by-citing-test", "hollow-scenario")
	root := repoRoot(t)
	reg := probeRegistry(t, root)
	sw4 := Entry{ID: "SW-004", Name: "SW_004", Text: "SW-004 title", Level: levelSpec, Status: "active",
		Scenarios: []Scenario{{Handle: "bounds"}, {Handle: "label"}}, Anchors: []Anchor{probeAnchors(t, root)["Clamp"]}}
	sw4.Refs = []Ref{
		{File: probeDir + "/probe_test.go", Func: "TestClampAndLabel", Subtest: []string{"clamp_holds"}, Handles: []string{"bounds"}},
		{File: probeDir + "/probe_test.go", Func: "TestClampAndLabel", Subtest: []string{"label"}, Handles: []string{"label"}},
	}
	reg.Entries = append(reg.Entries, sw4)
	rep, err := MutateRegistry(root, reg, MutateOptions{IDs: []string{"SW-004"}, GivenOnly: true, Workers: 1, Cache: "shared",
		Given: []GivenMutant{{Entry: "SW-004", File: probeDir + "/probe.go", Old: "return v\n", New: "return v + 1\n"}}})
	if err != nil {
		t.Fatal(err)
	}
	e := rep.Entries[0]
	if len(e.Mutants) != 1 || e.Mutants[0].Outcome != Killed || !slices.Equal(e.Mutants[0].Handles, []string{"bounds"}) {
		t.Fatalf("mutants %+v: want one, killed, credited to bounds alone", e.Mutants)
	}
	for _, s := range e.Scenarios {
		if (s.Handle == "bounds") == s.Hollow {
			t.Errorf("scenario %s hollow %v: want label hollow and bounds not", s.Handle, s.Hollow)
		}
	}
}

func TestAScenarioWhoseTestsKillNothingIsHollow(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "hollow-scenario")
	e := runProbe(t)
	hollow := map[string]bool{}
	for _, s := range e.Scenarios {
		hollow[s.Handle] = s.Hollow
	}
	if !hollow["label"] || hollow["bounds"] {
		t.Errorf("hollow %v, want label hollow and bounds not", hollow)
	}
	md := probeRun.Markdown()
	if !strings.Contains(md, "`label`: hollow") || !strings.Contains(md, "Survivors:") {
		t.Errorf("the markdown names no hollow scenario or survivor:\n%s", md)
	}
}

func TestAMutantThatDoesNotBuildIsCountedApart(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "unbuilt-apart")
	e := runProbe(t)
	for _, m := range e.Mutants {
		if m.Note == "does not build" && m.Outcome != Unbuilt {
			t.Errorf("the unbuildable mutant is %s, want unbuilt", m.Outcome)
		}
	}
	if e.Unbuilt < 1 || e.Killed+e.Survived+e.Unbuilt != len(e.Mutants) {
		t.Errorf("killed %d survived %d unbuilt %d of %d", e.Killed, e.Survived, e.Unbuilt, len(e.Mutants))
	}
}

func TestAGivenMutantRunsBesideTheOthers(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "given-mutants-run")
	e := runProbe(t)
	var found bool
	for _, m := range e.Mutants {
		if m.Note == "off by one inside" {
			found = true
			if m.Op != "given" || m.Func != "Clamp" || m.Outcome != Killed {
				t.Errorf("given mutant %+v: %s, want a given one in Clamp, killed", m.Mutant, m.Outcome)
			}
		}
	}
	if !found {
		t.Fatal("the given mutant did not run")
	}
	for _, m := range e.Mutants {
		if m.Note == "crashes the test binary" && (m.Outcome != Killed || !slices.Equal(m.Failed, []string{probeDir + "/TestClampHoldsTheBounds"})) {
			t.Errorf("a mutant that crashes the test binary: %s, failed %v; want killed, naming the test it crashed under", m.Outcome, m.Failed)
		}
	}
	root := repoRoot(t)
	plans, err := planMutants(root, probeRegistry(t, root), MutateOptions{GivenOnly: true, Given: []GivenMutant{
		{Entry: "SW-001", File: probeDir + "/probe.go", Old: "return v\n", New: "return v + 1\n"},
	}}, nil)
	if err != nil || len(plans) != 1 || len(plans[0].mutants) != 1 || plans[0].mutants[0].Op != "given" {
		t.Errorf("given only: %v, %+v: want the one given mutant alone", err, plans)
	}
	_, err = planMutants(root, probeRegistry(t, root), MutateOptions{Given: []GivenMutant{
		{Entry: "SW-001", File: probeDir + "/probe.go", Old: "return", New: "panic(0)"},
	}}, nil)
	if err == nil || !strings.Contains(err.Error(), probeDir+"/probe.go") || !strings.Contains(err.Error(), "occurs 5 times") {
		t.Errorf("a given text that occurs five times: %v, want a refusal naming the file", err)
	}
}

func TestAnAnchorOnAVariableIsListedNotMutated(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "unmutable-anchors-listed")
	e := runProbe(t)
	if len(e.Unmutable) != 1 || e.Unmutable[0].Decl != "Limit" {
		t.Errorf("unmutable %+v, want the anchor on Limit", e.Unmutable)
	}
	if !strings.Contains(probeRun.Markdown(), "anchors `Limit`, not a function") {
		t.Error("the markdown does not list the anchor on Limit")
	}
}

func TestAnEntryWhoseTestsFailUnmutatedRunsNoMutant(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "unsteady-not-mutated")
	t.Setenv("MUTATEPROBE_UNSTEADY", "1")
	root := repoRoot(t)
	rep, err := MutateRegistry(root, probeRegistry(t, root), MutateOptions{IDs: []string{"SW-002"}, Workers: 2, Cache: "shared"})
	if err != nil {
		t.Fatal(err)
	}
	e := rep.Entries[0]
	if !slices.Equal(e.Unsteady, []string{probeDir + "/TestUnsteady"}) || len(e.Mutants) != 0 {
		t.Errorf("unsteady %v, mutants %d: want TestUnsteady named and no mutant run", e.Unsteady, len(e.Mutants))
	}
	if !strings.Contains(rep.Markdown(), "Unsteady") {
		t.Error("the markdown does not say the entry is unsteady")
	}

	// Citing tests that do not build are unsteady too, saying why.
	reg := probeRegistry(t, root)
	reg.Entries[1].Refs = []Ref{{File: "chartertest/internal/nosuch/x_test.go", Func: "TestNothing", Handles: []string{"steady"}}}
	rep, err = MutateRegistry(root, reg, MutateOptions{IDs: []string{"SW-002"}, Workers: 1, Cache: "shared"})
	if err != nil {
		t.Fatal(err)
	}
	if u := rep.Entries[0].Unsteady; len(u) != 1 || !strings.Contains(u[0], "did not build or run unmutated") || !strings.Contains(u[0], "nosuch") || len(rep.Entries[0].Mutants) != 0 {
		t.Errorf("unsteady %v, mutants %d: want the reason, naming the missing package, and no mutant run", u, len(rep.Entries[0].Mutants))
	}
}

func TestARunLeavesTheTreeAsItWasOrFailsNamingTheFile(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "tree-untouched")
	root := repoRoot(t)
	// The shared run wrote nothing.
	runProbe(t)
	if len(probeRun.Dirty) != 0 {
		t.Errorf("the run over SW-001 left %v", probeRun.Dirty)
	}
	src, err := os.ReadFile(filepath.Join(root, probeDir, "probe.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), "if v < lo {") || strings.Contains(string(src), "_ = 0") {
		t.Error("probe.go was written by the run")
	}

	litter := filepath.Join(root, probeDir, "litter.txt")
	t.Cleanup(func() { _ = os.Remove(litter) })
	t.Setenv("MUTATEPROBE_LITTER", litter)
	rep, err := MutateRegistry(root, probeRegistry(t, root), MutateOptions{IDs: []string{"SW-003"}, Max: 1, Workers: 1, Cache: "shared"})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(rep.Dirty, probeDir+"/litter.txt") {
		t.Errorf("dirty %v, want the litter named", rep.Dirty)
	}
	if !strings.Contains(rep.Markdown(), "litter.txt") {
		t.Error("the markdown does not name the litter")
	}
}

func TestJudgeReadsTheStream(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "killed-by-citing-test", "unbuilt-apart")
	for _, c := range []struct {
		name, stream string
		deadline     bool
		want         MutantOutcome
		failed       []string
	}{
		{"pass", `{"Action":"run","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"pass","Package":"m/p","Test":"TestA"}`, false, Survived, nil},
		{"fail", `{"Action":"run","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"fail","Package":"m/p","Test":"TestA/sub"}` + "\n" + `{"Action":"fail","Package":"m/p","Test":"TestA"}`, false, Killed, []string{"m/p/TestA", "m/p/TestA/sub"}},
		{"build", `{"ImportPath":"m/p.test","Action":"build-fail"}` + "\n" + `{"Action":"fail","Package":"m/p","FailedBuild":"m/p.test"}`, false, Unbuilt, nil},
		{"nothing ran", ``, false, Unbuilt, nil},
		{"crash", `{"Action":"run","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"output","Package":"m/p","Test":"TestA","Output":"fatal error: sync: unlock of unlocked mutex\n"}` + "\n" + `{"Action":"fail","Package":"m/p"}`, false, Killed, []string{"m/p/TestA"}},
		{"crash after a pass", `{"Action":"run","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"pass","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"run","Package":"m/p","Test":"TestB"}` + "\n" + `{"Action":"fail","Package":"m/p"}`, false, Killed, []string{"m/p/TestB"}},
		{"not json", "PASS\nok  m/p\n" + `{"Action":"run","Package":"m/p","Test":"TestA"}`, false, Survived, nil},
		{"deadline", `{"Action":"run","Package":"m/p","Test":"TestA"}`, true, Killed, nil},
		{"timeout", `{"Action":"run","Package":"m/p","Test":"TestA"}` + "\n" + `{"Action":"output","Package":"m/p","Output":"panic: test timed out after 1s\n"}`, false, Killed, nil},
	} {
		got, _, failed := judge([]byte(c.stream), nil, c.deadline)
		if got != c.want || !slices.Equal(failed, c.failed) {
			t.Errorf("%s: %s %v, want %s %v", c.name, got, failed, c.want, c.failed)
		}
	}
}

func TestEveryOperatorMakesItsMutant(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "mutants-in-anchored-code")
	src := `package p

func other() bool { return 1 < 2 }

// f has one site for each operator.
//
//charter:realizes SW-001
func f(a, b int, err error) (int, error) {
	if a < b && b > 0 {
		g()
	}
	if a == b {
	}
	for i := 0; i < 3; i++ {
		if i == 1 {
			break
		}
		continue
	}
	ok := true
	_ = ok
	h := func() int { return a * b }
	_ = h
	if err != nil {
		return 0, err
	}
	return a + b, nil
}

func g() {}
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "p.go"), []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	line := strings.Count(src[:strings.Index(src, "//charter:realizes")], "\n") + 1
	fn, err := locateFunc(dir, Anchor{Kind: "realizes", File: "p.go", Line: line, Decl: "f"})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][]string{}
	for _, m := range mutantsOf("SW-001", fn, nil) {
		got[m.Op] = append(got[m.Op], m.Orig+" → "+m.Repl)
		mutated := string(m.apply([]byte(src)))
		if mutated == src || !strings.Contains(mutated, "func other() bool { return 1 < 2 }") {
			t.Errorf("%s %q changed nothing, or changed code outside f", m.Op, m.Orig)
		}
	}
	for op, want := range map[string][]string{
		"operator":         {"< → <=", "< → >=", "&& → ||", "> → >=", "> → <=", "== → !=", "++ → --", "* → /", "+ → -", "!= → =="},
		"negate-condition": {"a < b && b > 0 → !(a < b && b > 0)", "a == b → !(a == b)", "i < 3 → !(i < 3)", "err != nil → !(err != nil)"},
		"empty-branch":     {"g() → _ = 0", "break → _ = 0", "return 0, err → _ = 0"},
		"drop-call":        {"g() → _ = 0"},
		"nil-error":        {"err → nil"},
		"flip-bool":        {"true → false"},
		"int-plus-one":     {"0 → 1", "3 → 4", "1 → 2", "0 → 1"},
		"break-continue":   {"break → continue", "continue → break"},
	} {
		for _, w := range want {
			if !slices.Contains(got[op], w) {
				t.Errorf("%s lacks %q; it made %q", op, w, got[op])
			}
		}
	}
	if n := len(got["empty-branch"]); n != 3 {
		t.Errorf("empty-branch made %d mutants, want 3: an empty branch has nothing to empty", n)
	}

	// Lines outside the set make no mutant.
	for _, m := range mutantsOf("SW-001", fn, lineSet{{line + 2, line + 2}}) {
		if m.Line != line+2 {
			t.Errorf("mutant on line %d, want only %d", m.Line, line+2)
		}
	}

	// A cap draws the same mutants for the same seed, and keeps them in order.
	all := mutantsOf("SW-001", fn, nil)
	a, b := sample(all, 5, 7), sample(all, 5, 7)
	if len(a) != 5 || len(sample(all, 0, 7)) != len(all) {
		t.Fatalf("sample kept %d of %d, want 5; and all with no cap", len(a), len(all))
	}
	for i := range a {
		if a[i].Line != b[i].Line || a[i].Op != b[i].Op || a[i].Orig != b[i].Orig || (i > 0 && a[i].Line < a[i-1].Line) {
			t.Fatalf("two draws with one seed differ, or are out of order: %v / %v", a, b)
		}
	}
}

func TestOnlyActiveSpecsByTheirIDsAreMutated(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "mutants-in-anchored-code")
	root := repoRoot(t)
	reg := probeRegistry(t, root)
	if _, err := planMutants(root, reg, MutateOptions{IDs: []string{"SW-404"}}, nil); err == nil || !strings.Contains(err.Error(), "SW-404") {
		t.Errorf("an unknown id: %v, want it refused by name", err)
	}
	if _, err := planMutants(root, reg, MutateOptions{Given: []GivenMutant{{Entry: "SW-404", File: probeDir + "/probe.go", Old: "return v\n", New: "return 0\n"}}}, nil); err == nil || !strings.Contains(err.Error(), "SW-404") {
		t.Errorf("a given mutant for an unknown id: %v, want it refused by name", err)
	}
	reg.Entries[1].Status = "planned"
	plans, err := planMutants(root, reg, MutateOptions{IDs: []string{"SW-001", "SW-002"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0].entry.ID != "SW-001" {
		t.Errorf("%d plans, want SW-001 alone: SW-002 is planned", len(plans))
	}
	// A diff since HEAD on a tree whose probe is unchanged touches nothing.
	rep, err := MutateRegistry(root, reg, MutateOptions{Since: "HEAD", Cache: "shared"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Entries) != 0 {
		t.Errorf("a diff that changed no probe line mutated %d entries", len(rep.Entries))
	}
}

func TestAnEntryWhoseMutantsAllSurviveHasHollowScenarios(t *testing.T) {
	Verifies(t, claims.CH_261005_VFJR, "hollow-scenario")
	e := &Entry{ID: "SW-001", Scenarios: []Scenario{{Handle: "a"}, {Handle: "b"}}}
	p := plan{entry: e, tests: []citingTest{{pkg: "p", name: "TestA", cites: []Ref{{Handles: []string{"a"}}}}}}
	em := EntryMutation{Mutants: []MutantResult{{Outcome: Survived}, {Outcome: Survived}, {Outcome: Unbuilt}}}
	tally(&em, p)
	if em.Survived != 2 || em.Unbuilt != 1 || !em.Scenarios[0].Hollow || em.Scenarios[1].Hollow || em.Scenarios[1].Cited {
		t.Errorf("%+v: want a hollow (cited, killed nothing) and b neither cited nor hollow", em)
	}
	em = EntryMutation{Mutants: []MutantResult{{Outcome: Unbuilt}}}
	tally(&em, p)
	if em.Scenarios[0].Hollow {
		t.Error("a scenario is hollow though no mutant built")
	}
}

// A cache the caller names is the caller's: a run past the limit empties
// only the cache it made itself.
func TestPoolEmptiesOnlyItsOwnCache(t *testing.T) {
	for _, own := range []bool{true, false} {
		cache := t.TempDir()
		if err := os.WriteFile(filepath.Join(cache, "entry"), []byte("0123456789"), 0o600); err != nil {
			t.Fatal(err)
		}
		mo := MutateOptions{Workers: 1, CacheLimit: 1}
		if !own {
			mo.Cache = cache
		}
		r := &runner{mo: mo}
		r.pool(1, cache, func(int) {})
		_, err := os.Stat(cache)
		if own && err == nil {
			t.Errorf("the run's own cache past the limit was kept")
		}
		if !own && err != nil {
			t.Errorf("the caller's cache %s was removed: %v", cache, err)
		}
	}
}
