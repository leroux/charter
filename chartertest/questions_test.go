package chartertest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

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

// thirdQuestion is another open question, to follow sampleQuestions.
const thirdQuestion = "\n" +
	"## Q-261009-P3RX Does a request after a denied one spend a token?\n" +
	"Asked-by: agent, PR #453\n" +
	"Entry: SW-010/waited-request-passes\n" +
	"Provisional: yes\n"

// The question file holds what nobody ruled: a question is read with its
// parts, and one missing a part or about an entry the record does not hold
// is a form problem.
func TestTheQuestionFileHoldsWhatNobodyRuled(t *testing.T) {
	Verifies(t, claims.CH_261009_4B3B, "question-read", "unknown-entry-refused", "missing-parts-refused")
	pkg := newEntriesPackage()
	_, es, _ := parseArea("claims/refill.md", sampleArea)
	for _, e := range es {
		pkg.add(e)
	}
	qs, problems, err := parseQuestions("claims/questions.md", sampleQuestions, pkg)
	if err != nil || len(problems) != 0 {
		t.Fatalf("problems %q, err %v", problems, err)
	}
	if len(qs) != 2 {
		t.Fatalf("%d questions, want 2", len(qs))
	}
	q := qs[0]
	if q.ID != "Q-261009-M4TK" || q.Text != "Does a denied request re-read the file?" || q.AskedBy != "agent, PR #452" || q.Entry != "SW-010/burst-caps-tokens" || q.Provisional != "yes, as the text says the edge is the only check" || !q.Open() || q.Line != 3 {
		t.Errorf("question 1 = %+v", q)
	}
	if q := qs[1]; q.Answer != "2026-10-09, owner: yes, the whole file." || q.Open() {
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
		{head + "Asked-by: agent\nEntry: SW-010\nProvisional: yes\nLens: x\n", `"Lens" is not a field of a question`},
		{head + "Asked-by: agent\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n", "Asked-by appears twice"},
		{head + "Asked-by: agent\nEntry: SW-010\nProvisional: yes\n\n" + head + "Asked-by: agent\nEntry: SW-010\nProvisional: no\n", "Q-261009-M4TK is declared twice, at lines 1 and 6"},
		{"## Q-x Q?\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n", "is not an id"},
		{"## Questions\n", "a question's heading is \"## Q-ID question\""},
	} {
		_, problems, err := parseQuestions("claims/questions.md", tc.src, pkg)
		if err != nil {
			t.Fatal(err)
		}
		if !anyContains(problems, tc.want) {
			t.Errorf("want %q in %q for %q", tc.want, problems, tc.src)
		}
		for _, p := range problems {
			if !strings.HasPrefix(p, "claims/questions.md:") || !strings.Contains(p, ": form: ") {
				t.Errorf("problem %q does not name the file, the line and the form rule", p)
			}
		}
	}

	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", refillTest)
	writeFile(t, root, "claims/questions.md", sampleQuestions)
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a sound question file reported:\n%s", strings.Join(problems, "\n"))
	}
	writeFile(t, root, "claims/questions.md", "## Q-261009-M4TK Q?\nAsked-by: agent\nEntry: SW-999\nProvisional: yes\n")
	reports(t, root, "claims/questions.md:1: form: Q-261009-M4TK is about SW-999, which is not an entry")
}

// Lint says in one line, at the question file's first line, how many
// questions are open and which command lists them, whatever line the
// first open one is on, and says nothing of the questions when none is
// open.
func TestLintCountsTheOpenQuestionsInOneLine(t *testing.T) {
	Verifies(t, claims.CH_261009_4B3B, "open-listed")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", refillTest)
	answeredOnly := sampleQuestions[strings.Index(sampleQuestions, "## Q-261009-N7PR"):]
	for _, tc := range []struct {
		name, file string
		want       []string
	}{
		{"one open", sampleQuestions, []string{"claims/questions.md:1: lint: 1 question is open; go run ./cmd/claimsgen questions lists it"}},
		{"two open", sampleQuestions + thirdQuestion, []string{"claims/questions.md:1: lint: 2 questions are open; go run ./cmd/claimsgen questions lists them"}},
		{"none open", answeredOnly, nil},
	} {
		writeFile(t, root, "claims/questions.md", tc.file)
		lines, err := Lint(fixtureOptions(root))
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, l := range lines {
			if strings.HasPrefix(l, "claims/questions.md:") {
				got = append(got, l)
			}
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("%s: lint says of the question file %q, want %q", tc.name, got, tc.want)
		}
	}
}

// Questions reads the question file for the questions verb: every question
// in the order the file holds them, the open and the answered with their
// parts, none for a charter that keeps no file, and an error, never none,
// when the file cannot be read.
func TestQuestionsGivesEveryQuestionInFileOrder(t *testing.T) {
	Verifies(t, claims.CH_261009_4B3B, "questions-verb-lists")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	if qs, err := Questions(fixtureOptions(root)); err != nil || len(qs) != 0 {
		t.Fatalf("a charter without a question file gave %+v, %v", qs, err)
	}
	if qs, err := Questions(Options{Root: root}); err == nil {
		t.Errorf("options naming no entries package gave %+v and no error", qs)
	}
	if err := os.Mkdir(filepath.Join(root, "claims", "questions.md"), 0o750); err != nil {
		t.Fatal(err)
	}
	if qs, err := Questions(fixtureOptions(root)); err == nil {
		t.Errorf("a question file that cannot be read gave %+v and no error", qs)
	}
	if err := os.Remove(filepath.Join(root, "claims", "questions.md")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "claims/questions.md", sampleQuestions+thirdQuestion)
	qs, err := Questions(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, q := range qs {
		state := "open"
		if !q.Open() {
			state = "answered"
		}
		got = append(got, q.ID+" "+q.Entry+" "+state)
	}
	want := []string{
		"Q-261009-M4TK SW-010/burst-caps-tokens open",
		"Q-261009-N7PR SW-010 answered",
		"Q-261009-P3RX SW-010/waited-request-passes open",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("questions %q, want %q", got, want)
	}
	if q := qs[0]; q.Text != "Does a denied request re-read the file?" || q.AskedBy != "agent, PR #452" || q.Provisional != "yes, as the text says the edge is the only check" {
		t.Errorf("the open question = %+v", q)
	}
	if q := qs[1]; q.Answer != "2026-10-09, owner: yes, the whole file." {
		t.Errorf("the answered question = %+v", q)
	}
}
