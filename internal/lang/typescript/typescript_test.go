package typescript

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/gate"
	"github.com/leroux/charter/internal/record"
)

// repo is the charter's checkout, whose node_modules hold the compiler and
// vitest.
var repo, _ = filepath.Abs("../../..")

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
		rels = append(rels, rel)
	}
	slices.Sort(rels)
	set, err := Reader{Root: root, Compiler: filepath.Join(repo, "node_modules", "typescript"), Import: "@fx/claims", EnvFile: "web/env.ts", ClaimsDir: "claims", Table: "claims/charter.toml:5"}.Read(rels)
	if err != nil {
		t.Fatal(err)
	}
	return set
}

func refused(set *facts.Set) string {
	var ms []string
	for _, r := range set.Refused {
		ms = append(ms, r.Pos()+": "+r.Rule+": "+r.Message)
	}
	return strings.Join(ms, "\n")
}

const testFile = `import { describe, test } from "vitest";
import { claims, verifies } from "@fx/claims";

test("refills", () => {
  verifies(claims.SW_010, "waited-request-passes", "burst-caps-tokens");
});

test.skip("skipped", () => {
  verifies(claims.SW_010, "burst-caps-tokens");
});

test("later", () => {
  const x = 1;
  verifies(claims.SW_010);
});

test("by variable", () => {
  verifies(claims.SW_010, handle);
});

const handle = "x";
verifies(claims.SW_010);
`

func TestATypeScriptCitationOpensATest(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "typescript-citation-read", "disabled-ts-refused")
	set := read(t, map[string]string{"web/a.test.ts": testFile})
	if len(set.Citations) != 1 {
		t.Fatalf("citations = %+v\n%s", set.Citations, refused(set))
	}
	if c := set.Citations[0]; c.Pos() != "web/a.test.ts:5" || c.Symbol != "SW_010" || c.Test != "refills" || !slices.Equal(c.Handles, []string{"waited-request-passes", "burst-caps-tokens"}) {
		t.Errorf("citation = %+v", c)
	}
	got := refused(set)
	for _, want := range []string{
		"web/a.test.ts:9: rule 24: verifies sits in a test marked skip",
		"web/a.test.ts:14: rule 24: verifies is the leading statement of its test",
		"web/a.test.ts:18: rule 24: a scenario handle is a string literal",
		"web/a.test.ts:22: rule 24: verifies sits outside a test",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q in:\n%s", want, got)
		}
	}
}

func TestAVerifiesTheTestDeclaresIsRefused(t *testing.T) {
	claims.Verifies(t, claims.CH_023, "local-verifies-refused")
	set := read(t, map[string]string{"web/a.test.ts": `import { test } from "vitest";
import { claims } from "@fx/claims";

function verifies(..._: string[]) {}

test("refills", () => {
  verifies(claims.SW_010, "waited-request-passes");
});
`})
	if len(set.Citations) != 0 || !strings.Contains(refused(set), "web/a.test.ts:7: rule 24: verifies is declared in this file") {
		t.Errorf("citations %+v, refused:\n%s", set.Citations, refused(set))
	}
}

func TestATypeScriptDirectiveNamingAnUnknownIDFails(t *testing.T) {
	claims.Verifies(t, claims.CH_021, "typescript-anchor-resolves")
	root := t.TempDir()
	write(t, root, map[string]string{"claims/area.md": "# Area\n\n## SW-010 A rule\nStatus: active\n\nIt rules.\n\nWhy: someone needs it.\n\nScenarios:\n1. `rules`: it rules.\n"})
	rec, err := record.Read(root, "claims")
	if err != nil {
		t.Fatal(err)
	}
	set := read(t, map[string]string{"web/limit.ts": "// charter:realizes SW-010, SW-999\nexport function limit() {}\n"})
	if len(set.Anchors) != 1 || !slices.Equal(set.Anchors[0].IDs, []string{"SW-010", "SW-999"}) {
		t.Fatalf("anchors = %+v", set.Anchors)
	}
	ps := gate.Check(gate.Input{Record: rec, Facts: set})
	if !slices.Contains(ps, "web/limit.ts:1: rule 24: charter:realizes names SW-999, which is not an entry of the record") {
		t.Errorf("problems = %q", ps)
	}
}

func TestTypeScriptEnvironmentReadsAreFacts(t *testing.T) {
	set := read(t, map[string]string{
		"web/env.ts":  "export const RATE = \"FX_RATE\";\nconst hidden = \"FX_HIDDEN\";\n",
		"web/main.ts": "import * as env from \"./env\";\nimport { RATE } from \"./env.ts\";\n\nconst a = process.env[env.RATE];\nconst b = process.env[RATE];\nconst c = process.env.FX_RATE;\nconst d = process.env[name]; // charter:dynamic\nconst e = process.env[name];\nconst name = \"x\";\n",
	})
	if len(set.EnvConsts) != 1 || set.EnvConsts[0].Name != "env.RATE" || set.EnvConsts[0].Value != "FX_RATE" {
		t.Errorf("constants = %+v", set.EnvConsts)
	}
	r := set.EnvReads
	if len(r) != 5 || r[0].Const != "env.RATE" || r[1].Const != "env.RATE" || r[2].Literal != "FX_RATE" || !r[3].Dynamic || r[4].Const != "" || r[4].Literal != "" || r[4].Dynamic {
		t.Errorf("reads = %+v", r)
	}
}

