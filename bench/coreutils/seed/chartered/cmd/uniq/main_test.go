package main

import (
	"bytes"
	"strings"
	"testing"

	"example.com/coreutils/claims"
	"github.com/leroux/charter/chartertest"
)

// uniq runs the command on standard input and returns what it printed and
// its exit status.
func uniq(stdin string, args ...string) (stdout, stderr string, code int) {
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestUniqCountsLikeGNU(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_69ZV, "uniq-count-like-gnu")
	// GNU coreutils 9.4, LC_ALL=C: printf 'pear\napple\npear\nfig\napple\npear\n' | sort | uniq -c
	got, _, code := uniq("apple\napple\nfig\npear\npear\npear\n", "-c")
	if want := "      2 apple\n      1 fig\n      3 pear\n"; code != 0 || got != want {
		t.Errorf("uniq -c = %q (exit %d), want GNU's %q", got, code, want)
	}
}

func TestUniqRefusesAnUnknownFlag(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_WSGT, "unknown-flag-fails")
	out, stderr, code := uniq("a\na\n", "-cZ")
	if code == 0 || out != "" || !strings.Contains(stderr, "'Z'") {
		t.Errorf("uniq -cZ: exit %d, stdout %q, stderr %q; want a failure, no output, Z named", code, out, stderr)
	}
}

// merged is uniq's output for lines, one a line.
func merged(t *testing.T, args []string, lines ...string) string {
	t.Helper()
	out, stderr, code := uniq(strings.Join(lines, "\n")+"\n", args...)
	if code != 0 {
		t.Fatalf("uniq %q: exit %d: %s", args, code, stderr)
	}
	return out
}

func TestAdjacentLinesMerge(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4T6M, "adjacent-merged")
	if got := merged(t, nil, "a", "a", "b"); got != "a\nb\n" {
		t.Errorf("uniq = %q, want a b", got)
	}
}

func TestEqualLinesApartStayApart(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4T6M, "apart-stay")
	if got := merged(t, nil, "a", "b", "a"); got != "a\nb\na\n" {
		t.Errorf("uniq = %q, want a b a", got)
	}
}

func TestCaseMatters(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_4T6M, "case-matters")
	if got := merged(t, nil, "Apple", "apple"); got != "Apple\napple\n" {
		t.Errorf("uniq = %q, want both lines", got)
	}
}

func TestCountsInSevenColumns(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_8785, "count-seven-columns")
	if got := merged(t, []string{"-c"}, "a", "a", "b"); got != "      2 a\n      1 b\n" {
		t.Errorf("uniq -c = %q", got)
	}
}
