# Validity

Code: gjson.go

What gjson does with a document it cannot trust: Valid to ask, a
reader that does not validate and never panics, and the byte and
many-path forms of Get that answer as Get does.

## SW-261126-KNNZ Valid reports whether a document is exactly one JSON value, as the JSON grammar has it

Status: active
Serves: STK-261126-94D2

Valid and ValidBytes return true when the document, whitespace aside, is
one JSON value that the JSON grammar accepts, and false otherwise: a
trailing comma, a number with a leading zero, an unterminated string or
text after the value make it invalid.

Why: Get does not validate; Valid is the one place a caller learns
whether a document is well formed, and it must be strict for the answer
to be worth asking.

Scenarios:
1. `valid` (agent, provisional): `0`, `{"a":[1,2]}` and `"x"` are
   valid.
2. `invalid` (agent, provisional): `00`, `{"a":1,}`, `"x` and `1 2`
   are not.

## SW-261126-9CTW Get does not validate: on a malformed document it returns what it finds, and it never panics

Status: active
Kind: invariant
Serves: STK-261126-8FW8, STK-261126-94D2

Get reads only as far as the path needs and does not check the rest of
the document, so a malformed document may still give values; whatever
the document and the path, Get returns a result and does not panic.

Why: validating would make every lookup read the whole document, the
cost gjson exists to avoid; the caller who needs validation asks Valid,
and the one promise on bad input is that the program survives it.

Scenarios:
1. `lenient` (agent, provisional): `array.1` on
   `{"array":["a","b",0],}` returns `b`.
2. `no-panic` (agent, provisional): random bytes and the prefixes of a
   document, read with paths of every form, return without panicking.

## SW-261126-3FK5 GetBytes and GetMany return what Get returns for the same document and paths

Status: active
Serves: STK-261126-ZQJR

GetBytes on a byte slice returns what Get returns on the same text,
with strings that do not alias the slice; GetMany and GetManyBytes
return one result for each path, in order, each what Get returns for
it.

Why: the byte and many-path forms exist for speed, and a caller swaps
one for another only if the answers stay the same.

Scenarios:
1. `bytes-same` (agent, provisional): GetBytes and Get agree on every
   path of a document, and the string GetBytes returns survives the
   slice being overwritten.
2. `many-same` (agent, provisional): GetMany of three paths returns
   what three Gets return, in order.
