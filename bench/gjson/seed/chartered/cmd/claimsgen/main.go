// Command claimsgen is the charter's command line over gjson's charter, in
// claims/: the adopter's own main, passing its gate's options.
package main

import (
	"os"

	"github.com/leroux/charter/cli"
	"github.com/tidwall/gjson/claims/gate"
)

func main() {
	cli.Claimsgen(os.Args[1:], cli.Charter{Name: "product", Dir: "claims", Options: gate.Options})
}
