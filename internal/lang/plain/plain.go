// Package plain reads anchors in files no language reads, such as a
// firewall ruleset or a service unit: a comment line naming an entry.
package plain

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leroux/charter/internal/facts"
)

// Reader reads the files of one [[files]] table.
type Reader struct {
	Root string
	// Comment is the marker a comment line opens with, such as "#".
	Comment string
}

// Read returns the anchors in files, slash paths relative to the root: a
// comment line holding a realizes or concerns directive and a list of ids
// anchors the file it is in, and a directive that names no id is refused.
//
//charter:realizes CH-261127-45KM
func (r Reader) Read(files []string) (*facts.Set, error) {
	set := &facts.Set{}
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(bytes.NewReader(src))
		sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
		for n := 1; sc.Scan(); n++ {
			kind, rest := directive(sc.Text(), r.Comment)
			if kind == "" {
				continue
			}
			at := facts.At{File: rel, Line: n}
			ids := splitList(rest)
			if len(ids) == 0 {
				set.Refused = append(set.Refused, facts.Refused{At: at, Rule: "rule 24",
					Message: fmt.Sprintf("charter:%s names at least one spec", kind)})
				continue
			}
			set.Anchors = append(set.Anchors, facts.Anchor{At: at, Kind: kind, IDs: ids, Decl: "file " + rel, DeclKind: "file"})
		}
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", rel, err)
		}
	}
	return set, nil
}

// directive returns the kind and the rest of a directive line, or "".
func directive(line, comment string) (kind, rest string) {
	text, ok := strings.CutPrefix(strings.TrimSpace(line), comment)
	if !ok {
		return "", ""
	}
	text = strings.TrimLeft(text, " \t")
	for _, k := range []string{"realizes", "concerns"} {
		if r, ok := strings.CutPrefix(text, "charter:"+k); ok && (r == "" || r[0] == ' ' || r[0] == '\t') {
			return k, r
		}
	}
	return "", ""
}

func splitList(s string) []string {
	var out []string
	for _, f := range strings.Split(s, ",") {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
}
