// Package walk lists the source files of a charter's tree, by language: the
// whole repository under the root, nested modules included, stopping at a
// directory that holds another charter's configuration.
package walk

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Languages a walk lists files for.
const (
	Go         = "go"
	Python     = "python"
	TypeScript = "typescript"
)

// alwaysSkipped are base names every walk skips.
var alwaysSkipped = map[string]bool{".git": true, "node_modules": true, "vendor": true, "testdata": true}

// Options scope a walk.
type Options struct {
	Root string
	// Languages are the languages to list files for; a language not named
	// lists nothing.
	Languages []string
	// Ignore are path.Match patterns matched against the slash path relative
	// to the root and against the base name.
	Ignore []string
	// Only restricts the walk to these slash paths relative to the root.
	Only []string
	// Boundary is the slash path, relative to a directory, that marks it as
	// another charter's root, such as claims/charter.toml.
	Boundary string
	// Files lists, under each name, the files no language reads that match
	// its path.Match patterns, as Ignore's are matched.
	Files map[string][]string
}

// Tree is the files of each language, slash-relative to the root, sorted.
type Tree map[string][]string

// Walk lists the files.
//
//charter:realizes CH-036, CH-006, CH-261127-45KM
func Walk(o Options) (Tree, error) {
	want := map[string]bool{}
	for _, l := range o.Languages {
		want[l] = true
	}
	only := make([]string, len(o.Only))
	for i, p := range o.Only {
		only[i] = strings.Trim(path.Clean(p), "/")
	}
	tree := Tree{}
	err := filepath.WalkDir(o.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(o.Root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			base := d.Name()
			if alwaysSkipped[base] || strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_") || matchesAny(o.Ignore, rel) || !inScope(only, rel, true) {
				return filepath.SkipDir
			}
			if o.Boundary != "" {
				if _, err := os.Stat(filepath.Join(p, filepath.FromSlash(o.Boundary))); err == nil {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if matchesAny(o.Ignore, rel) || !inScope(only, rel, false) {
			return nil
		}
		if l := language(rel); l != "" && want[l] {
			tree[l] = append(tree[l], rel)
		} else if l == "" {
			// A file is read once, by the first table whose patterns match.
			for _, name := range sortedKeys(o.Files) {
				if matchesAny(o.Files[name], rel) {
					tree[name] = append(tree[name], rel)
					break
				}
			}
		}
		return nil
	})
	for _, files := range tree {
		sort.Strings(files)
	}
	return tree, err
}

// language is the language a file is written in, or "".
func language(rel string) string {
	switch {
	case strings.HasSuffix(rel, ".go"):
		return Go
	case strings.HasSuffix(rel, ".py"):
		return Python
	case strings.HasSuffix(rel, ".d.ts"):
		return ""
	case strings.HasSuffix(rel, ".ts"), strings.HasSuffix(rel, ".tsx"):
		return TypeScript
	}
	return ""
}

// inScope reports whether a path is within only: the path itself or
// something under it is named, or, for a directory, it lies on the way to
// something named, so the walk descends.
func inScope(only []string, rel string, isDir bool) bool {
	if len(only) == 0 {
		return true
	}
	for _, o := range only {
		if rel == o || strings.HasPrefix(rel, o+"/") {
			return true
		}
		if isDir && strings.HasPrefix(o, rel+"/") {
			return true
		}
	}
	return false
}

// MatchesAny reports whether a slash path, or its base name, matches one of
// the path.Match patterns.
func MatchesAny(patterns []string, rel string) bool { return matchesAny(patterns, rel) }

func matchesAny(patterns []string, rel string) bool {
	base := path.Base(rel)
	for _, pat := range patterns {
		if ok, _ := path.Match(pat, rel); ok {
			return true
		}
		if ok, _ := path.Match(pat, base); ok {
			return true
		}
		// A directory pattern covers what is under it.
		if strings.HasPrefix(rel, strings.TrimSuffix(pat, "/")+"/") {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string][]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
