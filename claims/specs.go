// Code generated from claims by charter gen. DO NOT EDIT.

// Package claims is the charter's record as Go: a symbol per entry and per
// observation, and the Verifies and Observes a test opens with.
package claims

import (
	"encoding/json"
	"os"
	"runtime"
	"sync"
	"testing"
)

// Status is an entry's status.
type Status string

// The statuses an entry may have.
const (
	Planned    Status = "planned"
	Active     Status = "active"
	Deprecated Status = "deprecated"
)

// Spec is an entry of the record: a need or a decision.
type Spec struct {
	ID, Title string
	Status    Status
	Kind      string
}

// Observation is a fact about a system the code does not own.
type Observation struct {
	Name, Subject, Status string
}

// Verifies opens a test that verifies an entry, naming by handle the
// scenarios it exercises. With CHARTER_EVENTS set, the test's outcome is
// recorded for each citation once it ends.
func Verifies(tb testing.TB, spec *Spec, handles ...string) {
	tb.Helper()
	if spec == nil {
		tb.Fatal("claims: Verifies needs a spec")
		return
	}
	record(tb, spec.ID, handles)
}

// Observes opens a probe that checks an observation against the live
// system.
func Observes(tb testing.TB, o *Observation) {
	tb.Helper()
	if o == nil {
		tb.Fatal("claims: Observes needs an observation")
		return
	}
	record(tb, o.Name, nil)
}

// event is one line of the events file.
type event struct {
	ID      string `json:"id"`
	Handle  string `json:"handle,omitempty"`
	Test    string `json:"test"`
	File    string `json:"file"`
	Outcome string `json:"outcome"`
}

var eventsMu sync.Mutex

func record(tb testing.TB, id string, handles []string) {
	path := os.Getenv("CHARTER_EVENTS")
	if path == "" {
		return
	}
	_, file, _, _ := runtime.Caller(2)
	if len(handles) == 0 {
		handles = []string{""}
	}
	tb.Cleanup(func() {
		outcome := "passed"
		switch {
		case tb.Failed():
			outcome = "failed"
		case tb.Skipped():
			outcome = "skipped"
		}
		eventsMu.Lock()
		defer eventsMu.Unlock()
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			tb.Errorf("claims: the events file: %v", err)
			return
		}
		defer f.Close()
		for _, h := range handles {
			line, _ := json.Marshal(event{ID: id, Handle: h, Test: tb.Name(), File: file, Outcome: outcome})
			// One write per line, opened to append: a line is never split
			// by another process's.
			if _, err := f.Write(append(line, '\n')); err != nil {
				tb.Errorf("claims: the events file: %v", err)
			}
		}
	})
}

// CH_001 is CH-001: The charter gate MUST let a reviewer judge a change by its entries and their tests, without reading the code
var CH_001 = &Spec{ID: "CH-001", Title: "The charter gate MUST let a reviewer judge a change by its entries and their tests, without reading the code", Status: Active}

// CH_002 is CH-002: The charter command MUST lead a builder from an id to its code and tests, and every gate problem MUST name its fix
var CH_002 = &Spec{ID: "CH-002", Title: "The charter command MUST lead a builder from an id to its code and tests, and every gate problem MUST name its fix", Status: Active}

// CH_006 is CH-006: The charter specifies itself, and its own gate is green
var CH_006 = &Spec{ID: "CH-006", Title: "The charter specifies itself, and its own gate is green", Status: Active}

// CH_010 is CH-010: The record reader MUST read an entry as a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes, in that order
var CH_010 = &Spec{ID: "CH-010", Title: "The record reader MUST read an entry as a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes, in that order", Status: Active}

// CH_011 is CH-011: The record reader MUST treat an entry of needs.md as a need, with no Kind, Serves or Guidance
var CH_011 = &Spec{ID: "CH-011", Title: "The record reader MUST treat an entry of needs.md as a need, with no Kind, Serves or Guidance", Status: Active}

