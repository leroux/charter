package chartertest

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// A mutant is a small change to an entry's realizing code: one or more
// edits to one file, each a byte range replaced by a text.

// Mutant is one change mutate runs the citing tests against.
type Mutant struct {
	// Entry is the id of the entry whose realizing code it changes.
	Entry string
	// File is slash-relative to the root; Line is where the change is.
	File string
	Line int
	// Func is the realizing function it changes, as the anchor names it.
	Func string
	// Op names the change: an operator below, or "given" for a mutant a
	// caller wrote.
	Op string
	// Orig is the text replaced and Repl what replaces it, shortened.
	Orig, Repl string
	// Note is a given mutant's note: the clause it aims at, say.
	Note string

	edits []mutantEdit
}

type mutantEdit struct {
	start, end int
	text       string
}

// apply returns src with the mutant's edits made.
func (m Mutant) apply(src []byte) []byte {
	edits := append([]mutantEdit(nil), m.edits...)
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := append([]byte(nil), src...)
	for _, e := range edits {
		out = append(out[:e.start:e.start], append([]byte(e.text), out[e.end:]...)...)
	}
	return out
}

// GivenMutant is a mutant a caller wrote: Old, which must occur exactly once
// in File, becomes New.
type GivenMutant struct {
	Entry string `json:"entry"`
	File  string `json:"file"`
	Old   string `json:"old"`
	New   string `json:"new"`
	Note  string `json:"note,omitempty"`
}

// realizingFunc is a function an anchor realizes, located in its file.
type realizingFunc struct {
	file string // slash-relative
	name string // as the anchor names it
	body *ast.BlockStmt
	fset *token.FileSet
	src  []byte
	from int // the function's first line, its doc comment included
	to   int // its last line
}

// errNotAFunction says an anchor is on a package, a type, a variable or a
// constant: nothing a mutant can change.
var errNotAFunction = errors.New("the anchor is not on a function")

// locateFunc finds the function whose doc comment holds the anchor's
// directive, or errNotAFunction.
func locateFunc(root string, a Anchor) (*realizingFunc, error) {
	if !strings.HasSuffix(a.File, ".go") {
		return nil, errNotAFunction
	}
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(a.File)))
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, a.File, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Doc == nil || fd.Body == nil {
			continue
		}
		if fset.Position(fd.Doc.Pos()).Line <= a.Line && a.Line <= fset.Position(fd.Doc.End()).Line {
			return &realizingFunc{
				file: a.File, name: a.Decl, body: fd.Body, fset: fset, src: src,
				from: fset.Position(fd.Pos()).Line, to: fset.Position(fd.End()).Line,
			}, nil
		}
	}
	return nil, errNotAFunction
}

// lineSet is a set of line ranges in one file; nil holds every line.
type lineSet [][2]int

func (s lineSet) has(line int) bool {
	if s == nil {
		return true
	}
	for _, r := range s {
		if r[0] <= line && line <= r[1] {
			return true
		}
	}
	return false
}

func (s lineSet) overlaps(from, to int) bool {
	if s == nil {
		return true
	}
	for _, r := range s {
		if r[0] <= to && from <= r[1] {
			return true
		}
	}
	return false
}

// swaps are the binary operators a mutant replaces, each with what it
// becomes: a relational one with its boundary and its opposite, == and !=
// with each other, && and || with each other, and arithmetic with its
// inverse.
var swaps = map[token.Token][]token.Token{ //nolint:exhaustive // only the operators a mutant swaps
	token.LSS: {token.LEQ, token.GEQ}, token.LEQ: {token.LSS, token.GTR},
	token.GTR: {token.GEQ, token.LEQ}, token.GEQ: {token.GTR, token.LSS},
	token.EQL: {token.NEQ}, token.NEQ: {token.EQL},
	token.LAND: {token.LOR}, token.LOR: {token.LAND},
	token.ADD: {token.SUB}, token.SUB: {token.ADD},
	token.MUL: {token.QUO}, token.QUO: {token.MUL},
}

