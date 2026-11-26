package chartertest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path"
	"sort"
	"strings"
)

// The report is what the gate cannot see: whether the tests that cite each
// scenario ran, and how. It reads the structured output of a test run, the
// attr events Verifies emits and the pass, fail and skip events that follow,
// and joins each citation to its test's result. A scenario is passed when
// a citing test passed and none failed, failed when one failed, skipped
// when every citing test was skipped, and none when no citing test ran.

// Outcome is a test's or a scenario's result in a report.
type Outcome string

const (
	Passed  Outcome = "pass"
	Failed  Outcome = "fail"
	Skipped Outcome = "skip"
	NoRun   Outcome = "none"
	// Waived: no citing test passed, and a waiver line accepts it.
	Waived Outcome = "waived"
)

// TestOutcome is one test's result, with the package it ran in.
type TestOutcome struct {
	Package, Test string
	Outcome       Outcome
}

// ScenarioReport is one scenario with the tests that cited it, and who
// ruled it.
type ScenarioReport struct {
	Handle     string
	Provenance string
	Tests      []TestOutcome
	Outcome    Outcome
}

// Count is how many scenarios, and how many of them passed.
type Count struct {
	Scenarios, Passed int
}

// Profile counts the active scenarios by who ruled them: the owner, a
// delegate, or an agent whose choice nobody ruled on.
type Profile struct {
	Owner, Delegate, Provisional Count
}

// add counts one active scenario under its provenance.
func (p *Profile) add(prov Provenance, passed bool) {
	c := &p.Provisional
	switch prov.Who {
	case "owner":
		c = &p.Owner
	case "delegate":
		c = &p.Delegate
	}
	c.Scenarios++
	if passed {
		c.Passed++
	}
}

// String renders the profile in a sentence.
func (p Profile) String() string {
	part := func(c Count, who string) string {
		return fmt.Sprintf("%d %s (%d passed)", c.Scenarios, who, c.Passed)
	}
	return "Active scenarios by who ruled them: " + part(p.Owner, "ruled by the owner") + ", " + part(p.Delegate, "by a delegate") + ", " + part(p.Provisional, "an agent's, provisional") + "."
}

// EntryReport is one spec with its scenarios' results and, for a citation
// that named no scenario, the tests that cited the spec itself.
type EntryReport struct {
	ID, Title, Status string
	Need              bool
	Scenarios         []ScenarioReport
	Tests             []TestOutcome
	Outcome           Outcome
}

// Report is the result of one run over every spec.
type Report struct {
	Entries []EntryReport
	// Unverified lists the active scenarios, as ID/handle, and the active
	// specs without scenarios, as ID, that no citing test passed: the
	// report's failures.
	Unverified []string
	// Profile counts the active scenarios by who ruled them.
	Profile Profile
}

// goEvent is one line of go test -json.
type goEvent struct {
	Action, Package, Test, Key, Value string
}

// RunGoTests runs the packages under go test -json in the root and returns
// the stream. A failing test is not an error: the stream carries it.
//
//charter:realizes CH-042
func RunGoTests(opts Options, packages ...string) ([]byte, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	if len(packages) == 0 {
		packages = []string{"./..."}
	}
	args := append([]string{"test", "-json", "-count=1"}, packages...)
	cmd := exec.Command("go", args...)
	cmd.Dir = r.Root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, fmt.Errorf("chartertest: go test: %v\n%s", err, strings.TrimSpace(stderr.String()))
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("chartertest: go test produced no stream: %v\n%s", err, strings.TrimSpace(stderr.String()))
		}
	}
	return out, nil
}

// ReadGoStream reads a go test -json stream: the claim attributes each test
// emitted, and each test's outcome.
//
//charter:realizes CH-042
func ReadGoStream(stream io.Reader) (claims map[string][]string, outcomes map[string]Outcome, err error) {
	claims, outcomes = map[string][]string{}, map[string]Outcome{}
	sc := bufio.NewScanner(stream)
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for sc.Scan() {
		var e goEvent
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil || e.Test == "" {
			continue
		}
		key := e.Package + "\x00" + e.Test
		switch e.Action {
		case "attr":
			if e.Key == "claim" {
				claims[key] = append(claims[key], e.Value)
			}
		case "pass":
			outcomes[key] = Passed
		case "fail":
			outcomes[key] = Failed
		case "skip":
			outcomes[key] = Skipped
		}
	}
	return claims, outcomes, sc.Err()
}

// BuildReport reads a go test -json stream and joins it to the charter the
// options name, honouring its waivers: a waived scenario or spec that no
// citing test passed is waived, not unverified.
//
//charter:realizes CH-042
func BuildReport(opts Options, stream io.Reader) (*Report, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	waivers, _, _, err := readWaivers(r)
	if err != nil {
		return nil, err
	}
	waived := func(key string) bool {
		for _, w := range waivers {
			if ok, _ := path.Match(w.Pattern, key); ok {
				return true
			}
		}
		return false
	}
	reg := registryOf(s)
	cl, oc, err := ReadGoStream(stream)
	if err != nil {
		return nil, err
	}
	return Build(reg, cl, oc, waived), nil
}

