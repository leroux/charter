# Arrays

Code: gjson.go

What the `#` component does in a path: the length of an array, a path
read from every element, and the queries that pick the elements that
match.

## SW-261126-6FGV `#` as the last component of a path returns the length of the array before it

Status: active
Serves: STK-261126-ZQJR

On an array, a final `#` returns the number of elements as a number;
an empty array gives `0`. On anything that is not an array, an object
among them, it gives a result that does not exist.

Why: the count is the question a caller asks of an array most often,
and answering it without materialising the elements keeps it as cheap
as any other lookup.

Scenarios:
1. `count` (agent, provisional): `children.#` on
   `{"children":["Sara","Alex","Jack"]}` returns the number `3`.
2. `count-empty` (agent, provisional): `#` on `[]` returns `0`.
3. `count-not-array` (agent, provisional): `a.#` on `{"a":{"x":1}}`
   returns a result that does not exist.

## SW-261126-V9J9 `#.path` reads the path from every element and returns the values found, in order, as an array

Status: active
Serves: STK-261126-ZQJR

A `#` followed by a dot and a path reads that path against each element
of the array and returns a JSON array of the values that exist, in the
elements' order; an element where the path finds nothing adds nothing,
and when no element has it the result is the empty array `[]`.

Why: pulling one field out of every record is the array lookup callers
need most, and skipping the records that lack it keeps the answer a
list of real values.

Scenarios:
1. `each` (agent, provisional): `friends.#.first` returns
   `["Dale","Roger","Jane"]`.
2. `each-skips-missing` (agent, provisional): `a.#.b` on
   `{"a":[{"b":1},{"c":2},{"b":3}]}` returns `[1,3]`.
3. `each-none` (agent, provisional): `a.#.z` on that document returns
   `[]`, a result that exists.

## SW-261126-2TYJ `#[query]` returns the first element that matches the query, and `#[query]#` every one as an array

Status: active
Serves: STK-261126-ZQJR

A `#[` path op value `]` component tests each element of the array:
the path is read against the element and its value compared with the
value by the operator. `#[...]` returns the first element that matches
and the rest of the path reads on from it; `#[...]#` returns a JSON
array of every match, with the rest of the path read from each, as
`#.` does. When no element matches, either form gives a result that
does not exist.

Why: finding a record by a field is how a caller reaches into a list of
objects without knowing positions; the two forms are the one-record and
every-record questions.

Scenarios:
1. `first` (agent, provisional): `friends.#[last=="Murphy"].first`
   returns `Dale`.
2. `all` (agent, provisional): `friends.#[last=="Murphy"]#.first`
   returns `["Dale","Jane"]`.
3. `no-match` (agent, provisional): `friends.#[age>100]` and
   `friends.#[age>100]#` each return a result that does not exist.

## SW-261126-XMB5 A query compares by the element's type, with `==`, `=`, `!=`, `<`, `<=`, `>`, `>=`, and matches strings with `%` and `!%`

Status: active
Serves: STK-261126-ZQJR

The value in a query is a quoted string, a number, `true` or `false`.
`=` and `==` are the same operator. A string element compares with the
value's text in byte order, and `%` matches it against a pattern where
`*` is any run and `?` any one character, `!%` its negation. A number
element compares numerically. A boolean element matches `=` and `!=`
against `true` or `false`, with true above false. An element of
another type matches nothing.

Why: one operator set over every JSON type keeps a query readable as a
condition, and comparing by the element's own type means `age>45`
means a number comparison where the data holds numbers.

Scenarios:
1. `strings` (agent, provisional): on the friends, `first>"Dale"`
   matches Roger and Jane, and `first=="Dale"` only Dale.
2. `numbers` (agent, provisional): `age>45` matches Roger and Jane,
   `age<=44` only Dale, `age!=44` the other two.
3. `like` (agent, provisional): `first%"D*"` matches Dale, and
   `first!%"D*"` Roger and Jane.
4. `booleans` (agent, provisional): `ok=true` matches the elements whose
   `ok` is true, and `ok!=true` those whose `ok` is false.
