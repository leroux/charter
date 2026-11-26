package lint

import (
	"errors"
	"fmt"
	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/record"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// The style a reader of the record stumbles on is the last thing lint
// names, by the bars of the README's "Style", so that an audit spends a
// reader only on what needs one: a word the record retired, a title over
// twenty words, a text over six sentences, a scenario sentence over thirty
// words, a Why that points at a ruling instead of giving the reason, and a
// test two decisions share and no need. What is read is the live entries,
// active and planned, of the area documents and the needs, and each
// document's header; a retired entry, a handle, Guidance and Changes are not. Lint counts what each check found,
// in one line naming the command that lists it; charter lint -style
// lists every finding at its file and line.
//
// The retired words are retired.txt in the entries directory, one a line:
//
//	quotas? -> bucket | quota defaults; new quota
//
// a pattern, RE2 matched case-insensitively at word boundaries; after ->
// the word to write instead; after | the phrases, separated by ;, whose
// line or sentence keeps the word. A blank line or one opening with # is
// no word. A word inside a code span, a double-quoted string or a fenced
// block is not read.

// retiredFile is the words the record retired, beside the area documents.
// The charter that keeps none names none.
const retiredFile = "retired.txt"

// retiredPath is the retired words' file, slash-relative to the root.
func retiredPath(dir string) string { return path.Join(dir, retiredFile) }

// The bars of the README's "Style".
const (
	titleWords    = 20
	textSentences = 6
	scenarioWords = 30
	// pointerWords is the length under which a Why that names a ruling is
	// taken to give no reason of its own.
	pointerWords = 25
)

var (
	// pointer is what a Why that points names: a story, a ruling's
	// number, an RFD, a review, a ruling. A number sign is no word
	// character, so #444 is matched without a boundary before it.
	pointer = regexp.MustCompile(`(?i)(?:\b(?:backlog \d+|R\d{1,2}|RFD|review|ruling)\b|#\d+\b)`)
	// codeSpan and quoted are what the retired words are not read in.
	codeSpan = regexp.MustCompile("`[^`]*`")
	quoted   = regexp.MustCompile(`"[^"]*"|“[^”]*”`)
	// listMarker opens a numbered item, whose period ends no sentence.
	listMarker = regexp.MustCompile(`^\d+[.)]\s+`)
)

// retiredWord is one line of the retired words' file.
type retiredWord struct {
	re    *regexp.Regexp
	write string
	// keep are the phrases whose line or sentence keeps the word,
	// lower-cased.
	keep []string
}

// readRetired reads the retired words' file; a missing file is no words. A
// line that is not "pattern -> word | phrase; phrase", or whose pattern
// does not compile, is a lint line of its own.
//
//charter:realizes CH-261126-WTMZ
func readRetired(root, dir string) ([]retiredWord, []string, error) {
	rel := retiredPath(dir)
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var words []retiredWord
	var problems []string
	for i, l := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		pattern, rest, ok := strings.Cut(l, "->")
		write, keep, _ := strings.Cut(rest, "|")
		pattern, write = strings.TrimSpace(pattern), strings.TrimSpace(write)
		if !ok || pattern == "" || write == "" {
			problems = append(problems, fmt.Sprintf("%s:%d: lint: a retired word is \"pattern -> the word to write\", then \"| phrase; phrase\" for the lines that keep it", rel, i+1))
			continue
		}
		re, err := regexp.Compile(`(?i)\b(?:` + pattern + `)\b`)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s:%d: lint: the pattern %q does not compile: %v", rel, i+1, pattern, err))
			continue
		}
		w := retiredWord{re: re, write: write}
		for _, p := range strings.Split(keep, ";") {
			if p = strings.TrimSpace(p); p != "" {
				w.keep = append(w.keep, lowerASCII(p))
			}
		}
		words = append(words, w)
	}
	return words, problems, nil
}

// retiredUse is one retired word in the record: its line, the word as
// written, and the word to write.
type retiredUse struct {
	line        int
	word, write string
}

