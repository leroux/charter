package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
)

func TestTheCharterFlagSelectsAGivenCharterByName(t *testing.T) {
	chartertest.Verifies(t, claims.CH_040, "charter-flag-selects")
	product := Charter{Name: "product", Dir: "claims"}
	lib := Charter{Name: "self", Dir: "lib/claims"}
	for _, tc := range []struct {
		name, want string
	}{{"", "product"}, {"product", "product"}, {"self", "self"}} {
		got, err := pick([]Charter{product, lib}, tc.name)
		if err != nil || got.Name != tc.want {
			t.Errorf("pick(%q) = %q, %v; want %q", tc.name, got.Name, err, tc.want)
		}
	}
	if _, err := pick([]Charter{product, lib}, "other"); err == nil {
		t.Error("a name the command was not given was accepted")
	}
	if _, err := pick(nil, ""); err == nil {
		t.Error("a command given no charter picked one")
	}
}

func TestACharterWithoutARootIsReadAboveTheWorkingDirectory(t *testing.T) {
	chartertest.Verifies(t, claims.CH_040, "root-from-working-directory")
	adopter := t.TempDir()
	if err := os.WriteFile(filepath.Join(adopter, "go.mod"), []byte("module example.com/adopter\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(adopter, "pkg", "deep")
	if err := os.MkdirAll(sub, 0o750); err != nil {
		t.Fatal(err)
	}
	t.Chdir(sub)
	got, err := pick([]Charter{{Name: "product", Dir: "claims", Options: chartertest.Options{Entries: "example.com/adopter/claims"}}}, "")
	if err != nil {
		t.Fatal(err)
	}
	if want, _ := filepath.EvalSymlinks(adopter); got.Options.Root != adopter && got.Options.Root != want {
		t.Errorf("root = %s, want the adopter's module %s, not the tooling's", got.Options.Root, adopter)
	}
	given := Charter{Name: "product", Options: chartertest.Options{Root: "/elsewhere"}}
	if got, _ := pick([]Charter{given}, ""); got.Options.Root != "/elsewhere" {
		t.Errorf("a given root became %s", got.Options.Root)
	}
}
