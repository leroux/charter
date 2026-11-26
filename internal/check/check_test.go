package check

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/config"
)

// fixture is a clean charter: a need, two decisions, an observation, the
// code that realizes them, the tests and the probe that cite them.
func fixture() map[string]string {
	return map[string]string{
		"go.mod": "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"claims/charter.toml": "name = \"fx\"\nstakeholder-tests = [\"cmd/*_test.go\"]\nprobes = [\"probes\"]\nenv-patterns = [\"^FX_[A-Z0-9_]+$\"]\n\n" +
			"[go]\nmodule = \"example.com/fx/claims\"\ndir = \".\"\ntest = \"go test ./...\"\nenv = \"env\"\n\n" +
			"[[surface]]\nkind = \"op\"\ncommand = \"echo refill\"\n",
		"claims/needs.md": "# Needs\n\n## STK-001 A caller that waits is let through again\nStatus: active\n\nA caller that waited is let through.\n\nWhy: a caller should not be punished for waiting.\n\nScenarios:\n1. `waits-and-passes`: a caller that waited passes.\n",
		"claims/refill.md": "# Refill\nCode: pkg/a\n\nHow the bucket refills.\n\n" +
			"## SW-010 The bucket refills at the rate\nStatus: active\nKind: invariant\nServes: STK-001\nAccounts: env:FX_RATE, op:refill\n\nAt each request the bucket reads the clock.\n\nWhy: a caller that waits is let through.\n\n" +
			"Scenarios:\n1. `waited-request-passes`: a request after a wait passes.\n2. `burst-caps-tokens`: a long wait earns no more than the burst.\n\n" +
			"## SW-011 A planned rule\nStatus: planned\nServes: STK-001\n\nNot yet built.\n\nWhy: someone needs it.\n",
		"claims/observations.md": "## ClockIsMonotonic\nSubject: runtime\nStatus: observed\nObserved: 2026-11-01\nSupports: SW-010\n\nThe clock never runs backwards.\n",
		"env/env.go":             "// Package env names the variables.\npackage env\n\n// Rate is the rate.\nconst Rate = \"FX_RATE\"\n",
		"pkg/a/a.go":             "package a\n\nimport (\n\t\"os\"\n\n\t\"example.com/fx/env\"\n)\n\n// Refill adds the tokens earned.\n//\n//charter:realizes SW-010\nfunc Refill() string { return os.Getenv(env.Rate) }\n",
		"pkg/a/a_test.go":        "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestRefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"waited-request-passes\", \"burst-caps-tokens\")\n}\n",
		"cmd/fx_test.go":         "package cmd\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestACallerWaits(t *testing.T) {\n\tclaims.Verifies(t, claims.STK_001, \"waits-and-passes\")\n}\n",
		"probes/clock_test.go":   "package probes\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestClock(t *testing.T) {\n\tclaims.Observes(t, claims.ClockIsMonotonic)\n}\n",
	}
}

