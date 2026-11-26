// Code cut from upstream 3b5bf6b..3b5bf6b gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestSubSelectors(t *testing.T) {
	json := `{
		"info": {
			"friends": [
				{
					"first": "Dale", "last": "Murphy", "kind": "Person",
					"cust1": true,
					"extra": [10,20,30],
					"details": {
						"city": "Tempe",
						"state": "Arizona"
					}
				},
				{
					"first": "Roger", "last": "Craig", "kind": "Person",
					"cust2": false,
					"extra": [40,50,60],
					"details": {
						"city": "Phoenix",
						"state": "Arizona"
					}
				}
			]
		}
	  }`
	assert(t, Get(json, "[]").String() == "[]")
	assert(t, Get(json, "{}").String() == "{}")
	res := Get(json, `{`+
		`abc:info.friends.0.first,`+
		`info.friends.1.last,`+
		`"a`+"\r"+`a":info.friends.0.kind,`+
		`"abc":info.friends.1.kind,`+
		`{123:info.friends.1.cust2},`+
		`[info.friends.#[details.city="Phoenix"]#|#]`+
		`}.@pretty.@ugly`).String()

	assert(t, Get(res, "abc").String() == "Dale")
	assert(t, Get(res, "last").String() == "Craig")
	assert(t, Get(res, "\"a\ra\"").String() == "Person")
	assert(t, Get(res, "@reverse.abc").String() == "Person")
	assert(t, Get(res, "_.123").String() == "false")
	assert(t, Get(res, "@reverse._.0").String() == "1")
	assert(t, Get(json, "info.friends.[0.first,1.extra.0]").String() ==
		`["Dale",40]`)
	assert(t, Get(json, "info.friends.#.[first,extra.0]").String() ==
		`[["Dale",10],["Roger",40]]`)
}
