// Package record reads a charter's record: the needs and the area documents
// of its claims directory, the observations, the questions, the waivers and
// the reading order. It knows nothing of code or of any language; what the
// record holds is checked against the tree by the gate.
package record

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// The files of the claims directory that are not area documents. Every other
// .md file there is one, needs.md holding the needs.
const (
	NeedsFile        = "needs.md"
	ReadmeFile       = "README.md"
	QuestionsFile    = "questions.md"
	ObservationsFile = "observations.md"
	PolicyFile       = "policy.md"
	GlossaryFile     = "glossary.md"
	AuditFile        = "audit.md"
	WaiversFile      = "waivers.txt"
	BaselineFile     = "baseline.txt"
	ReadingOrderFile = "reading-order.txt"
)

// Line is one line of a document with its number.
type Line struct {
	N    int
	Text string
}

// Scenario is one numbered scenario of an entry: a handle a test cites, the
// sentence naming an input and an observable outcome, and who ruled it.
type Scenario struct {
	Handle, Text string
	Line         int
	Provenance   Provenance `json:"provenance"`
}

// Change is one dated item of an entry's Changes.
type Change struct {
	Date, Text string
	Line       int
}

// Entry is a need or a decision.
type Entry struct {
	ID string
	// Symbol is the name each language's generated file declares it by:
	// the id with each hyphen an underscore.
	Symbol string
	Title  string
	// Status is planned, active or deprecated.
	Status string
	// Need is true for an entry of needs.md.
	Need bool
	// Prefix, Number, Day and Code are the id's parts: SW-118 has a number,
	// SW-261001-K7Q2 a day and a code.
	Prefix      string
	Number      int
	Day, Code   string
	Kind        string // invariant, bound, structure, or ""
	Serves, See []string
	Accounts    []string // surface keys, kind:name
	// Text is the entry's text; Why its reason; Guidance its advisory how.
	Text, Why, Guidance string
	Scenarios           []Scenario
	Changes             []Change
	// TextLines and WhyLines are the text and the Why as written, each line
	// with its number, the Why's label cut.
	TextLines, WhyLines []Line
	File                string // slash-relative to the root
	Line                int    // the heading's
	fields              map[string]bool
}

// Pos is the entry's heading as file:line.
func (e *Entry) Pos() string { return fmt.Sprintf("%s:%d", e.File, e.Line) }

// HasScenario reports whether the entry lists the handle.
func (e *Entry) HasScenario(handle string) bool {
	for _, s := range e.Scenarios {
		if s.Handle == handle {
			return true
		}
	}
	return false
}

// Area is the header of one area document, or of needs.md.
type Area struct {
	File  string // slash-relative to the root
	Title string
	Code  []string
	// Summary is the header's prose, paragraphs parted by a blank line.
	Summary string
	Prose   []Line
	Needs   bool
}

// Waiver is one line of the waiver file: an id or ID/handle, and the reason
// its missing test is accepted.
type Waiver struct {
	Target, Reason string
	Line           int
}

// Record is a charter's record as read.
type Record struct {
	// Dir is the claims directory, slash-relative to the root.
	Dir          string
	Areas        []*Area
	Entries      []*Entry // by prefix, then by id
	Observations []*Observation
	Questions    []Question
	Waivers      []Waiver
	// ReadingOrder names the area documents, slash-relative to the root, in
	// the order the reading-order file gives, or alphabetical without one.
	ReadingOrder []string
	// Problems are the form problems found, each "file:line: form: ...".
	Problems []string
	byID     map[string]*Entry
	obsByID  map[string]*Observation
}

// Entry returns the entry with the id, or nil.
func (r *Record) Entry(id string) *Entry { return r.byID[id] }

// Observation returns the observation with the name, or nil.
func (r *Record) Observation(name string) *Observation { return r.obsByID[name] }

// BySymbol returns the entry or the observation a generated symbol names:
// an entry's symbol, or an observation's name.
func (r *Record) BySymbol(sym string) (*Entry, *Observation) {
	if o := r.obsByID[sym]; o != nil {
		return nil, o
	}
	return r.byID[strings.ReplaceAll(sym, "_", "-")], nil
}

