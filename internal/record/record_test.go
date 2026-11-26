package record

import (
	"os"
	"path/filepath"
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

// writeRecord writes files into root/claims and reads the record there.
func writeRecord(t *testing.T, files map[string]string) *Record {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "claims")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r, err := Read(root, "claims")
	if err != nil {
		t.Fatal(err)
	}
	return r
}

const sampleQuestions = "# Questions\n\n" +
	"## Q-261009-M4TK Does a denied request re-read the file?\n" +
	"Asked-by: agent, PR #452\n" +
	"Entry: SW-010/burst-caps-tokens\n" +
	"Provisional: yes, as the text says\n" +
	"  the edge is the only check\n" +
	"\n" +
	"## Q-261009-N7PR Is the hash the whole file?\n" +
	"Asked-by: agent, PR #450\n" +
	"Entry: SW-010\n" +
	"Provisional: yes\n" +
	"Answer: 2026-10-09, owner: yes, the whole file.\n"

func TestTheQuestionFileHoldsWhatNobodyRuled(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_4B3B, "question-read", "unknown-entry-refused", "missing-parts-refused")
	r := writeRecord(t, map[string]string{"refill.md": sampleArea, "questions.md": sampleQuestions})
	if len(r.Problems) != 0 {
		t.Fatalf("problems %q", r.Problems)
	}
	if len(r.Questions) != 2 {
		t.Fatalf("%d questions, want 2", len(r.Questions))
	}
	q := r.Questions[0]
	if q.ID != "Q-261009-M4TK" || q.Text != "Does a denied request re-read the file?" || q.AskedBy != "agent, PR #452" || q.Entry != "SW-010/burst-caps-tokens" || q.Provisional != "yes, as the text says the edge is the only check" || !q.Open() || q.Line != 3 {
		t.Errorf("question 1 = %+v", q)
	}
	if q := r.Questions[1]; q.Answer != "2026-10-09, owner: yes, the whole file." || q.Open() {
		t.Errorf("question 2 = %+v", q)
	}
	head := "## Q-261009-M4TK Q?\n"
	for _, tc := range []struct{ src, want string }{
		{head + "Asked-by: agent\nEntry: SW-999\nProvisional: yes\n", "Q-261009-M4TK is about SW-999, which is not an entry"},
		{head + "Asked-by: agent\nEntry: SW-010/nope\nProvisional: yes\n", "Q-261009-M4TK is about SW-010/nope, a scenario SW-010 does not list"},
		{head + "Entry: SW-010\nProvisional: yes\n", "Q-261009-M4TK has no Asked-by"},
		{head + "Asked-by: agent\nProvisional: yes\n", "Q-261009-M4TK has no Entry"},
		{head + "Asked-by: agent\nEntry: SW-010\n", "Q-261009-M4TK has no Provisional"},
		{head + "Asked-by: agent\nEntry: SW-010\nProvisional: yes\nAnswer: yes\n", "an Answer opens with the date and who answered"},
		{head + "Asked-by: agent\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n", "Asked-by appears twice"},
		{head + "Asked-by: agent\nEntry: SW-010\nProvisional: yes\n\n" + head + "Asked-by: agent\nEntry: SW-010\nProvisional: no\n", "Q-261009-M4TK is declared twice, at lines 1 and 6"},
	} {
		r := writeRecord(t, map[string]string{"refill.md": sampleArea, "questions.md": tc.src})
		if !anyContains(r.Problems, "claims/questions.md:") || !anyContains(r.Problems, tc.want) {
			t.Errorf("want %q at claims/questions.md in %q", tc.want, r.Problems)
		}
	}
}

func TestANewIDIsItsDayAndACodeTheRecordDoesNotHold(t *testing.T) {
	claims.Verifies(t, claims.CH_040, "next-free")
	claims.Verifies(t, claims.CH_003, "next-free-id")
	r := writeRecord(t, map[string]string{"refill.md": strings.Replace(sampleArea, "## SW-010", "## SW-261002-BBBB", 1)})
	for _, tc := range []struct {
		name, prefix string
		draws        []int
		want         string
	}{
		{"the UTC date and the code drawn", "SW", []int{7, 25, 11, 20}, "SW-261002-K7Q2"},
		{"a draw the record holds is drawn again", "SW", []int{0, 0, 0, 0, 1, 1, 1, 1}, "SW-261002-CCCC"},
		{"a prefix nothing uses", "ZZ", []int{0, 0, 0, 0}, "ZZ-261002-BBBB"},
		{"the ends of the alphabet", "STK", []int{0, 27, 19, 20}, "STK-261002-B9Z2"},
	} {
		got, err := r.NewIDs(tc.prefix, 1, lateInTheWest, scripted(t, tc.draws...))
		if err != nil || len(got) != 1 || got[0] != tc.want {
			t.Errorf("%s: NewIDs = %v, %v; want %s", tc.name, got, err, tc.want)
		}
	}
	if ids, err := r.NewIDs("sw", 1, lateInTheWest, scripted(t)); err == nil {
		t.Errorf("NewIDs took the prefix sw and gave %v", ids)
	}
}

