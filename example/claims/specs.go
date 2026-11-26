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

// STK_261126_H2SW is STK-261126-H2SW: A caller gets a burst at once, then its sustained rate, and no more
var STK_261126_H2SW = &Spec{ID: "STK-261126-H2SW", Title: "A caller gets a burst at once, then its sustained rate, and no more", Status: Active}

// STK_261126_T2BD is STK-261126-T2BD: An operator sets the rate and the burst at start, and a bad setting stops the start
var STK_261126_T2BD = &Spec{ID: "STK-261126-T2BD", Title: "An operator sets the rate and the burst at start, and a bad setting stops the start", Status: Active}

// SW_261126_2XG8 is SW-261126-2XG8: The limiter reads time through its Clock alone
var SW_261126_2XG8 = &Spec{ID: "SW-261126-2XG8", Title: "The limiter reads time through its Clock alone", Status: Active, Kind: "structure"}

// SW_261126_39F6 is SW-261126-39F6: One key's requests never change another key's bucket
var SW_261126_39F6 = &Spec{ID: "SW-261126-39F6", Title: "One key's requests never change another key's bucket", Status: Active, Kind: "invariant"}

// SW_261126_CMTG is SW-261126-CMTG: The rate and the burst come from the environment at start, and an unreadable one stops the start
var SW_261126_CMTG = &Spec{ID: "SW-261126-CMTG", Title: "The rate and the burst come from the environment at start, and an unreadable one stops the start", Status: Active}

// SW_261126_VX8Z is SW-261126-VX8Z: A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request
var SW_261126_VX8Z = &Spec{ID: "SW-261126-VX8Z", Title: "A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request", Status: Active, Kind: "bound"}
