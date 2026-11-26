package chartertest

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// The earlier form of a spec is one markdown file in the spec directory,
// named for its id, <LENS>-<NNN>-<slug>.md, in the form clew uses:
//
//	**Title**
//	Retries back off and give up at the attempt budget
//
//	**Lens**: SW
//
//	**Status**: active
//
//	**Accounts**: env:NAME, op:GET /path
//
//	**Description**
//	One sentence per line.
//
//	**Rationale**
//	Why the decision was made.
//
//	**Verification Description**
//	What a test does to show it holds.
//
//	## Relations
//
//	**Realizes**
//
//	- [SYS-Name](SYS-Name.md)
//
//	## Changes
//
//	- **2026-09-26** — What changed and why.
//
// It is read while the tree is migrated into area documents (areas.go),
// and its Realizes count as Serves. Rule 21 holds the form, rule 23 the
// relations. The Go declaration of every spec, of either form, is generated
// into one file in the entries package (rule 22).

// UpdateSpecsEnv is the environment variable that makes Check write the
// generated Go instead of failing rule 22.
const UpdateSpecsEnv = "CHARTER_UPDATE_SPECS"

// generatedPrefix opens the first line of the Go generated from the specs.
const generatedPrefix = "// Code generated from "

// generatedFile is the generated Go's name in the entries package.
const generatedFile = "specs.go"

// Lenses are the lenses of the earlier form, in clew's words, which a spec
// file's name carries. An area entry's prefix is not a lens.
var Lenses = []struct{ ID, Description string }{
	{"STK", "Stakeholder spec: what a stakeholder, a person using the system or its operator, needs from the system."},
	{"SYS", "System spec: a capability or behaviour at the system level."},
	{"SW", "Software spec: observable, verifiable behaviour of a component."},
	{"ARCH", "Architecture spec: a structural rule or design decision."},
	{"NF", "Non-functional spec: a quality attribute such as performance or security."},
	{"CON", "Constraint: an invariant that must always hold."},
	{"ENT", "Domain entity: a data shape the system owns, with its attributes and invariants."},
}

// lensIndex orders lenses as clew lists them, an unknown lens last.
func lensIndex(s string) int {
	for i, l := range Lenses {
		if l.ID == s {
			return i
		}
	}
	return len(Lenses)
}

func isLens(s string) bool {
	for _, l := range Lenses {
		if l.ID == s {
			return true
		}
	}
	return false
}

var specStatuses = map[string]string{"planned": "Planned", "active": "Active", "deprecated": "Deprecated"}

// fieldOrder is the order a spec file presents its fields in.
var fieldOrder = []struct {
	name     string
	block    bool
	required bool
}{
	{"Title", true, true},
	{"Lens", false, true},
	{"Status", false, true},
	{"Accounts", false, false},
	{"Description", true, true},
	{"Rationale", true, true},
	{"Verification Description", true, true},
}

var (
	blockLabel  = regexp.MustCompile(`^\*\*([A-Za-z /]+)\*\*$`)
	inlineLabel = regexp.MustCompile(`^\*\*([A-Za-z /]+)\*\*:\s*(.*)$`)
	relationRow = regexp.MustCompile(`^- \[([^\]]+)\]\(([^)]+)\)$`)
	changeRow   = regexp.MustCompile(`^- \*\*\d{4}-\d{2}-\d{2}\*\* — \S`)
)

// specFileName is a spec file's name: SW-004-live-policy-decides.md.
var specFileName = regexp.MustCompile(`^([A-Z]+)-(\d+)-([a-z0-9]+(?:-[a-z0-9]+)*)\.md$`)

// symbolFor is the Go and TypeScript symbol of a spec: its id with each
// hyphen as an underscore, SW_004 or SW_261001_K7Q2. The title is the
// symbol's doc comment, so a title can be refined without renaming every
// citation.
//
//charter:realizes CH-020
func symbolFor(id string) string {
	return strings.ReplaceAll(id, "-", "_")
}

// specName reads a spec file's name.
//
//charter:realizes CH-014
func specName(base string) (id, lens string, number int, ok bool) {
	m := specFileName.FindStringSubmatch(base)
	if m == nil || !isLens(m[1]) {
		return "", "", 0, false
	}
	n, err := strconv.Atoi(m[2])
	if err != nil || n == 0 {
		return "", "", 0, false
	}
	return m[1] + "-" + m[2], m[1], n, true
}

