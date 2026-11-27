// Package docs reads a charter's docs: markdown files charter.toml names,
// whose prose is free and whose references to entries and to code are
// checked. It parses each doc's references, embeds, headings and curation
// block, says which do not resolve, rewrites the references and embeds to
// match the record and the code, hashes what a doc references, records a
// curation, and names a doc whose references changed since it was curated.
// It names no language: the declarations it resolves code against are the
// readers' facts.
package docs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/leroux/charter/internal/facts"
	"github.com/leroux/charter/internal/record"
)

// Rule is the gate rule a doc's problems are reported under.
const Rule = "rule 26"

// EntryRef is a link whose text opens with an entry's id.
type EntryRef struct {
	facts.At
	ID string
	// Rest is the link text after the id; Target the link's target.
	Rest, Target string
	// From and To are the link's byte offsets in its line.
	From, To int
}

// CodeRef is a link whose text is a name in backticks and whose target is a
// file under the root.
type CodeRef struct {
	facts.At
	Name string
	// Target is the link's target as written; File the file it names,
	// slash-relative to the root, without its fragment.
	Target, File string
	From, To     int
}

// Key is the name a code reference's hash is recorded under.
func (r CodeRef) Key() string { return r.File + "#" + r.Name }

// Embed is a charter:embed comment and, when it has one, the line of its
// closing comment.
type Embed struct {
	facts.At
	ID  string
	End int
}

// Curation is a doc's curation block: who curated it, when, and the hash of
// each reference as it was then, by entry id or by CodeRef.Key.
type Curation struct {
	// Line and End are the lines the block opens and closes on.
	Line   int
	End    int
	By     string
	Date   string
	Hashes map[string]string
}

// Doc is one doc as read.
type Doc struct {
	// File is slash-relative to the root.
	File      string
	Lines     []string
	EntryRefs []EntryRef
	CodeRefs  []CodeRef
	Embeds    []Embed
	// Headings are headings opening with an id: a doc declares no entry.
	Headings []facts.At
	Curation *Curation
}

var (
	link       = regexp.MustCompile(`\[([^\]\n]+)\]\(([^)\s]+)\)`)
	idAtStart  = regexp.MustCompile(`^([A-Z]+-[0-9]+(?:-[A-Z0-9]+)?)\b(.*)$`)
	codeText   = regexp.MustCompile("^`([A-Za-z_$][A-Za-z0-9_$]*(?:\\.[A-Za-z_$#][A-Za-z0-9_$]*)*)`$")
	heading    = regexp.MustCompile(`^#{1,6}\s+([A-Z]+-[0-9]+(?:-[A-Z0-9]+)?)\b`)
	embedOpen  = regexp.MustCompile(`^\s*<!--\s*charter:embed\s+(\S+)\s*-->\s*$`)
	embedClose = regexp.MustCompile(`^\s*<!--\s*/charter:embed\s*-->\s*$`)
	fence      = regexp.MustCompile("^ {0,3}(`{3,}|~{3,})")
	listItem   = regexp.MustCompile(`^\s*([-*+]|\d+[.)])\s`)
)

// codeSpans are the byte ranges of a line's inline code spans: a run of
// backticks opens one, and the next run of the same length closes it.
func codeSpans(line string) [][2]int {
	var out [][2]int
	for i := 0; i < len(line); {
		if line[i] != '`' {
			i++
			continue
		}
		j := i
		for j < len(line) && line[j] == '`' {
			j++
		}
		run := j - i
		closed := false
		for k := j; k < len(line); {
			if line[k] != '`' {
				k++
				continue
			}
			e := k
			for e < len(line) && line[e] == '`' {
				e++
			}
			if e-k == run {
				out = append(out, [2]int{i, e})
				i, closed = e, true
				break
			}
			k = e
		}
		if !closed {
			i = j
		}
	}
	return out
}

func inSpan(spans [][2]int, at int) bool {
	for _, s := range spans {
		if at >= s[0] && at < s[1] {
			return true
		}
	}
	return false
}

