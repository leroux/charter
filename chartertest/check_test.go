package chartertest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestCheckPassesOnTheCleanFixture(t *testing.T) {
	Verifies(t, claims.CH_037, "clean-fixture-passes")
	Verifies(t, claims.CH_001, "clean-tree-passes")
	root := copyFixture(t)
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("clean fixture reported:\n%s", strings.Join(problems, "\n"))
	}
}

// ruleCase is a mutation of the clean fixture and a problem it trips.
type ruleCase struct {
	name string
	edit edit
	want string
}

// ruleCases trip every rule of the gate, each named after its rule.
var ruleCases = []ruleCase{
	{"rule 1 function",
		edit{rel: "claims/mutate.go", write: "package claims\n\nfunc init() {}\n"},
		"rule 1: functions are not allowed"},
	{"rule 1 assignment from a test",
		edit{rel: "pkg/b/b_test.go", write: "package b\n\nimport \"example.com/fixture/claims\"\n\nfunc init() { claims.HealthIsJSON.Text = \"x\" }\n"},
		"rule 1: HealthIsJSON is assigned outside its declaration"},
	{"rule 1 a spec written in Go",
		edit{rel: "claims/hand.go", write: "package claims\n\nimport \"github.com/leroux/charter\"\n\nvar Hand = &charter.Spec{ID: \"SW-009\"}\n"},
		"rule 1: Hand is not &charter.Observation{}; a spec is written in a spec file"},
	{"rule 1 non-literal field",
		edit{rel: "claims/facts.go", old: `Subject:  "os",`, repl: `Subject:  subject,`},
		"rule 1: LockIsAdvisory.Subject must be a string literal"},
	{"rule 1 observation field",
		edit{rel: "claims/facts.go", old: `Observed: "2026-09-24",`, repl: `Observed: "2026-09-24", Accounts: nil,`},
		"rule 1: HealthIsJSON.Accounts is not a field of an observation"},
	{"rule 1 observation status",
		edit{rel: "claims/facts.go", old: `Status:   charter.Observed,`, repl: `Status:   charter.Active,`},
		"rule 1: HealthIsJSON.Status must be charter.Observed, charter.Documented, charter.Broken or charter.ObservationUnsure"},
	{"rule 2 observation id",
		edit{rel: "claims/facts.go", old: `ID:       "HealthIsJSON",`, repl: `ID:       "Health",`},
		"rule 2: HealthIsJSON has ID"},
	{"rule 6",
		edit{rel: "pkg/b/wasm.go", write: "//go:build wasm\n\npackage b\n"},
		"rule 6: tag:wasm"},
	{"rule 7 baseline gone",
		edit{rel: "claims/baseline.txt", add: "tag:nonexistent\tgone\n"},
		"rule 7: baseline lists tag:nonexistent, which no enumerator found"},
	{"rule 7 baseline accounted",
		edit{rel: "claims/baseline.txt", add: "env:FX_TIMEOUT\tstill listed\n"},
		"rule 7: baseline lists env:FX_TIMEOUT, which CON-002 now accounts for"},
	{"rule 7 accounts",
		edit{rel: "claims/specs/SW-001-bucket-survives-restart.md", old: "event:bucket.upsert", repl: "event:bucket.gone"},
		"rule 7: SW-001 accounts for event:bucket.gone"},
	{"rule 7 baseline malformed",
		edit{rel: "claims/baseline.txt", add: "tag:desktop no tab here\n"},
		"rule 7: \"tag:desktop no tab here\" has no reason"},
	{"rule 7 baseline duplicate",
		edit{rel: "claims/baseline.txt", add: "tag:desktop\tagain\n"},
		"rule 7: baseline lists tag:desktop twice"},
	{"rule 8 literal argument",
		edit{rel: "pkg/a/a.go", old: "os.Getenv(env.DataDir)", repl: `os.Getenv("FX_DATA_DIR")`},
		"rule 8: os.Getenv must name a constant"},
	{"rule 8 variable argument",
		edit{rel: "pkg/a/read.go", write: "package a\n\nimport \"os\"\n\n// Read reads any variable.\nfunc Read(name string) string { return os.Getenv(name) }\n"},
		"rule 8: os.Getenv must name a constant"},
	{"rule 8 read inside another call",
		edit{rel: "pkg/a/read.go", write: "package a\n\nimport (\n\t\"os\"\n\t\"strings\"\n)\n\n// Read reads any variable, trimmed.\nfunc Read(name string) string { return strings.TrimSpace(os.Getenv(name)) }\n"},
		"rule 8: os.Getenv must name a constant"},
	{"rule 8 duplicated name",
		edit{rel: "pkg/a/lit.go", write: "package a\n\nconst other = \"FX_DATA_DIR\"\n"},
		`rule 8: literal "FX_DATA_DIR" duplicates env.DataDir`},
	{"rule 8 undeclared pattern",
		edit{rel: "pkg/a/lit.go", write: "package a\n\nconst other = \"FX_OTHER\"\n"},
		`rule 8: literal "FX_OTHER" looks like an environment variable`},
	{"rule 9",
		edit{rel: "contract/events.json", write: "{}\n"},
		"rule 9: the event enumerator found nothing"},
	{"rule 10",
		edit{rel: "pkg/b/b_test.go", write: "package b\n\nimport . \"example.com/fixture/claims\"\n\nvar _ = SW_001\n"},
		"rule 10: the entries package must not be dot-imported"},
	{"rule 11 observed without a date",
		edit{rel: "claims/facts.go", old: `Observed: "2026-09-24",`, repl: `Observed: "",`},
		"rule 11: HealthIsJSON is observed and has no Observed date"},
	{"rule 11 date form",
		edit{rel: "claims/facts.go", old: `Observed: "2026-09-24",`, repl: `Observed: "yesterday",`},
		`rule 11: HealthIsJSON has Observed "yesterday"; want YYYY-MM-DD`},
	{"rule 11 date that is not on the calendar",
		edit{rel: "claims/facts.go", old: `Observed: "2026-09-24",`, repl: `Observed: "2026-02-30",`},
		`rule 11: HealthIsJSON has Observed "2026-02-30"; want YYYY-MM-DD`},
	{"rule 11 stray Observed on a documented observation",
		edit{rel: "claims/facts.go", old: `Source:   "https://man7.org/linux/man-pages/man2/flock.2.html",`, repl: `Source:   "https://man7.org/linux/man-pages/man2/flock.2.html", Observed: "someday",`},
		`rule 11: LockIsAdvisory has Observed "someday"; want YYYY-MM-DD`},
	{"rule 11 observed without a probe",
		edit{rel: "pkg/a/a_probe_test.go", old: "chartertest.Observes(t, claims.HealthIsJSON)", repl: "t.Log(\"nothing observed\")"},
		"rule 11: observed observation HealthIsJSON is cited by no probe"},
	{"rule 11 documented without a source",
		edit{rel: "claims/facts.go", old: `Source:   "https://man7.org/linux/man-pages/man2/flock.2.html",`, repl: `Source:   "",`},
		"rule 11: LockIsAdvisory is documented and has no Source"},
	{"rule 11 subject",
		edit{rel: "claims/facts.go", old: `Subject:  "os",`, repl: `Subject:  "",`},
		"rule 11: LockIsAdvisory has no Subject"},
	{"rule 12 observation cited from a test",
		edit{rel: "pkg/a/a_test.go", old: "chartertest.Verifies(t, claims.SW_002)", repl: "chartertest.Verifies(t, claims.SW_002)\n\t_ = claims.Suspect"},
		"rule 12: observation Suspect is referenced from pkg/a/a_test.go, which is not a probe file"},
	{"rule 12 probe cites a spec",
		edit{rel: "pkg/a/a_probe_test.go", add: "\nfunc TestProbeFold(t *testing.T) { _ = claims.CON_002 }\n"},
		"rule 12: probe file pkg/a/a_probe_test.go references entry CON_002; a probe cites observations only"},
	{"rule 13 broken under active",
		edit{rel: "claims/facts.go", old: `Status:   charter.Observed,`, repl: `Status:   charter.Broken,`},
		"rule 13: broken observation HealthIsJSON supports active spec CON-002"},
	{"rule 21 file name",
		edit{rel: "claims/specs/sw-4-bad.md", write: "**Title**\nx\n"},
		"rule 21: a spec file is named <LENS>-<NNN>-<slug>.md"},
	{"rule 21 declared twice",
		edit{rel: "claims/specs/SW-001-again.md", write: "**Title**\nAgain\n\n**Lens**: SW\n\n**Status**: planned\n\n**Description**\nx\n\n**Rationale**\nx\n\n**Verification Description**\nx\n"},
		"rule 21: SW-001 is declared twice"},
	{"rule 21 missing field",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "**Rationale**\nThe fixture needs it.\n\n", repl: ""},
		"rule 21: CON-001 has no Rationale"},
	{"rule 21 status",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "**Status**: active", repl: "**Status**: keep"},
		`rule 21: status "keep" is not planned, active or deprecated`},
	{"rule 21 lens",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "**Lens**: CON", repl: "**Lens**: SW"},
		"rule 21: Lens is SW but the file is named for CON"},
	{"rule 22 stale",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "The lock file is removed on close", repl: "The lock file goes on close"},
		"rule 22: the generated Go differs from the specs"},
	{"rule 23 unknown target",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "- [SW-001](SW-001-bucket-survives-restart.md)", repl: "- [SW-007](SW-007-nowhere.md)"},
		"rule 23: CON-001 realizes SW-007, which is not a spec"},
	{"rule 23 wrong file",
		edit{rel: "claims/specs/CON-001-lock-file-removed-on-close.md", old: "(SW-001-bucket-survives-restart.md)", repl: "(SW-001-old-slug.md)"},
		"rule 23: the link to SW-001 points at SW-001-old-slug.md; its file is SW-001-bucket-survives-restart.md"},
	{"rule 23 stakeholder realizes",
		edit{rel: "claims/specs/STK-001-come-back-to-everything.md", add: "\n## Relations\n\n**Realizes**\n\n- [STK-002](STK-002-run-anywhere.md)\n"},
		"rule 23: STK-001 is a stakeholder spec and realizes nothing"},
	{"rule 23 observation supports",
		edit{rel: "claims/facts.go", old: `Supports: []string{"CON-001"},`, repl: `Supports: []string{"CON-009"},`},
		"rule 23: observation LockIsAdvisory supports CON-009, which is not a spec"},
	{"rule 24 unknown",
		edit{rel: "pkg/a/a.go", old: "//charter:realizes SW-001, CON-001, CON-002", repl: "//charter:realizes SW-001, CON-001, CON-002\n//charter:concerns SW-009"},
		"rule 24: //charter:concerns names SW-009, which is not a spec id"},
	{"rule 24 planned",
		edit{rel: "pkg/a/a.go", old: "//charter:realizes SW-001, CON-001, CON-002", repl: "//charter:realizes SW-001, CON-001, CON-002, SW-003"},
		"rule 24: //charter:realizes names planned spec SW-003"},
	{"rule 24 typescript",
		edit{rel: "web/helper.ts", add: "\n// charter:realizes SW-008\nexport const x = 1;\n"},
		"web/helper.ts:7: rule 24: //charter:realizes names SW-008, which is not a spec id"},
	{"rule 25 not realized",
		edit{rel: "pkg/a/a.go", old: "//charter:realizes SW-001, CON-001, CON-002", repl: "//charter:realizes SW-001, CON-001"},
		"rule 25: active spec CON-002 is realized by no code; add //charter:realizes CON-002 where it is implemented"},
	{"rule 25 concerns alone",
		edit{rel: "pkg/a/a.go", old: "//charter:realizes SW-001, CON-001, CON-002", repl: "//charter:realizes SW-001, CON-001\n//charter:concerns CON-002"},
		"rule 25: active spec CON-002 is realized by no code"},
	{"rule 25 not verified",
		edit{rel: "pkg/a/a_test.go", old: "\tchartertest.Verifies(t, claims.STK_001)\n", repl: ""},
		"rule 25: active spec STK-001 is verified by no test"},
	{"rule 25 stale waiver",
		edit{rel: "claims/waivers.txt", write: "SW-001\tthe fixture waives nothing it needs\n"},
		"claims/waivers.txt:1: rule 25: the waiver for SW-001 absorbs no active spec or scenario that lacks a test"},
	{"rule 25 glob waiver",
		edit{rel: "claims/waivers.txt", write: "SW-*\tthe fixture waives nothing it needs\n"},
		`claims/waivers.txt:1: rule 25: the waiver "SW-*" is a glob`},
	{"rule 24 later statement",
		edit{rel: "pkg/a/a_test.go", old: "\tchartertest.Verifies(t, claims.CON_001)\n", repl: "\tt.Log(\"first\")\n\tchartertest.Verifies(t, claims.CON_001)\n"},
		"pkg/a/a_test.go:18: rule 24: Verifies is the leading statement of its test"},
	{"rule 24 not a test name",
		edit{rel: "pkg/a/a_test.go", old: "func TestLock(", repl: "func Testlock("},
		"rule 24: Verifies sits in Testlock, which go test never runs"},
	{"rule 24 a method",
		edit{rel: "pkg/a/a_test.go", old: "func TestLock(", repl: "type suite struct{}\n\nfunc (suite) TestLock("},
		"rule 24: Verifies sits in the method TestLock, which go test never runs"},
	{"rule 24 subtest in a closure",
		edit{rel: "pkg/a/a_test.go", old: "\tchartertest.Verifies(t, claims.SW_002)\n}", repl: "\t_ = func() {\n\t\tt.Run(\"lock\", func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.SW_002)\n\t\t})\n\t}\n}"},
		"pkg/a/a_test.go:20: rule 24: Verifies is the leading statement of its test"},
	{"rule 24 subtest under a branch",
		edit{rel: "pkg/a/a_test.go", old: "\tchartertest.Verifies(t, claims.SW_002)\n}", repl: "\tif false {\n\t\tt.Run(\"lock\", func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.SW_002)\n\t\t})\n\t}\n}"},
		"pkg/a/a_test.go:20: rule 24: Verifies is the leading statement of its test"},
	{"rule 24 directive outside a doc comment",
		edit{rel: "pkg/a/a.go", old: "func DataDir() string { return os.Getenv(env.DataDir) }", repl: "func DataDir() string {\n\t//charter:concerns SW-001\n\treturn os.Getenv(env.DataDir)\n}"},
		"rule 24: //charter:concerns is written in the doc comment of the declaration that carries the decision; this comment documents none"},
	{"rule 6 a deprecated entry accounts for nothing",
		edit{rel: "claims/specs/CON-002-fold-never-panics.md", old: "**Status**: active", repl: "**Status**: deprecated"},
		"rule 6: env:FX_TIMEOUT is neither accounted for by an entry nor listed in claims/baseline.txt"},
	{"rule 24 skipped typescript test",
		edit{rel: "web/b.spec.ts", old: `test("windows is not supported yet"`, repl: `test.skip("windows is not supported yet"`},
		"web/b.spec.ts:6: rule 24: verifies sits in a test that is skip, fixme, fail or todo"},
	{"rule 24 fixme typescript test",
		edit{rel: "web/b.spec.ts", old: `test("windows is not supported yet"`, repl: `test.fixme("windows is not supported yet"`},
		"rule 24: verifies sits in a test that is skip, fixme, fail or todo"},
	{"rule 24 expected-failure typescript test",
		edit{rel: "web/b.spec.ts", old: `test("windows is not supported yet"`, repl: `test.fail("windows is not supported yet"`},
		"rule 24: verifies sits in a test that is skip, fixme, fail or todo"},
	{"rule 24 local verifies",
		edit{rel: "web/a.test.ts", old: `import { verifies } from "./helper";`, repl: `const verifies = (..._: unknown[]): void => {};`},
		"web/a.test.ts:7: rule 24: verifies is declared in the test file itself"},
	{"rule 24 later typescript statement",
		edit{rel: "web/b.spec.ts", old: "  verifies(claims.SW_002);\n", repl: "  await Promise.resolve();\n  verifies(claims.SW_002);\n"},
		"web/b.spec.ts:7: rule 24: verifies is the leading statement of its test"},
	{"rule 24 typescript verifies under an if",
		edit{rel: "web/b.spec.ts", old: "  verifies(claims.SW_002);\n", repl: "  if (process.platform) {\n    verifies(claims.SW_002);\n  }\n"},
		"web/b.spec.ts:7: rule 24: verifies is the leading statement of its test"},
	{"rule 24 typescript verifies in a closure",
		edit{rel: "web/b.spec.ts", old: "  verifies(claims.SW_002);\n", repl: "  const cite = () => {\n    verifies(claims.SW_002);\n  };\n  cite();\n"},
		"web/b.spec.ts:7: rule 24: verifies is the leading statement of its test"},
	{"rule 24 skipped typescript describe",
		edit{rel: "web/a.test.ts", old: `describe("bucket"`, repl: `describe.skip("bucket"`},
		"web/a.test.ts:7: rule 24: verifies sits in a test that is skip, fixme, fail or todo"},
	{"rule 24 fixme typescript describe",
		edit{rel: "web/b.spec.ts", old: "test(\"windows is not supported yet\", async () => {\n  verifies(claims.SW_002);\n});\n",
			repl: "test.describe.fixme(\"platforms\", () => {\n  test(\"windows is not supported yet\", async () => {\n    verifies(claims.SW_002);\n  });\n});\n"},
		"web/b.spec.ts:7: rule 24: verifies sits in a test that is skip, fixme, fail or todo"},
	{"rule 25 waiver without a reason",
		edit{rel: "claims/waivers.txt", write: "SW-001\n"},
		"rule 25: the waiver for SW-001 gives no reason"},
}

