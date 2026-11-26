// Package export writes a charter as one JSON file the claims viewer
// reads: the areas in the reading order, the entries with their anchors and
// citations, the waivers, the observations and the gate's problems.
package export

import (
	"encoding/json"
	"io"
	"path"

	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/record"
)

// jsonRegistry is the registry as claims-viewer reads it: the entries of
// the record with every string and array written, so the viewer
// needs no defaults.
type jsonRegistry struct {
	// Charter is the charter's name; Root is the checkout's absolute path,
	// which the viewer turns a file:line into an editor link with.
	Charter  string       `json:"charter"`
	Source   string       `json:"source"`
	Root     string       `json:"root"`
	Problems []string     `json:"problems"`
	Areas    []jsonArea   `json:"areas"`
	Waivers  []jsonWaiver `json:"waivers"`
	Entries  []jsonEntry  `json:"entries"`
}

type jsonArea struct {
	File    string   `json:"file"`
	Title   string   `json:"title"`
	Code    []string `json:"code"`
	Needs   bool     `json:"needs"`
	Summary string   `json:"summary"`
}

type jsonWaiver struct {
	Pattern string `json:"pattern"`
	Reason  string `json:"reason"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

type jsonEntry struct {
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	Level    string    `json:"level"`
	Status   string    `json:"status"`
	Text     string    `json:"text"`
	File     string    `json:"file"`
	Pos      string    `json:"pos"`
	Supports []string  `json:"supports"`
	Accounts []string  `json:"accounts"`
	Refs     []jsonRef `json:"refs"`
	Subject  string    `json:"subject"`
	Source   string    `json:"source"`
	Observed string    `json:"observed"`
	// Specs only; empty for observations. Form is area for an entry of an
	// area document, file for a spec file; Serves and See are an area
	// entry's relations, Realizes, Concerns and Related a file entry's.
	Form         string         `json:"form"`
	Need         bool           `json:"need"`
	Lens         string         `json:"lens"`
	Kind         string         `json:"kind"`
	Serves       []string       `json:"serves"`
	See          []string       `json:"see"`
	Realizes     []string       `json:"realizes"`
	Concerns     []string       `json:"concerns"`
	Related      []string       `json:"related"`
	Description  string         `json:"description"`
	Rationale    string         `json:"rationale"`
	Guidance     string         `json:"guidance"`
	Verification string         `json:"verification"`
	Scenarios    []jsonScenario `json:"scenarios"`
	Changes      []jsonChange   `json:"changes"`
	Anchors      []jsonAnchor   `json:"anchors"`
}

type jsonScenario struct {
	Handle     string `json:"handle"`
	Text       string `json:"text"`
	Line       int    `json:"line"`
	Provenance string `json:"provenance"`
}

type jsonChange struct {
	Date string `json:"date"`
	Text string `json:"text"`
	Line int    `json:"line"`
}

type jsonAnchor struct {
	Kind string `json:"kind"`
	File string `json:"file"`
	Line int    `json:"line"`
	Decl string `json:"decl"`
}

type jsonRef struct {
	File    string   `json:"file"`
	Func    string   `json:"func"`
	Line    int      `json:"line"`
	Handles []string `json:"handles"`
}

// Write writes the charter as the viewer reads it to w: every area in the
// reading order, every entry with its anchors and citations, the waivers,
// the observations and the gate's problems.
func Write(w io.Writer, t *check.Tree, problems []string, source, root string) error {
	rec := t.Record
	out := jsonRegistry{
		Charter: t.Config.Name, Source: source, Root: root, Problems: nonNil(problems),
		Areas: []jsonArea{}, Waivers: []jsonWaiver{}, Entries: []jsonEntry{},
	}
	byFile := map[string]*record.Area{}
	for _, a := range rec.Areas {
		byFile[a.File] = a
	}
	var areas []*record.Area
	for _, a := range rec.Areas {
		if a.Needs {
			areas = append(areas, a)
		}
	}
	for _, f := range rec.ReadingOrder {
		areas = append(areas, byFile[f])
	}
	for _, a := range areas {
		out.Areas = append(out.Areas, jsonArea{File: a.File, Title: a.Title, Code: nonNil(a.Code), Needs: a.Needs, Summary: a.Summary})
	}
	for _, wv := range rec.Waivers {
		out.Waivers = append(out.Waivers, jsonWaiver{Pattern: wv.Target, Reason: wv.Reason, File: path.Join(rec.Dir, record.WaiversFile), Line: wv.Line})
	}
	refs := map[string][]jsonRef{}
	for _, c := range t.Facts.Citations {
		refs[c.Symbol] = append(refs[c.Symbol], jsonRef{File: c.File, Func: c.Test, Line: c.Line, Handles: nonNil(c.Handles)})
	}
	anchors := map[string][]jsonAnchor{}
	for _, a := range t.Facts.Anchors {
		for _, id := range a.IDs {
			anchors[id] = append(anchors[id], jsonAnchor{Kind: a.Kind, File: a.File, Line: a.Line, Decl: a.Decl})
		}
	}
	for _, e := range rec.Entries {
		scenarios := []jsonScenario{}
		for _, s := range e.Scenarios {
			scenarios = append(scenarios, jsonScenario{Handle: s.Handle, Text: s.Text, Line: s.Line, Provenance: s.Provenance.String()})
		}
		changes := []jsonChange{}
		for _, c := range e.Changes {
			changes = append(changes, jsonChange{Date: c.Date, Text: c.Text, Line: c.Line})
		}
		out.Entries = append(out.Entries, jsonEntry{
			Name: e.Symbol, ID: e.ID, Level: "spec", Status: e.Status, Text: e.Title,
			File: e.File, Pos: e.Pos(), Supports: []string{}, Accounts: nonNil(e.Accounts),
			Refs: nonNilRefs(refs[e.Symbol]), Form: "area", Need: e.Need, Lens: e.Prefix, Kind: e.Kind,
			Serves: nonNil(e.Serves), See: nonNil(e.See), Realizes: []string{}, Concerns: []string{}, Related: []string{},
			Description: e.Text, Rationale: e.Why, Guidance: e.Guidance,
			Scenarios: scenarios, Changes: changes, Anchors: nonNilAnchors(anchors[e.ID]),
		})
	}
	for _, o := range rec.Observations {
		out.Entries = append(out.Entries, jsonEntry{
			Name: o.Name, ID: o.Name, Level: "observation", Status: o.Status, Text: o.Text,
			File: o.File, Pos: o.Pos(), Supports: nonNil(o.Supports), Accounts: []string{},
			Refs: nonNilRefs(refs[o.Name]), Subject: o.Subject, Source: o.Source, Observed: o.Observed,
			Serves: []string{}, See: []string{}, Realizes: []string{}, Concerns: []string{}, Related: []string{},
			Scenarios: []jsonScenario{}, Changes: []jsonChange{}, Anchors: []jsonAnchor{},
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(out)
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func nonNilRefs(s []jsonRef) []jsonRef {
	if s == nil {
		return []jsonRef{}
	}
	return s
}

func nonNilAnchors(s []jsonAnchor) []jsonAnchor {
	if s == nil {
		return []jsonAnchor{}
	}
	return s
}
