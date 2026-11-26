// Package toml decodes TOML documents.
//
// It implements a subset of TOML 1.0 (https://toml.io/en/v1.0.0):
// comments, bare, quoted and dotted keys, single-line basic and literal
// strings, decimal integers, booleans and standard tables. A document that
// uses any other part of TOML is refused with an error that names it.
package toml

// Decode parses a TOML document and returns its root table. Tables are
// map[string]any; strings are string, integers int64 and booleans bool.
// A document that is not valid TOML, or that uses a part of TOML this
// package does not implement, is refused with a *ParseError and no table.
func Decode(data []byte) (doc map[string]any, err error) {
	p := &parser{src: data}
	defer func() {
		if r := recover(); r != nil {
			pe, ok := r.(*ParseError)
			if !ok {
				panic(r)
			}
			doc, err = nil, pe
		}
	}()
	p.parse()
	return p.root.export(), nil
}

// parser holds the state of one Decode: the source, the read position,
// the root table and the table key/value pairs currently go into.
type parser struct {
	src  []byte
	pos  int
	root *table
	cur  *table
}

// parse reads the whole document, one line at a time: blank lines and
// comments, table headers, and key/value pairs.
func (p *parser) parse() {
	p.checkEncoding()
	p.root = &table{origin: defined}
	p.cur = p.root
	for {
		p.skipSpace()
		if p.eof() {
			return
		}
		switch c := p.peek(); {
		case c == '#' || c == '\n' || c == '\r':
			p.lineEnd()
		case c == '[':
			p.header()
			p.lineEnd()
		default:
			p.keyValue()
			p.lineEnd()
		}
	}
}

// keyValue reads one key/value pair and stores it in the current table.
func (p *parser) keyValue() {
	at := p.pos
	path := p.key()
	p.skipSpace()
	if p.eof() || p.peek() != '=' {
		p.fail(p.pos, "expected '=' after the key")
	}
	p.pos++
	p.skipSpace()
	v := p.value()
	p.assign(path, v, at)
}