// tripRules applies each case named, or every case when none is, to a copy
// of the clean fixture of its own, and wants the problem it trips.
func tripRules(t *testing.T, names ...string) {
	t.Helper()
	ran := 0
	for _, c := range ruleCases {
		if len(names) > 0 && !slices.Contains(names, c.name) {
			continue
		}
		ran++
		t.Run(c.name, func(t *testing.T) {
			root := copyFixture(t)
			apply(t, root, c.edit)
			problems, err := Problems(fixtureOptions(root))
			if err != nil {
				t.Fatal(err)
			}
			if !anyContains(problems, c.want) {
				t.Errorf("want a problem containing %q, got:\n%s", c.want, strings.Join(problems, "\n"))
			}
		})
	}
	if len(names) > 0 && ran != len(names) {
		t.Fatalf("%d of the cases %q are rule cases", ran, names)
	}
}

// Every rule has a mutation of the clean fixture that trips it.
func TestCheckRules(t *testing.T) {
	Verifies(t, claims.CH_001, "every-rule-has-a-mutation")
	Verifies(t, claims.CH_004, "unaccounted-surface-fails", "baseline-holds-only-live-gaps", "env-read-through-a-constant")
	Verifies(t, claims.CH_005, "observed-fact-has-a-date-and-a-probe", "broken-fact-under-active-fails")
	tripRules(t)
}

