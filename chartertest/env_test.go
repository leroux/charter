package chartertest

import (
	"reflect"
	"strings"
	"testing"
)

func loadFixture(t *testing.T, root string) (*resolved, *tree) {
	t.Helper()
	r, err := fixtureOptions(root).resolve("unused")
	if err != nil {
		t.Fatal(err)
	}
	tr, err := loadTree(r)
	if err != nil {
		t.Fatal(err)
	}
	return r, tr
}

func TestEnvConstantsReadExportedStrings(t *testing.T) {
	root := copyFixture(t)
	_, tr := loadFixture(t, root)
	got := envConstants(tr, "env")
	want := map[string]string{"DataDir": "FX_DATA_DIR", "Timeout": "FX_TIMEOUT"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("envConstants = %v, want %v", got, want)
	}
}

func TestCheckEnvReadsOnTheCleanFixture(t *testing.T) {
	root := copyFixture(t)
	r, tr := loadFixture(t, root)
	consts := map[string]map[string]string{"example.com/fixture/env": envConstants(tr, "env")}
	if got := checkEnvReads(tr, r, consts); len(got) != 0 {
		t.Errorf("clean fixture reported %q", got)
	}
}

func TestCheckEnvReadsFlagsLiteralsAndVariables(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "pkg/a/read.go", "package a\n\nimport \"os\"\n\n// Read reads any variable.\nfunc Read(name string) string { return os.Getenv(name) }\n\nconst dup = \"FX_TIMEOUT\"\n\nconst undeclared = \"FX_OTHER\"\n")
	r, tr := loadFixture(t, root)
	consts := map[string]map[string]string{"example.com/fixture/env": envConstants(tr, "env")}
	got := checkEnvReads(tr, r, consts)
	for _, want := range []string{
		"pkg/a/read.go:6: rule 8: os.Getenv must name a constant",
		`pkg/a/read.go:8: rule 8: literal "FX_TIMEOUT" duplicates env.Timeout`,
		`pkg/a/read.go:10: rule 8: literal "FX_OTHER" looks like an environment variable`,
	} {
		if !anyContains(got, want) {
			t.Errorf("want %q in:\n%s", want, strings.Join(got, "\n"))
		}
	}
}
