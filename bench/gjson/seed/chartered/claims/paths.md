# Paths

Code: gjson.go

How Get reads a path: the keys and indices it is made of, the wildcards
and escapes a key may hold, the pipe that chains one path onto the
result of another, and the `..` prefix that reads JSON lines.

## SW-261126-T9SN A path is keys separated by dots, a key naming an object member and a number an array index

Status: active
Serves: STK-261126-ZQJR

Get reads a path from left to right, one dot-separated component at a
time: on an object a component names the member with that key, on an
array a component that is a whole number names the element at that
index, counting from zero. A component that names nothing, a key the
object lacks or an index past the end, ends the read with a result that
does not exist, and so does an empty path.

Why: dot paths are the whole interface; a caller predicts the answer
from the document alone, and a missing value is told apart from a
present one by Exists, not by a guess.

Scenarios:
1. `member` (agent, provisional): `name.last` on
   `{"name":{"first":"Tom","last":"Anderson"}}` returns `Anderson`.
2. `index` (agent, provisional): `children.1` on
   `{"children":["Sara","Alex","Jack"]}` returns `Alex`.
3. `missing` (agent, provisional): a key the object lacks, an index past
   the end of the array and a path through a string each return a result
   that does not exist.
4. `empty-path` (agent, provisional): the empty path returns a result
   that does not exist.

## SW-261126-4HQP A key with `*` or `?` matches the first member, in document order, whose key fits the pattern

Status: active
Serves: STK-261126-ZQJR

In a key, `*` matches any run of characters, empty too, and `?` any one
character. A key holding either is a pattern: Get follows the first
member of the object, in the order the document lists them, whose key
the pattern matches, and only that one; the result is that member's
value, never an array of every match. A pattern no key matches gives a
result that does not exist.

Why: a wildcard lets a caller reach a value whose key it knows only in
part; one value, the first, keeps a wildcard path returning the same
kind of result as any other path.

Scenarios:
1. `star` (agent, provisional): `child*.2` on
   `{"children":["Sara","Alex","Jack"]}` returns `Jack`.
2. `question` (agent, provisional): `c?ildren.0` on the same document
   returns `Sara`.
3. `first-match` (agent, provisional): `ab*` on `{"abc":1,"abd":2}`
   returns `1`, a number, not an array.
4. `no-match` (agent, provisional): `x*` on that document returns a
   result that does not exist.

## SW-261126-FX9X A backslash makes the next character of a path a literal part of the key

Status: active
Serves: STK-261126-ZQJR

A `\` in a path is dropped and the character after it is taken as part
of the key: an escaped dot does not end the key, and an escaped `*` or
`?` matches only itself.

Why: JSON keys may hold dots and wildcard characters, and a caller must
be able to name every key a document can have.

Scenarios:
1. `escaped-dot` (agent, provisional): `fav\.movie` on
   `{"fav.movie":"Deer Hunter"}` returns `Deer Hunter`.
2. `escaped-wildcard` (agent, provisional): `key\*v` matches the key
   `key*v` and not the key `keyXv`, and `keyv\?` the key `keyv?` alone.

## SW-261126-8RR3 A key matches exactly, case and all, and a key the document repeats gives its first value

Status: active
Kind: invariant
Serves: STK-261126-ZQJR

A key without a wildcard matches a member whose key is the same
characters, compared byte for byte after the document's escapes are
decoded: `name` does not match `Name`. When an object lists one key
twice, Get returns the first, and Map keeps that same one.

Why: JSON keys are case-sensitive and a document may hold both `id`
and `ID`; folding case would make one unreachable. Taking the first of
a repeated key costs nothing in a reader that stops at the first hit,
and Map agreeing with Get keeps the two views of a document the same.

Scenarios:
1. `case-sensitive` (agent, provisional): `name` on `{"Name":1}` returns
   a result that does not exist, and `Name` returns `1`.
2. `escaped-key` (agent, provisional): `a"b` matches a key the document
   writes as `"a\"b"`.
3. `duplicate-first` (agent, provisional): `name` on
   `{"name":"Alex","name":"Peter"}` returns `Alex`, and the result's
   Map holds `Alex` under `name`.

## SW-261126-SQKQ A pipe applies the path after it to the result of the path before it

Status: active
Serves: STK-261126-ZQJR

A `|` ends a path and starts another, read against the result of the
first as a document of its own: `friends|0|first` reads `friends`, then
`0` of that, then `first` of that. Setting DisableChaining makes `|` an
ordinary character of a key.

Why: a pipe lets a caller run a path on the output of a modifier or of
an array path, which a dot cannot express, since a dot after `#.first`
reads into each element.

Scenarios:
1. `pipe` (agent, provisional): `friends|0|first` returns the first
   friend's first name.
2. `pipe-after-modifier` (agent, provisional): `friends|@reverse|0|age`
   returns the last friend's age.
3. `chaining-disabled` (agent, provisional): with DisableChaining set,
   `a|b` returns the member whose key is `a|b`.

## SW-261126-ZQDQ A path starting with `..` reads the document as JSON lines, one value a line, as if they were an array

Status: active
Serves: STK-261126-ZQJR

When a path starts with `..`, Get treats the document as a sequence of
JSON values, one a line, and reads the rest of the path against them as
the elements of an array: `..#` is the number of lines, `..1` the
second, `..#.name` the `name` of each. ForEachLine calls a function
with each line's value in order until it returns false.

Why: logs and streams are JSON lines; reading them with the array syntax
a caller already knows saves splitting the text first.

Scenarios:
1. `lines-array` (agent, provisional): on seven lines, `..#` is `7`,
   `..0` the first line's value, `..2.name` a member of the third and
   `..#.name` the names of the lines that have one.
2. `for-each-line` (agent, provisional): ForEachLine visits the lines in
   order and stops when the function returns false.
