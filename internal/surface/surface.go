// Package surface enumerates a charter's declared surface kinds, by running
// each kind's command from the root, and reads and writes the baseline of
// surfaces that predate the charter. Environment variables, the one built-in
// kind, come from the readers' facts.
package surface

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/leroux/charter/internal/config"
)

// kindForm is a declared kind's name: a lower-case word.
var kindForm = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// Declared runs each declared kind's command from the root and returns the
// keys, kind:line, and the problems: a command that fails or prints a line
// twice, and a kind declared twice, named env or not a lower-case word.
//
//charter:realizes CH-261126-NTRR
func Declared(c *config.Config) ([]string, []string) {
	var keys, problems []string
	seen := map[string]bool{}
	for i, s := range c.Surfaces {
		at := fmt.Sprintf("%s:%d", c.Rel(), c.SurfaceLine(i))
		switch {
		case !kindForm.MatchString(s.Kind):
			problems = append(problems, fmt.Sprintf("%s: rule 9: the surface kind %q is not a lower-case word", at, s.Kind))
			continue
		case s.Kind == "env":
			problems = append(problems, fmt.Sprintf("%s: rule 9: the surface kind env is built in; name a declared kind otherwise", at))
			continue
		case seen[s.Kind]:
			problems = append(problems, fmt.Sprintf("%s: rule 9: the surface kind %s is declared twice", at, s.Kind))
			continue
		}
		seen[s.Kind] = true
		cmd := exec.Command("sh", "-c", s.Command)
		cmd.Dir = c.Root
		var out, errOut bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errOut
		if err := cmd.Run(); err != nil {
			problems = append(problems, fmt.Sprintf("%s: rule 9: the %s command %q failed (%v): %s", at, s.Kind, s.Command, err, strings.TrimSpace(errOut.String()+out.String())))
			continue
		}
		printed := map[string]bool{}
		for _, line := range strings.Split(out.String(), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if printed[line] {
				problems = append(problems, fmt.Sprintf("%s: rule 9: the %s command printed %s twice", at, s.Kind, line))
				continue
			}
			printed[line] = true
			keys = append(keys, s.Kind+":"+line)
		}
	}
	return keys, problems
}

// Line is one surface of the baseline: no entry accounts for it yet.
type Line struct {
	Key, Reason string
	Line        int
}

const header = "# Surfaces no entry accounts for yet: key, a tab, reason. This file only\n# shrinks: the gate fails a line whose surface is accounted for or gone.\n"

// InitialReason is the reason a rebuilt baseline gives a surface it adds.
const InitialReason = "unclaimed at baseline init"

// ReadBaseline reads the baseline at p; a missing file is an empty one.
// name is how problems name the file.
//
//charter:realizes CH-033
func ReadBaseline(p, name string) ([]Line, []string, error) {
	data, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var lines []Line
	var problems []string
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, reason, ok := strings.Cut(line, "\t")
		if !ok {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 7: %q has no reason; a tab separates the key from the reason", name, i+1, line))
			lines = append(lines, Line{Key: line, Line: i + 1})
			continue
		}
		lines = append(lines, Line{Key: strings.TrimSpace(key), Reason: strings.TrimSpace(reason), Line: i + 1})
	}
	return lines, problems, nil
}

// Rebuild keeps every line still unaccounted and present, with its reason,
// and adds every other unaccounted surface.
//
//charter:realizes CH-033
func Rebuild(old []Line, surfaces map[string]bool, accounted map[string]string) []Line {
	reasons := map[string]string{}
	for _, l := range old {
		if l.Reason != "" {
			reasons[l.Key] = l.Reason
		}
	}
	var keys []string
	for k := range surfaces {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []Line
	for _, k := range keys {
		if accounted[k] != "" {
			continue
		}
		reason := reasons[k]
		if reason == "" {
			reason = InitialReason
		}
		out = append(out, Line{Key: k, Reason: reason})
	}
	return out
}

// WriteBaseline writes the lines sorted by key.
func WriteBaseline(p string, lines []Line) error {
	sort.Slice(lines, func(i, j int) bool { return lines[i].Key < lines[j].Key })
	var b strings.Builder
	b.WriteString(header)
	for _, l := range lines {
		fmt.Fprintf(&b, "%s\t%s\n", l.Key, l.Reason)
	}
	return os.WriteFile(p, []byte(b.String()), 0o644)
}
