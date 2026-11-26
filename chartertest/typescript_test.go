package chartertest

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestTypeScriptFilePredicates(t *testing.T) {
	cases := []struct {
		rel        string
		file, test bool
	}{
		{"web/a.test.ts", true, true},
		{"web/b.spec.ts", true, true},
		{"web/c.test.tsx", true, true},
		{"web/d.spec.tsx", true, true},
		{"web/helper.ts", true, false},
		{"web/App.tsx", true, false},
		{"web/types.d.ts", false, false},
		{"web/setup.mjs", false, false},
		{"pkg/a/a_test.go", false, false},
	}
	for _, c := range cases {
		if got := isTSFile(c.rel); got != c.file {
			t.Errorf("isTSFile(%q) = %v, want %v", c.rel, got, c.file)
		}
		if got := isTSTest(c.rel); got != c.test {
			t.Errorf("isTSTest(%q) = %v, want %v", c.rel, got, c.test)
		}
	}
}

func TestWalkListsTypeScriptFilesWhenConfigured(t *testing.T) {
	Verifies(t, claims.CH_036, "typescript-listed")
	root := copyFixture(t)
	writeFile(t, root, "web/types.d.ts", "export type Nothing = never;\n")
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	tr, err := loadTree(r)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"claims/index.ts", "web/a.test.ts", "web/b.spec.ts", "web/helper.ts"}
	if !reflect.DeepEqual(tr.ts, want) {
		t.Errorf("ts = %v, want %v", tr.ts, want)
	}

	opts := fixtureOptions(root)
	opts.TypeScript = nil
	r, err = opts.resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	tr, err = loadTree(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.ts) != 0 {
		t.Errorf("ts = %v, want none without Options.TypeScript", tr.ts)
	}
}

func TestResolveTypeScriptDefaultsAndValidates(t *testing.T) {
	Verifies(t, claims.CH_037, "typescript-options-validated")
	root := copyFixture(t)
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	if r.TypeScript.Node != "node" || !strings.HasSuffix(r.TypeScript.Compiler, "typescript") {
		t.Errorf("TypeScript = %+v", r.TypeScript)
	}
	opts := fixtureOptions(root)
	opts.TypeScript = &TypeScript{Module: "claims/index.ts", Import: "@fixture/claims", Compiler: "web/node_modules/typescript"}
	r, err = opts.resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	if r.TypeScript.Compiler != root+"/web/node_modules/typescript" {
		t.Errorf("relative Compiler resolved to %q", r.TypeScript.Compiler)
	}
	for _, missing := range []TypeScript{
		{Import: "@fixture/claims", Compiler: "x"},
		{Module: "claims/index.ts", Compiler: "x"},
		{Module: "claims/index.ts", Import: "@fixture/claims"},
	} {
		opts.TypeScript = &missing
		if _, err := opts.resolve("unused"); err == nil || !strings.Contains(err.Error(), "Options.TypeScript.") {
			t.Errorf("resolve(%+v) = %v, want an error naming the missing field", missing, err)
		}
	}
}

