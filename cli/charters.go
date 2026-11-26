package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leroux/charter/chartertest"
)

// Charter is one charter a command can work on: a repository may keep more
// than one, such as its product's and a library's own.
type Charter struct {
	// Name selects the charter with -charter, such as "product" or "self".
	Name string
	// Dir is the entries directory, slash-relative to the repository root,
	// such as "claims"; the export names its source by it.
	Dir string
	// Options are the gate's options for the charter.
	Options chartertest.Options
}

// pick returns the charter named, or the first when name is empty, with
// its root resolved.
//
//charter:realizes CH-040
func pick(charters []Charter, name string) (Charter, error) {
	if len(charters) == 0 {
		return Charter{}, errors.New("no charter is configured")
	}
	c, found := charters[0], name == ""
	for _, cand := range charters {
		if !found && cand.Name == name {
			c, found = cand, true
		}
	}
	if !found {
		return Charter{}, fmt.Errorf("-charter is %s, not %q", names(charters), name)
	}
	if c.Options.Root == "" {
		root, err := workingRoot()
		if err != nil {
			return Charter{}, err
		}
		c.Options.Root = root
	}
	return c, nil
}

// workingRoot is the module root above the working directory. Left to
// chartertest, an empty root would resolve above the calling file, which
// for a command is this package, in the tooling's module, not the
// adopter's.
//
//charter:realizes CH-040
func workingRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for d := wd; ; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d, nil
		}
		if filepath.Dir(d) == d {
			return "", fmt.Errorf("no go.mod above the working directory %s; run the command inside the repository", wd)
		}
	}
}

// names lists the charters' names for a flag's help and its error.
func names(charters []Charter) string {
	n := make([]string, len(charters))
	for i, c := range charters {
		n[i] = fmt.Sprintf("%q (%s)", c.Name, c.Dir)
	}
	return strings.Join(n, " or ")
}
