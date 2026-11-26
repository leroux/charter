// Package gate holds a charter's record against the facts its readers
// found: every anchor and citation names what the record holds, every
// active entry is realized and each scenario cited, every surface is
// accounted for, and the observations carry their fields. It reads facts
// alone and imports no language, so a rule holds the same way in each.
//
//charter:realizes CH-261126-PQ7X
package gate

import (
	"fmt"
	"maps"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
	"github.com/leroux/charter/internal/surface"
)

// Input is everything the rules read.
type Input struct {
	Record *record.Record
	Facts  *facts.Set
	// StakeholderTests are patterns of the test files a need is verified
	// from; empty means any test file.
	StakeholderTests []string
	// Probes are the slash paths, relative to the root, that hold probes.
	Probes []string
	// EnvPatterns are regular expressions a literal naming an undeclared
	// variable matches.
	EnvPatterns []*regexp.Regexp
	// EnvModules are the constants modules configured, each "file:line" of
	// the configuration that names it, by module path.
	EnvModules map[string]string
	// Declared are the keys of the declared surface kinds.
	Declared []string
	// Baseline is the baseline's lines, and BaselineFile how problems name
	// it.
	Baseline     []surface.Line
	BaselineFile string
	// ClaimsSources are the source files under the claims directory, and
	// Generated the ones the generators write there.
	ClaimsSources []string
	Generated     []string
}

// Check applies every rule and returns the problems, sorted, each
// "file:line: rule N: what is wrong and what to do".
//
//charter:realizes CH-037
func Check(in Input) []string {
	c := &checker{Input: in, realized: map[string]bool{}, cites: map[string][]facts.Citation{}}
	for _, r := range in.Facts.Refused {
		c.add("%s: %s: %s", r.Pos(), r.Rule, r.Message)
	}
	c.handSources()
	c.anchors()
	c.citations()
	c.relations()
	c.coverage()
	c.observations()
	c.env()
	c.surfaces()
	sort.Strings(c.problems)
	return c.problems
}

type checker struct {
	Input
	problems []string
	realized map[string]bool             // entry id -> a realizes anchor names it
	cites    map[string][]facts.Citation // entry id -> its counted citations
}

func (c *checker) add(format string, args ...any) {
	c.problems = append(c.problems, fmt.Sprintf(format, args...))
}

// waiversPath and baselinePath name the record's files in problems.
func (c *checker) waiversPath() string { return path.Join(c.Record.Dir, record.WaiversFile) }

// handSources fails a source file under the claims directory that no
// generator writes: entries and observations are written in markdown.
//
//charter:realizes CH-025
func (c *checker) handSources() {
	for _, f := range c.ClaimsSources {
		if !slices.Contains(c.Generated, f) {
			c.add("%s:1: rule 1: %s is written by hand in the claims directory, which holds the record and the generated files; entries and observations are written in markdown", f, f)
		}
	}
}

// anchors holds each directive to an entry that exists and is not planned.
//
//charter:realizes CH-022
func (c *checker) anchors() {
	for _, a := range c.Facts.Anchors {
		for _, id := range a.IDs {
			e := c.Record.Entry(id)
			switch {
			case e == nil:
				c.add("%s: rule 24: charter:%s names %s, which is not an entry of the record", a.Pos(), a.Kind, id)
			case e.Status == "planned":
				c.add("%s: rule 24: charter:%s names planned entry %s; set it active as its work begins", a.Pos(), a.Kind, id)
			case a.Kind == "realizes":
				c.realized[id] = true
			}
		}
	}
}

// isProbe reports whether a test file is under a probe path.
//
//charter:realizes CH-032
func (c *checker) isProbe(file string) bool {
	for _, p := range c.Probes {
		p = strings.Trim(path.Clean(p), "/")
		if file == p || strings.HasPrefix(file, p+"/") {
			return true
		}
	}
	return false
}

