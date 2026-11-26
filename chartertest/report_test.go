package chartertest

import (
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// A go test -json stream as the fixture's tests would emit it: TestRestart
// cites SW-001 and two scenarios of SW-010 and passes; TestLock cites
// SW-002 and is skipped; TestFold cites a scenario and fails.
const goStream = `{"Action":"run","Package":"example.com/fixture/pkg/a","Test":"TestRestart"}
{"Action":"attr","Package":"example.com/fixture/pkg/a","Test":"TestRestart","Key":"claim","Value":"SW-001"}
{"Action":"attr","Package":"example.com/fixture/pkg/a","Test":"TestRestart","Key":"claim","Value":"SW-010/waited-request-passes"}
{"Action":"attr","Package":"example.com/fixture/pkg/a","Test":"TestRestart","Key":"claim","Value":"SW-010/burst-caps-tokens"}
{"Action":"pass","Package":"example.com/fixture/pkg/a","Test":"TestRestart","Elapsed":0.01}
{"Action":"run","Package":"example.com/fixture/pkg/a","Test":"TestLock"}
{"Action":"attr","Package":"example.com/fixture/pkg/a","Test":"TestLock","Key":"claim","Value":"SW-002"}
{"Action":"skip","Package":"example.com/fixture/pkg/a","Test":"TestLock","Elapsed":0}
{"Action":"run","Package":"example.com/fixture/pkg/a","Test":"TestFold"}
{"Action":"attr","Package":"example.com/fixture/pkg/a","Test":"TestFold","Key":"claim","Value":"SW-010/burst-caps-tokens"}
{"Action":"fail","Package":"example.com/fixture/pkg/a","Test":"TestFold","Elapsed":0.02}
{"Action":"pass","Package":"example.com/fixture/pkg/a","Elapsed":0.1}
`

func reportOfFixture(t *testing.T, stream string) *Report {
	t.Helper()
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	cl, oc, err := ReadGoStream(strings.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	return Build(reg, cl, oc, nil)
}

func entryOf(rep *Report, id string) *EntryReport {
	for i := range rep.Entries {
		if rep.Entries[i].ID == id {
			return &rep.Entries[i]
		}
	}
	return nil
}

// Each scenario gets the tests that cited it and the fold of their outcomes.
func TestReportJoinsGoResultsToScenarios(t *testing.T) {
	Verifies(t, claims.CH_042, "go-results-joined")
	rep := reportOfFixture(t, goStream)
	e := entryOf(rep, "SW-010")
	if e == nil || len(e.Scenarios) != 2 {
		t.Fatalf("SW-010 = %+v", e)
	}
	if s := e.Scenarios[0]; s.Handle != "waited-request-passes" || s.Outcome != Passed || len(s.Tests) != 1 || s.Tests[0].Test != "TestRestart" {
		t.Errorf("waited-request-passes = %+v", s)
	}
	if s := e.Scenarios[1]; s.Handle != "burst-caps-tokens" || s.Outcome != Failed || len(s.Tests) != 2 {
		t.Errorf("burst-caps-tokens = %+v", s)
	}
	if e.Outcome != Failed {
		t.Errorf("SW-010 outcome = %s, want fail", e.Outcome)
	}
	if e := entryOf(rep, "SW-001"); e == nil || e.Outcome != Passed || len(e.Tests) != 1 {
		t.Errorf("SW-001 = %+v", e)
	}
}

// A scenario whose every citing test was skipped, or that no test ran, is
// unverified; a spec without scenarios likewise; planned and deprecated
// specs are not counted.
func TestReportListsSkippedAndUnrunScenariosAsUnverified(t *testing.T) {
	Verifies(t, claims.CH_042, "skipped-scenario-fails")
	rep := reportOfFixture(t, goStream)
	want := []string{"SW-010/burst-caps-tokens", "STK-001", "CON-001", "CON-002"}
	if strings.Join(rep.Unverified, " ") != strings.Join(want, " ") {
		t.Errorf("unverified = %v, want %v", rep.Unverified, want)
	}
	if e := entryOf(rep, "SW-002"); e.Outcome != Skipped {
		t.Errorf("the deprecated, skipped spec = %+v", e)
	}
	for _, u := range rep.Unverified {
		if strings.HasPrefix(u, "SW-002") || strings.HasPrefix(u, "SW-003") {
			t.Errorf("an inactive spec is listed as unverified: %s", u)
		}
	}
}

// A waived scenario or spec that no test passed is waived, not unverified.
func TestReportHonoursWaivers(t *testing.T) {
	Verifies(t, claims.CH_042, "skipped-scenario-fails")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "claims/waivers.txt", "SW-010/burst-caps-tokens\tthe burst case is verified by review\nCON-001\tthe constraints are verified by review\nCON-002\tthe constraints are verified by review\n")
	rep, err := BuildReport(fixtureOptions(root), strings.NewReader(goStream))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(rep.Unverified, " ") != "STK-001" {
		t.Errorf("unverified = %v, want STK-001 alone", rep.Unverified)
	}
	if s := entryOf(rep, "SW-010").Scenarios[1]; s.Outcome != Waived {
		t.Errorf("the waived scenario = %+v", s)
	}
	if e := entryOf(rep, "CON-001"); e.Outcome != Waived {
		t.Errorf("the waived spec = %+v", e)
	}
}

// The markdown lists the unverified scenarios first, then each spec with
// its scenarios and their outcomes.
func TestReportRendersMarkdown(t *testing.T) {
	Verifies(t, claims.CH_042, "markdown-rendered")
	md := reportOfFixture(t, goStream).Markdown()
	for _, want := range []string{
		"4 active scenarios did not pass:",
		"- `SW-010/burst-caps-tokens`",
		"## SW-010 The bucket refills from its last request at the rate",
		"- `burst-caps-tokens` (agent, provisional): fail · TestFold (fail) · TestRestart (pass)",
		"## SW-001 A bucket survives a restart",
		"- TestRestart (pass)",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown lacks %q:\n%s", want, md)
		}
	}
}

