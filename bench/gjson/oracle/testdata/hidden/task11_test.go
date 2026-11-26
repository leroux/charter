// Code cut from upstream d7c940e..d7c940e gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestArrayCountRawOutput(t *testing.T) {
	assert(t, Get(`[1,2,3,4]`, "#").Raw == "4")
}
