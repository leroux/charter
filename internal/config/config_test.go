package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
)

// writeCharter writes root/claims/charter.toml with src.
func writeCharter(t *testing.T, root, src string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, Dir, File), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFindWalksUpToTheRoot(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_33GK, "found-walking-up")
	root := t.TempDir()
	writeCharter(t, root, "name = \"outer\"\n")
	deep := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	c, err := Find(deep)
	if err != nil {
		t.Fatal(err)
	}
	if c.Root != root || c.Name != "outer" {
		t.Errorf("Find from three below the root = root %s, name %q; want %s, \"outer\"", c.Root, c.Name, root)
	}
}

func TestTheNearestCharterWins(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_33GK, "nearest-wins")
	root := t.TempDir()
	writeCharter(t, root, "name = \"outer\"\n")
	inner := filepath.Join(root, "example")
	writeCharter(t, inner, "name = \"inner\"\n")
	work := filepath.Join(inner, "pkg")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	c, err := Find(work)
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "inner" || c.Root != inner {
		t.Errorf("Find inside the nested charter read %q at %s; want \"inner\" at %s", c.Name, c.Root, inner)
	}
}

func TestNoCharterNamesWhereTheSearchStarted(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_33GK, "none-found-named")
	dir := t.TempDir()
	_, err := Find(dir)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Find with no charter = %v; want ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), dir) {
		t.Errorf("the error %q does not name the directory the search started from, %s", err, dir)
	}
}

const everyKey = `name = "torii"
stakeholder-tests = ["proxy/internal/engine/sim_test.go", "dom0/tests/*"]
probes = ["proxy/internal/itest"]
ignore = ["build"]
only = ["proxy", "dom0"]
env-patterns = ["^TORII_[A-Z0-9_]+$"]
digest = "DIGEST.md"

[go]
module = "github.com/leroux/torii-qubes/claims"
dir = "proxy"
test = "go test ./..."
env = "proxy/internal/env"

[python]
dir = "dom0"
test = "python3 -m unittest tests.test_torii"
env = "dom0/torii/env.py"

[typescript]
dir = "web"
test = "pnpm test"
env = "web/src/env.ts"
compiler = "node_modules/typescript"
import = "@torii/claims"

[[surface]]
kind = "key"
command = "python3 -m torii keys"
`

func TestEveryKeyIsRead(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_33GK, "keys-read")
	root := t.TempDir()
	writeCharter(t, root, everyKey)
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	want := &Config{
		Root: root, Path: filepath.Join(root, Dir, File),
		Name:             "torii",
		StakeholderTests: []string{"proxy/internal/engine/sim_test.go", "dom0/tests/*"},
		Probes:           []string{"proxy/internal/itest"},
		Ignore:           []string{"build"},
		Only:             []string{"proxy", "dom0"},
		EnvPatterns:      []string{"^TORII_[A-Z0-9_]+$"},
		Digest:           "DIGEST.md",
		Go:               &Go{Language: Language{Dir: "proxy", Test: "go test ./...", Env: "proxy/internal/env"}, Module: "github.com/leroux/torii-qubes/claims"},
		Python:           &Language{Dir: "dom0", Test: "python3 -m unittest tests.test_torii", Env: "dom0/torii/env.py"},
		TypeScript:       &TypeScript{Language: Language{Dir: "web", Test: "pnpm test", Env: "web/src/env.ts"}, Compiler: "node_modules/typescript", Import: "@torii/claims"},
		Surfaces:         []Surface{{Kind: "key", Command: "python3 -m torii keys"}},
	}
	c.lines = nil
	if !reflect.DeepEqual(c, want) {
		t.Errorf("Load =\n%+v\nwant\n%+v", c, want)
	}
}

func TestBadKeysAreRefusedNamingTheLine(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_33GK, "unknown-key-refused")
	for _, tc := range []struct{ name, src, want string }{
		{"unknown key", "name = \"x\"\n\nlanguages = [\"go\"]\n", "3: unknown key \"languages\""},
		{"unknown key in a table", "[python]\ndir = \"dom0\"\ntest = \"t\"\nrunner = \"pytest\"\n", "4: unknown key \"python.runner\""},
		{"wrong type", "name = \"x\"\nignore = \"build\"\n", "2:"},
		{"no test command", "name = \"x\"\n\n[python]\ndir = \"dom0\"\n", "3: [python] has no test command"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.src)
			if err == nil {
				t.Fatalf("Parse accepted:\n%s", tc.src)
			}
			if !strings.HasPrefix(err.Error(), tc.want) {
				t.Errorf("Parse error = %q; want it to open %q", err, tc.want)
			}
		})
	}
}
