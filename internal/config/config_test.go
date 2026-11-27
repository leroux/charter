package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

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
	claims.Verifies(t, claims.CH_261126_33GK, "found-walking-up")
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
	claims.Verifies(t, claims.CH_261126_33GK, "nearest-wins")
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
	claims.Verifies(t, claims.CH_261126_33GK, "none-found-named")
	dir := t.TempDir()
	_, err := Find(dir)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Find with no charter = %v; want ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), dir) {
		t.Errorf("the error %q does not name the directory the search started from, %s", err, dir)
	}
}

const everyKey = `name = "acme"
stakeholder-tests = ["proxy/internal/engine/sim_test.go", "dom0/tests/*"]
probes = ["proxy/internal/itest"]
ignore = ["build"]
only = ["proxy", "dom0"]
env-patterns = ["^ACME_[A-Z0-9_]+$"]
digest = "DIGEST.md"

[go]
module = "github.com/leroux/acme-svc/claims"
dir = "proxy"
test = "go test ./..."
env = "proxy/internal/env"

[python]
dir = "dom0"
test = "python3 -m unittest tests.test_acme"
env = "dom0/acme/env.py"

[typescript]
dir = "web"
test = "pnpm test"
env = "web/src/env.ts"
compiler = "node_modules/typescript"
import = "@acme/claims"

[[surface]]
kind = "key"
command = "python3 -m acme keys"
`

func TestEveryKeyIsRead(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_33GK, "keys-read")
	root := t.TempDir()
	writeCharter(t, root, everyKey)
	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	want := &Config{
		Root: root, Path: filepath.Join(root, Dir, File),
		Name:             "acme",
		StakeholderTests: []string{"proxy/internal/engine/sim_test.go", "dom0/tests/*"},
		Probes:           []string{"proxy/internal/itest"},
		Ignore:           []string{"build"},
		Only:             []string{"proxy", "dom0"},
		EnvPatterns:      []string{"^ACME_[A-Z0-9_]+$"},
		Digest:           "DIGEST.md",
		Go:               &Go{Language: Language{Dir: "proxy", Test: "go test ./...", Env: "proxy/internal/env"}, Module: "github.com/leroux/acme-svc/claims"},
		Python:           &Language{Dir: "dom0", Test: "python3 -m unittest tests.test_acme", Env: "dom0/acme/env.py"},
		TypeScript:       &TypeScript{Language: Language{Dir: "web", Test: "pnpm test", Env: "web/src/env.ts"}, Compiler: "node_modules/typescript", Import: "@acme/claims"},
		Surfaces:         []Surface{{Kind: "key", Command: "python3 -m acme keys"}},
	}
	c.lines, c.surfaceLines = nil, nil
	if !reflect.DeepEqual(c, want) {
		t.Errorf("Load =\n%+v\nwant\n%+v", c, want)
	}
}

func TestBadKeysAreRefusedNamingTheLine(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_33GK, "unknown-key-refused")
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

func TestAFilesTableNeedsPatternsAndAMarker(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_45KM, "files-table-checked")
	for _, tc := range []struct{ name, src, want string }{
		{"no patterns", "name = \"x\"\n\n[[files]]\ncomment = \"#\"\n", "3: [[files]] has no patterns"},
		{"no marker", "name = \"x\"\n\n[[files]]\npatterns = [\"*.nft\"]\n\n[[files]]\npatterns = [\"*.sh\"]\n", "6: [[files]] has no comment marker"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.src)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Parse error = %v; want it to hold %q", err, tc.want)
			}
		})
	}
	c, err := Parse("name = \"x\"\n\n[[files]]\npatterns = [\"packaging/*/*.nft\"]\ncomment = \"#\"\n")
	if err != nil || len(c.Files) != 1 || c.Files[0].Comment != "#" || c.Files[0].Patterns[0] != "packaging/*/*.nft" {
		t.Errorf("Parse = %+v, %v", c, err)
	}
}
