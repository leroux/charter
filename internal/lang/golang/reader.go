// Package golang is the charter's Go language: a reader that turns Go source
// into facts, and a generator that writes the claims module's specs.go,
// whose Verifies and Observes are the marker a Go test calls.
package golang

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/leroux/charter/internal/facts"
)

// DefaultReaders are the calls that read the environment when the
// configuration names none.
var DefaultReaders = []string{"os.Getenv", "os.LookupEnv"}

// Reader reads Go files into facts.
type Reader struct {
	// Root is the charter's root; files are slash-relative to it.
	Root string
	// Claims is the import path of the generated claims module, and
	// ClaimsDir its directory, slash-relative to the root.
	Claims, ClaimsDir string
	// EnvDir is the constants package's directory, slash-relative to the
	// root, or "".
	EnvDir string
	// Readers are "importpath.Func" calls that read the environment.
	Readers []string
}

// helpers are the names the generated package declares beside the symbols:
// a selector of one of them is no reference to an entry.
var helpers = map[string]bool{"Verifies": true, "Observes": true, "Spec": true, "Observation": true, "Status": true, "Planned": true, "Active": true, "Deprecated": true}

// envName is the shape of an environment variable's name.
var envName = regexp.MustCompile(`^[A-Z][A-Z0-9]*(?:_[A-Z0-9]+)+$`)

// Read parses each file and returns its facts. Files under the claims
// directory are the generated module's and are not read.
//
//charter:realizes CH-261126-PQ7X
func (r Reader) Read(files []string) (*facts.Set, error) {
	set := &facts.Set{}
	fset := token.NewFileSet()
	envImport := ""
	if r.EnvDir != "" {
		envImport = importPath(r.Root, r.EnvDir)
	}
	readers := map[string]bool{}
	for _, rd := range r.Readers {
		readers[rd] = true
	}
	if len(readers) == 0 {
		for _, rd := range DefaultReaders {
			readers[rd] = true
		}
	}
	for _, rel := range files {
		if r.ClaimsDir != "" && (path.Dir(rel) == r.ClaimsDir || strings.HasPrefix(rel, r.ClaimsDir+"/")) {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(r.Root, filepath.FromSlash(rel)), nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", rel, err)
		}
		fr := fileReader{Reader: r, rel: rel, fset: fset, f: f, set: set, envImport: envImport, readers: readers}
		if strings.HasSuffix(rel, "_test.go") {
			set.TestFiles = append(set.TestFiles, rel)
			fr.citations()
			continue
		}
		fr.anchors()
		if r.EnvDir != "" && path.Dir(rel) == r.EnvDir {
			fr.envConsts()
			continue
		}
		fr.envReads()
	}
	return set, nil
}

type fileReader struct {
	Reader
	rel       string
	fset      *token.FileSet
	f         *ast.File
	set       *facts.Set
	envImport string
	readers   map[string]bool
}

func (fr *fileReader) at(p token.Pos) facts.At {
	return facts.At{File: fr.rel, Line: fr.fset.Position(p).Line}
}

func (fr *fileReader) refuse(p token.Pos, rule, format string, args ...any) {
	fr.set.Refused = append(fr.set.Refused, facts.Refused{At: fr.at(p), Rule: rule, Message: fmt.Sprintf(format, args...)})
}

// directive returns the kind and the rest of an anchor directive, accepting
// a space after the slashes as gofmt leaves it.
func directive(text string) (kind, rest string) {
	if r, ok := strings.CutPrefix(text, "// charter:"); ok {
		text = "//charter:" + r
	}
	for _, k := range []string{"realizes", "concerns"} {
		if r, ok := strings.CutPrefix(text, "//charter:"+k); ok && (r == "" || r[0] == ' ' || r[0] == '\t') {
			return k, r
		}
	}
	return "", ""
}