// A stakeholder spec that realizes something, and an observation that
// supports no spec, each fail.
func TestTheRelationRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_015, "stakeholder-realizes-nothing", "observation-supports-existing")
	tripRules(t, "rule 23 stakeholder realizes", "rule 23 observation supports")
}

// Generated Go that differs from the specs fails.
func TestStaleGeneratedGoTrips(t *testing.T) {
	Verifies(t, claims.CH_020, "stale-generated-fails")
	tripRules(t, "rule 22 stale")
}

// A TypeScript anchor naming no entry, and a directive outside a doc
// comment, each fail.
func TestTheAnchorFormRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_021, "typescript-anchor-resolves", "comment-elsewhere-refused")
	tripRules(t, "rule 24 typescript", "rule 24 directive outside a doc comment")
}

// An anchor naming no entry or a planned one fails, and a decision anchored
// by concerns alone is realized by no code.
func TestTheAnchorTargetRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_022, "unknown-id-refused", "planned-refused", "concerns-is-not-realized")
	tripRules(t, "rule 24 unknown", "rule 24 planned", "rule 25 concerns alone")
}

// A citation after another statement, under a block or in a closure, in a
// function go test never runs, in a disabled TypeScript test or suite, or
// through a verifies the file declares fails.
func TestTheCitationRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_023, "later-statement-refused", "non-test-name-refused", "disabled-ts-refused", "local-verifies-refused")
	tripRules(t, "rule 24 later statement", "rule 24 not a test name", "rule 24 skipped typescript test",
		"rule 24 fixme typescript test", "rule 24 expected-failure typescript test", "rule 24 local verifies",
		"rule 24 later typescript statement", "rule 24 typescript verifies under an if", "rule 24 typescript verifies in a closure",
		"rule 24 skipped typescript describe", "rule 24 fixme typescript describe")
}

