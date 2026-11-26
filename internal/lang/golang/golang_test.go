package golang

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
)

// write writes files under root, by slash path.
func write(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, src := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const fixtureArea = "# Refill\n\n" +
	"## SW-010 The bucket refills at the rate\n" +
	"Status: active\nKind: invariant\n\nAt each request the bucket reads the clock.\n\nWhy: a caller that waits is let through.\n\n" +
	"Scenarios:\n1. `waited-request-passes`: a request after a wait passes.\n2. `burst-caps-tokens`: a long wait earns no more than the burst.\n\n" +
	"## SW-011 A planned rule\nStatus: planned\n\nNot yet built.\n\nWhy: someone needs it.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\n\nNo longer.\n\nWhy: it was replaced.\n"

const fixtureObservations = "## LockIsAdvisory\nSubject: os\nStatus: documented\nSource: man flock\n\nflock locks are advisory.\n"

// readFixture reads Go files written into a fresh tree whose module is
// example.com/fx and whose claims module is example.com/fx/claims.
func readFixture(t *testing.T, files map[string]string) *facts.Set {
	t.Helper()
	root := t.TempDir()
	files["go.mod"] = "module example.com/fx\n\ngo 1.21\n"
	write(t, root, files)
	var gofiles []string
	for rel := range files {
		if strings.HasSuffix(rel, ".go") {
			gofiles = append(gofiles, rel)
		}
	}
	slices.Sort(gofiles)
	set, err := Reader{Root: root, Claims: "example.com/fx/claims", ClaimsDir: "claims", EnvDir: "env"}.Read(gofiles)
	if err != nil {
		t.Fatal(err)
	}
	return set
}

const importClaims = "import (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\n"

func TestAGoDirectiveAnchorsTheDeclarationItDocuments(t *testing.T) {
	claims.Verifies(t, claims.CH_021, "go-anchor-recorded", "spaced-directive-counts", "comment-elsewhere-refused")
	set := readFixture(t, map[string]string{
		"pkg/a/a.go": "// Package a refills.\n//\n//charter:realizes SW-010\npackage a\n\n" +
			"// Refill adds the tokens earned.\n//\n// charter:realizes SW-010, SW-012\nfunc (b *Bucket) Refill() {}\n\n" +
			"// Bucket holds tokens.\n//\n//charter:concerns SW-010\ntype Bucket struct{}\n\n" +
			"func f() {\n\t//charter:realizes SW-010\n}\n",
	})
	want := []facts.Anchor{
		{At: facts.At{File: "pkg/a/a.go", Line: 3}, Kind: "realizes", IDs: []string{"SW-010"}, Decl: "package a", DeclKind: "package"},
		{At: facts.At{File: "pkg/a/a.go", Line: 8}, Kind: "realizes", IDs: []string{"SW-010", "SW-012"}, Decl: "func Bucket.Refill", DeclKind: "func"},
		{At: facts.At{File: "pkg/a/a.go", Line: 13}, Kind: "concerns", IDs: []string{"SW-010"}, Decl: "type Bucket", DeclKind: "type"},
	}
	if len(set.Anchors) != len(want) {
		t.Fatalf("anchors = %+v", set.Anchors)
	}
	for i := range want {
		if a := set.Anchors[i]; a.At != want[i].At || a.Kind != want[i].Kind || !slices.Equal(a.IDs, want[i].IDs) || a.Decl != want[i].Decl || a.DeclKind != want[i].DeclKind {
			t.Errorf("anchor %d = %+v, want %+v", i, a, want[i])
		}
	}
	if len(set.Refused) != 1 || set.Refused[0].Pos() != "pkg/a/a.go:17" || !strings.Contains(set.Refused[0].Message, "this comment documents none") {
		t.Errorf("the directive inside a function: refused %+v", set.Refused)
	}
}

func TestAGoCitationIsReadFromTheTestItOpens(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "go-citation-read", "alias-followed", "reference-is-not-citation", "later-statement-refused")
	set := readFixture(t, map[string]string{
		"pkg/a/a_test.go": "package a\n\n" + strings.Replace(importClaims, "\"example.com/fx/claims\"", "spec \"example.com/fx/claims\"", 1) +
			"func TestRefill(t *testing.T) {\n\tspec.Verifies(t, spec.SW_010, \"waited-request-passes\", \"burst-caps-tokens\")\n\t_ = spec.SW_011\n" +
			"\tt.Run(\"sub\", func(t *testing.T) {\n\t\tspec.Verifies(t, spec.SW_012)\n\t})\n" +
			"\tfor range 2 {\n\t\tt.Run(\"each\", func(t *testing.T) {\n\t\t\tspec.Verifies(t, spec.SW_010, \"burst-caps-tokens\")\n\t\t})\n\t}\n}\n\n" +
			"func TestLater(t *testing.T) {\n\tt.Log(\"first\")\n\tspec.Verifies(t, spec.SW_010)\n}\n",
	})
	if len(set.Citations) != 3 {
		t.Fatalf("citations = %+v", set.Citations)
	}
	c := set.Citations[0]
	if c.Symbol != "SW_010" || c.Test != "TestRefill" || c.Pos() != "pkg/a/a_test.go:10" || !slices.Equal(c.Handles, []string{"waited-request-passes", "burst-caps-tokens"}) || c.Observes {
		t.Errorf("citation = %+v", c)
	}
	if set.Citations[1].Symbol != "SW_012" || set.Citations[2].Pos() != "pkg/a/a_test.go:17" {
		t.Errorf("the subtests' citations = %+v", set.Citations[1:])
	}
	if len(set.References) != 1 || set.References[0].Symbol != "SW_011" {
		t.Errorf("references = %+v; want the bare SW_011 alone", set.References)
	}
	if len(set.Refused) != 1 || set.Refused[0].Pos() != "pkg/a/a_test.go:24" || !strings.Contains(set.Refused[0].Message, "is the leading statement of its test") {
		t.Errorf("refused = %+v", set.Refused)
	}
	if !slices.Equal(set.TestFiles, []string{"pkg/a/a_test.go"}) {
		t.Errorf("test files = %v", set.TestFiles)
	}
}