// anchors reads the directives of a production file: one in the doc comment
// of a package clause or a declaration is an anchor on it, and one anywhere
// else is refused.
//
//charter:realizes CH-021
func (fr *fileReader) anchors() {
	owners := docOwners(fr.f)
	for _, g := range fr.f.Comments {
		for _, c := range g.List {
			kind, rest := directive(c.Text)
			if kind == "" {
				continue
			}
			owner, documents := owners[g]
			if !documents {
				fr.refuse(c.Pos(), "rule 24", "//charter:%s is written in the doc comment of the declaration that carries the decision; this comment documents none", kind)
				continue
			}
			ids := splitList(rest)
			if len(ids) == 0 {
				fr.refuse(c.Pos(), "rule 24", "//charter:%s names at least one spec", kind)
				continue
			}
			fr.set.Anchors = append(fr.set.Anchors, facts.Anchor{At: fr.at(c.Pos()), Kind: kind, IDs: ids, Decl: owner.name, DeclKind: owner.kind})
		}
	}
}

type owner struct{ name, kind string }

// docOwners maps each doc comment group to the declaration it documents.
func docOwners(f *ast.File) map[*ast.CommentGroup]owner {
	owners := map[*ast.CommentGroup]owner{}
	if f.Doc != nil {
		owners[f.Doc] = owner{"package " + f.Name.Name, "package"}
	}
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Doc == nil {
				continue
			}
			name := d.Name.Name
			if d.Recv != nil && len(d.Recv.List) == 1 {
				name = recvName(d.Recv.List[0].Type) + "." + name
			}
			owners[d.Doc] = owner{"func " + name, "func"}
		case *ast.GenDecl:
			kind := map[token.Token]string{token.TYPE: "type", token.VAR: "var", token.CONST: "const"}[d.Tok]
			for _, s := range d.Specs {
				var name string
				var doc *ast.CommentGroup
				switch s := s.(type) {
				case *ast.TypeSpec:
					name, doc = s.Name.Name, s.Doc
				case *ast.ValueSpec:
					name, doc = s.Names[0].Name, s.Doc
				}
				if doc == nil && len(d.Specs) == 1 {
					doc = d.Doc
				}
				if doc != nil && name != "" {
					owners[doc] = owner{kind + " " + name, kind}
				}
			}
		}
	}
	return owners
}

func recvName(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.StarExpr:
		return recvName(x.X)
	case *ast.IndexExpr:
		return recvName(x.X)
	case *ast.IndexListExpr:
		return recvName(x.X)
	case *ast.Ident:
		return x.Name
	}
	return ""
}

// envConsts reads the exported untyped string constants of the constants
// package: each names a variable the tree reads.
func (fr *fileReader) envConsts() {
	for _, d := range fr.f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, s := range gd.Specs {
			vs := s.(*ast.ValueSpec)
			if vs.Type != nil {
				continue
			}
			for i, n := range vs.Names {
				if !ast.IsExported(n.Name) || i >= len(vs.Values) {
					continue
				}
				if v, ok := stringLit(vs.Values[i]); ok {
					fr.set.EnvConsts = append(fr.set.EnvConsts, facts.EnvConst{At: fr.at(n.Pos()), Name: path.Base(fr.EnvDir) + "." + n.Name, Value: v})
				}
			}
		}
	}
}

// envReads reads each call of an environment reader, what it names, and
// each string literal shaped like a variable's name.
//
//charter:realizes CH-034
func (fr *fileReader) envReads() {
	byName, _ := imports(fr.f)
	dynamic := map[int]bool{}
	for _, g := range fr.f.Comments {
		for _, c := range g.List {
			if strings.HasPrefix(c.Text, "//charter:dynamic") || strings.HasPrefix(c.Text, "// charter:dynamic") {
				dynamic[fr.fset.Position(c.Pos()).Line] = true
			}
		}
	}
	ast.Inspect(fr.f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok || len(n.Args) == 0 {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || !fr.readers[byName[id.Name]+"."+sel.Sel.Name] {
				return true
			}
			read := facts.EnvRead{At: fr.at(n.Pos()), Dynamic: dynamic[fr.fset.Position(n.Pos()).Line]}
			switch a := n.Args[0].(type) {
			case *ast.SelectorExpr:
				if x, ok := a.X.(*ast.Ident); ok && fr.envImport != "" && byName[x.Name] == fr.envImport {
					read.Const = path.Base(fr.EnvDir) + "." + a.Sel.Name
				}
			case *ast.BasicLit:
				read.Literal, _ = stringLit(a)
			}
			fr.set.EnvReads = append(fr.set.EnvReads, read)
		case *ast.BasicLit:
			if v, ok := stringLit(n); ok && envName.MatchString(v) {
				fr.set.Literals = append(fr.set.Literals, facts.Literal{At: fr.at(n.Pos()), Value: v})
			}
		}
		return true
	})
}

