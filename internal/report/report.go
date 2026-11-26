// Package report runs each configured language's tests with CHARTER_EVENTS
// set, reads the line every citation's marker wrote with its test's
// outcome, and says per entry and per scenario which citing tests passed,
// failed or were skipped, and which active scenarios no citing test passed.
package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/env"
	"github.com/leroux/charter/internal/record"
)

// A scenario is passed when a citing test passed and none failed, failed
// when one failed, skipped when every citing test was skipped, and none
// when no citing test ran.

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

// TestOutcome is one test's result, with the file it is in.
type TestOutcome struct {
	File, Test string
	Outcome    Outcome
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
func (p *Profile) add(prov record.Provenance, passed bool) {
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

// Event is one line of the events file.
type Event struct {
	ID      string `json:"id"`
	Handle  string `json:"handle,omitempty"`
	Test    string `json:"test"`
	File    string `json:"file"`
	Outcome string `json:"outcome"`
}

// Run runs the test command of each language the charter names, in that
// language's directory, with CHARTER_EVENTS set to a file of its own, and
// builds the report from the lines written. A command that cannot start,
// or fails having written no line, is an error naming the language and its
// output; a failing test is not, the events carry it.
//
//charter:realizes CH-042
func Run(c *config.Config) (*Report, error) {
	rec, err := record.Read(c.Root, config.Dir, filepath.Base(c.Digest))
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp("", "charter-report-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	var events []Event
	for _, l := range languages(c) {
		path := filepath.Join(dir, l.name+".jsonl")
		cmd := exec.Command("sh", "-c", l.lang.Test)
		cmd.Dir = filepath.Join(c.Root, filepath.FromSlash(l.lang.Dir))
		cmd.Env = append(os.Environ(), env.Events+"="+path)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		runErr := cmd.Run()
		got, err := readEvents(path)
		if err != nil {
			return nil, err
		}
		if runErr != nil && len(got) == 0 {
			return nil, fmt.Errorf("the %s tests (%s, in %s) failed having recorded no citation: %v\n%s", l.name, l.lang.Test, l.lang.Dir, runErr, strings.TrimSpace(out.String()))
		}
		events = append(events, got...)
	}
	return Build(rec, events), nil
}

type named struct {
	name string
	lang *config.Language
}

func languages(c *config.Config) []named {
	var ls []named
	if c.Go != nil {
		ls = append(ls, named{"Go", &c.Go.Language})
	}
	if c.Python != nil {
		ls = append(ls, named{"Python", c.Python})
	}
	if c.TypeScript != nil {
		ls = append(ls, named{"TypeScript", &c.TypeScript.Language})
	}
	return ls
}

// readEvents reads the events file; a missing one is no events.
func readEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Event
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 1<<16), 1<<24)
	for s.Scan() {
		var ev Event
		if err := json.Unmarshal(s.Bytes(), &ev); err != nil {
			return nil, fmt.Errorf("the events file holds a line that is not an event: %q", s.Text())
		}
		out = append(out, ev)
	}
	return out, s.Err()
}

// Build joins a test run to the charter: for each spec in the registry, the
// tests that cited each scenario and their outcomes. waived, when not nil,
// says whether a scenario (ID/handle) or a spec (ID) is waived.
//
//charter:realizes CH-042
func Build(rec *record.Record, events []Event) *Report {
	waived := func(key string) bool {
		for _, w := range rec.Waivers {
			if w.Target == key && w.Reason != "" {
				return true
			}
		}
		return false
	}
	byValue := map[string][]TestOutcome{} // "ID" or "ID/handle" -> tests
	for _, ev := range events {
		key := ev.ID
		if ev.Handle != "" {
			key += "/" + ev.Handle
		}
		oc := map[string]Outcome{"passed": Passed, "failed": Failed, "skipped": Skipped}[ev.Outcome]
		if oc == "" {
			oc = NoRun
		}
		byValue[key] = append(byValue[key], TestOutcome{File: ev.File, Test: ev.Test, Outcome: oc})
	}
	for _, ts := range byValue {
		sort.Slice(ts, func(i, j int) bool {
			if ts[i].File != ts[j].File {
				return ts[i].File < ts[j].File
			}
			return ts[i].Test < ts[j].Test
		})
	}
	rep := &Report{}
	for _, e := range rec.Entries {
		er := EntryReport{ID: e.ID, Title: e.Title, Status: e.Status, Need: e.Need, Tests: byValue[e.ID]}
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
