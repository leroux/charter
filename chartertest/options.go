package chartertest

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// The charter's files in the entries package directory, beside the area
// documents: every .md file there is an area document but README.md, the
// digest, questions.md, policy.md, glossary.md and audit.md, and needs.md
// holds the needs.
const (
	// waiversFile has one line per spec id or scenario (SW-118/handle), a
	// tab, and the reason its missing test is accepted.
	waiversFile = "waivers.txt"
	// baselineFile lists the surfaces that predate the charter.
	baselineFile = "baseline.txt"
	// readingOrderFile lists the area documents one a line in the order the
	// digest reads them; a missing file means alphabetical order.
	readingOrderFile = "reading-order.txt"
)

// Options configures Check, Problems and Read. Entries is required.
type Options struct {
	// Root is the module root. Empty means the nearest directory holding a
	// go.mod above the calling file or, when there is none, above the
	// working directory, so a script outside the module works when run from
	// inside it.
	Root string
	// Entries is the import path of the entries package: the package that
	// holds the generated Go, the observations, and the area documents.
	Entries string
	// Only restricts every walk to these slash paths relative to Root, files
	// or directories. Empty means the whole tree. A charter that covers one
	// part of a module, such as the charter's own, names its part here.
	Only []string
	// EnvPackages are import paths whose exported string constants name
	// environment variables. Their values are the env surfaces.
	EnvPackages []string
	// EnvReaders are "importpath.Func" calls whose first argument must name
	// a constant of an env package. Empty means os.Getenv and os.LookupEnv.
	EnvReaders []string
	// EnvLiterals are regular expressions. A string literal matching one,
	// outside the env packages and outside test files, is an undeclared
	// environment variable.
	EnvLiterals []string
	// OpenAPI lists OpenAPI JSON files, relative to Root, whose paths and
	// methods are the operation surfaces.
	OpenAPI []string
	// Events is a JSON file, relative to Root, whose top-level keys are the
	// event kind surfaces. Empty means no event surfaces.
	Events string
	// Ignore are patterns skipped by every walk, matched with path.Match
	// against the slash path relative to Root and against the base name. A
	// pattern without a slash matches a base name at any depth, so name a
	// path such as web/node_modules to match one place. .git,
	// node_modules, vendor and testdata are always skipped, as is any
	// directory below Root whose name starts with . or _ or that holds its
	// own go.mod.
	Ignore []string
	// ProbeTag is the build tag every probe file carries. A test file whose
	// constraint includes it is a probe: it cites observations and nothing
	// else, and the tag is not a surface. Empty means no observations.
	ProbeTag string
	// Enumerators add surface kinds beyond the four built in.
	Enumerators []Enumerator
	// Specs is the directory of spec files in the earlier form, one file per
	// spec, relative to the entries package directory, read while they are
	// migrated into area documents. Empty means specs; a missing directory
	// means none.
	Specs string
	// StakeholderTests are patterns, matched with path.Match against the
	// slash path relative to Root and against the base name, of the test
	// files on the stakeholder's side: the command line, the server's routes,
	// the browser. A need is verified from those files only.
	// Empty means any test file verifies a need.
	StakeholderTests []string
	// TestPackages are the Go package patterns the report runs, relative
	// to Root, such as ./... Empty means ./...
	TestPackages []string
	// Regenerate is the command a stale generated file or TypeScript module
	// tells a person to run, such as "make claims". Empty names the update
	// environment variables.
	Regenerate string
	// TypeScript enables citations from TypeScript tests. Nil means none:
	// no TypeScript file is walked and node is never run.
	TypeScript *TypeScript
	// Digest is the digest file, relative to the entries package directory,
	// such as DIGEST.md: the charter rendered for a reader, written by
	// Generate and compared by the gate. Empty means none is kept.
	Digest string
}

// Enumerator produces the keys of one surface kind from the tree.
type Enumerator struct {
	// Name appears in rule 9 when the enumerator finds nothing.
	Name string
	// Keys returns every surface key of the kind, such as "metric:runs_total".
	Keys func(root string) ([]string, error)
}

// TypeScript configures citations from TypeScript tests: the module the
// gate renders from the entries, how tests import it, and the compiler the
// extractor parses with.
type TypeScript struct {
	// Module is the generated module, slash-relative to Root, such as
	// "claims/index.ts".
	Module string
	// Import is the specifier tests import the module by, such as
	// "@example/claims".
	Import string
	// Compiler is the typescript npm package directory, slash-relative to
	// Root or absolute, such as "web/node_modules/typescript".
	Compiler string
	// Install is the command that installs Compiler, named when it is
	// missing, such as "pnpm install". Empty names only the package.
	Install string
	// Node is the node binary. Empty means "node" on PATH.
	Node string
}

var defaultReaders = []string{"os.Getenv", "os.LookupEnv"}

// alwaysIgnored are base names every walk skips. loadTree also skips hidden
// and underscore directories and nested modules, as the go tool does.
var alwaysIgnored = map[string]bool{".git": true, "node_modules": true, "vendor": true, "testdata": true}

// resolved is Options with defaults applied and derived values computed.
type resolved struct {
	Options
	modulePath string
	entriesDir string            // slash path relative to Root
	areasDir   string            // slash path relative to Root
	specsDir   string            // slash path relative to Root
	envDirs    map[string]string // import path -> slash path relative to Root
	readers    map[string]bool
	literals   []*regexp.Regexp
}

