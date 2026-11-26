package chartertest

import (
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

const sampleArea = "# Refill\nCode: pkg/a\n\n" +
	"## SW-010 The bucket refills from its last request at the rate\n" +
	"Status: active\n" +
	"Kind: invariant\n" +
	"Serves: STK-001\n" +
	"See: SW-001\n" +
	"Accounts: env:FX_DATA_DIR\n" +
	"\n" +
	"At each request the bucket reads the clock.\n" +
	"When time passed, the earned tokens are added to the bucket.\n" +
	"\n" +
	"Why: a caller that waits is let through again at the rate.\n" +
	"\n" +
	"Scenarios:\n" +
	"1. `waited-request-passes`: a request after a wait of one token passes.\n" +
	"2. `burst-caps-tokens`: a long wait earns no more than the burst,\n" +
	"   and the excess is not kept.\n" +
	"\n" +
	"Guidance: no timer; the clock is read at a request only.\n" +
	"\n" +
	"Changes:\n" +
	"- 2026-09-28: a denied request reads the clock once more (backlog 034).\n"

// refillGo realizes SW-010 in the fixture.
const refillGo = "package a\n\n// Refill adds the tokens earned.\n//\n//charter:realizes SW-010\nfunc Refill() {}\n"

// refillTest verifies both scenarios of SW-010 from a Go test.
const refillTest = "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n" +
	"func TestRefillInserts(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_010, \"waited-request-passes\")\n}\n\n" +
	"func TestRefillKeepsTheBranch(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_010, \"burst-caps-tokens\")\n}\n"

func TestParseAreaReadsEveryPart(t *testing.T) {
	Verifies(t, claims.CH_010, "every-part-read")
	doc, entries, problems := parseArea("claims/refill.md", sampleArea)
	if len(problems) != 0 {
		t.Fatalf("problems: %v", problems)
	}
	if doc.Title != "Refill" || strings.Join(doc.Code, ",") != "pkg/a" || doc.Needs {
		t.Errorf("doc = %+v", doc)
	}
	if len(entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(entries))
	}
	e := entries[0]
	if e.ID != "SW-010" || e.Name != "SW_010" || e.Form != formArea || e.Need || e.Lens != "SW" || e.Number != 10 || e.Pos != "claims/refill.md:4" {
		t.Errorf("identity = %+v", e)
	}
	if e.Text != "The bucket refills from its last request at the rate" || e.Status != "active" || e.Kind != "invariant" {
		t.Errorf("title %q, status %q, kind %q", e.Text, e.Status, e.Kind)
	}
	if strings.Join(e.Serves, ",") != "STK-001" || strings.Join(e.See, ",") != "SW-001" || strings.Join(e.Accounts, ",") != "env:FX_DATA_DIR" {
		t.Errorf("serves %v, see %v, accounts %v", e.Serves, e.See, e.Accounts)
	}
	if !strings.HasPrefix(e.Description, "At each request") || strings.Count(e.Description, "\n") != 1 {
		t.Errorf("description = %q", e.Description)
	}
	if !strings.HasPrefix(e.Rationale, "a caller that waits") || !strings.HasPrefix(e.Guidance, "no timer") {
		t.Errorf("why = %q, guidance = %q", e.Rationale, e.Guidance)
	}
	if len(e.Scenarios) != 2 || e.Scenarios[0].Handle != "waited-request-passes" || e.Scenarios[0].Line != 17 ||
		e.Scenarios[1].Handle != "burst-caps-tokens" || e.Scenarios[1].Text != "a long wait earns no more than the burst, and the excess is not kept." {
		t.Errorf("scenarios = %+v", e.Scenarios)
	}
	if len(e.Changes) != 1 || e.Changes[0].Date != "2026-09-28" || !strings.HasPrefix(e.Changes[0].Text, "a denied request") {
		t.Errorf("changes = %+v", e.Changes)
	}
}

