// Package cli is the charter's command line, for an adopter's own commands:
// a repository that keeps a charter writes a main that passes its gate's
// options, so the verbs run on its charter with its configuration.
//
//	package main
//
//	import (
//		"os"
//
//		"example.com/product/claims/gate"
//		"github.com/leroux/charter/cli"
//	)
//
//	func main() {
//		cli.Claimsgen(os.Args[1:], cli.Charter{Name: "product", Dir: "claims", Options: gate.Options})
//	}
//
// Claimsgen holds the verbs a person or an agent runs; Claimsmap exports a
// charter as the JSON claims-viewer reads.
package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/leroux/charter/chartertest"
)

// Claimsgen is the charter's command line:
//
//	claimsgen [-charter self]          write the generated Go and the TypeScript module
//	claimsgen [-charter self] lint [-style | -baseline]
//	                                   the gate's problems, then how many questions
//	                                   are open, then what to look at that the gate
//	                                   does not fail on, then, for each style check
//	                                   that found anything, how many it found:
//	                                   retired words, long titles, texts and
//	                                   scenarios, a Why that only points, a test two
//	                                   decisions share; -style lists each finding
//	                                   at its line instead of the counts;
//	                                   -baseline names only the style counts that
//	                                   grew past style-baseline.txt, and nothing
//	                                   when none grew. Exits 1 when it names
//	                                   anything
//	claimsgen [-charter self] questions [-all]
//	                                   the open questions of the question file,
//	                                   a block each: the id and the entry it is
//	                                   about, then its text, the provisional
//	                                   answer and who asked; then how many are
//	                                   open. -all lists the answered ones too,
//	                                   each with its answer
//	claimsgen [-charter self] next PREFIX [COUNT]
//	                                   a new id with the prefix, such as
//	                                   SW-261001-K7Q2, or COUNT of them, one a line
//	claimsgen [-charter self] anchors ID...
//	                                   JSON, per id: its document, title, status,
//	                                   scenarios, the code that realizes or
//	                                   concerns it, the tests that verify it with
//	                                   the handles each names, and its relations
//	claimsgen [-charter self] report [-md] [-go FILE]
//	                                   runs the Go tests (or reads a saved
//	                                   go test -json stream) and reports, per
//	                                   entry and scenario, the citing tests and
//	                                   whether they passed; exits 1 when an
//	                                   active scenario did not pass
//	claimsgen [-charter self] digest     the charter read through, as markdown:
//	                                   the needs by area, a map of the areas,
//	                                   each area's decisions with their why
//	claimsgen [-charter self] diff REV   what a change did to the record that
//	                                   its Changes lines do not say, against
//	                                   the area documents at the git revision:
//	                                   an entry gone, unbuilt or retired, text
//	                                   or a scenario changed without a dated
//	                                   line, a ruling lost; exits 1 on any
//	claimsgen [-charter self] mutate [-since REV] [-max N] [-seed N]
//	          [-given FILE [-given-only]] [-workers N] [-timeout D] [-skip RE]
//	          [-gocache DIR|shared] [-md] [ID...]
//	                                   breaks the named entries' realizing
//	                                   functions, or those a diff since REV
//	                                   touched, and runs the citing tests
//	                                   against each break; reports per entry
//	                                   and scenario what they killed, as JSON
//	                                   or markdown; exits 1 only when the run
//	                                   left a file in the tree
//
// Without -charter it works on the first of the charters it is given; with
// -charter NAME on the one of that name. It takes the gates' options, not
// the entries packages, so it runs even when a generated file is stale or
// missing. A fatal error is printed and exits 1.
//
//charter:realizes CH-040
func Claimsgen(args []string, charters ...Charter) {
	fs := flag.NewFlagSet("claimsgen", flag.ExitOnError)
	which := fs.String("charter", "", "the charter to work on, by name: "+names(charters)+"; empty means the first")
	_ = fs.Parse(args)
	c, err := pick(charters, *which)
	if err != nil {
		fail(err)
	}
	opts := c.Options
	args = fs.Args()
	if len(args) == 0 {
		if err := chartertest.Generate(opts); err != nil {
			fail(err)
		}
		return
	}
	switch args[0] {
	case "lint":
		code, err := runLint(os.Stdout, opts, args[1:])
		if err != nil {
			fail(err)
		}
		if code != 0 {
			os.Exit(code)
		}
	case "questions":
		qfs := flag.NewFlagSet("questions", flag.ExitOnError)
		all := qfs.Bool("all", false, "list the answered questions too, each with its answer")
		_ = qfs.Parse(args[1:])
		if qfs.NArg() > 0 {
			fail(fmt.Errorf("questions [-all] takes no argument, not %q", qfs.Arg(0)))
		}
		qs, err := chartertest.Questions(opts)
		if err != nil {
			fail(err)
		}
		writeQuestions(os.Stdout, qs, *all)
	case "next":
		prefix, count, err := nextArgs(args[1:])
		if err != nil {
			fail(err)
		}
		ids, err := chartertest.NextIDs(opts, prefix, count)
		if err != nil {
			fail(err)
		}
		for _, id := range ids {
			fmt.Println(id)
		}
	case "report":
		rfs := flag.NewFlagSet("report", flag.ExitOnError)
		md := rfs.Bool("md", false, "render markdown instead of JSON")
		goFile := rfs.String("go", "", "a saved go test -json stream to read instead of running the tests")
		_ = rfs.Parse(args[1:])
		var stream []byte
		var err error
		if *goFile != "" {
			stream, err = os.ReadFile(*goFile)
		} else {
			stream, err = chartertest.RunGoTests(opts, opts.TestPackages...)
		}
		if err != nil {
			fail(err)
		}
		rep, err := chartertest.BuildReport(opts, bytes.NewReader(stream))
		if err != nil {
			fail(err)
		}
		if *md {
			fmt.Print(rep.Markdown())
		} else {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(rep); err != nil {
				fail(err)
			}
		}
		if len(rep.Unverified) > 0 {
			os.Exit(1)
		}
	case "digest":
		d, err := chartertest.Digest(opts)
		if err != nil {
			fail(err)
		}
		fmt.Print(d)
	case "diff":
		if len(args) != 2 {
			fail(fmt.Errorf("diff REV: the git revision to read the record at, such as origin/main"))
		}
		findings, err := chartertest.History(opts, args[1])
		if err != nil {
			fail(err)
		}
		for _, f := range findings {
			fmt.Println(f)
		}
		if len(findings) > 0 {
			os.Exit(1)
		}
	case "mutate":
		if err := mutate(opts, args[1:]); err != nil {
			fail(err)
		}
	case "anchors":
		reg, err := chartertest.Read(opts)
		if err != nil {
			fail(err)
		}
		type site struct {
			File    string   `json:"file"`
			Line    int      `json:"line,omitempty"`
			Decl    string   `json:"decl,omitempty"`
			Handles []string `json:"handles,omitempty"`
		}
		type result struct {
			File      string                 `json:"file"`
			Title     string                 `json:"title"`
			Status    string                 `json:"status"`
			Kind      string                 `json:"kind,omitempty"`
			Scenarios []chartertest.Scenario `json:"scenarios"`
			Realizes  []site                 `json:"realizes"`
			Concerns  []site                 `json:"concerns"`
			Verifies  []site                 `json:"verifies"`
			Relations struct {
				Serves   []string `json:"serves"`
				See      []string `json:"see"`
				Realizes []string `json:"realizes,omitempty"`
				Concerns []string `json:"concerns,omitempty"`
				Related  []string `json:"related,omitempty"`
			} `json:"relations"`
		}
		out := map[string]*result{}
		for _, id := range args[1:] {
			var e *chartertest.Entry
			for i := range reg.Entries {
				if reg.Entries[i].ID == id {
					e = &reg.Entries[i]
				}
			}
			if e == nil || e.Level != "spec" {
				fail(fmt.Errorf("%s is not a spec id", id))
			}
			r := &result{File: e.File, Title: e.Text, Status: e.Status, Kind: e.Kind, Scenarios: e.Scenarios, Realizes: []site{}, Concerns: []site{}, Verifies: []site{}}
			if r.Scenarios == nil {
				r.Scenarios = []chartertest.Scenario{}
			}
			for _, a := range e.Anchors {
				s := site{File: a.File, Line: a.Line, Decl: a.Decl}
				if a.Kind == "realizes" {
					r.Realizes = append(r.Realizes, s)
				} else {
					r.Concerns = append(r.Concerns, s)
				}
			}
			for _, ref := range e.Refs {
				r.Verifies = append(r.Verifies, site{File: ref.File, Line: ref.Line, Decl: ref.Func, Handles: ref.Handles})
			}
			r.Relations.Serves, r.Relations.See = e.Serves, e.See
			r.Relations.Realizes, r.Relations.Concerns, r.Relations.Related = e.Realizes, e.Concerns, e.Related
			if r.Relations.Serves == nil {
				r.Relations.Serves = []string{}
			}
			if r.Relations.See == nil {
				r.Relations.See = []string{}
			}
			out[id] = r
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fail(err)
		}
	default:
		fail(fmt.Errorf("unknown command %q; see go doc github.com/leroux/charter/cli.Claimsgen", args[0]))
	}
}

