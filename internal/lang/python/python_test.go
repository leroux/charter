package python

import (
	"bufio"
	"encoding/json"
	"errors"
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

func read(t *testing.T, files map[string]string) *facts.Set {
	t.Helper()
	root := t.TempDir()
	write(t, root, files)
	var rels []string
	for rel := range files {
		if strings.HasSuffix(rel, ".py") {
			rels = append(rels, rel)
		}
	}
	slices.Sort(rels)
	set, err := Reader{Root: root, Dir: "dom0", EnvFile: "dom0/torii/env.py", ClaimsDir: "claims", Table: "claims/charter.toml:3"}.Read(rels)
	if err != nil {
		t.Fatal(err)
	}
	return set
}

func messages(set *facts.Set) string {
	var ms []string
	for _, r := range set.Refused {
		ms = append(ms, r.Pos()+": "+r.Rule+": "+r.Message)
	}
	return strings.Join(ms, "\n")
}

func TestAPythonDirectiveAnchorsTheDefBelowIt(t *testing.T) {
	claims.Verifies(t, claims.CH_021, "python-anchor-recorded", "python-comment-elsewhere-refused")
	set := read(t, map[string]string{"dom0/torii/plan.py": `import functools

# Derives the policy from the machine file.
# charter:realizes SW-010, SW-012
@functools.cache
def derive(machine):
    # charter:realizes SW-010
    return machine


# charter:concerns SW-010
class Planner:
    # charter:realizes SW-010
    def step(self):
        pass

# charter:realizes SW-010

x = 1  # charter:realizes SW-010
`})
	if len(set.Anchors) != 3 {
		t.Fatalf("anchors = %+v\nrefused:\n%s", set.Anchors, messages(set))
	}
	if a := set.Anchors[0]; a.Pos() != "dom0/torii/plan.py:4" || a.Kind != "realizes" || !slices.Equal(a.IDs, []string{"SW-010", "SW-012"}) || a.Decl != "def derive" {
		t.Errorf("above the decorator: %+v", a)
	}
	if a := set.Anchors[1]; a.Kind != "concerns" || a.Decl != "class Planner" || a.DeclKind != "class" {
		t.Errorf("above a class: %+v", a)
	}
	if a := set.Anchors[2]; a.Decl != "def Planner.step" {
		t.Errorf("above a method: %+v", a)
	}
	got := messages(set)
	for _, line := range []string{"dom0/torii/plan.py:7:", "dom0/torii/plan.py:17:", "dom0/torii/plan.py:19:"} {
		if !strings.Contains(got, line+" rule 24: # charter:realizes is written on the comment lines directly above the def or class") {
			t.Errorf("want %s refused in:\n%s", line, got)
		}
	}
}

const unittestFile = `import unittest

import claims


class PlanTest(unittest.TestCase):
    def test_offline(self):
        """An offline qube gets no netvm."""
        claims.verifies(self, claims.SW_010, "waited-request-passes", "burst-caps-tokens")
        _ = claims.SW_011

    def test_later(self):
        x = 1
        claims.verifies(self, claims.SW_010)

    def helper(self):
        claims.verifies(self, claims.SW_010)

    def test_wrong_first(self):
        claims.verifies(None, claims.SW_010)

    def test_by_variable(self):
        claims.verifies(self, spec)


def test_pytest_style():
    claims.verifies(None, claims.SW_010, "burst-caps-tokens")


def test_handle_variable():
    claims.verifies(None, claims.SW_010, handle)
`

func TestAPythonCitationOpensAUnittestOrPytestTest(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "python-citation-read", "python-later-statement-refused")
	set := read(t, map[string]string{"dom0/tests/test_plan.py": unittestFile})
	if len(set.Citations) != 3 {
		t.Fatalf("citations = %+v\n%s", set.Citations, messages(set))
	}
	if c := set.Citations[0]; c.Pos() != "dom0/tests/test_plan.py:9" || c.Symbol != "SW_010" || c.Test != "test_offline" || !slices.Equal(c.Handles, []string{"waited-request-passes", "burst-caps-tokens"}) {
		t.Errorf("unittest citation = %+v", c)
	}
	if c := set.Citations[1]; c.Test != "test_pytest_style" || !slices.Equal(c.Handles, []string{"burst-caps-tokens"}) {
		t.Errorf("pytest citation = %+v", c)
	}
	if len(set.References) != 1 || set.References[0].Symbol != "SW_011" {
		t.Errorf("references = %+v", set.References)
	}
	got := messages(set)
	for _, want := range []string{
		"test_plan.py:14: rule 24: verifies is the leading statement of its test",
		"test_plan.py:17: rule 24: verifies sits outside a test",
		"test_plan.py:20: rule 24: a unittest test passes self as the first argument of verifies",
		"test_plan.py:23: rule 24: verifies names a spec by its symbol",
		"test_plan.py:31: rule 24: a scenario handle is a string literal",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q in:\n%s", want, got)
		}
	}
	if !slices.Equal(set.TestFiles, []string{"dom0/tests/test_plan.py"}) {
		t.Errorf("test files = %v", set.TestFiles)
	}
}