// retiredIn finds the retired words in prose, given as its lines with their
// numbers: each match outside code spans, double-quoted strings and fenced
// blocks, a match spanning lines too, unless the line or the sentence that
// holds it holds one of the word's phrases, a phrase wrapped onto the next
// line counting for both. A word is found once a line.
//
//charter:realizes CH-261126-WTMZ
func retiredIn(words []retiredWord, lines []record.Line) []retiredUse {
	if len(words) == 0 || len(lines) == 0 {
		return nil
	}
	text := proseLines(lines)
	flat := strings.Join(text, " ")
	starts := make([]int, len(text))
	for i, off := 0, 0; i < len(text); i++ {
		starts[i] = off
		off += len(text[i]) + 1
	}
	lineAt := func(at int) int {
		return sort.Search(len(starts), func(i int) bool { return starts[i] > at }) - 1
	}
	ends := sentenceEnds(codeAsWord(flat))
	sentenceAt := func(at int) (int, int) {
		k := sort.SearchInts(ends, at+1)
		from, to := 0, len(flat)
		if k > 0 {
			from = ends[k-1]
		}
		if k < len(ends) {
			to = ends[k]
		}
		return from, to
	}
	read := blank(quoted, blank(codeSpan, flat))
	lower := lowerASCII(flat)
	var uses []retiredUse
	seen := map[string]bool{}
	for _, w := range words {
		for _, m := range w.re.FindAllStringIndex(read, -1) {
			first, last := lineAt(m[0]), lineAt(m[1]-1)
			from, to := sentenceAt(m[0])
			if w.kept(lower, starts[first], starts[last]+len(text[last])) || w.kept(lower, from, to) {
				continue
			}
			word := flat[m[0]:m[1]]
			k := fmt.Sprintf("%d\x00%s", lines[first].N, lowerASCII(word))
			if seen[k] {
				continue
			}
			seen[k] = true
			uses = append(uses, retiredUse{line: lines[first].N, word: word, write: w.write})
		}
	}
	return uses
}

// kept reports whether one of the word's phrases lies, whole or in part,
// within lower[from:to].
func (w retiredWord) kept(lower string, from, to int) bool {
	for _, p := range w.keep {
		for i := 0; i < len(lower); {
			j := strings.Index(lower[i:], p)
			if j < 0 {
				break
			}
			if at := i + j; at < to && at+len(p) > from {
				return true
			}
			i += j + 1
		}
	}
	return false
}

// proseLines is the lines' text, trimmed, with a fenced block's lines,
// its fences among them, made empty, and a numbered item's marker blank.
func proseLines(lines []record.Line) []string {
	out := make([]string, len(lines))
	fence := ""
	for i, l := range lines {
		switch {
		case fence != "":
			if record.FenceCloses(l.Text, fence) {
				fence = ""
			}
		case record.FenceOpen(l.Text) != "":
			fence = record.FenceOpen(l.Text)
		default:
			t := strings.TrimSpace(l.Text)
			if m := listMarker.FindString(t); m != "" {
				t = strings.Repeat(" ", len(m)) + t[len(m):]
			}
			out[i] = t
		}
	}
	return out
}

// codeAsWord makes each code span of s one capital word as long, so a
// sentence may open with one and no period in one ends a sentence.
func codeAsWord(s string) string {
	return codeSpan.ReplaceAllStringFunc(s, func(m string) string { return "C" + strings.Repeat(" ", len(m)-1) })
}

// blank replaces each match of re in s with as many spaces, so the offsets
// of the rest hold.
func blank(re *regexp.Regexp, s string) string {
	return re.ReplaceAllStringFunc(s, func(m string) string { return strings.Repeat(" ", len(m)) })
}

// lowerASCII lower-cases the ASCII letters of s alone, so the offsets into
// s hold in the result.
func lowerASCII(s string) string {
	b := []byte(s)
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return string(b)
}

// countWords counts the words of a sentence: its whitespace-separated
// tokens that hold a letter or a digit.
//
//charter:realizes CH-261126-9K5J
func countWords(s string) int {
	n := 0
	for _, f := range strings.Fields(s) {
		if strings.IndexFunc(f, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }) >= 0 {
			n++
		}
	}
	return n
}

