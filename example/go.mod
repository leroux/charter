module example.com/ratelimit

go 1.27.1

require github.com/leroux/charter v0.0.0

// The example builds against the charter in this checkout, so a change to
// the tooling is exercised by an adopter in the same pull request.
replace github.com/leroux/charter => ../