const curatedOpen = "<!-- charter:curated"

// Parse reads a doc's text. References inside a fenced or an indented code
// block, an inline code span or an embed are not read.
//
//charter:realizes CH-261127-3NY3
func Parse(file, src string) *Doc {
	d := &Doc{File: file, Lines: strings.Split(src, "\n")}
	open, embed := "", -1 // the fence that opened the block, if one did
	prevBlank, indented, inList := true, false, false
	for i := 0; i < len(d.Lines); i++ {
		line, n := d.Lines[i], i+1
		if m := fence.FindStringSubmatch(line); m != nil {
			switch {
			case open == "":
				open = m[1]
			case m[1][0] == open[0] && len(m[1]) >= len(open) && strings.TrimSpace(line) == m[1]:
				open = ""
			}
			continue
		}
		if open != "" {
			continue
		}
		blank := strings.TrimSpace(line) == ""
		if listItem.MatchString(line) {
			inList = true
		} else if !blank && line[0] != ' ' && line[0] != '\t' {
			inList = false
		}
		// An indented code block: four spaces or a tab, after a blank line
		// or another such line, outside a list.
		isCode := !blank && !inList && (strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")) && (prevBlank || indented)
		indented = isCode || (indented && blank)
		prevBlank = blank
		if isCode {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), curatedOpen) {
			c := &Curation{Line: n, Hashes: map[string]string{}}
			for i++; i < len(d.Lines) && strings.TrimSpace(d.Lines[i]) != "-->"; i++ {
				l := strings.TrimSpace(d.Lines[i])
				switch {
				case strings.HasPrefix(l, "by:"):
					c.By = strings.TrimSpace(strings.TrimPrefix(l, "by:"))
				case strings.HasPrefix(l, "date:"):
					c.Date = strings.TrimSpace(strings.TrimPrefix(l, "date:"))
				default:
					if f := strings.Fields(l); len(f) == 2 {
						c.Hashes[f[0]] = f[1]
					}
				}
			}
			if i < len(d.Lines) {
				c.End = i + 1
			}
			d.Curation = c
			continue
		}
		if m := embedOpen.FindStringSubmatch(line); m != nil && embed < 0 {
			d.Embeds = append(d.Embeds, Embed{At: facts.At{File: file, Line: n}, ID: m[1]})
			embed = len(d.Embeds) - 1
			continue
		}
		if embedClose.MatchString(line) && embed >= 0 {
			d.Embeds[embed].End = n
			embed = -1
			continue
		}
		if embed >= 0 {
			continue
		}
		if m := heading.FindStringSubmatch(line); m != nil && record.IsID(m[1]) {
			d.Headings = append(d.Headings, facts.At{File: file, Line: n})
		}
		spans := codeSpans(line)
		for _, loc := range link.FindAllStringSubmatchIndex(line, -1) {
			if inSpan(spans, loc[0]) {
				continue
			}
			text, target := line[loc[2]:loc[3]], line[loc[4]:loc[5]]
			at := facts.At{File: file, Line: n}
			if m := idAtStart.FindStringSubmatch(text); m != nil && record.IsID(m[1]) {
				d.EntryRefs = append(d.EntryRefs, EntryRef{At: at, ID: m[1], Rest: m[2], Target: target, From: loc[0], To: loc[1]})
				continue
			}
			if m := codeText.FindStringSubmatch(text); m != nil {
				if f, ok := resolveTarget(file, target); ok {
					d.CodeRefs = append(d.CodeRefs, CodeRef{At: at, Name: m[1], Target: target, File: f, From: loc[0], To: loc[1]})
				}
			}
		}
	}
	return d
}

// resolveTarget is the slash path, relative to the root, a link's target
// names from the doc, or false for a URL, a fragment alone or a path
// outside the root.
func resolveTarget(doc, target string) (string, bool) {
	if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
		return "", false
	}
	t, _, _ := strings.Cut(target, "#")
	p := path.Clean(path.Join(path.Dir(doc), t))
	if strings.HasPrefix(p, "../") || p == ".." || path.IsAbs(p) {
		return "", false
	}
	return p, true
}

