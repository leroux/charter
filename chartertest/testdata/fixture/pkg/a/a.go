//go:build unix

// Package a reads the data directory.
//
//charter:realizes SW-001, CON-001, CON-002
package a

import (
	"os"

	"example.com/fixture/env"
)

// DataDir returns the configured data directory.
func DataDir() string { return os.Getenv(env.DataDir) }
