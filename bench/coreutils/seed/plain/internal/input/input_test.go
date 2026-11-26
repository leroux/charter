package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoNamesReadsStandardInput(t *testing.T) {
	srcs, err := Read(nil, strings.NewReader("piped\n"))
	if err != nil || len(srcs) != 1 || string(srcs[0].Data) != "piped\n" || srcs[0].Name != "" {
		t.Fatalf("Read(nil) = %+v, %v; want standard input's bytes, unnamed", srcs, err)
	}
}

func TestNamedFilesAreReadInOrder(t *testing.T) {
	dir := t.TempDir()
	a, b := filepath.Join(dir, "a"), filepath.Join(dir, "b")
	os.WriteFile(a, []byte("first\n"), 0o644)
	os.WriteFile(b, []byte("second\n"), 0o644)
	srcs, err := Read([]string{b, a}, strings.NewReader("ignored"))
	if err != nil || len(srcs) != 2 {
		t.Fatalf("Read = %+v, %v", srcs, err)
	}
	if string(srcs[0].Data) != "second\n" || string(srcs[1].Data) != "first\n" {
		t.Errorf("inputs %q, %q; want b's then a's", srcs[0].Data, srcs[1].Data)
	}
	if !srcs[0].Regular || srcs[0].Size != 7 {
		t.Errorf("b is regular %v of size %d; want a regular file of 7 bytes", srcs[0].Regular, srcs[0].Size)
	}
}

func lines(s string) []string {
	var out []string
	for _, l := range Lines([]byte(s)) {
		out = append(out, string(l))
	}
	return out
}

func TestAnUnterminatedLastLineIsALine(t *testing.T) {
	for _, in := range []string{"a\nb", "a\nb\n"} {
		if got := strings.Join(lines(in), "|"); got != "a|b" {
			t.Errorf("Lines(%q) = %q, want a|b", in, got)
		}
	}
}

func TestAnEmptyInputHasNoLines(t *testing.T) {
	if got := lines(""); len(got) != 0 {
		t.Errorf("Lines(\"\") = %q, want none", got)
	}
	if got := lines("\n"); len(got) != 1 || got[0] != "" {
		t.Errorf("Lines(\"\\n\") = %q, want one empty line", got)
	}
}
