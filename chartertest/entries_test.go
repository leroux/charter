package chartertest

import (
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

func parseFixture(t *testing.T, root string) (*resolved, *tree, *entriesPackage, []string) {
	t.Helper()
	r, tr := loadFixture(t, root)
	pkg := newEntriesPackage()
	_, problems, err := readAreas(r.Root, r.areasDir, pkg)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := readSpecs(r.Root, r.specsDir, pkg)
	if err != nil {
		t.Fatal(err)
	}
	problems = append(problems, ps...)
	problems = append(problems, parseObservations(tr, r, pkg)...)
	return r, tr, pkg, problems
}

func TestParseTheFixture(t *testing.T) {
	root := copyFixture(t)
	_, _, pkg, problems := parseFixture(t, root)
	if len(problems) != 0 {
		t.Fatalf("clean fixture reported %q", problems)
	}
	if pkg.name != "claims" || len(pkg.order) != 10 {
		t.Fatalf("package %q with %d entries, want claims with 10", pkg.name, len(pkg.order))
	}
	want := "STK_001 STK_002 SW_001 SW_002 SW_003 CON_001 CON_002 LockIsAdvisory HealthIsJSON Suspect"
	if got := strings.Join(pkg.order, " "); got != want {
		t.Errorf("order = %s\nwant    %s", got, want)
	}
	e := pkg.byID["SW-001"]
	if e == nil || e.Text != "A bucket survives a restart" || e.Status != "active" || strings.Join(e.Accounts, ",") != "env:FX_DATA_DIR,op:POST /keys,event:bucket.upsert,tag:unix,tag:windows" {
		t.Errorf("SW-001 = %+v", e)
	}
	if d := pkg.entries["LockIsAdvisory"]; d.Status != "documented" || d.Source == "" || strings.Join(d.Supports, ",") != "CON-001" {
		t.Errorf("LockIsAdvisory = %+v", d)
	}
}

func TestParseObservationsRejectsNonLiteralFields(t *testing.T) {
	root := copyFixture(t)
	replaceIn(t, root, "claims/facts.go", `Status:   charter.Observed,`, `Status:   charter.ObservationStatus(1),`)
	_, _, _, problems := parseFixture(t, root)
	if !anyContains(problems, "rule 1: HealthIsJSON.Status must be charter.Observed") {
		t.Errorf("problems = %q", problems)
	}
}

func TestCheckMutationsInTheEntriesPackage(t *testing.T) {
	Verifies(t, claims.CH_025, "assignment-refused")
	root := copyFixture(t)
	writeFile(t, root, "claims/gate_test.go", "package claims\n\nimport \"testing\"\n\nfunc TestMutate(t *testing.T) { HealthIsJSON.Text = \"x\" }\n")
	r, tr, pkg, _ := parseFixture(t, root)
	if got := checkMutations(tr, r, pkg); !anyContains(got, "rule 1: HealthIsJSON is assigned outside its declaration") {
		t.Errorf("checkMutations = %q", got)
	}
}