// build writes the files into a fresh root, generates the claims module
// and returns the charter.
func build(t *testing.T, files map[string]string) *config.Config {
	t.Helper()
	root := t.TempDir()
	for rel, src := range files {
		if src == "" {
			continue
		}
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
	if err := Write(c); err != nil {
		t.Fatal(err)
	}
	return c
}

func problems(t *testing.T, c *config.Config) []string {
	t.Helper()
	ps, err := Problems(c)
	if err != nil {
		t.Fatal(err)
	}
	return ps
}

// with is the fixture with the edits applied: a file's new content, or ""
// to remove it.
func with(edits map[string]string) map[string]string {
	f := fixture()
	for k, v := range edits {
		f[k] = v
	}
	return f
}

// reports checks that the tree's problems hold each wanted one.
func reports(t *testing.T, files map[string]string, want ...string) []string {
	t.Helper()
	ps := problems(t, build(t, files))
	for _, w := range want {
		found := false
		for _, p := range ps {
			found = found || strings.Contains(p, w)
		}
		if !found {
			t.Errorf("want a problem holding\n  %s\ngot:\n  %s", w, strings.Join(ps, "\n  "))
		}
	}
	return ps
}

func replace(t *testing.T, file, old, new string) map[string]string {
	t.Helper()
	f := fixture()
	if strings.Count(f[file], old) != 1 {
		t.Fatalf("%s holds %q %d times", file, old, strings.Count(f[file], old))
	}
	f[file] = strings.Replace(f[file], old, new, 1)
	return f
}

func TestTheCleanFixturePasses(t *testing.T) {
	chartertest.Verifies(t, claims.CH_037, "clean-fixture-passes")
	chartertest.Verifies(t, claims.CH_001, "clean-tree-passes")
	chartertest.Verifies(t, claims.CH_031, "need-without-anchor-passes", "onside-citation-counts")
	if ps := problems(t, build(t, fixture())); len(ps) != 0 {
		t.Errorf("the clean fixture:\n%s", strings.Join(ps, "\n"))
	}
}

func TestAMissingScenarioIsNamed(t *testing.T) {
	chartertest.Verifies(t, claims.CH_030, "scenario-unverified")
	chartertest.Verifies(t, claims.CH_001, "missing-scenario-named")
	reports(t, replace(t, "pkg/a/a_test.go", ", \"burst-caps-tokens\")", ")"),
		"claims/refill.md:18: rule 25: scenario burst-caps-tokens of SW-010 is verified by no test; cite it with verifies(claims.SW_010, \"burst-caps-tokens\") as a test's first statement, or waive SW-010/burst-caps-tokens in claims/waivers.txt with a reason")
}

func TestAnActiveDecisionNeedsARealizesAnchor(t *testing.T) {
	chartertest.Verifies(t, claims.CH_030, "not-realized")
	reports(t, replace(t, "pkg/a/a.go", "//charter:realizes SW-010", "//charter:concerns SW-010"),
		"claims/refill.md:6: rule 25: active entry SW-010 is realized by no code; add charter:realizes SW-010 where it is implemented")
	reports(t, replace(t, "pkg/a/a.go", "//\n//charter:realizes SW-010\n", ""),
		"active entry SW-010 is realized by no code")
}

func TestAConcernsAnchorRealizesNothing(t *testing.T) {
	chartertest.Verifies(t, claims.CH_022, "concerns-is-not-realized")
	reports(t, replace(t, "pkg/a/a.go", "//charter:realizes SW-010", "//charter:concerns SW-010"), "active entry SW-010 is realized by no code")
}

func TestAnAnchorNamesAnEntryThatExistsAndIsNotPlanned(t *testing.T) {
	chartertest.Verifies(t, claims.CH_022, "unknown-id-refused", "planned-refused")
	reports(t, replace(t, "pkg/a/a.go", "//charter:realizes SW-010", "//charter:realizes SW-010, SW-999, SW-011"),
		"pkg/a/a.go:11: rule 24: charter:realizes names SW-999, which is not an entry of the record",
		"pkg/a/a.go:11: rule 24: charter:realizes names planned entry SW-011; set it active as its work begins")
}

func TestAnEntryWithoutScenariosIsCitedWhole(t *testing.T) {
	chartertest.Verifies(t, claims.CH_030, "spec-unverified")
	f := replace(t, "claims/refill.md", "Status: planned\nServes: STK-001\n", "Status: active\nServes: STK-001\n")
	f["pkg/a/a.go"] = strings.Replace(f["pkg/a/a.go"], "SW-010", "SW-010, SW-011", 1)
	reports(t, f, "claims/refill.md:20: rule 25: active entry SW-011 is verified by no test")
}

func TestAWaiverAbsorbsAScenarioAndNeverAnAnchor(t *testing.T) {
	chartertest.Verifies(t, claims.CH_030, "waiver-absorbs-a-scenario", "waiver-never-absorbs-an-anchor", "idle-waiver-fails", "glob-refused")
	f := replace(t, "pkg/a/a_test.go", ", \"burst-caps-tokens\")", ")")
	f["claims/waivers.txt"] = "SW-010/burst-caps-tokens\tneeds a clock we cannot drive yet\n"
	if ps := problems(t, build(t, f)); len(ps) != 0 {
		t.Errorf("a waived scenario: %q", ps)
	}
	f["pkg/a/a.go"] = strings.Replace(f["pkg/a/a.go"], "//charter:realizes SW-010", "//charter:concerns SW-010", 1)
	reports(t, f, "active entry SW-010 is realized by no code")
	g := fixture()
	g["claims/waivers.txt"] = "SW-010/burst-caps-tokens\tidle\nSW-010/waited-request-passes\nSW-010/*\tall of it\n"
	reports(t, g,
		"claims/waivers.txt:1: rule 25: the waiver for SW-010/burst-caps-tokens absorbs no active entry or scenario that lacks a test",
		"claims/waivers.txt:2: rule 25: the waiver for SW-010/waited-request-passes gives no reason after a tab",
		"claims/waivers.txt:3: rule 25: the waiver \"SW-010/*\" is a glob")
}

func TestAHandleTheEntryDoesNotListIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.CH_023, "unknown-handle-refused")
	reports(t, replace(t, "pkg/a/a_test.go", "\"burst-caps-tokens\")", "\"burst-caps-tokens\", \"no-such\")"),
		"pkg/a/a_test.go:10: rule 24: SW-010 has no scenario \"no-such\"")
}

