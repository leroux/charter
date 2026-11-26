package chartertest

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestReadListsEveryEntryWithCitations(t *testing.T) {
	root := copyFixture(t)
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Problems) != 0 {
		t.Errorf("Problems = %v, want none", reg.Problems)
	}
	got := reg.Entry("SW_001")
	if got == nil || got.ID != "SW-001" || got.File != "claims/specs/SW-001-bucket-survives-restart.md" || !reflect.DeepEqual(got.Realizes, []string{"STK-001"}) {
		t.Errorf("Entry(SW-001) = %+v", got)
	}
	if w := reg.Entry("SW_002"); w == nil || w.Status != "deprecated" {
		t.Errorf("Entry(SW-002) = %+v", w)
	}
	if l := reg.Entry("LockIsAdvisory"); l == nil || l.Status != "documented" || l.Source != "https://man7.org/linux/man-pages/man2/flock.2.html" {
		t.Errorf("Entry(LockIsAdvisory) = %+v", l)
	}
	if reg.Entry("Nope") != nil {
		t.Error("Entry(Nope) found an entry that does not exist")
	}
}

func TestReadReportsLoadProblemsWithoutFailing(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "claims/extra.go", "package claims\n\nfunc helper() {}\n")
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Problems) != 1 {
		t.Fatalf("Problems = %v, want the one rule 1 problem", reg.Problems)
	}
	if len(reg.Entries) != 10 {
		t.Errorf("got %d entries, want the 10 that still parse", len(reg.Entries))
	}
}

// The areas come in the charter's reading order, which a renderer walks a
// reader through, with the prose that opens each.
func TestReadListsTheAreaDocumentsAndTheWaivers(t *testing.T) {
	Verifies(t, claims.CH_040, "registry-lists-everything")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "Code: pkg/a\n\n", "Code: pkg/a\n\nHow a run takes up\nan edit.\n\n", 1))
	writeFile(t, root, "claims/waits.md", "# Waits\n\nHow a run waits.\n")
	writeFile(t, root, "claims/reading-order.txt", "# waits before refill\nwaits.md\nrefill.md\n")
	writeFile(t, root, "claims/waivers.txt", "# why\nSW-010/burst-caps-tokens\tverified by review\nSW-011\n")
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	want := []Area{
		{File: "claims/waits.md", Title: "Waits", Summary: "How a run waits."},
		{File: "claims/refill.md", Title: "Refill", Code: []string{"pkg/a"}, Summary: "How a run takes up an edit."},
	}
	if !reflect.DeepEqual(reg.Areas, want) {
		t.Errorf("Areas = %+v, want %+v", reg.Areas, want)
	}
	wantW := []Waiver{{Pattern: "SW-010/burst-caps-tokens", Reason: "verified by review", File: "claims/waivers.txt", Line: 2}}
	if !reflect.DeepEqual(reg.Waivers, wantW) {
		t.Errorf("Waivers = %+v, want %+v", reg.Waivers, wantW)
	}
	if !slices.ContainsFunc(reg.Problems, func(p string) bool { return strings.HasPrefix(p, "claims/waivers.txt:3: ") }) {
		t.Errorf("Problems = %v, want the waiver with no reason", reg.Problems)
	}
}
