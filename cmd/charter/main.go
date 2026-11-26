// Command charter checks and keeps a charter: the record in claims/, the
// configuration in claims/charter.toml, and the tree around them. It works
// on the charter it finds above the working directory.
//
//	charter check        the gate: every problem, one a line; exit 1 on any
//	charter gen          write the generated files
//
//charter:realizes CH-040
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/leroux/charter/internal/check"
	"github.com/leroux/charter/internal/config"
)

const usage = `usage: charter <verb> [arguments]

  check    run the gate; print every problem and exit 1 on any
  gen      write the generated files from the record
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
	switch args[0] {
	case "check":
		return runCheck(c, stdout, stderr)
	case "gen":
		if err := check.Write(c); err != nil {
			fmt.Fprintf(stderr, "charter gen: %v\n", err)
			return 1
		}
		return 0
	}
	fmt.Fprintf(stderr, "charter: unknown verb %q\n%s", args[0], usage)
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
