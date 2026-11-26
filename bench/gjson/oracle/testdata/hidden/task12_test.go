// Code cut from upstream 1e964df..1e964df gjson_test.go by the gjson benchmark's extract. DO NOT EDIT.

package gjson

import (
	"testing"
)

func assert(t testing.TB, cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func TestParentSubQuery(t *testing.T) {
	var json = `{
		"topology": {
		  "instances": [
			{
			  "service_version": "1.2.3",
			  "service_locale": {"lang": "en"},
			  "service_roles": ["one", "two"]
			},
			{
			  "service_version": "1.2.4",
			  "service_locale": {"lang": "th"},
			  "service_roles": ["three", "four"]
			},
			{
			  "service_version": "1.2.2",
			  "service_locale": {"lang": "en"},
			  "service_roles": ["one"]
			}
		  ]
		}
	  }`
	res := Get(json, `topology.instances.#( service_roles.#(=="one"))#.service_version`)

	assert(t, res.String() == `["1.2.3","1.2.2"]`)
}
