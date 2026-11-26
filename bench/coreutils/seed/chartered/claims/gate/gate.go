// Package gate holds the options coreutils' charter is checked with, so
// the gate test and claimsgen read the same registry. It cannot live in
// claims, where the gate allows only entries.
package gate

import "github.com/leroux/charter/chartertest"

// Options configures the gate. Root is left empty: it resolves to the
// module above the caller.
var Options = chartertest.Options{
	Entries:    "example.com/coreutils/claims",
	Regenerate: "go run ./cmd/claimsgen",
	// A need is verified from the script author's side: a tool's command,
	// given argv and input and read back.
	StakeholderTests: []string{"cmd/*/*_test.go"},
	Digest:           "DIGEST.md",
}
