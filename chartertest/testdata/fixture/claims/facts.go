package claims

import "github.com/leroux/charter"

// Facts about the fixture's world.
var LockIsAdvisory = &charter.Observation{
	ID:       "LockIsAdvisory",
	Text:     "flock locks are advisory: a process that does not ask is not stopped.",
	Subject:  "os",
	Status:   charter.Documented,
	Source:   "https://man7.org/linux/man-pages/man2/flock.2.html",
	Supports: []string{"CON-001"},
}

var HealthIsJSON = &charter.Observation{
	ID:       "HealthIsJSON",
	Text:     "The upstream health endpoint answers JSON.",
	Subject:  "upstream",
	Status:   charter.Observed,
	Observed: "2026-09-24",
	Supports: []string{"CON-002"},
}

var Suspect = &charter.Observation{
	ID:      "Suspect",
	Text:    "The upstream keeps message order.",
	Subject: "upstream",
	Status:  charter.ObservationUnsure,
}