func TestPythonEnvironmentReadsAreFacts(t *testing.T) {
	claims.Verifies(t, claims.CH_034, "python-read-checked")
	set := read(t, map[string]string{
		"dom0/torii/env.py": "BUNDLE = \"TORII_BUNDLE\"\n_private = \"TORII_X\"\n",
		"dom0/torii/apply.py": "import os\nfrom torii import env\nfrom torii.env import BUNDLE\n\n" +
			"a = os.environ.get(env.BUNDLE)\nb = os.getenv(BUNDLE)\nc = os.environ.get(\"TORII_BUNDLE\")\nd = os.environ[name]\ne = os.getenv(name)  # charter:dynamic\n",
	})
	if len(set.EnvConsts) != 1 || set.EnvConsts[0].Name != "env.BUNDLE" || set.EnvConsts[0].Value != "TORII_BUNDLE" {
		t.Errorf("constants = %+v", set.EnvConsts)
	}
	if len(set.EnvReads) != 5 {
		t.Fatalf("reads = %+v", set.EnvReads)
	}
	r := set.EnvReads
	if r[0].Const != "env.BUNDLE" || r[1].Const != "env.BUNDLE" || r[2].Literal != "TORII_BUNDLE" || r[3].Const != "" || r[3].Literal != "" || r[3].Dynamic || !r[4].Dynamic {
		t.Errorf("reads = %+v", r)
	}
	if len(set.Literals) != 1 || set.Literals[0].Value != "TORII_BUNDLE" || set.Literals[0].Line != 7 {
		t.Errorf("literals = %+v", set.Literals)
	}
}

func TestAMissingPython3NamesTheTable(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_PQ7X, "parser-missing-named")
	old := Interpreter
	Interpreter = "python3-not-installed-here"
	defer func() { Interpreter = old }()
	root := t.TempDir()
	write(t, root, map[string]string{"a.py": "x = 1\n"})
	_, err := Reader{Root: root, Table: "claims/charter.toml:3"}.Read([]string{"a.py"})
	if !errors.Is(err, ErrNoInterpreter) || !strings.HasPrefix(err.Error(), "claims/charter.toml:3: ") {
		t.Errorf("err = %v", err)
	}
}

const fixtureArea = "# Refill\n\n" +
	"## SW-010 The bucket refills at the rate\nStatus: active\nKind: invariant\n\nAt each request the bucket reads the clock.\n\nWhy: a caller that waits is let through.\n\n" +
	"Scenarios:\n1. `waited-request-passes`: a request after a wait passes.\n2. `burst-caps-tokens`: a long wait earns no more than the burst.\n\n" +
	"## SW-011 A \"planned\" rule\nStatus: planned\n\nNot yet built.\n\nWhy: someone needs it.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\n\nNo longer.\n\nWhy: it was replaced.\n"

func generated(t *testing.T, root string) []byte {
	t.Helper()
	write(t, root, map[string]string{"claims/refill.md": fixtureArea, "claims/observations.md": "## LockIsAdvisory\nSubject: os\nStatus: documented\nSource: man flock\n\nflock locks are advisory.\n"})
	rec, err := record.Read(root, "claims")
	if err != nil || len(rec.Problems) != 0 {
		t.Fatalf("%v %q", err, rec.Problems)
	}
	src := Generate(rec, "charter gen")
	write(t, root, map[string]string{"claims/claims.py": string(src)})
	return src
}

