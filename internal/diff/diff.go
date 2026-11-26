// Package diff names what a change did to the record that its Changes lines
// do not say: it reads the area documents as they were at a git revision
// and as they are, and names an entry removed, unbuilt, or deprecated or
// changed with no new Changes line, and a scenario whose ruling was lost.
package diff

import (
	"fmt"
	"os/exec"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/record"
)

// Diff compares the record with its state at the git revision and returns
// the findings, one a line, each naming the entry.
//
//charter:realizes CH-261009-QQC8
func Diff(c *config.Config, rev string) ([]string, error) {
	cur, err := record.Read(c.Root, config.Dir, path.Base(c.Digest))
	if err != nil {
		return nil, err
	}
	notArea := []string{record.ReadmeFile, record.QuestionsFile, record.ObservationsFile, record.PolicyFile, record.GlossaryFile, record.AuditFile, path.Base(c.Digest)}
	names, err := git(c.Root, "ls-tree", "--name-only", rev, "--", config.Dir+"/")
	if err != nil {
		return nil, fmt.Errorf("reading the record at %s: %w", rev, err)
	}
	var docs []string
	for _, name := range strings.Split(strings.TrimSpace(names), "\n") {
		if strings.HasSuffix(name, ".md") && !slices.Contains(notArea, path.Base(name)) {
			docs = append(docs, name)
		}
	}
	sort.Strings(docs)
	prev := map[string]*record.Entry{}
	var order []string
	for _, name := range docs {
		src, err := git(c.Root, "show", rev+":"+name)
		if err != nil {
			return nil, fmt.Errorf("reading %s at %s: %w", name, rev, err)
		}
		_, es, _ := record.ParseArea(name, src)
		for _, e := range es {
			if prev[e.ID] == nil {
				prev[e.ID] = e
				order = append(order, e.ID)
			}
		}
	}
	now := map[string]*record.Entry{}
	for _, e := range cur.Entries {
		now[e.ID] = e
	}
	return compare(prev, now, order, rev), nil
}

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		return "", fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(ee.Stderr)))
	}
	return string(out), err
}

// compareHistory names what changed between the entries at the revision
// and the entries now that no Changes line records.
//
//charter:realizes CH-261009-QQC8
func compare(prev, cur map[string]*record.Entry, order []string, rev string) []string {
	var out []string
	note := func(e *record.Entry, format string, args ...any) {
		out = append(out, fmt.Sprintf("%s: history: %s", e.Pos(), fmt.Sprintf(format, args...)))
	}
	for _, id := range order {
		p := prev[id]
		c := cur[id]
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
func hasNewChange(prev, cur *record.Entry) bool {
	for _, c := range cur.Changes {
		if !slices.ContainsFunc(prev.Changes, func(p record.Change) bool { return p.Date == c.Date && p.Text == c.Text }) {
			return true
		}
	}
	return false
}

// whatChanged lists what differs between two readings of one entry that
// a Changes line should say: the title, the text, the kind, and each
// scenario added, removed or reworded. A status change is not listed; the
// status rules say which are allowed.
func whatChanged(prev, cur *record.Entry) []string {
	var out []string
	if prev.Title != cur.Title {
		out = append(out, "title")
	}
	if prev.Text != cur.Text {
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
		if !prev.HasScenario(cs.Handle) {
			out = append(out, "scenario "+cs.Handle+" added")
		}
	}
	return out
}
