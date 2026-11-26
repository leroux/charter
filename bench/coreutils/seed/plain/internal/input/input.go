// Package input reads what the tools read: standard input when no file is
// named, else each named file in turn, split into lines.
package input

import (
	"bytes"
	"io"
	"os"
)

// Source is one input a tool reads: its name as given on the command line,
// empty for standard input, and its bytes.
type Source struct {
	Name string
	Data []byte
	// Regular reports whether the input is a regular file, whose size is
	// known before it is read.
	Regular bool
	Size    int64
}

// Read reads the inputs a tool was given: standard input when names is
// empty, else each named file in order. It stops at the first file it
// cannot read and returns the sources read so far with the error.
func Read(names []string, stdin io.Reader) ([]Source, error) {
	if len(names) == 0 {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return nil, err
		}
		return []Source{{Data: data}}, nil
	}
	var out []Source
	for _, name := range names {
		src, err := ReadFile(name)
		if err != nil {
			return out, err
		}
		out = append(out, src)
	}
	return out, nil
}

// ReadFile reads one named file.
func ReadFile(name string) (Source, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Source{}, err
	}
	src := Source{Name: name, Data: data}
	if fi, err := os.Stat(name); err == nil && fi.Mode().IsRegular() {
		src.Regular, src.Size = true, fi.Size()
	}
	return src, nil
}

// Lines splits data into its lines, each without its newline. A final
// line that lacks a newline is still a line; an empty input has none.
func Lines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	data = bytes.TrimSuffix(data, []byte{'\n'})
	return bytes.Split(data, []byte{'\n'})
}
