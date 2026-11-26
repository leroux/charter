// Code cut from upstream 7660d0f..7660d0f gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func TestPipeDotMixing(t *testing.T) {
	json := `{
		"info": {
			"friends": [
				{"first": "Dale", "last": "Murphy", "age": 44},
				{"first": "Roger", "last": "Craig", "age": 68},
				{"first": "Jane", "last": "Murphy", "age": 47}
			]
		}
	  }`
	var res string
	res = Get(json, `info.friends.#[first="Dale"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `info|friends.#[first="Dale"].last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `info|friends.#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `info|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `@ugly|info|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `@ugly|info.@ugly|friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
	res = Get(json, `@ugly.info|@ugly.friends|#[first="Dale"]|last`).String()
	if res != "Murphy" {
		t.Fatalf("expected '%v', got '%v'", "Murphy", res)
	}
}
