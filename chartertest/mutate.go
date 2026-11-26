package chartertest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
)

// Mutation testing: mutate breaks an entry's realizing code in small ways
// and runs the tests that cite the entry against each break. A mutant one of
// them fails on is killed; one they all pass survives, and points at a
// scenario whose tests would not notice the decision broken.

// MutateOptions says what mutate runs.
type MutateOptions struct {
	// IDs are the entries to mutate. With Since, only those of them the diff
	// touched; empty with Since, every active entry the diff touched.
	IDs []string
	// Since is a git revision: when set, mutants are made only on the lines
	// changed between it and the working tree.
	Since string
	// Max caps the generated mutants per entry, drawn with Seed; 0 runs all.
	// Given mutants are never dropped.
	Max  int
	Seed int64
	// Given are mutants a caller wrote; with GivenOnly, the only ones run.
	Given     []GivenMutant
	GivenOnly bool
	// Workers is how many mutants run at once; 0 means the number of CPUs.
	Workers int
	// Timeout bounds one go test run; 0 means two minutes.
	Timeout time.Duration
	// Skip is passed to go test -skip, for tests that cannot run here.
	Skip string
	// Cache is the build cache the runs use. Empty means a cache of the
	// run's own, emptied past CacheLimit and removed at the end; "shared"
	// means the user's.
	Cache      string
	CacheLimit int64
	// Log, when not nil, receives a line of progress now and then.
	Log io.Writer
}

// MutantOutcome is what the citing tests made of a mutant.
type MutantOutcome string

const (
	// Killed: a citing test failed, the test binary crashed, or the run
	// timed out.
	Killed MutantOutcome = "killed"
	// Survived: every citing test passed.
	Survived MutantOutcome = "survived"
	// Unbuilt: the mutant does not build; it counts neither way.
	Unbuilt MutantOutcome = "unbuilt"
)

// MutantResult is one mutant and what happened to it.
type MutantResult struct {
	Mutant
	Outcome  MutantOutcome
	TimedOut bool `json:",omitempty"`
	// Failed are the citing tests that failed; Handles the scenarios they
	// name.
	Failed  []string `json:",omitempty"`
	Handles []string `json:",omitempty"`
	Seconds float64
	// Output is the tail of go test's error output when the mutant did not
	// build or the run did not start.
	Output string `json:",omitempty"`
}

// ScenarioMutation is one scenario's share of an entry's run.
type ScenarioMutation struct {
	Handle string
	// Cited reports a Go test citing the scenario; a scenario cited only
	// from TypeScript is not run.
	Cited bool
	// Killed counts the mutants a test citing it failed on.
	Killed int
	// Hollow: cited, mutants ran, and its tests killed none.
	Hollow bool
}

// EntryMutation is one entry's run.
type EntryMutation struct {
	ID, Title string
	// Tests are the citing Go tests, as package/Test.
	Tests []string
	// Unsteady are the citing tests that failed unmutated; when any did, no
	// mutant ran.
	Unsteady []string `json:",omitempty"`
	// Unmutable are realizing anchors on something other than a function.
	Unmutable []Anchor `json:",omitempty"`
	Killed    int
	Survived  int
	Unbuilt   int
	Scenarios []ScenarioMutation
	Mutants   []MutantResult
}

// Score is the share of built mutants killed, or -1 when none was built.
func (e EntryMutation) Score() float64 {
	if e.Killed+e.Survived == 0 {
		return -1
	}
	return float64(e.Killed) / float64(e.Killed+e.Survived)
}

// MutationReport is a whole run.
type MutationReport struct {
	Entries []EntryMutation
	// Dirty are the files the tree held after the run and not before.
	Dirty []string `json:",omitempty"`
}

// Mutate runs mutation testing over the charter the options name.
//
//charter:realizes CH-261005-VFJR
func Mutate(opts Options, mo MutateOptions) (*MutationReport, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	return MutateRegistry(r.Root, registryOf(s), mo)
}

