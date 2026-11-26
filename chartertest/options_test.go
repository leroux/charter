package chartertest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestResolveFindsTheModuleAndDirectories(t *testing.T) {
	Verifies(t, claims.CH_037, "root-resolved")
	root := copyFixture(t)
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	if r.modulePath != "example.com/fixture" || r.entriesDir != "claims" || r.envDirs["example.com/fixture/env"] != "env" {
		t.Errorf("resolved = %+v", r)
	}
	if !r.readers["os.Getenv"] || !r.readers["os.LookupEnv"] {
		t.Errorf("default readers missing: %v", r.readers)
	}
}

func TestResolveDefaultsRootToTheCallersModule(t *testing.T) {
	Verifies(t, claims.CH_037, "root-resolved")
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.Root = ""
	r, err := opts.resolve(filepath.Join(root, "pkg", "a"))
	if err != nil {
		t.Fatal(err)
	}
	if r.Root != root {
		t.Errorf("Root = %q, want %q", r.Root, root)
	}
}

func TestResolveFallsBackToTheWorkingDirectory(t *testing.T) {
	Verifies(t, claims.CH_037, "root-resolved")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	want, err := moduleRoot(wd)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Options{Entries: "github.com/leroux/charter"}.resolve("/nonexistent/dir")
	if err != nil {
		t.Fatal(err)
	}
	if r.Root != want {
		t.Errorf("Root = %q, want %q", r.Root, want)
	}
}

func TestResolveRejectsAPackageOutsideTheModule(t *testing.T) {
	Verifies(t, claims.CH_037, "root-resolved")
	root := copyFixture(t)
	opts := fixtureOptions(root)
	opts.Entries = "example.com/elsewhere/claims"
	if _, err := opts.resolve("unused"); err == nil {
		t.Error("want an error for an entries package outside the module")
	}
}

func TestIgnoredMatchesPathsAndBaseNames(t *testing.T) {
	Verifies(t, claims.CH_036, "ignore-patterns")
	r := &resolved{Options: Options{Ignore: []string{"docs", "*_gen.go", "web/node_modules"}}}
	for rel, want := range map[string]bool{
		"docs":                   true,
		"go/x/oas_client_gen.go": true,
		"web/node_modules":   true,
		"go/x/x.go":              false,
		"go/x/testdata":          true,
		"vendor":                 true,
		"go/x/docs_test.go":      false,
	} {
		if got := r.ignored(rel); got != want {
			t.Errorf("ignored(%q) = %v, want %v", rel, got, want)
		}
	}
}

// viaEntryPoint stands in for a public entry point such as Check, Problems
// or Read: it calls callerDir(1) to find the file that called it.
func viaEntryPoint() string { return callerDir(1) }

func TestCallerDirFindsTheCallingFile(t *testing.T) {
	want, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	if got := viaEntryPoint(); got != want {
		t.Fatalf("viaEntryPoint() = %q, want %q", got, want)
	}
}
