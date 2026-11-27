package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/leroux/charter/claims"
)

// docsArea is the record a doc tells a story over: an active decision, a
// planned one and a deprecated one.
const docsArea = "# Refill\nCode: pkg/a\n\nHow the bucket refills.\n\n" +
	"## SW-010 The bucket refills\nStatus: active\n\nIt refills at the rate.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `refills`: it refills.\n\n" +
	"## SW-011 A planned rule\nStatus: planned\n\nLater.\n\nWhy: someone needs it.\n\n" +
	"## SW-012 A retired rule\nStatus: deprecated\n\nGone.\n\nWhy: it was.\n\nChanges:\n- 2026-11-01: retired.\n"

// docsCharter is a charter of Go, Python and TypeScript whose docs are
// docs/*.md, with the doc given as docs/tour.md.
func docsCharter(t *testing.T, tour string) string {
	t.Helper()
	compiler, err := filepath.Abs("../../node_modules/typescript")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"go.mod": "module example.com/fx\n\ngo 1.21\n\nrequire example.com/fx/claims v0.0.0\n\nreplace example.com/fx/claims => ./claims\n",
		"claims/charter.toml": "name = \"fx\"\ndocs = [\"docs/*.md\"]\n\n[go]\nmodule = \"example.com/fx/claims\"\ndir = \".\"\ntest = \"go test ./...\"\n" +
			"\n[python]\ndir = \"py\"\ntest = \"python3 -m unittest\"\n" +
			"\n[typescript]\ndir = \"web\"\ntest = \"vitest run\"\ncompiler = \"" + compiler + "\"\nimport = \"@fx/claims\"\n",
		"claims/refill.md": docsArea,
		"pkg/a/a.go":       "package a\n\n// Refill refills.\n//\n//charter:realizes SW-010\nfunc Refill() {}\n\n// Bucket holds tokens.\ntype Bucket struct{}\n\n// Take spends one.\nfunc (b *Bucket) Take() {}\n",
		"pkg/a/a_test.go":  "package a\n\nimport (\n\t\"testing\"\n\n\t\"example.com/fx/claims\"\n)\n\nfunc TestRefill(t *testing.T) {\n\tclaims.Verifies(t, claims.SW_010, \"refills\")\n}\n",
		"py/acme/run.py":   "class Meter:\n    def tick(self):\n        return 1\n",
		"web/main.ts":      "export function spend(n: number): number {\n  return n - 1;\n}\n",
		"notes.txt":        "A note no language reads.\n",
		"README.md":        "Not a doc: [SW-999](x).\n",
		"docs/tour.md":     tour,
	}
	root := writeTree(t, files)
	var out, errOut bytes.Buffer
	if code := run([]string{"gen"}, root, &out, &errOut); code != 0 {
		t.Fatalf("gen: exit %d %s", code, errOut.String())
	}
	return root
}

// checked is charter check's output and exit status.
func checked(t *testing.T, root string) (string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run([]string{"check"}, root, &out, &errOut)
	if code == 2 {
		t.Fatalf("check: %s", errOut.String())
	}
	return out.String(), code
}

// warned is what charter check writes to standard error: its warnings.
func warned(t *testing.T, root string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	if code := run([]string{"check"}, root, &out, &errOut); code == 2 {
		t.Fatalf("check: %s", errOut.String())
	}
	return errOut.String()
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func writeFile(t *testing.T, root, rel, src string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(rel)), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

func edit(t *testing.T, root, rel, old, new string) {
	t.Helper()
	src := readFile(t, root, rel)
	if strings.Count(src, old) != 1 {
		t.Fatalf("%s holds %q %d times", rel, old, strings.Count(src, old))
	}
	writeFile(t, root, rel, strings.Replace(src, old, new, 1))
}

func TestADocIsAFileTheDocPatternsName(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_3NY3, "doc-by-pattern")
	root := docsCharter(t, "# Tour\n\nSee [SW-999](x).\n")
	out, _ := checked(t, root)
	if !strings.Contains(out, "docs/tour.md:3: rule 26: the reference names SW-999, which is not an entry of the record") {
		t.Errorf("a doc the patterns match is not read:\n%s", out)
	}
	if strings.Contains(out, "README.md") {
		t.Errorf("a file the patterns do not match is read as a doc:\n%s", out)
	}
	// The walk's only patterns do not narrow the docs, and a matched file
	// that is not markdown is no doc.
	edit(t, root, "claims/charter.toml", "docs = [\"docs/*.md\"]\n", "docs = [\"docs/*\"]\nonly = [\"pkg\"]\n")
	writeFile(t, root, "docs/note.txt", "[SW-997](x)\n")
	out, _ = checked(t, root)
	if !strings.Contains(out, "docs/tour.md:3: rule 26") {
		t.Errorf("only narrowed the docs:\n%s", out)
	}
	if strings.Contains(out, "note.txt") {
		t.Errorf("a file that is not markdown is read as a doc:\n%s", out)
	}
}