func TestParseAreaForm(t *testing.T) {
	Verifies(t, claims.CH_010, "heading-form", "fields-checked", "required-parts", "scenario-form", "parts-in-order", "change-form", "problems-locate")
	Verifies(t, claims.CH_002, "problems-name-file-line-and-fix")
	for _, tc := range []struct{ name, old, repl, want string }{
		{"heading", "## SW-010 The bucket", "## The bucket", `an entry's heading is "## ID Title"`},
		{"unknown field", "Kind: invariant", "Lens: SW", `"Lens" is not a field of an entry`},
		{"kind", "Kind: invariant", "Kind: rule", `kind "rule" is not invariant, bound or structure`},
		{"field twice", "Kind: invariant", "Kind: invariant\nKind: bound", "Kind appears twice"},
		{"no status", "Status: active\n", "", "SW-010 has no Status"},
		{"surface key", "Accounts: env:FX_DATA_DIR", "Accounts: health", `"health" is not a surface key`},
		{"no text", "At each request the bucket reads the clock.\nWhen time passed, the earned tokens are added to the bucket.\n\n", "", "SW-010 has no text"},
		{"no why", "Why: a caller that waits is let through again at the rate.\n\n", "", "SW-010 has no Why"},
		{"no scenarios", "Scenarios:\n1. `waited-request-passes`: a request after a wait of one token passes.\n2. `burst-caps-tokens`: a long wait earns no more than the burst,\n   and the excess is not kept.\n\n", "", "SW-010 is active and lists no Scenarios"},
		{"scenario form", "1. `waited-request-passes`: a request", "- waited-request-passes: a request", "a scenario is \"N. `handle`: sentence\""},
		{"handle form", "`waited-request-passes`", "`Waited Request`", `handle "Waited Request" is not lower-case kebab`},
		{"handle twice", "`burst-caps-tokens`", "`waited-request-passes`", `handle "waited-request-passes" appears twice`},
		{"order", "Scenarios:\n", "Guidance: early.\n\nScenarios:\n", "Scenarios is out of order"},
		{"text after a part", "Scenarios:\n", "Stray words.\n\nScenarios:\n", "text after Why"},
		{"change row", "- 2026-09-28: a denied request", "- 2026-09-28 — a denied request", `a change is "- YYYY-MM-DD: what changed and why"`},
		{"change date", "- 2026-09-28:", "- 2026-13-40:", "2026-13-40 is not a date"},
		{"heading inside", "Why: a caller", "### Why\n\nWhy: a caller", "a heading inside an entry"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(sampleArea, tc.old, tc.repl, 1)
			if src == sampleArea {
				t.Fatalf("the edit %q did not apply", tc.old)
			}
			_, _, problems := parseArea("claims/refill.md", src)
			if !anyContains(problems, tc.want) {
				t.Errorf("want %q in %v", tc.want, problems)
			}
			for _, p := range problems {
				if !strings.HasPrefix(p, "claims/refill.md:") || !strings.Contains(p, ": form: ") {
					t.Errorf("problem %q does not name the document, the line and the form rule", p)
				}
			}
		})
	}
}

// A status outside planned, active and deprecated is a form problem at the
// status line, and the three are none.
func TestAStatusOutsideTheThreeIsAFormProblem(t *testing.T) {
	Verifies(t, claims.CH_013, "status-form")
	for _, status := range []string{"planned", "active", "deprecated"} {
		if _, _, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "Status: active", "Status: "+status, 1)); len(problems) != 0 {
			t.Errorf("status %s: problems %v", status, problems)
		}
	}
	_, _, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "Status: active", "Status: keep", 1))
	if !anyContains(problems, `claims/refill.md:5: form: status "keep" is not planned, active or deprecated`) {
		t.Errorf("status keep: want a form problem at the status line, got %v", problems)
	}
}

