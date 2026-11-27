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
	"github.com/leroux/charter/internal/docs"
	"github.com/leroux/charter/internal/env"
	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/gate"
	"github.com/leroux/charter/internal/lang/golang"
	"github.com/leroux/charter/internal/lang/plain"
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
	// Docs are the docs charter.toml names, and Code what their code
	// references resolve against.
	Docs []*docs.Doc
	Code docs.Code
	// Problems found while reading, before the rules run.
	Problems []string
}

// DocsKind names the walk's list of the docs.
const DocsKind = "docs"

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
		Files:     filesKinds(c),
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
	for i, f := range c.Files {
		set, err := plain.Reader{Root: c.Root, Comment: f.Comment}.Read(t.Files[filesKind(i)])
		if err != nil {
			return nil, err
		}
		t.Facts.Add(set)
	}
	var langFiles []string
	for _, l := range languages(c) {
		langFiles = append(langFiles, t.Files[l]...)
	}
	t.Code = docs.NewCode(c.Root, langFiles, t.Facts.Decls)
	files, err := DocFiles(c)
	if err != nil {
		return nil, err
	}
	if t.Docs, err = docs.Read(c.Root, files); err != nil {
		return nil, err
	}
	return t, nil
}

// DocFiles lists the docs: every markdown file under the root that
// charter.toml's doc patterns match. The patterns are their own scope:
// the walk's only and ignore patterns do not narrow them.
//
//charter:realizes CH-261127-3NY3
func DocFiles(c *config.Config) ([]string, error) {
	if len(c.Docs) == 0 {
		return nil, nil
	}
	tree, err := walk.Walk(walk.Options{Root: c.Root, Boundary: config.Dir + "/" + config.File, Files: map[string][]string{DocsKind: c.Docs}})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, f := range tree[DocsKind] {
		if strings.HasSuffix(f, ".md") {
			out = append(out, f)
		}
	}
	return out, nil
}

// ReadDecls runs each configured language's reader over the files, found
// under root, which may be another directory than the charter's, as a
// revision's files written out, and returns their declarations.
func ReadDecls(c *config.Config, root string, files []string) ([]facts.Decl, error) {
	by := map[string][]string{}
	for _, f := range files {
		by[walk.Language(f)] = append(by[walk.Language(f)], f)
	}
	set := &facts.Set{}
	if c.Go != nil && len(by[walk.Go]) > 0 {
		s, err := golang.Reader{Root: root, Claims: c.Go.Module, ClaimsDir: config.Dir}.Read(by[walk.Go])
		if err != nil {
			return nil, err
		}
		set.Add(s)
	}
	if c.Python != nil && len(by[walk.Python]) > 0 {
		s, err := python.Reader{Root: root, Dir: c.Python.Dir, ClaimsDir: config.Dir, Table: fmt.Sprintf("%s:%d", c.Rel(), c.Line("python"))}.Read(by[walk.Python])
		if err != nil {
			return nil, err
		}
		set.Add(s)
	}
	if ts := c.TypeScript; ts != nil && len(by[walk.TypeScript]) > 0 {
		compiler := ts.Compiler
		if !filepath.IsAbs(compiler) {
			compiler = filepath.Join(c.Root, compiler)
		}
		s, err := typescript.Reader{Root: root, Compiler: compiler, Import: ts.Import, ClaimsDir: config.Dir, Table: fmt.Sprintf("%s:%d", c.Rel(), c.Line("typescript"))}.Read(by[walk.TypeScript])
		if err != nil {
			return nil, err
		}
		set.Add(s)
	}
	return set.Decls, nil
}

// filesKind names the walk's list of the i-th [[files]] table.
func filesKind(i int) string { return fmt.Sprintf("files.%d", i) }

func filesKinds(c *config.Config) map[string][]string {
	m := map[string][]string{}
	for i, f := range c.Files {
		m[filesKind(i)] = f.Patterns
	}
	return m
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

// Write writes every generated file, then rewrites each doc's references
// and embeds. A tree whose code cannot be read still has its docs' entry
// references and embeds made current.
//
//charter:realizes CH-020, CH-261127-H4B7
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
	if len(c.Docs) == 0 {
		return nil
	}
	t, err := Load(c)
	if err != nil {
		// The code cannot be read: the docs are still read, and their
		// references to entries and their embeds made current.
		files, werr := DocFiles(c)
		if werr != nil {
			return werr
		}
		ds, derr := docs.Read(c.Root, files)
		if derr != nil {
			return derr
		}
		t = &Tree{Docs: ds, Code: docs.NewCode(c.Root, nil, nil)}
	}
	for _, d := range t.Docs {
		out := docs.Rewrite(d, rec, t.Code)
		if out == strings.Join(d.Lines, "\n") {
			continue
		}
		if err := os.WriteFile(filepath.Join(c.Root, filepath.FromSlash(d.File)), []byte(out), 0o644); err != nil {
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
	problems = append(problems, docs.Problems(t.Docs, t.Record, t.Code)...)
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
