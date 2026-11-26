package chartertest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// The digest is the charter rendered for a reader who will not read code:
// the needs, grouped by the area that does most to meet each; a map of the
// areas and of which needs each serves; then, area by area in the reading
// order, every active decision's title, Kind and Why.

// UpdateDigestEnv, set to any value, makes the gate write the digest
// instead of comparing it.
const UpdateDigestEnv = "CHARTER_UPDATE_DIGEST"

// Digest renders the charter the options name for a reader, whether or not
// the options keep a digest file.
//
//charter:realizes CH-261005-6ZBS
func Digest(opts Options) (string, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return "", err
	}
	s, err := load(r)
	if err != nil {
		return "", err
	}
	out, _, err := renderDigest(r, s)
	return out, err
}

// digestArea is one area document as the digest reads it.
type digestArea struct {
	doc       *areaDoc
	decisions []*entry // active, in document order
	planned   int
	retired   int
}

// renderDigest renders the digest, and returns the problems of the reading
// order: an area document it leaves out, a line naming none.
//
//charter:realizes CH-261005-6ZBS
func renderDigest(r *resolved, s *state) (string, []string, error) {
	order, problems, err := readingOrder(r, s.docs)
	if err != nil {
		return "", nil, err
	}
	digestDir := path.Join(r.entriesDir, path.Dir(r.digestFile()))
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
	var needsDoc *areaDoc
	for _, d := range s.docs {
		if d.Needs {
			needsDoc = d
		}
	}
	// Entries in document order: by their position in their file.
	var specs []*entry
	for _, name := range s.pkg.order {
		if e := s.pkg.entries[name]; e.Level == levelSpec {
			specs = append(specs, e)
		}
	}
	sort.SliceStable(specs, func(i, j int) bool {
		if specs[i].File != specs[j].File {
			return specs[i].File < specs[j].File
		}
		return lineOf(specs[i].Pos) < lineOf(specs[j].Pos)
	})
	var needs []*entry
	for _, e := range specs {
		if e.isNeed() {
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
	var reach func(e *entry, seen map[string]bool) map[string]bool
	reach = func(e *entry, seen map[string]bool) map[string]bool {
		if got, ok := reaches[e.ID]; ok {
			return got
		}
		out := map[string]bool{}
		if seen[e.ID] {
			return out
		}
		seen[e.ID] = true
		for _, id := range append(append([]string(nil), e.Serves...), e.Realizes...) {
			p := s.pkg.byID[id]
			if p == nil {
				continue
			}
			if p.isNeed() {
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
	fmt.Fprintf(&b, "<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: %s -->\n\n", r.regenerate(UpdateDigestEnv))
	b.WriteString("# The charter, read through\n\n")
	fmt.Fprintf(&b, "What the system does, without its code: %d needs, what a person or the operator can rely on, and %d decisions in %d areas, how the system behaves at each boundary. ", needCount, decisions, len(areas))
	b.WriteString("Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.\n\n")
	b.WriteString("1. [What people need](#what-people-need)\n2. [The areas](#the-areas)\n3. [The decisions, area by area](#the-decisions-area-by-area)\n\n")

	b.WriteString("## What people need\n\n")
	b.WriteString("Each need is listed under the area that does most to meet it.\n\n")
	for _, a := range areas {
		var under []*entry
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
			fmt.Fprintf(&b, "- [%s](%s) %s\n", n.ID, entryLink(link, needsDoc, n), n.Text)
		}
		b.WriteString("\n")
	}
	var orphans []*entry
	for _, n := range needs {
		if home[n.ID] == nil {
			orphans = append(orphans, n)
		}
	}
	if len(orphans) > 0 {
		b.WriteString("**Needs no active decision serves yet**\n\n")
		for _, n := range orphans {
			fmt.Fprintf(&b, "- [%s](%s) %s\n", n.ID, entryLink(link, needsDoc, n), n.Text)
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
			fmt.Fprintf(&b, "- **[%s](%s)** %s%s", e.ID, entryLink(link, a.doc, e), e.Text, kind)
			if why := firstSentence(e.Rationale); why != "" {
				fmt.Fprintf(&b, "  \n  Why: %s", why)
			}
			b.WriteString("\n")
		}
		if a.planned+a.retired > 0 {
			fmt.Fprintf(&b, "\nNot shown: %d planned, %d deprecated.\n", a.planned, a.retired)
		}
	}
	return b.String(), problems, nil
}

// readingOrder lists the area documents, the needs aside, in the reading
// order file's order; one it leaves out goes last, in name order, and is a
// problem, as is a line naming no area document.
func readingOrder(r *resolved, docs []*areaDoc) ([]*areaDoc, []string, error) {
	rel := path.Join(r.entriesDir, readingOrderFile)
	byName := map[string]*areaDoc{}
	for _, d := range docs {
		if !d.Needs {
			byName[path.Base(d.File)] = d
		}
	}
	var order []*areaDoc
	var problems []string
	placed := map[string]bool{}
	data, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		return nil, nil, err
	default:
		for i, l := range strings.Split(string(data), "\n") {
			name := strings.TrimSpace(l)
			if name == "" || strings.HasPrefix(name, "#") {
				continue
			}
			d := byName[name]
			if d == nil || placed[name] {
				problems = append(problems, fmt.Sprintf("%s:%d: rule 26: %q is no area document, or is listed twice; remove the line", rel, i+1, name))
				continue
			}
			placed[name] = true
			order = append(order, d)
		}
	}
	var rest []string
	for name := range byName {
		if !placed[name] {
			rest = append(rest, name)
		}
	}
	sort.Strings(rest)
	for _, name := range rest {
		order = append(order, byName[name])
		if err == nil {
			problems = append(problems, fmt.Sprintf("%s: rule 26: the reading order leaves out %s; add it where a reader should meet it", rel, name))
		}
	}
	return order, problems, nil
}

// checkDigest applies rule 26: the reading order names every area document
// once, and the digest file equals the rendering. With write set it writes
// the digest instead of comparing it.
//
//charter:realizes CH-261005-6ZBS
func checkDigest(r *resolved, s *state, write bool) ([]string, error) {
	if r.Digest == "" {
		return nil, nil
	}
	want, problems, err := renderDigest(r, s)
	if err != nil {
		return nil, err
	}
	rel := path.Join(r.entriesDir, r.Digest)
	p := filepath.Join(r.Root, filepath.FromSlash(rel))
	if write {
		return problems, os.WriteFile(p, []byte(want), 0o600)
	}
	regen := "run " + r.regenerate(UpdateDigestEnv)
	got, err := os.ReadFile(p)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		problems = append(problems, fmt.Sprintf("%s: rule 26: the digest is missing; %s", rel, regen))
	case err != nil:
		return nil, err
	case string(got) != want:
		problems = append(problems, fmt.Sprintf("%s: rule 26: the digest differs from the documents; %s", rel, regen))
	}
	return problems, nil
}

// digestPath is the digest's slash path relative to the root, or "".
func (r *resolved) digestPath() string {
	if r.Digest == "" {
		return ""
	}
	return path.Join(r.entriesDir, r.Digest)
}

func (r *resolved) digestFile() string {
	if r.Digest == "" {
		return "DIGEST.md"
	}
	return r.Digest
}

// regenerate is the command a stale generated file names, the specs, the
// TypeScript module or the digest: the configured one, or the gate with the
// file's environment variable set.
func (r *resolved) regenerate(env string) string {
	if r.Regenerate != "" {
		return r.Regenerate
	}
	return fmt.Sprintf("%s=1 go test -count=1 ./%s", env, r.entriesDir)
}

// entryLink is the link to an entry's heading in its document.
func entryLink(link func(string) string, doc *areaDoc, e *entry) string {
	file := e.File
	if doc != nil && e.Form == formArea {
		file = doc.File
	}
	if e.Form != formArea {
		return link(file)
	}
	return link(file) + "#" + headingSlug(e.ID+" "+e.Text)
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
func areaName(d *areaDoc) string { return strings.TrimSuffix(path.Base(d.File), ".md") }

// lineOf is the line of a "file:line" position, or 0 when it has none.
func lineOf(pos string) int {
	_, n, _ := strings.Cut(pos, ":")
	line, _ := strconv.Atoi(n)
	return line
}