// MutateRegistry runs mutation testing over a registry read from the tree
// at root.
//
//charter:realizes CH-261005-VFJR
func MutateRegistry(root string, reg *Registry, mo MutateOptions) (*MutationReport, error) {
	if mo.Workers <= 0 {
		mo.Workers = runtime.NumCPU()
	}
	if mo.Timeout <= 0 {
		mo.Timeout = 2 * time.Minute
	}
	if mo.CacheLimit <= 0 {
		mo.CacheLimit = 4 << 30
	}
	var changed map[string]lineSet
	if mo.Since != "" {
		out, err := gitOutput(root, "diff", "-U0", "--no-color", "--no-ext-diff", mo.Since, "--", "*.go")
		if err != nil {
			return nil, err
		}
		changed = parseDiff(out)
	}
	plans, err := planMutants(root, reg, mo, changed)
	if err != nil {
		return nil, err
	}
	before, gitOK := treeFiles(root)

	work, err := os.MkdirTemp("", "claimsgen-mutate-")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(work) }()
	env := os.Environ()
	cache := ""
	switch mo.Cache {
	case "shared":
	case "":
		cache = filepath.Join(work, "gocache")
		env = append(env, "GOCACHE="+cache)
	default:
		cache = mo.Cache
		env = append(env, "GOCACHE="+cache)
	}
	module, err := modulePath(root)
	if err != nil {
		return nil, err
	}
	run := &runner{root: root, work: work, module: module, env: env, mo: mo}

	// Every entry's citing tests run unmutated first; then the mutants of
	// the steady entries run together, so a small entry leaves no worker
	// idle.
	type job struct{ plan, mutant int }
	var controls []job
	for i, p := range plans {
		if len(p.tests) > 0 {
			controls = append(controls, job{i, -1})
		}
	}
	unsteady := make([][]string, len(plans))
	run.pool(len(controls), cache, func(k int) {
		i := controls[k].plan
		ctl := run.one(plans[i], nil)
		if ctl.Outcome == Unbuilt {
			// A run that did not start, under load or short of disk, is
			// tried once more before the entry is called unsteady.
			ctl = run.one(plans[i], nil)
		}
		if ctl.Outcome != Survived {
			unsteady[i] = ctl.Failed
			if len(unsteady[i]) == 0 {
				unsteady[i] = []string{"(the citing tests did not build or run unmutated: " + ctl.Output + ")"}
			}
		}
	})
	var jobs []job
	results := make([][]MutantResult, len(plans))
	for i, p := range plans {
		if len(p.tests) == 0 || len(unsteady[i]) > 0 {
			continue
		}
		results[i] = make([]MutantResult, len(p.mutants))
		for j := range p.mutants {
			jobs = append(jobs, job{i, j})
		}
	}
	run.pool(len(jobs), cache, func(k int) {
		jb := jobs[k]
		m := plans[jb.plan].mutants[jb.mutant]
		results[jb.plan][jb.mutant] = run.one(plans[jb.plan], &m)
	})

	rep := &MutationReport{}
	for i, p := range plans {
		em := EntryMutation{ID: p.entry.ID, Title: p.entry.Text, Unmutable: p.unmutable, Unsteady: unsteady[i], Mutants: results[i]}
		for _, t := range p.tests {
			em.Tests = append(em.Tests, t.pkg+"/"+t.name)
		}
		tally(&em, p)
		rep.Entries = append(rep.Entries, em)
	}
	if gitOK {
		after, _ := treeFiles(root)
		for f := range after {
			if !before[f] {
				rep.Dirty = append(rep.Dirty, f)
			}
		}
		sort.Strings(rep.Dirty)
	}
	return rep, nil
}

// citingTest is a Go test that cites an entry, with its citations.
type citingTest struct {
	pkg, name string // pkg is the slash directory relative to the root
	cites     []Ref
}

// handles are the scenarios the test's citations name.
func (t citingTest) handles() []string {
	var out []string
	for _, c := range t.cites {
		out = append(out, c.Handles...)
	}
	return out
}

