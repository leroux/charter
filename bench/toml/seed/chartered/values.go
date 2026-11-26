package toml

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// value reads the value of a key/value pair.
//
//charter:realizes SW-261126-H72C
func (p *parser) value() any {
	if p.eof() || p.peek() == '\n' || p.peek() == '\r' || p.peek() == '#' {
		p.fail(p.pos, "expected a value after '='")
	}
	switch c := p.peek(); {
	case c == '"':
		if strings.HasPrefix(string(p.src[p.pos:]), `"""`) {
			p.fail(p.pos, "multi-line strings are not supported")
		}
		return p.basicString()
	case c == '\'':
		if strings.HasPrefix(string(p.src[p.pos:]), "'''") {
			p.fail(p.pos, "multi-line strings are not supported")
		}
		return p.literalString()
	case c == 't' || c == 'f':
		return p.boolean()
	case c == '+' || c == '-' || c >= '0' && c <= '9':
		return p.number()
	case c == '[':
		p.fail(p.pos, "arrays are not supported")
	case c == '{':
		p.fail(p.pos, "inline tables are not supported")
	case c == 'i' || c == 'n':
		p.fail(p.pos, "floats are not supported")
	}
	p.fail(p.pos, "expected a value, found %q", p.rest())
	return nil
}

// basicString reads a string in double quotes on one line, decoding its
// escapes. A control character other than tab is refused, as is an escape
// TOML 1.0 does not define or a \u or \U that names no Unicode scalar
// value.
//
//charter:realizes SW-261126-QZWX
func (p *parser) basicString() string {
	start := p.pos
	p.pos++ // '"'
	var b strings.Builder
	for {
		if p.eof() || p.peek() == '\n' || p.peek() == '\r' {
			p.fail(start, "unterminated string")
		}
		c := p.peek()
		switch {
		case c == '"':
			p.pos++
			return b.String()
		case c == '\\':
			p.escape(&b)
		case isControl(c):
			p.fail(p.pos, "control character %U in a string", rune(c))
		default:
			b.WriteByte(c)
			p.pos++
		}
	}
}

// escape decodes one escape sequence in a basic string.
//
//charter:realizes SW-261126-QZWX
func (p *parser) escape(b *strings.Builder) {
	at := p.pos
	p.pos++ // '\\'
	if p.eof() {
		p.fail(at, "unterminated string")
	}
	c := p.peek()
	p.pos++
	switch c {
	case 'b':
		b.WriteByte('\b')
	case 't':
		b.WriteByte('\t')
	case 'n':
		b.WriteByte('\n')
	case 'f':
		b.WriteByte('\f')
	case 'r':
		b.WriteByte('\r')
	case '"':
		b.WriteByte('"')
	case '\\':
		b.WriteByte('\\')
	case 'u', 'U':
		n := 4
		if c == 'U' {
			n = 8
		}
		if p.pos+n > len(p.src) {
			p.fail(at, "a \\%c escape needs %d hex digits", c, n)
		}
		hex := string(p.src[p.pos : p.pos+n])
		v, err := strconv.ParseUint(hex, 16, 32)
		if err != nil || strings.ContainsAny(hex, "+-") {
			p.fail(at, "a \\%c escape needs %d hex digits", c, n)
		}
		if !utf8.ValidRune(rune(v)) {
			p.fail(at, "\\%c%s is not a Unicode scalar value", c, hex)
		}
		b.WriteRune(rune(v))
		p.pos += n
	default:
		p.fail(at, "invalid escape \\%c", c)
	}
}

// literalString reads a string in single quotes on one line, exactly as
// written: a backslash is a backslash. A control character other than tab
// is refused.
//
//charter:realizes SW-261126-WZ4K
func (p *parser) literalString() string {
	start := p.pos
	p.pos++ // '\''
	for {
		if p.eof() || p.peek() == '\n' || p.peek() == '\r' {
			p.fail(start, "unterminated string")
		}
		c := p.peek()
		if c == '\'' {
			s := string(p.src[start+1 : p.pos])
			p.pos++
			return s
		}
		if isControl(c) {
			p.fail(p.pos, "control character %U in a string", rune(c))
		}
		p.pos++
	}
}

// boolean reads true or false, lower case only.
//
//charter:realizes SW-261126-WN2V
func (p *parser) boolean() bool {
	for _, w := range []string{"true", "false"} {
		if strings.HasPrefix(string(p.src[p.pos:]), w) {
			end := p.pos + len(w)
			if end == len(p.src) || !isBare(p.src[end]) {
				p.pos = end
				return w == "true"
			}
		}
	}
	p.fail(p.pos, "expected a value, found %q", p.rest())
	return false
}

var decimal = regexp.MustCompile(`^[+-]?(0|[1-9][0-9]*)$`)

// number reads an integer: decimal, with an optional sign, no leading
// zeros, within 64 bits. Every other number TOML has is refused by name.
//
//charter:realizes SW-261126-JD2N, SW-261126-H72C
func (p *parser) number() int64 {
	start := p.pos
	for !p.eof() && strings.IndexByte("0123456789abcdefABCDEFxoinIN_+-.:", p.peek()) >= 0 {
		p.pos++
	}
	tok := string(p.src[start:p.pos])
	if decimal.MatchString(tok) {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			p.fail(start, "integer %s does not fit in 64 bits", tok)
		}
		return v
	}
	unsigned := strings.TrimLeft(tok, "+-")
	switch {
	case strings.Contains(tok, ":") || len(unsigned) >= 10 && unsigned[4] == '-' && !strings.ContainsAny(tok[:1], "+-"):
		p.fail(start, "dates and times are not supported")
	case strings.HasPrefix(unsigned, "0x") || strings.HasPrefix(unsigned, "0o") || strings.HasPrefix(unsigned, "0b"):
		p.fail(start, "hexadecimal, octal and binary integers are not supported")
	case strings.ContainsAny(unsigned, ".eEin"):
		p.fail(start, "floats are not supported")
	case strings.Contains(tok, "_"):
		p.fail(start, "underscores in numbers are not supported")
	}
	p.fail(start, "invalid number %q", tok)
	return 0
}
