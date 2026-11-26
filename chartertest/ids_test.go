package chartertest

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/leroux/charter/claims"
)

// lateInTheWest is an instant whose date west of Greenwich, the first of
// October, is not its UTC date, the second.
var lateInTheWest = time.Date(2026, 10, 1, 22, 0, 0, 0, time.FixedZone("west", -5*60*60))

// scripted is a random source that gives the draws in order, each an index
// into the code alphabet, and fails the test when asked for more.
func scripted(t *testing.T, draws ...int) func(int) int {
	t.Helper()
	i := 0
	return func(n int) int {
		t.Helper()
		if n != len(codeAlphabet) {
			t.Fatalf("drew over %d, want the alphabet's %d", n, len(codeAlphabet))
		}
		if i == len(draws) {
			t.Fatalf("drew more than the %d scripted", len(draws))
		}
		i++
		return draws[i-1]
	}
}

// plannedEntry is a well-formed entry under the id, with nothing to cover.
func plannedEntry(id string) string {
	return "## " + id + " An entry\nStatus: planned\n\nWhat it rules.\n\nWhy: someone needed it.\n\n"
}

// An id is numbered or dated; a dated id's day is on the calendar and its
// code is four characters of the alphabet.
func TestAnIDIsNumberedOrDated(t *testing.T) {
	Verifies(t, claims.CH_012, "id-form")
	for id, want := range map[string]entry{
		"SW-010":         {Name: "SW_010", Lens: "SW", Number: 10},
		"SW-261001-K7Q2": {Name: "SW_261001_K7Q2", Lens: "SW", Day: "261001", Code: "K7Q2"},
		"CH-280229-ZY98": {Name: "CH_280229_ZY98", Lens: "CH", Day: "280229", Code: "ZY98"},
	} {
		_, entries, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## "+id, 1))
		if len(problems) != 0 || len(entries) != 1 {
			t.Errorf("%s: %d entries, problems %v", id, len(entries), problems)
			continue
		}
		e := entries[0]
		if e.ID != id || e.Name != want.Name || e.Lens != want.Lens || e.Number != want.Number || e.Day != want.Day || e.Code != want.Code {
			t.Errorf("%s parsed as %q: symbol %s, prefix %s, number %d, day %q, code %q", id, e.ID, e.Name, e.Lens, e.Number, e.Day, e.Code)
		}
	}
	for id, want := range map[string]string{
		"SW-261332-K7Q2":  "the date of SW-261332-K7Q2, 261332, is not a day on the calendar written YYMMDD",
		"SW-260229-K7Q2":  "the date of SW-260229-K7Q2, 260229, is not a day on the calendar",
		"SW-2610-K7Q2":    "the date of SW-2610-K7Q2, 2610, is not a day on the calendar",
		"SW-261001-K7Q":   "the code of SW-261001-K7Q, K7Q, is not 4 characters of BCDFGHJKMNPQRSTVWXYZ23456789",
		"SW-261001-K7Q22": "the code of SW-261001-K7Q22, K7Q22, is not 4 characters",
		"SW-261001-KAQ2":  "the code of SW-261001-KAQ2, KAQ2, is not 4 characters",
		"SW-261001-KLQ2":  "the code of SW-261001-KLQ2, KLQ2, is not 4 characters",
		"SW-261001-K0Q2":  "the code of SW-261001-K0Q2, K0Q2, is not 4 characters",
		"SW-261001-K1Q2":  "the code of SW-261001-K1Q2, K1Q2, is not 4 characters",
		"SW-261001-k7q2":  "the code of SW-261001-k7q2, k7q2, is not 4 characters",
		"SW-DRAFT":        `"SW-DRAFT" is not an id: an id is a prefix and a number, SW-118, or a prefix, the day it was taken and a code, SW-261001-K7Q2`,
	} {
		_, entries, problems := parseArea("claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## "+id, 1))
		if len(entries) != 0 {
			t.Errorf("%s yielded an entry", id)
		}
		if !anyContains(problems, "claims/refill.md:4: form: "+want) {
			t.Errorf("%s: want %q in:\n%s", id, want, strings.Join(problems, "\n"))
		}
	}
}