func TestAReferenceToAnEntryMustResolve(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_3NY3, "reference-resolves", "deprecated-reference-fails", "embed-resolves")
	root := docsCharter(t, "# Tour\n\nThe bucket: [SW-010](x), later [SW-011](x).\n")
	if out, code := checked(t, root); code != 0 {
		t.Errorf("references to an active and a planned entry fail:\n%s", out)
	}
	root = docsCharter(t, "# Tour\n\n[SW-999](x)\n\n[SW-012](x)\n\n<!-- charter:embed SW-998 -->\n<!-- /charter:embed -->\n")
	out, _ := checked(t, root)
	for _, want := range []string{
		"docs/tour.md:3: rule 26: the reference names SW-999, which is not an entry of the record",
		"docs/tour.md:5: rule 26: the reference names SW-012, which is deprecated",
		"docs/tour.md:7: rule 26: the embed names SW-998, which is not an entry of the record",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("want\n  %s\ngot\n%s", want, out)
		}
	}
}

func TestADocDeclaresNoEntry(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_3NY3, "doc-declares-nothing")
	root := docsCharter(t, "# Tour\n\n## SW-013 A rule written in a doc\n\nStatus: active\n")
	out, _ := checked(t, root)
	if !strings.Contains(out, "docs/tour.md:3: rule 26: a doc declares no entry") {
		t.Errorf("a heading with an id in a doc is not refused:\n%s", out)
	}
	var aout, aerr bytes.Buffer
	if code := run([]string{"anchors", "SW-013"}, root, &aout, &aerr); code != 2 {
		t.Errorf("an entry was read from a doc: %s", aout.String())
	}
}

func TestACodeReferenceNamesADeclarationOfItsFile(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_DW8R, "go-symbol-resolves", "python-symbol-resolves", "typescript-symbol-resolves", "method-by-type", "plain-file-exists", "prose-link-unread")
	good := "# Tour\n\n" +
		"[`Refill`](../pkg/a/a.go) and [`a.Refill`](../pkg/a/a.go#L1) and [`Bucket.Take`](../pkg/a/a.go).\n" +
		"[`Meter`](../py/acme/run.py) and [`Meter.tick`](../py/acme/run.py).\n" +
		"[`spend`](../web/main.ts) and [`main.spend`](../web/main.ts).\n" +
		"[`notes.txt`](../notes.txt).\n" +
		"[a prose link](../pkg/a/nowhere.go) and [`Far`](https://example.com/far.go).\n"
	if out, code := checked(t, docsCharter(t, good)); code != 0 {
		t.Errorf("references to declarations the files hold fail:\n%s", out)
	}
	bad := "# Tour\n\n" +
		"[`Drain`](../pkg/a/a.go)\n" +
		"[`Gauge`](../py/acme/run.py)\n" +
		"[`hoard`](../web/main.ts)\n" +
		"[`Bucket.Nope`](../pkg/a/a.go)\n" +
		"[`gone.txt`](../gone.txt)\n"
	out, _ := checked(t, docsCharter(t, bad))
	for _, want := range []string{
		"docs/tour.md:3: rule 26: pkg/a/a.go declares no Drain",
		"docs/tour.md:4: rule 26: py/acme/run.py declares no Gauge",
		"docs/tour.md:5: rule 26: web/main.ts declares no hoard",
		"docs/tour.md:6: rule 26: pkg/a/a.go declares no Bucket.Nope",
		"docs/tour.md:7: rule 26: the code reference to gone.txt names a file that is not there",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("want\n  %s\ngot\n%s", want, out)
		}
	}
	if strings.Count(out, "rule 26") != 5 {
		t.Errorf("want five problems, got:\n%s", out)
	}
}

