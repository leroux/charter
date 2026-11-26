package chartertest

import (
	"reflect"
	"testing"

	"github.com/leroux/charter/claims"
)

func TestEnumeratorsReadTheFixture(t *testing.T) {
	Verifies(t, claims.CH_033, "surfaces-enumerated")
	root := copyFixture(t)
	_, tr := loadFixture(t, root)
	if got, want := enumTags(tr, "fxprobe"), []string{"tag:desktop", "tag:unix", "tag:windows"}; !reflect.DeepEqual(got, want) {
		t.Errorf("tags = %v, want %v", got, want)
	}
	ops, err := enumOpenAPI(root, []string{"contract/openapi.json"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"op:GET /health", "op:POST /keys"}; !reflect.DeepEqual(ops, want) {
		t.Errorf("ops = %v, want %v", ops, want)
	}
	events, err := enumEvents(root, "contract/events.json")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"event:bucket.upsert"}; !reflect.DeepEqual(events, want) {
		t.Errorf("events = %v, want %v", events, want)
	}
}

func TestEnumTagsSkipsVersionTags(t *testing.T) {
	Verifies(t, claims.CH_033, "surfaces-enumerated")
	root := copyFixture(t)
	writeFile(t, root, "pkg/b/new.go", "//go:build go1.27 && !cgo\n\npackage b\n")
	_, tr := loadFixture(t, root)
	if got, want := enumTags(tr, "fxprobe"), []string{"tag:cgo", "tag:desktop", "tag:unix", "tag:windows"}; !reflect.DeepEqual(got, want) {
		t.Errorf("tags = %v, want %v", got, want)
	}
}

func TestLoadTreeSkipsHiddenAndNestedModules(t *testing.T) {
	Verifies(t, claims.CH_036, "hidden-and-nested-skipped")
	root := copyFixture(t)
	writeFile(t, root, "pkg/.hidden/h.go", "//go:build hiddentag\n\npackage hidden\n")
	writeFile(t, root, "nested/go.mod", "module example.com/nested\n\ngo 1.27\n")
	writeFile(t, root, "nested/n.go", "//go:build nestedtag\n\npackage nested\n")
	_, tr := loadFixture(t, root)
	if got, want := enumTags(tr, "fxprobe"), []string{"tag:desktop", "tag:unix", "tag:windows"}; !reflect.DeepEqual(got, want) {
		t.Errorf("tags = %v, want %v", got, want)
	}
}