func TestANeedIsVerifiedFromTheStakeholdersSide(t *testing.T) {
	chartertest.Verifies(t, claims.CH_031, "offside-citation-refused")
	f := fixture()
	f["pkg/a/need_test.go"] = strings.Replace(f["cmd/fx_test.go"], "package cmd", "package a", 1)
	f["cmd/fx_test.go"] = ""
	reports(t, f,
		"pkg/a/need_test.go:10: rule 25: need STK-001 is verified from pkg/a/need_test.go, which is not on the stakeholder's side",
		"scenario waits-and-passes of STK-001 is verified by no test")
}

func TestProbesAndTestsStayApart(t *testing.T) {
	chartertest.Verifies(t, claims.CH_032, "probe-reference-is-no-citation", "test-citing-an-observation-fails", "probe-by-path")
	chartertest.Verifies(t, claims.CH_005, "probes-and-tests-stay-apart")
	f := fixture()
	f["probes/clock_test.go"] = strings.Replace(f["probes/clock_test.go"], "claims.Observes(t, claims.ClockIsMonotonic)", "claims.Observes(t, claims.ClockIsMonotonic)\n\t_ = claims.SW_010", 1)
	f["probes/refill_test.go"] = strings.Replace(f["pkg/a/a_test.go"], "package a", "package probes", 1)
	f["pkg/a/obs_test.go"] = "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestObs(t *testing.T) {\n\t_ = claims.ClockIsMonotonic\n}\n"
	ps := reports(t, f,
		"probes/refill_test.go:10: rule 12: probe probes/refill_test.go cites entry SW-010; a probe cites observations only",
		"probes/clock_test.go:11: rule 12: probe probes/clock_test.go references entry SW-010",
		"pkg/a/obs_test.go:10: rule 12: observation ClockIsMonotonic is referenced from pkg/a/obs_test.go, which is not a probe")
	_ = ps
}

func TestObservationsNeedAProbePath(t *testing.T) {
	chartertest.Verifies(t, claims.CH_032, "probe-path-required")
	reports(t, replace(t, "claims/charter.toml", "probes = [\"probes\"]\n", ""),
		"claims/charter.toml:1: rule 11: the record holds observations and charter.toml names no probe path")
}

