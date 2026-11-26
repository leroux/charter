# cut

Code: cmd/cut

How cut -f selects fields from each line, what it does with a line that
has no delimiter, and which options it takes.

## SW-261126-972N cut -f prints the listed fields of each line, once each, in the order they stand

Status: active
Serves: STK-261126-69ZV

A field list is ranges separated by commas: N, N-M, N- (to the end of
the line) and -M (from the first), fields numbered from one. A line's
fields are the pieces between delimiters; cut prints the fields any
range selects, each once and in the order they stand in the line, joined
by the delimiter. A field past the end of the line is not printed, and a
line none of whose fields are selected prints as an empty line.

Why: GNU selects a set, not a sequence: -f3,1 is -f1,3, and a script
that wanted to reorder columns uses awk. A clone that reordered would
look more helpful and break every script that relied on GNU.

Scenarios:
1. `ascending-once` (agent, provisional): -f3,1 prints field one then
   field three, and -f2,2 prints field two once.
2. `ranges` (agent, provisional): -f2- prints from the second field on,
   -f-2 the first two, -f2-3 the second and third.
3. `past-the-end` (agent, provisional): -f5 on a line of two fields
   prints an empty line, and -f1,5 prints the first field alone.

## SW-261126-7RFJ A line without the delimiter is printed whole

Status: active
Kind: invariant
Serves: STK-261126-69ZV, STK-261126-3NDQ

cut prints a line that holds no delimiter as it stands, whatever the
field list.

Why: this is GNU's default, and it is what keeps header lines, comments
and blank lines in a script's output; GNU's -s is how a script asks to
drop them.

Scenarios:
1. `printed-whole` (agent, provisional): with -d: -f2, the line "plain"
   prints as "plain" and an empty line as an empty line.

## SW-261126-B85M cut takes a field list it must be given and a one-byte delimiter, a tab when none is given

Status: active
Serves: STK-261126-69ZV, STK-261126-WSGT

-f names the field list and is required; without it cut stops with
"you must specify a list of bytes, characters, or fields". -d names the
delimiter, one byte; a longer one stops cut with "the delimiter must be a
single character". Without -d the delimiter is a tab.

Why: GNU's defaults and refusals; a cut that guessed a delimiter or a
list would print columns nobody chose.

Scenarios:
1. `tab-default` (agent, provisional): without -d, "a\tb" -f2 prints b.
2. `single-byte` (agent, provisional): -d:: is refused with exit 1.
3. `list-required` (agent, provisional): cut without -f exits 1 with
   the message and prints nothing.
