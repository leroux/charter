// Code cut from upstream 1c258af..1c258af gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

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

func TestModifierInSubSelector(t *testing.T) {
	AddModifier("case", func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}
		if arg == "lower" {
			return strings.ToLower(json)
		}
		return json
	})
	json := `{"friends": [
		{"age": 44, "first": "Dale", "last": "Murphy"},
		{"age": 68, "first": "Roger", "last": "Craig"},
		{"age": 47, "first": "Jane", "last": "Murphy"}
	]}`

	res := Get(json, `friends.#.{age,first|@case:upper}|@ugly`)
	exp := `[{"age":44,"@case:upper":"DALE"},{"age":68,"@case:upper":"ROGER"},{"age":47,"@case:upper":"JANE"}]`
	assert(t, res.Raw == exp)

	res = Get(json, `{friends.#.{age,first:first|@case:upper}|0.first}`)
	exp = `{"first":"DALE"}`
	assert(t, res.Raw == exp)

}
