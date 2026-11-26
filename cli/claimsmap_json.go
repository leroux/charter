package cli

import (
	"encoding/json"
	"io"

	"github.com/leroux/charter/chartertest"
)

// jsonRegistry is the registry as claims-viewer reads it: the entries of
// chartertest.Registry with every string and array written, so the viewer
// needs no defaults.
type jsonRegistry struct {
	// Charter is product or self; Root is the checkout's absolute path,
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
	Handle string `json:"handle"`
	Text   string `json:"text"`
	Line   int    `json:"line"`
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

// writeJSON writes reg to w.
func writeJSON(w io.Writer, reg *chartertest.Registry, charter, source, root string) error {
	out := jsonRegistry{
		Charter: charter, Source: source, Root: root, Problems: nonNil(reg.Problems),
		Areas: []jsonArea{}, Waivers: []jsonWaiver{}, Entries: []jsonEntry{},
	}
	for _, a := range reg.Areas {
		out.Areas = append(out.Areas, jsonArea{File: a.File, Title: a.Title, Code: nonNil(a.Code), Needs: a.Needs, Summary: a.Summary})
	}
	for _, wv := range reg.Waivers {
		out.Waivers = append(out.Waivers, jsonWaiver{Pattern: wv.Pattern, Reason: wv.Reason, File: wv.File, Line: wv.Line})
	}
	for _, e := range reg.Entries {
		refs := []jsonRef{}
		for _, r := range e.Refs {
			refs = append(refs, jsonRef{File: r.File, Func: r.Func, Line: r.Line, Handles: nonNil(r.Handles)})
		}
		scenarios := []jsonScenario{}
		for _, s := range e.Scenarios {
			scenarios = append(scenarios, jsonScenario{Handle: s.Handle, Text: s.Text, Line: s.Line})
		}
		changes := []jsonChange{}
		for _, c := range e.Changes {
			changes = append(changes, jsonChange{Date: c.Date, Text: c.Text, Line: c.Line})
		}
		anchors := []jsonAnchor{}
		for _, a := range e.Anchors {
			anchors = append(anchors, jsonAnchor{Kind: a.Kind, File: a.File, Line: a.Line, Decl: a.Decl})
		}
		out.Entries = append(out.Entries, jsonEntry{
			Name: e.Name, ID: e.ID, Level: e.Level, Status: e.Status, Text: e.Text,
			File: e.File, Pos: e.Pos, Supports: nonNil(e.Supports), Accounts: nonNil(e.Accounts),
			Refs: refs, Subject: e.Subject, Source: e.Source, Observed: e.Observed,
			Form: e.Form, Need: e.Need, Lens: e.Lens, Kind: e.Kind,
			Serves: nonNil(e.Serves), See: nonNil(e.See),
			Realizes: nonNil(e.Realizes), Concerns: nonNil(e.Concerns), Related: nonNil(e.Related),
			Description: e.Description, Rationale: e.Rationale, Guidance: e.Guidance, Verification: e.Verification,
			Scenarios: scenarios, Changes: changes, Anchors: anchors,
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