// The header orients the reader and is not parsed beyond its title and
// Code line; a section heading in it is a form problem.
func TestTheHeaderIsProse(t *testing.T) {
	Verifies(t, claims.CH_010, "header-is-prose")
	header := "# Refill\nCode: pkg/a\n\n" +
		"A bucket reads the clock at each request.\n\n" +
		"- **burst**: the most tokens a bucket holds.\n" +
		"- **rate**: the tokens a bucket earns a second.\n\n" +
		"| Term | Is |\n|---|---|\n| rate | tokens a second |\n\n" +
		"```\n# a comment, not a title\n## SW-999 not an entry\nCode: not/this\nrate -> burst\n```\n\n" +
		"~~~text\n### nor this\n~~~\n\n"
	doc, entries, problems := parseArea("claims/refill.md", header)
	if len(problems) != 0 {
		t.Errorf("problems: %v", problems)
	}
	if len(entries) != 0 {
		t.Errorf("entries = %d, want none", len(entries))
	}
	if doc.Title != "Refill" || strings.Join(doc.Code, ",") != "pkg/a" {
		t.Errorf("doc = %+v", doc)
	}
	_, entries, problems = parseArea("claims/refill.md", header+sampleArea[len("# Refill\nCode: pkg/a\n\n"):])
	if len(problems) != 0 || len(entries) != 1 {
		t.Errorf("with an entry after the header: %d entries, problems %v", len(entries), problems)
	}
	_, _, problems = parseArea("claims/refill.md", "# Refill\n\n### Vocabulary\n\nprose\n\n"+sampleArea[len("# Refill\nCode: pkg/a\n\n"):])
	if !anyContains(problems, "claims/refill.md:3: form: a heading in the header") {
		t.Errorf("a section heading in the header: want a form problem, got %v", problems)
	}
}

// A line inside a fenced block in an entry is text, never a heading.
func TestAHeadingInAFenceIsText(t *testing.T) {
	Verifies(t, claims.CH_010, "heading-in-fence-is-text")
	for _, fenced := range []string{
		"```sh\n# a comment\n## SW-011 not an entry\nmake claims\n```\n",
		"~~~~\n# a comment\n~~~\n## still inside\n~~~~\n",
	} {
		src := strings.Replace(sampleArea, "\nWhy: a caller", "\n"+fenced+"\nWhy: a caller", 1)
		_, entries, problems := parseArea("claims/refill.md", src)
		if len(problems) != 0 {
			t.Errorf("%q: problems: %v", fenced, problems)
		}
		if len(entries) != 1 || !strings.Contains(entries[0].Description, "# a comment") || entries[0].Rationale == "" {
			t.Errorf("%q: entries = %+v", fenced, entries)
		}
	}
	// A fence among the field lines is no field line, but what it holds
	// is still text: its "## " line starts no entry.
	src := strings.Replace(sampleArea, "Kind: invariant\n", "```\n## SW-011 not an entry\n```\nKind: invariant\n", 1)
	_, entries, problems := parseArea("claims/refill.md", src)
	if len(entries) != 1 || entries[0].ID != "SW-010" {
		t.Errorf("a fence among the field lines: entries = %+v", entries)
	}
	if !anyContains(problems, "claims/refill.md:6: form: after the heading come field lines") {
		t.Errorf("a fence among the field lines: want a field-line problem on its opening line, got %v", problems)
	}
	// Once the fence closes, a heading is a heading again.
	src = strings.Replace(sampleArea, "\nWhy: a caller", "\n```\n# a comment\n```\n\n### Why\n\nWhy: a caller", 1)
	if _, _, problems := parseArea("claims/refill.md", src); !anyContains(problems, "a heading inside an entry") {
		t.Errorf("a heading after a closed fence: want a form problem, got %v", problems)
	}
}

// A need has no Kind, Serves or Guidance, and lives in needs.md.
func TestANeedHasNoKindServesOrGuidance(t *testing.T) {
	Verifies(t, claims.CH_011, "need-by-document", "need-fields-refused")
	src := strings.Replace(sampleArea, "## SW-010", "## STK-010", 1)
	_, entries, problems := parseArea("claims/needs.md", src)
	if len(entries) != 1 || !entries[0].Need || !entries[0].isNeed() {
		t.Fatalf("entries = %+v", entries)
	}
	for _, want := range []string{"a need has no Kind", "a need serves nothing", "a need has no Guidance"} {
		if !anyContains(problems, want) {
			t.Errorf("want %q in %v", want, problems)
		}
	}
}

