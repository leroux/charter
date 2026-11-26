package chartertest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// The digest fixture: two areas read b before a, three needs, and
// decisions that reach the needs directly and through one another.
const digestNeeds = "# What people need\n\n" +
	"## STK-010 A person keeps their work\n\nStatus: active\n\nA need.\n\nWhy: the need's why.\n\nScenarios:\n1. `kept`: it is kept.\n\n" +
	"## STK-011 A person sees their work\n\nStatus: active\n\nA need.\n\nWhy: the need's why.\n\nScenarios:\n1. `seen`: it is seen.\n\n" +
	"## STK-012 A person shares their work\n\nStatus: active\n\nA need.\n\nWhy: the need's why.\n\nScenarios:\n1. `shared`: it is shared.\n"

func decision(id, title, status, kind, serves, why string) string {
	s := "## " + id + " " + title + "\n\nStatus: " + status + "\n"
	if kind != "" {
		s += "Kind: " + kind + "\n"
	}
	return s + "Serves: " + serves + "\n\nThe decision.\n\nWhy: " + why + "\n\nScenarios:\n1. `holds`: it holds.\n\n"
}

var digestAreaA = "# Storage: where work is kept\nCode: pkg/a\n\nThe log keeps the work. It is read back on a restart.\n\n" +
	decision("SW-010", "The log keeps every event", "active", "invariant", "STK-010", "a lost event is lost work. It cannot come back.") +
	decision("SW-011", "A snapshot never outranks the log", "active", "", "SW-010", "a stale cache shows a state that never was.")

var digestAreaB = "# Views: what a person sees\nCode: pkg/b\n\nWhat the dashboard shows.\n\n" +
	decision("SW-020", "The view folds the log", "active", "", "STK-010, STK-011", "e.g. a refill shows the same bucket. That is the point.") +
	decision("SW-021", "A shared link opens the bucket", "planned", "", "STK-012", "not yet.") +
	decision("SW-022", "The old view", "deprecated", "", "STK-011", "gone.")

// digestRoot is the fixture with the digest's documents and options.
func digestRoot(t *testing.T) (string, Options) {
	t.Helper()
	root := copyFixture(t)
	writeFile(t, root, "claims/needs.md", digestNeeds)
	writeFile(t, root, "claims/a.md", digestAreaA)
	writeFile(t, root, "claims/b.md", digestAreaB)
	writeFile(t, root, "claims/reading-order.txt", "# b first\nb.md\na.md\n")
	opts := fixtureOptions(root)
	opts.Digest = "DIGEST.md"
	return root, opts
}

