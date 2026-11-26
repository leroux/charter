// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: go run ./cmd/claimsgen

package claims

import "github.com/leroux/charter"

// STK_261126_8FW8 is STK-261126-8FW8: Get never panics, whatever the document or the path
//
// Any string may be passed as the document and any string as the path:
// gjson returns a result, possibly one that does not exist, and never
// panics.
var STK_261126_8FW8 = &charter.Spec{ID: "STK-261126-8FW8", Title: "Get never panics, whatever the document or the path", Status: charter.Active}

// STK_261126_94D2 is STK-261126-94D2: A caller can tell whether a document is valid JSON before trusting it
//
// gjson reads without validating, so a caller who cannot trust a
// document's source asks gjson first whether it is valid JSON.
var STK_261126_94D2 = &charter.Spec{ID: "STK-261126-94D2", Title: "A caller can tell whether a document is valid JSON before trusting it", Status: charter.Active}

// STK_261126_TFYB is STK-261126-TFYB: A caller reads a result as the Go value it needs, and a big integer keeps every digit
//
// A result converts to a string, a number, a boolean, a slice or a map
// without the caller parsing JSON, and an integer read as an integer has
// every digit the document wrote, beyond what a float64 holds.
var STK_261126_TFYB = &charter.Spec{ID: "STK-261126-TFYB", Title: "A caller reads a result as the Go value it needs, and a big integer keeps every digit", Status: charter.Active}

// STK_261126_ZQJR is STK-261126-ZQJR: A path returns the value the documented syntax names, and nothing when it names none
//
// A caller writes a path in the syntax the README documents and gets the
// value it names in the document; when the path names no value, the
// caller gets a result that does not exist, never some other value.
var STK_261126_ZQJR = &charter.Spec{ID: "STK-261126-ZQJR", Title: "A path returns the value the documented syntax names, and nothing when it names none", Status: charter.Active}

// SW_261126_26XP is SW-261126-26XP: Exists is false only for a value the path did not find, and a JSON null exists
//
// Exists reports whether the result came from a value in the document: a
// member whose value is `null` exists, with Type Null, and so does an
// empty string, while a path that found nothing gives a result whose
// Exists is false and whose Type is also Null.
var SW_261126_26XP = &charter.Spec{ID: "SW-261126-26XP", Title: "Exists is false only for a value the path did not find, and a JSON null exists", Status: charter.Active}

// SW_261126_2TYJ is SW-261126-2TYJ: `#[query]` returns the first element that matches the query, and `#[query]#` every one as an array
//
// A `#[` path op value `]` component tests each element of the array:
// the path is read against the element and its value compared with the
// value by the operator. `#[...]` returns the first element that matches
// and the rest of the path reads on from it; `#[...]#` returns a JSON
// array of every match, with the rest of the path read from each, as
// `#.` does. When no element matches, either form gives a result that
// does not exist.
var SW_261126_2TYJ = &charter.Spec{ID: "SW-261126-2TYJ", Title: "`#[query]` returns the first element that matches the query, and `#[query]#` every one as an array", Status: charter.Active}

// SW_261126_3FK5 is SW-261126-3FK5: GetBytes and GetMany return what Get returns for the same document and paths
//
// GetBytes on a byte slice returns what Get returns on the same text,
// with strings that do not alias the slice; GetMany and GetManyBytes
// return one result for each path, in order, each what Get returns for
// it.
var SW_261126_3FK5 = &charter.Spec{ID: "SW-261126-3FK5", Title: "GetBytes and GetMany return what Get returns for the same document and paths", Status: charter.Active}

// SW_261126_4HQP is SW-261126-4HQP: A key with `*` or `?` matches the first member, in document order, whose key fits the pattern
//
// In a key, `*` matches any run of characters, empty too, and `?` any one
// character. A key holding either is a pattern: Get follows the first
// member of the object, in the order the document lists them, whose key
// the pattern matches, and only that one; the result is that member's
// value, never an array of every match. A pattern no key matches gives a
// result that does not exist.
var SW_261126_4HQP = &charter.Spec{ID: "SW-261126-4HQP", Title: "A key with `*` or `?` matches the first member, in document order, whose key fits the pattern", Status: charter.Active}

