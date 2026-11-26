package chartertest

import (
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"testing"
)

// UpdateEnv is the environment variable that makes Check rewrite the
// baseline with every surface no entry accounts for, instead of failing.
const UpdateEnv = "CHARTER_UPDATE_BASELINE"

const initialReason = "unclaimed at baseline init"

// state is one evaluation of the tree against the options.
type state struct {
	r         *resolved
	t         *tree
	env       map[string]map[string]string // import path -> constant name -> value
	pkg       *entriesPackage
	docs      []*areaDoc
	refs      map[string][]Ref // every reference to an entry from a test file
	cites     map[string][]Ref // the Verifies calls, by spec
	anchors   map[string][]Anchor
	questions []Question
	problems  []string
}

// load parses the tree, the env packages, the area documents, the spec
// files and the observations, and collects the citations and anchors. The
// form, the relations and the anchor rules are applied here.
func load(r *resolved) (*state, error) {
	t, err := loadTree(r)
	if err != nil {
		return nil, err
	}
	s := &state{r: r, t: t, env: map[string]map[string]string{}}
	for p, d := range r.envDirs {
		s.env[p] = envConstants(t, d)
	}
	s.pkg = newEntriesPackage()
	docs, ps, err := readAreas(r.Root, r.areasDir, s.pkg, r.notAreaDocuments()...)
	if err != nil {
		return nil, err
	}
	s.docs = docs
	s.problems = append(s.problems, ps...)
	ps, err = readSpecs(r.Root, r.specsDir, s.pkg)
	if err != nil {
		return nil, err
	}
	s.problems = append(s.problems, ps...)
	s.questions, ps, err = readQuestions(r, s.pkg)
	if err != nil {
		return nil, err
	}
	s.problems = append(s.problems, ps...)
	s.problems = append(s.problems, parseObservations(t, r, s.pkg)...)
	s.problems = append(s.problems, checkMutations(t, r, s.pkg)...)
	s.problems = append(s.problems, checkRelations(s.pkg)...)
	s.refs, s.cites, ps = collectRefs(t, r, s.pkg)
	s.problems = append(s.problems, ps...)
	s.anchors, ps, err = collectAnchors(t, r, s.pkg)
	if err != nil {
		return nil, err
	}
	s.problems = append(s.problems, ps...)
	if r.TypeScript != nil {
		tsRefs, tsCites, ps, err := collectTSRefs(t, r, s.pkg)
		if err != nil {
			return nil, err
		}
		s.problems = append(s.problems, ps...)
		for name, rs := range tsRefs {
			s.refs[name] = append(s.refs[name], rs...)
		}
		for name, rs := range tsCites {
			s.cites[name] = append(s.cites[name], rs...)
		}
	}
	return s, nil
}

// Check evaluates every rule and reports each problem as a test error.
//
//charter:realizes CH-037
func Check(tb testing.TB, opts Options) {
	tb.Helper()
	problems, err := problems(opts, callerDir(1))
	if err != nil {
		tb.Fatalf("charter check: %v", err)
	}
	for _, p := range problems {
		tb.Error(p)
	}
}

// Problems evaluates every rule and returns the problems, sorted. Skills
// and scripts use it where Check's test handle is not at hand.
//
//charter:realizes CH-037
func Problems(opts Options) ([]string, error) {
	return problems(opts, callerDir(1))
}

func problems(opts Options, callerDir string) ([]string, error) {
	r, err := opts.resolve(callerDir)
	if err != nil {
		return nil, err
	}
	if os.Getenv(UpdateSpecsEnv) != "" {
		if err := writeSpecs(r); err != nil {
			return nil, err
		}
	}
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	if err := s.evaluate(); err != nil {
		return nil, err
	}
	sort.Strings(s.problems)
	return s.problems, nil
}