// CH_012 is CH-012: An id MUST be declared once, and a new id MUST be its prefix, the day it was taken and a random code
var CH_012 = &Spec{ID: "CH-012", Title: "An id MUST be declared once, and a new id MUST be its prefix, the day it was taken and a random code", Status: Active}

// CH_013 is CH-013: An entry's Status MUST be planned, active or deprecated; every status MUST generate a symbol, and the gate MUST ask coverage of active entries only
var CH_013 = &Spec{ID: "CH-013", Title: "An entry's Status MUST be planned, active or deprecated; every status MUST generate a symbol, and the gate MUST ask coverage of active entries only", Status: Active}

// CH_015 is CH-015: The gate MUST fail a Serves or See that names no existing entry, or that names the entry itself
var CH_015 = &Spec{ID: "CH-015", Title: "The gate MUST fail a Serves or See that names no existing entry, or that names the entry itself", Status: Active}

// CH_020 is CH-020: Every entry has a generated symbol per language read, its id with underscores for hyphens, its title the doc comment
var CH_020 = &Spec{ID: "CH-020", Title: "Every entry has a generated symbol per language read, its id with underscores for hyphens, its title the doc comment", Status: Active}

// CH_021 is CH-021: Production code anchors an entry in a Go doc comment, a TypeScript line comment, or above a Python def
var CH_021 = &Spec{ID: "CH-021", Title: "Production code anchors an entry in a Go doc comment, a TypeScript line comment, or above a Python def", Status: Active}

// CH_022 is CH-022: An anchor names an entry that exists and is not planned, and concerns never counts toward coverage
var CH_022 = &Spec{ID: "CH-022", Title: "An anchor names an entry that exists and is not planned, and concerns never counts toward coverage", Status: Active}

// CH_023 is CH-023: A test cites an entry by a Verifies statement in a test function, naming by handle the scenarios it exercises
var CH_023 = &Spec{ID: "CH-023", Title: "A test cites an entry by a Verifies statement in a test function, naming by handle the scenarios it exercises", Status: Active}

// CH_024 is CH-024: Each language's verifies records each citation and its test's outcome as a JSON line in the events file
var CH_024 = &Spec{ID: "CH-024", Title: "Each language's verifies records each citation and its test's outcome as a JSON line in the events file", Status: Active}

// CH_025 is CH-025: The claims directory holds the record, charter.toml and generated files, and no test assigns to an entry
var CH_025 = &Spec{ID: "CH-025", Title: "The claims directory holds the record, charter.toml and generated files, and no test assigns to an entry", Status: Active}

// CH_030 is CH-030: An active decision is anchored and each scenario cited; a waiver absorbs a missing test, never a missing anchor
var CH_030 = &Spec{ID: "CH-030", Title: "An active decision is anchored and each scenario cited; a waiver absorbs a missing test, never a missing anchor", Status: Active, Kind: "invariant"}

// CH_031 is CH-031: A need needs no anchor and is verified only from the stakeholder's side
var CH_031 = &Spec{ID: "CH-031", Title: "A need needs no anchor and is verified only from the stakeholder's side", Status: Active}

// CH_032 is CH-032: A probe cites observations only, and a test never cites an observation
var CH_032 = &Spec{ID: "CH-032", Title: "A probe cites observations only, and a test never cites an observation", Status: Active}

// CH_033 is CH-033: Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces
var CH_033 = &Spec{ID: "CH-033", Title: "Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces", Status: Active, Kind: "invariant"}

// CH_034 is CH-034: An environment variable is read through a named constant, or under a dynamic directive, and no literal names one
var CH_034 = &Spec{ID: "CH-034", Title: "An environment variable is read through a named constant, or under a dynamic directive, and no literal names one", Status: Active, Kind: "invariant"}

// CH_035 is CH-035: Observations carry their status's fields, observed ones a date and a probe; one broken under an active entry fails
var CH_035 = &Spec{ID: "CH-035", Title: "Observations carry their status's fields, observed ones a date and a probe; one broken under an active entry fails", Status: Active}

