package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// wc runs the command as a script would, in a fresh directory holding
// files, and returns what it printed and its exit status.
func wc(t *testing.T, files map[string]string, stdin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	t.Chdir(t.TempDir())
	for name, data := range files {
		if err := os.WriteFile(name, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return out.String(), errOut.String(), code
}

// counted is wc's standard output, failing the test unless it exits 0.
func counted(t *testing.T, files map[string]string, stdin string, args ...string) string {
	t.Helper()
	out, stderr, code := wc(t, files, stdin, args...)
	if code != 0 {
		t.Fatalf("wc %q: exit %d: %s", args, code, stderr)
	}
	return out
}

func TestWcPrintsAFileLikeGNU(t *testing.T) {
	// GNU coreutils 9.4, LC_ALL=C: wc n1.txt
	got := counted(t, map[string]string{"n1.txt": "one two\nthree\tfour five\n\nsix"}, "", "n1.txt")
	if want := " 3  6 28 n1.txt\n"; got != want {
		t.Errorf("wc n1.txt = %q, want GNU's %q", got, want)
	}
}

func TestWcReadsStandardInputWhenNoFileIsNamed(t *testing.T) {
	if got := counted(t, nil, "a\nb\nc\n", "-l"); got != "3\n" {
		t.Errorf("wc -l on standard input = %q, want 3", got)
	}
}

func TestWcRefusesAnUnknownFlag(t *testing.T) {
	out, stderr, code := wc(t, nil, "a b\n", "-Z")
	if code == 0 || out != "" || !strings.Contains(stderr, "'Z'") {
		t.Errorf("wc -Z: exit %d, stdout %q, stderr %q; want a failure, no output, Z named", code, out, stderr)
	}
}

func TestLinesAreNewlines(t *testing.T) {
	if c := count([]byte("a\nb")); c.lines != 1 {
		t.Errorf("lines in \"a\\nb\" = %d, want 1", c.lines)
	}
}

func TestWordsHoldAPrintableByte(t *testing.T) {
	for in, want := range map[string]int{
		"\x01":         0,
		"a\x01b c":     2,
		"\x01a":        1,
		"x\ty\rz\vw\f": 4,
		"\xff\x80":     0,
	} {
		if c := count([]byte(in)); c.words != want {
			t.Errorf("words in %q = %d, want %d", in, c.words, want)
		}
	}
}

func TestBytesAreTheLength(t *testing.T) {
	if c := count([]byte("é\n\xff\n")); c.bytes != 5 {
		t.Errorf("bytes = %d, want 5", c.bytes)
	}
}

func TestNoOptionPrintsAllThree(t *testing.T) {
	if got := counted(t, nil, "a b\nc\n"); got != "      2       3       6\n" {
		t.Errorf("wc = %q, want lines, words, bytes", got)
	}
}

func TestTheOrderIsFixed(t *testing.T) {
	if got := counted(t, nil, "a b\nc\n", "-cl"); got != "      2       6\n" {
		t.Errorf("wc -cl = %q, want the lines then the bytes", got)
	}
}

func TestASingleCountIsBare(t *testing.T) {
	if got := counted(t, nil, strings.Repeat("x\n", 12), "-l"); got != "12\n" {
		t.Errorf("wc -l = %q, want 12 unpadded", got)
	}
}

func TestStandardInputPadsToSeven(t *testing.T) {
	if got := counted(t, nil, "x", "-wc"); got != "      1       1\n" {
		t.Errorf("wc -wc = %q, want seven columns", got)
	}
}

func TestAFilePadsToItsSizesDigits(t *testing.T) {
	got := counted(t, map[string]string{"f": "hello world\n"}, "", "f")
	if got != " 1  2 12 f\n" {
		t.Errorf("wc f = %q, want two columns", got)
	}
}

func TestTheNameFollowsTheCounts(t *testing.T) {
	if got := counted(t, map[string]string{"notes.txt": "a\nb\n"}, "", "-l", "notes.txt"); got != "2 notes.txt\n" {
		t.Errorf("wc -l notes.txt = %q", got)
	}
}
