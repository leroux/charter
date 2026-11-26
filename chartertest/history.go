package chartertest

import (
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"
)

// The gate reads one tree. The drift it cannot see is between two: an
// entry deleted with its anchors, a decision quietly unbuilt, a text or a
// scenario changed and nothing written in Changes, a ruling that was the
// owner's and is now an agent's. History reads the area documents as they
// were at a git revision and as they are, and names each.

// History compares the record with its state at the git revision and
// returns the findings, one a line, each naming the entry. An entry only
// set active, or changed with a Changes line the revision lacked, is not
// named.
//
//charter:realizes CH-261009-QQC8
func History(opts Options, rev string) ([]string, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	cur := newEntriesPackage()
	if _, _, err := readAreas(r.Root, r.areasDir, cur, r.notAreaDocuments()...); err != nil {
		return nil, err
	}
	prev := newEntriesPackage()
	names, err := gitOutput(r.Root, "ls-tree", "--name-only", rev, "--", r.areasDir+"/")
	if err != nil {
		return nil, fmt.Errorf("chartertest: reading the record at %s: %w", rev, err)
	}
	var docs []string
	for _, name := range strings.Split(strings.TrimSpace(string(names)), "\n") {
		if strings.HasSuffix(name, ".md") && path.Base(name) != readmeFile && !slices.Contains(r.notAreaDocuments(), name) {
			docs = append(docs, name)
		}
	}
	sort.Strings(docs)
	var entries []*entry
	for _, name := range docs {
		src, err := gitOutput(r.Root, "show", rev+":"+name)
		if err != nil {
			return nil, fmt.Errorf("chartertest: reading %s at %s: %w", name, rev, err)
		}
		_, es, _ := parseArea(name, string(src))
		entries = append(entries, es...)
	}
	for _, e := range entries {
		if prev.byID[e.ID] == nil {
			prev.add(e)
		}
	}
	return compareHistory(prev, cur, rev), nil
}

// compareHistory names what changed between the entries at the revision
// and the entries now that no Changes line records.
//
//charter:realizes CH-261009-QQC8
func compareHistory(prev, cur *entriesPackage, rev string) []string {
	var out []string
	note := func(e *entry, format string, args ...any) {
		out = append(out, fmt.Sprintf("%s: history: %s", e.Pos, fmt.Sprintf(format, args...)))
	}
	for _, name := range prev.order {
		p := prev.entries[name]
		if p.Level != levelSpec {
			continue
		}
		c := cur.byID[p.ID]
		if c == nil {
			note(p, "%s, in %s at %s, is gone; an entry is deprecated, never deleted", p.ID, p.File, rev)
			continue
		}
		recorded := hasNewChange(p, c)
		switch {
		case p.Status == "active" && c.Status == "planned":
			note(c, "%s went from active to planned; a built decision is deprecated with a Changes line, not unbuilt", c.ID)
		case c.Status == "deprecated" && p.Status != "deprecated" && !recorded:
			note(c, "%s was deprecated with no Changes line saying what replaces it", c.ID)
		}
		if changed := whatChanged(p, c); len(changed) > 0 && !recorded {
			note(c, "%s changed since %s (%s) with no Changes line saying what and why", c.ID, rev, strings.Join(changed, ", "))
		}
		for _, ps := range p.Scenarios {
			for _, cs := range c.Scenarios {
				if cs.Handle == ps.Handle && ps.Provenance.Ruled() && !cs.Provenance.Ruled() {
					note(c, "%s/%s lost its ruling: it was ruled (%s) and is now an agent's, provisional", c.ID, cs.Handle, ps.Provenance)
				}
			}
		}
	}
	return out
}

// hasNewChange reports whether cur has a Changes line prev had not.
func hasNewChange(prev, cur *entry) bool {
	for _, c := range cur.Changes {
		if !slices.ContainsFunc(prev.Changes, func(p Change) bool { return p.Date == c.Date && p.Text == c.Text }) {
			return true
		}
	}
	return false
}

// whatChanged lists what differs between two readings of one entry that
// a Changes line should say: the title, the text, the kind, and each
// scenario added, removed or reworded. A status change is not listed; the
// status rules say which are allowed.
func whatChanged(prev, cur *entry) []string {
	var out []string
	if prev.Text != cur.Text {
		out = append(out, "title")
	}
	if prev.Description != cur.Description {
		out = append(out, "text")
	}
	if prev.Kind != cur.Kind {
		out = append(out, "kind")
	}
	for _, ps := range prev.Scenarios {
		found := false
		for _, cs := range cur.Scenarios {
			if cs.Handle == ps.Handle {
				found = true
				if cs.Text != ps.Text {
					out = append(out, "scenario "+ps.Handle+" reworded")
				}
			}
		}
		if !found {
			out = append(out, "scenario "+ps.Handle+" removed")
		}
	}
	for _, cs := range cur.Scenarios {
		if !prev.hasScenario(cs.Handle) {
			out = append(out, "scenario "+cs.Handle+" added")
		}
	}
	return out
}
