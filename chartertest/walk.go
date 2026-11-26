package chartertest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// tree is every Go file under the root, parsed once, keyed by slash path
// relative to the root, and, when Options.TypeScript is set, the slash paths
// of every TypeScript file.
type tree struct {
	root  string
	fset  *token.FileSet
	files map[string]*ast.File
	order []string
	ts    []string
}

//charter:realizes CH-036
func loadTree(r *resolved) (*tree, error) {
	t := &tree{root: r.Root, fset: token.NewFileSet(), files: map[string]*ast.File{}}
	err := filepath.WalkDir(r.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(r.Root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			if r.ignored(rel) || !r.inScope(rel, true) {
				return filepath.SkipDir
			}
			// Hidden and underscore directories and nested modules are skipped, matching what the go tool builds.
			if base := d.Name(); strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_") {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(p, "go.mod")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}
		if r.ignored(rel) || !r.inScope(rel, false) {
			return nil
		}
		if r.TypeScript != nil && isTSFile(rel) {
			t.ts = append(t.ts, rel)
			return nil
		}
		if !strings.HasSuffix(rel, ".go") {
			return nil
		}
		f, err := parser.ParseFile(t.fset, p, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s: %w", rel, err)
		}
		t.files[rel] = f
		t.order = append(t.order, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(t.order)
	sort.Strings(t.ts)
	return t, nil
}

func isTest(rel string) bool { return strings.HasSuffix(rel, "_test.go") }

// dir is the slash directory of a relative file path, "." at the root.
func dir(rel string) string { return path.Dir(rel) }

// pos formats a position as rel:line.
func (t *tree) pos(p token.Pos) string {
	ps := t.fset.Position(p)
	rel, err := filepath.Rel(t.root, ps.Filename)
	if err != nil {
		rel = ps.Filename
	}
	return fmt.Sprintf("%s:%d", filepath.ToSlash(rel), ps.Line)
}

// imports maps each import name in a file to its path: the alias when
// given, else the last path element. Dot imports are returned apart.
func imports(f *ast.File) (byName map[string]string, dot []string) {
	byName = map[string]string{}
	for _, s := range f.Imports {
		p, err := strconv.Unquote(s.Path.Value)
		if err != nil {
			continue
		}
		switch {
		case s.Name == nil:
			byName[path.Base(p)] = p
		case s.Name.Name == ".":
			dot = append(dot, p)
		case s.Name.Name == "_":
		default:
			byName[s.Name.Name] = p
		}
	}
	return byName, dot
}

// aliasFor returns the name a file imports importPath under, or "".
func aliasFor(byName map[string]string, importPath string) string {
	for name, p := range byName {
		if p == importPath {
			return name
		}
	}
	return ""
}

// stringLit returns the value of a string literal expression.
func stringLit(x ast.Expr) (string, bool) {
	bl, ok := x.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", false
	}
	v, err := strconv.Unquote(bl.Value)
	return v, err == nil
}

// directiveLines returns the lines holding a comment that starts with the
// directive.
func (t *tree) directiveLines(f *ast.File, directive string) map[int]bool {
	lines := map[int]bool{}
	for _, g := range f.Comments {
		for _, c := range g.List {
			if strings.HasPrefix(c.Text, directive) {
				lines[t.fset.Position(c.Pos()).Line] = true
			}
		}
	}
	return lines
}
