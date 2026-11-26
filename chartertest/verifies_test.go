package chartertest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// TestVerifiesEmitsAttrEvents runs the probe package under go test -json
// and reads the attr events Verifies produces.
func TestVerifiesEmitsAttrEvents(t *testing.T) {
	Verifies(t, claims.CH_024, "attr-per-citation")
	out, err := exec.Command("go", "test", "-json", "-count=1", "./internal/attrprobe").Output()
	if err != nil {
		t.Fatalf("go test: %v\n%s", err, out)
	}
	type event struct{ Action, Test, Key, Value string }
	var attrs []event
	skipped := map[string]bool{}
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		var e event
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			continue
		}
		switch e.Action {
		case "attr":
			attrs = append(attrs, e)
		case "skip":
			skipped[e.Test] = true
		}
	}
	want := []event{
		{Action: "attr", Test: "TestProbe", Key: "claim", Value: "SW-001"},
		{Action: "attr", Test: "TestProbe/sub", Key: "claim", Value: "SW-002/skipped"},
		{Action: "attr", Test: "TestProbeObserves", Key: "observation", Value: "ProbeFact"},
	}
	if len(attrs) != len(want) {
		t.Fatalf("attr events = %v, want %v", attrs, want)
	}
	for i := range want {
		if attrs[i] != want[i] {
			t.Errorf("attr[%d] = %v, want %v", i, attrs[i], want[i])
		}
	}
	if !skipped["TestProbe/sub"] {
		t.Error("the skipped subtest did not report skip after its attr")
	}
}

// Verifies refuses a nil spec, the one wrong argument the compiler lets
// through.
func TestVerifiesRefusesANilSpec(t *testing.T) {
	Verifies(t, claims.CH_024, "nil-spec-refused")
	fake := &fatalRecorder{TB: t}
	call := Verifies // through a variable: the gate reads a direct call as a citation
	call(fake, nil)
	if !strings.Contains(fake.msg, "Verifies needs a spec") {
		t.Fatalf("Verifies did not fail on a nil spec; message %q", fake.msg)
	}
}

// fatalRecorder captures Fatalf instead of ending the test.
type fatalRecorder struct {
	testing.TB
	msg string
}

func (f *fatalRecorder) Helper() {}

func (f *fatalRecorder) Fatalf(format string, args ...any) { f.msg = fmt.Sprintf(format, args...) }

func (f *fatalRecorder) Fatal(args ...any) { f.msg = fmt.Sprint(args...) }
