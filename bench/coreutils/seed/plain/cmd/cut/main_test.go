package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// cut runs the command as a script would, in a fresh directory holding
// files, and returns what it printed and its exit status.
func cut(t *testing.T, files map[string]string, stdin string, args ...string) (stdout, stderr string, code int) {
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

// cutOK is cut's standard output, failing the test unless it exits 0.
func cutOK(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	out, stderr, code := cut(t, nil, stdin, args...)
	if code != 0 {
		t.Fatalf("cut %q: exit %d: %s", args, code, stderr)
	}
	return out
}

func TestCutPrintsFieldsLikeGNU(t *testing.T) {
	// GNU coreutils 9.4, LC_ALL=C: printf 'a:b:c\nplain\nx:y\n:\n' | cut -d: -f1,3
	if got := cutOK(t, "a:b:c\nplain\nx:y\n:\n", "-d:", "-f1,3"); got != "a:c\nplain\nx\n\n" {
		t.Errorf("cut -d: -f1,3 = %q, want GNU's", got)
	}
}

func TestCutReadsFilesInOrder(t *testing.T) {
	files := map[string]string{"one": "1:a\n1:b", "two": "2:a\n"}
	out, stderr, code := cut(t, files, "", "-d:", "-f2", "two", "one")
	if code != 0 || out != "a\na\nb\n" {
		t.Errorf("cut two one = %q (exit %d, %s), want two's lines then one's", out, code, stderr)
	}
}

func TestCutRefusesAnUnknownFlag(t *testing.T) {
	out, stderr, code := cut(t, nil, "a\tb\n", "-f1", "-Z")
	if code == 0 || out != "" || !strings.Contains(stderr, "'Z'") {
		t.Errorf("cut -Z: exit %d, stdout %q, stderr %q; want a failure, no output, Z named", code, out, stderr)
	}
}

func TestFieldsOnceInLineOrder(t *testing.T) {
	if got := cutOK(t, "a:b:c\n", "-d:", "-f3,1"); got != "a:c\n" {
		t.Errorf("-f3,1 = %q, want a:c", got)
	}
	if got := cutOK(t, "a:b:c\n", "-d:", "-f2,2"); got != "b\n" {
		t.Errorf("-f2,2 = %q, want b", got)
	}
}

func TestFieldRanges(t *testing.T) {
	for list, want := range map[string]string{"2-": "b:c:d\n", "-2": "a:b\n", "2-3": "b:c\n", "1,3-": "a:c:d\n"} {
		if got := cutOK(t, "a:b:c:d\n", "-d:", "-f"+list); got != want {
			t.Errorf("-f%s = %q, want %q", list, got, want)
		}
	}
}

func TestFieldsPastTheEnd(t *testing.T) {
	if got := cutOK(t, "a:b\n", "-d:", "-f5"); got != "\n" {
		t.Errorf("-f5 = %q, want an empty line", got)
	}
	if got := cutOK(t, "a:b\n", "-d:", "-f1,5"); got != "a\n" {
		t.Errorf("-f1,5 = %q, want a", got)
	}
}

func TestALineWithoutTheDelimiterIsPrintedWhole(t *testing.T) {
	if got := cutOK(t, "plain\n\nk:v\n:w\n", "-d:", "-f2"); got != "plain\n\nv\nw\n" {
		t.Errorf("cut -d: -f2 = %q, want plain and the empty line whole", got)
	}
}

func TestTheDelimiterIsATabByDefault(t *testing.T) {
	if got := cutOK(t, "a\tb\n", "-f2"); got != "b\n" {
		t.Errorf("-f2 = %q, want b", got)
	}
}

func TestTheDelimiterIsOneByte(t *testing.T) {
	out, stderr, code := cut(t, nil, "a::b\n", "-d::", "-f1")
	if code != 1 || out != "" || !strings.Contains(stderr, "the delimiter must be a single character") {
		t.Errorf("-d:: exit %d, stdout %q, stderr %q; want exit 1 and the message", code, out, stderr)
	}
}

func TestAFieldListIsRequired(t *testing.T) {
	out, stderr, code := cut(t, nil, "a\n", "-d:")
	if code != 1 || out != "" || !strings.Contains(stderr, "you must specify a list of bytes, characters, or fields") {
		t.Errorf("cut -d: exit %d, stdout %q, stderr %q; want exit 1 and the message", code, out, stderr)
	}
}