// Read reads the record in root/dir. skip names files of the directory, by
// base name, that are not area documents beyond those the record itself
// keeps, such as the generated digest.
//
//charter:realizes CH-261126-JC9B
func Read(root, dir string, skip ...string) (*Record, error) {
	r := &Record{Dir: dir, byID: map[string]*Entry{}, obsByID: map[string]*Observation{}}
	abs := filepath.Join(root, filepath.FromSlash(dir))
	des, err := os.ReadDir(abs)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	notArea := append([]string{ReadmeFile, QuestionsFile, ObservationsFile, PolicyFile, GlossaryFile, AuditFile}, skip...)
	var names []string
	for _, de := range des {
		if !de.IsDir() && strings.HasSuffix(de.Name(), ".md") && !slices.Contains(notArea, de.Name()) {
			names = append(names, de.Name())
		}
	}
	sort.Strings(names)
	var entries []*Entry
	for _, name := range names {
		rel := path.Join(dir, name)
		data, err := os.ReadFile(filepath.Join(abs, name))
		if err != nil {
			return nil, err
		}
		a, es, ps := parseArea(rel, string(data))
		r.Areas = append(r.Areas, a)
		entries = append(entries, es...)
		r.Problems = append(r.Problems, ps...)
	}
	r.addEntries(entries)

	if src, ok, err := readOptional(abs, ObservationsFile); err != nil {
		return nil, err
	} else if ok {
		os, ps := parseObservations(path.Join(dir, ObservationsFile), src)
		r.Problems = append(r.Problems, ps...)
		for _, o := range os {
			if r.obsByID[o.Name] == nil {
				r.obsByID[o.Name] = o
				r.Observations = append(r.Observations, o)
			}
		}
	}
	if src, ok, err := readOptional(abs, QuestionsFile); err != nil {
		return nil, err
	} else if ok {
		qs, ps := parseQuestions(path.Join(dir, QuestionsFile), src, r)
		r.Questions, r.Problems = qs, append(r.Problems, ps...)
	}
	if src, ok, err := readOptional(abs, WaiversFile); err != nil {
		return nil, err
	} else if ok {
		r.Waivers = parseWaivers(src)
	}
	src, ok, err := readOptional(abs, ReadingOrderFile)
	if err != nil {
		return nil, err
	}
	r.ReadingOrder, r.Problems = r.readingOrder(src, ok), append(r.Problems, r.orderProblems(src, ok)...)
	return r, nil
}

// addEntries orders the entries by prefix and then by id, and keeps each id
// once, a second declaration a form problem naming both places.
//
//charter:realizes CH-012
func (r *Record) addEntries(entries []*Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Prefix != entries[j].Prefix {
			return entries[i].Prefix < entries[j].Prefix
		}
		return Before(entries[i], entries[j])
	})
	for _, e := range entries {
		if other := r.byID[e.ID]; other != nil {
			r.Problems = append(r.Problems, fmt.Sprintf("%s: form: %s is declared twice, in %s and %s; %s", e.Pos(), e.ID, other.File, e.File, takeNewID(e.Prefix)))
			continue
		}
		r.byID[e.ID] = e
		r.Entries = append(r.Entries, e)
	}
}

func readOptional(dir, name string) (string, bool, error) {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if errors.Is(err, fs.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(data), true, nil
}

// parseWaivers reads the waiver file: a target, a tab and a reason a line,
// blank lines and lines opening with # skipped. The gate holds each line to
// its rules.
func parseWaivers(src string) []Waiver {
	var ws []Waiver
	for i, line := range strings.Split(src, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		target, reason, _ := strings.Cut(line, "\t")
		ws = append(ws, Waiver{Target: strings.TrimSpace(target), Reason: strings.TrimSpace(reason), Line: i + 1})
	}
	return ws
}

// readingOrder lists the area documents in the reading-order file's order,
// then any it leaves out alphabetically; without the file, alphabetically.
// needs.md is not an area of the order.
func (r *Record) readingOrder(src string, ok bool) []string {
	var all []string
	for _, a := range r.Areas {
		if !a.Needs {
			all = append(all, a.File)
		}
	}
	if !ok {
		return all
	}
	var out []string
	seen := map[string]bool{}
	for _, name := range orderLines(src) {
		rel := path.Join(r.Dir, name.Text)
		if slices.Contains(all, rel) && !seen[rel] {
			out = append(out, rel)
			seen[rel] = true
		}
	}
	for _, f := range all {
		if !seen[f] {
			out = append(out, f)
		}
	}
	return out
}

// orderProblems holds the reading-order file to the area documents: one it
// leaves out, a line naming none and a line naming one twice are problems.
//
//charter:realizes CH-261005-6ZBS
func (r *Record) orderProblems(src string, ok bool) []string {
	if !ok {
		return nil
	}
	rel := path.Join(r.Dir, ReadingOrderFile)
	var problems []string
	named := map[string]int{}
	areas := map[string]bool{}
	for _, a := range r.Areas {
		if !a.Needs {
			areas[path.Base(a.File)] = true
		}
	}
	for _, l := range orderLines(src) {
		switch {
		case !areas[l.Text]:
			problems = append(problems, fmt.Sprintf("%s:%d: form: %s names no area document of %s", rel, l.N, l.Text, r.Dir))
		case named[l.Text] > 0:
			problems = append(problems, fmt.Sprintf("%s:%d: form: %s is named already, at line %d", rel, l.N, l.Text, named[l.Text]))
		default:
			named[l.Text] = l.N
		}
	}
	var left []string
	for a := range areas {
		if named[a] == 0 {
			left = append(left, a)
		}
	}
	sort.Strings(left)
	for _, a := range left {
		problems = append(problems, fmt.Sprintf("%s:1: form: the reading order leaves out %s; add it where a reader should meet it", rel, a))
	}
	return problems
}

// orderLines are the reading-order file's names, blank lines and comments
// skipped.
func orderLines(src string) []Line {
	var out []Line
	for i, line := range strings.Split(src, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		out = append(out, Line{i + 1, t})
	}
	return out
}
