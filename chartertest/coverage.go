package chartertest

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Coverage: an active spec is covered when production code realizes it and
// a test verifies each of its scenarios; a spec without scenarios, in the
// earlier form or planned, is verified by any citation. Rule 25 fails an
// active spec that is not covered, unless a waiver accepts a missing test.
// A waiver never accepts missing code. A need is realized by the decisions
// that serve it rather than by code, so it needs tests and no anchor, and
// its tests are on the stakeholder's side. concerns never counts.

// waiver is one line of the waiver file.
type waiver struct {
	Pattern, Reason string
	Line            int
	used            bool
}

// readWaivers reads the waiver file; a missing file is no waivers.
//
//charter:realizes CH-030
func readWaivers(r *resolved) ([]*waiver, string, []string, error) {
	rel := path.Join(r.entriesDir, waiversFile)
	f, err := os.Open(filepath.Join(r.Root, filepath.FromSlash(rel)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, rel, nil, nil
	}
	if err != nil {
		return nil, rel, nil, err
	}
	defer f.Close()
	var ws []*waiver
	var problems []string
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		n++
		l := sc.Text()
		if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "#") {
			continue
		}
		pat, reason, _ := strings.Cut(l, "\t")
		pat, reason = strings.TrimSpace(pat), strings.TrimSpace(reason)
		if strings.ContainsAny(pat, "*?[") {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 25: the waiver %q is a glob; a waiver names one spec id, SW-118, or one scenario, SW-118/handle, so what it accepts is on its line", rel, n, pat))
			continue
		}
		if pat == "" {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 25: a waiver names a spec id, SW-118, or a scenario, SW-118/handle, then a tab and the reason", rel, n))
			continue
		}
		if reason == "" {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 25: the waiver for %s gives no reason after a tab", rel, n, pat))
			continue
		}
		ws = append(ws, &waiver{Pattern: pat, Reason: reason, Line: n})
	}
	return ws, rel, problems, sc.Err()
}

// checkCoverage applies rule 25, and the rule 24 half that needs the
// citations: a handle a test names is one the spec lists.
//
//charter:realizes CH-022, CH-030, CH-031
func (s *state) checkCoverage() ([]string, error) {
	waivers, rel, problems, err := readWaivers(s.r)
	if err != nil {
		return nil, err
	}
	waived := func(key string) bool {
		hit := false
		for _, w := range waivers {
			if ok, _ := path.Match(w.Pattern, key); ok {
				w.used, hit = true, true
			}
		}
		return hit
	}
	for _, name := range s.pkg.order {
		e := s.pkg.entries[name]
		if e.Level != levelSpec {
			continue
		}
		realized := false
		for _, a := range s.anchors[name] {
			if a.Kind == "realizes" {
				realized = true
			}
		}
		var cites []Ref
		for _, c := range s.cites[name] {
			if isProbe(s.t, c.File, s.r.ProbeTag) {
				continue
			}
			if e.isNeed() && !s.r.stakeholderSide(c.File) {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 25: need %s is verified from %s, which is not on the stakeholder's side; a need is verified from the person's side: the command line, the routes or the browser", c.File, c.Line, e.ID, c.File))
				continue
			}
			for _, h := range c.Handles {
				if !e.hasScenario(h) {
					problems = append(problems, fmt.Sprintf("%s:%d: rule 24: %s has no scenario %q", c.File, c.Line, e.ID, h))
				}
			}
			cites = append(cites, c)
		}
		if e.Status != "active" {
			continue
		}
		if !e.isNeed() && !realized {
			problems = append(problems, fmt.Sprintf("%s: rule 25: active spec %s is realized by no code; add //charter:realizes %s where it is implemented", e.Pos, e.ID, e.ID))
		}
		if len(e.Scenarios) == 0 {
			if len(cites) == 0 && !waived(e.ID) {
				problems = append(problems, fmt.Sprintf("%s: rule 25: active spec %s is verified by no test; cite it with Verifies, or waive it in %s with a reason", e.Pos, e.ID, rel))
			}
			continue
		}
		for _, sc := range e.Scenarios {
			cited := false
			for _, c := range cites {
				for _, h := range c.Handles {
					if h == sc.Handle {
						cited = true
					}
				}
			}
			if cited || waived(e.ID+"/"+sc.Handle) {
				continue
			}
			problems = append(problems, fmt.Sprintf("%s:%d: rule 25: scenario %s of %s is verified by no test; cite it with Verifies(t, claims.%s, %q), or waive %s/%s in %s with a reason", e.File, sc.Line, sc.Handle, e.ID, e.Name, sc.Handle, e.ID, sc.Handle, rel))
		}
	}
	for _, w := range waivers {
		if !w.used {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 25: the waiver for %s absorbs no active spec or scenario that lacks a test", rel, w.Line, w.Pattern))
		}
	}
	return problems, nil
}

// isGeneratedSpecFile reports whether a Go file was generated from the
// specs: its first comment, before the package clause, opens with
// generatedPrefix.
func isGeneratedSpecFile(f *ast.File) bool {
	return len(f.Comments) > 0 && len(f.Comments[0].List) > 0 &&
		strings.HasPrefix(f.Comments[0].List[0].Text, generatedPrefix) &&
		f.Comments[0].Pos() < f.Package
}
