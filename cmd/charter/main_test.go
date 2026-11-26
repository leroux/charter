package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leroux/charter/chartertest"
	"github.com/leroux/charter/claims"
)

const planned = "# Area\n\n## SW-010 A rule\nStatus: planned\n\nWhat it rules.\n\nWhy: someone needs it.\n"

func charterIn(t *testing.T, area string) string {
	t.Helper()
	root := t.TempDir()
	for rel, src := range map[string]string{"claims/charter.toml": "name = \"x\"\n", "claims/area.md": area} {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestCheckSetsItsExitStatus(t *testing.T) {
	chartertest.Verifies(t, claims.CH_037, "exit-status")
	var out, errOut bytes.Buffer
	if code := run([]string{"check"}, charterIn(t, planned), &out, &errOut); code != 0 || out.Len() != 0 {
		t.Errorf("a clean tree: exit %d, output %q %q", code, out.String(), errOut.String())
	}
	active := strings.Replace(planned, "Status: planned", "Status: active", 1) + "\nScenarios:\n1. `rules`: it rules.\n"
	out.Reset()
	if code := run([]string{"check"}, charterIn(t, active), &out, &errOut); code != 1 || !strings.Contains(out.String(), "claims/area.md:3: rule 25: active entry SW-010 is realized by no code") {
		t.Errorf("a tree with a problem: exit %d, output %q", code, out.String())
	}
}

func TestEveryVerbWorksOnTheCharterFoundAbove(t *testing.T) {
	chartertest.Verifies(t, claims.CH_040, "verbs-on-the-found-charter")
	root := charterIn(t, planned)
	deep := filepath.Join(root, "pkg", "a")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	for _, args := range [][]string{{"gen"}, {"check"}, {"lint"}, {"next", "SW"}, {"anchors", "SW-010"}, {"questions"}, {"export"}} {
		out.Reset()
		if code := run(args, deep, &out, &errOut); code != 0 {
			t.Errorf("charter %v from %s: exit %d, %s%s", args, deep, code, out.String(), errOut.String())
		}
	}
	if code := run([]string{"check"}, t.TempDir(), &out, &errOut); code != 2 || !strings.Contains(errOut.String(), "no claims/charter.toml") {
		t.Errorf("outside any charter: exit %d, %q", code, errOut.String())
	}
}

func TestQuestionsListsEachOpenQuestion(t *testing.T) {
	chartertest.Verifies(t, claims.CH_261009_4B3B, "questions-verb-lists")
	root := charterIn(t, planned)
	q := "# Questions\n\n## Q-261009-M4TK Does it rule?\nAsked-by: agent\nEntry: SW-010\nProvisional: yes\n\n## Q-261009-N7PR Is it whole?\nAsked-by: owner\nEntry: SW-010\nProvisional: no\nAnswer: 2026-10-09, owner: yes\n"
	if err := os.WriteFile(filepath.Join(root, "claims", "questions.md"), []byte(q), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	run([]string{"questions"}, root, &out, &errOut)
	if want := "Q-261009-M4TK  SW-010\n  Does it rule?\n  provisional: yes\n  asked by: agent\n\n1 open\n"; out.String() != want {
		t.Errorf("questions =\n%s\nwant\n%s%s", out.String(), want, errOut.String())
	}
	out.Reset()
	run([]string{"questions", "-all"}, root, &out, &errOut)
	if !strings.Contains(out.String(), "Q-261009-N7PR  SW-010") || !strings.Contains(out.String(), "  answer: 2026-10-09, owner: yes\n") || !strings.HasSuffix(out.String(), "1 open\n") {
		t.Errorf("questions -all =\n%s", out.String())
	}
}
