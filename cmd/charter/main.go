// Command charter checks and keeps a charter: the record in claims/, the
// configuration in claims/charter.toml, and the tree around them. It works
// on the charter it finds above the working directory.
//
//charter:realizes CH-040
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/diff"
	"github.com/leroux/charter/internal/export"
	"github.com/leroux/charter/internal/lint"
	"github.com/leroux/charter/internal/record"
	"github.com/leroux/charter/internal/report"
)

const usage = `usage: charter <verb> [arguments]

  check               run the gate; print every problem and exit 1 on any
  gen                 write the generated files and the digest from the record
  lint [-style]       the gate's problems, then advice; -style lists the style findings
  next PREFIX [N]     print N new ids with the prefix
  anchors ID...       an entry's code, tests and relations, as JSON
  questions [-all]    the open questions; -all the answered ones too
  report [-md]        run every language's tests and say which scenarios passed
  diff REV            name what the record changed since REV with no Changes line
  export [-o FILE]    the charter as one JSON file for the viewer
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
	case "lint":
		fs := flag.NewFlagSet("lint", flag.ContinueOnError)
		fs.SetOutput(stderr)
		style := fs.Bool("style", false, "list every style finding at its line")
		if err := fs.Parse(rest); err != nil {
			return 2
		}
		out, err := lint.Lint(c, *style)
		if err != nil {
			return fail(err)
		}
		return lines(out, true)
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
		return lines(findings, true)
	case "export":
		fs := flag.NewFlagSet("export", flag.ContinueOnError)
		fs.SetOutput(stderr)
		out := fs.String("o", "-", "the file to write, or - for standard output")
		if err := fs.Parse(rest); err != nil {
			return 2
		}
		t, problems, err := check.Run(c)
		if err != nil {
			return fail(err)
		}
		for _, p := range problems {
			fmt.Fprintf(stderr, "charter export: warning: %s\n", p)
		}
		var buf bytes.Buffer
		if err := export.Write(&buf, t, problems, source(c.Root), c.Root); err != nil {
			return fail(err)
		}
		if *out == "-" {
			_, err = stdout.Write(buf.Bytes())
		} else {
			err = os.WriteFile(*out, buf.Bytes(), 0o644)
		}
		if err != nil {
			return fail(err)
		}
		return 0
	}
	fmt.Fprintf(stderr, "charter: unknown verb %q\n%s", verb, usage)
	return 2
}

// runCheck prints each problem on a line of its own, sorted, and exits 1
// when it prints any; on a clean tree it prints nothing.
//
//charter:realizes CH-037
func runCheck(c *config.Config, stdout, stderr io.Writer) int {
	problems, err := check.Problems(c)
	if err != nil {
		fmt.Fprintf(stderr, "charter check: %v\n", err)
		return 2
	}
	for _, p := range problems {
		fmt.Fprintln(stdout, p)
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
// the tests that verify it with the handles each names, and its relations.
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
	}
	out := map[string]*result{}
	for _, id := range ids {
		e := t.Record.Entry(id)
		if e == nil {
			fmt.Fprintf(stderr, "charter anchors: %s is not an entry of the record\n", id)
			return 2
		}
		r := &result{File: e.File, Title: e.Title, Status: e.Status, Kind: e.Kind, Scenarios: e.Scenarios,
			Realizes: []site{}, Concerns: []site{}, Verifies: []site{}, Serves: e.Serves, See: e.See}
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
