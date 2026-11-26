module github.com/leroux/charter

go 1.27.1

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/leroux/charter/claims v0.0.0
)

// The claims module is the charter's own record as Go, written by charter gen.
replace github.com/leroux/charter/claims => ./claims
