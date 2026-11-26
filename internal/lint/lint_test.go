package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
)

const area = "# Refill\n\nHow the bucket refills.\n\n" +
	"## SW-010 The bucket refills at the rate\nStatus: active\nServes: STK-001\n\nAt each request the bucket reads the clock.\n\nWhy: a caller that waits is let through again, so a burst does not punish the patient.\n\n" +
	"Scenarios:\n1. `waited-request-passes`: a request after a wait passes.\n2. `burst-caps-tokens`: a long wait earns no more than the burst.\n\n" +
	"## SW-011 A rule that is planned\nStatus: planned\nServes: STK-001\n\nNot yet built.\n\nWhy: someone needs it, and the reason is long enough to stand on its own here.\n\nScenarios:\n1. `planned-rule`: it rules.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\nServes: STK-001\n\nNo longer.\n\nWhy: it was replaced by something better that serves the same need.\n"

const needs = "# Needs\n\n## STK-001 A caller that waits is let through\nStatus: active\n\nIt is let through.\n\nWhy: a caller should not be punished for waiting.\n\nScenarios:\n1. `waits`: a caller that waited passes.\n"

func files() map[string]string {
	return map[string]string{
		"go.mod":              "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"claims/charter.toml": "name = \"fx\"\n\n[go]\nmodule = \"example.com/fx/claims\"\ndir = \".\"\ntest = \"go test ./...\"\n",
		"claims/needs.md":     needs,
		"claims/refill.md":    area,
		"pkg/a/a.go":          "package a\n\n// Refill adds the tokens earned.\n//\n//charter:realizes SW-010\nfunc Refill() {}\n",
		"pkg/a/a_test.go":     "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestRefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"waited-request-passes\", \"burst-caps-tokens\")\n}\n\nfunc TestWaits(t *testing.T) {\n\tclaims.Verifies(t, claims.STK_001, \"waits\")\n}\n",
	}
}

