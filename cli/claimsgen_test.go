package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
)

// lintTree writes a charter of its own under a temporary root, each file
// slash-relative to it, and returns the options lint reads it with.
func lintTree(t *testing.T, files map[string]string) chartertest.Options {
	t.Helper()
	root := t.TempDir()
	files["go.mod"] = "module example.com/fx\n\ngo 1.22\n"
	for rel, content := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return chartertest.Options{Root: root, Entries: "example.com/fx/claims"}
}

// lint -baseline over a charter whose scenario sentences over thirty words
// grew past their baseline, from one to three, names that count at its line
// in the file, with both numbers and what to do, and exits 1, where it
// exited 0 before they grew; the title over twenty words, at its baseline,
// is not named.
func TestLintBaselineNamesAGrownCountAndExitsNonZero(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261009_SNX2, "grown-count-named")
	long := strings.TrimSpace(strings.Repeat("word ", 31))
	doc := func(scenarios int) string {
		d := "# Style\n\n" +
			"## SW-001 A decision whose title runs to twenty-one words, one past the bar of twenty that the style of the record sets\n" +
			"Status: planned\n\nIt rules.\n\nWhy: someone needed it.\n\nScenarios:\n"
		for i := 1; i <= scenarios; i++ {
			d += fmt.Sprintf("%d. `s%d`: %s.\n", i, i, long)
		}
		return d
	}
	opts := lintTree(t, map[string]string{
		"claims/style.md": doc(1),
		"claims/style-baseline.txt": "# The counts the style checks are held to.\n" +
			"titles-over-twenty-words 1\n" +
			"scenario-sentences-over-thirty-words 1\n",
	})
	var before strings.Builder
	if code, err := runLint(&before, opts, []string{"-baseline"}); err != nil || code != 0 || before.Len() != 0 {
		t.Fatalf("at the baseline: lint -baseline printed %q, exit %d, %v; want nothing and exit 0", before.String(), code, err)
	}
	if err := os.WriteFile(filepath.Join(opts.Root, "claims", "style.md"), []byte(doc(3)), 0o600); err != nil {
		t.Fatal(err)
	}
	var counted strings.Builder
	if _, err := runLint(&counted, opts, nil); err != nil || !strings.Contains(counted.String(), "lint: 1 title runs past twenty words") ||
		!strings.Contains(counted.String(), "lint: 3 scenario sentences run past thirty words") {
		t.Fatalf("lint should count one title and three scenario sentences past the bars; it printed\n%s%v", counted.String(), err)
	}
	var out strings.Builder
	code, err := runLint(&out, opts, []string{"-baseline"})
	want := "claims/style-baseline.txt:3: lint: scenario sentences over thirty words grew from 1 to 3; shorten or split them, or raise the count in the file with a reason\n"
	if err != nil || code != 1 || out.String() != want {
		t.Errorf("lint -baseline printed\n%s\nexit %d, %v; want\n%s\nexit 1", out.String(), code, err, want)
	}
}

// lint takes -style or -baseline, or neither, and refuses the two
// together.
func TestLintRefusesBaselineWithStyle(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261009_SNX2, "flags-exclusive")
	for _, tc := range []struct {
		args            []string
		style, baseline bool
	}{
		{nil, false, false},
		{[]string{"-style"}, true, false},
		{[]string{"-baseline"}, false, true},
	} {
		style, baseline, err := lintFlags(tc.args)
		if err != nil || style != tc.style || baseline != tc.baseline {
			t.Errorf("lintFlags(%q) = style %v, baseline %v, %v; want style %v, baseline %v", tc.args, style, baseline, err, tc.style, tc.baseline)
		}
	}
	for _, args := range [][]string{{"-style", "-baseline"}, {"-baseline", "-style"}} {
		var out strings.Builder
		if code, err := runLint(&out, chartertest.Options{}, args); err == nil || !strings.Contains(err.Error(), "not both") || out.Len() != 0 {
			t.Errorf("lint %s printed %q, exit %d, %v; want it refused before it runs", strings.Join(args, " "), out.String(), code, err)
		}
	}
}

// next takes a prefix and, after it, how many ids to give: one when
// nothing follows, and a count that is not a number above zero is refused.
func TestNextReadsItsPrefixAndCount(t *testing.T) {
	chartertest.Verifies(t, claims.CH_040, "next-count")
	for _, tc := range []struct {
		args   []string
		prefix string
		count  int
	}{
		{[]string{"SW"}, "SW", 1},
		{[]string{"SW", "1"}, "SW", 1},
		{[]string{"STK", "12"}, "STK", 12},
	} {
		prefix, count, err := nextArgs(tc.args)
		if err != nil || prefix != tc.prefix || count != tc.count {
			t.Errorf("nextArgs(%q) = %q, %d, %v; want %q, %d", tc.args, prefix, count, err, tc.prefix, tc.count)
		}
	}
	for _, args := range [][]string{{"SW", "0"}, {"SW", "-2"}, {"SW", "three"}, {"SW", "1.5"}, {}, {"SW", "2", "3"}} {
		if prefix, count, err := nextArgs(args); err == nil {
			t.Errorf("nextArgs(%q) = %q, %d; want it refused", args, prefix, count)
		}
	}
}

// questions prints a block per open question, the id and the entry it is
// about, then its text, the provisional answer and who asked, a blank line
// after each, and last how many are open; -all prints the answered ones in
// their place too, each with its answer.
func TestQuestionsListsTheOpenOnesThenHowMany(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261009_4B3B, "questions-verb-lists")
	qs := []chartertest.Question{
		{ID: "Q-261009-M4TK", Text: "Does a denied request re-read the file?", AskedBy: "agent, PR #452", Entry: "SW-010/burst-caps-tokens", Provisional: "yes, as the text says"},
		{ID: "Q-261009-N7PR", Text: "Is the hash the whole file?", AskedBy: "agent, PR #450", Entry: "SW-010", Provisional: "yes", Answer: "2026-10-09, owner: yes, the whole file."},
		{ID: "Q-261009-P3RX", Text: "Does a request after a denied one spend a token?", AskedBy: "agent, PR #453", Entry: "SW-010/waited-request-passes", Provisional: "yes"},
	}
	first := "Q-261009-M4TK  SW-010/burst-caps-tokens\n" +
		"  Does a denied request re-read the file?\n" +
		"  provisional: yes, as the text says\n" +
		"  asked by: agent, PR #452\n" +
		"\n"
	answered := "Q-261009-N7PR  SW-010\n" +
		"  Is the hash the whole file?\n" +
		"  provisional: yes\n" +
		"  asked by: agent, PR #450\n" +
		"  answer: 2026-10-09, owner: yes, the whole file.\n" +
		"\n"
	third := "Q-261009-P3RX  SW-010/waited-request-passes\n" +
		"  Does a request after a denied one spend a token?\n" +
		"  provisional: yes\n" +
		"  asked by: agent, PR #453\n" +
		"\n"
	for _, tc := range []struct {
		name string
		qs   []chartertest.Question
		all  bool
		want string
	}{
		{"open", qs, false, first + third + "2 open\n"},
		{"all", qs, true, first + answered + third + "2 open\n"},
		{"none", nil, false, "0 open\n"},
	} {
		var b strings.Builder
		writeQuestions(&b, tc.qs, tc.all)
		if b.String() != tc.want {
			t.Errorf("%s: printed\n%s\nwant\n%s", tc.name, b.String(), tc.want)
		}
	}
}
