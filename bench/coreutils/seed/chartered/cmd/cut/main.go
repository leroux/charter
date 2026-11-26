// Command cut prints selected fields of each line of its input, as GNU cut
// -f does in the C locale.
package main

import (
	"bytes"
	"errors"
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
	Valued: "df",
	Long:   map[string]byte{"delimiter": 'd', "fields": 'f'},
}

// span is one range of a field list, from lo to hi inclusive; hi is zero
// for a range open to the end of the line.
type span struct{ lo, hi int }

// list is a parsed field list: a field is selected when any span holds it.
type list []span

func (l list) has(n int) bool {
	for _, s := range l {
		if n >= s.lo && (s.hi == 0 || n <= s.hi) {
			return true
		}
	}
	return false
}

// parseList reads a field list: ranges N, N-M, N- and -M separated by
// commas, fields numbered from one.
//
//charter:realizes SW-261126-972N
func parseList(s string) (list, error) {
	var l list
	for _, part := range strings.Split(s, ",") {
		lo, hi, isRange := strings.Cut(part, "-")
		var sp span
		var err error
		switch {
		case !isRange:
			if sp.lo, err = field(part); err != nil {
				return nil, err
			}
			sp.hi = sp.lo
		case lo == "" && hi == "":
			return nil, errors.New("invalid range with no endpoint: -")
		default:
			sp.lo = 1
			if lo != "" {
				if sp.lo, err = field(lo); err != nil {
					return nil, err
				}
			}
			if hi != "" {
				if sp.hi, err = field(hi); err != nil {
					return nil, err
				}
				if sp.hi < sp.lo {
					return nil, errors.New("invalid decreasing range")
				}
			}
		}
		l = append(l, sp)
	}
	return l, nil
}

func field(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid field value '%s'", s)
	}
	if n == 0 {
		return 0, errors.New("fields are numbered from 1")
	}
	return n, nil
}

// fields cuts one line: the selected fields, each once and in the order
// they stand in the line whatever order the list names them in, joined by
// the delimiter. A line without the delimiter is printed whole.
//
//charter:realizes SW-261126-972N, SW-261126-7RFJ
func fields(line []byte, delim byte, l list) []byte {
	if bytes.IndexByte(line, delim) < 0 {
		return line
	}
	var out [][]byte
	for i, f := range bytes.Split(line, []byte{delim}) {
		if l.has(i + 1) {
			out = append(out, f)
		}
	}
	return bytes.Join(out, []byte{delim})
}

// config reads the options: the field list, which must be given, and the
// delimiter, a single byte, a tab when none is given.
//
//charter:realizes SW-261126-B85M
func config(given []opts.Opt) (delim byte, l list, err error) {
	delim = '\t'
	haveList := false
	for _, o := range given {
		switch o.Name {
		case 'd':
			if len(o.Value) != 1 {
				return 0, nil, errors.New("the delimiter must be a single character")
			}
			delim = o.Value[0]
		case 'f':
			if l, err = parseList(o.Value); err != nil {
				return 0, nil, err
			}
			haveList = true
		}
	}
	if !haveList {
		return 0, nil, errors.New("you must specify a list of bytes, characters, or fields")
	}
	return delim, l, nil
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	given, names, err := opts.Parse(args, spec)
	if err == nil {
		var delim byte
		var l list
		if delim, l, err = config(given); err == nil {
			srcs, err := input.Read(names, stdin)
			for _, s := range srcs {
				for _, ln := range input.Lines(s.Data) {
					stdout.Write(append(fields(ln, delim, l), '\n'))
				}
			}
			if err != nil {
				fmt.Fprintf(stderr, "cut: %v\n", err)
				return 1
			}
			return 0
		}
	}
	fmt.Fprintf(stderr, "cut: %v\nTry 'cut --help' for more information.\n", err)
	return 1
}