// credited are the handles a run that failed the test, id as
// importpath/Test, credits it with; failed holds every test and subtest
// that failed. A citation opening the test counts on any failure, one
// opening a subtest when that subtest failed, and every one when no
// subtest of the test failed, as when its own body failed or it crashed.
//
//charter:realizes CH-261005-VFJR
func (t citingTest) credited(id string, failed []string) []string {
	var subs [][]string
	for _, f := range failed {
		if rest, ok := strings.CutPrefix(f, id+"/"); ok {
			subs = append(subs, strings.Split(rest, "/"))
		}
	}
	var out []string
	for _, c := range t.cites {
		hit := len(subs) == 0 || len(c.Subtest) == 0
		for _, s := range subs {
			hit = hit || sameSubtest(c.Subtest, s)
		}
		if hit {
			out = append(out, c.Handles...)
		}
	}
	return out
}

// sameSubtest reports whether a citation's subtest is the one go test
// named; a segment "" stands for a title that was no literal, any name.
func sameSubtest(cited, named []string) bool {
	if len(cited) != len(named) {
		return false
	}
	for i := range cited {
		if cited[i] != "" && cited[i] != named[i] {
			return false
		}
	}
	return true
}

// plan is one entry's mutants and the tests they run against.
type plan struct {
	entry     *Entry
	tests     []citingTest
	mutants   []Mutant
	unmutable []Anchor
}

// planMutants chooses the entries and makes their mutants.
//
//charter:realizes CH-261005-VFJR
func planMutants(root string, reg *Registry, mo MutateOptions, changed map[string]lineSet) ([]plan, error) {
	want := map[string]bool{}
	for _, id := range mo.IDs {
		if e := entryByID(reg, id); e == nil || e.Level != levelSpec {
			return nil, fmt.Errorf("mutate: %s is not a spec id", id)
		}
		want[id] = true
	}
	given := map[string][]Mutant{}
	for _, g := range mo.Given {
		if e := entryByID(reg, g.Entry); e == nil || e.Level != levelSpec {
			return nil, fmt.Errorf("mutate: given mutant names %q, not a spec id", g.Entry)
		}
		m, err := givenMutant(root, g)
		if err != nil {
			return nil, fmt.Errorf("mutate: %w", err)
		}
		given[g.Entry] = append(given[g.Entry], m)
		want[g.Entry] = true
	}
	var plans []plan
	for i := range reg.Entries {
		e := &reg.Entries[i]
		if e.Level != levelSpec || e.Status != "active" {
			continue
		}
		if len(mo.IDs) > 0 || mo.Since == "" {
			if !want[e.ID] {
				continue
			}
		}
		p := plan{entry: e, tests: citingTests(e)}
		touched := len(given[e.ID]) > 0
		var generated []Mutant
		for _, a := range e.Anchors {
			if a.Kind != "realizes" {
				continue
			}
			fn, err := locateFunc(root, a)
			if errors.Is(err, errNotAFunction) {
				if strings.HasSuffix(a.File, ".go") {
					p.unmutable = append(p.unmutable, a)
				}
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("mutate: %s: %w", a.File, err)
			}
			var lines lineSet
			if changed != nil {
				lines = changed[a.File]
				if lines == nil {
					lines = lineSet{}
				}
				if !lines.overlaps(fn.from, fn.to) {
					continue
				}
			}
			touched = true
			generated = append(generated, mutantsOf(e.ID, fn, lines)...)
		}
		if changed != nil && !touched {
			continue
		}
		if mo.GivenOnly {
			generated = nil
		}
		p.mutants = append(sample(generated, mo.Max, mo.Seed), given[e.ID]...)
		plans = append(plans, p)
	}
	return plans, nil
}

func entryByID(reg *Registry, id string) *Entry {
	for i := range reg.Entries {
		if reg.Entries[i].ID == id {
			return &reg.Entries[i]
		}
	}
	return nil
}

