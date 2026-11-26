package walk

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/leroux/charter/claims"
)

func tree(t *testing.T, files ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, f := range files {
		p := filepath.Join(root, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

var all = []string{Go, Python, TypeScript}

func TestHiddenAndUnderscoreDirectoriesAreSkipped(t *testing.T) {
	claims.Verifies(t, claims.CH_036, "hidden-skipped")
	root := tree(t, "a.go", ".hidden/b.go", "_skip/c.go", "node_modules/d.ts", "vendor/e.go", "pkg/testdata/f.go")
	got, err := Walk(Options{Root: root, Languages: all})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got[Go], []string{"a.go"}) || len(got[TypeScript]) != 0 {
		t.Errorf("walk = %v", got)
	}
}

func TestANestedModuleIsWalked(t *testing.T) {
	claims.Verifies(t, claims.CH_036, "nested-module-walked")
	root := tree(t, "go.mod", "a.go", "proxy/go.mod", "proxy/internal/b.go", "proxy/internal/b_test.go")
	got, err := Walk(Options{Root: root, Languages: all, Boundary: "claims/charter.toml"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a.go", "proxy/internal/b.go", "proxy/internal/b_test.go"}; !slices.Equal(got[Go], want) {
		t.Errorf("walk = %v, want %v", got[Go], want)
	}
}

func TestAnotherChartersTreeIsSkipped(t *testing.T) {
	claims.Verifies(t, claims.CH_036, "other-charter-skipped")
	root := tree(t, "claims/charter.toml", "a.go", "example/claims/charter.toml", "example/limiter/l.go", "example/claims/specs.go")
	got, err := Walk(Options{Root: root, Languages: all, Boundary: "claims/charter.toml"})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got[Go], []string{"a.go"}) {
		t.Errorf("walk = %v; the example's charter is its own", got[Go])
	}
}

func TestIgnorePatternsMatchAPathOrABaseName(t *testing.T) {
	claims.Verifies(t, claims.CH_036, "ignore-patterns", "only-restricts")
	root := tree(t, "a.go", "gen/b.go", "pkg/zz_gen.go", "pkg/c.go", "web/d.ts")
	got, err := Walk(Options{Root: root, Languages: all, Ignore: []string{"gen", "zz_*.go"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a.go", "pkg/c.go"}; !slices.Equal(got[Go], want) {
		t.Errorf("ignored: %v, want %v", got[Go], want)
	}
	got, err = Walk(Options{Root: root, Languages: all, Only: []string{"pkg"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"pkg/c.go", "pkg/zz_gen.go"}; !slices.Equal(got[Go], want) || len(got[TypeScript]) != 0 {
		t.Errorf("only pkg: %v", got)
	}
}

func TestEachLanguageIsListedOnlyWhenNamed(t *testing.T) {
	claims.Verifies(t, claims.CH_036, "typescript-listed", "python-listed")
	root := tree(t, "a.go", "dom0/torii/plan.py", "dom0/tests/test_plan.py", "web/a.ts", "web/b.tsx", "web/types.d.ts")
	got, err := Walk(Options{Root: root, Languages: all})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"web/a.ts", "web/b.tsx"}; !slices.Equal(got[TypeScript], want) {
		t.Errorf("typescript = %v, want %v and no declaration file", got[TypeScript], want)
	}
	if want := []string{"dom0/tests/test_plan.py", "dom0/torii/plan.py"}; !slices.Equal(got[Python], want) {
		t.Errorf("python = %v, want %v", got[Python], want)
	}
}

func TestALanguageTheCharterDoesNotNameIsNotRead(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_PQ7X, "unnamed-language-unread")
	root := tree(t, "a.go", "dom0/torii/plan.py", "web/a.ts")
	got, err := Walk(Options{Root: root, Languages: []string{Go}})
	if err != nil {
		t.Fatal(err)
	}
	if len(got[Python]) != 0 || len(got[TypeScript]) != 0 {
		t.Errorf("a charter naming Go alone listed %v", got)
	}
}
