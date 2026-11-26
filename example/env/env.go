// Package env names every environment variable ratelimit reads. The gate
// counts each constant as a surface some decision must account for.
package env

const (
	// Rate is the tokens a key earns per second, a decimal number.
	Rate = "RATELIMIT_RATE"
	// Burst is the most tokens a key may hold, a whole number.
	Burst = "RATELIMIT_BURST"
)
