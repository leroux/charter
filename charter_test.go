package charter_test

import (
	"testing"

	"github.com/leroux/charter"
)

func TestStatusString(t *testing.T) {
	cases := map[charter.Status]string{
		charter.Planned: "planned", charter.Active: "active", charter.Deprecated: "deprecated",
		charter.Status(9): "Status(9)",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}

func TestObservationStatusString(t *testing.T) {
	cases := map[charter.ObservationStatus]string{
		charter.Observed: "observed", charter.Documented: "documented",
		charter.Broken: "broken", charter.ObservationUnsure: "unsure",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", int(s), got, want)
		}
	}
	if got := charter.ObservationStatus(9).String(); got != "ObservationStatus(9)" {
		t.Errorf("unknown status = %q", got)
	}
}
