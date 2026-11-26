package toml

import "strings"

// origin records how a table came to exist, which decides what may later
// add to it.
type origin int

const (
	// implicit: created as the parent of a header's table, as [a] is by
	// [a.b]. A header may still define it, once.
	implicit origin = iota
	// defined: defined by its own header, or the root table.
	defined
	// dotted: created by a dotted key, as a is by a.b = 1.
	dotted
)

// table is a table while it is parsed. Its entries are *table or values.
type table struct {
	origin  origin
	entries map[string]any
}

func (t *table) get(k string) (any, bool) {
	v, ok := t.entries[k]
	return v, ok
}

func (t *table) set(k string, v any) {
	if t.entries == nil {
		t.entries = map[string]any{}
	}
	t.entries[k] = v
}

// export turns the parsed table into what Decode returns.
func (t *table) export() map[string]any {
	out := make(map[string]any, len(t.entries))
	for k, v := range t.entries {
		if sub, ok := v.(*table); ok {
			out[k] = sub.export()
		} else {
			out[k] = v
		}
	}
	return out
}

// header reads a table header, [a.b.c], and makes its table the current
// one. Tables named on the way are created as implicit tables; the last
// is defined, which a table may be once.
func (p *parser) header() {
	at := p.pos
	p.pos++ // '['
	if !p.eof() && p.peek() == '[' {
		p.fail(at, "arrays of tables are not supported")
	}
	p.skipSpace()
	path := p.key()
	p.skipSpace()
	if p.eof() || p.peek() != ']' {
		p.fail(p.pos, "expected ']' to close the table header")
	}
	p.pos++
	name := strings.Join(path, ".")
	t := p.root
	for i, k := range path {
		v, ok := t.get(k)
		if !ok {
			nt := &table{origin: implicit}
			t.set(k, nt)
			v = nt
		}
		sub, isTable := v.(*table)
		if !isTable {
			p.fail(at, "key %q already holds a value and cannot be a table", strings.Join(path[:i+1], "."))
		}
		if i == len(path)-1 {
			switch {
			case !ok:
				sub.origin = defined
			case sub.origin == implicit:
				sub.origin = defined
			case sub.origin == defined:
				p.fail(at, "table [%s] is defined twice", name)
			case sub.origin == dotted:
				p.fail(at, "table [%s] was already defined by dotted keys", name)
			}
		}
		t = sub
	}
	p.cur = t
}

// assign stores a value under a key, relative to the current table. The
// parts of a dotted key before the last name tables, created as dotted
// tables when missing. A key that is already defined is refused.
func (p *parser) assign(path []string, v any, at int) {
	t := p.cur
	for i, k := range path[:len(path)-1] {
		e, ok := t.get(k)
		if !ok {
			nt := &table{origin: dotted}
			t.set(k, nt)
			t = nt
			continue
		}
		sub, isTable := e.(*table)
		if !isTable {
			p.fail(at, "key %q already holds a value and cannot be a table", strings.Join(path[:i+1], "."))
		}
		if sub.origin != dotted {
			p.fail(at, "table %q was defined by a header and cannot be extended by dotted keys", strings.Join(path[:i+1], "."))
		}
		t = sub
	}
	last := path[len(path)-1]
	if _, ok := t.get(last); ok {
		p.fail(at, "duplicate key %q", strings.Join(path, "."))
	}
	t.set(last, v)
}
