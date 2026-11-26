// Code cut from upstream c34bf81..c34bf81 gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestIssue141(t *testing.T) {
	json := `{"data": [{"q": 11, "w": 12}, {"q": 21, "w": 22}, {"q": 31, "w": 32} ], "sql": "some stuff here"}`
	assert(t, Get(json, "data.#").Int() == 3)
	assert(t, Get(json, "data.#.{q}|@ugly").Raw == `[{"q":11},{"q":21},{"q":31}]`)
	assert(t, Get(json, "data.#.q|@ugly").Raw == `[11,21,31]`)
}
