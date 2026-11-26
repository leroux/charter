package chartertest

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// copyFixture copies testdata/fixture into a temporary directory so a test
// may mutate it.
func copyFixture(t *testing.T) string {
	t.Helper()
	src := filepath.Join("testdata", "fixture")
	dst := t.TempDir()
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o750)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o600)
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}

// fixtureCompiler is the TypeScript compiler the extractor parses with:
// the repository's node_modules/typescript, which `pnpm install`
// at the repository root installs. The library tests need it as the gate does.
func fixtureCompiler() string {
	p, err := filepath.Abs(filepath.Join("..", "node_modules", "typescript"))
	if err != nil {
		panic(err)
	}
	return p
}

func fixtureOptions(root string) Options {
	return Options{
		Root:        root,
		Entries:     "example.com/fixture/claims",
		EnvPackages: []string{"example.com/fixture/env"},
		EnvLiterals: []string{`^FX_[A-Z_]+$`},
		OpenAPI:     []string{"contract/openapi.json"},
		Events:      "contract/events.json",
		ProbeTag:    "fxprobe",
		TypeScript:  &TypeScript{Module: "claims/index.ts", Import: "@fixture/claims", Compiler: fixtureCompiler()},
	}
}

// replaceIn replaces the one occurrence of old in the file.
func replaceIn(t *testing.T, root, rel, old, repl string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(data), old); n != 1 {
		t.Fatalf("%s: %q occurs %d times, want 1", rel, old, n)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(data), old, repl, 1)), 0o600); err != nil {
		t.Fatal(err)
	}
}

// writeFile creates or overwrites a file under root.
func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// appendTo appends content to a file under root.
func appendTo(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(data, content...), 0o600); err != nil {
		t.Fatal(err)
	}
}

// edit is one mutation of the fixture: a replacement in a file, a file
// written whole, content appended, or a file removed.
type edit struct {
	rel, old, repl string
	write, add     string
	rm             bool
}

func apply(t *testing.T, root string, e edit) {
	t.Helper()
	switch {
	case e.rm:
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(e.rel))); err != nil {
			t.Fatal(err)
		}
	case e.write != "":
		writeFile(t, root, e.rel, e.write)
	case e.add != "":
		appendTo(t, root, e.rel, e.add)
	default:
		replaceIn(t, root, e.rel, e.old, e.repl)
	}
}

func anyContains(list []string, want string) bool {
	for _, s := range list {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}

// noProblems checks the fixture the test built. The TypeScript module is
// rewritten for the tree first, as make claims would, so rule 15 checks the
// test's registry rather than the committed fixture module.
func noProblems(t *testing.T, root string) {
	t.Helper()
	t.Setenv(UpdateTypeScriptEnv, "1")
	problems, err := Problems(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("reported:\n%s", strings.Join(problems, "\n"))
	}
}