func TestObservationsCarryTheirStatussFields(t *testing.T) {
	chartertest.Verifies(t, claims.CH_035, "observed-needs-date-and-probe", "documented-needs-source", "subject-required", "broken-under-active")
	chartertest.Verifies(t, claims.CH_005, "observed-fact-has-a-date-and-a-probe", "broken-fact-under-active-fails")
	reports(t, replace(t, "claims/observations.md", "Observed: 2026-11-01\n", "Observed: 2026-13-01\n"), "claims/observations.md:1: rule 11: ClockIsMonotonic has Observed \"2026-13-01\"; want YYYY-MM-DD")
	reports(t, replace(t, "claims/observations.md", "Observed: 2026-11-01\n", ""), "rule 11: ClockIsMonotonic is observed and has no Observed date")
	reports(t, with(map[string]string{"probes/clock_test.go": ""}), "rule 11: observed observation ClockIsMonotonic is cited by no probe")
	reports(t, replace(t, "claims/observations.md", "Status: observed\nObserved: 2026-11-01\n", "Status: documented\n"), "rule 11: ClockIsMonotonic is documented and has no Source")
	reports(t, replace(t, "claims/observations.md", "Status: observed\nObserved: 2026-11-01\n", "Status: documented\nSource: man clock_gettime\nObserved: 2026-11-01\n"), "rule 11: ClockIsMonotonic is documented and carries an Observed date")
	reports(t, replace(t, "claims/observations.md", "Subject: runtime\n", ""), "rule 11: ClockIsMonotonic has no Subject")
	reports(t, replace(t, "claims/observations.md", "Status: observed\n", "Status: broken\n"), "claims/observations.md:1: rule 13: broken observation ClockIsMonotonic supports active entry SW-010")
	f := replace(t, "claims/observations.md", "Status: observed\n", "Status: broken\n")
	f["claims/observations.md"] = strings.Replace(f["claims/observations.md"], "Supports: SW-010", "Supports: SW-011", 1)
	for _, p := range problems(t, build(t, f)) {
		if strings.Contains(p, "rule 13") {
			t.Errorf("a broken observation under a planned entry: %s", p)
		}
	}
}

func TestRelationsNameEntriesTheRecordHolds(t *testing.T) {
	chartertest.Verifies(t, claims.CH_015, "unknown-target", "self-relation", "observation-supports-existing")
	reports(t, replace(t, "claims/refill.md", "Serves: STK-001\nAccounts", "Serves: STK-001, SW-999\nSee: SW-010\nAccounts"),
		"claims/refill.md:6: rule 23: SW-010 serves SW-999, which is not an entry",
		"claims/refill.md:6: rule 23: SW-010 sees itself")
	reports(t, replace(t, "claims/observations.md", "Supports: SW-010", "Supports: SW-010, SW-404"),
		"claims/observations.md:1: rule 23: observation ClockIsMonotonic supports SW-404, which is not an entry")
}

func TestEverySurfaceIsAccountedForOrBaselined(t *testing.T) {
	chartertest.Verifies(t, claims.CH_033, "surfaces-enumerated", "unaccounted-fails", "accounts-must-exist", "baseline-holds-live-gaps-once", "deprecated-accounts-nothing")
	chartertest.Verifies(t, claims.CH_004, "unaccounted-surface-fails", "baseline-holds-only-live-gaps")
	reports(t, replace(t, "claims/charter.toml", "command = \"echo refill\"", "command = \"printf 'refill\\\\ndrain\\\\n'\""),
		"rule 6: op:drain is neither accounted for by an entry nor listed in claims/baseline.txt")
	reports(t, replace(t, "claims/refill.md", "Accounts: env:FX_RATE, op:refill", "Accounts: env:FX_RATE, op:refill, op:gone"),
		"claims/refill.md:6: rule 7: SW-010 accounts for op:gone, which no enumerator found")
	f := fixture()
	f["claims/baseline.txt"] = "op:refill\taccounted now\nop:gone\tno longer there\nop:refill\ttwice\nenv:FX_RATE\n"
	reports(t, f,
		"claims/baseline.txt:1: rule 7: baseline lists op:refill, which SW-010 now accounts for",
		"claims/baseline.txt:2: rule 7: baseline lists op:gone, which no enumerator found",
		"claims/baseline.txt:3: rule 7: baseline lists op:refill twice",
		"claims/baseline.txt:4: rule 7: \"env:FX_RATE\" has no reason")
	g := replace(t, "claims/refill.md", "Status: planned\nServes: STK-001\n", "Status: deprecated\nServes: STK-001\nAccounts: op:refill\n")
	g["claims/refill.md"] = strings.Replace(g["claims/refill.md"], "Accounts: env:FX_RATE, op:refill\n", "Accounts: env:FX_RATE\n", 1)
	reports(t, g, "rule 6: op:refill is neither accounted for by an entry nor listed in claims/baseline.txt")
}