// parseSpec reads one spec file. It returns whatever parsed, with the rule 21
// problems found.
//
//charter:realizes CH-014
func parseSpec(rel, src string) (*entry, []string) {
	var problems []string
	bad := func(line int, format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: rule 21: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	id, lens, number, ok := specName(path.Base(rel))
	if !ok {
		id = strings.TrimSuffix(path.Base(rel), ".md")
		bad(1, "a spec file is named <LENS>-<NNN>-<slug>.md, with a lens of STK, SYS, SW, ARCH, NF, CON or ENT and a lower-case kebab slug")
	}
	e := &entry{ID: id, Level: levelSpec, Form: formFile, File: rel, Pos: rel + ":1", Lens: lens, Number: number, Links: map[string]string{}}
	e.Name = symbolFor(id)

	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	section := "fields"
	next := 0
	seen := map[string]bool{}
	var block *[]string
	var relation *[]string
	blocks := map[string]*[]string{}
	for i, l := range lines {
		n := i + 1
		trimmed := strings.TrimRight(l, " \t")
		if strings.HasPrefix(trimmed, "## ") {
			block = nil
			switch h := strings.TrimSpace(trimmed[3:]); {
			case h == "Relations" && section == "fields":
				section = "relations"
			case h == "Changes" && section != "changes":
				section = "changes"
			default:
				bad(n, "section %q is not Relations or Changes, or is out of order", h)
			}
			continue
		}
		switch section {
		case "relations":
			if trimmed == "" {
				continue
			}
			if m := blockLabel.FindStringSubmatch(trimmed); m != nil {
				switch m[1] {
				case "Realizes":
					relation = &e.Realizes
				case "Concerns":
					relation = &e.Concerns
				case "Related":
					relation = &e.Related
				default:
					bad(n, "relation %q is not Realizes, Concerns or Related", m[1])
					relation = nil
				}
				continue
			}
			m := relationRow.FindStringSubmatch(trimmed)
			switch {
			case m == nil || relation == nil:
				bad(n, "under Relations, write a Realizes, Concerns or Related label and then rows of - [ID](file.md)")
			case !isID(m[1]):
				bad(n, "a relation's link text is the target's id alone, not %q", m[1])
			default:
				*relation = append(*relation, m[1])
				e.Links[m[1]] = m[2]
			}
			continue
		case "changes":
			if strings.HasPrefix(trimmed, "- ") && !changeRow.MatchString(trimmed) {
				bad(n, "a change opens - **YYYY-MM-DD** — and says what changed and why")
			}
			continue
		}
		if m := inlineLabel.FindStringSubmatch(trimmed); m != nil {
			block = nil
			if !claimField(m[1], false, &next, seen, n, bad) {
				continue
			}
			v := strings.TrimSpace(m[2])
			switch m[1] {
			case "Lens":
				if v != lens {
					bad(n, "Lens is %s but the file is named for %s", v, lens)
				}
			case "Status":
				if specStatuses[v] == "" {
					bad(n, "status %q is not planned, active or deprecated", v)
				}
				e.Status = v
			case "Accounts":
				e.Accounts = splitList(v)
				for _, k := range e.Accounts {
					if kind, name, ok := strings.Cut(k, ":"); !ok || kind == "" || name == "" {
						bad(n, "%q is not a surface key such as env:NAME or op:GET /path", k)
					}
				}
			}
			continue
		}
		if m := blockLabel.FindStringSubmatch(trimmed); m != nil {
			block = nil
			if !claimField(m[1], true, &next, seen, n, bad) {
				continue
			}
			var b []string
			block = &b
			blocks[m[1]] = block
			continue
		}
		if block != nil {
			*block = append(*block, trimmed)
			continue
		}
		if trimmed != "" {
			bad(n, "text outside a field; the fields are **Title**, **Lens**, **Status**, **Accounts**, **Description**, **Rationale** and **Verification Description**")
		}
	}
	for _, f := range fieldOrder {
		if f.required && !seen[f.name] {
			bad(1, "%s has no %s", id, f.name)
		}
	}
	text := func(name string) string {
		b := blocks[name]
		if b == nil {
			return ""
		}
		out := *b
		for len(out) > 0 && out[0] == "" {
			out = out[1:]
		}
		for len(out) > 0 && out[len(out)-1] == "" {
			out = out[:len(out)-1]
		}
		return strings.Join(out, "\n")
	}
	e.Text = text("Title")
	if strings.Contains(e.Text, "\n") {
		bad(1, "the Title is one line")
	}
	e.Description, e.Rationale, e.Verification = text("Description"), text("Rationale"), text("Verification Description")
	for _, f := range []struct{ name, v string }{{"Title", e.Text}, {"Description", e.Description}, {"Rationale", e.Rationale}, {"Verification Description", e.Verification}} {
		if seen[f.name] && f.v == "" {
			bad(1, "%s is empty", f.name)
		}
	}
	return e, problems
}

// claimField records that a field label appeared and reports a label that is
// unknown, of the wrong form, repeated or out of order.
func claimField(name string, block bool, next *int, seen map[string]bool, line int, bad func(int, string, ...any)) bool {
	for i, f := range fieldOrder {
		if f.name != name {
			continue
		}
		switch {
		case f.block != block && block:
			bad(line, "write %s inline, **%s**: value", name, name)
		case f.block != block:
			bad(line, "write **%s** on its own line with the text below it", name)
		case seen[name]:
			bad(line, "%s appears twice", name)
		case i < *next:
			bad(line, "%s is out of order; the order is Title, Lens, Status, Accounts, Description, Rationale, Verification Description", name)
		default:
			seen[name] = true
			*next = i + 1
			return true
		}
		return false
	}
	bad(line, "%q is not a field of a spec", name)
	return false
}

func splitList(v string) []string {
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// readSpecs parses every spec file in the spec directory into pkg, each id
// once. A missing directory is no spec files.
//
//charter:realizes CH-012, CH-014
func readSpecs(root, specsDir string, pkg *entriesPackage) ([]string, error) {
	des, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(specsDir)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var specs []*entry
	var problems []string
	for _, de := range des {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(specsDir), de.Name()))
		if err != nil {
			return nil, err
		}
		e, ps := parseSpec(path.Join(specsDir, de.Name()), string(data))
		problems = append(problems, ps...)
		specs = append(specs, e)
	}
	sort.SliceStable(specs, func(i, j int) bool {
		if li, lj := lensIndex(specs[i].Lens), lensIndex(specs[j].Lens); li != lj {
			return li < lj
		}
		return specs[i].Number < specs[j].Number
	})
	for _, e := range specs {
		if other := pkg.byID[e.ID]; other != nil {
			problems = append(problems, fmt.Sprintf("%s: rule 21: %s is declared twice, in %s and %s; %s", e.Pos, e.ID, other.File, e.File, takeNewID(e)))
			continue
		}
		pkg.add(e)
	}
	return problems, nil
}