func TestNextGivesAsManyDifferentIDsAsAsked(t *testing.T) {
	claims.Verifies(t, claims.CH_040, "next-count")
	r := writeRecord(t, nil)
	got, err := r.NewIDs("SW", 2, lateInTheWest, scripted(t, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1))
	if want := []string{"SW-261002-BBBB", "SW-261002-CCCC"}; err != nil || !slices.Equal(got, want) {
		t.Errorf("NewIDs = %v, %v; want %v: a code drawn twice is drawn again", got, err, want)
	}
	for _, count := range []int{0, -1} {
		if ids, err := r.NewIDs("SW", count, lateInTheWest, scripted(t)); err == nil {
			t.Errorf("NewIDs with a count of %d gave %v", count, ids)
		}
	}
}

func TestAnIdIsDeclaredOnce(t *testing.T) {
	claims.Verifies(t, claims.CH_012, "declared-once")
	claims.Verifies(t, claims.CH_003, "id-declared-once")
	again := strings.Replace(sampleArea, "# Refill", "# Again", 1)
	r := writeRecord(t, map[string]string{"refill.md": sampleArea, "zz.md": again})
	if !anyContains(r.Problems, "claims/zz.md:4: form: SW-010 is declared twice, in claims/refill.md and claims/zz.md; the entry that has not landed takes a new id: charter next SW") {
		t.Errorf("problems = %q", r.Problems)
	}
	if len(r.Entries) != 1 {
		t.Errorf("%d entries, want the first declaration only", len(r.Entries))
	}
}

const sampleObservations = "# Observations\n\n" +
	"## LockIsAdvisory\n\n" +
	"Subject: os\n" +
	"Status: documented\n" +
	"Source: https://man7.org/linux/man-pages/man2/flock.2.html\n" +
	"Supports: SW-010\n\n" +
	"flock locks are advisory:\n" +
	"a process that does not ask is not stopped.\n\n" +
	"## HealthIsJSON\n" +
	"Subject: upstream\n" +
	"Status: observed\n" +
	"Observed: 2026-09-24\n\n" +
	"The upstream health endpoint answers JSON.\n"

func TestAnObservationIsReadWithItsFields(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_HK6S, "observation-read")
	r := writeRecord(t, map[string]string{"refill.md": sampleArea, "observations.md": sampleObservations})
	if len(r.Problems) != 0 {
		t.Fatalf("problems %q", r.Problems)
	}
	if len(r.Observations) != 2 {
		t.Fatalf("%d observations, want 2", len(r.Observations))
	}
	o := r.Observation("LockIsAdvisory")
	if o == nil || o.Subject != "os" || o.Status != "documented" || o.Source != "https://man7.org/linux/man-pages/man2/flock.2.html" ||
		!slices.Equal(o.Supports, []string{"SW-010"}) || o.Text != "flock locks are advisory: a process that does not ask is not stopped." || o.Pos() != "claims/observations.md:3" {
		t.Errorf("LockIsAdvisory = %+v", o)
	}
	if o := r.Observation("HealthIsJSON"); o == nil || o.Observed != "2026-09-24" || o.Status != "observed" {
		t.Errorf("HealthIsJSON = %+v", o)
	}
	if _, o := r.BySymbol("HealthIsJSON"); o == nil {
		t.Error("BySymbol does not find an observation by its name")
	}
}

func TestAnObservationsFormIsChecked(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_HK6S, "observation-form-checked")
	head := "## LockIsAdvisory\nSubject: os\n"
	for _, tc := range []struct{ src, want string }{
		{head + "Status: documented\nVendor: linux\n\nText.\n", `claims/observations.md:4: form: "Vendor" is not a field of an observation`},
		{head + "Status: documented\nSubject: kernel\n\nText.\n", "claims/observations.md:4: form: Subject appears twice"},
		{head + "Status: guessed\n\nText.\n", `claims/observations.md:3: form: status "guessed" is not observed, documented, broken or unsure`},
		{head + "Status: unsure\n\nText.\n\n" + head + "Status: unsure\n\nAgain.\n", "claims/observations.md:7: form: LockIsAdvisory is declared twice, at lines 1 and 7"},
		{"## lock-is-advisory\nStatus: unsure\n", "claims/observations.md:1: form: an observation's heading is \"## Name\""},
	} {
		r := writeRecord(t, map[string]string{"observations.md": tc.src})
		if !anyContains(r.Problems, tc.want) {
			t.Errorf("want %q in %q", tc.want, r.Problems)
		}
	}
}

