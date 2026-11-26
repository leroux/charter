// Code cut from upstream 6781e4e..6781e4e gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestParenQueries(t *testing.T) {
	json := `{
		"friends": [{"a":10},{"a":20},{"a":30},{"a":40}]
	}`
	assert(t, Get(json, "friends.#(a>9)#|#").Int() == 4)
	assert(t, Get(json, "friends.#(a>10)#|#").Int() == 3)
	assert(t, Get(json, "friends.#(a>40)#|#").Int() == 0)
}
