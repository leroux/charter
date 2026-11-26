package chartertest

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed refs.mjs
var refsScript string

// isTSFile reports whether rel is a TypeScript source file the walk lists:
// .ts or .tsx, never a .d.ts declaration.
//
//charter:realizes CH-036
func isTSFile(rel string) bool {
	return (strings.HasSuffix(rel, ".ts") || strings.HasSuffix(rel, ".tsx")) && !strings.HasSuffix(rel, ".d.ts")
}

// isTSTest reports whether rel is a TypeScript test file, the only kind
// whose citations count.
func isTSTest(rel string) bool {
	for _, s := range []string{".test.ts", ".test.tsx", ".spec.ts", ".spec.tsx"} {
		if strings.HasSuffix(rel, s) {
			return true
		}
	}
	return false
}

// tsFile is what the extractor reports for one file.
type tsFile struct {
	File        string         `json:"file"`
	Diagnostics []tsDiagnostic `json:"diagnostics"`
	Imports     []tsImport     `json:"imports"`
	Refs        []tsRef        `json:"refs"`
}

// tsDiagnostic is a parse error the compiler found.
type tsDiagnostic struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// tsImport is one import of the module. Form is named, namespace, default,
// side-effect or type. Local is the binding for named, namespace and
// default imports; Imported is the exported name a named import takes.
type tsImport struct {
	Line     int    `json:"line"`
	Form     string `json:"form"`
	Local    string `json:"local"`
	Imported string `json:"imported"`
}

// tsRef is one property access on a binding of claims, with the title of
// the enclosing test, or "" outside one; whether any test call encloses it;
// and, when it is the first argument of a call of the bare name verifies,
// the call's other arguments as handles, nil for one that is not a string
// literal.
type tsRef struct {
	Line    int       `json:"line"`
	Name    string    `json:"name"`
	Test    string    `json:"test"`
	InTest  bool      `json:"inTest"`
	Call    bool      `json:"call"`
	Handles []*string `json:"handles"`
	// Disabled says an enclosing test or suite is skip, fixme, fail or
	// todo: the body never runs as a test, so what it cites is no citation.
	Disabled bool `json:"disabled"`
	// LocalCallee says the verifies called is declared in the test file
	// itself rather than imported: a local one marks nothing.
	LocalCallee bool `json:"localCallee"`
	// Leading says the call is a statement of the test's own body with
	// only verifies calls before it, as markLeading reads a Go test: the
	// only place a citation counts, before anything can skip, branch or
	// enclose it.
	Leading bool `json:"leading"`
}

// extractTS runs the embedded extractor over files, slash paths relative
// to the root, and returns what it reports, in the same order. A missing
// node or compiler, a parse diagnostic, a failed run or unreadable output
// is an error: the gate cannot answer without the facts.
func extractTS(r *resolved, files []string) ([]tsFile, error) {
	ts := r.TypeScript
	node, err := exec.LookPath(ts.Node)
	if err != nil {
		return nil, fmt.Errorf("chartertest: %s not found on PATH; TypeScript citations need node", ts.Node)
	}
	compiler := ts.Compiler
	if !filepath.IsAbs(compiler) {
		compiler = filepath.Join(r.Root, compiler)
	}
	if _, err := os.Stat(filepath.Join(compiler, "package.json")); err != nil {
		if ts.Install != "" {
			return nil, fmt.Errorf("chartertest: no TypeScript compiler at %s; run %s", ts.Compiler, ts.Install)
		}
		return nil, fmt.Errorf("chartertest: no TypeScript compiler at %s; install the package Options.TypeScript.Compiler names", ts.Compiler)
	}
	if len(files) == 0 {
		return nil, nil
	}
	out, err := tsWorkerFor(node, compiler, ts.Import).extract(r.Root, files)
	if err != nil {
		return nil, fmt.Errorf("chartertest: the TypeScript extractor failed: %w", err)
	}
	var result []tsFile
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("chartertest: the TypeScript extractor wrote something other than JSON: %w", err)
	}
	for _, f := range result {
		for _, d := range f.Diagnostics {
			return nil, fmt.Errorf("parse %s:%d: %s", f.File, d.Line, d.Message)
		}
	}
	return result, nil
}

// tsWorker is one node process running the extractor: loading the
// compiler costs most of a second and parsing a few files milliseconds, so
// a process serves every extraction with the same node, compiler and
// import, one request at a time, until the program exits and closes its
// stdin.
type tsWorker struct {
	node, compiler, specifier string

	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr *lockedBuffer
}

