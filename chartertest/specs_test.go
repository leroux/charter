package chartertest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// reports runs the gate on root, with the generated files rewritten for
// the tree as make claims would, and fails unless a problem contains want.
func reports(t *testing.T, root, want string) {
	t.Helper()
	problems := check(t, root)
	if !anyContains(problems, want) {
		t.Fatalf("no problem contains %q; reported:\n%s", want, strings.Join(problems, "\n"))
	}
}

// check runs the gate on root with the generated Go and the TypeScript
// module rewritten for the tree first, as make claims would, so rules 15 and
// 22 check the test's charter rather than the committed fixture files.
func check(t *testing.T, root string) []string {
	t.Helper()
	t.Setenv(UpdateSpecsEnv, "1")
	t.Setenv(UpdateTypeScriptEnv, "1")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	return problems
}

const sampleSpec = `**Title**
Retries back off before any content arrives

**Lens**: SW

**Status**: active

**Accounts**: env:FX_TIMEOUT, op:GET /health

**Description**
A retryable failure is retried with a doubling delay.
A failure after content ends the call.

**Rationale**
Retrying after content would repeat what the person already saw.

**Verification Description**
A scripted endpoint fails twice and then answers; the call succeeds on the third attempt.

## Relations

**Realizes**

- [SYS-001](SYS-001-model-calls-are-reliable.md)

**Concerns**

- [CON-004](CON-004-nothing-is-sent-twice.md)

## Changes

- **2026-09-26** — Narrowed to failures before content.
`

func TestParseSpecReadsEveryField(t *testing.T) {
	Verifies(t, claims.CH_014, "file-form-read")
	e, problems := parseSpec("claims/specs/SW-014-retries-back-off-before-content.md", sampleSpec)
	if len(problems) != 0 {
		t.Fatalf("problems: %v", problems)
	}
	if e.ID != "SW-014" || e.Name != "SW_014" || e.Lens != "SW" || e.Number != 14 || e.Status != "active" || e.Form != formFile {
		t.Errorf("identity = %+v", e)
	}
	if e.Text != "Retries back off before any content arrives" || !strings.HasPrefix(e.Description, "A retryable failure") || strings.Count(e.Description, "\n") != 1 {
		t.Errorf("text = %q, description = %q", e.Text, e.Description)
	}
	if strings.Join(e.Accounts, ",") != "env:FX_TIMEOUT,op:GET /health" || strings.Join(e.Realizes, ",") != "SYS-001" || strings.Join(e.Concerns, ",") != "CON-004" {
		t.Errorf("accounts %v, realizes %v, concerns %v", e.Accounts, e.Realizes, e.Concerns)
	}
	if e.Links["SYS-001"] != "SYS-001-model-calls-are-reliable.md" {
		t.Errorf("links = %v", e.Links)
	}
}

func TestParseSpecForm(t *testing.T) {
	Verifies(t, claims.CH_014, "file-form-checked")
	for _, tc := range []struct{ name, old, repl, want string }{
		{"out of order", "**Status**: active\n\n", "", "Status out of order"},
		{"unknown field", "**Rationale**", "**Why**", `"Why" is not a field of a spec`},
		{"inline block", "**Description**\n", "**Description**: x\n", "write **Description** on its own line"},
		{"relation label", "**Concerns**", "**Supports**", `relation "Supports" is not Realizes, Concerns or Related`},
		{"link text", "- [CON-004](CON-004-nothing-is-sent-twice.md)", "- [the constraint](CON-004-nothing-is-sent-twice.md)", `link text is the target's id alone, not "the constraint"`},
		{"change row", "- **2026-09-26** — Narrowed", "- 2026-09-26: Narrowed", "a change opens - **YYYY-MM-DD** —"},
		{"text outside a field", "**Title**", "Stray words.\n\n**Title**", "text outside a field"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(sampleSpec, tc.old, tc.repl, 1)
			if src == sampleSpec {
				t.Fatalf("the edit %q did not apply", tc.old)
			}
			_, problems := parseSpec("claims/specs/SW-014-retries-back-off-before-content.md", src)
			if tc.name == "out of order" {
				tc.want = "Accounts is out of order"
				src2 := strings.Replace(sampleSpec, "**Accounts**: env:FX_TIMEOUT, op:GET /health\n\n", "", 1)
				src2 = strings.Replace(src2, "**Rationale**", "**Accounts**: env:FX_TIMEOUT\n\n**Rationale**", 1)
				_, problems = parseSpec("claims/specs/SW-014-retries-back-off-before-content.md", src2)
			}
			if !anyContains(problems, tc.want) {
				t.Errorf("want %q in %v", tc.want, problems)
			}
		})
	}
}

