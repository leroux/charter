// Code cut from upstream 8e88233..8e88233 gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"strings"
	"testing"

	"github.com/tidwall/pretty"
)

func TestModifier(t *testing.T) {
	json := `{"other":{"hello":"world"},"arr":[1,2,3,4,5,6]}`
	opts := *pretty.DefaultOptions
	opts.SortKeys = true
	exp := string(pretty.PrettyOptions([]byte(json), &opts))
	res := Get(json, `@pretty:{"sortKeys":true}`).String()
	if res != exp {
		t.Fatalf("expected '%v', got '%v'", exp, res)
	}
	res = Get(res, "@pretty|@reverse|@ugly").String()
	if res != json {
		t.Fatalf("expected '%v', got '%v'", json, res)
	}
	if res := Get(res, "@this").String(); res != json {
		t.Fatalf("expected '%v', got '%v'", json, res)
	}
	if res := Get(res, "other.@this").String(); res != `{"hello":"world"}` {
		t.Fatalf("expected '%v', got '%v'", json, res)
	}
	res = Get(res, "@pretty|@reverse|arr|@reverse|2").String()
	if res != "4" {
		t.Fatalf("expected '%v', got '%v'", "4", res)
	}
	AddModifier("case", func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}
		if arg == "lower" {
			return strings.ToLower(json)
		}
		return json
	})
	res = Get(json, "other|@case:upper").String()
	if res != `{"HELLO":"WORLD"}` {
		t.Fatalf("expected '%v', got '%v'", `{"HELLO":"WORLD"}`, res)
	}
}
