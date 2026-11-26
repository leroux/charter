// Package digest renders a charter as one document a person reads through:
// the needs under the area that does most to meet them, a map of the areas,
// then each area's active decisions with the first sentence of their why.
package digest

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/leroux/charter/internal/record"
)

type digestArea struct {
	doc       *record.Area
	decisions []*record.Entry
	planned   int
	retired   int
}

// Render is the digest of the record, written as the file digest, relative
// to the claims directory; regen is the command its header names.
//
//charter:realizes CH-261005-6ZBS
func Render(rec *record.Record, digest, regen string) string {
	byName := map[string]*record.Area{}
	for _, d := range rec.Areas {
		byName[d.File] = d
	}
	var order []*record.Area
	for _, f := range rec.ReadingOrder {
		order = append(order, byName[f])
	}
	digestDir := path.Join(rec.Dir, path.Dir(digest))
	link := func(file string) string {
		rel, err := filepath.Rel(filepath.FromSlash(digestDir), filepath.FromSlash(file))
		if err != nil {
			return file
		}
		return filepath.ToSlash(rel)
	}

	byFile := map[string]*digestArea{}
	var areas []*digestArea
	for _, d := range order {
		a := &digestArea{doc: d}
		byFile[d.File] = a
		areas = append(areas, a)
	}
	var needsDoc *record.Area
	for _, d := range rec.Areas {
		if d.Needs {
			needsDoc = d
		}
	}
	// Entries in document order: by their position in their file.
	specs := append([]*record.Entry(nil), rec.Entries...)
	sort.SliceStable(specs, func(i, j int) bool {
		if specs[i].File != specs[j].File {
			return specs[i].File < specs[j].File
		}
		return specs[i].Line < specs[j].Line
	})
	var needs []*record.Entry
	for _, e := range specs {
		if e.Need {
			if e.Status == "active" {
				needs = append(needs, e)
			}
			continue
		}
		a := byFile[e.File]
		if a == nil {
			continue
		}
		switch e.Status {
		case "active":
			a.decisions = append(a.decisions, e)
		case "planned":
			a.planned++
		default:
			a.retired++
		}
	}

	// The needs each active decision reaches through what it serves.
	reaches := map[string]map[string]bool{}
	var reach func(e *record.Entry, seen map[string]bool) map[string]bool
	reach = func(e *record.Entry, seen map[string]bool) map[string]bool {
		if got, ok := reaches[e.ID]; ok {
			return got
		}
		out := map[string]bool{}
		if seen[e.ID] {
			return out
		}
		seen[e.ID] = true
		for _, id := range e.Serves {
			p := rec.Entry(id)
			if p == nil {
				continue
			}
			if p.Need {
				out[p.ID] = true
				continue
			}
			for n := range reach(p, seen) {
				out[n] = true
			}
		}
		reaches[e.ID] = out
		return out
	}
	// count[need][area]: decisions of the area that reach the need.
	count := map[string]map[*digestArea]int{}
	for _, a := range areas {
		for _, e := range a.decisions {
			for n := range reach(e, map[string]bool{}) {
				if count[n] == nil {
					count[n] = map[*digestArea]int{}
				}
				count[n][a]++
			}
		}
	}
	// A need sits under the area that does most to meet it, the earlier in
	// the reading order on a tie.
	home := map[string]*digestArea{}
	for _, n := range needs {
		best := 0
		for _, a := range areas {
			if c := count[n.ID][a]; c > best {
				best, home[n.ID] = c, a
			}
		}
	}

	decisions, needCount := 0, len(needs)
	for _, a := range areas {
		decisions += len(a.decisions)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: %s -->\n\n", regen)
	b.WriteString("# The charter, read through\n\n")
	fmt.Fprintf(&b, "What the system does, without its code: %d needs, what a person or the operator can rely on, and %d decisions in %d areas, how the system behaves at each boundary. ", needCount, decisions, len(areas))
	b.WriteString("Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.\n\n")
	b.WriteString("1. [What people need](#what-people-need)\n2. [The areas](#the-areas)\n3. [The decisions, area by area](#the-decisions-area-by-area)\n\n")

	b.WriteString("## What people need\n\n")
	b.WriteString("Each need is listed under the area that does most to meet it.\n\n")
	for _, a := range areas {
		var under []*record.Entry
		for _, n := range needs {
			if home[n.ID] == a {
				under = append(under, n)
			}
		}
		if len(under) == 0 {
			continue
		}
		fmt.Fprintf(&b, "**%s**\n\n", a.doc.Title)
		for _, n := range under {
			fmt.Fprintf(&b, "- [%s](%s) %s\n", n.ID, entryLink(link, needsDoc, n), n.Title)
		}
		b.WriteString("\n")
	}
	var orphans []*record.Entry
	for _, n := range needs {
		if home[n.ID] == nil {
			orphans = append(orphans, n)
		}
	}
	if len(orphans) > 0 {
		b.WriteString("**Needs no active decision serves yet**\n\n")
		for _, n := range orphans {
			fmt.Fprintf(&b, "- [%s](%s) %s\n", n.ID, entryLink(link, needsDoc, n), n.Title)
		}
		b.WriteString("\n")
	}

	b.WriteString("## The areas\n\n")
	b.WriteString("| # | Area | What it holds | Decisions |\n|---|---|---|---|\n")
	for i, a := range areas {
		n := fmt.Sprint(len(a.decisions))
		if a.planned > 0 {
			n += fmt.Sprintf(", %d planned", a.planned)
		}
		fmt.Fprintf(&b, "| %d | [%s](%s) | %s | %s |\n", i+1, a.doc.Title, link(a.doc.File), cell(firstSentence(a.doc.Summary)), n)
	}
	b.WriteString("\nWhich areas serve each need, with how many of their decisions:\n\n")
	b.WriteString("| Need | Areas |\n|---|---|\n")
	for _, n := range needs {
		var parts []string
		for _, a := range areas {
			if c := count[n.ID][a]; c > 0 {
				parts = append(parts, fmt.Sprintf("%s %d", areaName(a.doc), c))
			}
		}
		if len(parts) == 0 {
			parts = []string{"none yet"}
		}
		fmt.Fprintf(&b, "| [%s](%s) | %s |\n", n.ID, entryLink(link, needsDoc, n), strings.Join(parts, " · "))
	}

	b.WriteString("\n## The decisions, area by area\n")
	for i, a := range areas {
		fmt.Fprintf(&b, "\n### %d. %s\n\n", i+1, a.doc.Title)
		if a.doc.Summary != "" {
			b.WriteString(oneParagraph(a.doc.Summary) + "\n\n")
		}
		for _, e := range a.decisions {
			kind := ""
			if e.Kind != "" {
				kind = " *(" + e.Kind + ")*"
			}
			fmt.Fprintf(&b, "- **[%s](%s)** %s%s", e.ID, entryLink(link, a.doc, e), e.Title, kind)
			if why := firstSentence(e.Why); why != "" {
				fmt.Fprintf(&b, "  \n  Why: %s", why)
			}
			b.WriteString("\n")
		}
		if a.planned+a.retired > 0 {
			fmt.Fprintf(&b, "\nNot shown: %d planned, %d deprecated.\n", a.planned, a.retired)
		}
	}
	return b.String()
}

// readingOrder lists the area documents, the needs aside, in the reading
// order file's order; one it leaves out goes last, in name order, and is a

// entryLink is the link from the digest to an entry's heading.
func entryLink(link func(string) string, doc *record.Area, e *record.Entry) string {
	file := e.File
	if doc != nil {
		file = doc.File
	}
	return link(file) + "#" + headingSlug(e.ID+" "+e.Title)
}

// headingSlug is the anchor GitHub gives a heading: lower case, letters,
// digits, spaces as hyphens, hyphens and underscores kept, the rest dropped.
func headingSlug(heading string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(heading) {
		switch {
		case r == ' ':
			b.WriteRune('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}

// firstSentence is a paragraph's first sentence, on one line.
func firstSentence(text string) string {
	text = oneParagraph(text)
	for i := 0; i+1 < len(text); i++ {
		if text[i] == '.' && text[i+1] == ' ' && !abbreviation(text[:i]) {
			return text[:i+1]
		}
	}
	return text
}

// abbreviation reports text ending in an abbreviation whose period ends no
// sentence.
func abbreviation(text string) bool {
	for _, a := range []string{"e.g", "i.e", "etc", "vs"} {
		if strings.HasSuffix(text, a) {
			return true
		}
	}
	return false
}

// oneParagraph is the text's first paragraph, its lines joined.
func oneParagraph(text string) string {
	para, _, _ := strings.Cut(strings.TrimSpace(text), "\n\n")
	return strings.Join(strings.Fields(para), " ")
}

// cell makes text safe in a table cell.
func cell(text string) string { return strings.ReplaceAll(text, "|", "\\|") }

// areaName is an area document's name without its extension.
func areaName(d *record.Area) string { return strings.TrimSuffix(path.Base(d.File), ".md") }
