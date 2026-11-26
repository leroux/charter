package toml

import (
	"testing"

	"example.com/toml/claims"
	"github.com/leroux/charter/chartertest"
)

func TestBareKeys(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Y5J5, "bare-chars")
	decodes(t, "a-B_9 = 1\n1234 = 2\n-_- = 3\n", map[string]any{"a-B_9": int64(1), "1234": int64(2), "-_-": int64(3)})
}

func TestQuotedKeysHoldAnyText(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Y5J5, "quoted-any")
	decodes(t, "\"a b.c\\u00e9\" = 1\n'x.y' = 2\n", map[string]any{"a b.cé": int64(1), "x.y": int64(2)})
}

func TestAnEmptyQuotedKeyIsAKey(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Y5J5, "empty-quoted")
	decodes(t, "\"\" = 1\n", map[string]any{"": int64(1)})
}

func TestABadBareKeyIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_Y5J5, "bad-bare-refused")
	refused(t, "a$ = 1\n")
	refused(t, "é = 1\n", "expected a key")
	refused(t, "= 1\n", "expected a key")
}

func TestDottedKeysNest(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_F9XD, "dotted-nests")
	decodes(t, "a.b.c = 1\n", map[string]any{"a": map[string]any{"b": map[string]any{"c": int64(1)}}})
}

func TestWhitespaceAroundADot(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_F9XD, "space-around-dot")
	decodes(t, "a . b = 1\na.\tc = 2\n", map[string]any{"a": map[string]any{"b": int64(1), "c": int64(2)}})
}

func TestAQuotedPartIsOneName(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_F9XD, "quoted-part")
	decodes(t, "site.\"example.com\" = true\n", map[string]any{"site": map[string]any{"example.com": true}})
}

func TestADuplicateKeyIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_KKT9, "duplicate-refused")
	refused(t, "a = 1\na = 1\n", "duplicate key")
	refused(t, "[t]\nx = 1\nx = 'y'\n", "duplicate key")
}

func TestSpellingsOfOneKeyAreOneKey(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_KKT9, "spelling-same-key-refused")
	refused(t, "name = 1\n\"name\" = 2\n", "duplicate key")
	refused(t, "'name' = 1\n\"n\\u0061me\" = 2\n", "duplicate key")
}

func TestADottedKeyIsDefinedOnce(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_KKT9, "dotted-duplicate-refused")
	refused(t, "a.b = 1\na.b = 2\n", "duplicate key")
}