// A citation in a method, or in a subtest that a closure or a branch
// encloses, sits where go test may never run it and fails.
func TestACitationGoTestNeverReachesTrips(t *testing.T) {
	Verifies(t, claims.CH_023, "later-statement-refused", "non-test-name-refused")
	tripRules(t, "rule 24 a method", "rule 24 subtest in a closure", "rule 24 subtest under a branch")
}

// A function, or a spec written in Go, in the entries package fails.
func TestTheEntriesPackageRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_025, "functions-refused", "spec-in-go-refused")
	tripRules(t, "rule 1 function", "rule 1 a spec written in Go")
}

// An active decision with no realizes anchor fails naming the directive to
// add, and so does a waiver that absorbs nothing, gives no reason or is a
// glob.
func TestTheCoverageRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_030, "not-realized", "idle-waiver-fails", "glob-refused")
	tripRules(t, "rule 25 not realized", "rule 25 stale waiver", "rule 25 waiver without a reason", "rule 25 glob waiver")
}

// A test file referencing an observation fails.
func TestAnObservationReferencedFromATestTrips(t *testing.T) {
	Verifies(t, claims.CH_032, "test-citing-an-observation-fails")
	tripRules(t, "rule 12 observation cited from a test")
}

// A surface no entry accounts for, an Accounts naming nothing enumerated,
// a baseline line that is not a live gap once, an empty enumerator and a
// deprecated entry's Accounts each fail.
func TestTheSurfaceRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_033, "unaccounted-fails", "accounts-must-exist", "baseline-holds-live-gaps-once", "empty-enumerator", "deprecated-accounts-nothing")
	tripRules(t, "rule 6", "rule 7 accounts", "rule 7 baseline gone", "rule 7 baseline accounted", "rule 7 baseline malformed",
		"rule 7 baseline duplicate", "rule 9", "rule 6 a deprecated entry accounts for nothing")
}

