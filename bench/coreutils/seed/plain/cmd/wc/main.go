// Command wc prints the newline, word and byte counts of its input, as GNU
// wc does in the C locale.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"example.com/coreutils/internal/input"
	"example.com/coreutils/internal/opts"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var spec = opts.Spec{
	Flags: "lwc",
	Long:  map[string]byte{"lines": 'l', "words": 'w', "bytes": 'c'},
}

// counts are what wc counts in one input.
type counts struct {
	lines, words, bytes int
}

// which counts to print, in the fixed order lines, words, bytes.
type which struct {
	lines, words, bytes bool
}

func (w which) n() int {
	n := 0
	for _, b := range []bool{w.lines, w.words, w.bytes} {
		if b {
			n++
		}
	}
	return n
}

// selection reads which counts the options ask for: those named, or all
// three when none is.
func selection(given []opts.Opt) which {
	var w which
	for _, o := range given {
		switch o.Name {
		case 'l':
			w.lines = true
		case 'w':
			w.words = true
		case 'c':
			w.bytes = true
		}
	}
	if w.n() == 0 {
		w = which{true, true, true}
	}
	return w
}

// count counts data's newlines, words and bytes. A word is a run of
// bytes that holds at least one printable byte, delimited by white space;
// a byte that is neither printable nor white space neither starts nor
// ends one.
func count(data []byte) counts {
	c := counts{bytes: len(data)}
	inWord := false
	for _, b := range data {
		switch {
		case b == '\n':
			c.lines++
			inWord = false
		case b == ' ' || b == '\t' || b == '\v' || b == '\f' || b == '\r':
			inWord = false
		case b >= 0x21 && b <= 0x7e:
			if !inWord {
				c.words++
				inWord = true
			}
		}
	}
	return c
}

// width is the width every count is right-aligned to: one when a single
// count of a single input is printed; seven when an input is standard
// input or not a regular file; else the digits of the inputs' total size.
func width(w which, srcs []input.Source) int {
	if w.n() == 1 && len(srcs) == 1 {
		return 1
	}
	min := 1
	var total int64
	for _, s := range srcs {
		if s.Regular {
			total += s.Size
		} else {
			min = 7
		}
	}
	n := len(strconv.FormatInt(total, 10))
	if n < min {
		n = min
	}
	return n
}

// line formats one input's counts: the selected counts in the order
// lines, words, bytes, each right-aligned to the width and separated by a
// space, then the input's name when it has one.
func line(c counts, w which, wid int, name string) string {
	var parts []string
	for _, f := range []struct {
		on bool
		n  int
	}{{w.lines, c.lines}, {w.words, c.words}, {w.bytes, c.bytes}} {
		if f.on {
			parts = append(parts, fmt.Sprintf("%*d", wid, f.n))
		}
	}
	if name != "" {
		parts = append(parts, name)
	}
	return strings.Join(parts, " ") + "\n"
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	given, names, err := opts.Parse(args, spec)
	if err != nil {
		fmt.Fprintf(stderr, "wc: %v\nTry 'wc --help' for more information.\n", err)
		return 1
	}
	w := selection(given)
	srcs, err := input.Read(names, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "wc: %v\n", err)
		return 1
	}
	wid := width(w, srcs)
	for _, s := range srcs {
		io.WriteString(stdout, line(count(s.Data), w, wid, s.Name))
	}
	return 0
}