// lintFlags reads the lint verb's flags: -style lists each finding of the
// style checks at its line, and -baseline holds their counts to the style
// baseline alone; the two together, or an argument, are refused.
//
//charter:realizes CH-261009-SNX2
func lintFlags(args []string) (style, baseline bool, err error) {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	s := fs.Bool("style", false, "list each finding of the style checks at its line instead of counting them")
	b := fs.Bool("baseline", false, "name only the style checks whose counts grew past style-baseline.txt in the entries directory")
	_ = fs.Parse(args)
	switch {
	case fs.NArg() > 0:
		return false, false, fmt.Errorf("lint [-style | -baseline] takes no argument, not %q", fs.Arg(0))
	case *s && *b:
		return false, false, fmt.Errorf("lint takes -style or -baseline, not both: -style lists the style findings, -baseline holds their counts to the style baseline")
	}
	return *s, *b, nil
}

// runLint runs the lint verb over the charter and prints a line for each
// thing it names: with no flag the gate's problems, the open questions, the
// advice and the style counts; with -style each style finding in place of
// its count; with -baseline the style counts that grew past the style
// baseline, and nothing else. It returns the exit status: 1 when it printed
// anything.
//
//charter:realizes CH-261009-SNX2
func runLint(w io.Writer, opts chartertest.Options, args []string) (int, error) {
	style, baseline, err := lintFlags(args)
	if err != nil {
		return 0, err
	}
	lint := chartertest.Lint
	switch {
	case style:
		lint = chartertest.LintStyle
	case baseline:
		lint = chartertest.LintBaseline
	}
	lines, err := lint(opts)
	if err != nil {
		return 0, err
	}
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	if len(lines) > 0 {
		return 1, nil
	}
	return 0, nil
}

