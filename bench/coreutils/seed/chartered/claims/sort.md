# sort

Code: cmd/sort

How sort orders the lines of its input: by bytes, by number with -n, and
in reverse with -r.

## SW-261126-YMCY sort orders whole lines by their bytes, and -r reverses the order

Status: active
Kind: invariant
Serves: STK-261126-69ZV

Without -n, lines compare as unsigned byte strings, as the C locale's
collation does: uppercase before lowercase, a prefix before the longer
line. -r reverses the comparison. sort reads every input before it
prints, the lines of all its files together.

Why: in the C locale GNU's order is the bytes' order; scripts rely on it
being the same on every machine, which is why they set LC_ALL=C.

Scenarios:
1. `bytewise` (agent, provisional): "b", "B", "a", "A", "ab" sort to A,
   B, a, ab, b.
2. `reverse` (agent, provisional): -r prints the same lines in the
   opposite order.

## SW-261126-GKC7 sort -n orders by the leading number, exactly, and breaks ties by bytes

Status: active
Serves: STK-261126-69ZV

-n reads each line's leading number: blanks, an optional minus, digits
and an optional point with more digits; a plus sign, an exponent or a
thousands separator ends it. A line that does not start with a number
reads as zero, as does negative zero. Numbers compare exactly, however
many digits they have. Lines whose numbers are equal compare by bytes,
and -r reverses that comparison too.

Why: GNU's -n is a string comparison of decimal numbers, not a float
parse, so huge numbers and 2.5 against 2.50 order the same on every
machine; the byte tie-break makes the output one order, not one of
several.

Scenarios:
1. `numeric-order` (agent, provisional): -1, 9, 10, 2.5 sort to -1,
   2.5, 9, 10.
2. `non-number-zero` (agent, provisional): "abc", "+5" and "-0" read as
   zero and sort between -1 and 1.
3. `ties-bytewise` (agent, provisional): "2.50" and "2.5" are equal
   numbers and sort by bytes, and -rn reverses the tie too.
4. `exact-digits` (agent, provisional): numbers of twenty digits that
   differ in the last order correctly.
