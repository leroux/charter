// Command claimsgen is the charter's command line over the charter's own
// charter, in claims/. See github.com/leroux/charter/cli.Claimsgen for the
// verbs; an adopter writes a main of its own like this one.
package main

import (
	"os"

	"github.com/leroux/charter/claims/gate"
	"github.com/leroux/charter/cli"
)

func main() {
	cli.Claimsgen(os.Args[1:], cli.Charter{Name: "self", Dir: "claims", Options: gate.Options})
}
