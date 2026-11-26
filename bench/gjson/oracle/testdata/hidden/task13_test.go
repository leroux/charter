// Code cut from upstream c5e72cd..c5e72cd gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestSingleModifier(t *testing.T) {
	var data = `{"@key": "value"}`
	assert(t, Get(data, "@key").String() == "value")
	assert(t, Get(data, "\\@key").String() == "value")
}