func TestTheUpdateVariableRebuildsTheBaseline(t *testing.T) {
	chartertest.Verifies(t, claims.CH_033, "update-rebuilds")
	f := replace(t, "claims/charter.toml", "command = \"echo refill\"", "command = \"printf 'refill\\\\ndrain\\\\nfill\\\\n'\"")
	f["claims/baseline.txt"] = "op:fill\tkept reason\n"
	c := build(t, f)
	t.Setenv("CHARTER_UPDATE_BASELINE", "1")
	if ps := problems(t, c); len(ps) != 0 {
		t.Errorf("with the update variable: %q", ps)
	}
	data, err := os.ReadFile(filepath.Join(c.Root, "claims", "baseline.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "op:drain\tunclaimed at baseline init\nop:fill\tkept reason\n") || strings.Contains(string(data), "op:refill") {
		t.Errorf("rebuilt baseline:\n%s", data)
	}
	first := string(data)
	problems(t, c)
	again, _ := os.ReadFile(filepath.Join(c.Root, "claims", "baseline.txt"))
	if string(again) != first {
		t.Errorf("a second run changed the baseline:\n%s", again)
	}
}

func TestAnEmptyConstantsModuleIsAConfigurationProblem(t *testing.T) {
	chartertest.Verifies(t, claims.CH_033, "empty-enumerator")
	f := replace(t, "env/env.go", "// Rate is the rate.\nconst Rate = \"FX_RATE\"\n", "")
	f["pkg/a/a.go"] = strings.Replace(strings.Replace(f["pkg/a/a.go"], "os.Getenv(env.Rate)", "\"\"", 1), "\t\"os\"\n\n\t\"example.com/fx/env\"\n", "", 1)
	f["claims/refill.md"] = strings.Replace(f["claims/refill.md"], "env:FX_RATE, ", "", 1)
	reports(t, f, "claims/charter.toml:10: rule 9: the constants module env declares no variable")
}

func TestAnEnvironmentReadGoesThroughAConstant(t *testing.T) {
	chartertest.Verifies(t, claims.CH_034, "literal-argument", "variable-argument", "duplicated-name", "undeclared-pattern", "dynamic-directive")
	chartertest.Verifies(t, claims.CH_004, "env-read-through-a-constant")
	reports(t, replace(t, "pkg/a/a.go", "os.Getenv(env.Rate)", "os.Getenv(\"FX_RATE\")"),
		"pkg/a/a.go:12: rule 8: an environment read names \"FX_RATE\" by a literal; read it through env.Rate",
		"pkg/a/a.go:12: rule 8: literal \"FX_RATE\" duplicates env.Rate; use the constant")
	reports(t, replace(t, "pkg/a/a.go", "os.Getenv(env.Rate)", "os.Getenv(env.Rate + \"\")"),
		"pkg/a/a.go:12: rule 8: an environment read must name a constant of a constants module")
	reports(t, replace(t, "pkg/a/a.go", "return os.Getenv(env.Rate) }", "return os.Getenv(env.Rate) + \"FX_BURST\" }"),
		"pkg/a/a.go:12: rule 8: literal \"FX_BURST\" looks like an environment variable; declare it in a constants module")
	f := replace(t, "pkg/a/a.go", "os.Getenv(env.Rate) }", "os.Getenv(env.Rate) + os.Getenv(name) } //charter:dynamic\n\nvar name = env.Rate\n")
	for _, p := range problems(t, build(t, f)) {
		if strings.Contains(p, "rule 8") {
			t.Errorf("a read under the dynamic directive: %s", p)
		}
	}
}

func TestAHandWrittenSourceInTheClaimsDirectoryIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.CH_025, "hand-written-source-refused")
	reports(t, with(map[string]string{"claims/facts.go": "package claims\n\nvar Extra = 1\n"}),
		"claims/facts.go:1: rule 1: claims/facts.go is written by hand in the claims directory, which holds the record and the generated files; entries and observations are written in markdown")
}