// SW_261126_4WX9 is SW-261126-4WX9: Less orders results by type, Null, False, Number, String, True, JSON, and within a type by value
//
// Less reports whether a result sorts before another: a lower Type first,
// in the order Null, False, Number, String, True, JSON; two numbers by
// value; two strings by bytes, or ignoring ASCII case when caseSensitive
// is false; anything else by its raw text.
var SW_261126_4WX9 = &charter.Spec{ID: "SW-261126-4WX9", Title: "Less orders results by type, Null, False, Number, String, True, JSON, and within a type by value", Status: charter.Active}

// SW_261126_4Z6X is SW-261126-4Z6X: A caller adds a modifier by name, and a name no modifier has gives a result that does not exist
//
// AddModifier registers a function from a JSON text and an argument to a
// JSON text under a name, replacing any modifier of that name, and the
// name is then usable as `@name` in any path; ModifierExists reports
// whether a name is registered. A path naming a modifier nobody
// registered gives a result that does not exist.
var SW_261126_4Z6X = &charter.Spec{ID: "SW-261126-4Z6X", Title: "A caller adds a modifier by name, and a name no modifier has gives a result that does not exist", Status: charter.Active}

// SW_261126_5VFD is SW-261126-5VFD: Int and Uint read every digit of a 64-bit integer, and Float reads a number or a numeric string
//
// Int and Uint parse an integer from the raw text when it is written as
// one, so every value of int64 and uint64 comes back exact; otherwise
// they truncate the float64. Float returns a number's float64, the
// number a numeric string spells, `1` for true and `0` for anything else.
var SW_261126_5VFD = &charter.Spec{ID: "SW-261126-5VFD", Title: "Int and Uint read every digit of a 64-bit integer, and Float reads a number or a numeric string", Status: charter.Active}

// SW_261126_6FGV is SW-261126-6FGV: `#` as the last component of a path returns the length of the array before it
//
// On an array, a final `#` returns the number of elements as a number;
// an empty array gives `0`. On anything that is not an array, an object
// among them, it gives a result that does not exist.
var SW_261126_6FGV = &charter.Spec{ID: "SW-261126-6FGV", Title: "`#` as the last component of a path returns the length of the array before it", Status: charter.Active}

// SW_261126_8RR3 is SW-261126-8RR3: A key matches exactly, case and all, and a key the document repeats gives its first value
//
// A key without a wildcard matches a member whose key is the same
// characters, compared byte for byte after the document's escapes are
// decoded: `name` does not match `Name`. When an object lists one key
// twice, Get returns the first, and Map keeps that same one.
var SW_261126_8RR3 = &charter.Spec{ID: "SW-261126-8RR3", Title: "A key matches exactly, case and all, and a key the document repeats gives its first value", Status: charter.Active, Kind: "invariant"}

// SW_261126_9CTW is SW-261126-9CTW: Get does not validate: on a malformed document it returns what it finds, and it never panics
//
// Get reads only as far as the path needs and does not check the rest of
// the document, so a malformed document may still give values; whatever
// the document and the path, Get returns a result and does not panic.
var SW_261126_9CTW = &charter.Spec{ID: "SW-261126-9CTW", Title: "Get does not validate: on a malformed document it returns what it finds, and it never panics", Status: charter.Active, Kind: "invariant"}

// SW_261126_FX9X is SW-261126-FX9X: A backslash makes the next character of a path a literal part of the key
//
// A `\` in a path is dropped and the character after it is taken as part
// of the key: an escaped dot does not end the key, and an escaped `*` or
// `?` matches only itself.
var SW_261126_FX9X = &charter.Spec{ID: "SW-261126-FX9X", Title: "A backslash makes the next character of a path a literal part of the key", Status: charter.Active}

