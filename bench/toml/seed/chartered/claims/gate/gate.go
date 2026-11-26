// Package gate holds the options the decoder's charter is checked with, so
// the gate test and claimsgen read the same registry. It cannot live in
// claims, where the gate allows only entries.
package gate

import "github.com/leroux/charter/chartertest"

// Options configures the gate. Root is left empty: it resolves to the
// module above the caller.
var Options = chartertest.Options{
	Entries:    "example.com/toml/claims",
	Regenerate: "go run ./cmd/claimsgen",
	// A need is verified from the config-file user's side: the decoder
	// command, fed a document on stdin and read back.
	StakeholderTests: []string{"cmd/toml-decoder/*_test.go"},
	Digest:           "DIGEST.md",
}
