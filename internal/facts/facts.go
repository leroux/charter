// Package facts is what a language's reader finds in source files and the
// gate's rules read: anchors, citations, references, environment reads and
// constants, literals that look like a variable's name, and the directives
// and calls a reader refused. A fact names no language; each carries the
// slash path, relative to the root, and the line it was found at.
package facts

import "fmt"

// At is where a fact was found.
type At struct {
	File string // slash-relative to the root
	Line int
}

// Pos is file:line.
func (a At) Pos() string { return fmt.Sprintf("%s:%d", a.File, a.Line) }

// Anchor is a realizes or concerns directive on a declaration of
// production code.
type Anchor struct {
	At
	// Kind is "realizes" or "concerns".
	Kind string
	IDs  []string
	// Decl names the declaration the directive sits on, such as "func
	// (*Limiter).Allow", "type Clock" or "package limiter".
	Decl string
	// DeclKind is func, type, var, const, class or package.
	DeclKind string
}

// Citation is a verifies call opening a test, or an observes call opening a
// probe: the symbol it names and the scenarios, by handle.
type Citation struct {
	At
	// Symbol is the generated symbol the call names: an entry's id with
	// underscores, or an observation's name.
	Symbol  string
	Handles []string
	// Test names the test the call opens.
	Test string
	// Observes is true for an observes call.
	Observes bool
}

// Reference is any other use of a generated symbol from a test file.
type Reference struct {
	At
	Symbol string
}

// EnvConst is a string constant of a language's constants module: an
// environment variable the tree reads.
type EnvConst struct {
	At
	// Name is the constant's name, Value the variable's.
	Name, Value string
}

// EnvRead is a read of the environment outside the constants modules and
// the tests. Const names the constant it reads through; Literal holds the
// name when it is a string literal instead; neither is a read through a
// variable. Dynamic is true when the line carries the dynamic directive.
type EnvRead struct {
	At
	Const   string
	Literal string
	Dynamic bool
}

// Literal is a string literal, outside the constants modules and the tests,
// that has the shape of an environment variable's name: capital letters,
// digits and underscores.
type Literal struct {
	At
	Value string
}

// Refused is a directive or a call a reader would not read as an anchor or
// a citation, with the rule it breaks and what to do.
type Refused struct {
	At
	// Rule is the gate rule the problem is reported under, such as "rule 24".
	Rule    string
	Message string
}

// Set is every fact of one tree.
type Set struct {
	Anchors    []Anchor
	Citations  []Citation
	References []Reference
	EnvConsts  []EnvConst
	EnvReads   []EnvRead
	Literals   []Literal
	Refused    []Refused
	// TestFiles are the test files read, slash-relative to the root.
	TestFiles []string
}

// Add appends every fact of o to s.
func (s *Set) Add(o *Set) {
	s.Anchors = append(s.Anchors, o.Anchors...)
	s.Citations = append(s.Citations, o.Citations...)
	s.References = append(s.References, o.References...)
	s.EnvConsts = append(s.EnvConsts, o.EnvConsts...)
	s.EnvReads = append(s.EnvReads, o.EnvReads...)
	s.Literals = append(s.Literals, o.Literals...)
	s.Refused = append(s.Refused, o.Refused...)
	s.TestFiles = append(s.TestFiles, o.TestFiles...)
}
