kind: neutral
Issue #149: please add a built-in `@this` modifier that returns the current element as-is, so a path can name the whole document (the equivalent of JSON Pointer `#/`). `Get(json, "@this")` returns the document, and `other.@this` on `{"other":{"hello":"world"}}` returns `{"hello":"world"}`. Please add it to the list of modifiers in the docs and add tests.