// CH_036 is CH-036: The walk covers every Go, Python and TypeScript file under the root, skipping another charter's tree and what charter.toml ignores
var CH_036 = &Spec{ID: "CH-036", Title: "The walk covers every Go, Python and TypeScript file under the root, skipping another charter's tree and what charter.toml ignores", Status: Active}

// CH_037 is CH-037: The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree
var CH_037 = &Spec{ID: "CH-037", Title: "The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree", Status: Active}

// CH_040 is CH-040: One charter command generates, checks, gives ids and lists anchors, and an adopter writes no main
var CH_040 = &Spec{ID: "CH-040", Title: "One charter command generates, checks, gives ids and lists anchors, and an adopter writes no main", Status: Active}

// CH_042 is CH-042: A report joins each test's markers to its scenarios and says which passed, failed, were skipped or have none
var CH_042 = &Spec{ID: "CH-042", Title: "A report joins each test's markers to its scenarios and says which passed, failed, were skipped or have none", Status: Active}

// CH_261005_6ZBS is CH-261005-6ZBS: The digest reads the charter in order: needs by area, an area map, then each area's decisions with their why
var CH_261005_6ZBS = &Spec{ID: "CH-261005-6ZBS", Title: "The digest reads the charter in order: needs by area, an area map, then each area's decisions with their why", Status: Active}

// CH_261009_3FHW is CH-261009-3FHW: The record MUST show who ruled each scenario, and charter diff MUST name any change to an entry the record does not state
var CH_261009_3FHW = &Spec{ID: "CH-261009-3FHW", Title: "The record MUST show who ruled each scenario, and charter diff MUST name any change to an entry the record does not state", Status: Active}

// CH_261009_4B3B is CH-261009-4B3B: The question file holds what nobody ruled, each question naming its entry, who asked and the provisional answer until answered
var CH_261009_4B3B = &Spec{ID: "CH-261009-4B3B", Title: "The question file holds what nobody ruled, each question naming its entry, who asked and the provisional answer until answered", Status: Active}

// CH_261009_QQC8 is CH-261009-QQC8: charter diff names what a change did to the record that its Changes lines do not say
var CH_261009_QQC8 = &Spec{ID: "CH-261009-QQC8", Title: "charter diff names what a change did to the record that its Changes lines do not say", Status: Active}

// CH_261009_VY5C is CH-261009-VY5C: A scenario says who ruled it: the owner with what it cites, a named delegate, or, unwritten, an agent, provisional
var CH_261009_VY5C = &Spec{ID: "CH-261009-VY5C", Title: "A scenario says who ruled it: the owner with what it cites, a named delegate, or, unwritten, an agent, provisional", Status: Active}

// CH_261126_33GK is CH-261126-33GK: charter finds claims/charter.toml by walking up from the working directory, and refuses a key it does not know
var CH_261126_33GK = &Spec{ID: "CH-261126-33GK", Title: "charter finds claims/charter.toml by walking up from the working directory, and refuses a key it does not know", Status: Active}

// CH_261126_5JMD is CH-261126-5JMD: The loop's skills ship as a plugin of this repository and read the layout from charter.toml alone
var CH_261126_5JMD = &Spec{ID: "CH-261126-5JMD", Title: "The loop's skills ship as a plugin of this repository and read the layout from charter.toml alone", Status: Active}

// CH_261126_HK6S is CH-261126-HK6S: The record reader MUST read an observation as a heading in observations.md, with its subject, status, date or source and the entries it supports
var CH_261126_HK6S = &Spec{ID: "CH-261126-HK6S", Title: "The record reader MUST read an observation as a heading in observations.md, with its subject, status, date or source and the entries it supports", Status: Active}

// CH_261126_JC9B is CH-261126-JC9B: The record reader MUST treat an area document's header as a title, a Code line and prose, and MUST NOT read a fenced line as a heading
var CH_261126_JC9B = &Spec{ID: "CH-261126-JC9B", Title: "The record reader MUST treat an area document's header as a title, a Code line and prose, and MUST NOT read a fenced line as a heading", Status: Active}

