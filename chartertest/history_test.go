package chartertest

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

func pkgOf(t *testing.T, src string) *entriesPackage {
	t.Helper()
	pkg := newEntriesPackage()
	_, es, problems := parseArea("claims/refill.md", src)
	if len(problems) != 0 {
		t.Fatalf("problems %q", problems)
	}
	for _, e := range es {
		pkg.add(e)
	}
	return pkg
}

// History names what a change did to the record that no Changes line says.
func TestHistoryNamesWhatTheChangesLinesDoNotSay(t *testing.T) {
	Verifies(t, claims.CH_261009_QQC8, "removed-named", "unbuilt-named", "unrecorded-change-named", "ruling-lost-named", "recorded-change-passes")
	Verifies(t, claims.CH_261009_3FHW, "unrecorded-change-named")
	base := strings.Replace(sampleArea, "1. `waited-request-passes`: a request", "1. `waited-request-passes` (owner, #444): a request", 1)
	prev := pkgOf(t, base)
	for _, tc := range []struct{ name, old, repl, want string }{
		{"removed", base[len("# Refill\nCode: pkg/a\n\n"):], "", "claims/refill.md:4: history: SW-010, in claims/refill.md at base, is gone; an entry is deprecated, never deleted"},
		{"unbuilt", "Status: active", "Status: planned", "history: SW-010 went from active to planned; a built decision is deprecated with a Changes line, not unbuilt"},
		{"deprecated without a line", "Status: active", "Status: deprecated", "history: SW-010 was deprecated with no Changes line saying what replaces it"},
		{"text changed", "At each request", "At every request", "history: SW-010 changed since base (text) with no Changes line saying what and why"},
		{"title changed", "The bucket refills from", "The bucket fills from", "changed since base (title)"},
		{"kind changed", "Kind: invariant", "Kind: bound", "changed since base (kind)"},
		{"scenario reworded", "a request after a wait of one token passes", "a request after a wait of one token passes first", "changed since base (scenario waited-request-passes reworded)"},
		{"scenario removed", "2. `burst-caps-tokens`: a long wait earns no more than the burst,\n   and the excess is not kept.\n", "", "changed since base (scenario burst-caps-tokens removed)"},
		{"scenario added", "   and the excess is not kept.\n", "   and the excess is not kept.\n3. `more`: more.\n", "changed since base (scenario more added)"},
		{"ruling lost", " (owner, #444)", "", "history: SW-010/waited-request-passes lost its ruling: it was ruled (owner, #444) and is now an agent's, provisional"},
	} {
		src := strings.Replace(base, tc.old, tc.repl, 1)
		if src == base {
			t.Fatalf("%s: the edit did not apply", tc.name)
		}
		got := compareHistory(prev, pkgOf(t, src), "base")
		if !anyContains(got, tc.want) {
			t.Errorf("%s: want %q in %q", tc.name, tc.want, got)
		}
	}
	recorded := strings.Replace(base, "At each request", "At every request", 1) + "- 2026-10-09: every edge, not each (backlog 081).\n"
	if got := compareHistory(prev, pkgOf(t, recorded), "base"); len(got) != 0 {
		t.Errorf("a recorded change was named: %q", got)
	}
	planned := pkgOf(t, strings.Replace(base, "Status: active", "Status: planned", 1))
	if got := compareHistory(planned, prev, "base"); len(got) != 0 {
		t.Errorf("setting an entry active was named: %q", got)
	}
	if got := compareHistory(prev, prev, "base"); len(got) != 0 {
		t.Errorf("an unchanged record was named: %q", got)
	}

	// Against a git revision.
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", base)
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-q", "-m", "base"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	writeFile(t, root, "claims/refill.md", strings.Replace(base, "At each request", "At every request", 1))
	got, err := History(fixtureOptions(root), "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !strings.Contains(got[0], "SW-010 changed since HEAD (text) with no Changes line") {
		t.Errorf("History = %q", got)
	}
}
