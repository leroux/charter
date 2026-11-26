# Needs

What a Go program that calls gjson can rely on. A need is verified from
the caller's side: the exported functions, called on a document and a
path as a program would call them.

## STK-261126-ZQJR A path returns the value the documented syntax names, and nothing when it names none

Status: active

A caller writes a path in the syntax the README documents and gets the
value it names in the document; when the path names no value, the
caller gets a result that does not exist, never some other value.

Why: gjson replaces decoding a whole document into structs; a caller who
cannot predict what a path returns would have to check every answer,
and the point of a one-line lookup is lost.

Scenarios:
1. `readme-paths` (agent, provisional): the README's path examples, on
   the README's document, return the values the README shows.
2. `missing-is-absent` (agent, provisional): a path to a key or an
   index the document does not have returns a result that does not
   exist.

## STK-261126-8FW8 Get never panics, whatever the document or the path

Status: active

Any string may be passed as the document and any string as the path:
gjson returns a result, possibly one that does not exist, and never
panics.

Why: callers feed gjson documents from the network and paths from
configuration; a panic in a lookup takes the whole program down with
it.

Scenarios:
1. `random-input` (agent, provisional): random bytes as the document,
   with random paths, never panic.
2. `truncated` (agent, provisional): a document cut short at every
   length, read with every path of its full form, never panics.

## STK-261126-TFYB A caller reads a result as the Go value it needs, and a big integer keeps every digit

Status: active

A result converts to a string, a number, a boolean, a slice or a map
without the caller parsing JSON, and an integer read as an integer has
every digit the document wrote, beyond what a float64 holds.

Why: the caller wants a Go value, not JSON text; an ID that loses its
last digits to a float silently names a different record.

Scenarios:
1. `as-go-values` (agent, provisional): one document's string, number,
   boolean, array and object each read back as the Go value they hold.
2. `big-integer-digits` (agent, provisional): the integer
   9007199254740993 reads back as that integer, and as that string.

## STK-261126-94D2 A caller can tell whether a document is valid JSON before trusting it

Status: active

gjson reads without validating, so a caller who cannot trust a
document's source asks gjson first whether it is valid JSON.

Why: a lenient reader is fast on good input and unpredictable on bad;
the caller decides which documents deserve the lenient read.

Scenarios:
1. `check-first` (agent, provisional): a well-formed document is
   reported valid and the same document with a trailing comma is not.
