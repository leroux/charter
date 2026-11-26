package chartertest

import (
	"fmt"
	"go/ast"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// A test verifies a spec by calling Verifies as a statement of a test
// function, naming the spec by its symbol and the scenarios it exercises by
// their handles:
//
//	chartertest.Verifies(t, claims.SW_118, "burst-caps-tokens")
//
// The gate reads the call statically: it is a citation of the spec and of
// each handle, and it counts toward coverage. A Verifies anywhere but inside
// a test function, a spec named by anything but its symbol, or a handle
// that is not a string literal is a rule 24 problem. Any other reference to
// an entry from a test file is a reference, not a citation: it is what a
// probe's Observes makes, and what rules 1 and 12 look at.

// Ref is one place a test cites or references an entry: the file,
// slash-relative to the root, the enclosing test function or title, or ""
// outside one, the line, and, for a citation, the scenario handles named.
// A Go citation opening a subtest names it in Subtest, as go test does below
// the test: the t.Run titles from the test down, spaces as underscores,
// split at each slash, and "" for a title that is no string literal.
type Ref struct {
	File, Func string
	Line       int
	Handles    []string
	Subtest    []string
}

// key identifies a reference site: one per file, function and line.
func (r Ref) key() string { return fmt.Sprintf("%s\x00%s\x00%d", r.File, r.Func, r.Line) }

// chartertestPath is this package's import path, which a test file imports
// to call Verifies.
var chartertestPath = reflect.TypeOf(Options{}).PkgPath()

// collectRefs walks every test file in the tree for Verifies calls, the
// citations, and for selector expressions that resolve to the entries
// package, the references. It also applies rule 10 and the assignment half
// of rule 1 to those files.
//
//charter:realizes CH-023, CH-025
func collectRefs(t *tree, r *resolved, pkg *entriesPackage) (refs, cites map[string][]Ref, problems []string) {
	refs, cites = map[string][]Ref{}, map[string][]Ref{}
	seen := map[*map[string][]Ref]map[string]bool{&refs: {}, &cites: {}}
	add := func(m *map[string][]Ref, name string, ref Ref) {
		k := name + "\x00" + ref.key()
		if seen[m][k] {
			return
		}
		seen[m][k] = true
		(*m)[name] = append((*m)[name], ref)
	}
	for _, rel := range t.order {
		if !isTest(rel) {
			continue
		}
		f := t.files[rel]
		byName, dots := imports(f)
		for _, d := range dots {
			if d == r.Entries {
				problems = append(problems, fmt.Sprintf("%s: rule 10: the entries package must not be dot-imported", t.pos(f.Pos())))
			}
		}
		alias := aliasFor(byName, r.Entries)
		ctAlias := aliasFor(byName, chartertestPath)
		if alias == "" && ctAlias == "" {
			continue
		}
		for _, decl := range f.Decls {
			fn := ""
			method := false
			// leading are the Verifies calls that open the test's body, or a
			// subtest's, one after another from its first statement: the only
			// place a citation counts, before anything can skip, branch or
			// enclose it. Each maps to the subtest it opens, nil for the test.
			leading := map[*ast.CallExpr][]string{}
			if fd, ok := decl.(*ast.FuncDecl); ok {
				fn, method = fd.Name.Name, fd.Recv != nil
				if fd.Body != nil {
					markRun(fd.Body, nil, leading, ctAlias, f.Name.Name)
				}
			}
			ast.Inspect(decl, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.AssignStmt, *ast.IncDecStmt:
					for _, name := range mutatedEntries(n, alias, pkg) {
						problems = append(problems, fmt.Sprintf("%s: rule 1: %s is assigned outside its declaration", t.pos(n.Pos()), name))
					}
				case *ast.CallExpr:
					if !isVerifiesCall(n, ctAlias, f.Name.Name) {
						return true
					}
					at := t.pos(n.Pos())
					line := t.fset.Position(n.Pos()).Line
					switch {
					case !strings.HasPrefix(fn, "Test"):
						problems = append(problems, fmt.Sprintf("%s: rule 24: Verifies sits outside a test function; it is the first statement of the test that verifies the spec", at))
						return true
					case method:
						problems = append(problems, fmt.Sprintf("%s: rule 24: Verifies sits in the method %s, which go test never runs; a test is a function", at, fn))
						return true
					case !isTestName(fn):
						problems = append(problems, fmt.Sprintf("%s: rule 24: Verifies sits in %s, which go test never runs; a test is Test followed by a capital, a digit or an underscore", at, fn))
						return true
					case !isLeading(leading, n):
						problems = append(problems, fmt.Sprintf("%s: rule 24: Verifies is the leading statement of its test, before anything runs, skips or branches; here something precedes or encloses it", at))
						return true
					}
					if len(n.Args) < 2 {
						return true
					}
					specSel, ok := n.Args[1].(*ast.SelectorExpr)
					var name string
					if ok {
						if id, ok := specSel.X.(*ast.Ident); ok && id.Name == alias {
							name = specSel.Sel.Name
						}
					}
					e := pkg.entries[name]
					if e == nil || e.Level != levelSpec {
						problems = append(problems, fmt.Sprintf("%s: rule 24: Verifies names a spec by its symbol, such as claims.SW_118", at))
						return true
					}
					ref := Ref{File: rel, Func: fn, Line: line, Subtest: leading[n]}
					for _, a := range n.Args[2:] {
						h, ok := stringLit(a)
						if !ok {
							problems = append(problems, fmt.Sprintf("%s: rule 24: a scenario handle is a string literal", at))
							continue
						}
						ref.Handles = append(ref.Handles, h)
					}
					add(&cites, name, ref)
				case *ast.SelectorExpr:
					id, ok := n.X.(*ast.Ident)
					if !ok || alias == "" || id.Name != alias {
						return true
					}
					if _, ok := pkg.entries[n.Sel.Name]; !ok {
						return true
					}
					add(&refs, n.Sel.Name, Ref{File: rel, Func: fn, Line: t.fset.Position(n.Pos()).Line})
				}
				return true
			})
		}
	}
	return refs, cites, problems
}