func TestTheReadingOrderNamesEveryAreaOnce(t *testing.T) {
	claims.Verifies(t, claims.CH_261005_6ZBS, "reading-order-complete")
	second := strings.Replace(strings.Replace(sampleArea, "# Refill", "# Limits", 1), "SW-010", "SW-011", 1)
	files := map[string]string{"refill.md": sampleArea, "limits.md": second}
	files["reading-order.txt"] = "# the order\nrefill.md\nlimits.md\n"
	if r := writeRecord(t, files); len(r.Problems) != 0 || !slices.Equal(r.ReadingOrder, []string{"claims/refill.md", "claims/limits.md"}) {
		t.Errorf("a complete order: problems %q, order %v", r.Problems, r.ReadingOrder)
	}
	for _, tc := range []struct{ order, want string }{
		{"refill.md\n", "claims/reading-order.txt:1: form: the reading order leaves out limits.md"},
		{"refill.md\nlimits.md\ngone.md\n", "claims/reading-order.txt:3: form: gone.md names no area document of claims"},
		{"refill.md\nlimits.md\nrefill.md\n", "claims/reading-order.txt:3: form: refill.md is named already, at line 1"},
	} {
		files["reading-order.txt"] = tc.order
		if r := writeRecord(t, files); !anyContains(r.Problems, tc.want) {
			t.Errorf("order %q: want %q in %q", tc.order, tc.want, r.Problems)
		}
	}
}

// The charter's own record reads with no form problem.
func TestTheChartersOwnRecordReadsClean(t *testing.T) {
	r, err := Read("../..", "claims", "DIGEST.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Problems) != 0 {
		t.Errorf("the charter's own record:\n%s", strings.Join(r.Problems, "\n"))
	}
	if r.Entry("CH-261126-HK6S") == nil || len(r.Questions) == 0 {
		t.Errorf("the record lacks CH-261126-HK6S or its questions")
	}
}

func TestEntriesAreOrderedByPrefixThenNumberThenDay(t *testing.T) {
	claims.Verifies(t, claims.CH_012, "ordered-by-date")
	entry := func(id string) string {
		return "## " + id + " An entry\nStatus: planned\n\nWhat it rules.\n\nWhy: someone needed it.\n\n"
	}
	r := writeRecord(t, map[string]string{
		"b.md": "# B\n\n" + entry("SW-261002-BBBB") + entry("SW-010") + entry("CH-002"),
		"a.md": "# A\n\n" + entry("SW-261001-ZZZZ") + entry("SW-002") + entry("SW-261002-BBBC"),
	})
	var got []string
	for _, e := range r.Entries {
		got = append(got, e.ID)
	}
	if want := []string{"CH-002", "SW-002", "SW-010", "SW-261001-ZZZZ", "SW-261002-BBBB", "SW-261002-BBBC"}; !slices.Equal(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

func TestBlankLinesMayOpenTheFieldsAndTheReadmeIsNoDocument(t *testing.T) {
	claims.Verifies(t, claims.CH_010, "blank-lines-before-fields", "readme-is-not-a-document")
	src := strings.Replace(sampleArea, "## SW-010 The bucket refills from its last request at the rate\n", "## SW-010 The bucket refills from its last request at the rate\n\n\n", 1)
	r := writeRecord(t, map[string]string{"refill.md": src, "README.md": "# How to write\n\n## Not an entry\n\nprose\n"})
	if len(r.Problems) != 0 || len(r.Entries) != 1 || r.Entries[0].Status != "active" || len(r.Areas) != 1 {
		t.Errorf("problems %q, %d entries, %d areas", r.Problems, len(r.Entries), len(r.Areas))
	}
}

func TestTheRecordGivesEachScenariosProvenance(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_3FHW, "provenance-in-the-registry")
	src := strings.Replace(sampleArea, "1. `waited-request-passes`: a request", "1. `waited-request-passes` (owner, #444): a request", 1)
	r := writeRecord(t, map[string]string{"refill.md": src})
	e := r.Entry("SW-010")
	if e == nil || e.Scenarios[0].Provenance != (Provenance{Who: "owner", Ref: "#444", Written: true}) || e.Scenarios[1].Provenance != (Provenance{Who: "agent"}) || e.Scenarios[1].Provenance.Written {
		t.Errorf("provenance: %+v", e)
	}
}