// CH_261126_NS7T is CH-261126-NS7T: An adopter MUST be able to keep a charter in any repository layout and language, with code that imports nothing of the tooling
var CH_261126_NS7T = &Spec{ID: "CH-261126-NS7T", Title: "An adopter MUST be able to keep a charter in any repository layout and language, with code that imports nothing of the tooling", Status: Active}

// CH_261126_NTRR is CH-261126-NTRR: A declared surface kind is a command run from the root whose output lines are its keys
var CH_261126_NTRR = &Spec{ID: "CH-261126-NTRR", Title: "A declared surface kind is a command run from the root whose output lines are its keys", Status: Active}

// CH_261126_P2SK is CH-261126-P2SK: The charter gate MUST fail a change that leaves a gap between the record and the code
var CH_261126_P2SK = &Spec{ID: "CH-261126-P2SK", Title: "The charter gate MUST fail a change that leaves a gap between the record and the code", Status: Active}

// CH_261126_PFG5 is CH-261126-PFG5: A repository's charter MUST write what the software must do once, above the code, where people and agents both read and change it
var CH_261126_PFG5 = &Spec{ID: "CH-261126-PFG5", Title: "A repository's charter MUST write what the software must do once, above the code, where people and agents both read and change it", Status: Active}

// CH_261126_PQ7X is CH-261126-PQ7X: A language is a reader, a generator and a marker, built in, and read only when charter.toml names it
var CH_261126_PQ7X = &Spec{ID: "CH-261126-PQ7X", Title: "A language is a reader, a generator and a marker, built in, and read only when charter.toml names it", Status: Active, Kind: "structure"}

// CH_261127_3NY3 is CH-261127-3NY3: A doc is a markdown file charter.toml names; it declares no entry, and its references to entries must resolve
var CH_261127_3NY3 = &Spec{ID: "CH-261127-3NY3", Title: "A doc is a markdown file charter.toml names; it declares no entry, and its references to entries must resolve", Status: Active}

// CH_261127_45KM is CH-261127-45KM: A file no language reads anchors an entry by a comment line, where charter.toml names it
var CH_261127_45KM = &Spec{ID: "CH-261127-45KM", Title: "A file no language reads anchors an entry by a comment line, where charter.toml names it", Status: Active}

// CH_261127_8MYR is CH-261127-8MYR: Which docs mention an entry is computed from the docs, and anchors and diff show it
var CH_261127_8MYR = &Spec{ID: "CH-261127-8MYR", Title: "Which docs mention an entry is computed from the docs, and anchors and diff show it", Status: Active}

// CH_261127_DW8R is CH-261127-DW8R: A doc references code by a link to the file declaring the named symbol, and one that file lacks fails
var CH_261127_DW8R = &Spec{ID: "CH-261127-DW8R", Title: "A doc references code by a link to the file declaring the named symbol, and one that file lacks fails", Status: Active}

// CH_261127_GXP7 is CH-261127-GXP7: charter curate records who curated a doc and a hash of what it references; check warns of one whose references changed
var CH_261127_GXP7 = &Spec{ID: "CH-261127-GXP7", Title: "charter curate records who curated a doc and a hash of what it references; check warns of one whose references changed", Status: Active}

// CH_261127_H4B7 is CH-261127-H4B7: charter gen keeps each doc's references and embeds current with the record and the code
var CH_261127_H4B7 = &Spec{ID: "CH-261127-H4B7", Title: "charter gen keeps each doc's references and embeds current with the record and the code", Status: Active}

// CH_261127_QRQS is CH-261127-QRQS: charter check warns, without failing, of a planned entry fully tested, a deprecated one still in the code, and a doc to read again
var CH_261127_QRQS = &Spec{ID: "CH-261127-QRQS", Title: "charter check warns, without failing, of a planned entry fully tested, a deprecated one still in the code, and a doc to read again", Status: Active}
