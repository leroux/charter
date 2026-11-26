package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/leroux/charter/chartertest"
)

// Claimsmap exports a charter as the JSON claims-viewer reads: the area
// documents, the waivers, and every spec and observation with its
// relations, scenarios, anchors and verifying tests.
//
//	claimsmap [-charter NAME] [-format json] [-o claims.json]
//
// Without -charter it exports the first of the charters it is given.
// -o - writes to stdout. Problems the gate would report while loading are
// printed to stderr, and the output is written anyway. It returns the exit
// code.
func Claimsmap(args []string, stderr io.Writer, charters ...Charter) int {
	fs := flag.NewFlagSet("claimsmap", flag.ContinueOnError)
	fs.SetOutput(stderr)
	out := fs.String("o", "claims.json", "output file, or - for stdout")
	format := fs.String("format", "json", "json, the registry for claims-viewer")
	which := fs.String("charter", "", "the charter to export, by name: "+names(charters)+"; empty means the first")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	say := func(format string, a ...any) { _, _ = fmt.Fprintf(stderr, "claimsmap: "+format+"\n", a...) }
	if *format != "json" {
		say("unknown format %q; want json", *format)
		return 1
	}
	c, err := pick(charters, *which)
	if err != nil {
		say("%v", err)
		return 1
	}
	opts := c.Options
	reg, err := chartertest.Read(opts)
	if err != nil {
		say("%v", err)
		return 1
	}
	for _, p := range reg.Problems {
		say("warning: %s", p)
	}
	var buf bytes.Buffer
	if err := writeJSON(&buf, reg, c.Name, source(c.Dir), root()); err != nil {
		say("%v", err)
		return 1
	}
	if *out == "-" {
		_, err = os.Stdout.Write(buf.Bytes())
	} else {
		err = os.WriteFile(*out, buf.Bytes(), 0o600)
	}
	if err != nil {
		say("%v", err)
		return 1
	}
	if *out != "-" {
		say("wrote %s (%d entries)", *out, len(reg.Entries))
	}
	return 0
}

func git(args ...string) string {
	b, err := exec.Command("git", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// root is the checkout's absolute path, or "" outside git.
func root() string { return git("rev-parse", "--show-toplevel") }

// source names the checkout the registry was read from, such as
// "claims @ main 1a2b3c4", or the directory alone outside git.
func source(dir string) string {
	sha := git("rev-parse", "--short", "HEAD")
	if sha == "" {
		return dir
	}
	s := dir + " @ "
	if b := git("branch", "--show-current"); b != "" {
		s += b + " "
	}
	s += sha
	if git("status", "--porcelain", "--", dir) != "" {
		s += " + uncommitted changes"
	}
	return s
}
