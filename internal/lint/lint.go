// Package lint names, after everything the gate fails on, what the gate
// does not fail: a decision that serves nothing, a planned entry each of
// whose scenarios is verified, a deprecated entry still anchored or cited, a
// decision realized only by a package clause, a citation that names no
// scenario, the open questions in one line, and the style a reader of the
// record would stumble on.
package lint

import (
	"fmt"
	"path"
	"slices"
	"sort"

	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
)

// Lint runs the gate and then the advice. With listStyle, every style
// finding is listed at its line; without it, each style check that found
// anything prints one line counting them, after the rest.
//
//charter:realizes CH-041
func Lint(c *config.Config, listStyle bool) ([]string, error) {
	t, problems, err := check.Run(c)
	if err != nil {
		return nil, err
	}
	out := append([]string{}, problems...)
	out = append(out, questions(t.Record)...)
	style, err := lintStyle(t)
	if err != nil {
		return nil, err
	}
	out = append(out, style.problems...)
	if listStyle {
		out = append(out, style.lines()...)
	}
	out = append(out, advise(t)...)
	sort.Strings(out)
	if !listStyle {
		out = append(out, style.counts()...)
	}
	return out, nil
}

// byEntry groups the anchors and the counted citations by the entry they
// name.
func byEntry(set *facts.Set, rec *record.Record) (map[string][]facts.Anchor, map[string][]facts.Citation) {
	anchors, cites := map[string][]facts.Anchor{}, map[string][]facts.Citation{}
	for _, a := range set.Anchors {
		for _, id := range a.IDs {
			anchors[id] = append(anchors[id], a)
		}
	}
	for _, c := range set.Citations {
		if c.Observes {
			continue
		}
		if e, _ := rec.BySymbol(c.Symbol); e != nil {
			cites[e.ID] = append(cites[e.ID], c)
		}
	}
	return anchors, cites
}

func advise(t *check.Tree) []string {
	var out []string
	anchors, cites := byEntry(t.Facts, t.Record)
	for _, e := range t.Record.Entries {
		as, cs := anchors[e.ID], cites[e.ID]
		switch e.Status {
		case "deprecated":
			for _, a := range as {
				out = append(out, fmt.Sprintf("%s: lint: %s is deprecated and still anchored here; strip it", a.Pos(), e.ID))
			}
			for _, c := range cs {
				out = append(out, fmt.Sprintf("%s: lint: %s is deprecated and still verified here; strip it", c.Pos(), e.ID))
			}
			continue
		case "planned":
			if len(e.Scenarios) > 0 {
				all := true
				for _, sc := range e.Scenarios {
					cited := false
					for _, c := range cs {
						cited = cited || slices.Contains(c.Handles, sc.Handle)
					}
					all = all && cited
				}
				if all {
					out = append(out, fmt.Sprintf("%s: lint: %s is planned and each of its scenarios is verified; set it active", e.Pos(), e.ID))
				}
			}
			continue
		}
		if !e.Need && len(e.Serves) == 0 {
			out = append(out, fmt.Sprintf("%s: lint: %s serves nothing; name the need or decision it serves", e.Pos(), e.ID))
		}
		if !e.Need && len(as) > 0 && e.Kind != "structure" {
			onlyPackages := true
			for _, a := range as {
				if a.Kind == "realizes" && a.DeclKind != "package" {
					onlyPackages = false
				}
			}
			if onlyPackages {
				out = append(out, fmt.Sprintf("%s: lint: %s is realized only by a package clause; anchor the declaration that carries the decision", e.Pos(), e.ID))
			}
		}
		if len(e.Scenarios) > 0 {
			for _, c := range cs {
				if len(c.Handles) == 0 {
					out = append(out, fmt.Sprintf("%s: lint: verifies names %s without a scenario; name the handles the test exercises", c.Pos(), e.ID))
				}
			}
		}
	}
	return out
}

// questions says in one line, at the question file's first line, how many
// questions are open and which command lists them; none when none is.
//
//charter:realizes CH-261009-4B3B
func questions(rec *record.Record) []string {
	open := 0
	for _, q := range rec.Questions {
		if q.Open() {
			open++
		}
	}
	at := path.Join(rec.Dir, record.QuestionsFile)
	switch open {
	case 0:
		return nil
	case 1:
		return []string{fmt.Sprintf("%s:1: lint: 1 question is open; charter questions lists it", at)}
	}
	return []string{fmt.Sprintf("%s:1: lint: %d questions are open; charter questions lists them", at, open)}
}

// sharedTests names, once, each test whose citations name scenarios of two
// or more live decisions and that cites no need, at its first citation.
//
//charter:realizes CH-261126-9K5J
func sharedTests(t *check.Tree) []string {
	type key struct{ file, test string }
	type test struct {
		line  int
		cited []string
		seen  map[string]bool
		specs map[string]bool
		need  bool
	}
	tests := map[key]*test{}
	var order []key
	for _, c := range t.Facts.Citations {
		e, _ := t.Record.BySymbol(c.Symbol)
		if e == nil || c.Observes || c.Test == "" {
			continue
		}
		k := key{c.File, c.Test}
		tt := tests[k]
		if tt == nil {
			tt = &test{line: c.Line, seen: map[string]bool{}, specs: map[string]bool{}}
			tests[k] = tt
			order = append(order, k)
		}
		tt.line = min(tt.line, c.Line)
		switch {
		case e.Need:
			tt.need = true
			continue
		case e.Status == "deprecated" || len(c.Handles) == 0:
			continue
		}
		tt.specs[e.ID] = true
		for _, h := range c.Handles {
			if cited := e.ID + "/" + h; !tt.seen[cited] {
				tt.seen[cited] = true
				tt.cited = append(tt.cited, cited)
			}
		}
	}
	var out []string
	for _, k := range order {
		tt := tests[k]
		if len(tt.specs) < 2 || tt.need {
			continue
		}
		out = append(out, fmt.Sprintf("%s:%d: lint: %s cites %s; two decisions share one test, so one may rule what the other states", k.file, tt.line, k.test, andList(tt.cited)))
	}
	return out
}