// Read parses each doc, the files given slash-relative to the root.
func Read(root string, files []string) ([]*Doc, error) {
	var out []*Doc
	for _, f := range files {
		src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(f)))
		if err != nil {
			return nil, err
		}
		out = append(out, Parse(f, string(src)))
	}
	return out, nil
}

// Code is what a doc's code references resolve against: the files the
// charter reads in a language, and the declarations its readers found.
type Code struct {
	Root string
	// Lang holds the files, slash-relative to the root, a language reads.
	Lang  map[string]bool
	Decls map[string][]facts.Decl
}

// NewCode indexes the readers' declarations by file.
func NewCode(root string, langFiles []string, decls []facts.Decl) Code {
	c := Code{Root: root, Lang: map[string]bool{}, Decls: map[string][]facts.Decl{}}
	for _, f := range langFiles {
		c.Lang[f] = true
	}
	for _, d := range decls {
		c.Decls[d.File] = append(c.Decls[d.File], d)
	}
	return c
}

// Resolve finds the declaration a code reference names. For a file no
// language reads, ok is the file's existing and decl is empty. A name
// resolves bare to a declaration of that name, or qualified by the file's
// package or the declaration's type: Accept, engine.Accept, Engine.Accept.
//
//charter:realizes CH-261127-DW8R
func (c Code) Resolve(r CodeRef) (decl *facts.Decl, ok bool) {
	if !c.Lang[r.File] {
		st, err := os.Stat(filepath.Join(c.Root, filepath.FromSlash(r.File)))
		return nil, err == nil && !st.IsDir()
	}
	for i, d := range c.Decls[r.File] {
		if r.Name == d.Name || r.Name == d.Package+"."+d.Name {
			return &c.Decls[r.File][i], true
		}
	}
	return nil, false
}

// Problems names each reference, embed and heading of the docs that the
// gate fails.
//
//charter:realizes CH-261127-3NY3, CH-261127-DW8R
func Problems(ds []*Doc, rec *record.Record, code Code) []string {
	var out []string
	add := func(at facts.At, format string, args ...any) {
		out = append(out, fmt.Sprintf("%s: %s: %s", at.Pos(), Rule, fmt.Sprintf(format, args...)))
	}
	entry := func(at facts.At, what, id string) {
		e := rec.Entry(id)
		switch {
		case e == nil:
			add(at, "%s names %s, which is not an entry of the record", what, id)
		case e.Status == "deprecated":
			add(at, "%s names %s, which is deprecated; name what replaces it, or drop the reference", what, id)
		}
	}
	for _, d := range ds {
		for _, h := range d.Headings {
			add(h, "a doc declares no entry; this heading opens with an id; write the entry in an area document under %s and reference it here", rec.Dir)
		}
		for _, r := range d.EntryRefs {
			entry(r.At, "the reference", r.ID)
		}
		for _, e := range d.Embeds {
			entry(e.At, "the embed", e.ID)
			if e.End == 0 {
				add(e.At, "the embed of %s has no closing <!-- /charter:embed --> comment", e.ID)
			}
		}
		for _, r := range d.CodeRefs {
			if _, ok := code.Resolve(r); ok {
				continue
			}
			if code.Lang[r.File] {
				add(r.At, "%s declares no %s; name a function, type, class, method or constant it declares, or point the link at the file that does", r.File, r.Name)
			} else {
				add(r.At, "the code reference to %s names a file that is not there", r.File)
			}
		}
	}
	sort.Strings(out)
	return out
}