// countSentences counts the sentences of a text, as sentenceEnds finds
// them: a period in a code span ends none, nor does a numbered item's
// marker or anything in a fenced block.
//
//charter:realizes CH-261126-9K5J
func countSentences(text string) int {
	var lines []record.Line
	for _, l := range strings.Split(text, "\n") {
		lines = append(lines, record.Line{Text: l})
	}
	return len(sentenceEnds(codeAsWord(strings.Join(proseLines(lines), " "))))
}

// sentenceEnds returns where each sentence of s ends: just past a full
// stop, a question mark or an exclamation mark, and the closing marks after
// it, that is followed by a space and a capital letter, or by the end; a
// last run that ends in none ends at the end. The period of i.e. or e.g.
// ends none.
//
//charter:realizes CH-261126-9K5J
func sentenceEnds(s string) []int {
	var ends []int
	last := 0
	for i := 0; i < len(s); i++ {
		if !strings.ContainsRune(".!?", rune(s[i])) || s[i] == '.' && abbreviated(s[:i+1]) {
			continue
		}
		j := i + 1
		for j < len(s) {
			r, size := utf8.DecodeRuneInString(s[j:])
			if !strings.ContainsRune(`"')]*_”’`, r) {
				break
			}
			j += size
		}
		if j == len(s) || s[j] == ' ' && opensWithCapital(strings.TrimLeft(s[j:], " ")) {
			ends = append(ends, j)
			last = j
		}
	}
	if strings.IndexFunc(s[last:], func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }) >= 0 {
		ends = append(ends, len(s))
	}
	return ends
}

// abbreviated reports text ending in i.e. or e.g., whose last period ends
// no sentence.
func abbreviated(text string) bool {
	l := lowerASCII(text)
	for _, a := range []string{"i.e.", "e.g."} {
		if strings.HasSuffix(l, a) {
			before, _ := utf8.DecodeLastRuneInString(l[:len(l)-len(a)])
			return before == utf8.RuneError || !unicode.IsLetter(before)
		}
	}
	return false
}

// opensWithCapital reports whether s, past any opening quote, parenthesis
// or emphasis, opens with a capital letter.
func opensWithCapital(s string) bool {
	for _, r := range s {
		if strings.ContainsRune(`"'([*_“‘`, r) {
			continue
		}
		return unicode.IsUpper(r)
	}
	return false
}

// styleCheck is what one style check found, each finding a lint line at
// its file and line, and how its count is said: at the file the count
// names, "1 title runs" when it found one and "4 titles run" when more. The
// style baseline counts it under its name, and lint -baseline says what
// grew and the fix when its count passes the baseline's.
type styleCheck struct {
	name, grew, fix string
	at, one, many   string
	found           []string
}

// styleFindings are what the style checks found, check by check in the
// order their counts are printed, and the retired words' file's own
// problems, which lint prints however the findings are.
type styleFindings struct {
	checks   []*styleCheck
	problems []string
}

// lines lists every finding, check by check.
func (f *styleFindings) lines() []string {
	var out []string
	for _, c := range f.checks {
		out = append(out, c.found...)
	}
	return out
}

// counts says, in one line for each check that found anything, how many
// findings it has, and names the command that lists them.
//
//charter:realizes CH-261126-9K5J
func (f *styleFindings) counts() []string {
	var out []string
	for _, c := range f.checks {
		switch n := len(c.found); n {
		case 0:
		case 1:
			out = append(out, fmt.Sprintf("%s:1: lint: 1 %s; charter lint -style lists it", c.at, c.one))
		default:
			out = append(out, fmt.Sprintf("%s:1: lint: %d %s; charter lint -style lists them", c.at, n, c.many))
		}
	}
	return out
}

