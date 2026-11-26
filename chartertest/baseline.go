package chartertest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
)

// baselineLine is one surface no entry accounts for yet, with its reason.
type baselineLine struct {
	Key, Reason string
	Line        int
}

const baselineHeader = "# Surfaces no spec accounts for yet: key, a tab, reason. This file only\n# shrinks: the gate fails a line whose surface is accounted for or gone.\n"

// readBaseline parses the file. A missing file is an empty baseline. name
// is the display name used in problem messages, such as "claims/baseline.txt".
//
//charter:realizes CH-033
func readBaseline(p, name string) ([]baselineLine, []string, error) {
	data, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var lines []baselineLine
	var problems []string
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, reason, ok := strings.Cut(line, "\t")
		if !ok {
			problems = append(problems, fmt.Sprintf("%s:%d: rule 7: %q has no reason; a tab separates the key from the reason", name, i+1, line))
			lines = append(lines, baselineLine{Key: line, Line: i + 1})
			continue
		}
		lines = append(lines, baselineLine{Key: strings.TrimSpace(key), Reason: strings.TrimSpace(reason), Line: i + 1})
	}
	return lines, problems, nil
}

// writeBaseline writes the lines sorted by key.
func writeBaseline(p string, lines []baselineLine) error {
	sort.Slice(lines, func(i, j int) bool { return lines[i].Key < lines[j].Key })
	var b strings.Builder
	b.WriteString(baselineHeader)
	for _, l := range lines {
		fmt.Fprintf(&b, "%s\t%s\n", l.Key, l.Reason)
	}
	return os.WriteFile(p, []byte(b.String()), 0o600)
}