func TestAGoCitationOutsideARunTestIsRefused(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "outside-a-test-refused", "non-test-name-refused", "literal-arguments")
	for _, tc := range []struct{ body, want string }{
		{"func helper(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010)\n}\n", "Verifies sits outside a test function"},
		{"func Testrefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010)\n}\n", "Verifies sits in Testrefill, which go test never runs"},
		{"type s struct{}\n\nfunc (s) TestRefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010)\n}\n", "Verifies sits in the method TestRefill"},
		{"func TestRefill(t *testing.T) {\n\tspec := claims.SW_010\n\tclaims.Verifies(t, spec)\n}\n", "is the leading statement"},
		{"func TestRefill(t *testing.T) {\n\tclaims.Verifies(t, nil)\n}\n", "Verifies names a spec by its symbol"},
		{"func TestRefill(t *testing.T) {\n\th := \"x\"\n\t_ = h\n}\n\nfunc TestB(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, handle)\n}\n\nvar handle = \"waited-request-passes\"\n", "a scenario handle is a string literal"},
	} {
		set := readFixture(t, map[string]string{"pkg/a/a_test.go": "package a\n\n" + importClaims + tc.body})
		found := false
		for _, r := range set.Refused {
			found = found || strings.Contains(r.Message, tc.want)
		}
		if !found {
			t.Errorf("want %q refused, got %+v", tc.want, set.Refused)
		}
	}
}

func TestATestMayNeitherAssignNorDotImportTheClaims(t *testing.T) {
	claims.Verifies(t, claims.CH_025, "assignment-refused", "dot-import-refused")
	set := readFixture(t, map[string]string{
		"pkg/a/a_test.go": "package a\n\n" + importClaims + "func TestA(t *testing.T) {\n\tclaims.SW_010.Title = \"x\"\n}\n",
		"pkg/b/b_test.go": "package b\n\nimport . \"example.com/fx/claims\"\n\nvar _ = SW_010\n",
	})
	var msgs []string
	for _, r := range set.Refused {
		msgs = append(msgs, r.Rule+": "+r.Message)
	}
	got := strings.Join(msgs, "\n")
	if !strings.Contains(got, "rule 1: SW_010 is assigned outside its declaration") || !strings.Contains(got, "rule 10: the claims package must not be dot-imported") {
		t.Errorf("refused:\n%s", got)
	}
}