func TestExtractorReportsImportsAndReferences(t *testing.T) {
	root := copyFixture(t)
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	got, err := extractTS(r, []string{"web/a.test.ts", "web/b.spec.ts", "web/helper.ts"})
	if err != nil {
		t.Fatal(err)
	}
	want := []tsFile{
		{File: "web/a.test.ts", Diagnostics: []tsDiagnostic{},
			Imports: []tsImport{{Line: 2, Form: "named", Local: "c", Imported: "claims"}},
			Refs: []tsRef{
				{Line: 7, Name: "SW_001", Test: "survives a restart", InTest: true, Call: true, Handles: []*string{}, Leading: true},
				{Line: 10, Name: "CON_002", Test: "folds %d", InTest: true, Call: true, Handles: []*string{}, Leading: true},
				{Line: 13, Name: "CON_001", Test: `locks ${"maybe"}`, InTest: true, Call: true, Handles: []*string{}, Leading: true},
			}},
		{File: "web/b.spec.ts", Diagnostics: []tsDiagnostic{},
			Imports: []tsImport{{Line: 2, Form: "named", Local: "claims", Imported: "claims"}},
			Refs:    []tsRef{{Line: 6, Name: "SW_002", Test: "windows is not supported yet", InTest: true, Call: true, Handles: []*string{}, Leading: true}}},
		{File: "web/helper.ts", Diagnostics: []tsDiagnostic{},
			Imports: []tsImport{{Line: 1, Form: "type"}},
			Refs:    []tsRef{}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("extractTS =\n%+v\nwant\n%+v", got, want)
	}
}

func TestExtractorReportsEveryImportForm(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "web/forms.test.ts", strings.Join([]string{
		`import * as ns from "@fixture/claims";`,
		`import dflt from "@fixture/claims";`,
		`import "@fixture/claims";`,
		`import { type Claim } from "@fixture/claims";`,
		`import { claims, other } from "@fixture/claims";`,
		`import { claims } from "elsewhere";`,
		`export const x = [ns, dflt, claims.CON_002, other];`,
		``,
	}, "\n"))
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	got, err := extractTS(r, []string{"web/forms.test.ts"})
	if err != nil {
		t.Fatal(err)
	}
	wantImports := []tsImport{
		{Line: 1, Form: "namespace", Local: "ns"},
		{Line: 2, Form: "default", Local: "dflt"},
		{Line: 3, Form: "side-effect"},
		{Line: 4, Form: "type"},
		{Line: 5, Form: "named", Local: "claims", Imported: "claims"},
		{Line: 5, Form: "named", Local: "other", Imported: "other"},
	}
	if !reflect.DeepEqual(got[0].Imports, wantImports) {
		t.Errorf("Imports = %+v, want %+v", got[0].Imports, wantImports)
	}
	wantRefs := []tsRef{{Line: 7, Name: "CON_002", Test: "", Handles: []*string{}}}
	if !reflect.DeepEqual(got[0].Refs, wantRefs) {
		t.Errorf("Refs = %+v, want %+v", got[0].Refs, wantRefs)
	}
}

func TestExtractorFindsTheEnclosingTestThroughStepsAndDescribes(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "web/nest.spec.ts", strings.Join([]string{
		`import { test } from "@example/test";`,
		`import { claims } from "@fixture/claims";`,
		`test.describe("outer", () => {`,
		`  test.beforeEach(() => { void claims.CON_001; });`,
		`  test.only("inner", async () => {`,
		`    await test.step("a step", async () => { void claims.SW_002; });`,
		`  });`,
		`  test(42 as unknown as string, () => { void claims.CON_002; });`,
		`});`,
		``,
	}, "\n"))
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	got, err := extractTS(r, []string{"web/nest.spec.ts"})
	if err != nil {
		t.Fatal(err)
	}
	want := []tsRef{
		{Line: 4, Name: "CON_001", Test: "", Handles: []*string{}},
		{Line: 6, Name: "SW_002", Test: "inner", InTest: true, Handles: []*string{}},
		{Line: 8, Name: "CON_002", Test: "", InTest: true, Handles: []*string{}},
	}
	if !reflect.DeepEqual(got[0].Refs, want) {
		t.Errorf("Refs = %+v, want %+v", got[0].Refs, want)
	}
}

func TestExtractorFailsOnAParseError(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "web/broken.test.ts", "const = ;\n")
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	_, err = extractTS(r, []string{"web/broken.test.ts"})
	if err == nil || !strings.Contains(err.Error(), "parse web/broken.test.ts:1") {
		t.Errorf("err = %v, want a parse error at web/broken.test.ts:1", err)
	}
}