// An STK id outside needs.md is a decision: active, it needs code that
// realizes it, and a Go test beside that code verifies it.
func TestAnSTKEntryOutsideNeedsIsADecision(t *testing.T) {
	Verifies(t, claims.CH_011, "need-by-document")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## STK-010", 1))
	writeFile(t, root, "pkg/a/refill_test.go", strings.ReplaceAll(refillTest, "claims.SW_010", "claims.STK_010"))
	problems := check(t, root)
	if !anyContains(problems, "claims/refill.md:4: rule 25: active spec STK-010 is realized by no code") {
		t.Errorf("an unanchored STK entry in an area document: want it reported unrealized, got:\n%s", strings.Join(problems, "\n"))
	}
	if anyContains(problems, "need STK-010") {
		t.Errorf("an STK entry in an area document is held to a need's rules:\n%s", strings.Join(problems, "\n"))
	}
	_, entries, _ := parseArea("claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## STK-010", 1))
	if len(entries) != 1 || entries[0].isNeed() {
		t.Errorf("entries = %+v, want one decision", entries)
	}
}

// An area entry is covered when code realizes it and a test verifies each
// of its scenarios by handle; the gate names the scenario that lacks one.
func TestAnAreaEntryIsCoveredPerScenario(t *testing.T) {
	Verifies(t, claims.CH_030, "scenario-unverified")
	Verifies(t, claims.CH_001, "missing-scenario-named")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", refillTest)
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a covered area entry reported:\n%s", strings.Join(problems, "\n"))
	}
	replaceIn(t, root, "pkg/a/refill_test.go", "\"burst-caps-tokens\"", "\"waited-request-passes\"")
	reports(t, root, "claims/refill.md:18: rule 25: scenario burst-caps-tokens of SW-010 is verified by no test; cite it with Verifies(t, claims.SW_010, \"burst-caps-tokens\")")
}

// A handle a test names is one the spec lists.
func TestAHandleMustBeAScenarioOfTheSpec(t *testing.T) {
	Verifies(t, claims.CH_023, "unknown-handle-refused")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, "\"burst-caps-tokens\"", "\"burst-caps-tokens\", \"denial-spends-nothing\"", 1))
	reports(t, root, "pkg/a/refill_test.go:15: rule 24: SW-010 has no scenario \"denial-spends-nothing\"")
}

// A Verifies sits inside a test function; anywhere else it is no citation.
func TestVerifiesSitsInsideATestFunction(t *testing.T) {
	Verifies(t, claims.CH_023, "outside-a-test-refused")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, "func TestRefillKeepsTheBranch(t *testing.T) {", "func helper(t *testing.T) {", 1))
	problems := check(t, root)
	for _, want := range []string{
		"pkg/a/refill_test.go:15: rule 24: Verifies sits outside a test function",
		"rule 25: scenario burst-caps-tokens of SW-010 is verified by no test",
	} {
		if !anyContains(problems, want) {
			t.Errorf("want %q in:\n%s", want, strings.Join(problems, "\n"))
		}
	}
}

// A scenario handle is a string literal, and the spec is named by its symbol.
func TestVerifiesArgumentsAreLiteral(t *testing.T) {
	Verifies(t, claims.CH_023, "literal-arguments")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, "\"burst-caps-tokens\"", "handle", 1))
	reports(t, root, "pkg/a/refill_test.go:15: rule 24: a scenario handle is a string literal")
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, "claims.SW_010, \"burst-caps-tokens\"", "spec, \"burst-caps-tokens\"", 1))
	reports(t, root, "pkg/a/refill_test.go:15: rule 24: Verifies names a spec by its symbol")
}

// A waiver absorbs one scenario, never a glob of them, never a missing anchor.
func TestAWaiverAbsorbsAScenario(t *testing.T) {
	Verifies(t, claims.CH_030, "waiver-absorbs-a-scenario", "glob-refused", "waiver-never-absorbs-an-anchor")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, "\"burst-caps-tokens\"", "\"waited-request-passes\"", 1))
	writeFile(t, root, "claims/waivers.txt", "SW-010/burst-caps-tokens\tthe burst case is verified by review\n")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a waived scenario reported:\n%s", strings.Join(problems, "\n"))
	}
	writeFile(t, root, "claims/waivers.txt", "SW-010/*\tthe burst case is verified by review\n")
	reports(t, root, "claims/waivers.txt:1: rule 25: the waiver \"SW-010/*\" is a glob")
	writeFile(t, root, "claims/waivers.txt", "SW-010/burst-caps-tokens\tthe burst case is verified by review\n")
	writeFile(t, root, "pkg/a/refill.go", "package a\n")
	reports(t, root, "rule 25: active spec SW-010 is realized by no code")
}

