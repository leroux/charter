// Command ratelimit reads one caller key a line from standard input and
// answers each with "allow KEY" or "deny KEY" on standard output, by the
// rate and burst RATELIMIT_RATE and RATELIMIT_BURST configure.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"example.com/ratelimit/env"
	"example.com/ratelimit/limiter"
)

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, os.Getenv, limiter.SystemClock{}))
}

// The defaults when the environment names none: one request a second, five
// at once.
const (
	defaultRate  = 1.0
	defaultBurst = 5
)

// config reads the rate and the burst from the environment, each its
// default when unset.
//
//charter:realizes SW-261126-CMTG
func config(getenv func(string) string) (rate float64, burst int, err error) {
	rate, burst = defaultRate, defaultBurst
	if v := getenv(env.Rate); v != "" {
		if rate, err = strconv.ParseFloat(v, 64); err != nil {
			return 0, 0, fmt.Errorf("%s=%q is not a number", env.Rate, v)
		}
	}
	if v := getenv(env.Burst); v != "" {
		if burst, err = strconv.Atoi(v); err != nil {
			return 0, 0, fmt.Errorf("%s=%q is not a whole number", env.Burst, v)
		}
	}
	if rate <= 0 {
		return 0, 0, fmt.Errorf("%s=%v must be above zero", env.Rate, rate)
	}
	if burst < 1 {
		return 0, 0, fmt.Errorf("%s=%d must be at least one", env.Burst, burst)
	}
	return rate, burst, nil
}

func run(in io.Reader, out, errOut io.Writer, getenv func(string) string, clock limiter.Clock) int {
	rate, burst, err := config(getenv)
	if err != nil {
		_, _ = fmt.Fprintln(errOut, "ratelimit:", err)
		return 2
	}
	l, err := limiter.New(rate, burst, clock)
	if err != nil {
		_, _ = fmt.Fprintln(errOut, "ratelimit:", err)
		return 2
	}
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		key := strings.TrimSpace(sc.Text())
		if key == "" {
			continue
		}
		verdict := "deny"
		if l.Allow(key) {
			verdict = "allow"
		}
		_, _ = fmt.Fprintln(out, verdict, key)
	}
	return 0
}
