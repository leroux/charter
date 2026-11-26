package main

import (
	"bytes"
	"strings"
	"testing"

	"example.com/coreutils/claims"
	"github.com/leroux/charter/chartertest"
)

// sorted runs the command on lines given on standard input and returns
// the lines it printed, failing the test unless it exits 0.
func sorted(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := run(args, strings.NewReader(stdin), &out, &errOut); code != 0 {
		t.Fatalf("sort %q: exit %d: %s", args, code, errOut.String())
	}
	return out.String()
}

func TestSortKeepsAnUnterminatedLastLine(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_3NDQ, "last-line-kept")
	if got := sorted(t, "b\na"); got != "a\nb\n" {
		t.Errorf("sort of \"b\\na\" = %q, want a and b, each with its newline", got)
	}
}

func TestSortRefusesAnUnknownFlag(t *testing.T) {
	chartertest.Verifies(t, claims.STK_261126_WSGT, "unknown-flag-fails")
	var out, errOut bytes.Buffer
	code := run([]string{"-nZ"}, strings.NewReader("2\n1\n"), &out, &errOut)
	if code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "'Z'") {
		t.Errorf("sort -nZ: exit %d, stdout %q, stderr %q; want exit 2, no output, Z named", code, out.String(), errOut.String())
	}
}

func TestLinesSortByTheirBytes(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_YMCY, "bytewise")
	if got := sorted(t, "b\nB\na\nA\nab\n"); got != "A\nB\na\nab\nb\n" {
		t.Errorf("sort = %q, want A B a ab b", got)
	}
}

func TestReverseReversesTheOrder(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_YMCY, "reverse")
	if got := sorted(t, "b\nB\na\nA\nab\n", "-r"); got != "b\nab\na\nB\nA\n" {
		t.Errorf("sort -r = %q, want b ab a B A", got)
	}
}

func TestNumbersSortByValue(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_GKC7, "numeric-order")
	if got := sorted(t, "10\n9\n-1\n2.5\n", "-n"); got != "-1\n2.5\n9\n10\n" {
		t.Errorf("sort -n = %q, want -1 2.5 9 10", got)
	}
}

func TestANonNumberReadsAsZero(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_GKC7, "non-number-zero")
	// The zeros tie and fall back to bytes: "+5" < "-0" < "abc".
	if got := sorted(t, "1\nabc\n-1\n+5\n-0\n", "-n"); got != "-1\n+5\n-0\nabc\n1\n" {
		t.Errorf("sort -n = %q, want -1, the zeros by bytes, 1", got)
	}
}

func TestEqualNumbersTieByBytes(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_GKC7, "ties-bytewise")
	if got := sorted(t, "2.50\n2.5\n", "-n"); got != "2.5\n2.50\n" {
		t.Errorf("sort -n = %q, want 2.5 then 2.50", got)
	}
	if got := sorted(t, "2.5\n2.50\n", "-rn"); got != "2.50\n2.5\n" {
		t.Errorf("sort -rn = %q, want 2.50 then 2.5", got)
	}
}

func TestNumbersCompareExactly(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_GKC7, "exact-digits")
	in := "99999999999999999999\n99999999999999999998\n-99999999999999999998\n-99999999999999999999\n"
	want := "-99999999999999999999\n-99999999999999999998\n99999999999999999998\n99999999999999999999\n"
	if got := sorted(t, in, "-n"); got != want {
		t.Errorf("sort -n = %q, want %q", got, want)
	}
}