// citingTests are an entry's Go citations, by package and test.
func citingTests(e *Entry) []citingTest {
	by := map[string]*citingTest{}
	var keys []string
	for _, ref := range e.Refs {
		if !strings.HasSuffix(ref.File, "_test.go") || ref.Func == "" {
			continue
		}
		pkg := path.Dir(ref.File)
		k := pkg + "/" + ref.Func
		t, ok := by[k]
		if !ok {
			t = &citingTest{pkg: pkg, name: ref.Func}
			by[k] = t
			keys = append(keys, k)
		}
		t.cites = append(t.cites, ref)
	}
	sort.Strings(keys)
	out := make([]citingTest, 0, len(keys))
	for _, k := range keys {
		out = append(out, *by[k])
	}
	return out
}

// runner runs mutants through go test -overlay.
type runner struct {
	root, work string
	module     string // the module path, joined to a directory for its import path
	env        []string
	mo         MutateOptions
	mu         sync.Mutex
	n          int
}

// pool runs n jobs on the workers, a batch at a time, emptying the run's
// own cache between batches when it has grown past the limit.
func (r *runner) pool(n int, cache string, do func(int)) {
	batch := r.mo.Workers * 8
	for start := 0; start < n; start += batch {
		end := min(start+batch, n)
		var wg sync.WaitGroup
		next := make(chan int)
		for range r.mo.Workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for k := range next {
					do(k)
				}
			}()
		}
		for k := start; k < end; k++ {
			next <- k
		}
		close(next)
		wg.Wait()
		if r.mo.Cache == "" && dirSize(cache) > r.mo.CacheLimit {
			_ = os.RemoveAll(cache)
		}
		if r.mo.Log != nil {
			_, _ = fmt.Fprintf(r.mo.Log, "mutate: %d/%d\n", end, n)
		}
	}
}