// citations reads a test file: each Verifies or Observes of the generated
// package that opens a test or a subtest it runs is a citation, one
// anywhere else is refused, and any other use of a symbol is a reference.
// A dot import of the package and an assignment to one of its symbols are
// refused.
//
//charter:realizes CH-023, CH-025
func (fr *fileReader) citations() {
	byName, dots := imports(fr.f)
	for _, d := range dots {
		if d == fr.Claims {
			fr.refuse(fr.f.Pos(), "rule 10", "the claims package must not be dot-imported")
		}
	}
	alias := aliasFor(byName, fr.Claims)
	if alias == "" {
		return
	}
	for _, decl := range fr.f.Decls {
		fn, method := "", false
		leading := map[*ast.CallExpr]bool{}
		if fd, ok := decl.(*ast.FuncDecl); ok {
			fn, method = fd.Name.Name, fd.Recv != nil
			if fd.Body != nil {
				markRun(fd.Body, leading, alias)
			}
		}
		ast.Inspect(decl, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.AssignStmt:
				for _, x := range n.Lhs {
					if name, ok := baseSymbol(x, alias); ok {
						fr.refuse(n.Pos(), "rule 1", "%s is assigned outside its declaration", name)
					}
				}
			case *ast.IncDecStmt:
				if name, ok := baseSymbol(n.X, alias); ok {
					fr.refuse(n.Pos(), "rule 1", "%s is assigned outside its declaration", name)
				}
			case *ast.CallExpr:
				call := markerName(n, alias)
				if call == "" {
					return true
				}
				switch {
				case !strings.HasPrefix(fn, "Test"):
					fr.refuse(n.Pos(), "rule 24", "%s sits outside a test function; it is the first statement of the test that verifies the spec", call)
					return false
				case method:
					fr.refuse(n.Pos(), "rule 24", "%s sits in the method %s, which go test never runs; a test is a function", call, fn)
					return false
				case !isTestName(fn):
					fr.refuse(n.Pos(), "rule 24", "%s sits in %s, which go test never runs; a test is Test followed by a capital, a digit or an underscore", call, fn)
					return false
				case !leading[n]:
					fr.refuse(n.Pos(), "rule 24", "%s is the leading statement of its test, before anything runs, skips or branches; here something precedes or encloses it", call)
					return false
				}
				if len(n.Args) < 2 {
					return false
				}
				sym := ""
				if sel, ok := n.Args[1].(*ast.SelectorExpr); ok {
					if id, ok := sel.X.(*ast.Ident); ok && id.Name == alias && !helpers[sel.Sel.Name] {
						sym = sel.Sel.Name
					}
				}
				if sym == "" {
					fr.refuse(n.Pos(), "rule 24", "%s names a spec by its symbol, such as %s.SW_118", call, alias)
					return false
				}
				c := facts.Citation{At: fr.at(n.Pos()), Symbol: sym, Test: fn, Observes: call == "Observes"}
				for _, a := range n.Args[2:] {
					h, ok := stringLit(a)
					if !ok {
						fr.refuse(n.Pos(), "rule 24", "a scenario handle is a string literal")
						continue
					}
					c.Handles = append(c.Handles, h)
				}
				fr.set.Citations = append(fr.set.Citations, c)
				return false
			case *ast.SelectorExpr:
				if id, ok := n.X.(*ast.Ident); ok && id.Name == alias && !helpers[n.Sel.Name] {
					fr.set.References = append(fr.set.References, facts.Reference{At: fr.at(n.Pos()), Symbol: n.Sel.Name})
				}
			}
			return true
		})
	}
}

