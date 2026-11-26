package toml

import (
	"testing"

	"example.com/toml/claims"
	"github.com/leroux/charter/chartertest"
)

func TestAHeaderDefinesItsTable(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_K75F, "header-defines")
	decodes(t, "top = 1\n[server]\nport = 8080\n[server.tls]\non = true\n", map[string]any{
		"top":    int64(1),
		"server": map[string]any{"port": int64(8080), "tls": map[string]any{"on": true}},
	})
}

func TestAParentIsDefinedLater(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_K75F, "parent-defined-later")
	decodes(t, "[a.b]\nx = 1\n[a]\ny = 2\n", map[string]any{
		"a": map[string]any{"b": map[string]any{"x": int64(1)}, "y": int64(2)},
	})
}

func TestATableDefinedTwiceIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_K75F, "defined-twice-refused")
	refused(t, "[a]\n[a]\n", "defined twice")
	refused(t, "[a.b]\n[a]\n[a.b]\n", "defined twice")
	refused(t, "[a]\n['a']\n", "defined twice")
}

func TestAHeaderCannotDefineADottedTable(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_596G, "header-over-dotted-refused")
	refused(t, "a.b = 1\n[a]\n", "dotted keys")
	refused(t, "[x]\na.b = 1\n[x.a]\n", "dotted keys")
}

func TestDottedKeysCannotExtendAHeadersTable(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_596G, "dotted-into-header-refused")
	refused(t, "[a.b]\nx = 1\n[a]\nb.c = 1\n", "defined by a header")
}

func TestASubtableUnderADottedTable(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_596G, "subtable-under-dotted")
	decodes(t, "[fruit]\napple.color = \"red\"\n[fruit.apple.texture]\nsmooth = true\n", map[string]any{
		"fruit": map[string]any{"apple": map[string]any{"color": "red", "texture": map[string]any{"smooth": true}}},
	})
}

func TestAHeaderCannotReuseAValue(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_DY2F, "header-over-value-refused")
	refused(t, "a = 1\n[a.b]\n", "holds a value")
	refused(t, "a = 1\n[a]\n", "holds a value")
}

func TestADottedKeyCannotReuseAValue(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_DY2F, "dotted-over-value-refused")
	refused(t, "a = 1\na.b = 2\n", "holds a value")
}