func TestGenKeepsADocsReferencesCurrent(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_H4B7, "title-kept-current", "target-kept-current", "line-kept-current", "bare-id-stays")
	root := docsCharter(t, "# Tour\n\nFirst [SW-010 an old title](wrong.md), bare [SW-010](wrong.md), and [`Refill`](../pkg/a/a.go).\n")
	got := readFile(t, root, "docs/tour.md")
	for _, want := range []string{
		"[SW-010: The bucket refills](../claims/refill.md#sw-010-the-bucket-refills)",
		"bare [SW-010](../claims/refill.md#sw-010-the-bucket-refills)",
		"[`Refill`](../pkg/a/a.go#L6)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("after gen, want %s in\n%s", want, got)
		}
	}
	// The title changes, the entry moves to another document, and the
	// declaration moves down its file.
	edit(t, root, "claims/refill.md", "## SW-010 The bucket refills\nStatus: active\n\nIt refills at the rate.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `refills`: it refills.\n\n", "")
	writeFile(t, root, "claims/bucket.md", "# Bucket\nCode: pkg/a\n\nThe bucket.\n\n## SW-010 The bucket refills at its rate\nStatus: active\n\nIt refills at the rate.\n\nWhy: a caller that waits is let through.\n\nScenarios:\n1. `refills`: it refills.\n")
	edit(t, root, "pkg/a/a.go", "package a\n", "package a\n\n// Rate is the rate.\nconst Rate = 1\n")
	var out, errOut bytes.Buffer
	if code := run([]string{"gen"}, root, &out, &errOut); code != 0 {
		t.Fatalf("gen: %s", errOut.String())
	}
	got = readFile(t, root, "docs/tour.md")
	for _, want := range []string{
		"[SW-010: The bucket refills at its rate](../claims/bucket.md#sw-010-the-bucket-refills-at-its-rate)",
		"bare [SW-010](../claims/bucket.md#sw-010-the-bucket-refills-at-its-rate)",
		"[`Refill`](../pkg/a/a.go#L9)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("after the change, want %s in\n%s", want, got)
		}
	}
}

func TestGenWritesEmbedsAndLeavesProseAlone(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_H4B7, "embed-written", "prose-untouched")
	prose := "# Tour\n\nA caller *waits*, then:\n\n```\n[SW-999](not a reference in a fence)\n```\n\n" +
		"~~~\n```\n[SW-999](still in the tilde fence)\n~~~\n\n    [SW-999](an indented code block)\n\nAnd ``[SW-999](x) in a double-backtick span``.\n"
	root := docsCharter(t, prose+"\n<!-- charter:embed SW-010 -->\nstale words\n<!-- /charter:embed -->\n\nAfter.\n")
	got := readFile(t, root, "docs/tour.md")
	want := prose + "\n<!-- charter:embed SW-010 -->\n> **[SW-010](../claims/refill.md#sw-010-the-bucket-refills) The bucket refills**\n>\n> It refills at the rate.\n<!-- /charter:embed -->\n\nAfter.\n"
	if got != want {
		t.Errorf("after gen the doc is\n%s\nwant\n%s", got, want)
	}
	// What is in a fence, an indented block or a code span is no reference.
	if out, code := checked(t, root); code != 0 {
		t.Errorf("a reference in code was read:\n%s", out)
	}
	// A second run writes nothing.
	old := time.Now().Add(-time.Hour).Truncate(time.Second)
	p := filepath.Join(root, "docs", "tour.md")
	if err := os.Chtimes(p, old, old); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"gen"}, root, &out, &errOut); code != 0 {
		t.Fatalf("gen: %s", errOut.String())
	}
	if st, err := os.Stat(p); err != nil || !st.ModTime().Equal(old) {
		t.Errorf("a second gen with nothing changed wrote the doc: %v", err)
	}
}

