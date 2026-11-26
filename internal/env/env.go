// Package env names every environment variable the charter command and the
// generated files read. The gate counts each constant as a surface some
// entry must account for.
package env

const (
	// Events is the file each generated verifies appends a line to, per
	// citation and outcome, once its test ends.
	Events = "CHARTER_EVENTS"
	// UpdateBaseline makes the gate rewrite the baseline with every
	// surface no entry accounts for, instead of failing.
	UpdateBaseline = "CHARTER_UPDATE_BASELINE"
)