// markLeading records the Verifies calls that open a body, one after
// another from its first statement.
//
//charter:realizes CH-023
func markLeading(body *ast.BlockStmt, sub []string, leading map[*ast.CallExpr][]string, ctAlias, pkgName string) {
	for _, st := range body.List {
		es, ok := st.(*ast.ExprStmt)
		if !ok {
			return
		}
		call, ok := es.X.(*ast.CallExpr)
		if !ok || !isVerifiesCall(call, ctAlias, pkgName) {
			return
		}
		leading[call] = sub
	}
}

// markRun marks the Verifies calls that open a body go test runs, a test's
// or a subtest's, and then those opening each subtest it runs: a t.Run
// that is a statement of the body, or of a range over cases in it. A t.Run
// a closure, a branch or any other block encloses may never run, and
// opens no subtest a citation counts in.
//
//charter:realizes CH-023
func markRun(body *ast.BlockStmt, sub []string, leading map[*ast.CallExpr][]string, ctAlias, pkgName string) {
	markLeading(body, sub, leading, ctAlias, pkgName)
	var runs func(stmts []ast.Stmt)
	runs = func(stmts []ast.Stmt) {
		for _, st := range stmts {
			switch st := st.(type) {
			case *ast.ExprStmt:
				if title, body := subtestBody(st.X); body != nil {
					markRun(body, append(slices.Clip(sub), title...), leading, ctAlias, pkgName)
				}
			case *ast.RangeStmt:
				runs(st.Body.List)
			}
		}
	}
	runs(body.List)
}

// subtestBody returns the body of the function literal a t.Run call runs
// and the subtest's name below its parent, as go test gives it, in
// segments: [""] when the title is no string literal. The body is nil for
// any other expression.
func subtestBody(x ast.Expr) ([]string, *ast.BlockStmt) {
	call, ok := x.(*ast.CallExpr)
	if !ok || len(call.Args) == 0 {
		return nil, nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Run" {
		return nil, nil
	}
	lit, ok := call.Args[len(call.Args)-1].(*ast.FuncLit)
	if !ok {
		return nil, nil
	}
	title, ok := stringLit(call.Args[0])
	if !ok {
		return []string{""}, lit.Body
	}
	return strings.Split(subtestName(title), "/"), lit.Body
}

// subtestName is the name go test gives a subtest titled so: each space an
// underscore and each unprintable rune escaped, as testing's rewrite does.
func subtestName(title string) string {
	var b strings.Builder
	for _, r := range title {
		switch {
		case unicode.IsSpace(r):
			b.WriteByte('_')
		case !strconv.IsPrint(r):
			q := strconv.QuoteRune(r)
			b.WriteString(q[1 : len(q)-1])
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// isLeading reports whether markRun marked the call.
func isLeading(leading map[*ast.CallExpr][]string, call *ast.CallExpr) bool {
	_, ok := leading[call]
	return ok
}

// isVerifiesCall reports whether a call is chartertest.Verifies through the
// alias the file imports it by, or a bare Verifies in this package's own
// tests.
//
//charter:realizes CH-023
func isVerifiesCall(call *ast.CallExpr, ctAlias, pkgName string) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		id, ok := fun.X.(*ast.Ident)
		return ok && fun.Sel.Name == "Verifies" && ctAlias != "" && id.Name == ctAlias
	case *ast.Ident:
		return fun.Name == "Verifies" && pkgName == "chartertest"
	}
	return false
}

// isTestName reports whether go test runs a function of the name: Test
// alone, or Test followed by anything but a lower-case letter.
//
//charter:realizes CH-023
func isTestName(fn string) bool {
	rest, ok := strings.CutPrefix(fn, "Test")
	if !ok {
		return false
	}
	if rest == "" {
		return true
	}
	r := []rune(rest)[0]
	return !unicode.IsLower(r)
}

// mutatedEntries returns the entries an assignment or an increment
// writes to, through any field, index or dereference. alias is the name the
// file imports the entries package under, or "" inside the package itself.
func mutatedEntries(n ast.Node, alias string, pkg *entriesPackage) []string {
	var targets []ast.Expr
	switch n := n.(type) {
	case *ast.AssignStmt:
		targets = n.Lhs
	case *ast.IncDecStmt:
		targets = []ast.Expr{n.X}
	}
	var names []string
	for _, x := range targets {
		if name, ok := baseEntry(x, alias, pkg); ok {
			names = append(names, name)
		}
	}
	return names
}

// baseEntry unwraps dereferences, indexes and field selections down to
// alias.Entry, or to a bare entry identifier when alias is "".
func baseEntry(x ast.Expr, alias string, pkg *entriesPackage) (string, bool) {
	for {
		switch e := x.(type) {
		case *ast.ParenExpr:
			x = e.X
		case *ast.StarExpr:
			x = e.X
		case *ast.IndexExpr:
			x = e.X
		case *ast.Ident:
			if alias != "" {
				return "", false
			}
			_, isEntry := pkg.entries[e.Name]
			return e.Name, isEntry
		case *ast.SelectorExpr:
			if id, ok := e.X.(*ast.Ident); ok && alias != "" && id.Name == alias {
				_, isEntry := pkg.entries[e.Sel.Name]
				return e.Sel.Name, isEntry
			}
			x = e.X
		default:
			return "", false
		}
	}
}
