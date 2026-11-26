module gjsonoracle

go 1.24

require github.com/tidwall/gjson v0.0.0

require (
	github.com/stretchr/testify v1.12.1 // indirect
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v1.0.0 // indirect
)

// ref/ is gjson at the last chosen upstream commit, 8e88233: the reference
// the grader runs.
replace github.com/tidwall/gjson => ./ref