// A dated id is anchored, cited from Go and from TypeScript and waived as a
// numbered one is, and the gate has nothing to say.
func TestADatedIDPassesTheGate(t *testing.T) {
	Verifies(t, claims.CH_012, "id-form")
	dated := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "SW-010", "SW-261001-K7Q2"), "SW_010", "SW_261001_K7Q2")
	}
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", dated(sampleArea))
	writeFile(t, root, "pkg/a/refill.go", dated(refillGo))
	writeFile(t, root, "pkg/a/refill_test.go", dated(strings.Replace(refillTest, "\"burst-caps-tokens\"", "\"waited-request-passes\"", 1)))
	writeFile(t, root, "claims/waivers.txt", "SW-261001-K7Q2/burst-caps-tokens\tthe burst case is verified by review\n")
	writeFile(t, root, "web/refill.test.ts", strings.Join([]string{
		`import { it } from "@example/test";`,
		`import { claims } from "@fixture/claims";`,
		`import { verifies } from "./helper";`,
		``,
		`it("inserts", () => { verifies(claims.SW_261001_K7Q2, "waited-request-passes"); });`,
		``,
	}, "\n"))
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a dated id reported:\n%s", strings.Join(problems, "\n"))
	}
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	e := reg.Entry("SW_261001_K7Q2")
	if e == nil || e.ID != "SW-261001-K7Q2" || len(e.Anchors) != 1 || len(e.Refs) != 3 {
		t.Fatalf("the registry's entry = %+v, want one anchor and three citations", e)
	}
	writeFile(t, root, "pkg/a/refill.go", "package a\n")
	reports(t, root, "rule 25: active spec SW-261001-K7Q2 is realized by no code")
}

// Within a prefix the numbered ids come first, by number, then the dated
// ones by day and then by code; the generated Go follows.
func TestEntriesAreOrderedByNumberThenByDayAndCode(t *testing.T) {
	Verifies(t, claims.CH_012, "ordered-by-date")
	root := t.TempDir()
	var doc strings.Builder
	doc.WriteString("# Order\n\n")
	for _, id := range []string{"SW-261001-K7Q2", "SW-010", "SW-260930-ZZZZ", "CON-261001-B222", "SW-261001-B222", "SW-009", "SW-261001-2222", "CON-007"} {
		doc.WriteString(plannedEntry(id))
	}
	writeFile(t, root, "claims/order.md", doc.String())
	pkg := newEntriesPackage()
	if _, problems, err := readAreas(root, "claims", pkg); err != nil || len(problems) != 0 {
		t.Fatalf("readAreas: %v, problems %v", err, problems)
	}
	want := []string{"CON_007", "CON_261001_B222", "SW_009", "SW_010", "SW_260930_ZZZZ", "SW_261001_2222", "SW_261001_B222", "SW_261001_K7Q2"}
	if !slices.Equal(pkg.order, want) {
		t.Errorf("order = %v, want %v", pkg.order, want)
	}
	src, err := renderSpecsGo(pkg, "claims", "claims", "make claims")
	if err != nil {
		t.Fatal(err)
	}
	at := -1
	for _, name := range want {
		decl := "\nvar " + name + " = &charter.Spec{ID: \"" + strings.ReplaceAll(name, "_", "-") + "\","
		next := strings.Index(string(src), decl)
		if next <= at {
			t.Errorf("%s is declared at %d, want after %d:\n%s", name, next, at, src)
		}
		at = next
	}
	module := renderTypeScript(pkg, "claims", "make claims")
	at = -1
	for _, name := range want {
		next := strings.Index(module, "\n  "+name+": \""+strings.ReplaceAll(name, "_", "-")+"\",")
		if next <= at {
			t.Errorf("the module maps %s at %d, want after %d:\n%s", name, next, at, module)
		}
		at = next
	}
}