// Slug is the fragment a code host gives a heading: lower case, letters,
// digits, hyphens and underscores kept, spaces made hyphens, the rest
// dropped.
func Slug(heading string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(heading) {
		switch {
		case r == ' ':
			b.WriteByte('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}

// entryTarget is the link from the doc to the entry's heading.
func entryTarget(doc string, e *record.Entry) string {
	rel, err := filepath.Rel(filepath.FromSlash(path.Dir(doc)), filepath.FromSlash(e.File))
	if err != nil {
		rel = e.File
	}
	return filepath.ToSlash(rel) + "#" + Slug(e.ID+" "+e.Title)
}

// Rewrite returns the doc's text with each entry reference's title and
// target, each code reference's line anchor and each embed's body made
// current. A bare id stays bare; a reference that does not resolve is left
// as it is. Nothing else changes.
//
//charter:realizes CH-261127-H4B7
func Rewrite(d *Doc, rec *record.Record, code Code) string {
	lines := append([]string{}, d.Lines...)
	type edit struct {
		from, to int
		text     string
	}
	edits := map[int][]edit{} // line index -> edits, applied from the right
	for _, r := range d.EntryRefs {
		e := rec.Entry(r.ID)
		if e == nil {
			continue
		}
		text := r.ID
		if strings.TrimSpace(r.Rest) != "" {
			text = r.ID + ": " + e.Title
		}
		edits[r.Line-1] = append(edits[r.Line-1], edit{r.From, r.To, "[" + text + "](" + entryTarget(d.File, e) + ")"})
	}
	for _, r := range d.CodeRefs {
		decl, ok := code.Resolve(r)
		if !ok || decl == nil {
			continue
		}
		t, _, _ := strings.Cut(r.Target, "#")
		edits[r.Line-1] = append(edits[r.Line-1], edit{r.From, r.To, "[`" + r.Name + "`](" + t + "#L" + strconv.Itoa(decl.Line) + ")"})
	}
	for i, es := range edits {
		sort.Slice(es, func(a, b int) bool { return es[a].from > es[b].from })
		for _, e := range es {
			lines[i] = lines[i][:e.from] + e.text + lines[i][e.to:]
		}
	}
	// Embeds, from the last, so earlier line numbers hold.
	for k := len(d.Embeds) - 1; k >= 0; k-- {
		em := d.Embeds[k]
		e := rec.Entry(em.ID)
		if e == nil || em.End == 0 {
			continue
		}
		body := embedBody(d.File, e)
		lines = append(lines[:em.Line], append(body, lines[em.End-1:]...)...)
	}
	return strings.Join(lines, "\n")
}

// embedBody is an entry as a quote: its id, linked, and title in bold, then
// its text.
func embedBody(doc string, e *record.Entry) []string {
	out := []string{fmt.Sprintf("> **[%s](%s) %s**", e.ID, entryTarget(doc, e), e.Title), ">"}
	for _, l := range strings.Split(strings.TrimSpace(e.Text), "\n") {
		if strings.TrimSpace(l) == "" {
			out = append(out, ">")
			continue
		}
		out = append(out, "> "+l)
	}
	return out
}

// Hash is a text's hash with its whitespace normalized.
func Hash(s string) string {
	sum := sha256.Sum256([]byte(strings.Join(strings.Fields(s), " ")))
	return hex.EncodeToString(sum[:])[:16]
}

// EntryHash is the hash of an entry's title, status and text: a new Changes
// line alone leaves it as it was.
//
//charter:realizes CH-261127-GXP7
func EntryHash(e *record.Entry) string {
	return Hash(e.Title + "\n" + e.Status + "\n" + e.Text)
}

// CodeHash is the hash of the declaration's source, or of the whole file
// for a file no language reads.
//
//charter:realizes CH-261127-GXP7
func (c Code) CodeHash(r CodeRef) (string, bool) {
	decl, ok := c.Resolve(r)
	if !ok {
		return "", false
	}
	src, err := os.ReadFile(filepath.Join(c.Root, filepath.FromSlash(r.File)))
	if err != nil {
		return "", false
	}
	if decl == nil {
		return Hash(string(src)), true
	}
	lines := strings.Split(string(src), "\n")
	from, to := decl.Line, decl.End
	if from < 1 || to > len(lines) || to < from {
		return "", false
	}
	return Hash(strings.Join(lines[from-1:to], "\n")), true
}

// Hashes are the current hashes of everything the doc references, by entry
// id or CodeRef.Key; a reference that does not resolve has none.
func Hashes(d *Doc, rec *record.Record, code Code) map[string]string {
	out := map[string]string{}
	for _, id := range d.Mentions() {
		if e := rec.Entry(id); e != nil {
			out[id] = EntryHash(e)
		}
	}
	for _, r := range d.CodeRefs {
		if h, ok := code.CodeHash(r); ok {
			out[r.Key()] = h
		}
	}
	return out
}

// Mentions are the ids the doc references or embeds, each once, sorted.
func (d *Doc) Mentions() []string {
	seen := map[string]bool{}
	var out []string
	for _, r := range d.EntryRefs {
		if !seen[r.ID] {
			seen[r.ID] = true
			out = append(out, r.ID)
		}
	}
	for _, e := range d.Embeds {
		if !seen[e.ID] {
			seen[e.ID] = true
			out = append(out, e.ID)
		}
	}
	sort.Strings(out)
	return out
}

// Curate returns the doc's text with its curation block, at its end,
// replaced by one naming the curator, the date and the hashes. It refuses a
// doc whose problems are not none.
//
//charter:realizes CH-261127-GXP7
func Curate(d *Doc, by, date string, hashes map[string]string) string {
	// Every byte before the block stays as it was: a block replaced keeps
	// the lines above it, and a new one follows the doc after a blank line.
	var body string
	if c := d.Curation; c != nil {
		body = strings.Join(d.Lines[:c.Line-1], "\n") + "\n"
		// What follows an old block stays, and the new block goes after it.
		if c.End > 0 && c.End < len(d.Lines) {
			if after := strings.Join(d.Lines[c.End:], "\n"); strings.TrimSpace(after) != "" {
				body += after
				if !strings.HasSuffix(body, "\n") {
					body += "\n"
				}
				body += "\n"
			}
		}
	} else {
		body = strings.Join(d.Lines, "\n")
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body += "\n"
	}
	keys := make([]string, 0, len(hashes))
	for k := range hashes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(body)
	b.WriteString(curatedOpen + "\nby: " + by + "\ndate: " + date + "\n")
	for _, k := range keys {
		b.WriteString(k + " " + hashes[k] + "\n")
	}
	b.WriteString("-->\n")
	return b.String()
}

// Stale names a drafted doc, or each reference of a curated one whose hash
// differs from the one recorded, as check's warnings.
//
//charter:realizes CH-261127-GXP7
func Stale(d *Doc, now map[string]string) []string {
	if d.Curation == nil {
		return []string{fmt.Sprintf("%s:1: warning: the doc is drafted: no one has curated it; read it against the record and the code, then run charter curate %s", d.File, d.File)}
	}
	var out []string
	at := func(key string) int {
		for _, r := range d.EntryRefs {
			if r.ID == key {
				return r.Line
			}
		}
		for _, e := range d.Embeds {
			if e.ID == key {
				return e.Line
			}
		}
		for _, r := range d.CodeRefs {
			if r.Key() == key {
				return r.Line
			}
		}
		return 1
	}
	keys := make([]string, 0, len(now))
	for k := range now {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		was, ok := d.Curation.Hashes[k]
		switch {
		case !ok:
			out = append(out, fmt.Sprintf("%s:%d: warning: the doc is stale: %s is referenced since it was curated on %s; read it again and run charter curate %s", d.File, at(k), k, d.Curation.Date, d.File))
		case was != now[k]:
			out = append(out, fmt.Sprintf("%s:%d: warning: the doc is stale: %s changed since it was curated on %s; read it again and run charter curate %s", d.File, at(k), k, d.Curation.Date, d.File))
		}
	}
	return out
}
