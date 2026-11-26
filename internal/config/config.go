// Package config finds a charter's claims/charter.toml and reads what it says
// about the tree: the languages the charter reads and how their tests run,
// the patterns that scope the walk, and the surface kinds it declares.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// Dir is the directory, under the root, that holds a charter's record.
const Dir = "claims"

// File is the configuration's name inside Dir.
const File = "charter.toml"

// Language is how a charter reads one language and runs its tests.
type Language struct {
	// Dir is where the language's test command runs, relative to the root.
	Dir string `toml:"dir"`
	// Test is the shell command that runs the language's tests.
	Test string `toml:"test"`
	// Env is the constants module naming the environment variables the
	// language reads: a package directory for Go, a file for Python and
	// TypeScript, relative to the root. Empty means none.
	Env string `toml:"env"`
}

// Go is the Go language's table: a Language and the import path of the
// module the generated symbols live in.
type Go struct {
	Language
	// Module is the import path of the claims module charter gen writes.
	Module string `toml:"module"`
}

// TypeScript is the TypeScript language's table: a Language, the
// typescript package the reader parses with and the specifier tests import
// the generated module by.
type TypeScript struct {
	Language
	// Compiler is the typescript npm package directory, relative to the root.
	Compiler string `toml:"compiler"`
	// Import is the specifier tests import claims/index.ts by.
	Import string `toml:"import"`
}

// Surface is a declared surface kind: each line its command prints is a key.
type Surface struct {
	Kind    string `toml:"kind"`
	Command string `toml:"command"`
}

// Config is a charter's claims/charter.toml, with Root resolved.
type Config struct {
	// Root is the absolute directory that holds claims/.
	Root string `toml:"-"`
	// Path is the configuration file's absolute path.
	Path string `toml:"-"`

	Name             string      `toml:"name"`
	Go               *Go         `toml:"go"`
	Python           *Language   `toml:"python"`
	TypeScript       *TypeScript `toml:"typescript"`
	StakeholderTests []string    `toml:"stakeholder-tests"`
	Probes           []string    `toml:"probes"`
	Ignore           []string    `toml:"ignore"`
	Only             []string    `toml:"only"`
	EnvPatterns      []string    `toml:"env-patterns"`
	Digest           string      `toml:"digest"`
	Surfaces         []Surface   `toml:"surface"`

	// lines maps a table's or a key's dotted name to its line, for problems.
	lines map[string]int
}

// Line returns the line of the table or key with the dotted name, such as
// "python" or "surface", or 1 when it is not written.
func (c *Config) Line(name string) int {
	if n, ok := c.lines[name]; ok {
		return n
	}
	return 1
}

// Rel returns the configuration file's path relative to the root.
func (c *Config) Rel() string { return Dir + "/" + File }

// ErrNotFound is wrapped by Find when no claims/charter.toml is above the
// directory it started from.
var ErrNotFound = errors.New("no claims/charter.toml")

// Find walks up from dir to the first directory holding claims/charter.toml
// and reads it; that directory is the root.
//
//charter:realizes CH-261126-33GK
func Find(dir string) (*Config, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for d := abs; ; d = filepath.Dir(d) {
		p := filepath.Join(d, Dir, File)
		if _, err := os.Stat(p); err == nil {
			return Load(d)
		}
		if filepath.Dir(d) == d {
			return nil, fmt.Errorf("%w in %s or any directory above it", ErrNotFound, abs)
		}
	}
}

// Load reads root/claims/charter.toml. An unknown key, a value of the wrong
// type and a language table without its test command are errors naming the
// line.
func Load(root string) (*Config, error) {
	p := filepath.Join(root, Dir, File)
	src, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	c, err := Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s/%s:%w", Dir, File, err)
	}
	c.Root, c.Path = root, p
	return c, nil
}

// Parse decodes a configuration's text. Its errors open with the line,
// "12: ...", for the caller to prefix with the file.
func Parse(src string) (*Config, error) {
	c := &Config{lines: keyLines(src)}
	md, err := toml.Decode(src, c)
	if err != nil {
		var pe toml.ParseError
		if errors.As(err, &pe) {
			return nil, fmt.Errorf("%d: %s", pe.Position.Line, pe.Message)
		}
		// A value of the wrong type is reported as "toml: line N (last key
		// "k"): ..."; the line opens the problem as for any other.
		if m := typeErrLine.FindStringSubmatch(err.Error()); m != nil {
			return nil, fmt.Errorf("%s: %s", m[1], m[2])
		}
		return nil, fmt.Errorf("1: %v", err)
	}
	var problems []string
	for _, k := range md.Undecoded() {
		name := k.String()
		problems = append(problems, fmt.Sprintf("%d: unknown key %q; the keys are name, go, python, typescript, stakeholder-tests, probes, ignore, only, env-patterns, digest and surface", c.Line(name), name))
	}
	for _, l := range []struct {
		name string
		lang *Language
	}{{"go", goLang(c.Go)}, {"python", c.Python}, {"typescript", tsLang(c.TypeScript)}} {
		if l.lang != nil && strings.TrimSpace(l.lang.Test) == "" {
			problems = append(problems, fmt.Sprintf("%d: [%s] has no test command; write test = \"...\", the command that runs its tests", c.Line(l.name), l.name))
		}
	}
	if c.Go != nil && c.Go.Module == "" {
		problems = append(problems, fmt.Sprintf("%d: [go] has no module; write module = \"...\", the import path of the generated claims module", c.Line("go")))
	}
	if c.TypeScript != nil && c.TypeScript.Compiler == "" {
		problems = append(problems, fmt.Sprintf("%d: [typescript] has no compiler; write compiler = \"...\", the typescript package directory", c.Line("typescript")))
	}
	if len(problems) > 0 {
		return nil, errors.New(strings.Join(problems, "\n"+Dir+"/"+File+":"))
	}
	return c, nil
}

func goLang(g *Go) *Language {
	if g == nil {
		return nil
	}
	return &g.Language
}

func tsLang(t *TypeScript) *Language {
	if t == nil {
		return nil
	}
	return &t.Language
}

var (
	typeErrLine = regexp.MustCompile(`^toml: line (\d+) (.*)$`)
	tableLine   = regexp.MustCompile(`^\s*\[\[?\s*([A-Za-z0-9_.-]+)\s*\]\]?`)
	keyLine     = regexp.MustCompile(`^\s*([A-Za-z0-9_-]+)\s*=`)
)

// keyLines maps each table's dotted name, and each key's dotted name under
// the table it is written in, to the first line it appears on.
func keyLines(src string) map[string]int {
	lines := map[string]int{}
	table := ""
	for i, line := range strings.Split(src, "\n") {
		n := i + 1
		if m := tableLine.FindStringSubmatch(line); m != nil {
			table = m[1]
			if _, ok := lines[table]; !ok {
				lines[table] = n
			}
			continue
		}
		if m := keyLine.FindStringSubmatch(line); m != nil {
			name := m[1]
			if table != "" {
				name = table + "." + name
			}
			if _, ok := lines[name]; !ok {
				lines[name] = n
			}
		}
	}
	return lines
}