func TestExtractorNeedsNodeAndTheCompiler(t *testing.T) {
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.TypeScript.Node = "chartertest-no-such-node"
	r, err := opts.resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := extractTS(r, nil); err == nil || !strings.Contains(err.Error(), "not found on PATH") {
		t.Errorf("err = %v, want node not found", err)
	}
	opts = fixtureOptions(root)
	opts.TypeScript.Compiler = t.TempDir()
	r, err = opts.resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := extractTS(r, nil); err == nil || !strings.Contains(err.Error(), "no TypeScript compiler at") {
		t.Errorf("err = %v, want compiler missing", err)
	}
}

// Every TypeScript rule has a mutation of the clean fixture that trips it,
// and a citation from TypeScript counts under rules 4 and 5.
func TestTypeScriptRules(t *testing.T) {
	cases := []struct {
		name string
		edit edit
		want string
	}{
		{"rule 16 namespace import",
			edit{rel: "web/a.test.ts", old: `import { claims as c } from "@fixture/claims";`, repl: `import * as c from "@fixture/claims";`},
			`web/a.test.ts:2: rule 16: import the module as { claims } from "@fixture/claims"`},
		{"rule 16 default import",
			edit{rel: "web/b.spec.ts", old: `import { claims } from "@fixture/claims";`, repl: `import claims from "@fixture/claims";`},
			`web/b.spec.ts:2: rule 16: import the module as { claims }`},
		{"rule 16 side-effect import",
			edit{rel: "web/b.spec.ts", old: `import { claims } from "@fixture/claims";`, repl: `import "@fixture/claims";`},
			`web/b.spec.ts:2: rule 16: import the module as { claims }`},
		{"rule 16 another named export",
			edit{rel: "web/b.spec.ts", old: `import { claims } from "@fixture/claims";`, repl: `import { claims, other } from "@fixture/claims";`},
			`web/b.spec.ts:2: rule 16: import the module as { claims }`},
		{"rule 16 value import outside a test",
			edit{rel: "web/prod.ts", write: "import { claims } from \"@fixture/claims\";\n\nexport const id = claims.CON_002;\n"},
			"web/prod.ts:1: rule 16: @fixture/claims is imported from test files only"},
		{"rule 17 unknown name",
			edit{rel: "web/a.test.ts", old: "c.CON_002", repl: "c.CON_0020"},
			"web/a.test.ts:10: rule 17: CON_0020 is not an entry"},
		{"rule 17 observation",
			edit{rel: "web/a.test.ts", old: "c.CON_001", repl: "c.Suspect"},
			"web/a.test.ts:13: rule 17: Suspect is an observation; a TypeScript test cites specs"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := copyFixture(t)
			apply(t, root, c.edit)
			problems, err := Problems(fixtureOptions(root))
			if err != nil {
				t.Fatal(err)
			}
			if !anyContains(problems, c.want) {
				t.Errorf("problems = %v\nwant one containing %q", problems, c.want)
			}
		})
	}
}

func TestAKeepEntryCitedOnlyFromTypeScriptPassesRule4(t *testing.T) {
	root := copyFixture(t)
	replaceIn(t, root, "pkg/a/a_test.go", "\tchartertest.Verifies(t, claims.CON_002)\n", "")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if anyContains(problems, "rule 25: active spec CON-002 is verified by no test") {
		t.Errorf("CON_002 is cited from web/a.test.ts, yet: %v", problems)
	}
}

func TestTypeScriptReferencesAreOffWithoutTheOption(t *testing.T) {
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.TypeScript = nil
	refs := refsOf(t, opts)
	for name, rs := range refs {
		for _, r := range rs {
			if strings.HasPrefix(r.File, "web/") {
				t.Errorf("%s is cited from %s without Options.TypeScript", name, r.File)
			}
		}
	}
}

