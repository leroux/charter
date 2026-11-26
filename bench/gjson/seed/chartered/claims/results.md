# Results

Code: gjson.go

What a caller reads off a Result: its type and raw text, and the
conversions to Go values, the walks over arrays and objects, and the
order Less puts results in.

## SW-261126-PKGR String returns a string unescaped, an integer as the document wrote it, and other numbers as the shortest float

Status: active
Serves: STK-261126-TFYB

A Result holds its Type (Null, False, Number, String, True or JSON),
its Raw text as it appears in the document, the decoded Str of a
string and the Num of a number. String returns a string's decoded text,
escapes such as `\n` and `é` resolved, and nothing trimmed; a
number written as an integer exactly as the document wrote it, every
digit kept; any other number as the shortest decimal that reads back
as its float64; `true` and `false` for the literals; the empty string
for null and for a result that does not exist; and the raw JSON of an
array or an object.

Why: String is the conversion callers use most; returning an integer's
own digits keeps an ID beyond 2^53 intact, which a float cannot hold,
and a string's text is the caller's data, whitespace included.

Scenarios:
1. `string-decoded` (agent, provisional): a string holding `é`,
   `\n` and spaces at both ends returns them decoded and untrimmed.
2. `integer-as-written` (agent, provisional): the number
   9007199254740993 returns `9007199254740993`.
3. `float-shortest` (agent, provisional): the number `1.50` returns
   `1.5` while its Raw stays `1.50`.
4. `literals-and-json` (agent, provisional): true, false and null
   return `true`, `false` and the empty string, and an object returns
   its raw JSON.

## SW-261126-5VFD Int and Uint read every digit of a 64-bit integer, and Float reads a number or a numeric string

Status: active
Serves: STK-261126-TFYB

Int and Uint parse an integer from the raw text when it is written as
one, so every value of int64 and uint64 comes back exact; otherwise
they truncate the float64. Float returns a number's float64, the
number a numeric string spells, `1` for true and `0` for anything else.

Why: integers in JSON are often IDs and counters wider than a float64's
53 bits, and a caller asking for an integer means the integer the
document wrote.

Scenarios:
1. `int64-exact` (agent, provisional): -9223372036854775808 reads back
   by Int exactly.
2. `uint64-exact` (agent, provisional): 18446744073709551615 reads back
   by Uint exactly.
3. `from-string` (agent, provisional): the string `"12"` reads as the
   integer 12 and the float 12, and true as the float 1.

## SW-261126-ZGV3 Bool is true for true, a number other than zero, and a string other than empty, `0` and `false`

Status: active
Serves: STK-261126-TFYB

Bool returns true for the literal true, for a number not equal to zero
and for a string other than the empty string, `"0"` and `"false"`, and
false for anything else: false, null, an array, an object and a result
that does not exist.

Why: configuration documents spell flags as booleans, numbers and
strings alike; a lenient reading lets a caller ask the one question
whatever the spelling.

Scenarios:
1. `literals` (agent, provisional): true is true, and false and null
   are false.
2. `numbers` (agent, provisional): 1 and -2 are true and 0 is false.
3. `strings` (agent, provisional): `"yes"` and `"true"` are true, and
   `""`, `"0"` and `"false"` are false.

## SW-261126-26XP Exists is false only for a value the path did not find, and a JSON null exists

Status: active
Serves: STK-261126-ZQJR

Exists reports whether the result came from a value in the document: a
member whose value is `null` exists, with Type Null, and so does an
empty string, while a path that found nothing gives a result whose
Exists is false and whose Type is also Null.

Why: null and absent mean different things to a caller patching a
record, and the Type alone cannot tell them apart.

Scenarios:
1. `null-exists` (agent, provisional): `a` on `{"a":null}` exists with
   Type Null.
2. `missing-not` (agent, provisional): `b` on that document does not
   exist, and its Type is Null too.

## SW-261126-NKTP Array, Map and ForEach walk an array's elements or an object's members in document order

Status: active
Serves: STK-261126-TFYB

Array returns an array's elements in order; on a result that does not
exist it returns an empty slice, and on any other value a slice of that
one value. Map returns an object's members by key, and an empty map for
anything that is not an object. ForEach calls a function with each key
and value of an object, or each value of an array with a key that does
not exist, in document order, until the function returns false; on a
single value it calls it once with that value.

Why: the walks are how a caller goes from one lookup to all the values
under it, and document order is the order the author meant.

Scenarios:
1. `array-of-scalar` (agent, provisional): Array on the number 1 returns
   one result holding 1, and on a result that does not exist none.
2. `map` (agent, provisional): Map on `{"a":1,"b":"x"}` returns the two
   members, and on a string an empty map.
3. `for-each-order` (agent, provisional): ForEach visits an object's
   members in document order and stops when the function returns false.

## SW-261126-NRBX Value returns the result as the Go value encoding/json would decode it into an interface

Status: active
Serves: STK-261126-TFYB

Value returns bool for true and false, float64 for a number, string
for a string, nil for null and for nothing, []interface{} for an array
and map[string]interface{} for an object, nested values converted the
same way.

Why: a caller handing a value on to code that takes interface{}
expects the shapes encoding/json produces, and gets them without a
second decode.

Scenarios:
1. `value-types` (agent, provisional): an object holding each JSON type
   returns a map whose values are bool, float64, string, nil and
   []interface{}.

## SW-261126-Z23S Parse reads a whole document as a result, and Result.Get reads a path inside one

Status: active
Serves: STK-261126-ZQJR

Parse returns the document's first value as a result, without a path;
Result.Get on that result, or on any result holding an object or an
array, reads a path against its raw JSON as Get would.

Why: a caller that reads many paths from one value parses once and
walks from there, instead of naming the whole path every time.

Scenarios:
1. `parse-value` (agent, provisional): Parse of `  {"a":1} ` returns a
   JSON result, of `12` a number and of `"x"` a string.
2. `result-get` (agent, provisional): Parse of
   `{"name":{"last":"Anderson"}}`, then Get of `name`, then Get of
   `last`, returns `Anderson`.

## SW-261126-4WX9 Less orders results by type, Null, False, Number, String, True, JSON, and within a type by value

Status: active
Serves: STK-261126-TFYB

Less reports whether a result sorts before another: a lower Type first,
in the order Null, False, Number, String, True, JSON; two numbers by
value; two strings by bytes, or ignoring ASCII case when caseSensitive
is false; anything else by its raw text.

Why: a caller sorting results of mixed types needs one total order that
does not depend on how the document spelled them.

Scenarios:
1. `by-type` (agent, provisional): null sorts before false, and false
   before true.
2. `by-value` (agent, provisional): `"ABC"` sorts before `"abc"` with
   case, not without it, and 123 before 456.
