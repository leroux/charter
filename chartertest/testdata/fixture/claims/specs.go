// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: CHARTER_UPDATE_SPECS=1 go test -count=1 ./claims

package claims

import "github.com/leroux/charter"

// STK_001 is STK-001: A person can come back and find everything as it was
//
// Everything a person left is there when they come back.
var STK_001 = &charter.Spec{ID: "STK-001", Title: "A person can come back and find everything as it was", Status: charter.Active}

// STK_002 is STK-002: An operator can run the limiter on another machine
//
// The limiter runs on any machine.
//
// Planned: STK-002 is ruled and not built; see claims/specs/STK-002-run-anywhere.md.
var STK_002 = &charter.Spec{ID: "STK-002", Title: "An operator can run the limiter on another machine", Status: charter.Planned}

// SW_001 is SW-001: A bucket survives a restart
//
// A key's bucket is there after the limiter restarts.
var SW_001 = &charter.Spec{ID: "SW-001", Title: "A bucket survives a restart", Status: charter.Active}

// SW_002 is SW-002: The limiter runs on Windows
//
// The limiter builds and runs on Windows.
//
// Deprecated: SW-002 is retired; see its Changes in claims/specs/SW-002-limiter-runs-on-windows.md. Add no new anchors.
var SW_002 = &charter.Spec{ID: "SW-002", Title: "The limiter runs on Windows", Status: charter.Deprecated}

// SW_003 is SW-003: A person can clear a quota's counters
//
// A clear empties the quota's counters.
//
// Planned: SW-003 is ruled and not built; see claims/specs/SW-003-clear-variables.md.
var SW_003 = &charter.Spec{ID: "SW-003", Title: "A person can clear a quota's counters", Status: charter.Planned}

// CON_001 is CON-001: The lock file is removed on close
//
// Closing the store removes its lock file.
var CON_001 = &charter.Spec{ID: "CON-001", Title: "The lock file is removed on close", Status: charter.Active}

// CON_002 is CON-002: The fold never panics
//
// Folding any log returns without panicking.
var CON_002 = &charter.Spec{ID: "CON-002", Title: "The fold never panics", Status: charter.Active}
