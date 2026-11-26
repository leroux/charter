package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

const planned = "# Area\n\n## SW-010 A rule\nStatus: planned\n\nWhat it rules.\n\nWhy: someone needs it.\n"

func charterIn(t *testing.T, area string) string {
	t.Helper()
	root := t.TempDir()
	for rel, src := range map[string]string{"claims/charter.toml": "name = \"x\"\n", "claims/area.md": area} {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestCheckSetsItsExitStatus(t *testing.T) {
	claims.Verifies(t, claims.CH_037, "exit-status")
	var out, errOut bytes.Buffer
	if code := run([]string{"check"}, charterIn(t, planned), &out, &errOut); code != 0 || out.Len() != 0 {
		t.Errorf("a clean tree: exit %d, output %q %q", code, out.String(), errOut.String())
	}
	active := strings.Replace(planned, "Status: planned", "Status: active", 1) + "\nScenarios:\n1. `rules`: it rules.\n"
	out.Reset()
	if code := run([]string{"check"}, charterIn(t, active), &out, &errOut); code != 1 || !strings.Contains(out.String(), "claims/area.md:3: rule 25: active entry SW-010 is realized by no code") {
		t.Errorf("a tree with a problem: exit %d, output %q", code, out.String())
	}
}

func TestEveryVerbWorksOnTheCharterFoundAbove(t *testing.T) {
	claims.Verifies(t, claims.CH_040, "verbs-on-the-found-charter")
	root := charterIn(t, planned)
	deep := filepath.Join(root, "pkg", "a")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	for _, args := range [][]string{{"gen"}, {"check"}, {"lint"}, {"next", "SW"}, {"anchors", "SW-010"}, {"questions"}, {"export"}} {
		out.Reset()
		if code := run(args, deep, &out, &errOut); code != 0 {
			t.Errorf("charter %v from %s: exit %d, %s%s", args, deep, code, out.String(), errOut.String())
		}
	}
	if code := run([]string{"check"}, t.TempDir(), &out, &errOut); code != 2 || !strings.Contains(errOut.String(), "no claims/charter.toml") {
		t.Errorf("outside any charter: exit %d, %q", code, errOut.String())
	}
}

func TestQuestionsListsEachOpenQuestion(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_4B3B, "questions-verb-lists")
	root := charterIn(t, planned)
	q := "# Questions\n\n## Q-261009-M4TK Does it rule?\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n\n## Q-261009-N7PR Is it whole?\nAsked-by: owner\nEntry: SW-010\nProvisional: no\nAnswer: 2026-10-09, owner: yes\n"
	if err := os.WriteFile(filepath.Join(root, "claims", "questions.md"), []byte(q), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	run([]string{"questions"}, root, &out, &errOut)
	if want := "Q-261009-M4TK  SW-010\n  Does it rule?\n  provisional: yes\n  asked by: agent\n\n1 open\n"; out.String() != want {
		t.Errorf("questions =\n%s\nwant\n%s%s", out.String(), want, errOut.String())
	}
	out.Reset()
	run([]string{"questions", "-all"}, root, &out, &errOut)
	if !strings.Contains(out.String(), "Q-261009-N7PR  SW-010") || !strings.Contains(out.String(), "  answer: 2026-10-09, owner: yes\n") || !strings.HasSuffix(out.String(), "1 open\n") {
		t.Errorf("questions -all =\n%s", out.String())
	}
}

// goCharter is a charter of Go: an active decision, its code and its test.
func goCharter(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, src := range map[string]string{
		"go.mod":              "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"claims/charter.toml": "name = \"fx\"\n\n[go]\nmodule = \"example.com/fx/claims\"\ndir = \".\"\ntest = \"go test ./...\"\n",
		"claims/refill.md":    "# Refill\nCode: pkg/a\n\nHow the bucket refills.\n\n## SW-010 The bucket refills\nStatus: active\nServes: SW-011\n\nIt refills.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `refills` (owner, #1): it refills.\n\n## SW-011 A planned rule\nStatus: planned\n\nLater.\n\nWhy: someone needs it.\n",
		"claims/waivers.txt":  "SW-010/refills\tidle on purpose\n",
		"pkg/a/a.go":          "package a\n\n// Refill refills.\n//\n//charter:realizes SW-010\nfunc Refill() {}\n",
		"pkg/a/a_test.go":     "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestRefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"refills\")\n}\n",
	} {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestGenWritesFromTheDocumentsAlone(t *testing.T) {
	claims.Verifies(t, claims.CH_040, "gen-from-documents-alone")
	root := goCharter(t)
	// No generated file yet, and a test that cannot compile without one.
	var out, errOut bytes.Buffer
	if code := run([]string{"gen"}, root, &out, &errOut); code != 0 {
		t.Fatalf("gen with no generated file: exit %d %s", code, errOut.String())
	}
	if err := os.WriteFile(filepath.Join(root, "claims", "specs.go"), []byte("package claims\n\nthis does not parse\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"gen"}, root, &out, &errOut); code != 0 {
		t.Fatalf("gen over a broken specs.go: exit %d %s", code, errOut.String())
	}
	specs, err := os.ReadFile(filepath.Join(root, "claims", "specs.go"))
	if err != nil || !strings.Contains(string(specs), "var SW_010 = &Spec{") || strings.Contains(string(specs), "this does not parse") {
		t.Errorf("gen did not rewrite specs.go from the documents: %v\n%s", err, specs)
	}
}

func TestAnchorsListsTheCodeAndTheTestsOfAnEntry(t *testing.T) {
	claims.Verifies(t, claims.CH_002, "anchors-and-citations-listed")
	root := goCharter(t)
	var out, errOut bytes.Buffer
	run([]string{"gen"}, root, &out, &errOut)
	out.Reset()
	if code := run([]string{"anchors", "SW-010"}, root, &out, &errOut); code != 0 {
		t.Fatalf("anchors: exit %d %s", code, errOut.String())
	}
	var got map[string]struct {
		Title    string
		Realizes []struct {
			File, Decl string
			Line       int
		}
		Verifies []struct {
			File, Test string
			Handles    []string
		}
		Serves []string
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	e := got["SW-010"]
	if e.Title != "The bucket refills" || len(e.Realizes) != 1 || e.Realizes[0].File != "pkg/a/a.go" || e.Realizes[0].Decl != "func Refill" || e.Realizes[0].Line != 5 ||
		len(e.Verifies) != 1 || e.Verifies[0].Test != "TestRefill" || len(e.Verifies[0].Handles) != 1 || e.Verifies[0].Handles[0] != "refills" || len(e.Serves) != 1 {
		t.Errorf("anchors SW-010 = %+v", e)
	}
}

func TestExportGivesEveryPartOfTheRecord(t *testing.T) {
	claims.Verifies(t, claims.CH_040, "registry-lists-everything")
	root := goCharter(t)
	var out, errOut bytes.Buffer
	run([]string{"gen"}, root, &out, &errOut)
	out.Reset()
	run([]string{"export"}, root, &out, &errOut)
	var got struct {
		Areas []struct {
			File, Title, Summary string
			Code                 []string
		}
		Waivers []struct{ Pattern, Reason string }
		Entries []struct {
			ID        string
			Scenarios []struct{ Handle, Provenance string }
			Anchors   []struct{ File string }
			Refs      []struct {
				Func    string
				Handles []string
			}
		}
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	if len(got.Areas) != 1 || got.Areas[0].Title != "Refill" || got.Areas[0].Code[0] != "pkg/a" || got.Areas[0].Summary != "How the bucket refills." {
		t.Errorf("areas = %+v", got.Areas)
	}
	if len(got.Waivers) != 1 || got.Waivers[0].Reason != "idle on purpose" {
		t.Errorf("waivers = %+v", got.Waivers)
	}
	e := got.Entries[0]
	if e.ID != "SW-010" || e.Scenarios[0].Handle != "refills" || e.Scenarios[0].Provenance != "owner, #1" || len(e.Anchors) != 1 || len(e.Refs) != 1 || e.Refs[0].Func != "TestRefill" || e.Refs[0].Handles[0] != "refills" {
		t.Errorf("entry = %+v", e)
	}
}

// writeTree writes files under a fresh root.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, src := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const sharedArea = "# Refill\n\nHow the bucket refills.\n\n## SW-010 The bucket refills\nStatus: planned\n\nIt refills.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `refills` (owner, #1): it refills.\n"

const sharedNeeds = "# Needs\n\n## SW-001 A caller that waits is let through\nStatus: planned\n\nIt is let through.\n\nWhy: waiting should not be punished.\n"

// exportedAll is every entry export gives, without where it was read.
func exportedAll(t *testing.T, root string) []map[string]any {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := run([]string{"export"}, root, &out, &errOut); code != 0 {
		t.Fatalf("export in %s: exit %d %s", root, code, errOut.String())
	}
	var got struct{ Entries []map[string]any }
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	for _, e := range got.Entries {
		for _, k := range []string{"refs", "anchors", "pos", "file"} {
			delete(e, k)
		}
	}
	return got.Entries
}

func exported(t *testing.T, root string) map[string]any {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := run([]string{"export"}, root, &out, &errOut); code != 0 {
		t.Fatalf("export in %s: exit %d %s", root, code, errOut.String())
	}
	var got struct{ Entries []map[string]any }
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	e := got.Entries[0]
	for _, k := range []string{"refs", "anchors", "pos", "file"} {
		delete(e, k)
	}
	return e
}

func TestOneRecordReadsAlikeInAnyLayout(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_PFG5, "one-form-everywhere")
	goRoot := writeTree(t, map[string]string{
		"go.mod":              "module example.com/a\n\ngo 1.21\n",
		"claims/charter.toml": "name = \"a\"\n\n[go]\nmodule = \"example.com/a/claims\"\ndir = \".\"\ntest = \"go test ./...\"\n",
		"claims/refill.md":    sharedArea,
		"claims/needs.md":     sharedNeeds,
	})
	pyRoot := writeTree(t, map[string]string{
		"claims/charter.toml": "name = \"b\"\n\n[python]\ndir = \"tool\"\ntest = \"python3 -m unittest\"\n",
		"claims/refill.md":    sharedArea,
		"claims/needs.md":     sharedNeeds,
		"tool/slug.py":        "x = 1\n",
	})
	a, b := exportedAll(t, goRoot), exportedAll(t, pyRoot)
	if len(a) != 2 {
		t.Fatalf("want a need and a decision, got %d entries", len(a))
	}
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Errorf("the record read differently:\n%s\n%s", ja, jb)
	}
}

func TestAChangedEntryKeepsItsIdAndItsHistory(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_PFG5, "reference-survives")
	root := writeTree(t, map[string]string{"claims/charter.toml": "name = \"a\"\n", "claims/refill.md": sharedArea})
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-q", "-m", "base"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	changed := strings.Replace(sharedArea, "It refills.", "It refills at the rate.", 1) + "\nChanges:\n- 2026-11-26: the rate named.\n"
	if err := os.WriteFile(filepath.Join(root, "claims", "refill.md"), []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"diff", "HEAD"}, root, &out, &errOut); code != 0 {
		t.Errorf("a recorded change: diff exit %d\n%s%s", code, out.String(), errOut.String())
	}
	e := exported(t, root)
	changes, _ := e["changes"].([]any)
	if e["id"] != "SW-010" || len(changes) != 1 || changes[0].(map[string]any)["date"] != "2026-11-26" {
		t.Errorf("after the change: %v", e)
	}
}

func TestAGapBetweenTheRecordAndTheCodeStopsTheChange(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_P2SK, "gap-stops-the-change")
	for _, tc := range []struct{ file, old, new, want string }{
		{"pkg/a/a.go", "//charter:realizes SW-010\n", "", "active entry SW-010 is realized by no code"},
		{"pkg/a/a_test.go", "\tclaims.Verifies(t, claims.SW_010, \"refills\")\n", "", "scenario refills of SW-010 is verified by no test"},
	} {
		root := goCharter(t)
		if err := os.Remove(filepath.Join(root, "claims", "waivers.txt")); err != nil {
			t.Fatal(err)
		}
		var out, errOut bytes.Buffer
		run([]string{"gen"}, root, &out, &errOut)
		if code := run([]string{"check"}, root, &out, &errOut); code != 0 {
			t.Fatalf("the charter before the gap: exit %d\n%s", code, out.String())
		}
		p := filepath.Join(root, filepath.FromSlash(tc.file))
		src, _ := os.ReadFile(p)
		if err := os.WriteFile(p, []byte(strings.Replace(string(src), tc.old, tc.new, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		out.Reset()
		if code := run([]string{"check"}, root, &out, &errOut); code != 1 || !strings.Contains(out.String(), tc.want) {
			t.Errorf("with %s changed: exit %d\n%s", tc.file, code, out.String())
		}
	}
}

func TestASurfaceNoEntryOwnsStopsTheChange(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_P2SK, "unowned-surface-stops-the-change")
	root := goCharter(t)
	if err := os.Remove(filepath.Join(root, "claims", "waivers.txt")); err != nil {
		t.Fatal(err)
	}
	toml := filepath.Join(root, "claims", "charter.toml")
	src, _ := os.ReadFile(toml)
	if err := os.WriteFile(toml, append(src, []byte("env = \"env\"\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "env"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "env", "env.go"), []byte("// Package env names the variables.\npackage env\n\n// Rate is new.\nconst Rate = \"FX_RATE\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	run([]string{"gen"}, root, &out, &errOut)
	if code := run([]string{"check"}, root, &out, &errOut); code != 1 || !strings.Contains(out.String(), "rule 6: env:FX_RATE is neither accounted for by an entry nor listed in claims/baseline.txt") {
		t.Errorf("a new variable no entry owns: exit %d\n%s%s", code, out.String(), errOut.String())
	}
}

func TestReportFailsOnAScenarioThatDidNotPass(t *testing.T) {
	claims.Verifies(t, claims.CH_042, "skipped-scenario-fails")
	root := goCharter(t)
	if err := os.Remove(filepath.Join(root, "claims", "waivers.txt")); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	run([]string{"gen"}, root, &out, &errOut)
	out.Reset()
	if code := run([]string{"report", "-md"}, root, &out, &errOut); code != 0 {
		t.Fatalf("every scenario passing: exit %d\n%s%s", code, out.String(), errOut.String())
	}
	test := filepath.Join(root, "pkg", "a", "a_test.go")
	src, _ := os.ReadFile(test)
	if err := os.WriteFile(test, []byte(strings.Replace(string(src), "\"refills\")\n", "\"refills\")\n\tt.Skip(\"not yet\")\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := run([]string{"report", "-md"}, root, &out, &errOut); code != 1 || !strings.Contains(out.String(), "SW-010/refills") {
		t.Errorf("a skipped scenario: exit %d\n%s%s", code, out.String(), errOut.String())
	}
}
