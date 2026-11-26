package chartertest

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/leroux/charter"
)

// entry is one spec or observation. A spec is read from its area document,
// or from its own file in the earlier form; an observation from its
// declaration in the entries package.
type entry struct {
	// Name is the Go symbol a test names it by: a spec's id with an
	// underscore, SW_118, or an observation's variable name.
	Name string
	// ID is a spec's id, SW-118, or an observation's name.
	ID    string
	Text  string // a spec's title or an observation's text
	Level string // spec or observation
	// Status is planned, active or deprecated for a spec; observed,
	// documented, broken or unsure for an observation.
	Status string

	// Specs only.
	Form   string // area, from an area document; file, from a spec file
	Need   bool   // an entry of needs.md
	Lens   string // the id's prefix: STK, SW, CH, or one of the earlier lenses
	Number int    // a numbered id's number; 0 for a dated id
	Day    string // a dated id's day, YYMMDD; "" for a numbered id
	Code   string // a dated id's code
	Kind   string // invariant, bound, structure, or ""
	// Serves and See are an area entry's relations, by id; Realizes,
	// Concerns and Related are a file entry's, and Links its link targets.
	Serves, See                 []string
	Links                       map[string]string
	Realizes, Concerns, Related []string
	Accounts                    []string // surface keys
	// Description is the decision's text; Rationale its Why; Guidance its
	// advisory how; Verification a file entry's prose in place of scenarios.
	Description, Rationale, Guidance string
	Verification                     string
	Scenarios                        []Scenario
	Changes                          []Change
	fields                           map[string]bool // field lines seen
	// textLines and whyLines are an area entry's text and Why as written,
	// each line with its number in the document, the Why's label cut.
	textLines, whyLines []bodyLine

	// Observations only.
	Subject, Source, Observed string
	Supports                  []string // spec ids

	File string // slash-relative to the root
	Pos  string // File:line
}

// isNeed reports whether the entry is a need: an entry of needs.md, or a
// stakeholder spec in the earlier form. An area entry is a need by its
// document alone, whatever its id.
//
//charter:realizes CH-011
func (e *entry) isNeed() bool { return e.Need || (e.Form == formFile && e.Lens == "STK") }

// hasScenario reports whether the entry lists the handle.
func (e *entry) hasScenario(handle string) bool {
	for _, s := range e.Scenarios {
		if s.Handle == handle {
			return true
		}
	}
	return false
}

const (
	levelObservation = "observation"
	levelSpec        = "spec"
)

var observationStatuses = map[string]string{"Observed": "observed", "Documented": "documented", "Broken": "broken", "ObservationUnsure": "unsure"}

// charterPath is the import path of the charter package, taken from its type
// so the gate recognises it wherever the library is imported from.
var charterPath = reflect.TypeOf(charter.Spec{}).PkgPath()

// entriesPackage is the registry: the specs read from their files and the
// observations read from the entries package, keyed by symbol.
type entriesPackage struct {
	name    string // the entries package's package clause
	entries map[string]*entry
	order   []string // specs by id, then observations as declared
	byID    map[string]*entry
}

func (p *entriesPackage) add(e *entry) {
	p.entries[e.Name] = e
	p.byID[e.ID] = e
	p.order = append(p.order, e.Name)
}

// parseObservations reads every hand-written, non-test file of the entries
// package; only observations may be declared there (rule 1). The file
// generated from the spec files is compared, not parsed.
//
//charter:realizes CH-025
func parseObservations(t *tree, r *resolved, pkg *entriesPackage) []string {
	var problems []string
	for _, rel := range t.order {
		if dir(rel) != r.entriesDir || isTest(rel) {
			continue
		}
		f := t.files[rel]
		pkg.name = f.Name.Name
		if isGeneratedSpecFile(f) {
			continue
		}
		byName, _ := imports(f)
		charterAlias := aliasFor(byName, charterPath)
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				problems = append(problems, fmt.Sprintf("%s: rule 1: functions are not allowed in the entries package", t.pos(d.Pos())))
			case *ast.GenDecl:
				if d.Tok == token.IMPORT {
					continue
				}
				if d.Tok != token.VAR {
					problems = append(problems, fmt.Sprintf("%s: rule 1: only observations and imports are allowed in the entries package", t.pos(d.Pos())))
					continue
				}
				for _, s := range d.Specs {
					e, ps := parseObservation(t, charterAlias, s.(*ast.ValueSpec))
					problems = append(problems, ps...)
					if e == nil {
						continue
					}
					e.File = rel
					if other := pkg.entries[e.Name]; other != nil {
						problems = append(problems, fmt.Sprintf("%s: rule 1: %s is declared twice", e.Pos, e.Name))
						continue
					}
					pkg.add(e)
				}
			}
		}
	}
	return problems
}

