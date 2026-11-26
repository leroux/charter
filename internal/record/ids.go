package record

import (
	"fmt"
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

// IDParts are the parts of an id: its prefix, and a numbered id's number or
// a dated id's day and code.
type IDParts struct {
	Prefix string
	Number int
	Day    string
	Code   string
}

// ParseID reads an id, returning the form to write when it is not one.
//
//charter:realizes CH-012
func ParseID(id string) (IDParts, string) {
	if m := numberedID.FindStringSubmatch(id); m != nil {
		n, err := strconv.Atoi(m[2])
		if err != nil {
			return IDParts{}, fmt.Sprintf("the number of %s is too long", id)
		}
		return IDParts{Prefix: m[1], Number: n}, ""
	}
	m := datedID.FindStringSubmatch(id)
	if m == nil {
		return IDParts{}, fmt.Sprintf("%q is not an id: an id is a prefix and a number, SW-118, or a prefix, the day it was taken and a code, SW-261001-K7Q2; charter next gives a new one", id)
	}
	if _, err := time.Parse(dayLayout, m[2]); err != nil || len(m[2]) != len(dayLayout) {
		return IDParts{}, fmt.Sprintf("the date of %s, %s, is not a day on the calendar written YYMMDD; charter next gives a new id", id, m[2])
	}
	if !codeForm.MatchString(m[3]) {
		return IDParts{}, fmt.Sprintf("the code of %s, %s, is not %d characters of %s; charter next gives a new id", id, m[3], codeLength, codeAlphabet)
	}
	return IDParts{Prefix: m[1], Day: m[2], Code: m[3]}, ""
}

// IsID reports whether s is an id of either kind.
func IsID(s string) bool {
	_, problem := ParseID(s)
	return problem == ""
}

// Before orders two entries of one prefix: the numbered ids first, by
// number, then the dated ones by day and then by code.
//
//charter:realizes CH-012
func Before(a, b *Entry) bool {
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
func takeNewID(prefix string) string {
	return fmt.Sprintf("the entry that has not landed takes a new id: charter next %s", prefix)
}

// NewIDs draws count ids with the prefix for the day of now, each code
// character by intn over the alphabet, drawing again for an id the record
// holds, a question holds, or an earlier draw gave.
//
//charter:realizes CH-040
func (r *Record) NewIDs(prefix string, count int, now time.Time, intn func(int) int) ([]string, error) {
	if !prefixForm.MatchString(prefix) {
		return nil, fmt.Errorf("a prefix is capital letters, such as SW, not %q", prefix)
	}
	if count < 1 {
		return nil, fmt.Errorf("the count of new ids is a number above zero, not %d", count)
	}
	held := map[string]bool{}
	for _, q := range r.Questions {
		held[q.ID] = true
	}
	day := now.UTC().Format(dayLayout)
	ids := make([]string, 0, count)
	for len(ids) < count {
		var code strings.Builder
		for range codeLength {
			code.WriteByte(codeAlphabet[intn(len(codeAlphabet))])
		}
		id := prefix + "-" + day + "-" + code.String()
		if r.byID[id] != nil || held[id] {
			continue
		}
		held[id] = true
		ids = append(ids, id)
	}
	return ids, nil
}
