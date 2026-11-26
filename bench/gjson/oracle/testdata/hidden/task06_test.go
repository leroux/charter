// Code cut from upstream d95cbca gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestMultiArrayEx(t *testing.T) {
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

	var res string

	res = Get(json, `info.friends.#[kind="Person"]#.kind|0`).String()
	if res != "Person" {
		t.Fatalf("expected '%v', got '%v'", "Person", res)
	}
	res = Get(json, `info.friends.#.kind|0`).String()
	if res != "Person" {
		t.Fatalf("expected '%v', got '%v'", "Person", res)
	}

	res = Get(json, `info.friends.#[kind="Person"]#.kind`).String()
	if res != `["Person","Person"]` {
		t.Fatalf("expected '%v', got '%v'", `["Person","Person"]`, res)
	}
	res = Get(json, `info.friends.#.kind`).String()
	if res != `["Person","Person"]` {
		t.Fatalf("expected '%v', got '%v'", `["Person","Person"]`, res)
	}

	res = Get(json, `info.friends.#[kind="Person"]#|kind`).String()
	if res != `` {
		t.Fatalf("expected '%v', got '%v'", ``, res)
	}
	res = Get(json, `info.friends.#|kind`).String()
	if res != `` {
		t.Fatalf("expected '%v', got '%v'", ``, res)
	}

	res = Get(json, `i*.f*.#[kind="Other"]#`).String()
	if res != `[]` {
		t.Fatalf("expected '%v', got '%v'", `[]`, res)
	}
}

func TestQueries(t *testing.T) {
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

	assert(t, Get(json, "i*.f*.#[extra.0<11].first").Exists())
	assert(t, Get(json, "i*.f*.#[extra.0<=11].first").Exists())
	assert(t, !Get(json, "i*.f*.#[extra.0<10].first").Exists())
	assert(t, Get(json, "i*.f*.#[extra.0<=10].first").Exists())
	assert(t, Get(json, "i*.f*.#[extra.0=10].first").Exists())
	assert(t, !Get(json, "i*.f*.#[extra.0=11].first").Exists())
	assert(t, Get(json, "i*.f*.#[extra.0!=10].first").String() == "Roger")
	assert(t, Get(json, "i*.f*.#[extra.0>10].first").String() == "Roger")
	assert(t, Get(json, "i*.f*.#[extra.0>=10].first").String() == "Dale")

	assert(t, Get(json, `i*.f*.#[extra.0<"11"].first`).Exists())
	assert(t, Get(json, `i*.f*.#[first>"Dale"].last`).String() == "Craig")
	assert(t, Get(json, `i*.f*.#[first>="Dale"].last`).String() == "Murphy")
	assert(t, Get(json, `i*.f*.#[first="Dale"].last`).String() == "Murphy")
	assert(t, Get(json, `i*.f*.#[first!="Dale"].last`).String() == "Craig")
	assert(t, !Get(json, `i*.f*.#[first<"Dale"].last`).Exists())
	assert(t, Get(json, `i*.f*.#[first<="Dale"].last`).Exists())
	assert(t, Get(json, `i*.f*.#[first%"Da*"].last`).Exists())
	assert(t, Get(json, `i*.f*.#[first%"Dale"].last`).Exists())
	assert(t, Get(json, `i*.f*.#[first%"*a*"]#|#`).String() == "1")
	assert(t, Get(json, `i*.f*.#[first%"*e*"]#|#`).String() == "2")
	assert(t, Get(json, `i*.f*.#[first!%"*e*"]#|#`).String() == "0")

	assert(t, Get(json, `i*.f*.#[cust1=true].first`).String() == "Dale")
	assert(t, Get(json, `i*.f*.#[cust2=false].first`).String() == "Roger")
	assert(t, Get(json, `i*.f*.#[cust1!=false].first`).String() == "Dale")
	assert(t, Get(json, `i*.f*.#[cust2!=true].first`).String() == "Roger")
	assert(t, !Get(json, `i*.f*.#[cust1>true].first`).Exists())
	assert(t, Get(json, `i*.f*.#[cust1>=true].first`).Exists())
	assert(t, !Get(json, `i*.f*.#[cust2<false].first`).Exists())
	assert(t, Get(json, `i*.f*.#[cust2<=false].first`).Exists())

}