// parseObservation reads one var declaration as an observation.
func parseObservation(t *tree, charterAlias string, vs *ast.ValueSpec) (*entry, []string) {
	at := t.pos(vs.Pos())
	if len(vs.Names) != 1 || len(vs.Values) != 1 {
		return nil, []string{fmt.Sprintf("%s: rule 1: a declaration names one observation with one value", at)}
	}
	name := vs.Names[0].Name
	shape := fmt.Sprintf("%s: rule 1: %s is not &charter.Observation{}; a spec is written in a spec file", at, name)
	un, ok := vs.Values[0].(*ast.UnaryExpr)
	if !ok || un.Op != token.AND {
		return nil, []string{shape}
	}
	lit, ok := un.X.(*ast.CompositeLit)
	if !ok {
		return nil, []string{shape}
	}
	if typ, ok := charterSelector(lit.Type, charterAlias); !ok || typ != "Observation" {
		return nil, []string{shape}
	}
	e := &entry{Name: name, Level: levelObservation, Pos: at}
	var problems []string
	bad := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s: rule 1: %s", at, fmt.Sprintf(format, args...)))
	}
	for _, el := range lit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			bad("%s: every field must be keyed", name)
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			bad("%s: every field must be keyed", name)
			continue
		}
		switch key.Name {
		case "ID", "Text", "Subject", "Source", "Observed":
			v, ok := stringLit(kv.Value)
			if !ok {
				bad("%s.%s must be a string literal", name, key.Name)
				continue
			}
			switch key.Name {
			case "ID":
				e.ID = v
			case "Text":
				e.Text = v
			case "Subject":
				e.Subject = v
			case "Source":
				e.Source = v
			default:
				e.Observed = v
			}
		case "Status":
			s, ok := charterSelector(kv.Value, charterAlias)
			if !ok || observationStatuses[s] == "" {
				bad("%s.Status must be charter.Observed, charter.Documented, charter.Broken or charter.ObservationUnsure", name)
				continue
			}
			e.Status = observationStatuses[s]
		case "Supports":
			cl, ok := kv.Value.(*ast.CompositeLit)
			if !ok {
				bad("%s.Supports must be a slice literal of spec ids", name)
				continue
			}
			for _, x := range cl.Elts {
				v, ok := stringLit(x)
				if !ok {
					bad("%s.Supports must list spec ids as string literals", name)
					continue
				}
				e.Supports = append(e.Supports, v)
			}
		default:
			bad("%s.%s is not a field of an observation", name, key.Name)
		}
	}
	if e.Status == "" {
		bad("%s has no Status", name)
	}
	if e.ID != name {
		problems = append(problems, fmt.Sprintf("%s: rule 2: %s has ID %q; the ID must equal the variable name", at, name, e.ID))
	}
	if e.Text == "" {
		problems = append(problems, fmt.Sprintf("%s: rule 2: %s has no Text", at, name))
	}
	return e, problems
}

// charterSelector returns Name for an expression charter.Name, where
// charter is the alias the file imports the charter package under.
func charterSelector(x ast.Expr, charterAlias string) (string, bool) {
	sel, ok := x.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Name != charterAlias {
		return "", false
	}
	return sel.Sel.Name, true
}

// envConstRef resolves alias.Name to the value of an env package constant.
func envConstRef(x ast.Expr, byName map[string]string, envConsts map[string]map[string]string) (string, bool) {
	sel, ok := x.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	consts, ok := envConsts[byName[id.Name]]
	if !ok {
		return "", false
	}
	v, ok := consts[sel.Sel.Name]
	return v, ok
}

// checkMutations applies rule 1 to the entries package's own test files:
// an entry may not be assigned there either.
//
//charter:realizes CH-025
func checkMutations(t *tree, r *resolved, pkg *entriesPackage) []string {
	var problems []string
	for _, rel := range t.order {
		if dir(rel) != r.entriesDir || !isTest(rel) {
			continue
		}
		ast.Inspect(t.files[rel], func(n ast.Node) bool {
			for _, name := range mutatedEntries(n, "", pkg) {
				problems = append(problems, fmt.Sprintf("%s: rule 1: %s is assigned outside its declaration", t.pos(n.Pos()), name))
			}
			return true
		})
	}
	return problems
}
