package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"example.com/ratelimit/claims"
)

// stillClock never moves, so a run's requests all fall at one instant.
type stillClock struct{}

func (stillClock) Now() time.Time { return time.Unix(1_700_000_000, 0) }

// drive runs the command as a script in front of a service would: keys in,
// one a line, and a verdict line out for each.
func drive(t *testing.T, vars map[string]string, keys ...string) (verdicts []string, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	getenv := func(k string) string { return vars[k] }
	code = run(strings.NewReader(strings.Join(keys, "\n")+"\n"), &out, &errOut, getenv, stillClock{})
	if s := strings.TrimSpace(out.String()); s != "" {
		verdicts = strings.Split(s, "\n")
	}
	return verdicts, errOut.String(), code
}

// verdictsOf drives the command and fails the test unless it exits 0.
func verdictsOf(t *testing.T, vars map[string]string, keys ...string) []string {
	t.Helper()
	got, stderr, code := drive(t, vars, keys...)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, stderr)
	}
	return got
}

func TestAQuietCallerBurstsAndTheNextIsDenied(t *testing.T) {
	claims.Verifies(t, claims.STK_261126_H2SW, "burst-then-deny")
	got := verdictsOf(t, nil, "alice", "alice", "alice", "alice", "alice", "alice")
	want := []string{"allow alice", "allow alice", "allow alice", "allow alice", "allow alice", "deny alice"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("verdicts %q, want %q", got, want)
	}
}

func TestASpentCallerDoesNotStopAnother(t *testing.T) {
	claims.Verifies(t, claims.STK_261126_H2SW, "callers-apart")
	keys := []string{"alice", "alice", "alice", "alice", "alice", "alice", "alice", "bob"}
	got := verdictsOf(t, nil, keys...)
	if got[len(got)-1] != "allow bob" {
		t.Errorf("bob's verdict is %q after alice spent her burst, want allow", got[len(got)-1])
	}
}

func TestAnOperatorSetsTheBurst(t *testing.T) {
	claims.Verifies(t, claims.STK_261126_T2BD, "burst-set")
	got := verdictsOf(t, map[string]string{"RATELIMIT_BURST": "2"}, "a", "a", "a")
	if got[2] != "deny a" {
		t.Errorf("verdicts %q, want the third denied under a burst of 2", got)
	}
}

func TestABadSettingStopsTheStart(t *testing.T) {
	claims.Verifies(t, claims.STK_261126_T2BD, "bad-setting-stops")
	got, stderr, code := drive(t, map[string]string{"RATELIMIT_RATE": "fast"}, "a")
	if code == 0 || len(got) != 0 || !strings.Contains(stderr, "RATELIMIT_RATE") {
		t.Errorf("exit %d, verdicts %q, stderr %q; want a failing exit, no verdict and the setting named", code, got, stderr)
	}
}

func TestTheDefaultsAreOneASecondAndABurstOfFive(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_CMTG, "defaults")
	rate, burst, err := config(func(string) string { return "" })
	if err != nil || rate != 1 || burst != 5 {
		t.Errorf("config() = %v, %v, %v; want 1, 5", rate, burst, err)
	}
}

func TestTheEnvironmentSetsTheLimits(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_CMTG, "set-from-env")
	for _, tc := range []struct {
		rate, burst string
		wantRate    float64
		wantBurst   int
	}{
		{"2.5", "7", 2.5, 7},
		{"0.001", "1", 0.001, 1}, // the least of each
	} {
		vars := map[string]string{"RATELIMIT_RATE": tc.rate, "RATELIMIT_BURST": tc.burst}
		rate, burst, err := config(func(k string) string { return vars[k] })
		if err != nil || rate != tc.wantRate || burst != tc.wantBurst {
			t.Errorf("config() with %s, %s = %v, %v, %v; want %v, %v", tc.rate, tc.burst, rate, burst, err, tc.wantRate, tc.wantBurst)
		}
	}
}

func TestAnUnreadableSettingIsRefusedBeforeAnyKey(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_CMTG, "unreadable-refused")
	for _, tc := range []struct {
		name, value string
	}{
		{"RATELIMIT_RATE", "fast"},
		{"RATELIMIT_BURST", "2.5"},
		{"RATELIMIT_RATE", "0"},
		{"RATELIMIT_BURST", "0"},
	} {
		got, stderr, code := drive(t, map[string]string{tc.name: tc.value}, "a")
		if code != 2 || len(got) != 0 || !strings.Contains(stderr, tc.name) {
			t.Errorf("%s=%s: exit %d, verdicts %q, stderr %q; want exit 2, no verdict, the setting named", tc.name, tc.value, code, got, stderr)
		}
	}
}
