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

// styleEntry is a planned decision for the style fixtures: its heading,
// text and Why, then its scenarios as written after their numbers.
func styleEntry(id, title, text, why string, scenarios ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s %s\nStatus: planned\nServes: STK-001\n\n%s\n\nWhy: %s\n", id, title, text, why)
	if len(scenarios) > 0 {
		b.WriteString("\nScenarios:\n")
		for i, s := range scenarios {
			fmt.Fprintf(&b, "%d. %s\n", i+1, s)
		}
	}
	return b.String() + "\n"
}

// wordsLong is a sentence of n words, one of them written as two around a
// dash that is no word.
func wordsLong(n int) string {
	w := make([]string, n)
	for i := range w {
		w[i] = fmt.Sprintf("w%d", i+1)
	}
	return strings.Join(w[:n-1], " ") + " — " + w[n-1]
}

// lineIn is the number of the first line of doc holding text.
func lineIn(t *testing.T, doc, text string) int {
	t.Helper()
	for i, l := range strings.Split(doc, "\n") {
		if strings.Contains(l, text) {
			return i + 1
		}
	}
	t.Fatalf("no line holds %q", text)
	return 0
}

// styleLines runs lint -style over the tree, every finding of the style
// checks listed, and returns its lines holding marker, in lint's order.
func styleLines(t *testing.T, root, marker string) []string {
	t.Helper()
	lines, err := LintStyle(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, l := range lines {
		if strings.Contains(l, marker) {
			got = append(got, l)
		}
	}
	return got
}

// listsThem is how a style check's count names the command that lists it.
const listsThem = "go run ./cmd/claimsgen lint -style lists"

// wantCounts fails the test unless lint, the style checks counted, ends
// with exactly the count lines want, in that order, and lists no finding.
func wantCounts(t *testing.T, root string, want ...string) {
	t.Helper()
	lines, err := Lint(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	var counts []string
	for _, l := range lines {
		if strings.Contains(l, listsThem) {
			counts = append(counts, l)
		}
		if strings.Contains(l, "(README, Style)") || strings.Contains(l, "is retired in") || strings.Contains(l, "two decisions share one test") {
			t.Errorf("lint listed a finding without -style: %s", l)
		}
	}
	if !slices.Equal(counts, want) || len(lines) < len(want) || !slices.Equal(lines[len(lines)-len(want):], want) {
		t.Errorf("lint said\n%s\nwant it to end with the views\n%s", strings.Join(lines, "\n"), strings.Join(want, "\n"))
	}
}

// wantLines fails the test unless got holds exactly want, in any order.
func wantLines(t *testing.T, what string, got, want []string) {
	t.Helper()
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("%s: lint said\n%s\nwant\n%s", what, strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// styleWords is an area document that uses the fixture's retired words
// everywhere lint reads and everywhere it does not.
const styleWords = "# Style\nCode: pkg/a\n\n" +
	"A quota opens the header.\n" +
	"The `quota` in code and the \"quota\" quoted are not read.\n" +
	"The quota defaults are kept. So is this quota.\n" +
	"\n" +
	"## SW-020 A quota keeps its read model\n" +
	"Status: planned\n" +
	"Serves: STK-001\n" +
	"\n" +
	"The quota defaults apply, and a quotas panel lists each\n" +
	"quota it holds. Quotas end, as other buckets do.\n" +
	"A quota, and another quota, end with it. Quotas\n" +
	"upsert their counters in the read\n" +
	"model.\n" +
	"\n" +
	"Why: a quota upserted once is read once.\n" +
	"\n" +
	"Scenarios:\n" +
	"1. `quota-kept` (owner, #12): a quota survives a restart.\n" +
	"2. `defaults-kept`: the quota\n" +
	"   defaults hold.\n" +
	"\n" +
	"Guidance: a quota is not read here.\n" +
	"\n" +
	"Changes:\n" +
	"- 2026-10-09: a quota was not read here.\n" +
	"\n" +
	"## SW-021 A retired quota stays as it was\n" +
	"Status: deprecated\n" +
	"Serves: STK-001\n" +
	"\n" +
	"A quota is kept.\n" +
	"\n" +
	"Why: a quota.\n"

// A retired word is counted, and listed under -style where a reader meets
// it, with the word to write, and nowhere a reader is meant to see it as
// written.
func TestLintNamesARetiredWord(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "claims/style.md", styleWords)
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "retired-word-named")
		wantCounts(t, root, "claims/retired.txt:1: lint: 11 retired words are used; "+listsThem+" them")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "retired-word-named")
		listsRetiredWords(t, root)
	})
}