// markerName returns "Verifies" or "Observes" for a call of the generated
// package's marker through alias, else "".
func markerName(call *ast.CallExpr, alias string) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Name != alias || (sel.Sel.Name != "Verifies" && sel.Sel.Name != "Observes") {
		return ""
	}
	return sel.Sel.Name
}

// markRun marks the marker calls that open a body go test runs, a test's or
// a subtest's, one after another from its first statement, and then those
// opening each subtest it runs: a t.Run that is a statement of the body, or
// of a range over cases in it.
//
//charter:realizes CH-023
func markRun(body *ast.BlockStmt, leading map[*ast.CallExpr]bool, alias string) {
	for _, st := range body.List {
		es, ok := st.(*ast.ExprStmt)
		if !ok {
			break
		}
		call, ok := es.X.(*ast.CallExpr)
		if !ok || markerName(call, alias) == "" {
			break
		}
		leading[call] = true
	}
	var runs func(stmts []ast.Stmt)
	runs = func(stmts []ast.Stmt) {
		for _, st := range stmts {
			switch st := st.(type) {
			case *ast.ExprStmt:
				if b := subtestBody(st.X); b != nil {
					markRun(b, leading, alias)
				}
			case *ast.RangeStmt:
				runs(st.Body.List)
			}
		}
	}
	runs(body.List)
}

// subtestBody returns the body of the function literal a t.Run call runs,
// or nil for any other expression.
func subtestBody(x ast.Expr) *ast.BlockStmt {
	call, ok := x.(*ast.CallExpr)
	if !ok || len(call.Args) == 0 {
		return nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Run" {
		return nil
	}
	lit, ok := call.Args[len(call.Args)-1].(*ast.FuncLit)
	if !ok {
		return nil
	}
	return lit.Body
}

// isTestName reports whether go test runs a function of the name: Test
// alone, or Test followed by anything but a lower-case letter.
func isTestName(fn string) bool {
	rest, ok := strings.CutPrefix(fn, "Test")
	if !ok {
		return false
	}
	if rest == "" {
		return true
	}
	return !unicode.IsLower([]rune(rest)[0])
}

// baseSymbol unwraps dereferences, indexes and field selections down to
// alias.Symbol.
func baseSymbol(x ast.Expr, alias string) (string, bool) {
	for {
		switch e := x.(type) {
		case *ast.ParenExpr:
			x = e.X
		case *ast.StarExpr:
			x = e.X
		case *ast.IndexExpr:
			x = e.X
		case *ast.SelectorExpr:
			if id, ok := e.X.(*ast.Ident); ok && id.Name == alias {
				return e.Sel.Name, !helpers[e.Sel.Name]
			}
			x = e.X
		default:
			return "", false
		}
	}
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

func aliasFor(byName map[string]string, importPath string) string {
	for name, p := range byName {
		if p == importPath {
			return name
		}
	}
	return ""
}

func stringLit(x ast.Expr) (string, bool) {
	bl, ok := x.(*ast.BasicLit)
	if !ok || bl.Kind != token.STRING {
		return "", false
	}
	v, err := strconv.Unquote(bl.Value)
	return v, err == nil
}

func splitList(v string) []string {
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// importPath is the import path of the package in dir, slash-relative to
// root: the module path of the nearest go.mod above it, and the rest.
func importPath(root, dir string) string {
	for d := dir; ; d = path.Dir(d) {
		if mod := modulePath(filepath.Join(root, filepath.FromSlash(d), "go.mod")); mod != "" {
			rest := strings.TrimPrefix(strings.TrimPrefix(dir, d), "/")
			if rest == "" || d == dir {
				return mod
			}
			return mod + "/" + rest
		}
		if d == "." || d == "/" {
			return ""
		}
	}
}

func modulePath(gomod string) string {
	f, err := os.Open(gomod)
	if err != nil {
		return ""
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(s.Text()), "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}