// one runs the plan's citing tests against a mutant, or unmutated when m is
// nil.
func (r *runner) one(p plan, m *Mutant) MutantResult {
	res := MutantResult{}
	var args []string
	if m != nil {
		res.Mutant = *m
		r.mu.Lock()
		r.n++
		dir := filepath.Join(r.work, fmt.Sprint(r.n))
		r.mu.Unlock()
		abs := filepath.Join(r.root, filepath.FromSlash(m.File))
		src, err := os.ReadFile(abs)
		if err == nil {
			err = os.MkdirAll(dir, 0o750)
		}
		mutated := filepath.Join(dir, filepath.Base(abs))
		if err == nil {
			err = os.WriteFile(mutated, m.apply(src), 0o600) //nolint:gosec // the run's own temporary directory, named by a counter
		}
		overlay := filepath.Join(dir, "overlay.json")
		if err == nil {
			var ov []byte
			ov, err = json.Marshal(map[string]map[string]string{"Replace": {abs: mutated}})
			if err == nil {
				err = os.WriteFile(overlay, ov, 0o600)
			}
		}
		if err != nil {
			res.Outcome, res.Failed = Unbuilt, []string{err.Error()}
			return res
		}
		defer func() { _ = os.RemoveAll(dir) }()
		args = append(args, "-overlay", overlay)
	}
	names := map[string]bool{}
	pkgs := map[string]bool{}
	for _, t := range p.tests {
		names[regexp.QuoteMeta(t.name)] = true
		pkgs["./"+t.pkg] = true
	}
	args = append(args, "-vet=off", "-count=1", "-json", "-timeout", r.mo.Timeout.String(), "-run", "^("+strings.Join(slices.Sorted(maps.Keys(names)), "|")+")$")
	if r.mo.Skip != "" {
		args = append(args, "-skip", r.mo.Skip)
	}
	args = append(args, slices.Sorted(maps.Keys(pkgs))...)
	ctx, cancel := context.WithTimeout(context.Background(), r.mo.Timeout+time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", append([]string{"test"}, args...)...)
	cmd.Dir, cmd.Env = r.root, r.env
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	begin := time.Now()
	stream, _ := cmd.Output()
	res.Seconds = time.Since(begin).Round(100 * time.Millisecond).Seconds()
	outcome, timedOut, failed := judge(stream, stderr.Bytes(), ctx.Err() != nil)
	res.Outcome, res.TimedOut = outcome, timedOut
	if outcome == Unbuilt {
		res.Output = tail(stderr.String()+" "+buildOutput(stream), 400)
	}
	handles := map[string]bool{}
	for _, f := range failed {
		for _, t := range p.tests {
			if path.Join(r.module, t.pkg)+"/"+t.name == f {
				res.Failed = append(res.Failed, t.pkg+"/"+t.name)
				for _, h := range t.credited(f, failed) {
					handles[h] = true
				}
			}
		}
	}
	res.Handles = slices.Sorted(maps.Keys(handles))
	return res
}

// judge reads a go test -json stream: the tests and subtests that failed,
// as importpath/Test and importpath/Test/sub, and whether the run built,
// passed or timed out.
//
//charter:realizes CH-261005-VFJR
func judge(stream, stderr []byte, deadline bool) (MutantOutcome, bool, []string) {
	failed := map[string]bool{}
	// running holds the top-level tests started and not yet ended: when a
	// package fails with no build failure while one runs, the binary
	// crashed under it, and go test -json gives that test no fail of its own.
	running := map[string]bool{}
	crashed := false
	built := false
	unbuilt := false
	timedOut := deadline
	for _, line := range bytes.Split(stream, []byte("\n")) {
		var e struct {
			Action, Package, Test, Output, FailedBuild string
		}
		if json.Unmarshal(line, &e) != nil {
			continue
		}
		top := ""
		if e.Test != "" {
			top = e.Package + "/" + strings.SplitN(e.Test, "/", 2)[0]
		}
		switch {
		case e.Action == "build-fail":
			unbuilt = true
		case e.Action == "fail" && e.FailedBuild != "":
			unbuilt = true
		case e.Action == "fail" && e.Test != "":
			failed[top] = true
			failed[e.Package+"/"+e.Test] = true
			delete(running, top)
		case e.Action == "fail":
			crashed = true
			for t := range running {
				failed[t] = true
			}
		case e.Action == "run":
			built = true
			if !strings.Contains(e.Test, "/") {
				running[top] = true
			}
		case (e.Action == "pass" || e.Action == "skip") && e.Test != "" && !strings.Contains(e.Test, "/"):
			delete(running, top)
		case e.Action == "output" && strings.Contains(e.Output, "panic: test timed out"):
			timedOut = true
		}
	}
	if bytes.Contains(stderr, []byte("setup failed")) || bytes.Contains(stderr, []byte("build failed")) {
		unbuilt = true
	}
	switch {
	case len(failed) > 0 || timedOut || (crashed && built && !unbuilt):
		return Killed, timedOut, slices.Sorted(maps.Keys(failed))
	case unbuilt || !built:
		return Unbuilt, false, nil
	}
	return Survived, false, nil
}

// tally counts an entry's outcomes and its scenarios' kills.
//
//charter:realizes CH-261005-VFJR
func tally(em *EntryMutation, p plan) {
	cited := map[string]bool{}
	for _, t := range p.tests {
		for _, h := range t.handles() {
			cited[h] = true
		}
	}
	kills := map[string]int{}
	for _, m := range em.Mutants {
		switch m.Outcome {
		case Killed:
			em.Killed++
			for _, h := range m.Handles {
				kills[h]++
			}
		case Survived:
			em.Survived++
		case Unbuilt:
			em.Unbuilt++
		}
	}
	ran := em.Killed+em.Survived > 0
	for _, s := range p.entry.Scenarios {
		sm := ScenarioMutation{Handle: s.Handle, Cited: cited[s.Handle], Killed: kills[s.Handle]}
		sm.Hollow = sm.Cited && ran && sm.Killed == 0
		em.Scenarios = append(em.Scenarios, sm)
	}
}

// gitOutput runs a git command in root and returns its output, or an
// error that carries git's stderr.
func gitOutput(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// treeFiles is the set of files git reports untracked or changed, and
// whether git could say.
func treeFiles(root string) (map[string]bool, bool) {
	out, err := gitOutput(root, "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return nil, false
	}
	files := map[string]bool{}
	for _, l := range strings.Split(string(out), "\n") {
		if len(l) > 3 {
			files[l[3:]] = true
		}
	}
	return files, true
}

func dirSize(dir string) int64 {
	var n int64
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if info, err := d.Info(); err == nil {
				n += info.Size()
			}
		}
		return nil
	})
	return n
}