func TestSpecNames(t *testing.T) {
	Verifies(t, claims.CH_014, "file-name-checked")
	for _, tc := range []struct {
		base string
		id   string
		ok   bool
	}{
		{"SW-004-live-policy.md", "SW-004", true},
		{"SW-261001-K7Q2-dated.md", "", false},
		{"SYS-012-x.md", "SYS-012", true},
		{"SW-000-zero.md", "", false},
		{"XX-004-unknown-lens.md", "", false},
		{"SW-004-Upper.md", "", false},
		{"SW-TMP-004-draft.md", "", false},
	} {
		id, _, _, ok := specName(tc.base)
		if ok != tc.ok || id != tc.id {
			t.Errorf("specName(%q) = %q, %v; want %q, %v", tc.base, id, ok, tc.id, tc.ok)
		}
	}
}

// The generated Go declares every spec of any status with a bare symbol, a
// dated id with each hyphen an underscore, a planned one marked planned and
// a deprecated one deprecated. The order they come in is CH-012's.
func TestTheGeneratedGoDeclaresEverySpec(t *testing.T) {
	Verifies(t, claims.CH_020, "bare-symbol-declared", "planned-and-deprecated-marked")
	Verifies(t, claims.CH_003, "retired-entry-still-builds")
	root := copyFixture(t)
	data, err := os.ReadFile(filepath.Join(root, "claims", "specs.go"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(data)
	for _, want := range []string{
		"// Code generated from claims by the charter gate. DO NOT EDIT.",
		`var SW_001 = &charter.Spec{ID: "SW-001", Title: "A bucket survives a restart", Status: charter.Active}`,
		"// Deprecated: SW-002 is retired",
		"// Planned: SW-003 is ruled and not built",
		`var SW_003 = &charter.Spec{ID: "SW-003", Title: "A person can clear a quota's counters", Status: charter.Planned}`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("specs.go lacks %q", want)
		}
	}
	if strings.Contains(src, "Lens:") {
		t.Error("a lens was generated")
	}

	pkg := newEntriesPackage()
	writeFile(t, root, "claims/dated.md", "# Dated\n\n"+plannedEntry("SW-261001-K7Q2")+plannedEntry("SW-010"))
	if _, problems, err := readAreas(root, "claims", pkg); err != nil || len(problems) != 0 {
		t.Fatalf("readAreas: %v, problems %v", err, problems)
	}
	dated, err := renderSpecsGo(pkg, "claims", "claims", "make claims")
	if err != nil {
		t.Fatal(err)
	}
	for _, decl := range []string{
		`var SW_010 = &charter.Spec{ID: "SW-010", Title: "An entry", Status: charter.Planned}`,
		`var SW_261001_K7Q2 = &charter.Spec{ID: "SW-261001-K7Q2", Title: "An entry", Status: charter.Planned}`,
	} {
		if !strings.Contains(string(dated), decl) {
			t.Errorf("the generated Go lacks %s:\n%s", decl, dated)
		}
	}
}

// A deprecated entry's declaration carries a deprecation notice naming the
// document it is retired in.
func TestADeprecatedEntryIsDeclaredWithANotice(t *testing.T) {
	Verifies(t, claims.CH_020, "planned-and-deprecated-marked")
	root := copyFixture(t)
	writeFile(t, root, "claims/gone.md", "# Gone\n\n"+strings.Replace(plannedEntry("SW-010"), "Status: planned", "Status: deprecated", 1))
	pkg := newEntriesPackage()
	if _, problems, err := readAreas(root, "claims", pkg); err != nil || len(problems) != 0 {
		t.Fatalf("readAreas: %v, problems %v", err, problems)
	}
	if _, err := readSpecs(root, "claims/specs", pkg); err != nil {
		t.Fatal(err)
	}
	src, err := renderSpecsGo(pkg, "claims", "claims", "make claims")
	if err != nil {
		t.Fatal(err)
	}
	for _, notice := range []string{
		"// Deprecated: SW-010 is retired; see its Changes in claims/gone.md. Add no new anchors.\n" +
			`var SW_010 = &charter.Spec{ID: "SW-010", Title: "An entry", Status: charter.Deprecated}`,
		"// Deprecated: SW-002 is retired; see its Changes in claims/specs/SW-002-limiter-runs-on-windows.md. Add no new anchors.\n" +
			`var SW_002 = &charter.Spec{ID: "SW-002", Title: "The limiter runs on Windows", Status: charter.Deprecated}`,
	} {
		if !strings.Contains(string(src), notice) {
			t.Errorf("the generated Go lacks the notice:\n%s\nin:\n%s", notice, src)
		}
	}
}

// A directive on a package clause is the package's anchor, recorded with
// its kind, file, line and declaration.
func TestAPackageDirectiveIsThePackagesAnchor(t *testing.T) {
	Verifies(t, claims.CH_021, "go-anchor-recorded")
	reg, err := Read(fixtureOptions(copyFixture(t)))
	if err != nil {
		t.Fatal(err)
	}
	want := Anchor{Kind: "realizes", File: "pkg/a/a.go", Line: 5, Decl: "package a"}
	for _, name := range []string{"SW_001", "CON_001", "CON_002"} {
		if e := reg.Entry(name); e == nil || len(e.Anchors) != 1 || e.Anchors[0] != want {
			t.Errorf("%s anchors = %+v, want %+v", name, e, want)
		}
	}
}

func TestReadListsSpecsWithRelationsAnchorsAndCitations(t *testing.T) {
	Verifies(t, claims.CH_040, "registry-lists-everything")
	Verifies(t, claims.CH_002, "anchors-and-citations-listed")
	root := copyFixture(t)
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	e := reg.Entry("SW_001")
	if e == nil || e.ID != "SW-001" || e.Lens != "SW" || len(e.Realizes) != 1 || len(e.Anchors) != 1 || e.Anchors[0].Decl != "package a" {
		t.Fatalf("SW-001 = %+v", e)
	}
	if len(e.Refs) != 2 || e.Refs[0].Line != 11 || e.Refs[1].File != "web/a.test.ts" {
		t.Errorf("SW-001 refs = %v, want the Go test and the TypeScript test", e.Refs)
	}
	if p := reg.Entry("SW_003"); p == nil || p.Status != "planned" {
		t.Errorf("the planned spec is missing from the registry: %+v", p)
	}
}

// A test may cite a planned spec: the symbol is generated, and coverage asks
// nothing of a planned spec.
func TestAPlannedSpecMayBeCited(t *testing.T) {
	Verifies(t, claims.CH_013, "planned-cited-freely")
	Verifies(t, claims.CH_002, "planned-entry-cited-first")
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "chartertest.Verifies(t, claims.SW_002)", "chartertest.Verifies(t, claims.SW_002)\n\tchartertest.Verifies(t, claims.SW_003)")
	if problems := check(t, root); len(problems) != 0 {
		t.Errorf("citing a planned spec reported:\n%s", strings.Join(problems, "\n"))
	}
}
