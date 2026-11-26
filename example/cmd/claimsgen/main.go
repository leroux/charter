// Command claimsgen is the charter's command line over ratelimit's charter,
// in claims/: the adopter's own main, passing its gate's options.
package main

import (
	"os"

	"example.com/ratelimit/claims/gate"
	"github.com/leroux/charter/cli"
)

func main() {
	cli.Claimsgen(os.Args[1:], cli.Charter{Name: "product", Dir: "claims", Options: gate.Options})
}
