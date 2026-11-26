# Tables

Code: tables.go

How the decoder builds tables: headers, the tables a dotted key makes,
and the rules that keep each table defined in one place.

## SW-261126-K75F A header defines its table once, and the tables on its path exist until then

Status: active
Kind: invariant
Serves: STK-261126-Y8PN, STK-261126-HR5J

A header [a.b.c] defines the table c and makes it the table the pairs
below it go into; the tables a and a.b on its path are created if
missing. A table created that way may be defined later by its own header,
once; a table already defined by a header is refused the second time.

Why: TOML 1.0 lets a file define a table's subtables before or after the
table itself, and forbids defining a table twice, so a setting belongs to
one section a reader can find.

Scenarios:
1. `header-defines` (agent, provisional): pairs under [server] and under
   [server.tls] go into those tables.
2. `parent-defined-later` (agent, provisional): [a.b] then [a] is valid,
   with a holding b and a's own pairs.
3. `defined-twice-refused` (agent, provisional): [a] twice is refused, as
   is [a.b] twice with [a] between.

## SW-261126-596G A table made by dotted keys and a table made by a header are each defined one way only

Status: active
Kind: invariant
Serves: STK-261126-HR5J

The tables on the path of a dotted key are created as dotted tables in
the current table. A header may not define a dotted table, though it may
define a new table under one, and a dotted key may not reach into a table
a header defined or created on its path.

Why: TOML 1.0 forbids defining a table both ways, so each table is
written in one place; the subtable under a dotted table is the
specification's own example of what stays valid.

Scenarios:
1. `header-over-dotted-refused` (agent, provisional): `a.b = 1` then
   [a] is refused.
2. `dotted-into-header-refused` (agent, provisional): [a.b] then [a]
   holding `b.c = 1` is refused.
3. `subtable-under-dotted` (agent, provisional): [fruit] with
   `apple.color = "red"`, then [fruit.apple.texture], is valid.

## SW-261126-DY2F A key that holds a value is not a table

Status: active
Kind: invariant
Serves: STK-261126-HR5J

A key defined with a value cannot be used as a table afterwards, by a
header or by a dotted key; the attempt is refused.

Why: a setting is either a value or a section, never both; letting one
replace the other would drop a setting the file wrote.

Scenarios:
1. `header-over-value-refused` (agent, provisional): `a = 1` then [a.b]
   is refused.
2. `dotted-over-value-refused` (agent, provisional): `a = 1` then
   `a.b = 2` is refused.