// citations resolves each citation and reference: a verifies names an
// entry and handles it lists, an observes an observation, from a probe; a
// need is cited from the stakeholder's side.
//
//charter:realizes CH-023, CH-031, CH-032
func (c *checker) citations() {
	for _, ct := range c.Facts.Citations {
		e, o := c.Record.BySymbol(ct.Symbol)
		probe := c.isProbe(ct.File)
		switch {
		case ct.Observes && o == nil:
			c.add("%s: rule 24: observes names %s, which is no observation of the record", ct.Pos(), ct.Symbol)
		case ct.Observes && !probe:
			c.add("%s: rule 12: observation %s is cited from %s, which is not a probe; a test never cites an observation", ct.Pos(), ct.Symbol, ct.File)
		case ct.Observes:
		case e == nil:
			c.add("%s: rule 24: verifies names claims.%s, which is no entry of the record; name an entry by its symbol, such as claims.SW_118", ct.Pos(), ct.Symbol)
		case probe:
			c.add("%s: rule 12: probe %s cites entry %s; a probe cites observations only, and its citations count for nothing", ct.Pos(), ct.File, e.ID)
		case e.Need && len(c.StakeholderTests) > 0 && !matchesAny(c.StakeholderTests, ct.File):
			c.add("%s: rule 25: need %s is verified from %s, which is not on the stakeholder's side; a need is verified from the tests charter.toml names in stakeholder-tests", ct.Pos(), e.ID, ct.File)
		default:
			ok := true
			for _, h := range ct.Handles {
				if !e.HasScenario(h) {
					c.add("%s: rule 24: %s has no scenario %q", ct.Pos(), e.ID, h)
					ok = false
				}
			}
			if ok {
				c.cites[e.ID] = append(c.cites[e.ID], ct)
			}
		}
	}
	for _, r := range c.Facts.References {
		e, o := c.Record.BySymbol(r.Symbol)
		probe := c.isProbe(r.File)
		switch {
		case o != nil && !probe:
			c.add("%s: rule 12: observation %s is referenced from %s, which is not a probe", r.Pos(), r.Symbol, r.File)
		case e != nil && probe:
			c.add("%s: rule 12: probe %s references entry %s; a probe cites observations only", r.Pos(), r.File, e.ID)
		}
	}
}

// relations holds Serves, See and an observation's Supports to entries the
// record holds, and nothing to itself.
//
//charter:realizes CH-015
func (c *checker) relations() {
	for _, e := range c.Record.Entries {
		for _, rel := range []struct {
			name string
			ids  []string
		}{{"serves", e.Serves}, {"sees", e.See}} {
			for _, id := range rel.ids {
				switch {
				case id == e.ID:
					c.add("%s: rule 23: %s %s itself", e.Pos(), e.ID, rel.name)
				case c.Record.Entry(id) == nil:
					c.add("%s: rule 23: %s %s %s, which is not an entry", e.Pos(), e.ID, rel.name, id)
				}
			}
		}
	}
	for _, o := range c.Record.Observations {
		for _, id := range o.Supports {
			if c.Record.Entry(id) == nil {
				c.add("%s: rule 23: observation %s supports %s, which is not an entry", o.Pos(), o.Name, id)
			}
		}
	}
}

