package chartertest

import (
	"fmt"
	"go/build/constraint"
	"time"
)

// isDate reports whether v is a real calendar date written YYYY-MM-DD.
func isDate(v string) bool {
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}

// isProbe reports whether the file's build constraint names the probe tag.
// A probe is a test file compiled only under that tag.
//
//charter:realizes CH-032
func isProbe(t *tree, rel, probeTag string) bool {
	if probeTag == "" || !isTest(rel) {
		return false
	}
	f := t.files[rel]
	for _, g := range f.Comments {
		if g.Pos() > f.Package {
			break
		}
		for _, c := range g.List {
			if !constraint.IsGoBuild(c.Text) {
				continue
			}
			expr, err := constraint.Parse(c.Text)
			if err != nil {
				continue
			}
			if namesTag(expr, probeTag) {
				return true
			}
		}
	}
	return false
}

// namesTag reports whether the constraint requires the tag on some path:
// it looks through and and or, and never into a negation.
func namesTag(x constraint.Expr, tag string) bool {
	switch x := x.(type) {
	case *constraint.TagExpr:
		return x.Tag == tag
	case *constraint.AndExpr:
		return namesTag(x.X, tag) || namesTag(x.Y, tag)
	case *constraint.OrExpr:
		return namesTag(x.X, tag) || namesTag(x.Y, tag)
	default:
		return false
	}
}

// checkObservations applies rules 11, 12 and 13. Without a probe tag no
// file can be a probe, so it reports that one configuration problem in place
// of the checks that ask whether a file is a probe.
//
//charter:realizes CH-032, CH-035
func (s *state) checkObservations() []string {
	var problems []string
	hasObservation := false
	tagged := s.r.ProbeTag != ""
	for _, name := range s.pkg.order {
		e := s.pkg.entries[name]
		isObs := e.Level == levelObservation
		if isObs {
			hasObservation = true
		}
		if tagged {
			for _, ref := range s.refs[name] {
				probe := isProbe(s.t, ref.File, s.r.ProbeTag)
				switch {
				case isObs && !probe:
					problems = append(problems, fmt.Sprintf("%s: rule 12: observation %s is referenced from %s, which is not a probe file", e.Pos, name, ref.File))
				case !isObs && probe:
					problems = append(problems, fmt.Sprintf("%s: rule 12: probe file %s references entry %s; a probe cites observations only", e.Pos, ref.File, name))
				}
			}
		}
		if !isObs {
			continue
		}
		if e.Subject == "" {
			problems = append(problems, fmt.Sprintf("%s: rule 11: %s has no Subject", e.Pos, name))
		}
		switch e.Status {
		case "observed":
			if e.Observed == "" {
				problems = append(problems, fmt.Sprintf("%s: rule 11: %s is observed and has no Observed date", e.Pos, name))
			}
			if tagged && !s.citedByProbe(name) {
				problems = append(problems, fmt.Sprintf("%s: rule 11: observed observation %s is cited by no probe", e.Pos, name))
			}
		case "documented":
			if e.Source == "" {
				problems = append(problems, fmt.Sprintf("%s: rule 11: %s is documented and has no Source", e.Pos, name))
			}
		case "broken":
			for _, p := range e.Supports {
				if parent := s.pkg.byID[p]; parent != nil && parent.Status == "active" {
					problems = append(problems, fmt.Sprintf("%s: rule 13: broken observation %s supports active spec %s", e.Pos, name, p))
				}
			}
		}
		if e.Observed != "" && !isDate(e.Observed) {
			problems = append(problems, fmt.Sprintf("%s: rule 11: %s has Observed %q; want YYYY-MM-DD", e.Pos, name, e.Observed))
		}
	}
	if hasObservation && !tagged {
		problems = append(problems, "rule 11: observations need Options.ProbeTag; set it to the build tag probe files carry")
	}
	return problems
}

// citedByProbe reports whether any reference to the entry is from a probe.
func (s *state) citedByProbe(name string) bool {
	for _, ref := range s.refs[name] {
		if isProbe(s.t, ref.File, s.r.ProbeTag) {
			return true
		}
	}
	return false
}