var (
	tsWorkersMu sync.Mutex
	tsWorkers   = map[string]*tsWorker{}
)

// tsWorkerFor returns the worker for node, compiler and import, started at
// its first request.
func tsWorkerFor(node, compiler, specifier string) *tsWorker {
	key := node + "\x00" + compiler + "\x00" + specifier
	tsWorkersMu.Lock()
	defer tsWorkersMu.Unlock()
	w, ok := tsWorkers[key]
	if !ok {
		w = &tsWorker{node: node, compiler: compiler, specifier: specifier}
		tsWorkers[key] = w
	}
	return w
}

// extract sends one request and returns the records it was answered with,
// starting the process first if it is not running. A process that fails is
// dropped, and the next request starts another.
func (w *tsWorker) extract(root string, files []string) (json.RawMessage, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stdin == nil {
		if err := w.start(); err != nil {
			return nil, err
		}
	}
	req, err := json.Marshal(struct {
		Root  string   `json:"root"`
		Files []string `json:"files"`
	}{root, files})
	if err != nil {
		return nil, err
	}
	var reply struct {
		Records json.RawMessage `json:"records"`
		Error   string          `json:"error"`
	}
	if _, err := w.stdin.Write(append(req, '\n')); err != nil {
		return nil, w.fail(err)
	}
	line, err := w.stdout.ReadBytes('\n')
	if err != nil {
		return nil, w.fail(err)
	}
	if err := json.Unmarshal(line, &reply); err != nil {
		return line, nil // not a reply; the caller reports what it got
	}
	if reply.Error != "" {
		return nil, errors.New(reply.Error)
	}
	return reply.Records, nil
}

