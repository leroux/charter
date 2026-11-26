package record

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

func anyContains(list []string, want string) bool {
	for _, s := range list {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}

func TestParseAreaReadsEveryPart(t *testing.T) {
	claims.Verifies(t, claims.CH_010, "every-part-read")
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
	if e.ID != "SW-010" || e.Symbol != "SW_010" || e.Need || e.Prefix != "SW" || e.Number != 10 || e.Pos() != "claims/refill.md:4" {
		t.Errorf("identity = %+v", e)
	}
	if e.Title != "The bucket refills from its last request at the rate" || e.Status != "active" || e.Kind != "invariant" {
		t.Errorf("title %q, status %q, kind %q", e.Title, e.Status, e.Kind)
	}
	if strings.Join(e.Serves, ",") != "STK-001" || strings.Join(e.See, ",") != "SW-001" || strings.Join(e.Accounts, ",") != "env:FX_DATA_DIR" {
		t.Errorf("serves %v, see %v, accounts %v", e.Serves, e.See, e.Accounts)
	}
	if !strings.HasPrefix(e.Text, "At each request") || strings.Count(e.Text, "\n") != 1 {
		t.Errorf("description = %q", e.Text)
	}
	if !strings.HasPrefix(e.Why, "a caller that waits") || !strings.HasPrefix(e.Guidance, "no timer") {
		t.Errorf("why = %q, guidance = %q", e.Why, e.Guidance)
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
	claims.Verifies(t, claims.CH_010, "heading-form", "fields-checked", "required-parts", "scenario-form", "parts-in-order", "change-form", "problems-locate")
	claims.Verifies(t, claims.CH_002, "problems-name-file-line-and-fix")
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
	claims.Verifies(t, claims.CH_013, "status-form")
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
	claims.Verifies(t, claims.CH_261126_JC9B, "header-is-prose")
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
	claims.Verifies(t, claims.CH_261126_JC9B, "heading-in-fence-is-text")
	for _, fenced := range []string{
		"```sh\n# a comment\n## SW-011 not an entry\nmake claims\n```\n",
		"~~~~\n# a comment\n~~~\n## still inside\n~~~~\n",
	} {
		src := strings.Replace(sampleArea, "\nWhy: a caller", "\n"+fenced+"\nWhy: a caller", 1)
		_, entries, problems := parseArea("claims/refill.md", src)
		if len(problems) != 0 {
			t.Errorf("%q: problems: %v", fenced, problems)
		}
		if len(entries) != 1 || !strings.Contains(entries[0].Text, "# a comment") || entries[0].Why == "" {
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
	claims.Verifies(t, claims.CH_011, "need-by-document", "need-fields-refused")
	src := strings.Replace(sampleArea, "## SW-010", "## STK-010", 1)
	_, entries, problems := parseArea("claims/needs.md", src)
	if len(entries) != 1 || !entries[0].Need {
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

func TestAScenarioCarriesWhoRuledIt(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_VY5C, "owner-read", "delegate-read", "default-is-provisional", "malformed-refused")
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

// An id is numbered or dated; a dated id's day is on the calendar and its
// code is four characters of the alphabet.
func TestAnIDIsNumberedOrDated(t *testing.T) {
	claims.Verifies(t, claims.CH_012, "id-form")
	for id, want := range map[string]Entry{
		"SW-010":         {Symbol: "SW_010", Prefix: "SW", Number: 10},
		"SW-261001-K7Q2": {Symbol: "SW_261001_K7Q2", Prefix: "SW", Day: "261001", Code: "K7Q2"},
		"CH-280229-ZY98": {Symbol: "CH_280229_ZY98", Prefix: "CH", Day: "280229", Code: "ZY98"},
	} {
		_, entries, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## "+id, 1))
		if len(problems) != 0 || len(entries) != 1 {
			t.Errorf("%s: %d entries, problems %v", id, len(entries), problems)
			continue
		}
		e := entries[0]
		if e.ID != id || e.Symbol != want.Symbol || e.Prefix != want.Prefix || e.Number != want.Number || e.Day != want.Day || e.Code != want.Code {
			t.Errorf("%s parsed as %q: symbol %s, prefix %s, number %d, day %q, code %q", id, e.ID, e.Symbol, e.Prefix, e.Number, e.Day, e.Code)
		}
	}
	for id, want := range map[string]string{
		"SW-261332-K7Q2":  "the date of SW-261332-K7Q2, 261332, is not a day on the calendar written YYMMDD",
		"SW-260229-K7Q2":  "the date of SW-260229-K7Q2, 260229, is not a day on the calendar",
		"SW-2610-K7Q2":    "the date of SW-2610-K7Q2, 2610, is not a day on the calendar",
		"SW-261001-K7Q":   "the code of SW-261001-K7Q, K7Q, is not 4 characters of BCDFGHJKMNPQRSTVWXYZ23456789",
		"SW-261001-K7Q22": "the code of SW-261001-K7Q22, K7Q22, is not 4 characters",
		"SW-261001-KAQ2":  "the code of SW-261001-KAQ2, KAQ2, is not 4 characters",
		"SW-261001-KLQ2":  "the code of SW-261001-KLQ2, KLQ2, is not 4 characters",
		"SW-261001-K0Q2":  "the code of SW-261001-K0Q2, K0Q2, is not 4 characters",
		"SW-261001-K1Q2":  "the code of SW-261001-K1Q2, K1Q2, is not 4 characters",
		"SW-261001-k7q2":  "the code of SW-261001-k7q2, k7q2, is not 4 characters",
		"SW-DRAFT":        `"SW-DRAFT" is not an id: an id is a prefix and a number, SW-118, or a prefix, the day it was taken and a code, SW-261001-K7Q2`,
	} {
		_, entries, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## "+id, 1))
		if len(entries) != 0 {
			t.Errorf("%s yielded an entry", id)
		}
		if !anyContains(problems, "claims/refill.md:4: form: "+want) {
			t.Errorf("%s: want %q in:\n%s", id, want, strings.Join(problems, "\n"))
		}
	}
}

// A dated id is anchored, cited from Go and from TypeScript and waived as a