// listsRetiredWords checks the retired words lint -style lists in the
// styleWords document, and the retired words' file's own problems.
func listsRetiredWords(t *testing.T, root string) {
	t.Helper()
	at := func(text string) int { return lineIn(t, styleWords, text) }
	retired := func(line int, word, part, write string) string {
		return fmt.Sprintf("claims/style.md:%d: lint: %q is retired in %s; write %s (claims/retired.txt)", line, word, part, write)
	}
	wantLines(t, "the record", styleLines(t, root, "(claims/retired.txt)"), []string{
		retired(at("A quota opens"), "quota", "the header", "bucket"),
		retired(at("## SW-020"), "quota", "the title of SW-020", "bucket"),
		retired(at("## SW-020"), "read model", "the title of SW-020", "the views"),
		retired(at("quota it holds"), "Quotas", "the text of SW-020", "bucket"),
		retired(at("A quota, and"), "quota", "the text of SW-020", "bucket"),
		retired(at("A quota, and"), "Quotas", "the text of SW-020", "bucket"),
		retired(at("upsert their"), "upsert", "the text of SW-020", "counter update"),
		retired(at("upsert their"), "read model", "the text of SW-020", "the views"),
		retired(at("Why: a quota"), "quota", "the Why of SW-020", "bucket"),
		retired(at("Why: a quota"), "upserted", "the Why of SW-020", "counter update"),
		retired(at("`quota-kept`"), "quota", "scenario `quota-kept` of SW-020", "bucket"),
	})

	// A line whose pattern does not compile, or that lacks the arrow, the
	// pattern or the word to write, is named itself.
	appendTo(t, root, "claims/retired.txt", "(unclosed -> anything\nno arrow here\nno word ->\n-> no pattern\n")
	got := styleLines(t, root, "claims/retired.txt:")
	form := "lint: a retired word is \"pattern -> the word to write\", then \"| phrase; phrase\" for the lines that keep it"
	if len(got) != 4 || !strings.HasPrefix(got[0], "claims/retired.txt:6: lint: the pattern \"(unclosed\" does not compile") ||
		got[1] != "claims/retired.txt:7: "+form || got[2] != "claims/retired.txt:8: "+form || got[3] != "claims/retired.txt:9: "+form {
		t.Errorf("a bad pattern and three malformed lines: lint said %q", got)
	}

	// A charter that keeps no retired words is told of none, and a file
	// that cannot be read is an error, never no words.
	if err := os.Remove(filepath.Join(root, "claims", "retired.txt")); err != nil {
		t.Fatal(err)
	}
	if got := styleLines(t, root, "retired"); len(got) != 0 {
		t.Errorf("without the file: lint said %q", got)
	}
	if err := os.Mkdir(filepath.Join(root, "claims", "retired.txt"), 0o750); err != nil {
		t.Fatal(err)
	}
	if lines, err := LintStyle(fixtureOptions(root)); err == nil {
		t.Errorf("a retired-words file that cannot be read gave %d lines and no error", len(lines))
	}
}

// A title over twenty words is counted, and listed under -style with its
// length; one of twenty, or a retired entry's, is not.
func TestLintNamesATitleOverTwentyWords(t *testing.T) {
	root := copyFixture(t)
	doc := "# Titles\n\n" +
		strings.Replace(styleEntry("SW-022", wordsLong(25), "It rules.", "someone needed it."), "Status: planned", "Status: deprecated", 1) +
		styleEntry("SW-023", wordsLong(21), "It rules.", "someone needed it.") +
		styleEntry("SW-024", wordsLong(20), "It rules.", "someone needed it.")
	writeFile(t, root, "claims/titles.md", doc)
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-title-named")
		wantCounts(t, root, "claims/README.md:1: lint: 1 title runs past twenty words; "+listsThem+" it")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-title-named")
		wantLines(t, "titles", styleLines(t, root, "twenty is the bar"), []string{
			fmt.Sprintf("claims/titles.md:%d: lint: the title of SW-023 runs to 21 words; twenty is the bar (README, Style)", lineIn(t, doc, "## SW-023")),
		})
	})
}