func TestATypeScriptRefIsCountedOncePerTest(t *testing.T) {
	root := copyFixture(t)
	replaceIn(t, root, "web/a.test.ts", "verifies(c.SW_001);", "verifies(c.SW_001); verifies(c.SW_001);")
	refs := refsOf(t, fixtureOptions(root))
	if got := refs["SW_001"]; len(got) != 2 {
		t.Errorf("refs = %v, want the Go ref and one TypeScript ref", got)
	}
}

func TestExtractorWithNoFilesReportsNothing(t *testing.T) {
	root := copyFixture(t)
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	got, err := extractTS(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("extractTS(nil) = %+v, want nothing", got)
	}
}

func TestRenderTypeScriptMatchesTheFixtureModule(t *testing.T) {
	Verifies(t, claims.CH_020, "typescript-module-rendered")
	root := copyFixture(t)
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	s, err := load(r)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(root, "claims", "index.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if got := renderTypeScript(s.pkg, r.entriesDir, r.regenerate(UpdateTypeScriptEnv)); got != string(want) {
		t.Errorf("renderTypeScript =\n%s\nwant\n%s", got, want)
	}
}

func TestRenderTypeScriptEscapesACommentCloser(t *testing.T) {
	Verifies(t, claims.CH_020, "typescript-module-rendered")
	root := copyFixture(t)
	replaceIn(t, root, "claims/specs/CON-002-fold-never-panics.md", "The fold never panics\n", "The fold never panics */ or does it\n")
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	s, err := load(r)
	if err != nil {
		t.Fatal(err)
	}
	got := renderTypeScript(s.pkg, r.entriesDir, r.regenerate(UpdateTypeScriptEnv))
	if !strings.Contains(got, `/** The fold never panics *\/ or does it */`) {
		t.Errorf("rendered:\n%s", got)
	}
}

func TestRule15(t *testing.T) {
	Verifies(t, claims.CH_020, "stale-generated-fails")
	cases := []struct {
		name string
		edit edit
		want string
	}{
		{"stale module",
			edit{rel: "claims/index.ts", add: "// stale\n"},
			"claims/index.ts: rule 15: the TypeScript module differs from the registry; run CHARTER_UPDATE_TYPESCRIPT=1 go test -count=1 ./claims"},
		{"missing module",
			edit{rel: "claims/index.ts", rm: true},
			"claims/index.ts: rule 15: the TypeScript module is missing; run CHARTER_UPDATE_TYPESCRIPT=1 go test -count=1 ./claims"},
		{"entry added",
			edit{rel: "claims/specs/CON-002-fold-never-panics.md", old: "The fold never panics\n", repl: "The fold never panics at all\n"},
			"claims/index.ts: rule 15: the TypeScript module differs from the registry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := copyFixture(t)
			apply(t, root, c.edit)
			problems, err := Problems(fixtureOptions(root))
			if err != nil {
				t.Fatal(err)
			}
			if !anyContains(problems, c.want) {
				t.Errorf("problems = %v\nwant one containing %q", problems, c.want)
			}
		})
	}
}

func TestUpdateWritesTheTypeScriptModuleAndIsAFixedPoint(t *testing.T) {
	Verifies(t, claims.CH_020, "update-writes")
	root := copyFixture(t)
	appendTo(t, root, "claims/index.ts", "// stale\n")
	t.Setenv(UpdateTypeScriptEnv, "1")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("update reported: %v", problems)
	}
	t.Setenv(UpdateTypeScriptEnv, "")
	problems, err = Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("after update: %v", problems)
	}
	got, err := os.ReadFile(filepath.Join(root, "claims", "index.ts"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "fixture", "claims", "index.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("update wrote:\n%s\nwant the committed fixture module", got)
	}
}

func TestReadDoesNotApplyRule15(t *testing.T) {
	root := copyFixture(t)
	appendTo(t, root, "claims/index.ts", "// stale\n")
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if anyContains(reg.Problems, "rule 15") {
		t.Errorf("Read reported rule 15: %v", reg.Problems)
	}
}