// evaluate applies the rules that need the whole tree once it is loaded:
// environment reads, surfaces, coverage, observations, the generated files
// and the baseline.
//
//charter:realizes CH-033
func (s *state) evaluate() error {
	s.problems = append(s.problems, checkEnvReads(s.t, s.r, s.env)...)

	surfaces, err := s.enumerate()
	if err != nil {
		return err
	}

	accounted := map[string]string{} // key -> spec
	for _, name := range s.pkg.order {
		e := s.pkg.entries[name]
		// A deprecated entry owns nothing: what it accounted for is a live
		// entry's to claim, or the baseline's.
		if e.Level != levelSpec || e.Status == "deprecated" {
			continue
		}
		for _, k := range e.Accounts {
			if !surfaces[k] {
				s.problems = append(s.problems, fmt.Sprintf("%s: rule 7: %s accounts for %s, which no enumerator found", e.Pos, e.ID, k))
			}
			if _, ok := accounted[k]; !ok {
				accounted[k] = e.ID
			}
		}
	}
	ps, err := s.checkCoverage()
	if err != nil {
		return err
	}
	s.problems = append(s.problems, ps...)
	s.problems = append(s.problems, s.checkObservations()...)

	ps, err = checkGenerated(s.r, s.pkg, false)
	if err != nil {
		return err
	}
	s.problems = append(s.problems, ps...)

	if s.r.TypeScript != nil {
		ps, err := checkTypeScriptModule(s.r, s.pkg, os.Getenv(UpdateTypeScriptEnv) != "")
		if err != nil {
			return err
		}
		s.problems = append(s.problems, ps...)
	}

	ps, err = checkDigest(s.r, s, os.Getenv(UpdateDigestEnv) != "")
	if err != nil {
		return err
	}
	s.problems = append(s.problems, ps...)

	display := path.Join(s.r.entriesDir, baselineFile)
	baselinePath := filepath.Join(s.r.Root, filepath.FromSlash(s.r.entriesDir), baselineFile)
	lines, ps, err := readBaseline(baselinePath, display)
	if err != nil {
		return err
	}
	if os.Getenv(UpdateEnv) != "" {
		lines = rebuildBaseline(lines, surfaces, accounted)
		if err := writeBaseline(baselinePath, lines); err != nil {
			return err
		}
		ps = nil
	}
	s.problems = append(s.problems, ps...)
	listed := map[string]bool{}
	for _, l := range lines {
		if listed[l.Key] {
			s.problems = append(s.problems, fmt.Sprintf("%s:%d: rule 7: baseline lists %s twice", display, l.Line, l.Key))
		}
		listed[l.Key] = true
		switch {
		case !surfaces[l.Key]:
			s.problems = append(s.problems, fmt.Sprintf("%s:%d: rule 7: baseline lists %s, which no enumerator found", display, l.Line, l.Key))
		case accounted[l.Key] != "":
			s.problems = append(s.problems, fmt.Sprintf("%s:%d: rule 7: baseline lists %s, which %s now accounts for", display, l.Line, l.Key, accounted[l.Key]))
		}
	}
	for _, k := range slices.Sorted(maps.Keys(surfaces)) {
		if accounted[k] == "" && !listed[k] {
			s.problems = append(s.problems, fmt.Sprintf("rule 6: %s is neither accounted for by an entry nor listed in %s", k, display))
		}
	}
	return nil
}

// enumerate runs every enumerator and applies rule 9.
//
//charter:realizes CH-033
func (s *state) enumerate() (map[string]bool, error) {
	type kind struct {
		name string
		keys []string
	}
	var kinds []kind
	if len(s.r.EnvPackages) > 0 {
		var keys []string
		for _, consts := range s.env {
			for _, v := range consts {
				keys = append(keys, "env:"+v)
			}
		}
		kinds = append(kinds, kind{"env", keys})
	}
	kinds = append(kinds, kind{"tag", enumTags(s.t, s.r.ProbeTag)})
	if len(s.r.OpenAPI) > 0 {
		keys, err := enumOpenAPI(s.r.Root, s.r.OpenAPI)
		if err != nil {
			return nil, err
		}
		kinds = append(kinds, kind{"op", keys})
	}
	if s.r.Events != "" {
		keys, err := enumEvents(s.r.Root, s.r.Events)
		if err != nil {
			return nil, err
		}
		kinds = append(kinds, kind{"event", keys})
	}
	for _, e := range s.r.Enumerators {
		keys, err := e.Keys(s.r.Root)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name, err)
		}
		kinds = append(kinds, kind{e.Name, keys})
	}
	surfaces := map[string]bool{}
	for _, k := range kinds {
		// Zero build tags is a valid tree, while the other kinds are
		// configured by path, so an empty result there is a misconfiguration.
		if len(k.keys) == 0 && k.name != "tag" {
			s.problems = append(s.problems, fmt.Sprintf("rule 9: the %s enumerator found nothing; check Options", k.name))
		}
		for _, key := range k.keys {
			surfaces[key] = true
		}
	}
	return surfaces, nil
}

// rebuildBaseline keeps every line that is still unaccounted and present,
// with its reason, and adds every other unaccounted surface.
//
//charter:realizes CH-033
func rebuildBaseline(old []baselineLine, surfaces map[string]bool, accounted map[string]string) []baselineLine {
	reasons := map[string]string{}
	for _, l := range old {
		if l.Reason != "" {
			reasons[l.Key] = l.Reason
		}
	}
	var out []baselineLine
	for _, k := range slices.Sorted(maps.Keys(surfaces)) {
		if accounted[k] != "" {
			continue
		}
		reason := reasons[k]
		if reason == "" {
			reason = initialReason
		}
		out = append(out, baselineLine{Key: k, Reason: reason})
	}
	return out
}

// writeSpecs writes the Go generated from the specs.
func writeSpecs(r *resolved) error {
	s, err := load(r)
	if err != nil {
		return err
	}
	_, err = checkGenerated(r, s.pkg, true)
	return err
}

// Generate writes everything the gate compares to a rendering: the Go
// generated from the specs, the digest, then the TypeScript module. A command such as
// cmd/claimsgen calls it, so the specs can be regenerated even when the
// entries package does not compile.
//
//charter:realizes CH-040
func Generate(opts Options) error {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return err
	}
	if err := writeSpecs(r); err != nil {
		return err
	}
	s, err := load(r)
	if err != nil {
		return err
	}
	if _, err := checkDigest(r, s, true); err != nil {
		return err
	}
	if r.TypeScript == nil {
		return nil
	}
	_, err = checkTypeScriptModule(r, s.pkg, true)
	return err
}
