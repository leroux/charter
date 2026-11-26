// Package gate holds the options ratelimit's charter is checked with, so
// the gate test and claimsgen read the same registry. It cannot live in
// claims, where the gate allows only entries.
package gate

import "github.com/leroux/charter/chartertest"

// Options configures the gate. Root is left empty: it resolves to the
// module above the caller.
var Options = chartertest.Options{
	Entries:     "example.com/ratelimit/claims",
	EnvPackages: []string{"example.com/ratelimit/env"},
	// A string that looks like one of ratelimit's variables, outside env,
	// is a variable nobody accounted for.
	EnvLiterals: []string{`^RATELIMIT_[A-Z0-9_]+$`},
	Regenerate:  "go run ./cmd/claimsgen",
	// A need is verified from the caller's side: the command, driven as a
	// person or a script would drive it.
	StakeholderTests: []string{"cmd/ratelimit/*_test.go"},
	Digest:           "DIGEST.md",
}