// A new id is its prefix, the UTC date and four characters of the alphabet,
// and a draw the tree holds is drawn again.
func TestANewIDIsItsDayAndACodeTheTreeDoesNotHold(t *testing.T) {
	Verifies(t, claims.CH_040, "next-free")
	Verifies(t, claims.CH_003, "next-free-id")
	pkg := newEntriesPackage()
	pkg.add(&entry{ID: "SW-261002-BBBB", Name: "SW_261002_BBBB", Level: levelSpec, Lens: "SW", Day: "261002", Code: "BBBB"})
	for _, tc := range []struct {
		name, prefix string
		draws        []int
		want         string
	}{
		{"the UTC date and the code drawn", "SW", []int{7, 25, 11, 20}, "SW-261002-K7Q2"},
		{"a draw the tree holds is drawn again", "SW", []int{0, 0, 0, 0, 1, 1, 1, 1}, "SW-261002-CCCC"},
		{"a prefix nothing uses", "ZZ", []int{0, 0, 0, 0}, "ZZ-261002-BBBB"},
		{"the ends of the alphabet", "STK", []int{0, 27, 19, 20}, "STK-261002-B9Z2"},
	} {
		got := newIDs(pkg, tc.prefix, 1, lateInTheWest, scripted(t, tc.draws...))
		if len(got) != 1 || got[0] != tc.want {
			t.Errorf("%s: newIDs = %v, want %s", tc.name, got, tc.want)
		}
	}

	// From a tree: what its documents hold is what a draw may not give. The
	// first draw is the id the area document declares, so the second is
	// taken; once the document no longer declares it, the first is.
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "## SW-010", "## SW-261002-BBBB", 1))
	opts := fixtureOptions(root)
	got, err := nextIDs(opts, ".", "SW", 1, lateInTheWest, scripted(t, 0, 0, 0, 0, 1, 1, 1, 1))
	if err != nil || len(got) != 1 || got[0] != "SW-261002-CCCC" {
		t.Errorf("nextIDs over a tree declaring SW-261002-BBBB = %v, %v; want SW-261002-CCCC, the second draw", got, err)
	}
	writeFile(t, root, "claims/refill.md", sampleArea)
	got, err = nextIDs(opts, ".", "SW", 1, lateInTheWest, scripted(t, 0, 0, 0, 0))
	if err != nil || len(got) != 1 || got[0] != "SW-261002-BBBB" {
		t.Errorf("nextIDs over a tree that no longer declares it = %v, %v; want SW-261002-BBBB, the first draw", got, err)
	}
}

// Next, over the tree, gives an id of today's UTC date that the charter
// does not hold, and refuses a prefix that is not capital letters.
func TestNextGivesAnIDTheCharterDoesNotHold(t *testing.T) {
	Verifies(t, claims.CH_040, "next-free")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	opts := fixtureOptions(root)
	before := time.Now().UTC().Format(dayLayout)
	ids, err := NextIDs(opts, "SW", 1)
	after := time.Now().UTC().Format(dayLayout)
	if err != nil || len(ids) != 1 {
		t.Fatalf("NextIDs = %v, %v", ids, err)
	}
	parts, problem := parseID(ids[0])
	if problem != "" || parts.Lens != "SW" || parts.Code == "" || (parts.Day != before && parts.Day != after) {
		t.Errorf("NextIDs gave %s: %+v, %s; want SW, today's UTC date and a code", ids[0], parts, problem)
	}
	held := newEntriesPackage()
	if _, _, err := readAreas(root, "claims", held); err != nil {
		t.Fatal(err)
	}
	if _, err := readSpecs(root, "claims/specs", held); err != nil {
		t.Fatal(err)
	}
	if held.byID["SW-010"] == nil || held.byID["SW-001"] == nil || held.byID[ids[0]] != nil {
		t.Errorf("the tree holds SW-010 %v, SW-001 %v, %s %v; want the first two only", held.byID["SW-010"] != nil, held.byID["SW-001"] != nil, ids[0], held.byID[ids[0]] != nil)
	}
	if ids, err := NextIDs(opts, "sw", 1); err == nil {
		t.Errorf("NextIDs took the prefix sw and gave %v", ids)
	}
}

// Next gives as many ids as it is asked for, each different, and refuses a
// count that is not above zero.
func TestNextGivesAsManyDifferentIDsAsAsked(t *testing.T) {
	Verifies(t, claims.CH_040, "next-count")
	got := newIDs(newEntriesPackage(), "SW", 2, lateInTheWest, scripted(t, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1))
	if want := []string{"SW-261002-BBBB", "SW-261002-CCCC"}; !slices.Equal(got, want) {
		t.Errorf("newIDs = %v, want %v: a code drawn twice is drawn again", got, want)
	}
	opts := fixtureOptions(copyFixture(t))
	ids, err := NextIDs(opts, "SW", 5)
	if err != nil || len(ids) != 5 {
		t.Fatalf("NextIDs = %v, %v; want five ids", ids, err)
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if _, problem := parseID(id); problem != "" || seen[id] {
			t.Errorf("of %v, %s is repeated or no id: %s", ids, id, problem)
		}
		seen[id] = true
	}
	for _, count := range []int{0, -1} {
		if ids, err := NextIDs(opts, "SW", count); err == nil {
			t.Errorf("NextIDs with a count of %d gave %v", count, ids)
		}
	}
}