func digestOf(t *testing.T, opts Options) string {
	t.Helper()
	d, err := Digest(opts)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// section is the digest from one heading to the next of its level.
func section(d, heading string) string {
	_, rest, ok := strings.Cut(d, heading+"\n")
	if !ok {
		return ""
	}
	level := heading[:strings.Index(heading, " ")+1]
	if i := strings.Index(rest, "\n"+level); i >= 0 {
		return rest[:i]
	}
	return rest
}

func inOrder(t *testing.T, text string, parts ...string) {
	t.Helper()
	at := 0
	for _, p := range parts {
		i := strings.Index(text[at:], p)
		if i < 0 {
			t.Fatalf("%q not found after offset %d in:\n%s", p, at, text)
		}
		at += i + len(p)
	}
}

func TestTheDigestListsEachNeedUnderTheAreaThatServesItMost(t *testing.T) {
	Verifies(t, claims.CH_261005_6ZBS, "needs-under-their-area")
	_, opts := digestRoot(t)
	needs := section(digestOf(t, opts), "## What people need")
	// STK-010 is reached by SW-010 and, through it, SW-011 in a (2), and
	// by SW-020 in b (1): it sits under a. STK-011 only b reaches. STK-012
	// only a planned decision serves, so no active one reaches it.
	inOrder(t, needs, "**Views: what a person sees**", "STK-011", "**Storage: where work is kept**", "STK-010", "**Needs no active decision serves yet**", "STK-012")
	if strings.Count(needs, "[STK-010]") != 1 {
		t.Errorf("STK-010 listed more than once:\n%s", needs)
	}
	// On a tie the earlier area in the reading order wins.
	root, opts := digestRoot(t)
	writeFile(t, root, "claims/a.md", strings.Replace(digestAreaA, "Serves: SW-010", "Serves: STK-011", 1))
	writeFile(t, root, "claims/b.md", strings.Replace(digestAreaB, "Serves: STK-010, STK-011", "Serves: STK-010, STK-011, STK-012", 1))
	needs = section(digestOf(t, opts), "## What people need")
	// STK-010: a 1, b 1; STK-011: a 1, b 1. b reads first.
	views := section(needs+"\n**", "**Views: what a person sees**")
	if !strings.Contains(views, "STK-010") || !strings.Contains(views, "STK-011") {
		t.Errorf("on a tie the needs belong to b, the earlier area:\n%s", needs)
	}
}

func TestTheDigestMapsTheAreasAndTheNeedsTheyServe(t *testing.T) {
	Verifies(t, claims.CH_261005_6ZBS, "area-map")
	_, opts := digestRoot(t)
	areas := section(digestOf(t, opts), "## The areas")
	inOrder(t, areas,
		"| 1 | [Views: what a person sees](b.md) | What the dashboard shows. | 1, 1 planned |",
		"| 2 | [Storage: where work is kept](a.md) | The log keeps the work. | 2 |",
		"[STK-010](needs.md#stk-010-a-person-keeps-their-work) | b 1 · a 2 |",
		"[STK-011](needs.md#stk-011-a-person-sees-their-work) | b 1 |",
		"[STK-012](needs.md#stk-012-a-person-shares-their-work) | none yet |")
}

func TestTheDigestReadsEachAreasDecisionsInOrderWithTheirWhy(t *testing.T) {
	Verifies(t, claims.CH_261005_6ZBS, "decisions-in-order")
	_, opts := digestRoot(t)
	d := section(digestOf(t, opts), "## The decisions, area by area")
	inOrder(t, d,
		"### 1. Views: what a person sees", "What the dashboard shows.",
		"**[SW-020](b.md#sw-020-the-view-folds-the-log)** The view folds the log", "Why: e.g. a refill shows the same bucket.",
		"Not shown: 1 planned, 1 deprecated.",
		"### 2. Storage: where work is kept",
		"**[SW-010](a.md#sw-010-the-log-keeps-every-event)** The log keeps every event *(invariant)*", "Why: a lost event is lost work.\n",
		"**[SW-011](a.md#sw-011-a-snapshot-never-outranks-the-log)** A snapshot never outranks the log")
	for _, absent := range []string{"SW-021", "SW-022", "That is the point", "It cannot come back"} {
		if strings.Contains(d, absent) {
			t.Errorf("the decisions part holds %q: a planned or deprecated entry, or a Why past its first sentence", absent)
		}
	}
}

func TestTheReadingOrderNamesEveryAreaOnce(t *testing.T) {
	Verifies(t, claims.CH_261005_6ZBS, "reading-order-complete")
	root, opts := digestRoot(t)
	writeFile(t, root, "claims/reading-order.txt", "b.md\nc.md\nb.md\n")
	got := digestProblems(t, opts)
	for _, want := range []string{
		`claims/reading-order.txt:2: rule 26: "c.md" is no area document, or is listed twice; remove the line`,
		`claims/reading-order.txt:3: rule 26: "b.md" is no area document, or is listed twice; remove the line`,
		"claims/reading-order.txt: rule 26: the reading order leaves out a.md; add it where a reader should meet it",
	} {
		if !anyContains(got, want) {
			t.Errorf("problems lack %q:\n%s", want, strings.Join(got, "\n"))
		}
	}
	// The area left out still reads, last.
	inOrder(t, digestOf(t, opts), "### 1. Views", "### 2. Storage")
}

func TestAStaleDigestFailsTheGateAndGenerationWritesIt(t *testing.T) {
	Verifies(t, claims.CH_261005_6ZBS, "stale-digest-fails")
	root, opts := digestRoot(t)
	if got := digestProblems(t, opts); !anyContains(got, "claims/DIGEST.md: rule 26: the digest is missing; run CHARTER_UPDATE_DIGEST=1 go test -count=1 ./claims") {
		t.Errorf("a missing digest: %v", got)
	}
	if err := Generate(opts); err != nil {
		t.Fatal(err)
	}
	if got := digestProblems(t, opts); len(got) != 0 {
		t.Errorf("after generation: %v", got)
	}
	written, err := os.ReadFile(filepath.Join(root, "claims", "DIGEST.md"))
	if err != nil || string(written) != digestOf(t, opts) {
		t.Fatalf("generation wrote %q (%v), want the rendering", written, err)
	}
	replaceIn(t, root, "claims/a.md", "The log keeps every event", "The log keeps each event")
	if got := digestProblems(t, opts); !anyContains(got, "claims/DIGEST.md: rule 26: the digest differs from the documents") {
		t.Errorf("a stale digest: %v", got)
	}
	opts.Regenerate = "make claims"
	if got := digestProblems(t, opts); !anyContains(got, "differs from the documents; run make claims") {
		t.Errorf("the command named: %v", got)
	}
}

// digestProblems are the gate's rule 26 problems.
func digestProblems(t *testing.T, opts Options) []string {
	t.Helper()
	all, err := Problems(opts)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, p := range all {
		if strings.Contains(p, "rule 26") {
			out = append(out, p)
		}
	}
	return out
}
