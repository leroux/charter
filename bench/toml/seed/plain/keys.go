package toml

// key reads a key: one or more simple keys joined by dots, with optional
// whitespace around each dot. It returns the parts in order.
func (p *parser) key() []string {
	var parts []string
	for {
		parts = append(parts, p.simpleKey())
		p.skipSpace()
		if p.eof() || p.peek() != '.' {
			return parts
		}
		p.pos++
		p.skipSpace()
	}
}

// simpleKey reads one part of a key: a bare key of ASCII letters, digits,
// underscores and dashes, or a basic or literal string, which may be
// empty and may hold any character a string may.
func (p *parser) simpleKey() string {
	if p.eof() {
		p.fail(p.pos, "expected a key")
	}
	switch c := p.peek(); {
	case c == '"':
		return p.basicString()
	case c == '\'':
		return p.literalString()
	case isBare(c):
		start := p.pos
		for !p.eof() && isBare(p.peek()) {
			p.pos++
		}
		return string(p.src[start:p.pos])
	default:
		p.fail(p.pos, "expected a key, found %q", p.rest())
		return ""
	}
}

func isBare(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-'
}
