package toml

import (
	"fmt"
	"unicode/utf8"
)

// ParseError is why a document was refused, and where: the line and the
// column, both counted from 1, of the first character of the problem.
// Columns count characters, not bytes.
//
//charter:realizes SW-261126-RVMG
type ParseError struct {
	Line   int
	Column int
	Msg    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("toml: line %d, column %d: %s", e.Line, e.Column, e.Msg)
}

// fail stops the parse with a *ParseError at the byte offset at.
//
//charter:realizes SW-261126-RVMG
func (p *parser) fail(at int, format string, args ...any) {
	line, col := 1, 1
	for i := 0; i < at && i < len(p.src); {
		r, n := utf8.DecodeRune(p.src[i:])
		if r == '\n' {
			line, col = line+1, 1
		} else {
			col++
		}
		i += n
	}
	panic(&ParseError{Line: line, Column: col, Msg: fmt.Sprintf(format, args...)})
}