// SW_261126_H629 is SW-261126-H629: A path starting with `@name` runs that modifier and reads the rest of the path from its output
//
// A path, or a path after a pipe, that starts with `@` and a modifier's
// name runs the modifier on the JSON so far; an argument may follow the
// name after a `:`, as JSON or plain characters. The output is read as a
// document of its own by the rest of the path after a `|`. The built-in
// modifiers are `@reverse`, which reverses an array's elements or an
// object's members and leaves any other value as it was, `@ugly`, which
// removes all whitespace, and `@pretty`, which indents, taking the
// options `sortKeys`, `indent`, `prefix` and `width` as a JSON object.
var SW_261126_H629 = &charter.Spec{ID: "SW-261126-H629", Title: "A path starting with `@name` runs that modifier and reads the rest of the path from its output", Status: charter.Active}

// SW_261126_KNNZ is SW-261126-KNNZ: Valid reports whether a document is exactly one JSON value, as the JSON grammar has it
//
// Valid and ValidBytes return true when the document, whitespace aside, is
// one JSON value that the JSON grammar accepts, and false otherwise: a
// trailing comma, a number with a leading zero, an unterminated string or
// text after the value make it invalid.
var SW_261126_KNNZ = &charter.Spec{ID: "SW-261126-KNNZ", Title: "Valid reports whether a document is exactly one JSON value, as the JSON grammar has it", Status: charter.Active}

// SW_261126_NKTP is SW-261126-NKTP: Array, Map and ForEach walk an array's elements or an object's members in document order
//
// Array returns an array's elements in order; on a result that does not
// exist it returns an empty slice, and on any other value a slice of that
// one value. Map returns an object's members by key, and an empty map for
// anything that is not an object. ForEach calls a function with each key
// and value of an object, or each value of an array with a key that does
// not exist, in document order, until the function returns false; on a
// single value it calls it once with that value.
var SW_261126_NKTP = &charter.Spec{ID: "SW-261126-NKTP", Title: "Array, Map and ForEach walk an array's elements or an object's members in document order", Status: charter.Active}

// SW_261126_NRBX is SW-261126-NRBX: Value returns the result as the Go value encoding/json would decode it into an interface
//
// Value returns bool for true and false, float64 for a number, string
// for a string, nil for null and for nothing, []interface{} for an array
// and map[string]interface{} for an object, nested values converted the
// same way.
var SW_261126_NRBX = &charter.Spec{ID: "SW-261126-NRBX", Title: "Value returns the result as the Go value encoding/json would decode it into an interface", Status: charter.Active}

// SW_261126_PKGR is SW-261126-PKGR: String returns a string unescaped, an integer as the document wrote it, and other numbers as the shortest float
//
// A Result holds its Type (Null, False, Number, String, True or JSON),
// its Raw text as it appears in the document, the decoded Str of a
// string and the Num of a number. String returns a string's decoded text,
// escapes such as `\n` and `é` resolved, and nothing trimmed; a
// number written as an integer exactly as the document wrote it, every
// digit kept; any other number as the shortest decimal that reads back
// as its float64; `true` and `false` for the literals; the empty string
// for null and for a result that does not exist; and the raw JSON of an
// array or an object.
var SW_261126_PKGR = &charter.Spec{ID: "SW-261126-PKGR", Title: "String returns a string unescaped, an integer as the document wrote it, and other numbers as the shortest float", Status: charter.Active}

// SW_261126_SQKQ is SW-261126-SQKQ: A pipe applies the path after it to the result of the path before it
//
// A `|` ends a path and starts another, read against the result of the
// first as a document of its own: `friends|0|first` reads `friends`, then
// `0` of that, then `first` of that. Setting DisableChaining makes `|` an
// ordinary character of a key.
var SW_261126_SQKQ = &charter.Spec{ID: "SW-261126-SQKQ", Title: "A pipe applies the path after it to the result of the path before it", Status: charter.Active}