// A text over six sentences is counted, and listed under -style with its
// length; what ends a sentence is a stop before a capital or the end.
func TestLintNamesATextOverSixSentences(t *testing.T) {
	root := copyFixture(t)
	six := "It reads its file, e.g. Flow A, the file `a.Md`, its notes, etc. and\n" +
		"i.e. Its own. Does it stop? It never stops! It writes \"Done.\" Then\n" +
		"it waits (for the person.) It ends there"
	doc := "# Texts\n\n" +
		styleEntry("SW-025", "A text of seven sentences, the last with no stop", six+". And it is gone", "someone needed it.") +
		styleEntry("SW-026", "A text of six sentences, the last with no stop", six, "someone needed it.")
	writeFile(t, root, "claims/texts.md", doc)
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-text-named")
		wantCounts(t, root, "claims/README.md:1: lint: 1 text runs past six sentences; "+listsThem+" it")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-text-named")
		wantLines(t, "texts", styleLines(t, root, "six is the bar"), []string{
			fmt.Sprintf("claims/texts.md:%d: lint: the text of SW-025 runs to 7 sentences; six is the bar (README, Style)", lineIn(t, doc, "## SW-025")),
		})
	})
}

// A scenario sentence over thirty words is counted, and listed under
// -style at its line; the handle and who ruled it are not words of the
// sentence.
func TestLintNamesAScenarioOverThirtyWords(t *testing.T) {
	root := copyFixture(t)
	long := strings.Fields(wordsLong(31))
	doc := "# Scenarios\n\n" + styleEntry("SW-027", "A decision with two scenarios", "It rules.", "someone needed it.",
		"`long-one` (owner, backlog 077): "+strings.Join(long[:12], " ")+"\n   "+strings.Join(long[12:], " ")+".",
		"`thirty` (delegate alice, #12 R3): "+wordsLong(30)+".",
	)
	writeFile(t, root, "claims/scenarios.md", doc)
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-scenario-named")
		wantCounts(t, root, "claims/README.md:1: lint: 1 scenario sentence runs past thirty words; "+listsThem+" it")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "long-scenario-named")
		wantLines(t, "scenarios", styleLines(t, root, "thirty is the bar"), []string{
			fmt.Sprintf("claims/scenarios.md:%d: lint: scenario `long-one` of SW-027 runs to 31 words; thirty is the bar (README, Style)", lineIn(t, doc, "`long-one`")),
		})
	})
}

// A short Why that names a story, a ruling or a review is counted, and
// listed under -style; a long one, or one that names none, is the reason
// itself.
func TestLintNamesAWhyThatOnlyPoints(t *testing.T) {
	root := copyFixture(t)
	doc := "# Whys\n\n" +
		styleEntry("SW-050", "A decision with a story for a Why", "It rules.", "#12 R3.") +
		styleEntry("SW-051", "A decision with an issue for a Why", "It rules.", "as #444 asked.") +
		styleEntry("SW-052", "A decision with a ruling for a Why", "It rules.", "the owner's ruling.") +
		styleEntry("SW-053", "A decision with a reason for a Why", "It rules.", "the record cannot hold two answers to one question.") +
		styleEntry("SW-054", "A decision with a Why of twenty-five words", "It rules.", "a reviewer reads the decision first, and a review that finds no reason\nin it has nothing to weigh the change against when it lands.")
	writeFile(t, root, "claims/whys.md", doc)
	why := func(id string) string {
		return fmt.Sprintf("claims/whys.md:%d: lint: the Why of %s points at a ruling and gives no reason (README, Style)", lineIn(t, doc, "## "+id), id)
	}
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "pointing-why-named")
		wantCounts(t, root, "claims/README.md:1: lint: 3 Whys only point at a ruling; "+listsThem+" them")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "pointing-why-named")
		wantLines(t, "Whys", styleLines(t, root, "gives no reason"), []string{why("SW-050"), why("SW-051"), why("SW-052")})
	})
}

