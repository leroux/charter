// Code cut from upstream 5c2e4b3..5c2e4b3 gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"strings"
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestChainedModifierStringArgs(t *testing.T) {

	AddModifier("push", func(json, arg string) string {
		json = strings.TrimSpace(json)
		if len(json) < 2 || !Parse(json).IsArray() {
			return json
		}
		json = strings.TrimSpace(json[1 : len(json)-1])
		if len(json) == 0 {
			return "[" + arg + "]"
		}
		return "[" + json + "," + arg + "]"
	})
	res := Get("[]", `@push:"2"|@push:"3"|@push:{"a":"b","c":["e","f"]}|@push:true|@push:10.23`)
	assert(t, res.String() == `["2","3",{"a":"b","c":["e","f"]},true,10.23]`)
}
