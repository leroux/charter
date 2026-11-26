// Command uniq prints its input with each run of adjacent identical lines
// written once, as GNU uniq does in the C locale; -c prefixes each with the
// run's length.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"example.com/coreutils/internal/input"
	"example.com/coreutils/internal/opts"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var spec = opts.Spec{
	Flags: "c",
	Long:  map[string]byte{"count": 'c'},
}

// run of adjacent lines that compare equal: its first line and its length.
type group struct {
	line  []byte
	count int
}

// groups merges each run of adjacent lines that are byte for byte the same
// into one group. Lines equal but apart stay apart.
func groups(lines [][]byte) []group {
	var out []group
	for _, l := range lines {
		if n := len(out); n > 0 && bytes.Equal(out[n-1].line, l) {
			out[n-1].count++
			continue
		}
		out = append(out, group{line: l, count: 1})
	}
	return out
}

// format writes one group: its line, after its count right-aligned in
// seven columns and a space when counting.
func format(g group, counting bool) []byte {
	var b bytes.Buffer
	if counting {
		fmt.Fprintf(&b, "%7d ", g.count)
	}
	b.Write(g.line)
	b.WriteByte('\n')
	return b.Bytes()
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	given, names, err := opts.Parse(args, spec)
	if err == nil && len(names) > 1 {
		err = fmt.Errorf("extra operand '%s'", names[1])
	}
	if err != nil {
		fmt.Fprintf(stderr, "uniq: %v\nTry 'uniq --help' for more information.\n", err)
		return 1
	}
	counting := len(given) > 0
	srcs, err := input.Read(names, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "uniq: %v\n", err)
		return 1
	}
	for _, g := range groups(input.Lines(srcs[0].Data)) {
		stdout.Write(format(g, counting))
	}
	return 0
}