// renderSpecsGo is the Go generated from the specs: a declaration for every
// spec, by id, of either form and any status, so a test can cite a planned
// spec before it is built. Rule 22 compares the file to it.
//
//charter:realizes CH-013, CH-020
func renderSpecsGo(pkg *entriesPackage, pkgName, source, regen string) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s%s by the charter gate. DO NOT EDIT.\n", generatedPrefix, source)
	fmt.Fprintf(&b, "// Regenerate: %s\n\n", regen)
	fmt.Fprintf(&b, "package %s\n", pkgName)
	imported := false
	for _, name := range pkg.order {
		e := pkg.entries[name]
		if e.Level != levelSpec || specStatuses[e.Status] == "" {
			continue
		}
		if !imported {
			fmt.Fprintf(&b, "\nimport %q\n", charterPath)
			imported = true
		}
		b.WriteString("\n")
		fmt.Fprintf(&b, "// %s is %s: %s\n", e.Name, e.ID, e.Text)
		if e.Description != "" {
			b.WriteString("//\n")
			writeComment(&b, e.Description)
		}
		switch e.Status {
		case "planned":
			fmt.Fprintf(&b, "//\n// Planned: %s is ruled and not built; see %s.\n", e.ID, e.File)
		case "deprecated":
			fmt.Fprintf(&b, "//\n// Deprecated: %s is retired; see its Changes in %s. Add no new anchors.\n", e.ID, e.File)
		}
		kind := ""
		if e.Kind != "" {
			kind = fmt.Sprintf(", Kind: %q", e.Kind)
		}
		fmt.Fprintf(&b, "var %s = &charter.Spec{ID: %q, Title: %q, Status: charter.%s%s}\n",
			e.Name, e.ID, e.Text, specStatuses[e.Status], kind)
	}
	return format.Source(b.Bytes())
}