func TestTheGeneratedModuleEqualsTheFixture(t *testing.T) {
	claims.Verifies(t, claims.CH_020, "python-module-rendered")
	got := generated(t, t.TempDir())
	golden := filepath.Join("testdata", "claims.py")
	if os.Getenv("CHARTER_UPDATE_GOLDEN") != "" {
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("claims.py differs from %s; regenerate it with CHARTER_UPDATE_GOLDEN=1 if the change is meant", golden)
	}
	for _, line := range strings.Split(string(got), "\n") {
		if strings.HasPrefix(line, "import ") && !slices.Contains([]string{"import json", "import os", "import sys", "import threading"}, line) {
			t.Errorf("claims.py imports beyond the standard library: %s", line)
		}
	}
	for _, want := range []string{"SW_010 = Spec(\"SW-010\", \"The bucket refills at the rate\", \"active\", \"invariant\")", "SW_011 = Spec(\"SW-011\", \"A \\\"planned\\\" rule\", \"planned\", \"\")", "#: Deprecated: SW-012 is retired", "LockIsAdvisory = Observation(\"LockIsAdvisory\", \"os\", \"documented\")"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("claims.py lacks %s", want)
		}
	}
	if strings.Contains(string(got), "At each request") {
		t.Error("claims.py holds an entry's text")
	}
}

type line struct{ ID, Handle, Test, File, Outcome string }

func events(t *testing.T, p string) []line {
	t.Helper()
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []line
	s := bufio.NewScanner(f)
	for s.Scan() {
		var l line
		if err := json.Unmarshal(s.Bytes(), &l); err != nil {
			t.Fatalf("not a JSON line: %q", s.Text())
		}
		out = append(out, l)
	}
	return out
}

func TestAPythonTestRecordsItsOutcomeInTheEventsFile(t *testing.T) {
	claims.Verifies(t, claims.CH_024, "event-per-citation", "outcome-recorded", "unset-records-nothing")
	root := t.TempDir()
	generated(t, root)
	write(t, root, map[string]string{
		"dom0/tests/test_unit.py": "import unittest\n\nimport claims\n\n\nclass T(unittest.TestCase):\n" +
			"    def test_passes(self):\n        claims.verifies(self, claims.SW_010, \"waited-request-passes\")\n\n" +
			"    def test_fails(self):\n        claims.verifies(self, claims.SW_010, \"burst-caps-tokens\")\n        self.fail(\"on purpose\")\n\n" +
			"    def test_skips(self):\n        claims.verifies(self, claims.SW_011)\n        self.skipTest(\"on purpose\")\n",
		"dom0/tests/test_py.py":  "import pytest\n\nimport claims\n\n\ndef test_passes():\n    claims.verifies(None, claims.SW_010, \"waited-request-passes\")\n\n\ndef test_fails():\n    claims.verifies(None, claims.SW_010, \"burst-caps-tokens\")\n    assert False\n\n\ndef test_skips():\n    claims.verifies(None, claims.SW_011)\n    pytest.skip(\"on purpose\")\n",
		"dom0/tests/conftest.py": "pytest_plugins = [\"claims\"]\n",
	})
	env := []string{"PYTHONPATH=" + filepath.Join(root, "claims")}
	unitEvents := filepath.Join(t.TempDir(), "unit.jsonl")
	run(t, filepath.Join(root, "dom0"), append(env, "CHARTER_EVENTS="+unitEvents), "python3", "-m", "unittest", "tests.test_unit")
	want := map[string]string{"test_passes/waited-request-passes": "passed", "test_fails/burst-caps-tokens": "failed", "test_skips/": "skipped"}
	got := map[string]string{}
	for _, l := range events(t, unitEvents) {
		got[l.Test[strings.LastIndex(l.Test, ".")+1:]+"/"+l.Handle] = l.Outcome
		if !strings.HasSuffix(l.File, "test_unit.py") || (l.ID != "SW-010" && l.ID != "SW-011") {
			t.Errorf("line %+v", l)
		}
	}
	if !mapsEqual(got, want) {
		t.Errorf("unittest outcomes = %v, want %v", got, want)
	}
	if _, err := exec.LookPath("pytest"); err == nil {
		pyEvents := filepath.Join(t.TempDir(), "py.jsonl")
		run(t, filepath.Join(root, "dom0"), append(env, "CHARTER_EVENTS="+pyEvents), "python3", "-m", "pytest", "-q", "-p", "no:cacheprovider", "tests/test_py.py")
		got = map[string]string{}
		for _, l := range events(t, pyEvents) {
			got[l.Test[strings.LastIndex(l.Test, "::")+2:]+"/"+l.Handle] = l.Outcome
		}
		if !mapsEqual(got, want) {
			t.Errorf("pytest outcomes = %v, want %v", got, want)
		}
	}
	run(t, filepath.Join(root, "dom0"), env, "python3", "-m", "unittest", "tests.test_unit")
	if matches, _ := filepath.Glob(filepath.Join(root, "**", "*.jsonl")); len(matches) != 0 {
		t.Errorf("with CHARTER_EVENTS unset a file was written: %v", matches)
	}
}

func run(t *testing.T, dir string, env []string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