//charter:realizes CH-037
func (o Options) resolve(callerDir string) (*resolved, error) {
	r := &resolved{Options: o}
	if r.Root == "" {
		root, err := defaultRoot(callerDir)
		if err != nil {
			return nil, err
		}
		r.Root = root
	}
	if r.Entries == "" {
		return nil, errors.New("chartertest: Options.Entries is required")
	}
	mod, err := modulePath(r.Root)
	if err != nil {
		return nil, err
	}
	r.modulePath = mod
	if r.entriesDir, err = packageDir(mod, r.Entries); err != nil {
		return nil, err
	}
	for i, o := range r.Only {
		r.Only[i] = strings.Trim(path.Clean(o), "/")
	}
	r.envDirs = map[string]string{}
	for _, p := range r.EnvPackages {
		d, err := packageDir(mod, p)
		if err != nil {
			return nil, err
		}
		r.envDirs[p] = d
	}
	readers := r.EnvReaders
	if len(readers) == 0 {
		readers = defaultReaders
	}
	r.readers = map[string]bool{}
	for _, rd := range readers {
		r.readers[rd] = true
	}
	for _, pat := range r.EnvLiterals {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("chartertest: EnvLiterals %q: %w", pat, err)
		}
		r.literals = append(r.literals, re)
	}
	if ts := r.TypeScript; ts != nil {
		switch {
		case ts.Module == "":
			return nil, errors.New("chartertest: Options.TypeScript.Module is required")
		case ts.Import == "":
			return nil, errors.New("chartertest: Options.TypeScript.Import is required")
		case ts.Compiler == "":
			return nil, errors.New("chartertest: Options.TypeScript.Compiler is required")
		}
		cp := *ts
		if cp.Node == "" {
			cp.Node = "node"
		}
		if !filepath.IsAbs(cp.Compiler) {
			cp.Compiler = filepath.Join(r.Root, filepath.FromSlash(cp.Compiler))
		}
		r.TypeScript = &cp
	}
	r.areasDir = r.entriesDir
	if r.Specs == "" {
		r.Specs = "specs"
	}
	r.specsDir = path.Join(r.entriesDir, r.Specs)
	return r, nil
}

// defaultRoot is the module root above callerDir or, failing that, above
// the working directory.
func defaultRoot(callerDir string) (string, error) {
	if root, err := moduleRoot(callerDir); err == nil {
		return root, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("chartertest: no go.mod above %s, and no working directory: %w", callerDir, err)
	}
	if root, err := moduleRoot(wd); err == nil {
		return root, nil
	}
	return "", fmt.Errorf("chartertest: no go.mod above %s or above the working directory %s", callerDir, wd)
}

// moduleRoot walks up from dir to the nearest directory holding a go.mod.
func moduleRoot(dir string) (string, error) {
	for d := dir; ; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d, nil
		}
		if filepath.Dir(d) == d {
			return "", fmt.Errorf("chartertest: no go.mod above %s", dir)
		}
	}
}

// modulePath reads the module line of root's go.mod.
func modulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("chartertest: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", fmt.Errorf("chartertest: no module line in %s/go.mod", root)
}

// packageDir maps an import path inside the module to a slash path relative
// to the root.
func packageDir(modulePath, importPath string) (string, error) {
	if importPath == modulePath {
		return ".", nil
	}
	rest, ok := strings.CutPrefix(importPath, modulePath+"/")
	if !ok {
		return "", fmt.Errorf("chartertest: %s is not inside module %s", importPath, modulePath)
	}
	return rest, nil
}

// ignored reports whether a slash path relative to Root is skipped.
//
//charter:realizes CH-036
func (r *resolved) ignored(rel string) bool {
	return alwaysIgnored[path.Base(rel)] || matchesAny(r.Ignore, rel)
}

// inScope reports whether a slash path relative to Root is within Only: the
// path itself or something under it is named, or, for a directory, it lies
// on the way to something named, so the walk descends.
//
//charter:realizes CH-036
func (r *resolved) inScope(rel string, isDir bool) bool {
	if len(r.Only) == 0 {
		return true
	}
	for _, o := range r.Only {
		if rel == o || strings.HasPrefix(rel, o+"/") {
			return true
		}
		if isDir && strings.HasPrefix(o, rel+"/") {
			return true
		}
	}
	return false
}

// stakeholderSide reports whether a test file, slash-relative to Root, is on
// the stakeholder's side: it matches a StakeholderTests pattern, or none is
// configured.
//
//charter:realizes CH-031
func (r *resolved) stakeholderSide(rel string) bool {
	return len(r.StakeholderTests) == 0 || matchesAny(r.StakeholderTests, rel)
}

// matchesAny reports whether a slash path, or its base name, matches one of
// the path.Match patterns.
func matchesAny(patterns []string, rel string) bool {
	base := path.Base(rel)
	for _, pat := range patterns {
		if ok, _ := path.Match(pat, rel); ok {
			return true
		}
		if ok, _ := path.Match(pat, base); ok {
			return true
		}
	}
	return false
}

// callerDir is the directory of the source file `skip` frames above
// callerDir's caller: 1 is the caller's caller, which is what a public
// entry point passes to find the test file that called it.
func callerDir(skip int) string {
	_, file, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
