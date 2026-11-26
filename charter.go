// Package charter holds the types of a charter: specs, each recording one
// decision about what the system does or one need a person has of it, and
// the observations beneath them about the systems outside the code.
//
// A spec is written in an area document in the adopter's entries directory;
// its Go declaration is generated from it, so a test can verify it with
// chartertest.Verifies and the compiler holds the reference. Production code
// names the specs it realizes with //charter:realizes and //charter:concerns
// directives, which the gate resolves. Observations are written in Go in
// the entries package.
package charter

import "fmt"

// Status is a spec's implementation state.
type Status int

const (
	// Planned is ruled and not built. Its symbol is generated, so a test can
	// be written first; the gate asks nothing of it.
	Planned Status = iota + 1
	// Active is built and in service. The gate requires code to realize it
	// and a test to verify each of its scenarios.
	Active
	// Deprecated is retired. Its symbol is still generated, marked
	// deprecated, so existing anchors keep building while they are stripped.
	Deprecated
)

// String returns the lower-case name of the status.
func (s Status) String() string {
	switch s {
	case Planned:
		return "planned"
	case Active:
		return "active"
	case Deprecated:
		return "deprecated"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

// Spec is the generated declaration of a spec. ID is its id, such as SW-118;
// Title states the decision or the need; Kind is invariant, bound or
// structure when the decision is one of those, else empty. Everything else
// about the spec, its scenarios included, is in its document.
type Spec struct {
	ID, Title string
	Status    Status
	Kind      string
}

// ObservationStatus is the state of a fact about an external system. It is
// not a Status: a spec's status says whether the code does something, while
// a fact is observed, read, broken or unclassified.
type ObservationStatus int

const (
	// Observed: a probe cites it and last passed on the Observed date.
	Observed ObservationStatus = iota + 1
	// Documented: read in the subject's documentation; Source says where.
	Documented
	// Broken: a probe failed or the subject changed; the specs resting on
	// it are suspect until a person rules.
	Broken
	// ObservationUnsure: noted, not yet classified.
	ObservationUnsure
)

// String returns the lower-case name of the status.
func (s ObservationStatus) String() string {
	switch s {
	case Observed:
		return "observed"
	case Documented:
		return "documented"
	case Broken:
		return "broken"
	case ObservationUnsure:
		return "unsure"
	default:
		return fmt.Sprintf("ObservationStatus(%d)", int(s))
	}
}

// Observation states a fact about something outside the codebase that a
// spec rests on: a vendor API, a server, a runtime, a model. Subject is the
// adopter's name for that thing. It records what was seen or read, and when;
// it is not a promise the code keeps. Supports names, by id, the specs that
// rest on it; the gate resolves them.
type Observation struct {
	ID, Text string
	Subject  string
	Status   ObservationStatus
	// Source is where the fact was read: a URL or a path in the repository.
	Source string
	// Observed is the date the probe last passed, YYYY-MM-DD.
	Observed string
	Supports []string
}
