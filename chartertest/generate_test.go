package chartertest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// Generate writes the generated Go and the TypeScript module from the
// documents alone, so a stale or missing file is no obstacle.
func TestGenerateWritesBothFiles(t *testing.T) {
	Verifies(t, claims.CH_040, "generate-writes-both")
	root := copyFixture(t)
	if err := os.Remove(filepath.Join(root, "claims", "specs.go")); err != nil {
		t.Fatal(err)
	}
	appendTo(t, root, "claims/index.ts", "// stale\n")
	if err := Generate(fixtureOptions(root)); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "claims", "specs.go"))
	if err != nil || !strings.Contains(string(data), `var SW_001 = &charter.Spec{ID: "SW-001"`) {
		t.Errorf("specs.go = %q, %v", data, err)
	}
	got, err := os.ReadFile(filepath.Join(root, "claims", "index.ts"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "fixture", "claims", "index.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("index.ts was not rewritten to the fixture module")
	}
}
