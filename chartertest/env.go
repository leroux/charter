package chartertest

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strconv"
)

// envConstants returns the exported, untyped string constants declared in
// the package at dirRel: name -> value. A constant with a declared type is
// not a variable's name.
func envConstants(t *tree, dirRel string) map[string]string {
	out := map[string]string{}
	for _, rel := range t.order {
		if dir(rel) != dirRel || isTest(rel) {
			continue
		}
		for _, d := range t.files[rel].Decls {
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
						out[n.Name] = v
					}
				}
			}
		}
	}
	return out
}

// checkEnvReads applies rule 8 to every non-test file outside the env
// packages: a call to an env reader names an env constant, no literal
// duplicates a declared name, and no literal matches EnvLiterals.
//
//charter:realizes CH-034
func checkEnvReads(t *tree, r *resolved, envConsts map[string]map[string]string) []string {
	names := map[string]string{} // value -> pkg.Const
	for p, consts := range envConsts {
		for n, v := range consts {
			names[v] = path.Base(p) + "." + n
		}
	}
	envDirs := map[string]bool{}
	for _, d := range r.envDirs {
		envDirs[d] = true
	}
	var problems []string
	for _, rel := range t.order {
		if isTest(rel) || envDirs[dir(rel)] {
			continue
		}
		f := t.files[rel]
		byName, _ := imports(f)
		dynamic := t.directiveLines(f, "//charter:dynamic")
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				sel, ok := n.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				id, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				callee := byName[id.Name] + "." + sel.Sel.Name
				if !r.readers[callee] || len(n.Args) == 0 {
					return true
				}
				if _, ok := envConstRef(n.Args[0], byName, envConsts); ok {
					return true
				}
				if dynamic[t.fset.Position(n.Pos()).Line] {
					return true
				}
				problems = append(problems, fmt.Sprintf("%s: rule 8: %s must name a constant of an env package; a read named by configuration carries //charter:dynamic", t.pos(n.Pos()), callee))
			case *ast.BasicLit:
				if n.Kind != token.STRING {
					return true
				}
				v, err := strconv.Unquote(n.Value)
				if err != nil {
					return true
				}
				if c, ok := names[v]; ok {
					problems = append(problems, fmt.Sprintf("%s: rule 8: literal %q duplicates %s; use the constant", t.pos(n.Pos()), v, c))
					return true
				}
				for _, re := range r.literals {
					if re.MatchString(v) {
						problems = append(problems, fmt.Sprintf("%s: rule 8: literal %q looks like an environment variable; declare it in an env package", t.pos(n.Pos()), v))
						break
					}
				}
			}
			return true
		})
	}
	return problems
}