func lint(t *testing.T, fs map[string]string, listStyle bool) []string {
	t.Helper()
	root := t.TempDir()
	for rel, src := range fs {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	c, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := check.Write(c); err != nil {
		t.Fatal(err)
	}
	out, err := Lint(c, listStyle)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func edit(t *testing.T, file, old, new string) map[string]string {
	t.Helper()
	fs := files()
	if strings.Count(fs[file], old) != 1 {
		t.Fatalf("%s holds %q %d times", file, old, strings.Count(fs[file], old))
	}
	fs[file] = strings.Replace(fs[file], old, new, 1)
	return fs
}

func has(t *testing.T, out []string, want string) {
	t.Helper()
	for _, l := range out {
		if strings.Contains(l, want) {
			return
		}
	}
	t.Errorf("want a line holding\n  %s\nin:\n  %s", want, strings.Join(out, "\n  "))
}

func TestTheCleanFixtureLintsClean(t *testing.T) {
	if out := lint(t, files(), false); len(out) != 0 {
		t.Errorf("lint of the clean fixture:\n%s", strings.Join(out, "\n"))
	}
}

func TestLintAdvises(t *testing.T) {
	claims.Verifies(t, claims.CH_041, "serves-nothing", "planned-verified", "deprecated-still-anchored", "package-clause-only", "no-scenario-named", "gate-problems-first")
	has(t, lint(t, edit(t, "claims/refill.md", "Status: active\nServes: STK-001\n", "Status: active\n"), false), "claims/refill.md:5: lint: SW-010 serves nothing")
	fs := files()
	fs["pkg/a/a_test.go"] += "\nfunc TestPlanned(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_011, \"planned-rule\")\n}\n"
	has(t, lint(t, fs, false), "claims/refill.md:17: lint: SW-011 is planned and each of its scenarios is verified; set it active")
	has(t, lint(t, edit(t, "pkg/a/a.go", "//charter:realizes SW-010", "//charter:realizes SW-010, SW-012"), false), "pkg/a/a.go:5: lint: SW-012 is deprecated and still anchored here; strip it")
	pkg := edit(t, "pkg/a/a.go", "package a\n\n// Refill adds the tokens earned.\n//\n//charter:realizes SW-010\nfunc Refill() {}\n", "// Package a refills.\n//\n//charter:realizes SW-010\npackage a\n")
	has(t, lint(t, pkg, false), "lint: SW-010 is realized only by a package clause")
	fs = files()
	fs["pkg/a/b_test.go"] = "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestWhole(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010)\n}\n"
	has(t, lint(t, fs, false), "pkg/a/b_test.go:10: lint: verifies names SW-010 without a scenario")
	has(t, lint(t, edit(t, "pkg/a/a_test.go", ", \"burst-caps-tokens\")", ")"), false), "rule 25: scenario burst-caps-tokens of SW-010 is verified by no test")
}

func TestTheStyleIsCountedAndListed(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_9K5J, "long-title-named", "long-text-named", "long-scenario-named", "pointing-why-named")
	fs := files()
	fs["claims/retired.txt"] = "quotas? -> bucket | quota defaults\n"
	fs["claims/refill.md"] = strings.Replace(area, "At each request the bucket reads the clock.", "At each request the quota reads the clock. One. Two. Three. Four. Five. Six. `e.g. code.`", 1)
	fs["claims/refill.md"] = strings.Replace(fs["claims/refill.md"], "## SW-010 The bucket refills at the rate", "## SW-010 The bucket refills at the rate it was given when the limiter started and holds the rest for later every single time", 1)
	fs["claims/refill.md"] = strings.Replace(fs["claims/refill.md"], "2. `burst-caps-tokens`: a long wait earns no more than the burst.", "2. `burst-caps-tokens`: a long wait earns no more than the burst, however long it lasts, and whatever the caller does in the meantime, so the next request after it is let through once and then denied.", 1)
	fs["claims/refill.md"] = strings.Replace(fs["claims/refill.md"], "Why: someone needs it, and the reason is long enough to stand on its own here.", "Why: as ruled in backlog 003.", 1)
	counts := lint(t, fs, false)
	for _, want := range []string{
		"claims/retired.txt:1: lint: 1 retired word is used; charter lint -style lists it",
		"claims/README.md:1: lint: 1 title runs past twenty words",
		"claims/README.md:1: lint: 1 text runs past six sentences",
		"claims/README.md:1: lint: 1 scenario sentence runs past thirty words",
		"claims/README.md:1: lint: 1 Why only points at a ruling",
	} {
		has(t, counts, want)
	}
	listed := lint(t, fs, true)
	for _, want := range []string{
		"claims/refill.md:9: lint: \"quota\" is retired in the text of SW-010; write bucket",
		"claims/refill.md:5: lint: the title of SW-010 runs to 22 words",
		"claims/refill.md:5: lint: the text of SW-010 runs to 8 sentences",
		"claims/refill.md:15: lint: scenario `burst-caps-tokens` of SW-010 runs to 34 words",
		"claims/refill.md:17: lint: the Why of SW-011 points at a ruling and gives no reason",
	} {
		has(t, listed, want)
	}
}

func TestATestTwoDecisionsShareIsNamed(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_9K5J, "shared-test-named")
	fs := edit(t, "claims/refill.md", "## SW-011 A rule that is planned\nStatus: planned", "## SW-011 A rule that is planned\nStatus: active")
	fs["pkg/a/a.go"] = strings.Replace(fs["pkg/a/a.go"], "SW-010", "SW-010, SW-011", 1)
	fs["pkg/a/a_test.go"] = strings.Replace(fs["pkg/a/a_test.go"], "\"burst-caps-tokens\")\n", "\"burst-caps-tokens\")\n\tclaims.Verifies(t, claims.SW_011, \"planned-rule\")\n", 1)
	has(t, lint(t, fs, true), "pkg/a/a_test.go:10: lint: TestRefill cites SW-010/waited-request-passes, SW-010/burst-caps-tokens and SW-011/planned-rule; two decisions share one test")
	fs["pkg/a/a_test.go"] = strings.Replace(fs["pkg/a/a_test.go"], "\tclaims.Verifies(t, claims.SW_011, \"planned-rule\")\n", "\tclaims.Verifies(t, claims.SW_011, \"planned-rule\")\n\tclaims.Verifies(t, claims.STK_001, \"waits\")\n", 1)
	for _, l := range lint(t, fs, true) {
		if strings.Contains(l, "two decisions share one test") {
			t.Errorf("a shared test that also cites a need: %s", l)
		}
	}
}

func TestLintCountsTheOpenQuestionsInOneLine(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_4B3B, "open-listed")
	fs := files()
	q := "## Q-261009-M4TK A question?\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n\n"
	fs["claims/questions.md"] = "# Questions\n\n" + q + strings.Replace(q, "M4TK", "N7PR", 1)
	has(t, lint(t, fs, false), "claims/questions.md:1: lint: 2 questions are open; charter questions lists them")
	fs["claims/questions.md"] = "# Questions\n\n" + q
	has(t, lint(t, fs, false), "claims/questions.md:1: lint: 1 question is open; charter questions lists it")
	fs["claims/questions.md"] = "# Questions\n\n" + strings.Replace(q, "Provisional: yes\n", "Provisional: yes\nAnswer: 2026-11-26, owner: yes\n", 1)
	for _, l := range lint(t, fs, false) {
		if strings.Contains(l, "questions.md") {
			t.Errorf("no question open: %s", l)
		}
	}
}

func TestARetiredWordIsNamedUnlessAPhraseKeepsIt(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_WTMZ, "retired-word-named")
	fs := files()
	fs["claims/retired.txt"] = "quotas? -> bucket | quota defaults\n"
	// The first sentence uses the word; the second holds it in code and in
	// one of its phrases, which keep it.
	fs["claims/refill.md"] = strings.Replace(area, "At each request the bucket reads the clock.", "At each request the quota reads the clock.\nIts `quota` field keeps the quota defaults.", 1)
	has(t, lint(t, fs, false), "claims/retired.txt:1: lint: 1 retired word is used; charter lint -style lists it")
	has(t, lint(t, fs, true), "claims/refill.md:9: lint: \"quota\" is retired in the text of SW-010; write bucket")
}