// sharedTestGo cites, test by test, a decision and no scenario; two
// decisions, one scenario twice; one decision; two decisions and a need;
// and a retired decision and a live one.
const sharedTestGo = "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fixture/claims\"\n\t\"github.com/leroux/charter/chartertest\"\n)\n\n" +
	"func TestNamesNoScenario(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_031)\n}\n\n" +
	"func TestShared(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_030, \"a\")\n\tchartertest.Verifies(t, claims.SW_030, \"a\")\n\tchartertest.Verifies(t, claims.SW_031, \"b\")\n}\n\n" +
	"func TestOneDecision(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_030, \"a\", \"c\")\n}\n\n" +
	"func TestTwoDecisionsAndANeed(t *testing.T) {\n\tchartertest.Verifies(t, claims.STK_030, \"x\")\n\tchartertest.Verifies(t, claims.SW_030, \"c\")\n\tchartertest.Verifies(t, claims.SW_031, \"b\")\n}\n\n" +
	"func TestARetiredDecisionAndALiveOne(t *testing.T) {\n\tchartertest.Verifies(t, claims.SW_032, \"d\")\n\tchartertest.Verifies(t, claims.SW_031, \"b\")\n}\n"

// sharedTestTS cites two decisions from one TypeScript test.
const sharedTestTS = "import { it } from \"@example/test\";\nimport { claims } from \"@fixture/claims\";\nimport { verifies } from \"./helper\";\n\n" +
	"it(\"shares\", () => {\n  verifies(claims.SW_030, \"c\");\n  verifies(claims.SW_031, \"b\");\n});\n" +
	"it(\"keeps to one\", () => {\n  verifies(claims.SW_031, \"b\");\n});\n"

// A test whose citations name scenarios of two decisions and no need is
// counted, and listed once under -style, at the test, with what it cites; a
// test that also cites a need is neither.
func TestLintNamesATestTwoDecisionsShare(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "claims/shared.md", "# Shared\n\n"+
		styleEntry("SW-030", "One decision", "It rules.", "someone needed it.", "`a`: one thing holds.", "`c`: another holds.")+
		styleEntry("SW-031", "Another decision", "It rules.", "someone needed it.", "`b`: a third holds.")+
		strings.Replace(styleEntry("SW-032", "A retired decision", "It ruled.", "someone needed it.", "`d`: a fourth held."), "Status: planned", "Status: deprecated", 1))
	writeFile(t, root, "claims/needs.md", "# Needs\n\n## STK-030 A person relies on it\nStatus: planned\n\nIt holds.\n\nWhy: someone needs it.\n\nScenarios:\n1. `x`: it holds.\n")
	writeFile(t, root, "pkg/a/shared_test.go", sharedTestGo)
	writeFile(t, root, "web/shared.test.ts", sharedTestTS)
	t.Run("counted", func(t *testing.T) {
		Verifies(t, claims.CH_041, "shared-test-named")
		wantCounts(t, root, "claims/README.md:1: lint: 2 tests are shared by two or more decisions; "+listsThem+" them")
	})
	t.Run("listed", func(t *testing.T) {
		Verifies(t, claims.CH_041, "shared-test-named")
		share := "; two decisions share one test, so one may rule what the other states"
		wantLines(t, "tests", styleLines(t, root, "two decisions share one test"), []string{
			fmt.Sprintf("pkg/a/shared_test.go:%d: lint: TestShared cites SW-030/a and SW-031/b%s", lineIn(t, sharedTestGo, "func TestShared"), share),
			fmt.Sprintf("web/shared.test.ts:%d: lint: test \"shares\" cites SW-030/c and SW-031/b%s", lineIn(t, sharedTestTS, "claims.SW_030"), share),
		})
	})
}
