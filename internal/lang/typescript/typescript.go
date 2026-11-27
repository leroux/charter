// Package typescript is the charter's TypeScript language: a reader that
// runs node over the project's own typescript compiler to turn TypeScript
// source into facts, and a generator that writes claims/index.ts, whose
// verifies and observes are the marker a test calls.
package typescript

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
)

//go:embed reader.mjs
var readerScript string

// Node is the node binary the reader runs.
var Node = "node"

// Reader reads TypeScript files into facts.
type Reader struct {
	Root string
	// Compiler is the typescript package directory, relative to the root.
	Compiler string
	// Import is the specifier tests import the generated module by.
	Import string
	// EnvFile is the constants module, slash-relative to the root, or "".
	EnvFile string
	// ClaimsDir is the claims directory, whose generated module is not read.
	ClaimsDir string
	// Table names the configuration's TypeScript table in a problem.
	Table string
}

// ErrNoParser is wrapped when node or the compiler is missing.
var ErrNoParser = errors.New("the TypeScript parser cannot run")

// IsTest reports whether a TypeScript file is a test.
func IsTest(rel string) bool {
	for _, s := range []string{".test.ts", ".test.tsx", ".spec.ts", ".spec.tsx"} {
		if strings.HasSuffix(rel, s) {
			return true
		}
	}
	return false
}