func writeComment(b *bytes.Buffer, text string) {
	for _, l := range strings.Split(text, "\n") {
		if l == "" {
			b.WriteString("//\n")
		} else {
			fmt.Fprintf(b, "// %s\n", l)
		}
	}
}

// checkGenerated applies rule 22 to the generated Go: it equals the
// rendering. With write set it writes the file instead.
//
//charter:realizes CH-020
func checkGenerated(r *resolved, pkg *entriesPackage, write bool) ([]string, error) {
	pkgName := pkg.name
	if pkgName == "" {
		pkgName = path.Base(r.entriesDir)
	}
	regen := r.regenerate(UpdateSpecsEnv)
	want, err := renderSpecsGo(pkg, pkgName, r.entriesDir, regen)
	if err != nil {
		return []string{fmt.Sprintf("%s: rule 21: the generated Go does not format: %v", r.entriesDir, err)}, nil
	}
	rel := path.Join(r.entriesDir, generatedFile)
	p := filepath.Join(r.Root, filepath.FromSlash(rel))
	if write {
		return nil, os.WriteFile(p, want, 0o600)
	}
	got, err := os.ReadFile(p)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return []string{fmt.Sprintf("%s:1: rule 22: the generated Go is missing; run %s", rel, regen)}, nil
	case err != nil:
		return nil, err
	case !bytes.Equal(got, want):
		return []string{fmt.Sprintf("%s:1: rule 22: the generated Go differs from the specs; run %s", rel, regen)}, nil
	}
	return nil, nil
}

// checkRelations applies rule 23: every relation names an existing spec, no
// spec names itself, and a need serves and realizes nothing. In the earlier
// form a link points at the target's file, unless the target has moved into
// an area document, when the id alone is checked. Observations' Supports
// name existing specs too.
//
//charter:realizes CH-015
func checkRelations(pkg *entriesPackage) []string {
	var problems []string
	for _, name := range pkg.order {
		e := pkg.entries[name]
		if e.Level == levelObservation {
			for _, id := range e.Supports {
				if t := pkg.byID[id]; t == nil || t.Level != levelSpec {
					problems = append(problems, fmt.Sprintf("%s: rule 23: observation %s supports %s, which is not a spec", e.Pos, name, id))
				}
			}
			continue
		}
		if e.Form == formFile && e.Lens == "STK" && len(e.Realizes) > 0 {
			problems = append(problems, fmt.Sprintf("%s: rule 23: %s is a stakeholder spec and realizes nothing; name related needs under Related", e.Pos, e.ID))
		}
		var rels []struct {
			kind string
			ids  []string
		}
		if e.Form == formArea {
			rels = append(rels, struct {
				kind string
				ids  []string
			}{"serves", e.Serves}, struct {
				kind string
				ids  []string
			}{"sees", e.See})
		} else {
			rels = append(rels, struct {
				kind string
				ids  []string
			}{"realizes", e.Realizes}, struct {
				kind string
				ids  []string
			}{"concerns", e.Concerns}, struct {
				kind string
				ids  []string
			}{"is related to", e.Related})
		}
		for _, rel := range rels {
			for _, id := range rel.ids {
				t := pkg.byID[id]
				switch {
				case id == e.ID:
					problems = append(problems, fmt.Sprintf("%s: rule 23: %s %s itself", e.Pos, e.ID, rel.kind))
				case t == nil || t.Level != levelSpec:
					problems = append(problems, fmt.Sprintf("%s: rule 23: %s %s %s, which is not a spec", e.Pos, e.ID, rel.kind, id))
				case e.Form == formFile && t.Form == formFile && e.Links[id] != path.Base(t.File):
					problems = append(problems, fmt.Sprintf("%s: rule 23: the link to %s points at %s; its file is %s", e.Pos, id, e.Links[id], path.Base(t.File)))
				}
			}
		}
	}
	return problems
}

