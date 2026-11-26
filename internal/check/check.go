// Package check runs a charter's gate: it reads the record, walks the tree,
// has each configured language's reader turn its files into facts, runs the
// declared surface kinds, compares the generated files with their
// rendering, and hands all of it to the gate's rules.
package check

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/digest"
	"github.com/leroux/charter/internal/env"
	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/gate"
	"github.com/leroux/charter/internal/lang/golang"
	"github.com/leroux/charter/internal/lang/python"
	"github.com/leroux/charter/internal/lang/typescript"
	"github.com/leroux/charter/internal/record"
	"github.com/leroux/charter/internal/surface"
	"github.com/leroux/charter/internal/walk"
)

// Regen is the command a stale generated file names.
const Regen = "charter gen"

// Tree is one charter read whole.
type Tree struct {
	Config *config.Config
	Record *record.Record
	Facts  *facts.Set
	// Files are the walked files by language.
	Files walk.Tree
	// Problems found while reading, before the rules run.
	Problems []string
}

// Load reads the charter's record and its tree's facts.
func Load(c *config.Config) (*Tree, error) {
	skip := []string{}
	if c.Digest != "" {
		skip = append(skip, path.Base(c.Digest))
	}
	rec, err := record.Read(c.Root, config.Dir, skip...)
	if err != nil {
		return nil, err
	}
	t := &Tree{Config: c, Record: rec, Facts: &facts.Set{}}
	t.Problems = append(t.Problems, rec.Problems...)
	t.Files, err = walk.Walk(walk.Options{
		Root:      c.Root,
		Languages: languages(c),
		Ignore:    c.Ignore,
		Only:      c.Only,
		Boundary:  config.Dir + "/" + config.File,
	})
	if err != nil {
		return nil, err
	}
	if c.Go != nil {
		set, err := golang.Reader{
			Root: c.Root, Claims: c.Go.Module, ClaimsDir: config.Dir,
			EnvDir: strings.Trim(c.Go.Env, "/"),
		}.Read(t.Files[walk.Go])
		if err != nil {
			return nil, err
		}
		t.Facts.Add(set)
	}
	if c.Python != nil {
		set, err := python.Reader{
			Root: c.Root, Dir: c.Python.Dir, EnvFile: strings.Trim(c.Python.Env, "/"),
			ClaimsDir: config.Dir, Table: fmt.Sprintf("%s:%d", c.Rel(), c.Line("python")),
		}.Read(t.Files[walk.Python])
		if err != nil {
			return nil, err
		}
		t.Facts.Add(set)
	}
	if ts := c.TypeScript; ts != nil {
		set, err := typescript.Reader{
			Root: c.Root, Compiler: ts.Compiler, Import: ts.Import, EnvFile: strings.Trim(ts.Env, "/"),
			ClaimsDir: config.Dir, Table: fmt.Sprintf("%s:%d", c.Rel(), c.Line("typescript")),
		}.Read(t.Files[walk.TypeScript])
		if err != nil {
			return nil, err
		}
		t.Facts.Add(set)
	}
	return t, nil
}

func languages(c *config.Config) []string {
	var ls []string
	if c.Go != nil {
		ls = append(ls, walk.Go)
	}
	if c.Python != nil {
		ls = append(ls, walk.Python)
	}
	if c.TypeScript != nil {
		ls = append(ls, walk.TypeScript)
	}
	return ls
}

// Generated renders every file charter gen writes, by slash path relative
// to the root.
func Generated(c *config.Config, rec *record.Record) (map[string][]byte, error) {
	out := map[string][]byte{}
	if c.Go != nil {
		gomod, specs, err := golang.Generate(rec, c.Go.Module, Regen)
		if err != nil {
			return nil, err
		}
		out[config.Dir+"/go.mod"] = gomod
		out[config.Dir+"/specs.go"] = specs
	}
	if c.Python != nil {
		out[config.Dir+"/"+python.ClaimsModule+".py"] = python.Generate(rec, Regen)
	}
	if c.TypeScript != nil {
		out[config.Dir+"/index.ts"] = typescript.Generate(rec, Regen)
	}
	if c.Digest != "" {
		out[path.Join(config.Dir, c.Digest)] = []byte(digest.Render(rec, c.Digest, Regen))
	}
	return out, nil
}

