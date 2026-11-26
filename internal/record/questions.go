package record

import (
	"fmt"
	"regexp"
	"strings"
)

// The question file, questions.md beside the area documents, holds what
// nobody ruled: an agent that had to choose something no entry decides
// writes the question, who asked, the entry it is about and the provisional
// answer it took, and goes on. A person answers in place, opening the Answer
// with the date; the question stays, answered, as the record of the ruling.

// Question is one entry of the question file.
type Question struct {
	ID, Text    string
	AskedBy     string
	Entry       string // an id, or ID/handle
	Provisional string
	Answer      string // "" while open
	File        string
	Line        int
}

// Open reports whether nobody has answered yet.
func (q Question) Open() bool { return q.Answer == "" }

var (
	questionHeading = regexp.MustCompile(`^## (Q-\S+) (\S.*)$`)
	questionField   = regexp.MustCompile(`^([A-Z][a-z]+(?:-[a-z]+)?):(?: (.*))?$`)
	answerOpening   = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}), (\S.*)$`)
)

// parseQuestions reads the question file's text against the record's
// entries, each form problem with its file and line.
//
//charter:realizes CH-261009-4B3B
func parseQuestions(rel, src string, r *Record) ([]Question, []string) {
	var qs []Question
	var problems []string
	bad := func(line int, format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: form: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	ids := map[string]int{}
	var cur *Question
	fields := map[string]bool{}
	last := "" // the field a continuation line extends
	finish := func() {
		if cur == nil {
			return
		}
		for _, f := range []string{"Asked-by", "Entry", "Provisional"} {
			if !fields[f] {
				bad(cur.Line, "%s has no %s; a question says who asked, which entry it is about, and the provisional answer taken", cur.ID, f)
			}
		}
		if cur.Entry != "" {
			id, handle, _ := strings.Cut(cur.Entry, "/")
			e := r.byID[id]
			switch {
			case e == nil:
				bad(cur.Line, "%s is about %s, which is not an entry", cur.ID, id)
			case handle != "" && !e.HasScenario(handle):
				bad(cur.Line, "%s is about %s/%s, a scenario %s does not list", cur.ID, id, handle, id)
			}
		}
		qs = append(qs, *cur)
		cur, fields, last = nil, map[string]bool{}, ""
	}
	for i, l := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		n := i + 1
		line := strings.TrimRight(l, " \t")
		if strings.HasPrefix(line, "## ") {
			finish()
			m := questionHeading.FindStringSubmatch(line)
			if m == nil {
				bad(n, "a question's heading is \"## Q-ID question\", the id from charter next Q")
				continue
			}
			parts, problem := ParseID(m[1])
			if problem != "" {
				bad(n, "%s", problem)
				continue
			}
			if parts.Prefix != "Q" {
				bad(n, "a question's id has the prefix Q, not %s", parts.Prefix)
			}
			if first, ok := ids[m[1]]; ok {
				bad(n, "%s is declared twice, at lines %d and %d; %s", m[1], first, n, takeNewID("Q"))
			}
			ids[m[1]] = n
			cur = &Question{ID: m[1], Text: strings.TrimSpace(m[2]), File: rel, Line: n}
			continue
		}
		if cur == nil || line == "" {
			continue
		}
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			if last != "" {
				cur.extend(last, strings.TrimSpace(line))
			}
			continue
		}
		m := questionField.FindStringSubmatch(line)
		if m == nil {
			bad(n, "a question's lines are its fields, Asked-by, Entry, Provisional and Answer, each \"Name: value\"; a continuation line is indented")
			continue
		}
		name, value := m[1], strings.TrimSpace(m[2])
		if fields[name] {
			bad(n, "%s appears twice", name)
			continue
		}
		fields[name] = true
		last = name
		switch name {
		case "Asked-by":
			cur.AskedBy = value
		case "Entry":
			cur.Entry = value
		case "Provisional":
			cur.Provisional = value
		case "Answer":
			if m := answerOpening.FindStringSubmatch(value); m == nil || !isDate(m[1]) {
				bad(n, "an Answer opens with the date and who answered: \"Answer: 2026-10-09, owner: yes\"")
			}
			cur.Answer = value
		default:
			bad(n, "%q is not a field of a question; the fields are Asked-by, Entry, Provisional and Answer", name)
		}
	}
	finish()
	return qs, problems
}

// extend appends a continuation line to the field last read.
func (q *Question) extend(field, text string) {
	switch field {
	case "Asked-by":
		q.AskedBy += " " + text
	case "Entry":
		q.Entry += " " + text
	case "Provisional":
		q.Provisional += " " + text
	case "Answer":
		q.Answer += " " + text
	}
}
