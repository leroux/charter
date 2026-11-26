module example.com/ratelimit

go 1.27.1

require example.com/ratelimit/claims v0.0.0

// The claims module is ratelimit's record as Go, written by charter gen.
replace example.com/ratelimit/claims => ./claims
