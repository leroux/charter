package toml

import (
	"testing"
)

func TestInvalidUTF8IsRefused(t *testing.T) {
	refused(t, "a = \"caf\xe9\"\n", "UTF-8")
	refused(t, "# \xff\n", "UTF-8")
}

func TestAByteOrderMarkIsSkipped(t *testing.T) {
	decodes(t, "\xEF\xBB\xBFa = 1\n", map[string]any{"a": int64(1)})
}

func TestTabsAreWhitespace(t *testing.T) {
	decodes(t, "a\t=\t\"x\"\t\n\t[t]\t\n", map[string]any{"a": "x", "t": map[string]any{}})
}

func TestCRLFEndsALine(t *testing.T) {
	decodes(t, "a = 1\r\n# c\r\n[t]\r\nb = true\r\n", map[string]any{"a": int64(1), "t": map[string]any{"b": true}})
}

func TestALoneCarriageReturnIsRefused(t *testing.T) {
	refused(t, "a = 1\rb = 2\n", "carriage return")
	refused(t, "a = 1 # c\rb = 2\n")
}

func TestCommentsAddNothing(t *testing.T) {
	decodes(t, "# top\na = 1 # after a pair\n[t] # after a header\n  # indented\n#\n", map[string]any{"a": int64(1), "t": map[string]any{}})
}

func TestAControlCharacterInACommentIsRefused(t *testing.T) {
	refused(t, "# a\x00b\n", "control character")
	refused(t, "a = 1 # \x7f\n", "control character")
	decodes(t, "# a\tb\n", map[string]any{})
}

func TestAHashInAStringIsText(t *testing.T) {
	decodes(t, "a = \"#1\" # c\nb = '# no'\n", map[string]any{"a": "#1", "b": "# no"})
}

func TestTwoPairsOnOneLineAreRefused(t *testing.T) {
	refused(t, "a = 1 b = 2\n", "end of the line")
}

func TestTextAfterAHeaderIsRefused(t *testing.T) {
	refused(t, "[t] x\n", "end of the line")
}
