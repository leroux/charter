// Package gate holds the options the charter's own charter is checked with:
// the decisions in the area documents of claims/, realized by the charter
// library (the root package, chartertest) and the commands under cmd/, and
// verified by their tests. It cannot live in the entries package, where
// rule 1 allows only observations.
package gate

import "github.com/leroux/charter/chartertest"

// Options configures the charter's own gate. The walk covers the tooling:
// the module's Go code. The viewer and the example
// adopter are not the tooling; the example is a module of its own and is
// skipped as every nested module is.
var Options = chartertest.Options{
	Entries: "github.com/leroux/charter/claims",
	// The gate's own environment variables are the exported string
	// constants of chartertest, read there.
	EnvPackages: []string{"github.com/leroux/charter/chartertest"},
	EnvLiterals: []string{`^CHARTER_[A-Z0-9_]+$`},
	Regenerate:  "make claims",
	// The report runs the tooling's tests.
	TestPackages: []string{"./..."},
	// The attr and mutate probes cite specs they declare themselves.
	Ignore: []string{"chartertest/internal", "claims-viewer"},
}