// The stream of a real run is read: the attr probe's own tests.
func TestReportReadsARealRun(t *testing.T) {
	Verifies(t, claims.CH_042, "go-results-joined")
	out, err := RunGoTests(Options{Entries: "github.com/leroux/charter/claims"}, "./chartertest/internal/attrprobe")
	if err != nil {
		t.Fatal(err)
	}
	cl, oc, err := ReadGoStream(strings.NewReader(string(out)))
	if err != nil {
		t.Fatal(err)
	}
	var probe, sub string
	for k := range cl {
		if strings.HasSuffix(k, "\x00TestProbe") {
			probe = k
		}
		if strings.HasSuffix(k, "\x00TestProbe/sub") {
			sub = k
		}
	}
	if probe == "" || sub == "" {
		t.Fatalf("the probe's tests were not found in the stream: claims = %v\n%s", cl, out)
	}
	if cl[probe][0] != "SW-001" || oc[probe] != Passed || cl[sub][0] != "SW-002/skipped" || oc[sub] != Skipped {
		t.Errorf("claims = %v, outcomes = %v", cl, oc)
	}
}

// The report says who ruled each scenario and counts the active ones by it.
func TestReportShowsWhoRuled(t *testing.T) {
	Verifies(t, claims.CH_042, "provenance-and-profile")
	Verifies(t, claims.CH_261009_3FHW, "report-counts-who-ruled")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", strings.Replace(sampleArea, "1. `waited-request-passes`: a request", "1. `waited-request-passes` (owner, #444): a request", 1))
	reg, err := Read(fixtureOptions(root))
	if err != nil {
		t.Fatal(err)
	}
	cl, oc, err := ReadGoStream(strings.NewReader(goStream))
	if err != nil {
		t.Fatal(err)
	}
	rep := Build(reg, cl, oc, nil)
	e := entryOf(rep, "SW-010")
	if e.Scenarios[0].Provenance != "owner, #444" || e.Scenarios[1].Provenance != "agent, provisional" {
		t.Errorf("provenance in the report: %+v", e.Scenarios)
	}
	if rep.Profile.Owner != (Count{Scenarios: 1, Passed: 1}) || rep.Profile.Provisional != (Count{Scenarios: 1}) || rep.Profile.Delegate != (Count{}) {
		t.Errorf("profile = %+v", rep.Profile)
	}
	md := rep.Markdown()
	for _, want := range []string{
		"Active scenarios by who ruled them: 1 ruled by the owner (1 passed), 0 by a delegate (0 passed), 1 an agent's, provisional (0 passed).",
		"- `waited-request-passes` (owner, #444): pass · TestRestart (pass)",
		"- `burst-caps-tokens` (agent, provisional): fail",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown lacks %q:\n%s", want, md)
		}
	}
}