// A read by a literal or a variable, a literal naming a declared variable
// and a literal matching a pattern each fail.
func TestTheEnvironmentRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_034, "literal-argument", "variable-argument", "duplicated-name", "undeclared-pattern")
	tripRules(t, "rule 8 literal argument", "rule 8 variable argument", "rule 8 read inside another call",
		"rule 8 duplicated name", "rule 8 undeclared pattern")
}

// An observation without the fields its status needs, or broken under an
// active entry, fails.
func TestTheObservationRulesTrip(t *testing.T) {
	Verifies(t, claims.CH_035, "observed-needs-date-and-probe", "documented-needs-source", "subject-required", "broken-under-active")
	tripRules(t, "rule 11 observed without a date", "rule 11 date form", "rule 11 date that is not on the calendar",
		"rule 11 observed without a probe", "rule 11 documented without a source", "rule 11 stray Observed on a documented observation",
		"rule 11 subject", "rule 13 broken under active")
}

// An active spec cited by neither a Go nor a TypeScript test trips rule 25,
// whether or not a TypeScript citation of some other spec remains.
func TestAnActiveSpecCitedByNoTestIsReported(t *testing.T) {
	Verifies(t, claims.CH_030, "spec-unverified")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.SW_001)\n", "")
	replaceIn(t, root, "web/a.test.ts", "verifies(c.SW_001);", "verifies(c.CON_002);")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(problems, "rule 25: active spec SW-001 is verified by no test") {
		t.Errorf("want a problem naming SW-001, got:\n%s", strings.Join(problems, "\n"))
	}
}

