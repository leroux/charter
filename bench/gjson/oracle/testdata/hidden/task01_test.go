// Code cut from upstream 89b1979..89b1979 gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func TestArrayEx(t *testing.T) {
	s := `
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
	res := Get(s, "@ugly|#.c.#[a=10.11]").String()
	if res != `[{"a":10.11}]` {
		t.Fatalf("expected '%v', got '%v'", `[{"a":10.11}]`, res)
	}
	res = Get(s, "@ugly|#.c.#").String()
	if res != `[1,1]` {
		t.Fatalf("expected '%v', got '%v'", `[1,1]`, res)
	}
}