// SW_261126_T9SN is SW-261126-T9SN: A path is keys separated by dots, a key naming an object member and a number an array index
//
// Get reads a path from left to right, one dot-separated component at a
// time: on an object a component names the member with that key, on an
// array a component that is a whole number names the element at that
// index, counting from zero. A component that names nothing, a key the
// object lacks or an index past the end, ends the read with a result that
// does not exist, and so does an empty path.
var SW_261126_T9SN = &charter.Spec{ID: "SW-261126-T9SN", Title: "A path is keys separated by dots, a key naming an object member and a number an array index", Status: charter.Active}

// SW_261126_V9J9 is SW-261126-V9J9: `#.path` reads the path from every element and returns the values found, in order, as an array
//
// A `#` followed by a dot and a path reads that path against each element
// of the array and returns a JSON array of the values that exist, in the
// elements' order; an element where the path finds nothing adds nothing,
// and when no element has it the result is the empty array `[]`.
var SW_261126_V9J9 = &charter.Spec{ID: "SW-261126-V9J9", Title: "`#.path` reads the path from every element and returns the values found, in order, as an array", Status: charter.Active}

// SW_261126_XMB5 is SW-261126-XMB5: A query compares by the element's type, with `==`, `=`, `!=`, `<`, `<=`, `>`, `>=`, and matches strings with `%` and `!%`
//
// The value in a query is a quoted string, a number, `true` or `false`.
// `=` and `==` are the same operator. A string element compares with the
// value's text in byte order, and `%` matches it against a pattern where
// `*` is any run and `?` any one character, `!%` its negation. A number
// element compares numerically. A boolean element matches `=` and `!=`
// against `true` or `false`, with true above false. An element of
// another type matches nothing.
var SW_261126_XMB5 = &charter.Spec{ID: "SW-261126-XMB5", Title: "A query compares by the element's type, with `==`, `=`, `!=`, `<`, `<=`, `>`, `>=`, and matches strings with `%` and `!%`", Status: charter.Active}

// SW_261126_Z23S is SW-261126-Z23S: Parse reads a whole document as a result, and Result.Get reads a path inside one
//
// Parse returns the document's first value as a result, without a path;
// Result.Get on that result, or on any result holding an object or an
// array, reads a path against its raw JSON as Get would.
var SW_261126_Z23S = &charter.Spec{ID: "SW-261126-Z23S", Title: "Parse reads a whole document as a result, and Result.Get reads a path inside one", Status: charter.Active}

// SW_261126_Z3CC is SW-261126-Z3CC: DisableModifiers makes `@` an ordinary character of a key
//
// With DisableModifiers set, Get runs no modifier and a component
// starting with `@` is a key like any other.
var SW_261126_Z3CC = &charter.Spec{ID: "SW-261126-Z3CC", Title: "DisableModifiers makes `@` an ordinary character of a key", Status: charter.Active}

// SW_261126_ZGV3 is SW-261126-ZGV3: Bool is true for true, a number other than zero, and a string other than empty, `0` and `false`
//
// Bool returns true for the literal true, for a number not equal to zero
// and for a string other than the empty string, `"0"` and `"false"`, and
// false for anything else: false, null, an array, an object and a result
// that does not exist.
var SW_261126_ZGV3 = &charter.Spec{ID: "SW-261126-ZGV3", Title: "Bool is true for true, a number other than zero, and a string other than empty, `0` and `false`", Status: charter.Active}

// SW_261126_ZQDQ is SW-261126-ZQDQ: A path starting with `..` reads the document as JSON lines, one value a line, as if they were an array
//
// When a path starts with `..`, Get treats the document as a sequence of
// JSON values, one a line, and reads the rest of the path against them as
// the elements of an array: `..#` is the number of lines, `..1` the
// second, `..#.name` the `name` of each. ForEachLine calls a function
// with each line's value in order until it returns false.
var SW_261126_ZQDQ = &charter.Spec{ID: "SW-261126-ZQDQ", Title: "A path starting with `..` reads the document as JSON lines, one value a line, as if they were an array", Status: charter.Active}
