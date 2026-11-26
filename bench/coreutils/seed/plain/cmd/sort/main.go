// Command sort prints the lines of its input in order, as GNU sort does in
// the C locale: by bytes, or by number with -n, reversed with -r.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"

	"example.com/coreutils/internal/input"
	"example.com/coreutils/internal/opts"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var spec = opts.Spec{
	Flags: "nr",
	Long:  map[string]byte{"numeric-sort": 'n', "reverse": 'r'},
}

// order is how lines compare: by number or by bytes, and whether reversed.
type order struct {
	numeric, reverse bool
}

// compare orders two lines. By bytes, the lines compare as unsigned byte
// strings. By number, the leading numbers compare and lines whose numbers
// are equal fall back to the byte comparison. Reversal reverses the whole
// comparison, the fallback with it.
func (o order) compare(a, b []byte) int {
	c := 0
	if o.numeric {
		c = numcompare(a, b)
	}
	if c == 0 {
		c = bytes.Compare(a, b)
	}
	if o.reverse {
		c = -c
	}
	return c
}

// number is a line's leading number, as -n reads it: its sign, its
// integer digits without leading zeros and its fraction digits without
// trailing zeros.
type number struct {
	neg      bool
	whole    []byte
	fraction []byte
}

// parseNumber reads a line's leading number: blanks, an optional minus,
// digits, and an optional point with more digits. A line that does not
// start with one reads as zero, and so does negative zero.
func parseNumber(s []byte) number {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	var n number
	if i < len(s) && s[i] == '-' {
		n.neg = true
		i++
	}
	start := i
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	n.whole = bytes.TrimLeft(s[start:i], "0")
	if i < len(s) && s[i] == '.' {
		i++
		start = i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		n.fraction = bytes.TrimRight(s[start:i], "0")
	}
	if len(n.whole) == 0 && len(n.fraction) == 0 {
		n.neg = false
	}
	return n
}

// numcompare compares the leading numbers of two lines exactly, digit by
// digit, however long they are.
func numcompare(a, b []byte) int {
	x, y := parseNumber(a), parseNumber(b)
	if x.neg != y.neg {
		if x.neg {
			return -1
		}
		return 1
	}
	c := magnitude(x, y)
	if x.neg {
		c = -c
	}
	return c
}

func magnitude(x, y number) int {
	if len(x.whole) != len(y.whole) {
		if len(x.whole) < len(y.whole) {
			return -1
		}
		return 1
	}
	if c := bytes.Compare(x.whole, y.whole); c != 0 {
		return c
	}
	return bytes.Compare(x.fraction, y.fraction)
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	given, names, err := opts.Parse(args, spec)
	if err != nil {
		fmt.Fprintf(stderr, "sort: %v\nTry 'sort --help' for more information.\n", err)
		return 2
	}
	var o order
	for _, g := range given {
		switch g.Name {
		case 'n':
			o.numeric = true
		case 'r':
			o.reverse = true
		}
	}
	srcs, err := input.Read(names, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "sort: %v\n", err)
		return 2
	}
	var lines [][]byte
	for _, s := range srcs {
		lines = append(lines, input.Lines(s.Data)...)
	}
	slices.SortStableFunc(lines, o.compare)
	for _, l := range lines {
		stdout.Write(append(l, '\n'))
	}
	return 0
}
