package chartertest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// sharedOnceGo cites two decisions from one test, and no need.
const sharedOnceGo = "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n" +
	"func TestShared(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_030, \"a\")\n\tchartertest.Verifies(t, claims.SW_031, \"b\")\n}\n"

// atTheBaseline writes into the fixture a record whose style counts are
// the fixture baseline's: eleven retired words, a title over twenty words,
// a scenario sentence over thirty words and a test two decisions share,
// and no text over six sentences or Why that only points.
func atTheBaseline(t *testing.T, root string) {
	t.Helper()
	writeFile(t, root, "claims/style.md", styleWords)
	writeFile(t, root, "claims/titles.md", "# Titles\n\n"+styleEntry("SW-023", wordsLong(21), "It rules.", "someone needed it."))
	writeFile(t, root, "claims/scenarios.md", "# Scenarios\n\n"+
		styleEntry("SW-027", "A decision with a long scenario", "It rules.", "someone needed it.", "`long-one`: "+wordsLong(31)+"."))
	writeFile(t, root, "claims/shared.md", "# Shared\n\n"+
		styleEntry("SW-030", "One decision", "It rules.", "someone needed it.", "`a`: one thing holds.")+
		styleEntry("SW-031", "Another decision", "It rules.", "someone needed it.", "`b`: another holds."))
	writeFile(t, root, "pkg/a/shared_test.go", sharedOnceGo)
}

// With every style count at or under the baseline's, lint -baseline prints
// nothing, though the gate fails on the tree and a question is open.
func TestLintBaselineIsSilentWithinTheBaseline(t *testing.T) {
	Verifies(t, claims.CH_261009_SNX2, "within-baseline-silent")
	root := copyFixture(t)
	atTheBaseline(t, root)
	writeFile(t, root, "claims/questions.md", "# Questions\n\n"+
		"## Q-261009-M4TK Does a quota keep its defaults?\nAsked-by: agent, PR #1\nEntry: SW-020\nProvisional: yes\n")

	// Lint names the gate's problems, the open question and the four counts
	// the baseline holds.
	wantCounts(t, root,
		"claims/retired.txt:1: lint: 11 retired words are used; "+listsThem+" them",
		"claims/README.md:1: lint: 1 title runs past twenty words; "+listsThem+" it",
		"claims/README.md:1: lint: 1 scenario sentence runs past thirty words; "+listsThem+" it",
		"claims/README.md:1: lint: 1 test is shared by two or more decisions; "+listsThem+" it",
	)
	lines, err := Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if !anyContains(lines, "rule 22: the generated Go differs") || !anyContains(lines, "claims/questions.md:1: lint: 1 question is open") {
		t.Fatalf("the tree should fail the gate and hold an open question; lint said\n%s", strings.Join(lines, "\n"))
	}

	if got, err := LintBaseline(fixtureOptions(root)); err != nil || len(got) != 0 {
		t.Errorf("at the baseline: lint -baseline said %q, %v; want nothing", got, err)
	}
	// Under it is silent too: a count is held under its baseline, not to it.
	if err := os.Remove(filepath.Join(root, "claims", "scenarios.md")); err != nil {
		t.Fatal(err)
	}
	if got, err := LintBaseline(fixtureOptions(root)); err != nil || len(got) != 0 {
		t.Errorf("under the baseline: lint -baseline said %q, %v; want nothing", got, err)
	}
}

// A check the baseline leaves out is held to zero, so its first finding is
// named as grown from 0 to 1, at the file's first line; a missing baseline
// is an error, never a baseline of zeros.
func TestLintBaselineReadsAMissingCheckAsZero(t *testing.T) {
	Verifies(t, claims.CH_261009_SNX2, "missing-check-read-as-zero")
	root := copyFixture(t)
	if got, err := LintBaseline(fixtureOptions(root)); err != nil || len(got) != 0 {
		t.Fatalf("the clean fixture: lint -baseline said %q, %v; want nothing", got, err)
	}
	writeFile(t, root, "claims/whys.md", "# Whys\n\n"+styleEntry("SW-050", "A decision with a story for a Why", "It rules.", "#12 R3."))
	want := []string{`claims/style-baseline.txt:1: lint: Whys that only point at a ruling grew from 0 to 1; give the reason in each, or add "pointing-whys 1" to the file with a reason`}
	if got, err := LintBaseline(fixtureOptions(root)); err != nil || !slices.Equal(got, want) {
		t.Errorf("a Why that only points, which the baseline leaves out: lint -baseline said\n%s\n%v\nwant\n%s", strings.Join(got, "\n"), err, strings.Join(want, "\n"))
	}

	if err := os.Remove(filepath.Join(root, "claims", "style-baseline.txt")); err != nil {
		t.Fatal(err)
	}
	if got, err := LintBaseline(fixtureOptions(root)); err == nil || !strings.Contains(err.Error(), "claims/style-baseline.txt: there is no style baseline") {
		t.Errorf("without the file: lint -baseline said %q, %v; want an error saying there is no baseline", got, err)
	}
}

// A line of the baseline that names no style check, names one counted
// already, or holds anything but a whole count after the name is a lint
// line at that line, and the only one there: the Why that only points is
// not named as grown past a count its line does not give.
func TestLintBaselineNamesABadLine(t *testing.T) {
	Verifies(t, claims.CH_261009_SNX2, "bad-line-named")
	root := copyFixture(t)
	writeFile(t, root, "claims/whys.md", "# Whys\n\n"+styleEntry("SW-050", "A decision with a story for a Why", "It rules.", "#12 R3."))
	writeFile(t, root, "claims/style-baseline.txt", "# A bad line of each kind, after a good one.\n"+
		"retired-words 0\n"+
		"long-titles 3\n"+
		"pointing-whys many\n"+
		"texts-over-six-sentences\n"+
		"retired-words 4\n"+
		"shared-tests -1\n"+
		"titles-over-twenty-words 1 2\n")
	form := func(line int, written string) string {
		return fmt.Sprintf(`claims/style-baseline.txt:%d: lint: a line is "check count", the count a whole number, not %q`, line, written)
	}
	want := []string{
		`claims/style-baseline.txt:3: lint: "long-titles" is not a style check; the checks are retired-words, titles-over-twenty-words, texts-over-six-sentences, scenario-sentences-over-thirty-words, pointing-whys and shared-tests`,
		form(4, "pointing-whys many"),
		form(5, "texts-over-six-sentences"),
		"claims/style-baseline.txt:6: lint: retired-words is counted at line 2 already; a check has one count",
		form(7, "shared-tests -1"),
		form(8, "titles-over-twenty-words 1 2"),
	}
	if got, err := LintBaseline(fixtureOptions(root)); err != nil || !slices.Equal(got, want) {
		t.Errorf("lint -baseline said\n%s\n%v\nwant\n%s", strings.Join(got, "\n"), err, strings.Join(want, "\n"))
	}
}
