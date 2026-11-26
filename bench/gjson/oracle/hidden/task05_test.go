// Code cut from upstream e91a88b..e91a88b gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

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

}
