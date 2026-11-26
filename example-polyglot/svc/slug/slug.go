// Package slug turns a title into a slug.
package slug

import (
	"strings"
	"unicode"
)

// Of is the slug of title: its letters and digits in lower case, each run
// of them joined by one hyphen.
//
//charter:realizes SW-261126-4NQQ
func Of(title string) string {
	words := strings.FieldsFunc(title, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	return strings.ToLower(strings.Join(words, "-"))
}
