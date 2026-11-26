// Package gate holds the options gjson's charter is checked with, so the
// gate test and claimsgen read the same registry. It cannot live in
// claims, where the gate allows only entries.
package gate

import "github.com/leroux/charter/chartertest"

// Options configures the gate. Root is left empty: it resolves to the
// module above the caller.
var Options = chartertest.Options{
	Entries:    "github.com/tidwall/gjson/claims",
	Regenerate: "go run ./cmd/claimsgen",
	Digest:     "DIGEST.md",
}