func TestEnvironmentReadsAreFacts(t *testing.T) {
	claims.Verifies(t, claims.CH_034, "literal-argument", "variable-argument", "dynamic-directive")
	set := readFixture(t, map[string]string{
		"env/env.go": "package env\n\nconst (\n\tRate = \"FX_RATE\"\n\ttyped string = \"FX_TYPED\"\n)\n",
		"cmd/a/a.go": "package main\n\nimport (\n\t\"os\"\n\n\t\"example.com/fx/env\"\n)\n\nfunc main() {\n\t_ = os.Getenv(env.Rate)\n\t_ = os.Getenv(\"FX_BURST\")\n\tv := \"x\"\n\t_ = os.Getenv(v)\n\t_ = os.Getenv(v) //charter:dynamic\n}\n",
	})
	if len(set.EnvConsts) != 1 || set.EnvConsts[0].Name != "env.Rate" || set.EnvConsts[0].Value != "FX_RATE" {
		t.Errorf("constants = %+v", set.EnvConsts)
	}
	if len(set.EnvReads) != 4 {
		t.Fatalf("reads = %+v", set.EnvReads)
	}
	if r := set.EnvReads[0]; r.Const != "env.Rate" {
		t.Errorf("a read through the constant = %+v", r)
	}
	if r := set.EnvReads[1]; r.Literal != "FX_BURST" || r.Const != "" {
		t.Errorf("a read by literal = %+v", r)
	}
	if r := set.EnvReads[2]; r.Const != "" || r.Literal != "" || r.Dynamic {
		t.Errorf("a read through a variable = %+v", r)
	}
	if r := set.EnvReads[3]; !r.Dynamic {
		t.Errorf("a read under the directive = %+v", r)
	}
	if len(set.Literals) != 1 || set.Literals[0].Value != "FX_BURST" {
		t.Errorf("literals = %+v", set.Literals)
	}
}

// generated writes the claims module of the fixture record into root.
func generated(t *testing.T, root string) (gomod, specs []byte) {
	t.Helper()
	write(t, root, map[string]string{"claims/refill.md": fixtureArea, "claims/observations.md": fixtureObservations})
	rec, err := record.Read(root, "claims")
	if err != nil || len(rec.Problems) != 0 {
		t.Fatalf("record: %v %q", err, rec.Problems)
	}
	gomod, specs, err = Generate(rec, "example.com/fx/claims", "charter gen")
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, map[string]string{"claims/go.mod": string(gomod), "claims/specs.go": string(specs)})
	return gomod, specs
}

func TestTheGeneratedModuleDeclaresEveryEntry(t *testing.T) {
	claims.Verifies(t, claims.CH_020, "bare-symbol-declared", "planned-and-deprecated-marked", "title-only-comment", "go-module-self-contained")
	gomod, specs := generated(t, t.TempDir())
	src := string(specs)
	for _, want := range []string{
		"// SW_010 is SW-010: The bucket refills at the rate\nvar SW_010 = &Spec{ID: \"SW-010\", Title: \"The bucket refills at the rate\", Status: Active, Kind: \"invariant\"}",
		"var SW_011 = &Spec{ID: \"SW-011\", Title: \"A planned rule\", Status: Planned}",
		"// Deprecated: SW-012 is retired; see its Changes in claims/refill.md. Add no new anchors.\nvar SW_012",
		"var LockIsAdvisory = &Observation{Name: \"LockIsAdvisory\", Subject: \"os\", Status: \"documented\"}",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("specs.go lacks:\n%s\n\n%s", want, src)
		}
	}
	if strings.Contains(src, "At each request") {
		t.Error("specs.go holds an entry's text; a symbol's comment is its title")
	}
	if strings.Contains(string(gomod), "require") || !strings.Contains(string(gomod), "module example.com/fx/claims\n\ngo "+GoVersion) {
		t.Errorf("go.mod = %s", gomod)
	}
	for _, imp := range []string{"\"encoding/json\"", "\"os\"", "\"runtime\"", "\"sync\"", "\"testing\""} {
		src = strings.Replace(src, imp, "", 1)
	}
	if i := strings.Index(src, "import ("); i < 0 || strings.Contains(src[i:i+strings.Index(src[i:], ")")], "\"") {
		t.Errorf("specs.go imports beyond the standard library:\n%s", src)
	}
}

// goTest runs go test in dir with the environment added.
func goTest(t *testing.T, dir string, env ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), append([]string{"GOFLAGS=-mod=mod", "GOWORK=off", "GOTOOLCHAIN=local"}, env...)...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

type line struct {
	ID, Handle, Test, File, Outcome string
}

func readEvents(t *testing.T, path string) []line {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []line
	s := bufio.NewScanner(f)
	for s.Scan() {
		var l line
		if err := json.Unmarshal(s.Bytes(), &l); err != nil {
			t.Fatalf("a line of the events file is not a JSON object: %q: %v", s.Text(), err)
		}
		out = append(out, l)
	}
	return out
}