func TestAMissingCompilerNamesTheTable(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_PQ7X, "parser-missing-named")
	root := t.TempDir()
	write(t, root, map[string]string{"a.ts": "export const a = 1;\n"})
	_, err := Reader{Root: root, Compiler: "node_modules/typescript", Table: "claims/charter.toml:5"}.Read([]string{"a.ts"})
	if err == nil || !strings.HasPrefix(err.Error(), "claims/charter.toml:5: ") || !strings.Contains(err.Error(), "no typescript package at node_modules/typescript") {
		t.Errorf("err = %v", err)
	}
}

const fixtureArea = "# Refill\n\n" +
	"## SW-010 The bucket refills at the rate */\nStatus: active\n\nAt each request the bucket reads the clock.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `waited-request-passes`: a request after a wait passes.\n2. `burst-caps-tokens`: a long wait earns no more.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\n\nNo longer.\n\nWhy: it was replaced.\n"

func generated(t *testing.T, root string) []byte {
	t.Helper()
	write(t, root, map[string]string{"claims/refill.md": fixtureArea})
	rec, err := record.Read(root, "claims")
	if err != nil || len(rec.Problems) != 0 {
		t.Fatalf("%v %q", err, rec.Problems)
	}
	src := Generate(rec, "charter gen")
	write(t, root, map[string]string{"claims/index.ts": string(src)})
	return src
}

func TestTheGeneratedModuleEqualsTheFixture(t *testing.T) {
	claims.Verifies(t, claims.CH_020, "typescript-module-rendered")
	got := generated(t, t.TempDir())
	golden := filepath.Join("testdata", "index.ts")
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
		t.Errorf("index.ts differs from %s", golden)
	}
	for _, w := range []string{"/** The bucket refills at the rate *\\/ */\n  SW_010: \"SW-010\",", "/** @deprecated SW-012 is retired"} {
		if !strings.Contains(string(got), w) {
			t.Errorf("index.ts lacks %q", w)
		}
	}
}

func TestATypeScriptTestRecordsItsOutcomeInTheEventsFile(t *testing.T) {
	claims.Verifies(t, claims.CH_024, "event-per-citation", "outcome-recorded", "unset-records-nothing")
	vitest := filepath.Join(repo, "claims-viewer", "node_modules", ".bin", "vitest")
	if _, err := os.Stat(vitest); err != nil {
		t.Fatalf("vitest is not installed at %s; run pnpm install", vitest)
	}
	root := t.TempDir()
	generated(t, root)
	write(t, root, map[string]string{"limit.test.ts": `import { verifies, claims } from "./claims/index.ts";

test("passes", () => {
  verifies(claims.SW_010, "waited-request-passes");
});

test("fails", () => {
  verifies(claims.SW_010, "burst-caps-tokens");
  expect(1).toBe(2);
});

test("skips", (ctx) => {
  verifies(claims.SW_010);
  ctx.skip();
});
`})
	events := filepath.Join(t.TempDir(), "events.jsonl")
	run := func(env ...string) string {
		cmd := exec.Command(vitest, "run", "--globals", "--root", root)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), env...)
		out, _ := cmd.CombinedOutput()
		return string(out)
	}
	out := run("CHARTER_EVENTS=" + events)
	f, err := os.Open(events)
	if err != nil {
		t.Fatalf("no events file: %v\n%s", err, out)
	}
	defer f.Close()
	got := map[string]string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		var l struct{ ID, Handle, Test, File, Outcome string }
		if err := json.Unmarshal(s.Bytes(), &l); err != nil {
			t.Fatalf("not a JSON line: %q", s.Text())
		}
		got[l.Test+"/"+l.Handle] = l.Outcome
		if l.ID != "SW-010" || !strings.HasSuffix(l.File, "limit.test.ts") {
			t.Errorf("line %+v", l)
		}
	}
	want := map[string]string{"passes/waited-request-passes": "passed", "fails/burst-caps-tokens": "failed", "skips/": "skipped"}
	if len(got) != len(want) {
		t.Errorf("outcomes = %v, want %v\n%s", got, want, out)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
	if err := os.Remove(events); err != nil {
		t.Fatal(err)
	}
	run()
	if _, err := os.Stat(events); err == nil {
		t.Error("with CHARTER_EVENTS unset the events file was written")
	}
}
