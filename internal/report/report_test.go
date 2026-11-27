package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
)

var repo, _ = filepath.Abs("../..")

const area = "# Refill\n\n" +
	"## SW-010 The bucket refills at the rate\nStatus: active\n\nIt refills.\n\nWhy: a caller that waits is let through.\n\n" +
	"Scenarios:\n1. `go-passes` (owner, #1): passes in Go.\n2. `python-passes`: passes in Python.\n3. `typescript-passes`: passes in TypeScript.\n4. `one-fails`: one citing test fails.\n5. `all-skipped`: every citing test is skipped.\n6. `never-cited`: no test cites it.\n7. `waived`: waived.\n\n" +
	"## SW-011 A planned rule\nStatus: planned\n\nNot yet.\n\nWhy: someone needs it.\n\nScenarios:\n1. `later`: later.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\n\nNo longer.\n\nWhy: it was replaced.\n\nScenarios:\n1. `gone`: gone.\n"

func build(t *testing.T, extra map[string]string) *config.Config {
	t.Helper()
	root := t.TempDir()
	vitest := filepath.Join(repo, "node_modules", ".bin", "vitest")
	files := map[string]string{
		"go.mod": "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"claims/charter.toml": "name = \"fx\"\n\n[go]\nmodule = \"example.com/fx/claims\"\ndir = \".\"\ntest = \"GOFLAGS=-mod=mod GOWORK=off go test -count=1 ./...\"\n\n" +
			"[python]\ndir = \"py\"\ntest = \"PYTHONPATH=../claims python3 -m unittest -q\"\n\n" +
			"[typescript]\ndir = \"web\"\ntest = \"" + vitest + " run --globals --root .\"\ncompiler = \"" + filepath.Join(repo, "node_modules", "typescript") + "\"\nimport = \"../claims/index.ts\"\n",
		"claims/refill.md":   area,
		"claims/waivers.txt": "SW-010/waived\tno clock to drive yet\n",
		"pkg/a/a_test.go": "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\n" +
			"func TestGo(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"go-passes\", \"one-fails\")\n}\n\n" +
			"func TestSkipped(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"all-skipped\")\n\tt.Skip(\"on purpose\")\n}\n",
		"py/test_refill.py":  "import unittest\n\nimport claims\n\n\nclass T(unittest.TestCase):\n    def test_python(self):\n        claims.verifies(self, claims.SW_010, \"python-passes\", \"one-fails\")\n        self.fail(\"on purpose\")\n\n    def test_skipped(self):\n        claims.verifies(self, claims.SW_010, \"all-skipped\")\n        self.skipTest(\"on purpose\")\n",
		"web/refill.test.ts": "import { claims, verifies } from \"../claims/index.ts\";\n\ntest(\"typescript\", () => {\n  verifies(claims.SW_010, \"typescript-passes\");\n});\n",
	}
	for k, v := range extra {
		files[k] = v
	}
	for rel, src := range files {
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
	return c
}

func TestTheReportJoinsEveryLanguagesOutcomes(t *testing.T) {
	claims.Verifies(t, claims.CH_042, "results-joined", "skipped-scenario-fails", "markdown-rendered", "provenance-and-profile")
	claims.Verifies(t, claims.CH_261009_3FHW, "report-counts-who-ruled")
	rep, err := Run(build(t, nil))
	if err != nil {
		t.Fatal(err)
	}
	var e *EntryReport
	for i := range rep.Entries {
		if rep.Entries[i].ID == "SW-010" {
			e = &rep.Entries[i]
		}
	}
	if e == nil {
		t.Fatal("no SW-010 in the report")
	}
	want := map[string]Outcome{"go-passes": Passed, "python-passes": Failed, "typescript-passes": Passed, "one-fails": Failed, "all-skipped": Skipped, "never-cited": NoRun, "waived": Waived}
	for _, s := range e.Scenarios {
		if s.Outcome != want[s.Handle] {
			t.Errorf("%s = %s (%+v), want %s", s.Handle, s.Outcome, s.Tests, want[s.Handle])
		}
	}
	unverified := strings.Join(rep.Unverified, " ")
	if unverified != "SW-010/python-passes SW-010/one-fails SW-010/all-skipped SW-010/never-cited" {
		t.Errorf("unverified = %s; a planned entry is never", unverified)
	}
	if rep.Profile.Owner.Scenarios != 1 || rep.Profile.Owner.Passed != 1 || rep.Profile.Provisional.Scenarios != 6 {
		t.Errorf("profile = %+v", rep.Profile)
	}
	md := rep.Markdown()
	if !strings.HasPrefix(md, "Active scenarios by who ruled them: 1 ruled by the owner (1 passed)") {
		t.Errorf("markdown opens:\n%s", md[:min(len(md), 200)])
	}
	i, j := strings.Index(md, "4 active scenarios did not pass"), strings.Index(md, "## SW-010")
	if i < 0 || j < i || !strings.Contains(md, "- `go-passes` (owner, #1): pass · TestGo (pass)") {
		t.Errorf("markdown:\n%s", md)
	}
}

func TestATestCommandThatCannotRunIsNamed(t *testing.T) {
	claims.Verifies(t, claims.CH_042, "command-failure-named")
	c := build(t, nil)
	c.Python.Test = "python3 -c 'import sys; print(\"no such suite\"); sys.exit(4)'"
	_, err := Run(c)
	if err == nil || !strings.Contains(err.Error(), "the Python tests") || !strings.Contains(err.Error(), "no such suite") {
		t.Errorf("err = %v", err)
	}
}
