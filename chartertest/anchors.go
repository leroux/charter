package chartertest

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Production code names the specs it answers to with directives in a
// comment, on a package clause's doc or a declaration's, by id:
//
//	//charter:realizes SW-118, SW-119
//	//charter:concerns CON-001
//
// A space after the slashes is accepted too, as gofmt leaves it. TypeScript
// production files carry the same directives in line comments. realizes
// marks code that implements a spec; concerns marks code that must respect
// a spec without implementing it, the usual case for an invariant or a
// structural decision. A realizes anchor is half of a spec's coverage;
// concerns never counts. Rule 24 holds each anchor to a spec that exists
// and is not planned.

const (
	realizesDirective = "//charter:realizes"
	concernsDirective = "//charter:concerns"
)

// directiveText returns a comment's text with "// charter:" normalised to
// "//charter:", so both spellings match the directives.
//
//charter:realizes CH-021
func directiveText(c *ast.Comment) string {
	if rest, ok := strings.CutPrefix(c.Text, "// charter:"); ok {
		return "//charter:" + rest
	}
	return c.Text
}

// Anchor is one place production code names a spec.
type Anchor struct {
	// Kind is realizes or concerns.
	Kind string
	// File is slash-relative to the root; Line is the directive's line.
	File string
	Line int
	// Decl is what the directive documents: "package name", a function or
	// method name, a type, var or const name, or "" elsewhere.
	Decl string
}

// collectAnchors reads every anchor in the tree's non-test Go files outside
// the entries package and applies rule 24.
//
//charter:realizes CH-021
func collectAnchors(t *tree, r *resolved, pkg *entriesPackage) (map[string][]Anchor, []string, error) {
	anchors := map[string][]Anchor{}
	var problems []string
	for _, rel := range t.order {
		if isTest(rel) || dir(rel) == r.entriesDir {
			continue
		}
		f := t.files[rel]
		decls := docOwners(f)
		for _, g := range f.Comments {
			for _, c := range g.List {
				kind := ""
				rest := ""
				for _, d := range []string{realizesDirective, concernsDirective} {
					if r, ok := strings.CutPrefix(directiveText(c), d); ok && (r == "" || r[0] == ' ' || r[0] == '\t') {
						kind, rest = strings.TrimPrefix(d, "//charter:"), r
					}
				}
				if kind == "" {
					continue
				}
				at := t.pos(c.Pos())
				line := t.fset.Position(c.Pos()).Line
				if _, documents := decls[g]; !documents {
					problems = append(problems, fmt.Sprintf("%s: rule 24: //charter:%s is written in the doc comment of the declaration that carries the decision; this comment documents none", at, kind))
					continue
				}
				names := splitList(rest)
				if len(names) == 0 {
					problems = append(problems, fmt.Sprintf("%s: rule 24: //charter:%s names at least one spec", at, kind))
					continue
				}
				for _, n := range names {
					e, ps := anchorTarget(pkg, at, kind, n)
					problems = append(problems, ps...)
					if e == nil {
						continue
					}
					anchors[e.Name] = append(anchors[e.Name], Anchor{Kind: kind, File: rel, Line: line, Decl: decls[g]})
				}
			}
		}
	}
	tsAnchors, ps, err := collectTSAnchors(t, r, pkg)
	problems = append(problems, ps...)
	for n, as := range tsAnchors {
		anchors[n] = append(anchors[n], as...)
	}
	return anchors, problems, err
}

// anchorTarget resolves one id a directive names and applies rule 24.
//
//charter:realizes CH-022
func anchorTarget(pkg *entriesPackage, at, kind, id string) (*entry, []string) {
	e := pkg.byID[id]
	switch {
	case e == nil || e.Level != levelSpec:
		return nil, []string{fmt.Sprintf("%s: rule 24: //charter:%s names %s, which is not a spec id", at, kind, id)}
	case e.Status == "planned":
		return nil, []string{fmt.Sprintf("%s: rule 24: //charter:%s names planned spec %s; set it active as its work begins", at, kind, id)}
	}
	return e, nil
}

// tsDirective finds an anchor directive in a TypeScript line comment.
var tsDirective = regexp.MustCompile(`^\s*//\s*charter:(realizes|concerns)\b(.*)$`)

// collectTSAnchors reads the directives in the tree's TypeScript files that
// are not tests. Decl is empty: the directive's line says where it is.
//
//charter:realizes CH-021
func collectTSAnchors(t *tree, r *resolved, pkg *entriesPackage) (map[string][]Anchor, []string, error) {
	anchors := map[string][]Anchor{}
	var problems []string
	for _, rel := range t.ts {
		if isTSTest(rel) || (r.TypeScript != nil && rel == r.TypeScript.Module) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(t.root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, nil, err
		}
		for i, l := range strings.Split(string(data), "\n") {
			m := tsDirective.FindStringSubmatch(l)
			if m == nil {
				continue
			}
			at := fmt.Sprintf("%s:%d", rel, i+1)
			names := splitList(m[2])
			if len(names) == 0 {
				problems = append(problems, fmt.Sprintf("%s: rule 24: //charter:%s names at least one spec", at, m[1]))
			}
			for _, n := range names {
				e, ps := anchorTarget(pkg, at, m[1], n)
				problems = append(problems, ps...)
				if e != nil {
					anchors[e.Name] = append(anchors[e.Name], Anchor{Kind: m[1], File: rel, Line: i + 1})
				}
			}
		}
	}
	return anchors, problems, nil
}

// docOwners maps each doc comment group to the declaration it documents.
//
//charter:realizes CH-021
func docOwners(f *ast.File) map[*ast.CommentGroup]string {
	owners := map[*ast.CommentGroup]string{}
	if f.Doc != nil {
		owners[f.Doc] = "package " + f.Name.Name
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
			owners[d.Doc] = name
		case *ast.GenDecl:
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
					owners[doc] = name
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
