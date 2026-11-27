// Command charter checks and keeps a charter: the record in claims/, the
// configuration in claims/charter.toml, and the tree around them. It works
// on the charter it finds above the working directory.
//
//charter:realizes CH-040
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/diff"
	"github.com/leroux/charter/internal/docs"
	"github.com/leroux/charter/internal/record"
	"github.com/leroux/charter/internal/report"
)

const usage = `usage: charter <verb> [arguments]

  check               run the gate; print every problem and exit 1 on any, warnings on stderr
  gen                 write the generated files and the digest from the record
  next PREFIX [N]     print N new ids with the prefix
  anchors ID...       an entry's code, tests and relations, as JSON
  questions [-all]    the open questions; -all the answered ones too
  report [-md]        run every language's tests and say which scenarios passed
  diff REV            name what the record changed since REV with no Changes line
  curate [-by NAME] FILE...
                      record that a person read each doc, with a hash of what it references
`

func main() {
	os.Exit(run(os.Args[1:], ".", os.Stdout, os.Stderr))
}

// run runs one verb from the directory dir and returns the exit status.
func run(args []string, dir string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	c, err := config.Find(dir)
	if err != nil {
		fmt.Fprintf(stderr, "charter: %v\n", err)
		return 2
	}
	verb, rest := args[0], args[1:]
	fail := func(err error) int {
		fmt.Fprintf(stderr, "charter %s: %v\n", verb, err)
		return 2
	}
	lines := func(ls []string, failOnAny bool) int {
		for _, l := range ls {
			fmt.Fprintln(stdout, l)
		}
		if failOnAny && len(ls) > 0 {
			return 1
		}
		return 0
	}
	switch verb {
	case "check":
		return runCheck(c, stdout, stderr)
	case "gen":
		if err := check.Write(c); err != nil {
			return fail(err)
		}
		return 0
	case "next":
		prefix, count, err := nextArgs(rest)
		if err != nil {
			return fail(err)
		}
		rec, err := record.Read(c.Root, config.Dir)
		if err != nil {
			return fail(err)
		}
		ids, err := rec.NewIDs(prefix, count, time.Now(), rand.IntN)
		if err != nil {
			return fail(err)
		}
		return lines(ids, false)
	case "anchors":
		return anchors(c, rest, stdout, stderr)
	case "questions":
		fs := flag.NewFlagSet("questions", flag.ContinueOnError)
		fs.SetOutput(stderr)
		all := fs.Bool("all", false, "list the answered questions too, each with its answer")
		if err := fs.Parse(rest); err != nil {
			return 2
		}
		rec, err := record.Read(c.Root, config.Dir)
		if err != nil {
			return fail(err)
		}
		writeQuestions(stdout, rec.Questions, *all)
		return 0
	case "report":
		fs := flag.NewFlagSet("report", flag.ContinueOnError)
		fs.SetOutput(stderr)
		md := fs.Bool("md", false, "render markdown instead of JSON")
		if err := fs.Parse(rest); err != nil {
			return 2
		}
		rep, err := report.Run(c)
		if err != nil {
			return fail(err)
		}
		if *md {
			fmt.Fprint(stdout, rep.Markdown())
		} else {
			enc := json.NewEncoder(stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(rep); err != nil {
				return fail(err)
			}
		}
		if len(rep.Unverified) > 0 {
			return 1
		}
		return 0
	case "diff":
		if len(rest) != 1 {
			return fail(fmt.Errorf("diff REV: the git revision to read the record at, such as origin/main"))
		}
		findings, err := diff.Diff(c, rest[0])
		if err != nil {
			return fail(err)
		}
		status := lines(findings, true)
		if len(c.Docs) > 0 {
			// The docs to read again are a list, never a finding: when they
			// cannot be listed, diff says so and its status stands.
			again, err := docsToReadAgain(c, rest[0])
			if err != nil {
				fmt.Fprintf(stderr, "charter diff: the docs to read again are not listed: %v\n", err)
			}
			lines(again, false)
		}
		return status
	case "curate":
		return curate(c, rest, stdout, stderr)
	}
	fmt.Fprintf(stderr, "charter: unknown verb %q\n%s", verb, usage)
	return 2
}

// runCheck prints each problem on a line of its own, sorted, and exits 1
// when it prints any; on a clean tree it prints nothing. Its warnings go
// to standard error and never change its exit.
//
//charter:realizes CH-037, CH-261127-QRQS
func runCheck(c *config.Config, stdout, stderr io.Writer) int {
	t, problems, err := check.Run(c)
	if err != nil {
		fmt.Fprintf(stderr, "charter check: %v\n", err)
		return 2
	}
	for _, p := range problems {
		fmt.Fprintln(stdout, p)
	}
	for _, w := range check.Warnings(t) {
		fmt.Fprintln(stderr, w)
	}
	if len(problems) > 0 {
		return 1
	}
	return 0
}

func nextArgs(args []string) (prefix string, count int, err error) {
	switch len(args) {
	case 1:
		return args[0], 1, nil
	case 2:
		count, err := strconv.Atoi(args[1])
		if err != nil || count < 1 {
			return "", 0, fmt.Errorf("next PREFIX [COUNT]: the count is a number above zero, not %q", args[1])
		}
		return args[0], count, nil
	}
	return "", 0, fmt.Errorf("next PREFIX [COUNT]")
}

// anchors prints, per id, the entry's document, title, status, kind and
// scenarios, the code that realizes or concerns it with the declarations,
// the tests that verify it with the handles each names, its relations, the
// docs that mention it and its changes.
//
//charter:realizes CH-040
func anchors(c *config.Config, ids []string, stdout, stderr io.Writer) int {
	t, err := check.Load(c)
	if err != nil {
		fmt.Fprintf(stderr, "charter anchors: %v\n", err)
		return 2
	}
	type site struct {
		File    string   `json:"file"`
		Line    int      `json:"line"`
		Decl    string   `json:"decl,omitempty"`
		Test    string   `json:"test,omitempty"`
		Handles []string `json:"handles,omitempty"`
	}
	type result struct {
		File      string            `json:"file"`
		Title     string            `json:"title"`
		Status    string            `json:"status"`
		Kind      string            `json:"kind,omitempty"`
		Scenarios []record.Scenario `json:"scenarios"`
		Realizes  []site            `json:"realizes"`
		Concerns  []site            `json:"concerns"`
		Verifies  []site            `json:"verifies"`
		Serves    []string          `json:"serves"`
		See       []string          `json:"see"`
		Docs      []docSite         `json:"docs"`
		Changes   []record.Change   `json:"changes"`
	}
	out := map[string]*result{}
	for _, id := range ids {
		e := t.Record.Entry(id)
		if e == nil {
			fmt.Fprintf(stderr, "charter anchors: %s is not an entry of the record\n", id)
			return 2
		}
		r := &result{File: e.File, Title: e.Title, Status: e.Status, Kind: e.Kind, Scenarios: e.Scenarios,
			Realizes: []site{}, Concerns: []site{}, Verifies: []site{}, Serves: e.Serves, See: e.See, Docs: docMentions(t, id), Changes: e.Changes}
		for _, a := range t.Facts.Anchors {
			for _, aid := range a.IDs {
				if aid != id {
					continue
				}
				s := site{File: a.File, Line: a.Line, Decl: a.Decl}
				if a.Kind == "realizes" {
					r.Realizes = append(r.Realizes, s)
				} else {
					r.Concerns = append(r.Concerns, s)
				}
			}
		}
		for _, ct := range t.Facts.Citations {
			if ct.Symbol == e.Symbol && !ct.Observes {
				r.Verifies = append(r.Verifies, site{File: ct.File, Line: ct.Line, Test: ct.Test, Handles: ct.Handles})
			}
		}
		out[id] = r
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(stderr, "charter anchors: %v\n", err)
		return 2
	}
	return 0
}

// docsToReadAgain is diff's list of the docs whose references changed
// since the revision.
//
//charter:realizes CH-261127-8MYR
func docsToReadAgain(c *config.Config, rev string) ([]string, error) {
	t, err := check.Load(c)
	if err != nil {
		return nil, err
	}
	return diff.Docs(c, rev, t)
}

// docSite is a doc's mention of an entry.
type docSite struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

// docMentions are the docs that reference or embed the entry, at the line
// of each mention.
//
//charter:realizes CH-261127-8MYR
func docMentions(t *check.Tree, id string) []docSite {
	out := []docSite{}
	for _, d := range t.Docs {
		for _, r := range d.EntryRefs {
			if r.ID == id {
				out = append(out, docSite{File: d.File, Line: r.Line})
			}
		}
		for _, e := range d.Embeds {
			if e.ID == id {
				out = append(out, docSite{File: d.File, Line: e.Line})
			}
		}
	}
	return out
}

// curate records, at the end of each doc named, who curated it, the date
// and the hash of everything it references. A doc whose references do not
// resolve, or that charter.toml does not name, is refused and not written.
//
//charter:realizes CH-261127-GXP7
func curate(c *config.Config, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("curate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	by := fs.String("by", "owner", "who read the doc against the record and the code")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "charter curate: name the docs to curate")
		return 2
	}
	t, err := check.Load(c)
	if err != nil {
		fmt.Fprintf(stderr, "charter curate: %v\n", err)
		return 2
	}
	byFile := map[string]*docs.Doc{}
	for _, d := range t.Docs {
		byFile[d.File] = d
	}
	date := time.Now().UTC().Format("2006-01-02")
	status := 0
	for _, arg := range fs.Args() {
		rel, err := docPath(c.Root, arg)
		d := byFile[rel]
		if err != nil || d == nil {
			fmt.Fprintf(stderr, "charter curate: %s is not a doc charter.toml names\n", arg)
			status = 1
			continue
		}
		if ps := docs.Problems([]*docs.Doc{d}, t.Record, t.Code); len(ps) > 0 {
			for _, p := range ps {
				fmt.Fprintln(stderr, p)
			}
			fmt.Fprintf(stderr, "charter curate: %s has references that do not resolve; nothing written\n", rel)
			status = 1
			continue
		}
		out := docs.Curate(d, *by, date, docs.Hashes(d, t.Record, t.Code))
		if err := os.WriteFile(filepath.Join(c.Root, filepath.FromSlash(rel)), []byte(out), 0o644); err != nil {
			fmt.Fprintf(stderr, "charter curate: %v\n", err)
			return 2
		}
		fmt.Fprintf(stdout, "%s curated by %s on %s\n", rel, *by, date)
	}
	return status
}

// docPath is a path given on the command line, relative to the working
// directory, as a slash path relative to the root.
func docPath(root, arg string) (string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

// writeQuestions prints a block per open question, or with all per
// question: the id and the entry it is about, then, indented, the text,
// the provisional answer, who asked and, once answered, the answer; a blank
// line after each block, and last how many are open.
//
//charter:realizes CH-261009-4B3B
func writeQuestions(w io.Writer, qs []record.Question, all bool) {
	open := 0
	for _, q := range qs {
		if q.Open() {
			open++
		} else if !all {
			continue
		}
		fmt.Fprintf(w, "%s  %s\n  %s\n  provisional: %s\n  asked by: %s\n", q.ID, q.Entry, q.Text, q.Provisional, q.AskedBy)
		if !q.Open() {
			fmt.Fprintf(w, "  answer: %s\n", q.Answer)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintf(w, "%d open\n", open)
}

// source names the checkout the charter was read from, such as
// "claims @ main 1a2b3c4", or the directory alone outside git.
func source(root string) string {
	git := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		b, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(b))
	}
	sha := git("rev-parse", "--short", "HEAD")
	if sha == "" {
		return config.Dir
	}
	s := config.Dir + " @ "
	if b := git("branch", "--show-current"); b != "" {
		s += b + " "
	}
	s += sha
	if git("status", "--porcelain", "--", config.Dir) != "" {
		s += " + uncommitted changes"
	}
	return s
}
