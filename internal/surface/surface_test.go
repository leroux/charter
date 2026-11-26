package surface

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/config"
)

func charter(t *testing.T, src string) *config.Config {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "claims"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "claims", "charter.toml"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestEachLineADeclaredCommandPrintsIsAKey(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NTRR, "lines-are-keys")
	c := charter(t, "[[surface]]\nkind = \"key\"\ncommand = \"printf 'outbound.internet\\\\n\\\\nproxy\\\\n'\"\n")
	keys, problems := Declared(c)
	if len(problems) != 0 || !slices.Equal(keys, []string{"key:outbound.internet", "key:proxy"}) {
		t.Errorf("keys %v, problems %q", keys, problems)
	}
	// The command runs from the root.
	c = charter(t, "[[surface]]\nkind = \"key\"\ncommand = \"cat keys.txt\"\n")
	if err := os.WriteFile(filepath.Join(c.Root, "keys.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if keys, problems := Declared(c); len(problems) != 0 || !slices.Equal(keys, []string{"key:a"}) {
		t.Errorf("from the root: keys %v, problems %q", keys, problems)
	}
}

func TestAFailingCommandIsNamed(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NTRR, "failing-command-named")
	c := charter(t, "\n[[surface]]\nkind = \"key\"\ncommand = \"echo no schema >&2; exit 3\"\n")
	_, problems := Declared(c)
	if len(problems) != 1 || !strings.HasPrefix(problems[0], "claims/charter.toml:2: rule 9: the key command \"echo no schema >&2; exit 3\" failed") || !strings.Contains(problems[0], "no schema") {
		t.Errorf("problems = %q", problems)
	}
}

func TestADuplicateKeyIsNamed(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NTRR, "duplicate-key-named")
	c := charter(t, "[[surface]]\nkind = \"key\"\ncommand = \"printf 'a\\\\na\\\\n'\"\n")
	keys, problems := Declared(c)
	if len(problems) != 1 || !strings.Contains(problems[0], "the key command printed a twice") || !slices.Equal(keys, []string{"key:a"}) {
		t.Errorf("keys %v, problems %q", keys, problems)
	}
}

func TestAKindsFormIsChecked(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NTRR, "kind-form-checked")
	for _, tc := range []struct{ src, want string }{
		{"[[surface]]\nkind = \"key\"\ncommand = \"true\"\n[[surface]]\nkind = \"key\"\ncommand = \"true\"\n", "the surface kind key is declared twice"},
		{"[[surface]]\nkind = \"env\"\ncommand = \"true\"\n", "the surface kind env is built in"},
		{"[[surface]]\nkind = \"Build-Tag\"\ncommand = \"true\"\n", "the surface kind \"Build-Tag\" is not a lower-case word"},
	} {
		_, problems := Declared(charter(t, tc.src))
		if len(problems) != 1 || !strings.HasPrefix(problems[0], "claims/charter.toml:1: rule 9: ") || !strings.Contains(problems[0], tc.want) {
			t.Errorf("want %q at its line, got %q", tc.want, problems)
		}
	}
}