func TestAStaleGeneratedFileNamesCharterGen(t *testing.T) {
	chartertest.Verifies(t, claims.CH_020, "stale-generated-fails", "gen-writes")
	c := build(t, fixture())
	specs := filepath.Join(c.Root, "claims", "specs.go")
	before, _ := os.ReadFile(specs)
	if err := Write(c); err != nil {
		t.Fatal(err)
	}
	if after, _ := os.ReadFile(specs); string(after) != string(before) {
		t.Error("a second charter gen changed specs.go")
	}
	if err := os.WriteFile(specs, append(before, []byte("// edited\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(c.Root, "claims", "go.mod")); err != nil {
		t.Fatal(err)
	}
	ps := problems(t, c)
	if !strings.Contains(strings.Join(ps, "\n"), "claims/specs.go:1: rule 15: claims/specs.go is stale: it differs from what the record renders; run charter gen") ||
		!strings.Contains(strings.Join(ps, "\n"), "claims/go.mod:1: rule 15: claims/go.mod is missing; run charter gen") {
		t.Errorf("problems = %q", ps)
	}
	if after, _ := os.ReadFile(specs); !strings.HasSuffix(string(after), "// edited\n") {
		t.Error("charter check wrote a generated file")
	}
}

func TestEveryRuleHasAMutationThatTripsIt(t *testing.T) {
	chartertest.Verifies(t, claims.CH_001, "every-rule-has-a-mutation")
	tripped := map[string]bool{}
	for _, f := range []map[string]string{
		with(map[string]string{"claims/facts.go": "package claims\n"}),
		with(map[string]string{"pkg/a/dot_test.go": "package a\n\nimport . \"example.com/fx/claims\"\n\nvar _ = SW_010\n"}),    // rule 10                                                         // rule 1
		replace(t, "pkg/a/a.go", "os.Getenv(env.Rate)", "os.Getenv(\"FX_RATE\")"),                                              // rule 8
		replace(t, "claims/charter.toml", "command = \"echo refill\"", "command = \"exit 1\""),                                 // rule 9
		replace(t, "probes/clock_test.go", "claims.Observes(t, claims.ClockIsMonotonic)", "claims.Verifies(t, claims.SW_010)"), // rule 11, 12
		replace(t, "claims/observations.md", "Status: observed\n", "Status: broken\n"),                                         // rule 13
		replace(t, "claims/refill.md", "Serves: STK-001\nAccounts", "Serves: SW-999\nAccounts"),                                // rule 23
		replace(t, "pkg/a/a_test.go", "\"burst-caps-tokens\")", "\"no-such\")"),                                                // rule 24, 25
		replace(t, "claims/refill.md", "Accounts: env:FX_RATE, op:refill", "Accounts: env:FX_RATE"),                            // rule 6
		replace(t, "claims/refill.md", "Accounts: env:FX_RATE, op:refill", "Accounts: env:FX_RATE, op:refill, op:x"),           // rule 7
	} {
		for _, p := range problems(t, build(t, f)) {
			if i := strings.Index(p, "rule "); i >= 0 {
				n := strings.Fields(p[i+5:])[0]
				tripped["rule "+strings.TrimSuffix(n, ":")] = true
			}
		}
	}
	c := build(t, fixture())
	if err := os.WriteFile(filepath.Join(c.Root, "claims", "specs.go"), []byte("package claims\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, p := range problems(t, c) {
		if strings.Contains(p, "rule 15") {
			tripped["rule 15"] = true
		}
	}
	for _, r := range gateRules(t) {
		if !tripped[r] {
			t.Errorf("%s: no mutation of the clean fixture trips it", r)
		}
	}
}

// gateRules are the rule numbers the gate's source reports problems under.
func gateRules(t *testing.T) []string {
	t.Helper()
	rules := map[string]bool{}
	for _, dir := range []string{"../gate", ".", "../surface", "../lang/golang"} {
		pkgs, err := parser.ParseDir(token.NewFileSet(), dir, func(fi os.FileInfo) bool { return !strings.HasSuffix(fi.Name(), "_test.go") }, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range pkgs {
			for name := range p.Files {
				src, _ := os.ReadFile(name)
				for _, part := range strings.Split(string(src), "rule ")[1:] {
					n := ""
					for _, r := range part {
						if r < '0' || r > '9' {
							break
						}
						n += string(r)
					}
					if _, err := strconv.Atoi(n); err == nil {
						rules["rule "+n] = true
					}
				}
			}
		}
	}
	var out []string
	for r := range rules {
		out = append(out, r)
	}
	return out
}

func TestTheGatesRulesImportNoLanguage(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_PQ7X, "rules-name-no-language")
	pkgs, err := parser.ParseDir(token.NewFileSet(), "../gate", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pkgs {
		for name, f := range p.Files {
			for _, imp := range f.Imports {
				v, _ := strconv.Unquote(imp.Path.Value)
				if strings.Contains(v, "/internal/lang/") || strings.HasPrefix(v, "go/") {
					t.Errorf("%s imports %s; the gate's rules read facts alone", name, v)
				}
			}
		}
	}
}

func TestTheSameRuleSaysTheSameInEveryLanguage(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_PQ7X, "same-rule-every-language")
	compiler, _ := filepath.Abs("../../node_modules/typescript")
	want := "claims/refill.md:18: rule 25: scenario burst-caps-tokens of SW-010 is verified by no test; cite it with verifies(claims.SW_010, \"burst-caps-tokens\") as a test's first statement, or waive SW-010/burst-caps-tokens in claims/waivers.txt with a reason"
	for _, tc := range []struct {
		lang  string
		edits map[string]string
	}{
		{"go", map[string]string{"pkg/a/a_test.go": strings.Replace(fixture()["pkg/a/a_test.go"], ", \"burst-caps-tokens\")", ")", 1)}},
		{"python", map[string]string{
			"pkg/a/a_test.go":        strings.Replace(fixture()["pkg/a/a_test.go"], ", \"burst-caps-tokens\")", ")", 1),
			"claims/charter.toml":    fixture()["claims/charter.toml"] + "\n[python]\ndir = \"py\"\ntest = \"python3 -m unittest\"\n",
			"py/tests/test_burst.py": "import unittest\n\nimport claims\n\n\nclass T(unittest.TestCase):\n    def test_burst(self):\n        claims.verifies(self, claims.SW_010)\n",
		}},
		{"typescript", map[string]string{
			"pkg/a/a_test.go":     strings.Replace(fixture()["pkg/a/a_test.go"], ", \"burst-caps-tokens\")", ")", 1),
			"claims/charter.toml": fixture()["claims/charter.toml"] + "\n[typescript]\ndir = \"web\"\ntest = \"vitest run\"\ncompiler = \"" + compiler + "\"\nimport = \"@fx/claims\"\n",
			"web/burst.test.ts":   "import { claims, verifies } from \"@fx/claims\";\n\ntest(\"burst\", () => {\n  verifies(claims.SW_010);\n});\n",
		}},
	} {
		ps := problems(t, build(t, with(tc.edits)))
		if len(ps) != 1 || ps[0] != want {
			t.Errorf("%s: problems =\n%s\nwant exactly\n%s", tc.lang, strings.Join(ps, "\n"), want)
		}
	}
}