// mutantsOf lists the mutants of one realizing function for an entry, on
// the lines the set holds.
//
//charter:realizes CH-261005-VFJR
func mutantsOf(entry string, fn *realizingFunc, lines lineSet) []Mutant {
	var out []Mutant
	off := func(p token.Pos) int { return fn.fset.Position(p).Offset }
	line := func(p token.Pos) int { return fn.fset.Position(p).Line }
	add := func(op string, at token.Pos, edits ...mutantEdit) {
		if !lines.has(line(at)) {
			return
		}
		// Orig and Repl show the span the edits cover, before and after.
		from, to := edits[0].start, edits[len(edits)-1].end
		span := Mutant{edits: append([]mutantEdit(nil), edits...)}
		for i := range span.edits {
			span.edits[i].start -= from
			span.edits[i].end -= from
		}
		orig := fn.src[from:to]
		out = append(out, Mutant{
			Entry: entry, File: fn.file, Line: line(at), Func: fn.name, Op: op,
			Orig: shorten(string(orig)), Repl: shorten(string(span.apply(orig))), edits: edits,
		})
	}
	ast.Inspect(fn.body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			// A closure is the function's code too.
			return true
		case *ast.BinaryExpr:
			for _, t := range swaps[x.Op] {
				p := off(x.OpPos)
				add("operator", x.OpPos, mutantEdit{p, p + len(x.Op.String()), t.String()})
			}
		case *ast.IfStmt:
			c := x.Cond
			add("negate-condition", c.Pos(), mutantEdit{off(c.Pos()), off(c.Pos()), "!("}, mutantEdit{off(c.End()), off(c.End()), ")"})
			if len(x.Body.List) > 0 {
				add("empty-branch", x.Body.Pos(), mutantEdit{off(x.Body.Lbrace) + 1, off(x.Body.Rbrace), "\n_ = 0\n"})
			}
		case *ast.ForStmt:
			if x.Cond != nil {
				c := x.Cond
				add("negate-condition", c.Pos(), mutantEdit{off(c.Pos()), off(c.Pos()), "!("}, mutantEdit{off(c.End()), off(c.End()), ")"})
			}
		case *ast.ExprStmt:
			if _, ok := x.X.(*ast.CallExpr); ok {
				add("drop-call", x.Pos(), mutantEdit{off(x.Pos()), off(x.End()), "_ = 0"})
			}
		case *ast.ReturnStmt:
			for _, r := range x.Results {
				if isErrorExpr(fn.src, off, r) {
					add("nil-error", r.Pos(), mutantEdit{off(r.Pos()), off(r.End()), "nil"})
				}
			}
		case *ast.Ident:
			if x.Name == "true" || x.Name == "false" {
				flip := map[string]string{"true": "false", "false": "true"}[x.Name]
				add("flip-bool", x.Pos(), mutantEdit{off(x.Pos()), off(x.End()), flip})
			}
		case *ast.BasicLit:
			if x.Kind == token.INT {
				if v, err := strconv.ParseInt(x.Value, 0, 64); err == nil {
					add("int-plus-one", x.Pos(), mutantEdit{off(x.Pos()), off(x.End()), strconv.FormatInt(v+1, 10)})
				}
			}
		case *ast.IncDecStmt:
			flip := "--"
			if x.Tok == token.DEC {
				flip = "++"
			}
			p := off(x.TokPos)
			add("operator", x.TokPos, mutantEdit{p, p + 2, flip})
		case *ast.BranchStmt:
			switch {
			case x.Tok == token.BREAK && x.Label == nil:
				add("break-continue", x.Pos(), mutantEdit{off(x.Pos()), off(x.End()), "continue"})
			case x.Tok == token.CONTINUE && x.Label == nil:
				add("break-continue", x.Pos(), mutantEdit{off(x.Pos()), off(x.End()), "break"})
			}
		}
		return true
	})
	return out
}