// coverage fails an active decision no realizes anchor names, and an active
// entry's scenario no citation names, unless a waiver accepts the missing
// test; a waiver never accepts missing code, and is never a glob.
//
//charter:realizes CH-030
func (c *checker) coverage() {
	type waiver struct {
		record.Waiver
		used bool
	}
	var ws []*waiver
	for _, w := range c.Record.Waivers {
		switch {
		case strings.ContainsAny(w.Target, "*?["):
			c.add("%s:%d: rule 25: the waiver %q is a glob; a waiver names one entry, SW-118, or one scenario, SW-118/handle, so what it accepts is on its line", c.waiversPath(), w.Line, w.Target)
		case w.Target == "":
			c.add("%s:%d: rule 25: a waiver names an entry, SW-118, or a scenario, SW-118/handle, then a tab and the reason", c.waiversPath(), w.Line)
		case w.Reason == "":
			c.add("%s:%d: rule 25: the waiver for %s gives no reason after a tab", c.waiversPath(), w.Line, w.Target)
		default:
			ws = append(ws, &waiver{Waiver: w})
		}
	}
	waived := func(key string) bool {
		hit := false
		for _, w := range ws {
			if w.Target == key {
				w.used, hit = true, true
			}
		}
		return hit
	}
	for _, e := range c.Record.Entries {
		if e.Status != "active" {
			continue
		}
		if !e.Need && !c.realized[e.ID] {
			c.add("%s: rule 25: active entry %s is realized by no code; add charter:realizes %s where it is implemented", e.Pos(), e.ID, e.ID)
		}
		cites := c.cites[e.ID]
		if len(e.Scenarios) == 0 {
			if len(cites) == 0 && !waived(e.ID) {
				c.add("%s: rule 25: active entry %s is verified by no test; cite it with verifies as a test's first statement, or waive it in %s with a reason", e.Pos(), e.ID, c.waiversPath())
			}
			continue
		}
		for _, sc := range e.Scenarios {
			cited := false
			for _, ct := range cites {
				cited = cited || slices.Contains(ct.Handles, sc.Handle)
			}
			if cited || waived(e.ID+"/"+sc.Handle) {
				continue
			}
			c.add("%s:%d: rule 25: scenario %s of %s is verified by no test; cite it with verifies(claims.%s, %q) as a test's first statement, or waive %s/%s in %s with a reason", e.File, sc.Line, sc.Handle, e.ID, e.Symbol, sc.Handle, e.ID, sc.Handle, c.waiversPath())
		}
	}
	for _, w := range ws {
		if !w.used {
			c.add("%s:%d: rule 25: the waiver for %s absorbs no active entry or scenario that lacks a test", c.waiversPath(), w.Line, w.Target)
		}
	}
}

// observations holds each observation to its status's fields, an observed
// one to a probe that cites it, and a broken one to no active entry.
//
//charter:realizes CH-035
func (c *checker) observations() {
	if len(c.Record.Observations) > 0 && len(c.Probes) == 0 {
		c.add("%s:1: rule 11: the record holds observations and charter.toml names no probe path; set probes to the paths that hold them", path.Join(c.Record.Dir, "charter.toml"))
	}
	probed := map[string]bool{}
	for _, ct := range c.Facts.Citations {
		if ct.Observes && c.isProbe(ct.File) {
			probed[ct.Symbol] = true
		}
	}
	for _, o := range c.Record.Observations {
		if o.Subject == "" {
			c.add("%s: rule 11: %s has no Subject", o.Pos(), o.Name)
		}
		if o.Observed != "" && !isDate(o.Observed) {
			c.add("%s: rule 11: %s has Observed %q; want YYYY-MM-DD", o.Pos(), o.Name, o.Observed)
		}
		switch o.Status {
		case "observed":
			if o.Observed == "" {
				c.add("%s: rule 11: %s is observed and has no Observed date", o.Pos(), o.Name)
			}
			if len(c.Probes) > 0 && !probed[o.Name] {
				c.add("%s: rule 11: observed observation %s is cited by no probe", o.Pos(), o.Name)
			}
		case "documented":
			if o.Source == "" {
				c.add("%s: rule 11: %s is documented and has no Source", o.Pos(), o.Name)
			}
			if o.Observed != "" {
				c.add("%s: rule 11: %s is documented and carries an Observed date; a date is a probe's", o.Pos(), o.Name)
			}
		case "broken":
			for _, id := range o.Supports {
				if e := c.Record.Entry(id); e != nil && e.Status == "active" {
					c.add("%s: rule 13: broken observation %s supports active entry %s", o.Pos(), o.Name, id)
				}
			}
		}
	}
}

