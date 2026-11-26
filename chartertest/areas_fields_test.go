package chartertest

import (
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// The field lines may follow the heading after blank lines; the first blank
// line after a field ends them.
func TestFieldsMayFollowABlankLine(t *testing.T) {
	Verifies(t, claims.CH_010, "blank-lines-before-fields")
	src := strings.Replace(sampleArea, "## SW-010 The bucket refills from its last request at the rate\n", "## SW-010 The bucket refills from its last request at the rate\n\n", 1)
	_, entries, problems := parseArea("claims/refill.md", src)
	if len(problems) != 0 {
		t.Fatalf("problems: %v", problems)
	}
	if e := entries[0]; e.Status != "active" || e.Kind != "invariant" || !strings.HasPrefix(e.Description, "At each request") {
		t.Errorf("entry = %+v", e)
	}
}