// A need is verified only from the stakeholder's side when the options name
// that side.
func TestANeedIsVerifiedFromTheStakeholderSide(t *testing.T) {
	Verifies(t, claims.CH_031, "offside-citation-refused", "onside-citation-counts")
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.StakeholderTests = []string{"web/*"}
	t.Setenv(UpdateSpecsEnv, "1")
	t.Setenv(UpdateTypeScriptEnv, "1")
	problems, err := Problems(opts)
	if err != nil {
		t.Fatal(err)
	}
	want := "pkg/a/a_test.go:13: rule 25: need STK-001 is verified from pkg/a/a_test.go, which is not on the stakeholder's side"
	if !anyContains(problems, want) {
		t.Errorf("want %q in:\n%s", want, strings.Join(problems, "\n"))
	}
	if !anyContains(problems, "rule 25: active spec STK-001 is verified by no test") {
		t.Errorf("the offside citation still counted:\n%s", strings.Join(problems, "\n"))
	}
	opts.StakeholderTests = []string{"pkg/a/*_test.go"}
	problems, err = Problems(opts)
	if err != nil {
		t.Fatal(err)
	}
	if anyContains(problems, "STK-001") {
		t.Errorf("a citation from the stakeholder's side was refused:\n%s", strings.Join(problems, "\n"))
	}
}

// A TypeScript test verifies scenarios by handle too.
func TestATypeScriptTestVerifiesScenarios(t *testing.T) {
	Verifies(t, claims.CH_023, "typescript-citation-read")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "web/refill.test.ts", strings.Join([]string{
		`import { it } from "@example/test";`,
		`import { claims } from "@fixture/claims";`,
		`import { verifies } from "./helper";`,
		``,
		`it("inserts", () => { verifies(claims.SW_010, "waited-request-passes"); });`,
		`it("keeps the branch", () => { verifies(claims.SW_010, "burst-caps-tokens"); });`,
		`it("does both", async () => {`,
		`  verifies(claims.SW_010, "waited-request-passes");`,
		`  verifies(claims.SW_010, "burst-caps-tokens");`,
		`  await Promise.resolve();`,
		`});`,
		``,
	}, "\n"))
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("scenarios verified from TypeScript reported:\n%s", strings.Join(problems, "\n"))
	}
	replaceIn(t, root, "web/refill.test.ts", `"burst-caps-tokens"); });`, `handle); });`)
	reports(t, root, "web/refill.test.ts:6: rule 24: a scenario handle is a string literal")
	writeFile(t, root, "web/refill.test.ts", strings.Join([]string{
		`import { claims } from "@fixture/claims";`,
		`import { verifies } from "./helper";`,
		``,
		`verifies(claims.SW_010, "burst-caps-tokens");`,
		``,
	}, "\n"))
	reports(t, root, "web/refill.test.ts:4: rule 24: verifies sits outside a test")
}

// Two entries with one id, across documents or forms, fail, and the problem
// says which of them takes a new id.
func TestAnIdIsDeclaredOnce(t *testing.T) {
	Verifies(t, claims.CH_012, "declared-once")
	Verifies(t, claims.CH_003, "id-declared-once")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## SW-001", 1))
	reports(t, root, "rule 21: SW-001 is declared twice, in claims/refill.md and claims/specs/SW-001-bucket-survives-restart.md; the entry that has not landed takes a new id: claimsgen next SW")
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "claims/again.md", "# Again\n\n"+strings.SplitN(sampleArea, "\n\n", 2)[1])
	reports(t, root, "claims/refill.md:4: form: SW-010 is declared twice, in claims/again.md and claims/refill.md; the entry that has not landed takes a new id: claimsgen next SW")
	apply(t, root, edit{rel: "claims/again.md", rm: true})
	writeFile(t, root, "claims/refill.md", sampleArea+"\n"+strings.SplitN(sampleArea, "\n\n", 2)[1])
	reports(t, root, "form: SW-010 is declared twice, in claims/refill.md and claims/refill.md; the entry that has not landed takes a new id: claimsgen next SW")
}