type record_ struct {
	File        string `json:"file"`
	Diagnostics []struct {
		Line    int    `json:"line"`
		Message string `json:"message"`
	} `json:"diagnostics"`
	Imports []struct {
		Line     int    `json:"line"`
		Form     string `json:"form"`
		Local    string `json:"local"`
		Imported string `json:"imported"`
	} `json:"imports"`
	Refs []struct {
		Line        int       `json:"line"`
		Name        string    `json:"name"`
		Test        string    `json:"test"`
		InTest      bool      `json:"inTest"`
		Call        bool      `json:"call"`
		Handles     []*string `json:"handles"`
		Disabled    bool      `json:"disabled"`
		LocalCallee bool      `json:"localCallee"`
		Callee      string    `json:"callee"`
		Leading     bool      `json:"leading"`
	} `json:"refs"`
	Consts []struct {
		Line  int    `json:"line"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"consts"`
	EnvReads []struct {
		Line    int    `json:"line"`
		Const   string `json:"const"`
		Literal string `json:"literal"`
		Dynamic bool   `json:"dynamic"`
	} `json:"envReads"`
	Literals []struct {
		Line  int    `json:"line"`
		Value string `json:"value"`
	} `json:"literals"`
	Decls []struct {
		Name string `json:"name"`
		Line int    `json:"line"`
		End  int    `json:"end"`
	} `json:"decls"`
}

var directive = regexp.MustCompile(`^\s*//\s*charter:(realizes|concerns)\b(.*)$`)

// Read runs the reader over the files, and decides from its records which
// citations count and which are refused, and which environment reads go
// through a constant.
//
//charter:realizes CH-261126-PQ7X, CH-023, CH-034, CH-261127-DW8R
func (r Reader) Read(files []string) (*facts.Set, error) {
	var keep []string
	for _, f := range files {
		if r.ClaimsDir == "" || !strings.HasPrefix(f, r.ClaimsDir+"/") {
			keep = append(keep, f)
		}
	}
	set := &facts.Set{}
	if len(keep) == 0 {
		return set, nil
	}
	recs, err := r.extract(keep)
	if err != nil {
		return nil, err
	}
	for _, f := range recs {
		for _, d := range f.Diagnostics {
			return nil, fmt.Errorf("parse %s:%d: %s", f.File, d.Line, d.Message)
		}
		at := func(line int) facts.At { return facts.At{File: f.File, Line: line} }
		refuse := func(line int, rule, format string, args ...any) {
			set.Refused = append(set.Refused, facts.Refused{At: at(line), Rule: rule, Message: fmt.Sprintf(format, args...)})
		}
		test := IsTest(f.File)
		for _, imp := range f.Imports {
			switch {
			case imp.Form == "type":
			case !test:
				refuse(imp.Line, "rule 16", "%s is imported from test files only", r.Import)
			case imp.Form != "named" || (imp.Imported != "claims" && imp.Imported != "verifies" && imp.Imported != "observes"):
				refuse(imp.Line, "rule 16", "import { claims, verifies } from %q", r.Import)
			}
		}
		if !test {
			if err := r.anchors(f.File, set); err != nil {
				return nil, err
			}
			for _, k := range f.Consts {
				set.EnvConsts = append(set.EnvConsts, facts.EnvConst{At: at(k.Line), Name: k.Name, Value: k.Value})
			}
			for _, e := range f.EnvReads {
				set.EnvReads = append(set.EnvReads, facts.EnvRead{At: at(e.Line), Const: e.Const, Literal: e.Literal, Dynamic: e.Dynamic})
			}
			for _, l := range f.Literals {
				set.Literals = append(set.Literals, facts.Literal{At: at(l.Line), Value: l.Value})
			}
			pkg := strings.TrimSuffix(path.Base(f.File), path.Ext(f.File))
			for _, d := range f.Decls {
				set.Decls = append(set.Decls, facts.Decl{At: at(d.Line), End: d.End, Name: d.Name, Package: pkg})
			}
			continue
		}
		set.TestFiles = append(set.TestFiles, f.File)
		for _, ref := range f.Refs {
			if !ref.Call {
				set.References = append(set.References, facts.Reference{At: at(ref.Line), Symbol: ref.Name})
				continue
			}
			switch {
			case !ref.InTest:
				refuse(ref.Line, "rule 24", "verifies sits outside a test; it is the first statement of the test that verifies the spec")
				continue
			case ref.Disabled:
				refuse(ref.Line, "rule 24", "verifies sits in a test marked skip, fixme, fail or todo, which never runs as a test")
				continue
			case ref.LocalCallee:
				refuse(ref.Line, "rule 24", "verifies is declared in this file; import it with the claims from %q", r.Import)
				continue
			case !ref.Leading:
				refuse(ref.Line, "rule 24", "verifies is the leading statement of its test, before anything runs, skips or branches; here something precedes or encloses it")
				continue
			}
			c := facts.Citation{At: at(ref.Line), Symbol: ref.Name, Test: ref.Test, Observes: ref.Callee == "observes"}
			ok := true
			for _, h := range ref.Handles {
				if h == nil {
					refuse(ref.Line, "rule 24", "a scenario handle is a string literal")
					ok = false
					continue
				}
				c.Handles = append(c.Handles, *h)
			}
			if ok {
				set.Citations = append(set.Citations, c)
			}
		}
	}
	return set, nil
}

// anchors reads the directives of a production file's line comments.
//
//charter:realizes CH-021
func (r Reader) anchors(rel string, set *facts.Set) error {
	data, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
	if err != nil {
		return err
	}
	for i, l := range strings.Split(string(data), "\n") {
		m := directive.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		at := facts.At{File: rel, Line: i + 1}
		var ids []string
		for _, p := range strings.Split(m[2], ",") {
			if p = strings.TrimSpace(p); p != "" {
				ids = append(ids, p)
			}
		}
		if len(ids) == 0 {
			set.Refused = append(set.Refused, facts.Refused{At: at, Rule: "rule 24", Message: fmt.Sprintf("// charter:%s names at least one spec", m[1])})
			continue
		}
		set.Anchors = append(set.Anchors, facts.Anchor{At: at, Kind: m[1], IDs: ids})
	}
	return nil
}

// extract runs the reader script once over the files.
func (r Reader) extract(files []string) ([]record_, error) {
	node, err := exec.LookPath(Node)
	if err != nil {
		return nil, fmt.Errorf("%s: the charter reads TypeScript and %w: node is not on PATH", r.Table, ErrNoParser)
	}
	compiler := r.Compiler
	if !filepath.IsAbs(compiler) {
		compiler = filepath.Join(r.Root, compiler)
	}
	if _, err := os.Stat(filepath.Join(compiler, "package.json")); err != nil {
		return nil, fmt.Errorf("%s: the charter reads TypeScript and %w: no typescript package at %s", r.Table, ErrNoParser, r.Compiler)
	}
	dir, err := os.MkdirTemp("", "charter-ts-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	script := filepath.Join(dir, "reader.mjs")
	if err := os.WriteFile(script, []byte(readerScript), 0o600); err != nil {
		return nil, err
	}
	req, err := json.Marshal(map[string]any{"root": r.Root, "files": files})
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(node, script, compiler, r.Import, r.EnvFile)
	cmd.Stdin = bytes.NewReader(append(req, '\n'))
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("the TypeScript reader failed: %v\n%s", err, strings.TrimSpace(errOut.String()))
	}
	sc := bufio.NewScanner(&out)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	if !sc.Scan() || sc.Text() != "ready" {
		return nil, fmt.Errorf("the TypeScript reader did not start: %s", strings.TrimSpace(out.String()+errOut.String()))
	}
	if !sc.Scan() {
		return nil, errors.New("the TypeScript reader wrote no reply")
	}
	var reply struct {
		Records []record_ `json:"records"`
		Error   string    `json:"error"`
	}
	if err := json.Unmarshal(sc.Bytes(), &reply); err != nil {
		return nil, fmt.Errorf("the TypeScript reader wrote something other than JSON: %w", err)
	}
	if reply.Error != "" {
		return nil, errors.New(reply.Error)
	}
	return reply.Records, nil
}

// Generate renders index.ts: the claims, mapping each symbol to its id or
// an observation's name, and verifies and observes. regen is the command
// named in its header.
//
//charter:realizes CH-020
func Generate(rec *record.Record, regen string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated from %s by %s. DO NOT EDIT.\n", rec.Dir, regen)
	b.WriteString(marker)
	b.WriteString("\n/** The record's entries and observations, each symbol mapped to its id or name. */\nexport const claims = {\n")
	for _, e := range rec.Entries {
		fmt.Fprintf(&b, "  /** %s */\n", comment(e.Title))
		if e.Status == "deprecated" {
			fmt.Fprintf(&b, "  /** @deprecated %s is retired; see its Changes in %s. */\n", e.ID, e.File)
		}
		fmt.Fprintf(&b, "  %s: %q,\n", e.Symbol, e.ID)
	}
	for _, o := range rec.Observations {
		fmt.Fprintf(&b, "  /** An observation of %s: %s */\n  %s: %q,\n", comment(o.Subject), comment(o.Text), o.Name, o.Name)
	}
	b.WriteString("} as const;\n")
	return b.Bytes()
}

// comment is text for a doc comment: one line, a comment closer escaped.
func comment(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "*/", "*\\/")
}