// Write writes every generated file.
//
//charter:realizes CH-020
func Write(c *config.Config) error {
	rec, err := record.Read(c.Root, config.Dir, path.Base(c.Digest))
	if err != nil {
		return err
	}
	files, err := Generated(c, rec)
	if err != nil {
		return err
	}
	for rel, data := range files {
		if err := os.WriteFile(filepath.Join(c.Root, filepath.FromSlash(rel)), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// Problems runs the gate over the charter and returns every problem,
// sorted.
//
//charter:realizes CH-037
func Problems(c *config.Config) ([]string, error) {
	_, problems, err := Run(c)
	return problems, err
}

// Run loads the charter and runs the gate, returning the tree it read and
// every problem, sorted.
func Run(c *config.Config) (*Tree, []string, error) {
	t, err := Load(c)
	if err != nil {
		return nil, nil, err
	}
	problems := append([]string{}, t.Problems...)

	generated, err := Generated(c, t.Record)
	if err != nil {
		return nil, nil, err
	}
	var genFiles []string
	for rel, want := range generated {
		genFiles = append(genFiles, rel)
		got, err := os.ReadFile(filepath.Join(c.Root, filepath.FromSlash(rel)))
		switch {
		case errors.Is(err, fs.ErrNotExist):
			problems = append(problems, fmt.Sprintf("%s:1: rule 15: %s is missing; run %s", rel, rel, Regen))
		case err != nil:
			return nil, nil, err
		case !bytes.Equal(got, want):
			problems = append(problems, fmt.Sprintf("%s:1: rule 15: %s is stale: it differs from what the record renders; run %s", rel, rel, Regen))
		}
	}

	in := gate.Input{
		Record:           t.Record,
		Facts:            t.Facts,
		StakeholderTests: c.StakeholderTests,
		Probes:           c.Probes,
		EnvModules:       map[string]string{},
		BaselineFile:     config.Dir + "/" + record.BaselineFile,
		Generated:        genFiles,
	}
	for _, f := range allFiles(t.Files) {
		if strings.HasPrefix(f, config.Dir+"/") {
			in.ClaimsSources = append(in.ClaimsSources, f)
		}
	}
	for _, p := range c.EnvPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 9: env-patterns %q is not a regular expression: %v", c.Rel(), c.Line("env-patterns"), p, err))
			continue
		}
		in.EnvPatterns = append(in.EnvPatterns, re)
	}
	if c.Go != nil && c.Go.Env != "" {
		in.EnvModules[strings.Trim(c.Go.Env, "/")] = fmt.Sprintf("%s:%d", c.Rel(), c.Line("go.env"))
	}
	if c.Python != nil && c.Python.Env != "" {
		in.EnvModules[strings.Trim(c.Python.Env, "/")] = fmt.Sprintf("%s:%d", c.Rel(), c.Line("python.env"))
	}
	if c.TypeScript != nil && c.TypeScript.Env != "" {
		in.EnvModules[strings.Trim(c.TypeScript.Env, "/")] = fmt.Sprintf("%s:%d", c.Rel(), c.Line("typescript.env"))
	}
	declared, ps := surface.Declared(c)
	in.Declared, problems = declared, append(problems, ps...)

	baselinePath := filepath.Join(c.Root, config.Dir, record.BaselineFile)
	lines, ps, err := surface.ReadBaseline(baselinePath, in.BaselineFile)
	if err != nil {
		return nil, nil, err
	}
	if os.Getenv(env.UpdateBaseline) != "" {
		surfaces, accounted := gate.Surfaces(in)
		lines = surface.Rebuild(lines, surfaces, accounted)
		if err := surface.WriteBaseline(baselinePath, lines); err != nil {
			return nil, nil, err
		}
		ps = nil
	}
	in.Baseline, problems = lines, append(problems, ps...)

	problems = append(problems, gate.Check(in)...)
	sort.Strings(problems)
	return t, problems, nil
}

func allFiles(t walk.Tree) []string {
	var out []string
	for _, fs := range t {
		out = append(out, fs...)
	}
	return out
}
