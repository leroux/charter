package diff

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
	"github.com/leroux/charter/internal/config"
	"github.com/leroux/charter/internal/record"
)

const sampleArea = "# Refill\nCode: pkg/a\n\n" +
	"## SW-010 The bucket refills from its last request at the rate\n" +
	"Status: active\n" +
	"Kind: invariant\n" +
	"Serves: STK-001\n" +
	"See: SW-001\n" +
	"Accounts: env:FX_DATA_DIR\n" +
	"\n" +
	"At each request the bucket reads the clock.\n" +
	"When time passed, the earned tokens are added to the bucket.\n" +
	"\n" +
	"Why: a caller that waits is let through again at the rate.\n" +
	"\n" +
	"Scenarios:\n" +
	"1. `waited-request-passes`: a request after a wait of one token passes.\n" +
	"2. `burst-caps-tokens`: a long wait earns no more than the burst,\n" +
	"   and the excess is not kept.\n" +
	"\n" +
	"Guidance: no timer; the clock is read at a request only.\n" +
	"\n" +
	"Changes:\n" +
	"- 2026-09-28: a denied request reads the clock once more (backlog 034).\n"

func mapOf(t *testing.T, src string) map[string]*record.Entry {
	t.Helper()
	_, es, problems := record.ParseArea("claims/refill.md", src)
	if len(problems) != 0 {
		t.Fatalf("problems %q", problems)
	}
	m := map[string]*record.Entry{}
	for _, e := range es {
		m[e.ID] = e
	}
	return m
}

func anyContains(list []string, want string) bool {
	for _, s := range list {
		if strings.Contains(s, want) {
			return true
		}
	}
	return false
}

func writeFile(t *testing.T, root, rel, src string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

// History names what a change did to the record that no Changes line says.
func TestHistoryNamesWhatTheChangesLinesDoNotSay(t *testing.T) {
	claims.Verifies(t, claims.CH_261009_QQC8, "removed-named", "unbuilt-named", "unrecorded-change-named", "ruling-lost-named", "recorded-change-passes")
	claims.Verifies(t, claims.CH_261009_3FHW, "unrecorded-change-named")
	base := strings.Replace(sampleArea, "1. `waited-request-passes`: a request", "1. `waited-request-passes` (owner, #444): a request", 1)
	prev := mapOf(t, base)
	order := []string{"SW-010"}
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
		got := compare(prev, mapOf(t, src), order, "base")
		if !anyContains(got, tc.want) {
			t.Errorf("%s: want %q in %q", tc.name, tc.want, got)
		}
	}
	recorded := strings.Replace(base, "At each request", "At every request", 1) + "- 2026-10-09: every edge, not each (backlog 081).\n"
	if got := compare(prev, mapOf(t, recorded), order, "base"); len(got) != 0 {
		t.Errorf("a recorded change was named: %q", got)
	}
	planned := mapOf(t, strings.Replace(base, "Status: active", "Status: planned", 1))
	if got := compare(planned, prev, order, "base"); len(got) != 0 {
		t.Errorf("setting an entry active was named: %q", got)
	}
	if got := compare(prev, prev, order, "base"); len(got) != 0 {
		t.Errorf("an unchanged record was named: %q", got)
	}

	// Against a git revision.
	root := t.TempDir()
	writeFile(t, root, "claims/charter.toml", "name = \"fx\"\n")
	writeFile(t, root, "claims/refill.md", base)
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-q", "-m", "base"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	writeFile(t, root, "claims/refill.md", strings.Replace(base, "At each request", "At every request", 1))
	c, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Diff(c, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !strings.Contains(got[0], "SW-010 changed since HEAD (text) with no Changes line") {
		t.Errorf("History = %q", got)
	}
}