// isErrorExpr reports a returned expression that is plainly an error: err,
// a name ending in Err, or a call that makes one.
func isErrorExpr(src []byte, off func(token.Pos) int, e ast.Expr) bool {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name == "err" || strings.HasSuffix(x.Name, "Err")
	case *ast.CallExpr:
		fun := string(src[off(x.Fun.Pos()):off(x.Fun.End())])
		return strings.HasPrefix(fun, "fmt.Errorf") || strings.HasPrefix(fun, "errors.") || strings.HasSuffix(fun, "Errorf") || strings.HasSuffix(fun, "Error")
	}
	return false
}

// givenMutant turns a caller's replacement into a mutant, refusing one
// whose text does not occur exactly once in its file.
//
//charter:realizes CH-261005-VFJR
func givenMutant(root string, g GivenMutant) (Mutant, error) {
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(g.File)))
	if err != nil {
		return Mutant{}, fmt.Errorf("given mutant for %s: %w", g.Entry, err)
	}
	if n := bytes.Count(src, []byte(g.Old)); g.Old == "" || n != 1 {
		return Mutant{}, fmt.Errorf("given mutant for %s: its old text occurs %d times in %s, not once", g.Entry, n, g.File)
	}
	at := bytes.Index(src, []byte(g.Old))
	return Mutant{
		Entry: g.Entry, File: g.File, Line: bytes.Count(src[:at], []byte("\n")) + 1, Func: enclosingFunc(g.File, src, at), Op: "given",
		Orig: shorten(g.Old), Repl: shorten(g.New), Note: g.Note,
		edits: []mutantEdit{{at, at + len(g.Old), g.New}},
	}, nil
}

// enclosingFunc names the function holding the offset, as an anchor would
// name it, or "".
func enclosingFunc(name string, src []byte, at int) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, src, 0)
	if err != nil {
		return ""
	}
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fset.Position(fd.Pos()).Offset > at || at >= fset.Position(fd.End()).Offset {
			continue
		}
		if fd.Recv != nil && len(fd.Recv.List) == 1 {
			return recvName(fd.Recv.List[0].Type) + "." + fd.Name.Name
		}
		return fd.Name.Name
	}
	return ""
}

// sample keeps at most limit of the mutants, drawn with the seed, in their
// order; limit 0 keeps all.
func sample(ms []Mutant, limit int, seed int64) []Mutant {
	if limit <= 0 || len(ms) <= limit {
		return ms
	}
	idx := rand.New(rand.NewSource(seed)).Perm(len(ms))[:limit]
	sort.Ints(idx)
	out := make([]Mutant, 0, limit)
	for _, i := range idx {
		out = append(out, ms[i])
	}
	return out
}

func shorten(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 80 {
		return s[:77] + "..."
	}
	return s
}

// parseDiff reads git diff -U0 output: per file, slash-relative, the line
// ranges of the new side that the diff added or changed.
func parseDiff(diff []byte) map[string]lineSet {
	out := map[string]lineSet{}
	var file string
	sc := bufio.NewScanner(bytes.NewReader(diff))
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for sc.Scan() {
		l := sc.Text()
		switch {
		case strings.HasPrefix(l, "+++ "):
			file = ""
			if p := strings.TrimPrefix(l, "+++ "); strings.HasPrefix(p, "b/") {
				file = p[2:]
			}
		case strings.HasPrefix(l, "@@ ") && file != "":
			// @@ -a[,b] +c[,d] @@
			fields := strings.Fields(l)
			if len(fields) < 3 || !strings.HasPrefix(fields[2], "+") {
				continue
			}
			start, count := fields[2][1:], "1"
			if s, c, ok := strings.Cut(start, ","); ok {
				start, count = s, c
			}
			from, err1 := strconv.Atoi(start)
			n, err2 := strconv.Atoi(count)
			if err1 != nil || err2 != nil || n == 0 {
				continue
			}
			out[file] = append(out[file], [2]int{from, from + n - 1})
		}
	}
	return out
}
