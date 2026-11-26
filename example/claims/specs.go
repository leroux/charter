// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: go run ./cmd/claimsgen

package claims

import "github.com/leroux/charter"

// STK_261126_H2SW is STK-261126-H2SW: A caller gets a burst at once, then its sustained rate, and no more
//
// A caller that has been quiet may make a burst of requests at once; after
// that it is let through at its sustained rate, and a request beyond either
// is denied. One caller's requests never use up another's.
var STK_261126_H2SW = &charter.Spec{ID: "STK-261126-H2SW", Title: "A caller gets a burst at once, then its sustained rate, and no more", Status: charter.Active}

// STK_261126_T2BD is STK-261126-T2BD: An operator sets the rate and the burst at start, and a bad setting stops the start
//
// The operator sets the sustained rate and the burst when starting
// ratelimit, without rebuilding it; a setting it cannot read stops it at the
// start with the setting named, rather than limiting by a value nobody
// chose.
var STK_261126_T2BD = &charter.Spec{ID: "STK-261126-T2BD", Title: "An operator sets the rate and the burst at start, and a bad setting stops the start", Status: charter.Active}

// SW_261126_2XG8 is SW-261126-2XG8: The limiter reads time through its Clock alone
//
// The limiter takes a Clock when it is made and reads the time from it and
// nowhere else; the command passes the wall clock.
var SW_261126_2XG8 = &charter.Spec{ID: "SW-261126-2XG8", Title: "The limiter reads time through its Clock alone", Status: charter.Active, Kind: "structure"}

// SW_261126_39F6 is SW-261126-39F6: One key's requests never change another key's bucket
//
// Each key's bucket is its own: spending, earning and denial on one key
// leave every other key's bucket as it was, under concurrent use too.
var SW_261126_39F6 = &charter.Spec{ID: "SW-261126-39F6", Title: "One key's requests never change another key's bucket", Status: charter.Active, Kind: "invariant"}

// SW_261126_CMTG is SW-261126-CMTG: The rate and the burst come from the environment at start, and an unreadable one stops the start
//
// RATELIMIT_RATE is the tokens a key earns per second, a decimal number,
// and RATELIMIT_BURST the most a key holds, a whole number; unset, they are
// one and five. A value that is not a number of its form, a rate not above
// zero or a burst below one stops the command before it reads a key, with
// exit status 2 and the setting named.
var SW_261126_CMTG = &charter.Spec{ID: "SW-261126-CMTG", Title: "The rate and the burst come from the environment at start, and an unreadable one stops the start", Status: charter.Active}

// SW_261126_VX8Z is SW-261126-VX8Z: A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request
//
// A key the limiter has not seen gets a bucket holding the burst. The
// bucket earns the rate in tokens per second of the clock's time, never
// holding more than the burst. A request is let through when the bucket
// holds a whole token and spends it; otherwise it is denied and spends
// nothing.
var SW_261126_VX8Z = &charter.Spec{ID: "SW-261126-VX8Z", Title: "A key's bucket starts full, earns the rate, holds at most the burst and pays a token a request", Status: charter.Active, Kind: "bound"}