// marker is the generated module's own code. Under vitest with globals, a
// citation and its test's outcome are appended to the file CHARTER_EVENTS
// names, as a JSON line, once the test finishes.
//
//charter:realizes CH-024
const marker = `
import { appendFileSync } from "node:fs";

type Finished = (fn: (ctx: { task: { name: string; file?: { filepath?: string }; result?: { state?: string } } }) => void) => void;

function record(id: string, handles: string[]): void {
  const path = (globalThis as { process?: { env?: Record<string, string | undefined> } }).process?.env?.["CHARTER_EVENTS"];
  const onTestFinished = (globalThis as { onTestFinished?: Finished }).onTestFinished;
  if (!path || typeof onTestFinished !== "function") return;
  onTestFinished((ctx) => {
    const state = ctx.task.result?.state;
    const outcome = state === "pass" ? "passed" : state === "skip" ? "skipped" : "failed";
    for (const handle of handles.length ? handles : [""]) {
      const line: Record<string, string> = { id, test: ctx.task.name, file: ctx.task.file?.filepath ?? "", outcome };
      if (handle) line.handle = handle;
      appendFileSync(path, JSON.stringify(line) + "\n");
    }
  });
}

/** Opens a test that verifies an entry, naming by handle the scenarios it exercises. */
export function verifies(id: string, ...handles: string[]): void {
  if (typeof id !== "string" || id === "") throw new TypeError("claims: verifies needs a spec");
  record(id, handles);
}

/** Opens a probe that checks an observation against the live system. */
export function observes(name: string): void {
  if (typeof name !== "string" || name === "") throw new TypeError("claims: observes needs an observation");
  record(name, []);
}
`