// A verifies reached through a property, x.verifies(...), is no citation:
// the claims access in it is a reference and counts for no scenario.
func TestAMemberVerifiesIsAReference(t *testing.T) {
	Verifies(t, claims.CH_023, "reference-is-not-citation")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.SW_001)\n", "")
	replaceIn(t, root, "web/a.test.ts", `import { verifies } from "./helper";`, `import * as helper from "./helper";`)
	replaceIn(t, root, "web/a.test.ts", "verifies(c.SW_001);", "helper.verifies(c.SW_001);")
	replaceIn(t, root, "web/a.test.ts", "verifies(c.CON_002);", "helper.verifies(c.CON_002);")
	replaceIn(t, root, "web/a.test.ts", "verifies(c.CON_001);", "helper.verifies(c.CON_001);")
	reports(t, root, "rule 25: active spec SW-001 is verified by no test")
}

// A waiver accepts a missing test and never missing code.
func TestAWaiverAcceptsAMissingTestOnly(t *testing.T) {
	Verifies(t, claims.CH_030, "waiver-never-absorbs-an-anchor")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.CON_001)\n", "")
	replaceIn(t, root, "web/a.test.ts", "verifies(c.CON_001);", "")
	writeFile(t, root, "claims/waivers.txt", "CON-001\tverified by review\n")
	noProblems(t, root)
	replaceIn(t, root, "pkg/a/a.go", "//charter:realizes SW-001, CON-001, CON-002", "//charter:realizes SW-001, CON-002")
	reports(t, root, "rule 25: active spec CON-001 is realized by no code")
}

