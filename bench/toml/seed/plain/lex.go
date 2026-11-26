package toml

import "unicode/utf8"

// checkEncoding refuses a document that is not UTF-8 and skips a leading
// byte order mark.
func (p *parser) checkEncoding() {
	if len(p.src) >= 3 && p.src[0] == 0xEF && p.src[1] == 0xBB && p.src[2] == 0xBF {
		p.src = p.src[3:]
	}
	for i := 0; i < len(p.src); {
		r, n := utf8.DecodeRune(p.src[i:])
		if r == utf8.RuneError && n <= 1 {
			p.fail(i, "the document is not valid UTF-8")
		}
		i += n
	}
}

func (p *parser) eof() bool  { return p.pos >= len(p.src) }
func (p *parser) peek() byte { return p.src[p.pos] }

// skipSpace skips whitespace, which in TOML is space and tab only.
func (p *parser) skipSpace() {
	for !p.eof() && (p.peek() == ' ' || p.peek() == '\t') {
		p.pos++
	}
}

// lineEnd reads what may follow a key/value pair or a header: whitespace,
// an optional comment, then a newline (LF or CRLF) or the end of the
// document. Anything else on the line is refused.
func (p *parser) lineEnd() {
	p.skipSpace()
	if p.eof() {
		return
	}
	if p.peek() == '#' {
		p.comment()
		if p.eof() {
			return
		}
	}
	switch {
	case p.peek() == '\n':
		p.pos++
	case p.peek() == '\r' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '\n':
		p.pos += 2
	case p.peek() == '\r':
		p.fail(p.pos, "a carriage return must be followed by a newline")
	default:
		p.fail(p.pos, "expected the end of the line, found %q", p.rest())
	}
}

// comment skips a comment, from # to the end of the line. A control
// character other than tab is refused inside it.
func (p *parser) comment() {
	p.pos++ // '#'
	for !p.eof() {
		c := p.peek()
		if c == '\n' || (c == '\r' && p.pos+1 < len(p.src) && p.src[p.pos+1] == '\n') {
			return
		}
		if isControl(c) {
			p.fail(p.pos, "control character %U in a comment", rune(c))
		}
		p.pos++
	}
}

// isControl reports whether c is a control character TOML forbids in
// comments and strings: U+0000 to U+001F other than tab, and U+007F.
func isControl(c byte) bool {
	return (c < 0x20 && c != '\t') || c == 0x7F
}

// rest is the remainder of the current line, for error messages.
func (p *parser) rest() string {
	end := p.pos
	for end < len(p.src) && p.src[end] != '\n' && end-p.pos < 20 {
		end++
	}
	return string(p.src[p.pos:end])
}