// Markdown renders the run for a person: a line per entry, then per entry
// its unsteady tests, hollow scenarios, unmutated anchors and survivors.
//
//charter:realizes CH-261005-VFJR
func (r *MutationReport) Markdown() string {
	var b strings.Builder
	b.WriteString("# Mutation\n\n")
	if len(r.Dirty) > 0 {
		b.WriteString("The run left files in the tree, which fails it:\n\n")
		for _, f := range r.Dirty {
			fmt.Fprintf(&b, "- `%s`\n", f)
		}
		b.WriteString("\n")
	}
	b.WriteString("| Entry | Killed | Survived | Unbuilt | Score | Hollow scenarios |\n|---|---|---|---|---|---|\n")
	for _, e := range r.Entries {
		score := "–"
		switch {
		case len(e.Unsteady) > 0:
			score = "unsteady"
		case len(e.Tests) == 0:
			score = "no Go tests"
		case e.Score() >= 0:
			score = fmt.Sprintf("%.0f%%", 100*e.Score())
		}
		var hollow []string
		for _, s := range e.Scenarios {
			if s.Hollow {
				hollow = append(hollow, "`"+s.Handle+"`")
			}
		}
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %s | %s |\n", e.ID, e.Killed, e.Survived, e.Unbuilt, score, strings.Join(hollow, ", "))
	}
	for _, e := range r.Entries {
		fmt.Fprintf(&b, "\n## %s %s\n\n", e.ID, e.Title)
		if len(e.Unsteady) > 0 {
			fmt.Fprintf(&b, "Unsteady: these citing tests fail unmutated, so no mutant ran: %s.\n\n", strings.Join(e.Unsteady, ", "))
		}
		if len(e.Tests) == 0 {
			b.WriteString("No Go test cites it, so no mutant ran.\n\n")
		}
		for _, a := range e.Unmutable {
			fmt.Fprintf(&b, "- Not mutated: `%s:%d` anchors `%s`, not a function.\n", a.File, a.Line, a.Decl)
		}
		for _, s := range e.Scenarios {
			switch {
			case s.Hollow:
				fmt.Fprintf(&b, "- `%s`: hollow, its tests killed no mutant.\n", s.Handle)
			case s.Cited:
				fmt.Fprintf(&b, "- `%s`: killed %d.\n", s.Handle, s.Killed)
			default:
				fmt.Fprintf(&b, "- `%s`: no Go test cites it.\n", s.Handle)
			}
		}
		var survivors []MutantResult
		for _, m := range e.Mutants {
			if m.Outcome == Survived {
				survivors = append(survivors, m)
			}
		}
		if len(survivors) > 0 {
			b.WriteString("\nSurvivors:\n\n")
			for _, m := range survivors {
				note := ""
				if m.Note != "" {
					note = " — " + m.Note
				}
				fmt.Fprintf(&b, "- `%s:%d` %s, %s: `%s` → `%s`%s\n", m.File, m.Line, m.Func, m.Op, m.Orig, m.Repl, note)
			}
		}
	}
	return b.String()
}

// tail is the last n bytes of s, on one line.
func tail(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > n {
		return "…" + s[len(s)-n:]
	}
	return s
}

// buildOutput is what a go test -json stream says outside any test: the
// build's errors and a package's setup failure.
func buildOutput(stream []byte) string {
	var b strings.Builder
	for _, line := range bytes.Split(stream, []byte("\n")) {
		var e struct{ Action, Test, Output string }
		if json.Unmarshal(line, &e) == nil && e.Test == "" && (e.Action == "build-output" || e.Action == "output") {
			b.WriteString(e.Output)
		}
	}
	return b.String()
}
