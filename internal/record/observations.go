package record

import (
	"fmt"
	"regexp"
	"strings"
)

// Observation is a fact about a system the code does not own, that entries
// rest on:
//
//	## LockIsAdvisory
//
//	Subject: os
//	Status: documented
//	Source: https://man7.org/linux/man-pages/man2/flock.2.html
//	Supports: CON-001
//
//	flock locks are advisory: a process that does not ask is not stopped.
type Observation struct {
	// Name is an upper-camel word, unique among the observations; it is the
	// symbol each language's generated file declares the observation by.
	Name string
	// Text is what was seen or read.
	Text    string
	Subject string
	// Status is observed, documented, broken or unsure.
	Status   string
	Observed string // the date a probe last saw it true
	Source   string // where it was read
	Supports []string
	File     string
	Line     int
}

// Pos is the observation's heading as file:line.
func (o *Observation) Pos() string { return fmt.Sprintf("%s:%d", o.File, o.Line) }

var (
	observationHeading = regexp.MustCompile(`^## (\S+)\s*$`)
	observationName    = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
)

// observationStatuses are the states of a fact.
var observationStatuses = map[string]bool{"observed": true, "documented": true, "broken": true, "unsure": true}

// parseObservations reads observations.md: one heading per observation, its
// field lines, then its text. An unknown field, a field twice, a status
// outside the four and a name declared twice are form problems.
//
//charter:realizes CH-261126-HK6S
func parseObservations(rel, src string) ([]*Observation, []string) {
	var out []*Observation
	var problems []string
	bad := func(line int, format string, args ...any) {
		problems = append(problems, fmt.Sprintf("%s:%d: form: %s", rel, line, fmt.Sprintf(format, args...)))
	}
	seen := map[string]int{}
	var cur *Observation
	var fields map[string]bool
	var text []string
	inFields := false
	finish := func() {
		if cur == nil {
			return
		}
		cur.Text = strings.TrimSpace(strings.Join(text, " "))
		if !fields["Status"] {
			bad(cur.Line, "%s has no Status; it is observed, documented, broken or unsure", cur.Name)
		}
		out = append(out, cur)
		cur, text = nil, nil
	}
	for i, l := range strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n") {
		n := i + 1
		line := strings.TrimRight(l, " \t")
		if strings.HasPrefix(line, "## ") {
			finish()
			m := observationHeading.FindStringSubmatch(line)
			if m == nil || !observationName.MatchString(m[1]) {
				bad(n, "an observation's heading is \"## Name\", the name an upper-camel word such as LockIsAdvisory")
				continue
			}
			if first, ok := seen[m[1]]; ok {
				bad(n, "%s is declared twice, at lines %d and %d", m[1], first, n)
			}
			seen[m[1]] = n
			cur = &Observation{Name: m[1], File: rel, Line: n}
			fields = map[string]bool{}
			inFields = true
			continue
		}
		if cur == nil {
			continue
		}
		if inFields {
			if line == "" {
				if len(fields) > 0 {
					inFields = false
				}
				continue
			}
			if m := fieldLine.FindStringSubmatch(line); m != nil {
				name, value := m[1], strings.TrimSpace(m[2])
				if fields[name] {
					bad(n, "%s appears twice", name)
					continue
				}
				fields[name] = true
				switch name {
				case "Subject":
					cur.Subject = value
				case "Status":
					if !observationStatuses[value] {
						bad(n, "status %q is not observed, documented, broken or unsure", value)
					}
					cur.Status = value
				case "Observed":
					cur.Observed = value
				case "Source":
					cur.Source = value
				case "Supports":
					cur.Supports = splitList(value)
				default:
					bad(n, "%q is not a field of an observation; the fields are Subject, Status, Observed, Source and Supports", name)
				}
				continue
			}
			inFields = false
		}
		if line != "" {
			text = append(text, strings.TrimSpace(line))
		}
	}
	finish()
	return out, problems
}