// A stakeholder spec needs a test and no anchor.
func TestAStakeholderSpecNeedsNoAnchor(t *testing.T) {
	Verifies(t, claims.CH_031, "need-without-anchor-passes")
	root := copyFixture(t)
	noProblems(t, root)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.STK_001)\n", "")
	reports(t, root, "rule 25: active spec STK-001 is verified by no test")
}

// A reference from a probe file is not a citation, and is a rule 12
// problem of its own.
func TestProbeReferenceIsNoCitation(t *testing.T) {
	Verifies(t, claims.CH_032, "probe-reference-is-no-citation")
	Verifies(t, claims.CH_005, "probes-and-tests-stay-apart")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.CON_002)\n", "")
	replaceIn(t, root, "web/a.test.ts", "verifies(c.CON_002);", "verifies(c.SW_001);")
	appendTo(t, root, "pkg/a/a_probe_test.go", "\nfunc TestProbeFold(t *testing.T) { _ = claims.CON_002 }\n")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"rule 25: active spec CON-002 is verified by no test",
		"rule 12: probe file pkg/a/a_probe_test.go references entry CON_002",
	} {
		if !anyContains(problems, want) {
			t.Errorf("missing %q in:\n%s", want, strings.Join(problems, "\n"))
		}
	}
}

// Broken under a planned or a deprecated spec is nobody's problem yet:
// nothing active rests on the observation.
func TestBrokenUnderInactiveIsNoProblem(t *testing.T) {
	Verifies(t, claims.CH_035, "broken-under-active")
	Verifies(t, claims.CH_005, "broken-fact-under-active-fails")
	for _, parent := range []string{"SW-003", "SW-002"} {
		t.Run(parent, func(t *testing.T) {
			root := copyFixture(t)
			replaceIn(t, root, "claims/facts.go", `Supports: []string{"CON-001"},`, `Supports: []string{"`+parent+`"},`)
			replaceIn(t, root, "claims/facts.go", "Status:   charter.Documented,", "Status:   charter.Broken,")
			problems, err := Problems(fixtureOptions(root))
			if err != nil {
				t.Fatal(err)
			}
			if anyContains(problems, "rule 13") {
				t.Fatalf("rule 13 tripped under %s:\n%s", parent, strings.Join(problems, "\n"))
			}
		})
	}
}

// A negated probe tag does not make a probe: a file built without the tag
// is an ordinary test, and citing an observation from it is rule 12.
func TestNegatedProbeTagIsNoProbe(t *testing.T) {
	Verifies(t, claims.CH_032, "negated-tag-is-no-probe")
	root := copyFixture(t)
	writeFile(t, root, "pkg/b/not_probe_test.go", `//go:build !fxprobe

package b

import (
	"testing"

	"example.com/fixture/claims"
)

func TestNot(t *testing.T) { _ = claims.Suspect }
`)
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	want := "rule 12: observation Suspect is referenced from pkg/b/not_probe_test.go, which is not a probe file"
	if !anyContains(problems, want) {
		t.Fatalf("want a problem containing %q, got:\n%s", want, strings.Join(problems, "\n"))
	}
}