func curated(t *testing.T, root string) (string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run([]string{"curate", "-by", "pat", filepath.Join(root, "docs", "tour.md")}, root, &out, &errOut)
	return out.String() + errOut.String(), code
}

const tour = "# Tour\n\nThe bucket, [SW-010](x), refills in [`Refill`](../pkg/a/a.go),\n[`Meter.tick`](../py/acme/run.py) counts and [`spend`](../web/main.ts) spends.\n"

func TestCurateRecordsWhoAndWhatTheDocRestedOn(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_GXP7, "curate-records", "curate-touches-only-its-block")
	root := docsCharter(t, tour)
	before := readFile(t, root, "docs/tour.md")
	if out, code := curated(t, root); code != 0 {
		t.Fatalf("curate: exit %d %s", code, out)
	}
	got := readFile(t, root, "docs/tour.md")
	if !strings.HasPrefix(got, before) {
		t.Errorf("curating changed the doc above its block:\n%s", got)
	}
	block := got[len(before):]
	date := time.Now().UTC().Format("2006-01-02")
	for _, want := range []string{"\n<!-- charter:curated\nby: pat\ndate: " + date + "\n", "\nSW-010 ", "\npkg/a/a.go#Refill ", "\npy/acme/run.py#Meter.tick ", "\nweb/main.ts#spend ", "\n-->\n"} {
		if !strings.Contains(block, want) {
			t.Errorf("the block lacks %q:\n%s", want, block)
		}
	}
	if out, code := curated(t, root); code != 0 {
		t.Fatalf("curate again: exit %d %s", code, out)
	}
	if again := readFile(t, root, "docs/tour.md"); again != got || strings.Count(again, "charter:curated") != 1 {
		t.Errorf("curating again did not replace the block:\n%s", again)
	}
	// Prose written after the block is kept, and the new block follows it.
	writeFile(t, root, "docs/tour.md", got+"\nWritten after the block.\n")
	if out, code := curated(t, root); code != 0 {
		t.Fatalf("curate with prose after the block: exit %d %s", code, out)
	}
	after := readFile(t, root, "docs/tour.md")
	if !strings.HasPrefix(after, before) || !strings.Contains(after, "Written after the block.") || strings.Count(after, "charter:curated") != 1 || !strings.HasSuffix(after, "-->\n") {
		t.Errorf("curating dropped or misplaced what followed the block:\n%s", after)
	}
}

func TestCurateRefusesADocWithABrokenReference(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_GXP7, "curate-refuses-broken")
	root := docsCharter(t, "# Tour\n\n[SW-999](x)\n")
	before := readFile(t, root, "docs/tour.md")
	out, code := curated(t, root)
	if code != 1 || !strings.Contains(out, "docs/tour.md:3: rule 26: the reference names SW-999") {
		t.Errorf("curate of a broken doc: exit %d %s", code, out)
	}
	if readFile(t, root, "docs/tour.md") != before {
		t.Errorf("curate wrote a doc it refused")
	}
}

func TestLintNamesADocWhoseReferencesChanged(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_GXP7, "stale-on-entry-change", "changes-line-not-stale", "stale-on-go-change", "stale-on-python-change", "stale-on-typescript-change", "drafted-named", "stale-never-fails")
	root := docsCharter(t, tour)
	if out := warned(t, root); !strings.Contains(out, "docs/tour.md:1: warning: the doc is drafted") {
		t.Errorf("an uncurated doc is not named drafted:\n%s", out)
	}
	if out, code := checked(t, root); code != 0 {
		t.Errorf("a drafted doc fails the gate:\n%s", out)
	}
	if out, code := curated(t, root); code != 0 {
		t.Fatalf("curate: %s", out)
	}
	if out := warned(t, root); strings.Contains(out, "docs/tour.md") {
		t.Errorf("a doc curated a moment ago is named:\n%s", out)
	}
	// A Changes line alone, and whitespace, change nothing a doc rests on.
	edit(t, root, "claims/refill.md", "1. `refills`: it refills.\n", "1. `refills`: it refills.\n\nChanges:\n- 2026-11-27: noted.\n")
	edit(t, root, "pkg/a/a.go", "func Refill() {}", "func Refill()   {}")
	if out := warned(t, root); strings.Contains(out, "docs/tour.md") {
		t.Errorf("a Changes line or whitespace made the doc stale:\n%s", out)
	}
	edit(t, root, "claims/refill.md", "It refills at the rate.", "It refills at twice the rate.")
	edit(t, root, "pkg/a/a.go", "func Refill()   {}", "func Refill() { println() }")
	edit(t, root, "py/acme/run.py", "return 1", "return 2")
	edit(t, root, "web/main.ts", "return n - 1;", "return n - 2;")
	out := warned(t, root)
	for _, want := range []string{
		"docs/tour.md:3: warning: the doc is stale: SW-010 changed since it was curated",
		"docs/tour.md:3: warning: the doc is stale: pkg/a/a.go#Refill changed since it was curated",
		"docs/tour.md:4: warning: the doc is stale: py/acme/run.py#Meter.tick changed since it was curated",
		"docs/tour.md:4: warning: the doc is stale: web/main.ts#spend changed since it was curated",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("want\n  %s\ngot\n%s", want, out)
		}
	}
	if out, code := checked(t, root); code != 0 {
		t.Errorf("a stale doc fails the gate:\n%s", out)
	}
}