// Serves and See resolve, and a need serves nothing.
func TestServesAndSeeResolve(t *testing.T) {
	Verifies(t, claims.CH_015, "unknown-target", "self-relation")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "See: SW-001", "See: SW-099", 1))
	reports(t, root, "claims/refill.md:4: rule 23: SW-010 sees SW-099, which is not a spec")
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "Serves: STK-001", "Serves: SW-010", 1))
	reports(t, root, "rule 23: SW-010 serves itself")
}

// Lint says what the gate does not fail on.
func TestLintAdvises(t *testing.T) {
	Verifies(t, claims.CH_041, "serves-nothing", "planned-verified", "deprecated-still-anchored", "package-clause-only", "no-scenario-named")
	root := copyFixture(t)
	area := strings.Replace(sampleArea, "Serves: STK-001\n", "", 1)
	writeFile(t, root, "claims/refill.md", area)
	writeFile(t, root, "pkg/a/refill.go", "// Package a is the fixture.\n//\n//charter:realizes SW-010\npackage a\n")
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(refillTest, ", \"waited-request-passes\"", "", 1))
	problems, err := Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"claims/refill.md:4: lint: SW-010 serves nothing",
		"claims/refill.md:4: lint: SW-010 is realized only by a package clause",
		"pkg/a/refill_test.go:11: lint: Verifies names SW-010 without a scenario",
	} {
		if !anyContains(problems, want) {
			t.Errorf("want %q in:\n%s", want, strings.Join(problems, "\n"))
		}
	}
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "Status: active", "Status: planned", 1))
	writeFile(t, root, "pkg/a/refill.go", "package a\n")
	writeFile(t, root, "pkg/a/refill_test.go", refillTest)
	problems, err = Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(problems, "lint: SW-010 is planned and each of its scenarios is verified; set it active") {
		t.Errorf("want the planned advice in:\n%s", strings.Join(problems, "\n"))
	}
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "Status: active", "Status: deprecated", 1))
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	problems, err = Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(problems, "pkg/a/refill.go:5: lint: SW-010 is deprecated and still anchored here") {
		t.Errorf("want the strip advice in:\n%s", strings.Join(problems, "\n"))
	}
}

// The walk stays within Only.
func TestOnlyRestrictsTheWalk(t *testing.T) {
	Verifies(t, claims.CH_036, "only-restricts")
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.Only = []string{"pkg/b", "claims", "env"}
	t.Setenv(UpdateSpecsEnv, "1")
	t.Setenv(UpdateTypeScriptEnv, "1")
	problems, err := Problems(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(problems, "rule 25: active spec SW-001 is realized by no code") || !anyContains(problems, "rule 25: active spec SW-001 is verified by no test") {
		t.Errorf("pkg/a was walked despite Only:\n%s", strings.Join(problems, "\n"))
	}
}

// A Go directive with a space after the slashes is a directive too.
func TestASpacedDirectiveIsAnAnchor(t *testing.T) {
	Verifies(t, claims.CH_021, "spaced-directive-counts")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a.go", "//charter:realizes SW-001, CON-001, CON-002", "// charter:realizes SW-001, CON-001, CON-002")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a spaced directive reported:\n%s", strings.Join(problems, "\n"))
	}
}

