package attrprobe

import (
	"testing"

	"github.com/leroux/charter"
	"github.com/leroux/charter/chartertest"
)

var probeEntry = &charter.Spec{ID: "SW-001", Title: "The probe runs.", Status: charter.Active}

var subEntry = &charter.Spec{ID: "SW-002", Title: "The probe's subtest is skipped.", Status: charter.Active}

func TestProbe(t *testing.T) {
	chartertest.Verifies(t, probeEntry)
	t.Run("sub", func(t *testing.T) {
		chartertest.Verifies(t, subEntry, "skipped")
		t.Skip("planned")
	})
}

var fact = &charter.Observation{ID: "ProbeFact", Text: "The probe's upstream answers.", Subject: "upstream", Status: charter.Documented, Source: "doc"}

func TestProbeObserves(t *testing.T) {
	chartertest.Observes(t, fact)
}
