package chartertest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// The style baseline holds the style checks' counts to a bar that only
// tightens. style-baseline.txt in the entries directory gives, a line per
// check, the check's name and the count it is held to:
//
//	scenario-sentences-over-thirty-words 83
//
// A blank line or one opening with # is no count, and a check the file
// leaves out is held to zero. claimsgen lint -baseline runs the style
// checks alone and names each whose count grew past its baseline; a count
// is lowered in the file as what it counts is fixed.

// styleBaselineFile is the style baseline, beside the area documents.
const styleBaselineFile = "style-baseline.txt"

// styleBaselinePath is the style baseline, slash-relative to the root.
func (r *resolved) styleBaselinePath() string {
	return path.Join(r.entriesDir, styleBaselineFile)
}

// wholeCount is a count of the style baseline: a whole number in digits.
var wholeCount = regexp.MustCompile(`^[0-9]+$`)

// baselineCount is the count the style baseline holds a check to and the
// line it is written on; bad when no count can be read from that line, so
// the check is held to nothing.
type baselineCount struct {
	count, line int
	bad         bool
}

// lintLine is a lint line and the line of the file it names.
type lintLine struct {
	n    int
	text string
}

// LintBaseline loads the charter, runs the style checks alone and holds
// each count to the style baseline, as claimsgen lint -baseline does: it
// returns a line for each check whose count grew past its baseline and one
// for each line of the baseline that is not a style check and a whole
// count, in the order of the file's lines, and nothing when no count grew.
// The gate's problems are not part of it. A missing baseline is an error.
//
//charter:realizes CH-261009-SNX2
func LintBaseline(opts Options) ([]string, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	return lintBaseline(r)
}

// lintBaseline is LintBaseline over a resolved charter.
func lintBaseline(r *resolved) ([]string, error) {
	rel := r.styleBaselinePath()
	data, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%s: there is no style baseline; write one, a line \"check count\" per style check, with the counts lint gives", rel)
	}
	if err != nil {
		return nil, err
	}
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	style, err := lintStyle(s)
	if err != nil {
		return nil, err
	}
	return style.overBaseline(rel, string(data)), nil
}

// overBaseline holds each check's count to the style baseline rel, whose
// text is src. It names each check whose count exceeds its baseline, at the
// check's line or, when the file leaves the check out, at the first line,
// with both counts and what to do; and each line of the file that is not a
// style check, counted once, and a whole count. The lines come in the
// order of the file's, and none comes when no count grew.
//
//charter:realizes CH-261009-SNX2
func (f *styleFindings) overBaseline(rel, src string) []string {
	base, lines := parseStyleBaseline(rel, src, f.checks)
	for _, c := range f.checks {
		b, n := base[c.name], len(c.found)
		switch {
		case b.bad || n <= b.count:
		case b.line == 0:
			lines = append(lines, lintLine{1, fmt.Sprintf("%s:1: lint: %s grew from 0 to %d; %s, or add %q to the file with a reason", rel, c.grew, n, c.fix, fmt.Sprintf("%s %d", c.name, n))})
		default:
			lines = append(lines, lintLine{b.line, fmt.Sprintf("%s:%d: lint: %s grew from %d to %d; %s, or raise the count in the file with a reason", rel, b.line, c.grew, b.count, n, c.fix)})
		}
	}
	sort.SliceStable(lines, func(i, j int) bool { return lines[i].n < lines[j].n })
	var out []string
	for _, l := range lines {
		out = append(out, l.text)
	}
	return out
}

// parseStyleBaseline reads the style baseline rel, whose text is src: the
// count each check is held to, with its line, and a lint line for each
// line that names no check of checks, names one a second time, or holds
// anything but a whole count after the name. A blank line, or one opening
// with #, is no count.
//
//charter:realizes CH-261009-SNX2
func parseStyleBaseline(rel, src string, checks []*styleCheck) (map[string]baselineCount, []lintLine) {
	known := map[string]bool{}
	var names []string
	for _, c := range checks {
		known[c.name] = true
		names = append(names, c.name)
	}
	base := map[string]baselineCount{}
	var lines []lintLine
	bad := func(n int, format string, args ...any) {
		lines = append(lines, lintLine{n, fmt.Sprintf("%s:%d: lint: %s", rel, n, fmt.Sprintf(format, args...))})
	}
	for i, l := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		n := i + 1
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		fields := strings.Fields(l)
		name := fields[0]
		if !known[name] {
			bad(n, "%q is not a style check; the checks are %s", name, andList(names))
			continue
		}
		if at := base[name].line; at != 0 {
			bad(n, "%s is counted at line %d already; a check has one count", name, at)
			continue
		}
		count := -1
		if len(fields) == 2 && wholeCount.MatchString(fields[1]) {
			if c, err := strconv.Atoi(fields[1]); err == nil {
				count = c
			}
		}
		if count < 0 {
			bad(n, "a line is \"check count\", the count a whole number, not %q", l)
			base[name] = baselineCount{line: n, bad: true}
			continue
		}
		base[name] = baselineCount{count: count, line: n}
	}
	return base, lines
}