// Observations need a probe tag to be checked against. Without one, that
// is the single problem: no file can be a probe, so nothing is blamed on one.
func TestObservationsNeedAProbeTag(t *testing.T) {
	Verifies(t, claims.CH_032, "tag-required")
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.ProbeTag = ""
	problems, err := Problems(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(problems, "rule 11: observations need Options.ProbeTag") {
		t.Fatalf("no configuration problem reported:\n%s", strings.Join(problems, "\n"))
	}
	if anyContains(problems, "rule 12") || anyContains(problems, "cited by no probe") {
		t.Fatalf("probe-dependent problems reported without a probe tag:\n%s", strings.Join(problems, "\n"))
	}
}

// A tree without build tags is valid: rule 9 does not apply to tags.
func TestTagsMayBeAbsent(t *testing.T) {
	Verifies(t, claims.CH_033, "empty-enumerator")
	root := copyFixture(t)
	for _, rel := range []string{"pkg/a/a.go", "pkg/b/b.go", "pkg/b/win.go"} {
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatal(err)
		}
	}
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if anyContains(problems, "rule 9: the tag enumerator") {
		t.Errorf("rule 9 fired for tags:\n%s", strings.Join(problems, "\n"))
	}
	if !anyContains(problems, "rule 7") {
		t.Errorf("want the stale tag:desktop baseline line reported, got:\n%s", strings.Join(problems, "\n"))
	}
}

func TestDynamicDirectiveExemptsARead(t *testing.T) {
	Verifies(t, claims.CH_034, "dynamic-directive")
	root := copyFixture(t)
	writeFile(t, root, "pkg/a/read.go", "package a\n\nimport \"os\"\n\n// Read reads the variable the configuration names.\nfunc Read(name string) string {\n\treturn os.Getenv(name) //charter:dynamic the name comes from config\n}\n")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if anyContains(problems, "rule 8") {
		t.Errorf("the directive did not exempt the read:\n%s", strings.Join(problems, "\n"))
	}
}

func TestBaselineUpdateIsAFixedPoint(t *testing.T) {
	Verifies(t, claims.CH_033, "update-rebuilds")
	root := copyFixture(t)
	baseline := filepath.Join(root, "claims", "baseline.txt")
	if err := os.Remove(baseline); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "pkg/b/wasm.go", "//go:build wasm\n\npackage b\n")

	t.Setenv(UpdateEnv, "1")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("update run reported:\n%s", strings.Join(problems, "\n"))
	}
	written, err := os.ReadFile(baseline)
	if err != nil {
		t.Fatal(err)
	}
	want := baselineHeader + "op:GET /health\t" + initialReason + "\ntag:desktop\t" + initialReason + "\ntag:wasm\t" + initialReason + "\n"
	if string(written) != want {
		t.Fatalf("baseline:\n%s\nwant:\n%s", written, want)
	}

	t.Setenv(UpdateEnv, "")
	problems, err = Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("second run reported:\n%s", strings.Join(problems, "\n"))
	}
	again, err := os.ReadFile(baseline)
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != want {
		t.Fatal("the second run changed the baseline")
	}
}

func TestBaselineUpdateKeepsReasons(t *testing.T) {
	Verifies(t, claims.CH_033, "update-rebuilds")
	root := copyFixture(t)
	writeFile(t, root, "pkg/b/wasm.go", "//go:build wasm\n\npackage b\n")
	t.Setenv(UpdateEnv, "1")
	if _, err := Problems(fixtureOptions(root)); err != nil {
		t.Fatal(err)
	}
	written, err := os.ReadFile(filepath.Join(root, "claims", "baseline.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{"op:GET /health\tnobody claims the health check yet", "tag:desktop\tpackaging goes with R2", "tag:wasm\t" + initialReason} {
		if !strings.Contains(string(written), line+"\n") {
			t.Errorf("baseline lacks %q:\n%s", line, written)
		}
	}
}
