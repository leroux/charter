package chartertest

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// An id is a prefix and then a number, SW-118, or the day it was taken and
// a code, SW-261001-K7Q2. The numbered ids are the ones given before ids
// were dated; every new id is dated.

// codeAlphabet is what a dated id's code is drawn from: the consonants and
// digits that spell no word and are not read as one another.
const codeAlphabet = "BCDFGHJKMNPQRSTVWXYZ23456789"

// codeLength is the length of a dated id's code.
const codeLength = 4

// dayLayout is a dated id's day: the UTC date as year, month and day.
const dayLayout = "060102"

var (
	prefixForm = regexp.MustCompile(`^[A-Z]+$`)
	numberedID = regexp.MustCompile(`^([A-Z]+)-(\d+)$`)
	datedID    = regexp.MustCompile(`^([A-Z]+)-(\d+)-([^-]+)$`)
	codeForm   = regexp.MustCompile(fmt.Sprintf(`^[%s]{%d}$`, codeAlphabet, codeLength))
)

// idParts are the parts of an id: its prefix, and a numbered id's number or
// a dated id's day and code.
type idParts struct {
	Lens   string
	Number int
	Day    string
	Code   string
}

// parseID reads an id, returning the form to write when it is not one.
//
//charter:realizes CH-012
func parseID(id string) (idParts, string) {
	if m := numberedID.FindStringSubmatch(id); m != nil {
		n, err := strconv.Atoi(m[2])
		if err != nil {
			return idParts{}, fmt.Sprintf("the number of %s is too long", id)
		}
		return idParts{Lens: m[1], Number: n}, ""
	}
	m := datedID.FindStringSubmatch(id)
	if m == nil {
		return idParts{}, fmt.Sprintf("%q is not an id: an id is a prefix and a number, SW-118, or a prefix, the day it was taken and a code, SW-261001-K7Q2; claimsgen next gives a new one", id)
	}
	if _, err := time.Parse(dayLayout, m[2]); err != nil || len(m[2]) != len(dayLayout) {
		return idParts{}, fmt.Sprintf("the date of %s, %s, is not a day on the calendar written YYMMDD; claimsgen next gives a new id", id, m[2])
	}
	if !codeForm.MatchString(m[3]) {
		return idParts{}, fmt.Sprintf("the code of %s, %s, is not %d characters of %s; claimsgen next gives a new id", id, m[3], codeLength, codeAlphabet)
	}
	return idParts{Lens: m[1], Day: m[2], Code: m[3]}, ""
}

// isID reports whether s is an id of either kind.
func isID(s string) bool {
	_, problem := parseID(s)
	return problem == ""
}

// idBefore orders two entries of one prefix: the numbered ids first, by
// number, then the dated ones by day and then by code.
//
//charter:realizes CH-012
func idBefore(a, b *entry) bool {
	if (a.Day == "") != (b.Day == "") {
		return a.Day == ""
	}
	if a.Day == "" {
		return a.Number < b.Number
	}
	if a.Day != b.Day {
		return a.Day < b.Day
	}
	return a.Code < b.Code
}

// takeNewID is what a declared-twice problem ends with: which of the two
// entries gives way, and where its new id comes from.
func takeNewID(e *entry) string {
	return fmt.Sprintf("the entry that has not landed takes a new id: claimsgen next %s", e.Lens)
}

// NextIDs returns count new ids with the prefix: each the UTC date and a
// code drawn at random, different from one another and from every id the
// charter's documents hold. It reads the documents alone, so it works on a
// tree whose generated files are stale and whose tests do not parse.
//
//charter:realizes CH-040
func NextIDs(opts Options, prefix string, count int) ([]string, error) {
	return nextIDs(opts, callerDir(1), prefix, count, time.Now(), rand.IntN)
}

// nextIDs is NextIDs for the instant now and the random source intn, with
// the root resolved from the directory from.
func nextIDs(opts Options, from, prefix string, count int, now time.Time, intn func(int) int) ([]string, error) {
	if !prefixForm.MatchString(prefix) {
		return nil, fmt.Errorf("a prefix is capital letters, such as SW, not %q", prefix)
	}
	if count < 1 {
		return nil, fmt.Errorf("the count of new ids is a number above zero, not %d", count)
	}
	r, err := opts.resolve(from)
	if err != nil {
		return nil, err
	}
	pkg := newEntriesPackage()
	if _, _, err := readAreas(r.Root, r.areasDir, pkg, r.notAreaDocuments()...); err != nil {
		return nil, err
	}
	if _, err := readSpecs(r.Root, r.specsDir, pkg); err != nil {
		return nil, err
	}
	return newIDs(pkg, prefix, count, now, intn), nil
}

// newIDs draws count ids with the prefix for the day of now, each code
// character by intn over the alphabet, drawing again for an id the entries
// hold or an earlier draw gave.
func newIDs(pkg *entriesPackage, prefix string, count int, now time.Time, intn func(int) int) []string {
	day := now.UTC().Format(dayLayout)
	drawn := map[string]bool{}
	ids := make([]string, 0, count)
	for len(ids) < count {
		var code strings.Builder
		for range codeLength {
			code.WriteByte(codeAlphabet[intn(len(codeAlphabet))])
		}
		id := prefix + "-" + day + "-" + code.String()
		if pkg.byID[id] != nil || drawn[id] {
			continue
		}
		drawn[id] = true
		ids = append(ids, id)
	}
	return ids
}
