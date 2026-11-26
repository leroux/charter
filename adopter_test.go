package charter_test

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
)

// goIn runs the go command in dir and returns its combined output.
func goIn(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOWORK=off")
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return out.String(), err
}

func TestTheExampleAdopterPassesItsOwnGate(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_NS7T, "example-gate-passes")
	if out, err := goIn(t, "example", "test", "-count=1", "./..."); err != nil {
		t.Fatalf("the example's tests and gate failed: %v\n%s", err, out)
	}
	if out, err := goIn(t, "example", "run", "./cmd/claimsgen", "lint"); err != nil {
		t.Fatalf("the example's own claimsgen lint named something: %v\n%s", err, out)
	}
	// Nothing of the tooling is copied in: the example's Go imports the
	// tooling, and declares no package of it.
	err := filepath.WalkDir("example", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") {
			return err
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, pkg := range []string{"package chartertest", "package charter\n", "package cli"} {
			if strings.Contains(string(src), pkg) {
				t.Errorf("%s declares %q: the example copied the tooling instead of requiring it", p, strings.TrimSpace(pkg))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestACopyOfTheExampleMissingACitationFailsNamingTheScenario(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261126_NS7T, "example-untested-scenario-named")
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	if err := os.CopyFS(dst, os.DirFS("example")); err != nil {
		t.Fatal(err)
	}
	// The copy requires this checkout's tooling where it now stands.
	edit(t, filepath.Join(dst, "go.mod"), "=> ../", "=> "+root)
	// The test of time-from-clock no longer cites it.
	edit(t, filepath.Join(dst, "limiter", "limiter_test.go"),
		`chartertest.Verifies(t, claims.SW_261126_2XG8, "time-from-clock")`, `_ = claims.SW_261126_2XG8`)
	out, err := goIn(t, dst, "test", "-count=1", "./claims")
	if err == nil {
		t.Fatalf("the gate passed a scenario no test cites:\n%s", out)
	}
	if !strings.Contains(out, "scenario time-from-clock of SW-261126-2XG8 is verified by no test") {
		t.Errorf("the gate failed without naming the scenario:\n%s", out)
	}
}

// edit replaces the one occurrence of old in the file.
func edit(t *testing.T, path, old, new string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(b), old) != 1 {
		t.Fatalf("%s holds %q %d times, want once", path, old, strings.Count(string(b), old))
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(b), old, new, 1)), 0o600); err != nil {
		t.Fatal(err)
	}
}