func TestBacklinksAreReadFromTheDocs(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_8MYR, "anchors-lists-docs")
	root := docsCharter(t, "# Tour\n\nOne\n[SW-010](x).\n")
	writeFile(t, root, "docs/more.md", "# More\n\n<!-- charter:embed SW-010 -->\n<!-- /charter:embed -->\n")
	var out, errOut bytes.Buffer
	if code := run([]string{"anchors", "SW-010"}, root, &out, &errOut); code != 0 {
		t.Fatalf("anchors: %s", errOut.String())
	}
	var got map[string]struct {
		Docs []struct {
			File string
			Line int
		}
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	ds := got["SW-010"].Docs
	if len(ds) != 2 || ds[0].File != "docs/more.md" || ds[0].Line != 3 || ds[1].File != "docs/tour.md" || ds[1].Line != 4 {
		t.Errorf("anchors gives the docs %+v", ds)
	}
}

func TestDiffListsTheDocsToReadAgain(t *testing.T) {
	claims.Verifies(t, claims.CH_261127_8MYR, "diff-lists-docs")
	root := docsCharter(t, "# Tour\n\n[SW-010](x) and [`Refill`](../pkg/a/a.go).\n")
	for _, args := range [][]string{{"init", "-q"}, {"add", "-A"}, {"-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-qm", "start"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"diff", "HEAD"}, root, &out, &errOut); code != 0 || out.Len() != 0 {
		t.Errorf("nothing changed: exit %d %q %s", code, out.String(), errOut.String())
	}
	// The entry's text changes with a Changes line, so diff names nothing
	// of the record, and the code the doc points at changes.
	edit(t, root, "claims/refill.md", "It refills at the rate.", "It refills at twice the rate.")
	edit(t, root, "claims/refill.md", "1. `refills`: it refills.\n", "1. `refills`: it refills.\n\nChanges:\n- 2026-11-27: twice the rate.\n")
	edit(t, root, "pkg/a/a.go", "func Refill() {}", "func Refill() { println() }")
	out.Reset()
	code := run([]string{"diff", "HEAD"}, root, &out, &errOut)
	for _, want := range []string{"docs/tour.md: docs: SW-010 changed since HEAD; read the doc again", "docs/tour.md: docs: pkg/a/a.go#Refill changed since HEAD; read the doc again"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("want %s in\n%s%s", want, out.String(), errOut.String())
		}
	}
	if code != 0 {
		t.Errorf("the docs to read again changed diff's exit to %d", code)
	}
	// Code that cannot be read leaves the docs unlisted and the status as
	// it stands.
	writeFile(t, root, "pkg/a/broken.go", "package a\n\nthis does not parse\n")
	out.Reset()
	errOut.Reset()
	if code := run([]string{"diff", "HEAD"}, root, &out, &errOut); code != 0 || !strings.Contains(errOut.String(), "the docs to read again are not listed") {
		t.Errorf("with unreadable code diff exits %d: %s%s", code, out.String(), errOut.String())
	}
}
