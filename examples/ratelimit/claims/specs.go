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

// STK_261126_H2SW is STK-261126-H2SW: ratelimit MUST let a quiet caller burst, then hold it to its sustained rate, and MUST NOT let one caller spend another's requests
var STK_261126_H2SW = &Spec{ID: "STK-261126-H2SW", Title: "ratelimit MUST let a quiet caller burst, then hold it to its sustained rate, and MUST NOT let one caller spend another's requests", Status: Active}

// STK_261126_T2BD is STK-261126-T2BD: ratelimit MUST take its rate and burst from the operator at start, and MUST NOT start on a setting it cannot read
var STK_261126_T2BD = &Spec{ID: "STK-261126-T2BD", Title: "ratelimit MUST take its rate and burst from the operator at start, and MUST NOT start on a setting it cannot read", Status: Active}

// SW_261126_2XG8 is SW-261126-2XG8: The limiter MUST read the time from the Clock it is given and nowhere else
var SW_261126_2XG8 = &Spec{ID: "SW-261126-2XG8", Title: "The limiter MUST read the time from the Clock it is given and nowhere else", Status: Active, Kind: "structure"}

// SW_261126_39F6 is SW-261126-39F6: The limiter MUST NOT let one key's requests change another key's bucket
var SW_261126_39F6 = &Spec{ID: "SW-261126-39F6", Title: "The limiter MUST NOT let one key's requests change another key's bucket", Status: Active, Kind: "invariant"}

// SW_261126_CMTG is SW-261126-CMTG: The command MUST take the rate and the burst from the environment at start, and MUST NOT start on one it cannot read
var SW_261126_CMTG = &Spec{ID: "SW-261126-CMTG", Title: "The command MUST take the rate and the burst from the environment at start, and MUST NOT start on one it cannot read", Status: Active}

// SW_261126_VX8Z is SW-261126-VX8Z: The limiter MUST give a new key a full bucket that earns the rate, holds at most the burst and pays a token for each request let through
var SW_261126_VX8Z = &Spec{ID: "SW-261126-VX8Z", Title: "The limiter MUST give a new key a full bucket that earns the rate, holds at most the burst and pays a token for each request let through", Status: Active, Kind: "bound"}