// Build joins a test run to the charter: for each spec in the registry, the
// tests that cited each scenario and their outcomes. waived, when not nil,
// says whether a scenario (ID/handle) or a spec (ID) is waived.
//
//charter:realizes CH-042
func Build(reg *Registry, claims map[string][]string, outcomes map[string]Outcome, waived func(string) bool) *Report {
	if waived == nil {
		waived = func(string) bool { return false }
	}
	byValue := map[string][]TestOutcome{} // "ID" or "ID/handle" -> tests
	for key, values := range claims {
		pkg, test, _ := strings.Cut(key, "\x00")
		oc, ok := outcomes[key]
		if !ok {
			oc = NoRun
		}
		for _, v := range values {
			byValue[v] = append(byValue[v], TestOutcome{Package: pkg, Test: test, Outcome: oc})
		}
	}
	for _, ts := range byValue {
		sort.Slice(ts, func(i, j int) bool {
			if ts[i].Package != ts[j].Package {
				return ts[i].Package < ts[j].Package
			}
			return ts[i].Test < ts[j].Test
		})
	}
	rep := &Report{}
	for _, e := range reg.Entries {
		if e.Level != levelSpec {
			continue
		}
		er := EntryReport{ID: e.ID, Title: e.Text, Status: e.Status, Need: e.Need, Tests: byValue[e.ID]}
		for _, s := range e.Scenarios {
			key := e.ID + "/" + s.Handle
			tests := byValue[key]
			sr := ScenarioReport{Handle: s.Handle, Provenance: s.Provenance.String(), Tests: tests, Outcome: outcomeOf(tests)}
			if sr.Outcome != Passed && waived(key) {
				sr.Outcome = Waived
			}
			er.Scenarios = append(er.Scenarios, sr)
			if e.Status == "active" {
				rep.Profile.add(s.Provenance, sr.Outcome == Passed)
				if sr.Outcome != Passed && sr.Outcome != Waived {
					rep.Unverified = append(rep.Unverified, key)
				}
			}
		}
		if len(e.Scenarios) == 0 {
			er.Outcome = outcomeOf(er.Tests)
			if er.Outcome != Passed && waived(e.ID) {
				er.Outcome = Waived
			}
			if e.Status == "active" && er.Outcome != Passed && er.Outcome != Waived {
				rep.Unverified = append(rep.Unverified, e.ID)
			}
		} else {
			er.Outcome = Passed
			for _, s := range er.Scenarios {
				er.Outcome = worse(er.Outcome, s.Outcome)
			}
		}
		rep.Entries = append(rep.Entries, er)
	}
	return rep
}

// outcomeOf folds the outcomes of the tests that cited one thing.
func outcomeOf(tests []TestOutcome) Outcome {
	if len(tests) == 0 {
		return NoRun
	}
	out := Skipped
	for _, t := range tests {
		switch t.Outcome {
		case Failed:
			return Failed
		case Passed:
			out = Passed
		case NoRun:
			if out == Skipped {
				out = NoRun
			}
		}
	}
	return out
}

// worse orders outcomes: fail, then none, then skip, then waived, then pass.
func worse(a, b Outcome) Outcome {
	rank := map[Outcome]int{Failed: 0, NoRun: 1, Skipped: 2, Waived: 3, Passed: 4}
	if rank[b] < rank[a] {
		return b
	}
	return a
}

// Markdown renders the report for a pull request: who ruled the active
// scenarios, then each spec with its scenarios and their outcomes, the
// unverified ones first.
//
//charter:realizes CH-042
func (r *Report) Markdown() string {
	var b strings.Builder
	b.WriteString(r.Profile.String() + "\n\n")
	if len(r.Unverified) == 0 {
		b.WriteString("Every active scenario passed.\n")
	} else {
		fmt.Fprintf(&b, "%d active scenarios did not pass:\n\n", len(r.Unverified))
		for _, u := range r.Unverified {
			fmt.Fprintf(&b, "- `%s`\n", u)
		}
		b.WriteString("\n")
	}
	for _, e := range r.Entries {
		fmt.Fprintf(&b, "## %s %s\n\n", e.ID, e.Title)
		fmt.Fprintf(&b, "%s, %s\n\n", e.Status, e.Outcome)
		for _, s := range e.Scenarios {
			fmt.Fprintf(&b, "- `%s` (%s): %s", s.Handle, s.Provenance, s.Outcome)
			for _, t := range s.Tests {
				fmt.Fprintf(&b, " · %s (%s)", t.Test, t.Outcome)
			}
			b.WriteString("\n")
		}
		if len(e.Scenarios) == 0 {
			for _, t := range e.Tests {
				fmt.Fprintf(&b, "- %s (%s)\n", t.Test, t.Outcome)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