// nextArgs reads what follows next: a prefix, then how many ids to give,
// one when no count is written.
//
//charter:realizes CH-040
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

// writeQuestions prints a block per open question, or with all per
// question: the id and the entry it is about, then, indented, the text,
// the provisional answer, who asked and, once answered, the answer; a blank
// line after each block, and last how many are open.
//
//charter:realizes CH-261009-4B3B
func writeQuestions(w io.Writer, qs []chartertest.Question, all bool) {
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

// mutate reads the mutate verb's flags, runs it and prints the result.
//
//charter:realizes CH-261005-VFJR
func mutate(opts chartertest.Options, args []string) error {
	fs := flag.NewFlagSet("mutate", flag.ExitOnError)
	since := fs.String("since", "", "a git revision: mutate only the lines changed since it, in the entries whose functions they touch")
	limit := fs.Int("max", 0, "at most this many generated mutants per entry, drawn with -seed; 0 runs all")
	seed := fs.Int64("seed", 1, "the seed -max draws with")
	given := fs.String("given", "", `a JSON file of mutants to run too: [{"entry","file","old","new","note"}]`)
	givenOnly := fs.Bool("given-only", false, "run the -given mutants alone, none generated")
	workers := fs.Int("workers", 0, "mutants run at once; 0 is the number of CPUs")
	timeout := fs.Duration("timeout", 2*time.Minute, "the bound on one go test run")
	skip := fs.String("skip", "", "passed to go test -skip")
	cache := fs.String("gocache", "", `the build cache: empty for one of the run's own, "shared" for the user's, or a directory`)
	md := fs.Bool("md", false, "render markdown instead of JSON")
	_ = fs.Parse(args)
	mo := chartertest.MutateOptions{
		IDs: fs.Args(), Since: *since, Max: *limit, Seed: *seed, Workers: *workers,
		Timeout: *timeout, Skip: *skip, Cache: *cache, Log: os.Stderr, GivenOnly: *givenOnly,
	}
	if mo.Since == "" && len(mo.IDs) == 0 && *given == "" {
		return fmt.Errorf("mutate: name entries, or -since a revision, or -given mutants")
	}
	if *given != "" {
		data, err := os.ReadFile(*given)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &mo.Given); err != nil {
			return fmt.Errorf("mutate: %s: %w", *given, err)
		}
	}
	rep, err := chartertest.Mutate(opts, mo)
	if err != nil {
		return err
	}
	if *md {
		fmt.Print(rep.Markdown())
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			return err
		}
	}
	if len(rep.Dirty) > 0 {
		return fmt.Errorf("mutate: the run left %d files in the tree: %s", len(rep.Dirty), strings.Join(rep.Dirty, ", "))
	}
	return nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "claimsgen:", err)
	os.Exit(1)
}
