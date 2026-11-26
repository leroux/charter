package toml

import (
	"testing"
)

func TestBareKeys(t *testing.T) {
	decodes(t, "a-B_9 = 1\n1234 = 2\n-_- = 3\n", map[string]any{"a-B_9": int64(1), "1234": int64(2), "-_-": int64(3)})
}

func TestQuotedKeysHoldAnyText(t *testing.T) {
	decodes(t, "\"a b.c\\u00e9\" = 1\n'x.y' = 2\n", map[string]any{"a b.cé": int64(1), "x.y": int64(2)})
}

func TestAnEmptyQuotedKeyIsAKey(t *testing.T) {
	decodes(t, "\"\" = 1\n", map[string]any{"": int64(1)})
}

func TestABadBareKeyIsRefused(t *testing.T) {
	refused(t, "a$ = 1\n")
	refused(t, "é = 1\n", "expected a key")
	refused(t, "= 1\n", "expected a key")
}

func TestDottedKeysNest(t *testing.T) {
	decodes(t, "a.b.c = 1\n", map[string]any{"a": map[string]any{"b": map[string]any{"c": int64(1)}}})
}

func TestWhitespaceAroundADot(t *testing.T) {
	decodes(t, "a . b = 1\na.\tc = 2\n", map[string]any{"a": map[string]any{"b": int64(1), "c": int64(2)}})
}

func TestAQuotedPartIsOneName(t *testing.T) {
	decodes(t, "site.\"example.com\" = true\n", map[string]any{"site": map[string]any{"example.com": true}})
}

func TestADuplicateKeyIsRefused(t *testing.T) {
	refused(t, "a = 1\na = 1\n", "duplicate key")
	refused(t, "[t]\nx = 1\nx = 'y'\n", "duplicate key")
}

func TestSpellingsOfOneKeyAreOneKey(t *testing.T) {
	refused(t, "name = 1\n\"name\" = 2\n", "duplicate key")
	refused(t, "'name' = 1\n\"n\\u0061me\" = 2\n", "duplicate key")
}

func TestADottedKeyIsDefinedOnce(t *testing.T) {
	refused(t, "a.b = 1\na.b = 2\n", "duplicate key")
}