// lintStyle checks the style the record's reader would stumble on, by the
// bars of the README's "Style": in the live entries of the area documents
// and the needs, a retired word in a title, a text, a Why or a scenario's
// sentence, a title over twenty words, a text over six sentences, a
// scenario sentence over thirty words, and a Why under twenty-five words
// that names a ruling; a retired word in a document's header; and a test
// two decisions share and no need. Each finding names the file and the
// line: the heading's, or the scenario's, or the line a retired word is
// on. The retired words' count names their file, the others the README.
//
//charter:realizes CH-261126-9K5J, CH-261126-WTMZ
func lintStyle(t *check.Tree) (*styleFindings, error) {
	dir := t.Record.Dir
	words, problems, err := readRetired(t.Config.Root, dir)
	if err != nil {
		return nil, err
	}
	readme := path.Join(dir, record.ReadmeFile)
	used := &styleCheck{name: "retired-words", grew: "retired words", fix: "write the words " + retiredPath(dir) + " gives instead",
		at: retiredPath(dir), one: "retired word is used", many: "retired words are used"}
	titles := &styleCheck{name: "titles-over-twenty-words", grew: "titles over twenty words", fix: "shorten them",
		at: readme, one: "title runs past twenty words", many: "titles run past twenty words"}
	texts := &styleCheck{name: "texts-over-six-sentences", grew: "texts over six sentences", fix: "shorten them or split their decisions",
		at: readme, one: "text runs past six sentences", many: "texts run past six sentences"}
	scenarios := &styleCheck{name: "scenario-sentences-over-thirty-words", grew: "scenario sentences over thirty words", fix: "shorten or split them",
		at: readme, one: "scenario sentence runs past thirty words", many: "scenario sentences run past thirty words"}
	whys := &styleCheck{name: "pointing-whys", grew: "Whys that only point at a ruling", fix: "give the reason in each",
		at: readme, one: "Why only points at a ruling", many: "Whys only point at a ruling"}
	shared := &styleCheck{name: "shared-tests", grew: "tests two decisions share", fix: "give each decision a test of its own",
		at: readme, one: "test is shared by two or more decisions", many: "tests are shared by two or more decisions"}
	retired := func(file, part string, lines []record.Line) {
		for _, u := range retiredIn(words, lines) {
			used.found = append(used.found, fmt.Sprintf("%s:%d: lint: %q is retired in %s; write %s (%s)", file, u.line, u.word, part, u.write, retiredPath(dir)))
		}
	}
	for _, d := range t.Record.Areas {
		retired(d.File, "the header", d.Prose)
	}
	for _, e := range t.Record.Entries {
		if e.Status != "active" && e.Status != "planned" {
			continue
		}
		retired(e.File, "the title of "+e.ID, []record.Line{{N: e.Line, Text: e.Title}})
		retired(e.File, "the text of "+e.ID, e.TextLines)
		retired(e.File, "the Why of "+e.ID, e.WhyLines)
		for _, sc := range e.Scenarios {
			retired(e.File, fmt.Sprintf("scenario `%s` of %s", sc.Handle, e.ID), []record.Line{{N: sc.Line, Text: sc.Text}})
		}
		if n := countWords(e.Title); n > titleWords {
			titles.found = append(titles.found, fmt.Sprintf("%s: lint: the title of %s runs to %d words; twenty is the bar (README, Style)", e.Pos(), e.ID, n))
		}
		if n := countSentences(e.Text); n > textSentences {
			texts.found = append(texts.found, fmt.Sprintf("%s: lint: the text of %s runs to %d sentences; six is the bar (README, Style)", e.Pos(), e.ID, n))
		}
		for _, sc := range e.Scenarios {
			if n := countWords(sc.Text); n > scenarioWords {
				scenarios.found = append(scenarios.found, fmt.Sprintf("%s:%d: lint: scenario `%s` of %s runs to %d words; thirty is the bar (README, Style)", e.File, sc.Line, sc.Handle, e.ID, n))
			}
		}
		if countWords(e.Why) < pointerWords && pointer.MatchString(e.Why) {
			whys.found = append(whys.found, fmt.Sprintf("%s: lint: the Why of %s points at a ruling and gives no reason (README, Style)", e.Pos(), e.ID))
		}
	}
	shared.found = sharedTests(t)
	return &styleFindings{checks: []*styleCheck{used, titles, texts, scenarios, whys, shared}, problems: problems}, nil
}

// andList joins items as a sentence does: "a", "a and b", "a, b and c".
func andList(items []string) string {
	if len(items) < 2 {
		return strings.Join(items, "")
	}
	return strings.Join(items[:len(items)-1], ", ") + " and " + items[len(items)-1]
}
