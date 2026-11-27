package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
)

// The examples are charters of their own beside the tooling: ratelimit, Go
// alone, and slugs, a nested Go module and Python with no go.mod at its root.
var examples = []string{"../../examples/ratelimit", "../../examples/slugs"}

func charter(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(args, dir, &out, &errOut)
	return out.String() + errOut.String(), code
}

func TestTheExamplesPassTheirOwnGate(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NS7T, "example-gate-passes", "polyglot-example-passes")
	for _, dir := range examples {
		// Clean: no problem, and no warning either.
		if out, code := charter(t, dir, "check"); code != 0 || out != "" {
			t.Errorf("charter check in %s: exit %d\n%s", dir, code, out)
		}
	}
}

// copyExample copies an example into a fresh directory, applying the edit
// to one file.
func copyExample(t *testing.T, dir, file, old, new string) string {
	t.Helper()
	dst := t.TempDir()
	if err := os.CopyFS(dst, os.DirFS(dir)); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dst, filepath.FromSlash(file))
	src, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(src), old) != 1 {
		t.Fatalf("%s holds %q %d times", file, old, strings.Count(string(src), old))
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(src), old, new, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	return dst
}

func TestACopyMissingACitationFailsNamingTheScenario(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NS7T, "example-untested-scenario-named", "polyglot-untested-scenario-named")
	for _, tc := range []struct{ dir, file, old, new, want string }{
		{"../../examples/ratelimit", "limiter/limiter_test.go", `claims.Verifies(t, claims.SW_261126_2XG8, "time-from-clock")`, `_ = claims.SW_261126_2XG8`,
			"scenario time-from-clock of SW-261126-2XG8 is verified by no test"},
		{"../../examples/slugs", "tool/tests/test_slug.py", `claims.verifies(self, claims.SW_261126_X6HB, "separator-from-env")`, `claims.verifies(self, claims.SW_261126_X6HB)`,
			"scenario separator-from-env of SW-261126-X6HB is verified by no test"},
	} {
		out, code := charter(t, copyExample(t, tc.dir, tc.file, tc.old, tc.new), "check")
		if code != 1 || !strings.Contains(out, tc.want) {
			t.Errorf("a copy of %s missing a citation: exit %d\n%s", tc.dir, code, out)
		}
	}
}

func TestTheExamplesImportNothingOfTheTooling(t *testing.T) {
	claims.Verifies(t, claims.CH_261126_NS7T, "nothing-of-the-tooling-imported")
	for _, dir := range examples {
		err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if !strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, ".py") && !strings.HasSuffix(p, ".ts") && filepath.Base(p) != "go.mod" {
				return nil
			}
			src, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			if strings.Contains(string(src), "github.com/leroux/charter") {
				t.Errorf("%s names the tooling; an adopter imports nothing of it", p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTheChartersOwnGateIsGreenOverTheToolingAlone(t *testing.T) {
	claims.Verifies(t, claims.CH_006, "self-gate-walks-the-tooling-only")
	c, err := config.Find(".")
	if err != nil {
		t.Fatal(err)
	}
	tree, problems, err := check.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range problems {
		t.Error(p)
	}
	for _, files := range tree.Files {
		for _, f := range files {
			if strings.HasPrefix(f, "examples/") {
				t.Errorf("the charter's own walk read %s, another charter's", f)
			}
		}
	}
}
