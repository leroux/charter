package chartertest

import (
	"encoding/json"
	"fmt"
	"go/build/constraint"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// enumTags collects every tag named in a build constraint, version tags
// and the probe tag excluded: a probe is a kind of test, not a variant of
// the product.
//
//charter:realizes CH-033
func enumTags(t *tree, probeTag string) []string {
	set := map[string]bool{}
	for _, rel := range t.order {
		f := t.files[rel]
		for _, g := range f.Comments {
			if g.Pos() > f.Package {
				break
			}
			for _, c := range g.List {
				if !constraint.IsGoBuild(c.Text) {
					continue
				}
				expr, err := constraint.Parse(c.Text)
				if err != nil {
					continue
				}
				collectTags(expr, set, probeTag)
			}
		}
	}
	return slices.Sorted(maps.Keys(set))
}

func collectTags(x constraint.Expr, set map[string]bool, probeTag string) {
	switch x := x.(type) {
	case *constraint.TagExpr:
		if !strings.HasPrefix(x.Tag, "go1.") && x.Tag != probeTag {
			set["tag:"+x.Tag] = true
		}
	case *constraint.NotExpr:
		collectTags(x.X, set, probeTag)
	case *constraint.AndExpr:
		collectTags(x.X, set, probeTag)
		collectTags(x.Y, set, probeTag)
	case *constraint.OrExpr:
		collectTags(x.X, set, probeTag)
		collectTags(x.Y, set, probeTag)
	}
}

// enumOpenAPI collects "op:METHOD /path" for every operation in the files.
//
//charter:realizes CH-033
func enumOpenAPI(root string, files []string) ([]string, error) {
	set := map[string]bool{}
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, fmt.Errorf("openapi: %w", err)
		}
		var doc struct {
			Paths map[string]map[string]json.RawMessage `json:"paths"`
		}
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("openapi %s: %w", rel, err)
		}
		for p, methods := range doc.Paths {
			for m := range methods {
				switch m {
				case "get", "post", "put", "delete", "patch", "head", "options":
					set["op:"+strings.ToUpper(m)+" "+p] = true
				}
			}
		}
	}
	return slices.Sorted(maps.Keys(set)), nil
}

// enumEvents collects "event:kind" for every top-level key in the file.
//
//charter:realizes CH-033
func enumEvents(root, rel string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return nil, fmt.Errorf("events: %w", err)
	}
	var kinds map[string]json.RawMessage
	if err := json.Unmarshal(data, &kinds); err != nil {
		return nil, fmt.Errorf("events %s: %w", rel, err)
	}
	set := map[string]bool{}
	for k := range kinds {
		set["event:"+k] = true
	}
	return slices.Sorted(maps.Keys(set)), nil
}
