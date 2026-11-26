// Code cut from upstream 00b15d7..ebf0b3f gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestQueryArrayValues(t *testing.T) {
	json := `{
		"artists": [
			["Bob Dylan"],
			"John Lennon",
			"Mick Jagger",
			"Elton John",
			"Michael Jackson",
			"John Smith",
			true,
			123,
			456,
			false,
			null
		]
	}`
	assert(t, Get(json, `a*.#[0="Bob Dylan"]#|#`).String() == "1")
	assert(t, Get(json, `a*.#[0="Bob Dylan 2"]#|#`).String() == "0")
	assert(t, Get(json, `a*.#[%"John*"]#|#`).String() == "2")
	assert(t, Get(json, `a*.#[_%"John*"]#|#`).String() == "0")
	assert(t, Get(json, `a*.#[="123"]#|#`).String() == "1")
}
