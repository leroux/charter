// Code cut from upstream b877bd4..b877bd4 gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func TestChaining(t *testing.T) {
	json := `{
		"info": {
			"friends": [
				{"first": "Dale", "last": "Murphy", "age": 44},
				{"first": "Roger", "last": "Craig", "age": 68},
				{"first": "Jane", "last": "Murphy", "age": 47}
			]
		}
	  }`
	res := Get(json, "info.friends|0|first").String()
	if res != "Dale" {
		t.Fatalf("expected '%v', got '%v'", "Dale", res)
	}
	res = Get(json, "info.friends|@reverse|0|age").String()
	if res != "47" {
		t.Fatalf("expected '%v', got '%v'", "47", res)
	}
	res = Get(json, "@ugly|i\\nfo|friends.0.first").String()
	if res != "Dale" {
		t.Fatalf("expected '%v', got '%v'", "Dale", res)
	}
}

func TestArrayEx(t *testing.T) {
	json := `
	[
		{
			"c":[
				{"a":10.11}
			]
		}, {
			"c":[
				{"a":11.11}
			]
		}
	]`
	res := Get(json, "@ugly|#.c.#[a=10.11]").String()
	if res != `[{"a":10.11}]` {
		t.Fatalf("expected '%v', got '%v'", `[{"a":10.11}]`, res)
	}
	res = Get(json, "@ugly|#.c.#").String()
	if res != `[1,1]` {
		t.Fatalf("expected '%v', got '%v'", `[1,1]`, res)
	}
	res = Get(json, "@reverse|0|c|0|a").String()
	if res != "11.11" {
		t.Fatalf("expected '%v', got '%v'", "11.11", res)
	}
	res = Get(json, "#.c|#").String()
	if res != "2" {
		t.Fatalf("expected '%v', got '%v'", "2", res)
	}
}
