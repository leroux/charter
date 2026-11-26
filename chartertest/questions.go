package chartertest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// The question file, questions.md beside the area documents, holds what
// nobody ruled: an agent that had to choose something no entry decides
// writes the question, who asked, the entry it is about and the
// provisional answer it took, and goes on. A person answers in place,
// opening the Answer with the date; the question stays, answered, as the
// record of the ruling. Lint counts the open ones in one line, and the
// questions verb of claimsgen lists them.
//
//	## Q-261009-M4TK Does a reset go on to the keys under a removed one?
//	Asked-by: agent, PR #452
//	Entry: SW-221/reset-passes-on
//	Provisional: yes, as SW-226 says "every key under it"
//	Answer: 2026-10-09, owner: yes

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

const questionsFile = "questions.md"

// policyFile is the owner's policy beside the area documents: by class,
// where a provisional scenario may land and whether a verdict of the judge
// blocks. People and the judge read it; the gate does not parse it.
const policyFile = "policy.md"

// policyPath is the policy file, slash-relative to the root.
func (r *resolved) policyPath() string { return path.Join(r.entriesDir, policyFile) }

// glossaryFile is the record's vocabulary beside the area documents: each
// term defined once, at the boundary, with the entry that rules it. People
// read it; the gate does not parse it.
const glossaryFile = "glossary.md"

// glossaryPath is the glossary, slash-relative to the root.
func (r *resolved) glossaryPath() string { return path.Join(r.entriesDir, glossaryFile) }

// auditFile is the audit's report beside the area documents: what the
// audit found where the record, its tests and its code disagree. People
// read it; the gate does not parse it.
const auditFile = "audit.md"

// auditPath is the audit's report, slash-relative to the root.
func (r *resolved) auditPath() string { return path.Join(r.entriesDir, auditFile) }

// notAreaDocuments are the markdown files beside the area documents that
// the reader skips, README.md aside: the digest, the question file, the
// policy, the glossary and the audit's report.
//
//charter:realizes CH-010
func (r *resolved) notAreaDocuments() []string {
	return []string{r.digestPath(), r.questionsPath(), r.policyPath(), r.glossaryPath(), r.auditPath()}
}

var (
	questionHeading = regexp.MustCompile(`^## (Q-\S+) (\S.*)$`)
	questionField   = regexp.MustCompile(`^([A-Z][a-z]+(?:-[a-z]+)?):(?: (.*))?$`)
	answerOpening   = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}), (\S.*)$`)
)

// questionsPath is the question file, slash-relative to the root.
func (r *resolved) questionsPath() string { return path.Join(r.entriesDir, questionsFile) }

// Questions reads the question file of the charter the options name with
// the gate's reader, and returns every question in the order the file
// holds them, open and answered; a missing file is no questions. The
// entries bear only on the form problems, which are the gate's to report
// and lint's, so it reads none, and it works on a tree the gate fails.
//
//charter:realizes CH-261009-4B3B
func Questions(opts Options) ([]Question, error) {
	r, err := opts.resolve(callerDir(1))
	if err != nil {
		return nil, err
	}
	qs, _, err := readQuestions(r, newEntriesPackage())
	return qs, err
}

// readQuestions parses the question file against the entries, reporting
// each form problem with its file and line. A missing file is no
// questions.
//
//charter:realizes CH-261009-4B3B
func readQuestions(r *resolved, pkg *entriesPackage) ([]Question, []string, error) {
	rel := r.questionsPath()
	data, err := os.ReadFile(filepath.Join(r.Root, filepath.FromSlash(rel)))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	return parseQuestions(rel, string(data), pkg)
}

// parseQuestions reads the question file's text.
//
//charter:realizes CH-261009-4B3B
func parseQuestions(rel, src string, pkg *entriesPackage) ([]Question, []string, error) {
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
			e := pkg.byID[id]
			switch {
			case e == nil || e.Level != levelSpec:
				bad(cur.Line, "%s is about %s, which is not an entry", cur.ID, id)
			case handle != "" && !e.hasScenario(handle):
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
				bad(n, "a question's heading is \"## Q-ID question\", the id from claimsgen next Q")
				continue
			}
			parts, problem := parseID(m[1])
			if problem != "" {
				bad(n, "%s", problem)
				continue
			}
			if parts.Lens != "Q" {
				bad(n, "a question's id has the prefix Q, not %s", parts.Lens)
			}
			if first, ok := ids[m[1]]; ok {
				bad(n, "%s is declared twice, at lines %d and %d; %s", m[1], first, n, takeNewID(&entry{Lens: "Q"}))
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
	return qs, problems, nil
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

// lintQuestions says in one line, at the question file's first line, how
// many questions are open for the owner to answer, and names the command
// that lists them; it says nothing when none is open.
//
//charter:realizes CH-261009-4B3B
func lintQuestions(s *state) []string {
	open := 0
	for _, q := range s.questions {
		if q.Open() {
			open++
		}
	}
	switch open {
	case 0:
		return nil
	case 1:
		return []string{fmt.Sprintf("%s:1: lint: 1 question is open; go run ./cmd/claimsgen questions lists it", s.r.questionsPath())}
	}
	return []string{fmt.Sprintf("%s:1: lint: %d questions are open; go run ./cmd/claimsgen questions lists them", s.r.questionsPath(), open)}
}
