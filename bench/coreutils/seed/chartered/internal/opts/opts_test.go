package opts

import (
	"fmt"
	"strings"
	"testing"

	"example.com/coreutils/claims"
	"github.com/leroux/charter/chartertest"
)

var spec = Spec{Flags: "lw", Valued: "df", Long: map[string]byte{"fields": 'f', "lines": 'l'}}

// parse renders what Parse made of args: the options, then the operands.
func parse(t *testing.T, args ...string) string {
	t.Helper()
	opts, operands, err := Parse(args, spec)
	if err != nil {
		t.Fatalf("Parse(%q): %v", args, err)
	}
	var b strings.Builder
	for _, o := range opts {
		fmt.Fprintf(&b, "-%c%s ", o.Name, o.Value)
	}
	return b.String() + "| " + strings.Join(operands, " ")
}

func TestClusteredOptions(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "clustered")
	if got := parse(t, "-lw", "f"); got != "-l -w | f" {
		t.Errorf("-lw f parsed as %q", got)
	}
}

func TestOptionsAfterOperands(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "after-operands")
	if got := parse(t, "a", "-l", "b"); got != "-l | a b" {
		t.Errorf("a -l b parsed as %q", got)
	}
}

func TestDoubleDashEndsOptions(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "double-dash")
	if got := parse(t, "-w", "--", "-l", "-"); got != "-w | -l -" {
		t.Errorf("-w -- -l - parsed as %q", got)
	}
}

func TestAValueAttachedOrNext(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "value-attached-or-next")
	a, b := parse(t, "-d:", "-f1"), parse(t, "-d", ":", "-f", "1")
	if a != "-d: -f1 | " || a != b {
		t.Errorf("-d: -f1 parsed as %q, -d : -f 1 as %q; want the same", a, b)
	}
}

func TestLongForms(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "long-forms")
	a, b := parse(t, "--fields=2"), parse(t, "--fields", "2")
	if a != "-f2 | " || b != a {
		t.Errorf("--fields=2 parsed as %q, --fields 2 as %q; want -f2", a, b)
	}
}

func TestAnUnknownOptionIsRefused(t *testing.T) {
	chartertest.Verifies(t, claims.SW_261126_2QDT, "unknown-refused")
	for args, want := range map[string]string{
		"-lZ":   "invalid option -- 'Z'",
		"--bad": "unrecognized option '--bad'",
		"-d":    "option requires an argument -- 'd'",
	} {
		if _, _, err := Parse([]string{args}, spec); err == nil || err.Error() != want {
			t.Errorf("Parse(%s) error %v, want %q", args, err, want)
		}
	}
}
