# Keys

Code: keys.go, tables.go

How the decoder reads a key and decides which table it names: bare and
quoted keys, dotted keys, and the rule that a key is defined once.

## SW-261126-Y5J5 A key is bare or quoted, and a quoted key may hold anything a string may

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

A bare key is one or more ASCII letters, digits, underscores and dashes,
and is a string even when it is all digits. A quoted key is a basic or a
literal string, decoded as the string would be, and may be empty. Any
other character where a key begins is refused.

Why: these are TOML 1.0's keys; a decoder that took more bare characters
would accept files others refuse, and one that took fewer would refuse
valid ones.

Scenarios:
1. `bare-chars` (agent, provisional): keys of letters, digits, dashes and
   underscores, and one of digits alone, decode as their text.
2. `quoted-any` (agent, provisional): a quoted key with spaces, a dot and
   an escape decodes as one key holding that text.
3. `empty-quoted` (agent, provisional): an empty quoted key is a key.
4. `bad-bare-refused` (agent, provisional): a bare key with a character
   outside the bare set, ASCII or not, is refused.

## SW-261126-F9XD A dotted key names a path of tables, with whitespace allowed around each dot

Status: active
Serves: STK-261126-Y8PN

A key of parts joined by dots names the last part in the table the parts
before it name, each relative to the current table. Whitespace around a
dot is ignored; a dot inside a quoted part is part of its name.

Why: dotted keys are how TOML 1.0 writes nested settings without a
header for each table.

Scenarios:
1. `dotted-nests` (agent, provisional): `a.b.c = 1` decodes to a table a
   holding a table b holding c.
2. `space-around-dot` (agent, provisional): `a . b = 1` is the same key
   as `a.b = 1`.
3. `quoted-part` (agent, provisional): `site."example.com" = true` is one
   key named example.com in the table site.

## SW-261126-KKT9 A key is defined once in its table, however it is spelt

Status: active
Kind: invariant
Serves: STK-261126-HR5J

A key that is already defined in a table is refused the second time,
whether it is written the same way, bare against quoted, or as the last
part of a dotted key. No definition overwrites or merges with another.

Why: TOML 1.0 forbids defining a key twice, and a file that does is
almost always a mistake (two people setting the same option); a decoder
that let one win would decide silently which setting the program runs
with.

Scenarios:
1. `duplicate-refused` (agent, provisional): a key defined twice in the
   root table is refused.
2. `spelling-same-key-refused` (agent, provisional): `name` and `"name"`
   in one table are the same key and refused.
3. `dotted-duplicate-refused` (agent, provisional): `a.b = 1` then
   `a.b = 2` is refused.
