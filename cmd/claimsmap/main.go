// Command claimsmap exports the charter's own charter, in claims/, as the
// JSON claims-viewer reads. See github.com/leroux/charter/cli.Claimsmap.
package main

import (
	"os"

	"github.com/leroux/charter/claims/gate"
	"github.com/leroux/charter/cli"
)

func main() {
	os.Exit(cli.Claimsmap(os.Args[1:], os.Stderr, cli.Charter{Name: "self", Dir: "claims", Options: gate.Options}))
}