// A scenario says who ruled it; nothing written is an agent's, provisional.
func TestAScenarioCarriesWhoRuledIt(t *testing.T) {
	Verifies(t, claims.CH_261009_VY5C, "owner-read", "delegate-read", "default-is-provisional", "malformed-refused")
	Verifies(t, claims.CH_261009_3FHW, "provenance-in-the-registry")
	src := strings.Replace(sampleArea, "1. `waited-request-passes`: a request", "1. `waited-request-passes` (owner, #444): a request", 1)
	src = strings.Replace(src, "2. `burst-caps-tokens`: a long wait", "2. `burst-caps-tokens` (delegate alice, #12 R3): a long wait", 1)
	src = strings.Replace(src, "   and the excess is not kept.\n", "   and the excess is not kept.\n3. `denial-spends-nothing` (agent, provisional): a denied request spends nothing.\n4. `state-carries`: the tokens carry across.\n", 1)
	_, entries, problems := parseArea("claims/refill.md", src)
	if len(problems) != 0 || len(entries) != 1 {
		t.Fatalf("entries %d, problems %q", len(entries), problems)
	}
	want := []Provenance{
		{Who: "owner", Ref: "#444", Written: true},
		{Who: "delegate", Name: "alice", Ref: "#12 R3", Written: true},
		{Who: "agent", Written: true},
		{Who: "agent"},
	}
	for i, w := range want {
		if got := entries[0].Scenarios[i].Provenance; got != w {
			t.Errorf("scenario %d provenance = %+v, want %+v", i+1, got, w)
		}
	}
	for i, w := range []string{"owner, #444", "delegate alice, #12 R3", "agent, provisional", "agent, provisional"} {
		if got := entries[0].Scenarios[i].Provenance.String(); got != w {
			t.Errorf("scenario %d renders %q, want %q", i+1, got, w)
		}
	}
	if !entries[0].Scenarios[0].Provenance.Ruled() || !entries[0].Scenarios[1].Provenance.Ruled() || entries[0].Scenarios[3].Provenance.Ruled() {
		t.Error("Ruled: the owner's and the delegate's are, the default is not")
	}

	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", src)
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	e := reg.Entry("SW_010")
	if e == nil || len(e.Scenarios) != 4 || e.Scenarios[0].Provenance.Who != "owner" || e.Scenarios[0].Provenance.Ref != "#444" || e.Scenarios[3].Provenance.Written {
		t.Errorf("the registry's provenance: %+v", e)
	}

	for _, tc := range []struct{ repl, want string }{
		{"(reviewer, #444)", `"reviewer" is not owner, delegate or agent`},
		{"(owner bob, #444)", "the owner's provenance is (owner) or (owner, #444), not (owner bob, #444)"},
		{"(delegate, #12)", "a delegate's provenance names the delegate"},
		{"(agent, #444)", "an agent's scenario is provisional until a person rules it"},
	} {
		bad := strings.Replace(src, "(owner, #444)", tc.repl, 1)
		_, _, problems := parseArea("claims/refill.md", bad)
		if !anyContains(problems, tc.want) {
			t.Errorf("%s: want %q in %q", tc.repl, tc.want, problems)
		}
		for _, p := range problems {
			if !strings.HasPrefix(p, "claims/refill.md:17: form: ") {
				t.Errorf("%s: the problem does not name the scenario's line: %q", tc.repl, p)
			}
		}
	}
}

// A Verifies that opens a subtest's body is a citation; one that another
// statement precedes is not.
func TestAVerifiesOpeningASubtestCounts(t *testing.T) {
	Verifies(t, claims.CH_023, "later-statement-refused")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	sub := "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n" +
		"func TestRefill(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_010, \"waited-request-passes\")\n\tt.Run(\"branch\", func(t *testing.T) {\n\t\tchartertest.Verifies(t, claims.SW_010, \"burst-caps-tokens\")\n\t})\n}\n"
	writeFile(t, root, "pkg/a/refill_test.go", sub)
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a citation opening a subtest reported:\n%s", strings.Join(problems, "\n"))
	}
	writeFile(t, root, "pkg/a/refill_test.go", strings.Replace(sub, "\t\tchartertest.Verifies(t, claims.SW_010, \"burst-caps-tokens\")", "\t\tt.Log(\"first\")\n\t\tchartertest.Verifies(t, claims.SW_010, \"burst-caps-tokens\")", 1))
	reports(t, root, "pkg/a/refill_test.go:14: rule 24: Verifies is the leading statement of its test")
}

// A subtest run from a range over the cases, or from another subtest, is
// a subtest the test runs: a Verifies opening it counts.
func TestAVerifiesOpeningASubtestInALoopCounts(t *testing.T) {
	Verifies(t, claims.CH_023, "later-statement-refused")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n"+
		"func TestRefill(t *testing.T) {\n\tfor _, name := range []string{\"a\", \"b\"} {\n\t\tt.Run(name, func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.SW_010, \"waited-request-passes\")\n\t\t})\n\t}\n"+
		"\tt.Run(\"outer\", func(t *testing.T) {\n\t\tt.Run(\"inner\", func(t *testing.T) {\n\t\t\tchartertest.Verifies(t, claims.SW_010, \"burst-caps-tokens\")\n\t\t})\n\t})\n}\n")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a citation opening a subtest in a loop or a subtest reported:\n%s", strings.Join(problems, "\n"))
	}
}
