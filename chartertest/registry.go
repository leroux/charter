package chartertest

import "slices"

// Registry is the entries package as read from source, each entry with the
// tests that cite it. Renderers and scripts use it to show the charter.
type Registry struct {
	// Entries are the specs by id, then the observations as declared.
	Entries []Entry
	// Problems are what loading found: form, relations, anchors and the
	// entries package's rules. The rules that need the whole tree's
	// surfaces are Check's to report.
	Problems []string
	// Areas are the area documents: needs.md, then the others in the
	// reading order, which a renderer walks a reader through.
	Areas []Area
	// Waivers are the lines of the waiver file, with the problems a line
	// has among Problems.
	Waivers []Waiver
}

// Area is the header of one area document.
type Area struct {
	File  string   // slash-relative to the root
	Title string   // the document's title, its # heading
	Code  []string // the code the area covers, from its Code line
	Needs bool     // needs.md
	// Summary is the header's prose, what the document holds, its
	// paragraphs apart by a blank line.
	Summary string
}

// Waiver is one line of the waiver file: a spec id or a scenario as
// ID/handle, and the reason it is accepted untested.
type Waiver struct {
	Pattern, Reason string
	File            string // slash-relative to the root
	Line            int
}

// Entry is one spec or observation of a Registry.
type Entry struct {
	// Name is the Go symbol: SW_118 for a spec, the variable name for an
	// observation.
	Name string
	// ID is a spec's id, SW-118, or an observation's name.
	ID string
	// Text is a spec's title or an observation's text.
	Text  string
	Level string // spec or observation
	// Status is planned, active or deprecated; for an observation observed,
	// documented, broken or unsure.
	Status string
	// Form is area for an entry of an area document, file for a spec file.
	Form string
	// Need reports a need: an entry of needs.md or a stakeholder spec.
	Need bool
	// Lens is the id's prefix.
	Lens string
	// Kind is invariant, bound or structure, or "".
	Kind string
	// Serves and See are an area entry's relations, by id.
	Serves, See []string
	// Realizes, Concerns and Related are a file spec's relations, by id.
	Realizes, Concerns, Related []string
	// Description is the decision's text, Rationale its Why, Guidance its
	// advisory how; Verification is a file spec's prose.
	Description, Rationale, Guidance, Verification string
	// Scenarios are the numbered scenarios, each with its handle.
	Scenarios []Scenario
	// Changes are the dated entries of Changes.
	Changes  []Change
	Subject  string   // observations only
	Source   string   // observations only
	Observed string   // observations only, YYYY-MM-DD
	Supports []string // observations only: spec ids
	Accounts []string // surface keys, such as env:RATELIMIT_RATE
	// File is the document or the observation's Go file; Pos is File:line.
	File, Pos string
	// Refs are the tests that verify the spec, with the handles each names;
	// for an observation, the probes that cite it.
	Refs []Ref
	// Anchors are the places production code realizes or concerns a spec.
	Anchors []Anchor
}

// Read loads the entries package and the test references without applying
// the rules that need the whole tree, so it never touches the baseline.
//
//charter:realizes CH-040
func Read(opts Options) (*Registry, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	reg := registryOf(s)
	order, orderProblems, err := readingOrder(r, s.docs)
	if err != nil {
		return nil, err
	}
	area := func(d *areaDoc) Area {
		return Area{File: d.File, Title: d.Title, Code: d.Code, Needs: d.Needs, Summary: d.Summary}
	}
	for _, d := range s.docs {
		if d.Needs {
			reg.Areas = append(reg.Areas, area(d))
		}
	}
	for _, d := range order {
		reg.Areas = append(reg.Areas, area(d))
	}
	waivers, rel, problems, err := readWaivers(r)
	if err != nil {
		return nil, err
	}
	for _, w := range waivers {
		reg.Waivers = append(reg.Waivers, Waiver{Pattern: w.Pattern, Reason: w.Reason, File: rel, Line: w.Line})
	}
	reg.Problems = slices.Sorted(slices.Values(append(append(reg.Problems, orderProblems...), problems...)))
	return reg, nil
}

// registryOf is the Registry of a loaded state.
func registryOf(s *state) *Registry {
	reg := &Registry{Problems: slices.Sorted(slices.Values(s.problems))}
	for _, name := range s.pkg.order {
		e := s.pkg.entries[name]
		refs := s.cites[name]
		if e.Level == levelObservation {
			refs = s.refs[name]
		}
		reg.Entries = append(reg.Entries, Entry{
			Name: e.Name, ID: e.ID, Text: e.Text, Level: e.Level, Status: e.Status,
			Form: e.Form, Need: e.isNeed(), Lens: e.Lens, Kind: e.Kind,
			Serves: e.Serves, See: e.See,
			Realizes: e.Realizes, Concerns: e.Concerns, Related: e.Related,
			Description: e.Description, Rationale: e.Rationale, Guidance: e.Guidance, Verification: e.Verification,
			Scenarios: e.Scenarios, Changes: e.Changes,
			Subject: e.Subject, Source: e.Source, Observed: e.Observed, Supports: e.Supports, Accounts: e.Accounts,
			File: e.File, Pos: e.Pos, Refs: refs, Anchors: s.anchors[name],
		})
	}
	return reg
}

// Entry returns the entry named name, or nil.
func (r *Registry) Entry(name string) *Entry {
	for i := range r.Entries {
		if r.Entries[i].Name == name {
			return &r.Entries[i]
		}
	}
	return nil
}