// env holds every read of the environment to a constant of a constants
// module, or to the dynamic directive, and refuses a literal that names a
// variable.
//
//charter:realizes CH-034
func (c *checker) env() {
	declared := map[string]string{} // value -> constant
	for _, k := range c.Facts.EnvConsts {
		declared[k.Value] = k.Name
	}
	for _, r := range c.Facts.EnvReads {
		if r.Const != "" || r.Dynamic {
			continue
		}
		if r.Literal != "" {
			if name, ok := declared[r.Literal]; ok {
				c.add("%s: rule 8: an environment read names %q by a literal; read it through %s", r.Pos(), r.Literal, name)
				continue
			}
			c.add("%s: rule 8: an environment read names %q by a literal; declare it in a constants module and read it through the constant", r.Pos(), r.Literal)
			continue
		}
		c.add("%s: rule 8: an environment read must name a constant of a constants module; a read named by configuration carries charter:dynamic", r.Pos())
	}
	for _, l := range c.Facts.Literals {
		if name, ok := declared[l.Value]; ok {
			c.add("%s: rule 8: literal %q duplicates %s; use the constant", l.Pos(), l.Value, name)
			continue
		}
		for _, re := range c.EnvPatterns {
			if re.MatchString(l.Value) {
				c.add("%s: rule 8: literal %q looks like an environment variable; declare it in a constants module", l.Pos(), l.Value)
				break
			}
		}
	}
	for _, mod := range slices.Sorted(maps.Keys(c.EnvModules)) {
		found := false
		for _, k := range c.Facts.EnvConsts {
			found = found || k.File == mod || path.Dir(k.File) == mod
		}
		if !found {
			c.add("%s: rule 9: the constants module %s declares no variable", c.EnvModules[mod], mod)
		}
	}
}

// Surfaces returns every surface key, the environment variables' and the
// declared kinds', and the live entry that accounts for each one accounted.
//
//charter:realizes CH-033
func Surfaces(in Input) (surfaces map[string]bool, accounted map[string]string) {
	surfaces, accounted = map[string]bool{}, map[string]string{}
	for _, k := range in.Facts.EnvConsts {
		surfaces["env:"+k.Value] = true
	}
	for _, k := range in.Declared {
		surfaces[k] = true
	}
	for _, e := range in.Record.Entries {
		// A deprecated entry owns nothing: what it accounted for is a live
		// entry's to claim, or the baseline's.
		if e.Status == "deprecated" {
			continue
		}
		for _, k := range e.Accounts {
			if _, ok := accounted[k]; !ok {
				accounted[k] = e.ID
			}
		}
	}
	return surfaces, accounted
}

// surfaces holds every surface to an entry's Accounts or a baseline line,
// every Accounts to a surface, and the baseline to live, unaccounted
// surfaces, each once.
//
//charter:realizes CH-033
func (c *checker) surfaces() {
	surfaces, accounted := Surfaces(c.Input)
	for _, e := range c.Record.Entries {
		if e.Status == "deprecated" {
			continue
		}
		for _, k := range e.Accounts {
			if !surfaces[k] {
				c.add("%s: rule 7: %s accounts for %s, which no enumerator found", e.Pos(), e.ID, k)
			}
		}
	}
	listed := map[string]bool{}
	for _, l := range c.Baseline {
		if listed[l.Key] {
			c.add("%s:%d: rule 7: baseline lists %s twice", c.BaselineFile, l.Line, l.Key)
		}
		listed[l.Key] = true
		switch {
		case !surfaces[l.Key]:
			c.add("%s:%d: rule 7: baseline lists %s, which no enumerator found", c.BaselineFile, l.Line, l.Key)
		case accounted[l.Key] != "":
			c.add("%s:%d: rule 7: baseline lists %s, which %s now accounts for", c.BaselineFile, l.Line, l.Key, accounted[l.Key])
		}
	}
	for _, k := range slices.Sorted(maps.Keys(surfaces)) {
		if accounted[k] == "" && !listed[k] {
			c.add("rule 6: %s is neither accounted for by an entry nor listed in %s", k, c.BaselineFile)
		}
	}
}

func matchesAny(patterns []string, rel string) bool {
	base := path.Base(rel)
	for _, pat := range patterns {
		if ok, _ := path.Match(pat, rel); ok {
			return true
		}
		if ok, _ := path.Match(pat, base); ok {
			return true
		}
	}
	return false
}

func isDate(v string) bool {
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}