// start runs node over the embedded script and waits for it to load the
// compiler. The script's directory goes once node has loaded it.
func (w *tsWorker) start() error {
	dir, err := os.MkdirTemp("", "chartertest-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	script := filepath.Join(dir, "refs.mjs")
	if err := os.WriteFile(script, []byte(refsScript), 0o600); err != nil {
		return err
	}
	cmd := exec.Command(w.node, script, w.compiler, w.specifier)
	w.stderr = &lockedBuffer{}
	cmd.Stderr = w.stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	w.cmd, w.stdin, w.stdout = cmd, stdin, bufio.NewReader(stdout)
	ready, err := w.stdout.ReadString('\n')
	if err != nil || ready != "ready\n" {
		if err == nil {
			err = fmt.Errorf("it wrote %q before it was ready", ready)
		}
		return w.fail(err)
	}
	return nil
}

// fail stops the process, so the next request starts another, and returns
// err with what the process wrote to stderr.
func (w *tsWorker) fail(err error) error {
	_ = w.stdin.Close()
	werr := w.cmd.Wait()
	w.stdin, w.stdout = nil, nil
	if werr != nil {
		err = werr
	}
	return fmt.Errorf("%v\n%s", err, strings.TrimSpace(w.stderr.String()))
}

// lockedBuffer is a bytes.Buffer the process's stderr copier writes while
// a request reads it.
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

// collectTSRefs runs the extractor over the tree's TypeScript files and
// applies rules 16, 17 and the rule 24 half about verifies calls. The
// references and the citations come back keyed by entry name, in file
// order, one per site.
//
//charter:realizes CH-023
func collectTSRefs(t *tree, r *resolved, pkg *entriesPackage) (refs, cites map[string][]Ref, problems []string, err error) {
	files, err := extractTS(r, t.ts)
	if err != nil {
		return nil, nil, nil, err
	}
	spec := r.TypeScript.Import
	refs, cites = map[string][]Ref{}, map[string][]Ref{}
	seen := map[string]bool{}
	add := func(m map[string][]Ref, name string, ref Ref) {
		k := name + "\x00" + ref.key()
		if seen[k] {
			return
		}
		seen[k] = true
		m[name] = append(m[name], ref)
	}
	for _, f := range files {
		test := isTSTest(f.File)
		for _, imp := range f.Imports {
			switch {
			case imp.Form == "type":
			case !test:
				problems = append(problems, fmt.Sprintf("%s:%d: rule 16: %s is imported from test files only", f.File, imp.Line, spec))
			case imp.Form != "named" || imp.Imported != "claims":
				problems = append(problems, fmt.Sprintf("%s:%d: rule 16: import the module as { claims } from %q", f.File, imp.Line, spec))
			}
		}
		if !test {
			continue
		}
		for _, ref := range f.Refs {
			e, ok := pkg.entries[ref.Name]
			switch {
			case !ok:
				problems = append(problems, fmt.Sprintf("%s:%d: rule 17: %s is not an entry", f.File, ref.Line, ref.Name))
				continue
			case e.Level == levelObservation:
				problems = append(problems, fmt.Sprintf("%s:%d: rule 17: %s is an observation; a TypeScript test cites specs", f.File, ref.Line, ref.Name))
				continue
			}
			rf := Ref{File: f.File, Func: ref.Test, Line: ref.Line}
			if !ref.Call {
				add(refs, ref.Name, rf)
				continue
			}
			if !ref.InTest {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 24: verifies sits outside a test; it is the first statement of the test that verifies the spec", f.File, ref.Line))
				continue
			}
			if !ref.Leading {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 24: verifies is the leading statement of its test, before anything runs, skips or branches; here something precedes or encloses it", f.File, ref.Line))
				continue
			}
			if ref.Disabled {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 24: verifies sits in a test that is skip, fixme, fail or todo, whose body proves nothing; it is no citation", f.File, ref.Line))
				continue
			}
			if ref.LocalCallee {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 24: verifies is declared in the test file itself; a citation calls the helper the suite imports, which marks the running test", f.File, ref.Line))
				continue
			}
			for _, h := range ref.Handles {
				if h == nil {
					problems = append(problems, fmt.Sprintf("%s:%d: rule 24: a scenario handle is a string literal", f.File, ref.Line))
					continue
				}
				rf.Handles = append(rf.Handles, *h)
			}
			add(cites, ref.Name, rf)
		}
	}
	return refs, cites, problems, nil
}

// UpdateTypeScriptEnv is the environment variable that makes Check write
// the TypeScript module instead of failing rule 15 when it differs.
const UpdateTypeScriptEnv = "CHARTER_UPDATE_TYPESCRIPT"

// renderTypeScript is the TypeScript module for the specs: every spec, its
// symbol mapped to its id, grouped by prefix, each with its title as a doc
// comment, under a header naming regen, the command that regenerates it.
// Rule 15 compares the file to this byte for byte, so the form is fixed.
//
//charter:realizes CH-020
func renderTypeScript(pkg *entriesPackage, entriesDir, regen string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "// Generated from %s by the charter gate. Do not edit.\n", entriesDir)
	fmt.Fprintf(&b, "// Regenerate: %s\n\n", regen)
	b.WriteString("export const claims = {\n")
	lens, started := "", false
	for _, name := range pkg.order {
		e := pkg.entries[name]
		if e.Level != levelSpec {
			continue
		}
		if !started || e.Lens != lens {
			if started {
				b.WriteString("\n")
			}
			fmt.Fprintf(&b, "  // %s\n", e.Lens)
			lens, started = e.Lens, true
		}
		text := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", " "), "*/", `*\/`)
		switch e.Status {
		case "deprecated":
			text = "@deprecated " + text
		case "planned":
			text = "Planned: " + text
		}
		fmt.Fprintf(&b, "  /** %s */\n  %s: %q,\n", text, name, e.ID)
	}
	b.WriteString("} as const;\n\n/** A spec id: what a test verifies. */\nexport type Claim = (typeof claims)[keyof typeof claims];\n")
	return b.String()
}

// checkTypeScriptModule applies rule 15: the module equals the rendering.
// With write set it writes the module instead.
//
//charter:realizes CH-020
func checkTypeScriptModule(r *resolved, pkg *entriesPackage, write bool) ([]string, error) {
	want := renderTypeScript(pkg, r.entriesDir, r.regenerate(UpdateTypeScriptEnv))
	rel := r.TypeScript.Module
	p := filepath.Join(r.Root, filepath.FromSlash(rel))
	if write {
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			return nil, err
		}
		return nil, os.WriteFile(p, []byte(want), 0o600)
	}
	regen := "run " + r.regenerate(UpdateTypeScriptEnv)
	got, err := os.ReadFile(p)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return []string{fmt.Sprintf("%s: rule 15: the TypeScript module is missing; %s", rel, regen)}, nil
	case err != nil:
		return nil, err
	case string(got) != want:
		return []string{fmt.Sprintf("%s: rule 15: the TypeScript module differs from the registry; %s", rel, regen)}, nil
	}
	return nil, nil
}