// Lint loads the charter and returns what the gate fails on, and after it
// what a person should look at that the gate does not fail on: how many
// questions are open, a decision that serves nothing, a planned entry whose
// scenarios are all verified, a deprecated entry still anchored or cited,
// an active decision realized only by a package clause, and a citation of
// a scenario'd entry that names no scenario, in the order of their files
// and lines. Last come the style checks (style.go), one line for each that
// found anything, counting what it found. The advisory lines are labelled
// "lint".
//
//charter:realizes CH-041
func Lint(opts Options) ([]string, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	return lint(r, false)
}

// LintStyle is Lint with every finding of the style checks listed, in the
// order of their files and lines among the rest, in place of the counts:
// what claimsgen lint -style prints.
//
//charter:realizes CH-041
func LintStyle(opts Options) ([]string, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	return lint(r, true)
}

// lint is Lint over a resolved charter, the style checks' findings listed
// when listStyle is set and counted when it is not.
//
//charter:realizes CH-041
func lint(r *resolved, listStyle bool) ([]string, error) {
	s, err := load(r)
	if err != nil {
		return nil, err
	}
	// Everything the gate fails on comes first, the whole of it.
	if err := s.evaluate(); err != nil {
		return nil, err
	}
	problems := append([]string{}, s.problems...)
	problems = append(problems, lintQuestions(s)...)
	style, err := lintStyle(s)
	if err != nil {
		return nil, err
	}
	problems = append(problems, style.problems...)
	if listStyle {
		problems = append(problems, style.lines()...)
	}
	for _, name := range s.pkg.order {
		e := s.pkg.entries[name]
		if e.Level != levelSpec {
			continue
		}
		anchors, cites := s.anchors[name], s.cites[name]
		switch e.Status {
		case "deprecated":
			for _, a := range anchors {
				problems = append(problems, fmt.Sprintf("%s:%d: lint: %s is deprecated and still anchored here; strip it", a.File, a.Line, e.ID))
			}
			for _, c := range cites {
				problems = append(problems, fmt.Sprintf("%s:%d: lint: %s is deprecated and still verified here; strip it", c.File, c.Line, e.ID))
			}
			continue
		case "planned":
			if len(e.Scenarios) > 0 {
				all := true
				for _, sc := range e.Scenarios {
					cited := false
					for _, c := range cites {
						for _, h := range c.Handles {
							cited = cited || h == sc.Handle
						}
					}
					all = all && cited
				}
				if all {
					problems = append(problems, fmt.Sprintf("%s: lint: %s is planned and each of its scenarios is verified; set it active", e.Pos, e.ID))
				}
			}
			continue
		}
		if e.Form == formArea && !e.isNeed() && len(e.Serves) == 0 {
			problems = append(problems, fmt.Sprintf("%s: lint: %s serves nothing; name the need or decision it serves", e.Pos, e.ID))
		}
		if !e.isNeed() && len(anchors) > 0 && e.Kind != "structure" {
			onlyPackages := true
			for _, a := range anchors {
				if a.Kind == "realizes" && !strings.HasPrefix(a.Decl, "package ") {
					onlyPackages = false
				}
			}
			if onlyPackages {
				problems = append(problems, fmt.Sprintf("%s: lint: %s is realized only by a package clause; anchor the declaration that carries the decision", e.Pos, e.ID))
			}
		}
		if len(e.Scenarios) > 0 {
			for _, c := range cites {
				if len(c.Handles) == 0 {
					problems = append(problems, fmt.Sprintf("%s:%d: lint: Verifies names %s without a scenario; name the handles the test exercises", c.File, c.Line, e.ID))
				}
			}
		}
	}
	sort.Strings(problems)
	if !listStyle {
		problems = append(problems, style.counts()...)
	}
	return problems, nil
}

func newEntriesPackage() *entriesPackage {
	return &entriesPackage{entries: map[string]*entry{}, byID: map[string]*entry{}}
}