func TestAGoTestRecordsItsOutcomeInTheEventsFile(t *testing.T) {
	claims.Verifies(t, claims.CH_024, "event-per-citation", "outcome-recorded", "unset-records-nothing", "concurrent-lines-whole", "nil-spec-refused")
	if testing.Short() {
		t.Skip("builds and runs a module")
	}
	root := t.TempDir()
	generated(t, root)
	var parallel strings.Builder
	for i := range 40 {
		parallel.WriteString("\tt.Run(\"p" + string(rune('a'+i%26)) + string(rune('a'+i/26)) + "\", func(t *testing.T) {\n\t\tclaims.Verifies(t, claims.SW_010, \"waited-request-passes\", \"burst-caps-tokens\")\n\t\tt.Parallel()\n\t})\n")
	}
	write(t, root, map[string]string{
		"go.mod": "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"pkg/a/a_test.go": "package a\n\n" + importClaims +
			"func TestPasses(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"waited-request-passes\")\n}\n\n" +
			"func TestFails(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"burst-caps-tokens\")\n\tif claims.SW_010.ID != \"\" && false {\n\t\treturn\n\t}\n\tif claims.SW_010 != nil {\n\t\tt.Error(\"failed on purpose\")\n\t}\n}\n\n" +
			"func TestSkips(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_011)\n\tt.Skip(\"skipped on purpose\")\n}\n\n" +
			"func TestParallel(t *testing.T) {\n" + parallel.String() + "}\n",
		"pkg/b/b_test.go": "package b\n\n" + importClaims + "func TestNil(t *testing.T) {\n\tclaims.Verifies(t, nil)\n}\n",
	})
	events := filepath.Join(t.TempDir(), "events.jsonl")
	out, _ := goTest(t, root, "CHARTER_EVENTS="+events)
	if !strings.Contains(out, "claims: Verifies needs a spec") {
		t.Errorf("a nil spec did not fail its test:\n%s", out)
	}
	lines := readEvents(t, events)
	find := func(test, handle string) *line {
		for i := range lines {
			if lines[i].Test == test && lines[i].Handle == handle {
				return &lines[i]
			}
		}
		return nil
	}
	if l := find("TestPasses", "waited-request-passes"); l == nil || l.ID != "SW-010" || l.Outcome != "passed" || !strings.HasSuffix(l.File, "pkg/a/a_test.go") {
		t.Errorf("TestPasses: %+v", l)
	}
	if l := find("TestFails", "burst-caps-tokens"); l == nil || l.Outcome != "failed" {
		t.Errorf("TestFails: %+v", l)
	}
	if l := find("TestSkips", ""); l == nil || l.ID != "SW-011" || l.Outcome != "skipped" {
		t.Errorf("TestSkips: %+v", l)
	}
	n := 0
	for _, l := range lines {
		if strings.HasPrefix(l.Test, "TestParallel/") {
			n++
		}
	}
	if n != 80 {
		t.Errorf("the parallel subtests wrote %d whole lines, want 80", n)
	}

	unset, _ := goTest(t, root)
	if strings.Contains(unset, "events file") {
		t.Errorf("with CHARTER_EVENTS unset a test touched the events file:\n%s", unset)
	}
	if !strings.Contains(unset, "--- FAIL: TestFails") || strings.Contains(unset, "--- FAIL: TestPasses") {
		t.Errorf("with CHARTER_EVENTS unset the tests did not run as they would:\n%s", unset)
	}
}

func TestAGoCitationABlockEnclosesOrAVariableNamesIsRefused(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "later-statement-refused", "literal-arguments")
	for _, tc := range []struct{ body, want string }{
		{"func TestRefill(t *testing.T) {\n\tif true {\n\t\tclaims.Verifies(t, claims.SW_010)\n\t}\n}\n", "Verifies is the leading statement of its test"},
		{"func TestRefill(t *testing.T) {\n\tclaims.Verifies(t, spec)\n}\n\nvar spec = claims.SW_010\n", "Verifies names a spec by its symbol"},
	} {
		set := readFixture(t, map[string]string{"pkg/a/a_test.go": "package a\n\n" + importClaims + tc.body})
		if len(set.Citations) != 0 {
			t.Errorf("counted %+v for\n%s", set.Citations, tc.body)
		}
		found := false
		for _, r := range set.Refused {
			found = found || strings.Contains(r.Message, tc.want)
		}
		if !found {
			t.Errorf("want %q refused, got %+v", tc.want, set.Refused)
		}
	}
}
