package record

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
)

// An area document is one markdown file of the claims directory, holding the
// decisions of one area; needs.md holds the needs. Its header is a title and,
// optionally, the code the area covers, then prose that orients the reader;
// then one entry per heading: "## ID Title", the field lines, a blank line,
// the text, then Why, Scenarios, Guidance and Changes, each a paragraph
// opened by its label, in that order.

// statuses an entry may have.
var statuses = map[string]bool{"planned": true, "active": true, "deprecated": true}

// kinds a decision may carry.
var kinds = map[string]bool{"invariant": true, "bound": true, "structure": true}

var (
	// entryHeading takes the heading's first word as its id when it opens
	// as one does, with a prefix and a hyphen; parseID holds it to the form.
	entryHeading = regexp.MustCompile(`^## ([A-Z]+-\S*) (\S.*)$`)
	anyHeading   = regexp.MustCompile(`^#{1,6} `)
	fieldLine    = regexp.MustCompile(`^([A-Z][a-z]+):(?: (.*))?$`)
	scenarioItem = regexp.MustCompile("^\\d+\\. `([^`]*)`(?: \\(([^()]*)\\))?: (\\S.*)$")
	changeItem   = regexp.MustCompile(`^- (\d{4}-\d{2}-\d{2}): (\S.*)$`)
	handleForm   = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// bodyParts are the labelled parts of an entry's body, in order.
var bodyParts = []string{"Why", "Scenarios", "Guidance", "Changes"}

// Symbol is the name a generated file declares an entry by: its id with
// each hyphen an underscore.
func Symbol(id string) string { return strings.ReplaceAll(id, "-", "_") }

// splitList splits a comma list, dropping empty items.
func splitList(v string) []string {
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// isDate reports whether v is a calendar date, YYYY-MM-DD.
func isDate(v string) bool {
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}

// Provenance is who ruled a scenario: the owner, a named delegate, or an
// agent whose choice nobody has ruled on yet. It is written after the
// handle: `handle` (owner, #444), `handle` (delegate name, #12 R3)
// or `handle` (agent, provisional); nothing written is the last.
type Provenance struct {
	// Who is owner, delegate or agent.
	Who string `json:"who"`
	// Name is the delegate's name.
	Name string `json:"name,omitempty"`
	// Ref is what the ruling cites: #444, #12 R3, a date.
	Ref string `json:"ref,omitempty"`
	// Written says the provenance is in the document, not the default.
	Written bool `json:"written"`
}

// String renders the provenance as it is written between the parentheses.
func (p Provenance) String() string {
	switch p.Who {
	case "owner":
		if p.Ref != "" {
			return "owner, " + p.Ref
		}
		return "owner"
	case "delegate":
		if p.Ref != "" {
			return "delegate " + p.Name + ", " + p.Ref
		}
		return "delegate " + p.Name
	}
	return "agent, provisional"
}

// Ruled reports whether a person ruled the scenario.
func (p Provenance) Ruled() bool { return p.Who == "owner" || p.Who == "delegate" }

// provisional is the provenance of a scenario with none written.
var provisional = Provenance{Who: "agent"}

const provenanceForm = "(owner, #444), (delegate name, #12) or (agent, provisional)"

// parseProvenance reads what is written between the parentheses after a
// scenario's handle, returning the form to write when it is not one.
//
//charter:realizes CH-261009-VY5C
func parseProvenance(s string) (Provenance, string) {
	if strings.TrimSpace(s) == "" {
		return provisional, ""
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	head := strings.Fields(parts[0])
	rest := strings.Join(parts[1:], ", ")
	if len(head) == 0 {
		return provisional, fmt.Sprintf("a provenance says who ruled the scenario: %s", provenanceForm)
	}
	p := Provenance{Who: head[0], Written: true}
	switch head[0] {
	case "owner":
		if len(head) > 1 {
			return provisional, fmt.Sprintf("the owner's provenance is (owner) or (owner, #444), not (%s)", s)
		}
		p.Ref = rest
	case "delegate":
		if len(head) != 2 {
			return provisional, fmt.Sprintf("a delegate's provenance names the delegate, (delegate name) or (delegate name, #12), not (%s)", s)
		}
		p.Name, p.Ref = head[1], rest
	case "agent":
		if len(head) > 1 || (rest != "" && rest != "provisional") {
			return provisional, fmt.Sprintf("an agent's scenario is provisional until a person rules it: (agent, provisional), not (%s)", s)
		}
	default:
		return provisional, fmt.Sprintf("%q is not owner, delegate or agent; a provenance is %s", head[0], provenanceForm)
	}
	return p, ""
}

// parseArea reads one area document: its header and its entries.
//
//charter:realizes CH-010, CH-011
func parseArea(rel, src string) (*Area, []*Entry, []string) {
	doc := &Area{File: rel, Needs: path.Base(rel) == NeedsFile}
	var entries []*Entry
	var problems []string
	bad := func(line int, format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: form: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	var cur *Entry
	var body []Line
	inFields := false
	flush := func() {
		if cur == nil {
			return
		}
		problems = append(problems, parseBody(cur, body, rel)...)
		entries = append(entries, cur)
		cur, body = nil, nil
	}
	fence := "" // the open fenced block's marker, "" outside one
	for i, l := range lines {
		n := i + 1
		line := strings.TrimRight(l, " \t")
		// A fenced block is text: no line in it is a heading, nor the
		// header's title or Code line.
		inFence := fence != ""
		if !inFence {
			fence = fenceOpen(line)
			inFence = fence != ""
		} else if fenceCloses(line, fence) {
			fence = ""
		}
		if inFence {
			if cur != nil {
				// A fence among the field lines ends them: it is no field
				// line, and what it holds is the text's.
				if inFields {
					bad(n, "after the heading come field lines, \"Status: active\" and the like, then a blank line before the text")
					inFields = false
				}
				body = append(body, Line{n, line})
			}
			continue
		}
		if strings.HasPrefix(line, "## ") {
			flush()
			m := entryHeading.FindStringSubmatch(line)
			if m == nil {
				bad(n, "an entry's heading is \"## ID Title\", such as \"## SW-118 The bucket refills at the rate\"")
				continue
			}
			id := m[1]
			parts, problem := ParseID(id)
			if problem != "" {
				bad(n, "%s", problem)
				continue
			}
			cur = &Entry{
				ID: id, Symbol: Symbol(id), Title: strings.TrimSpace(m[2]),
				Need: doc.Needs, Prefix: parts.Prefix, Number: parts.Number, Day: parts.Day, Code: parts.Code,
				File: rel, Line: n,
			}
			inFields = true
			continue
		}
		if cur == nil {
			// The document's header: a title, Code, and any prose.
			if t, ok := strings.CutPrefix(line, "# "); ok && doc.Title == "" {
				doc.Title = strings.TrimSpace(t)
				continue
			}
			if anyHeading.MatchString(line) {
				bad(n, "a heading in the header; the header is a title, a Code line and prose, and only \"## ID Title\" headings follow, one per entry")
				continue
			}
			if m := fieldLine.FindStringSubmatch(line); m != nil && m[1] == "Code" {
				doc.Code = splitList(m[2])
				continue
			}
			if line == "" && doc.Summary != "" && !strings.HasSuffix(doc.Summary, "\n\n") {
				doc.Summary += "\n\n"
			} else if line != "" {
				if doc.Summary != "" && !strings.HasSuffix(doc.Summary, "\n\n") {
					doc.Summary += " "
				}
				doc.Summary += strings.TrimSpace(line)
				doc.Prose = append(doc.Prose, Line{n, line})
			}
			continue
		}
		if anyHeading.MatchString(line) {
			bad(n, "a heading inside an entry; only \"## ID Title\" headings are used, one per entry")
			continue
		}
		if inFields {
			// Blank lines before the first field are allowed; the first
			// blank line after one ends the fields.
			if line == "" {
				if len(cur.fields) > 0 {
					inFields = false
				}
				continue
			}
			m := fieldLine.FindStringSubmatch(line)
			if m == nil {
				bad(n, "after the heading come field lines, \"Status: active\" and the like, then a blank line before the text")
				inFields = false
				body = append(body, Line{n, line})
				continue
			}
			problems = append(problems, cur.setField(rel, n, m[1], strings.TrimSpace(m[2]))...)
			continue
		}
		body = append(body, Line{n, line})
	}
	flush()
	doc.Summary = strings.TrimSpace(doc.Summary)
	return doc, entries, problems
}

// fenceOpen returns the marker of the fenced block a line opens, three or
// more backticks or tildes after at most three spaces, or "".
func fenceOpen(line string) string {
	t := strings.TrimLeft(line, " ")
	if len(line)-len(t) > 3 || len(t) < 3 || (t[0] != '`' && t[0] != '~') {
		return ""
	}
	k := 0
	for k < len(t) && t[k] == t[0] {
		k++
	}
	if k < 3 || (t[0] == '`' && strings.Contains(t[k:], "`")) {
		return ""
	}
	return t[:k]
}

// fenceCloses reports whether a line closes the block marker opened: the
// same character, at least as many, after at most three spaces and with
// nothing after.
func fenceCloses(line, marker string) bool {
	t := strings.TrimLeft(line, " ")
	if len(line)-len(t) > 3 {
		return false
	}
	k := 0
	for k < len(t) && t[k] == marker[0] {
		k++
	}
	return k >= len(marker) && strings.TrimSpace(t[k:]) == ""
}

// setField records one field line of an entry and reports a field that is
// unknown, repeated, malformed, or not a need's.
//
//charter:realizes CH-010, CH-011, CH-013
func (e *Entry) setField(rel string, line int, name, value string) []string {
	var problems []string
	bad := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: form: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	if e.fields == nil {
		e.fields = map[string]bool{}
	}
	if e.fields[name] {
		bad("%s appears twice", name)
		return problems
	}
	e.fields[name] = true
	switch name {
	case "Status":
		if !statuses[value] {
			bad("status %q is not planned, active or deprecated", value)
		}
		e.Status = value
	case "Kind":
		if e.Need {
			bad("a need has no Kind")
		} else if !kinds[value] {
			bad("kind %q is not invariant, bound or structure", value)
		}
		e.Kind = value
	case "Serves":
		if e.Need {
			bad("a need serves nothing; it is realized by the decisions that serve it")
		}
		e.Serves = splitList(value)
	case "See":
		e.See = splitList(value)
	case "Accounts":
		e.Accounts = splitList(value)
		for _, k := range e.Accounts {
			if kind, sname, ok := strings.Cut(k, ":"); !ok || kind == "" || sname == "" {
				bad("%q is not a surface key such as env:NAME or op:GET /path", k)
			}
		}
	default:
		bad("%q is not a field of an entry; the fields are Status, Kind, Serves, See and Accounts", name)
	}
	return problems
}

// parseBody reads an entry's body: the decision text, then Why, Scenarios,
// Guidance and Changes, each a paragraph opened by its label, in that order.
//
//charter:realizes CH-010
func parseBody(e *Entry, body []Line, rel string) []string {
	var problems []string
	bad := func(line int, format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: form: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	// Paragraphs: runs of non-blank lines.
	var paras [][]Line
	var cur []Line
	for _, l := range body {
		if strings.TrimSpace(l.Text) == "" {
			if len(cur) > 0 {
				paras = append(paras, cur)
				cur = nil
			}
			continue
		}
		cur = append(cur, l)
	}
	if len(cur) > 0 {
		paras = append(paras, cur)
	}
	last := -1 // index into bodyParts of the last labelled part seen
	var text []string
	for _, p := range paras {
		first := p[0]
		label, rest, ok := strings.Cut(first.Text, ":")
		part := -1
		if ok {
			for i, bp := range bodyParts {
				if label == bp {
					part = i
				}
			}
		}
		if part < 0 {
			if last >= 0 {
				bad(first.N, "text after %s; the decision's text comes first, then Why, Scenarios, Guidance and Changes", bodyParts[last])
				continue
			}
			for _, l := range p {
				text = append(text, l.Text)
			}
			text = append(text, "")
			e.TextLines = append(e.TextLines, p...)
			continue
		}
		if part <= last {
			if part == last {
				bad(first.N, "%s appears twice", label)
			} else {
				bad(first.N, "%s is out of order; the order is Why, Scenarios, Guidance, Changes", label)
			}
			continue
		}
		last = part
		rest = strings.TrimSpace(rest)
		switch label {
		case "Why", "Guidance":
			var ls []string
			var at []Line
			if rest != "" {
				ls = append(ls, rest)
				at = append(at, Line{first.N, rest})
			}
			for _, l := range p[1:] {
				ls = append(ls, strings.TrimSpace(l.Text))
				at = append(at, Line{l.N, strings.TrimSpace(l.Text)})
			}
			v := strings.Join(ls, "\n")
			if v == "" {
				bad(first.N, "%s is empty", label)
			}
			if label == "Why" {
				e.Why = v
				e.WhyLines = at
			} else {
				if e.Need {
					bad(first.N, "a need has no Guidance")
				}
				e.Guidance = v
			}
		case "Scenarios":
			if rest != "" {
				bad(first.N, "Scenarios: is followed by its numbered items on the next lines")
			}
			handles := map[string]bool{}
			for _, l := range p[1:] {
				if strings.HasPrefix(l.Text, " ") || strings.HasPrefix(l.Text, "\t") {
					if k := len(e.Scenarios); k > 0 {
						e.Scenarios[k-1].Text += " " + strings.TrimSpace(l.Text)
					}
					continue
				}
				m := scenarioItem.FindStringSubmatch(l.Text)
				if m == nil {
					bad(l.N, "a scenario is \"N. `handle`: sentence\", or \"N. `handle` (who ruled): sentence\" with the provenance %s; a continuation line is indented", provenanceForm)
					continue
				}
				switch {
				case !handleForm.MatchString(m[1]):
					bad(l.N, "handle %q is not lower-case kebab", m[1])
				case handles[m[1]]:
					bad(l.N, "handle %q appears twice", m[1])
				}
				handles[m[1]] = true
				prov, problem := parseProvenance(m[2])
				if problem != "" {
					bad(l.N, "%s", problem)
				}
				e.Scenarios = append(e.Scenarios, Scenario{Handle: m[1], Text: m[3], Line: l.N, Provenance: prov})
			}
			if len(e.Scenarios) == 0 {
				bad(first.N, "Scenarios lists none")
			}
		case "Changes":
			if rest != "" {
				bad(first.N, "Changes: is followed by its dated items on the next lines")
			}
			for _, l := range p[1:] {
				if strings.HasPrefix(l.Text, " ") || strings.HasPrefix(l.Text, "\t") {
					if k := len(e.Changes); k > 0 {
						e.Changes[k-1].Text += " " + strings.TrimSpace(l.Text)
					}
					continue
				}
				m := changeItem.FindStringSubmatch(l.Text)
				if m == nil {
					bad(l.N, "a change is \"- YYYY-MM-DD: what changed and why\"")
					continue
				}
				if !isDate(m[1]) {
					bad(l.N, "%s is not a date", m[1])
				}
				e.Changes = append(e.Changes, Change{Date: m[1], Text: m[2], Line: l.N})
			}
		}
	}
	e.Text = strings.TrimSpace(strings.Join(text, "\n"))
	line := e.Line // the heading's
	if !e.fields["Status"] {
		bad(line, "%s has no Status", e.ID)
	}
	if e.Text == "" {
		bad(line, "%s has no text: what is observable, at which boundary", e.ID)
	}
	if !strings.Contains(strings.Join(paraLabels(paras), " "), "Why") {
		bad(line, "%s has no Why", e.ID)
	}
	if e.Status == "active" && len(e.Scenarios) == 0 {
		bad(line, "%s is active and lists no Scenarios; each clause of the text gets one", e.ID)
	}
	return problems
}

// paraLabels returns each paragraph's opening label, or "" for prose.
func paraLabels(paras [][]Line) []string {
	var out []string
	for _, p := range paras {
		label, _, ok := strings.Cut(p[0].Text, ":")
		if ok {
			for _, bp := range bodyParts {
				if label == bp {
					out = append(out, label)
				}
			}
		}
	}
	return out
}
