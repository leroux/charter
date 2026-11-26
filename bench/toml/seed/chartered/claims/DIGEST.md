<!-- Generated from the charter's documents. DO NOT EDIT. Regenerate: go run ./cmd/claimsgen -->

# The charter, read through

What the system does, without its code: 2 needs, what a person or the operator can rely on, and 18 decisions in 6 areas, how the system behaves at each boundary. Read the needs, then the map of the areas, then the areas in order. Each decision is its title, the sentence it rules, and its why; its id links to the whole entry, with its scenarios.

1. [What people need](#what-people-need)
2. [The areas](#the-areas)
3. [The decisions, area by area](#the-decisions-area-by-area)

## What people need

Each need is listed under the area that does most to meet it.

**Values**

- [STK-261126-Y8PN](needs.md#stk-261126-y8pn-a-valid-document-decodes-to-exactly-the-values-toml-10-gives-it) A valid document decodes to exactly the values TOML 1.0 gives it
- [STK-261126-HR5J](needs.md#stk-261126-hr5j-a-document-that-is-not-valid-toml-is-refused-with-where) A document that is not valid TOML is refused, with where

## The areas

| # | Area | What it holds | Decisions |
|---|---|---|---|
| 1 | [Lexing](lexing.md) | How the decoder reads the characters of a document before it reads any key or value: the encoding, whitespace and newlines, comments, and what may follow a pair or a header on its line. | 4 |
| 2 | [Keys](keys.md) | How the decoder reads a key and decides which table it names: bare and quoted keys, dotted keys, and the rule that a key is defined once. | 3 |
| 3 | [Values](values.md) | How the decoder reads the value of a pair: the forms it has built, which decode to exactly what TOML 1.0 says, and the forms it has not, which it refuses by name. | 5 |
| 4 | [Tables](tables.md) | How the decoder builds tables: headers, the tables a dotted key makes, and the rules that keep each table defined in one place. | 3 |
| 5 | [Errors](errors.md) | What a refusal carries back: where the problem is, in what form, and the rule that a refused document gives the program nothing. | 2 |
| 6 | [The decoder command](command.md) | How a program, toml-test above all, drives the decoder: the command toml-decoder, its input and its output. | 1 |

Which areas serve each need, with how many of their decisions:

| Need | Areas |
|---|---|
| [STK-261126-Y8PN](needs.md#stk-261126-y8pn-a-valid-document-decodes-to-exactly-the-values-toml-10-gives-it) | lexing 3 · keys 2 · values 4 · tables 1 · command 1 |
| [STK-261126-HR5J](needs.md#stk-261126-hr5j-a-document-that-is-not-valid-toml-is-refused-with-where) | lexing 4 · keys 2 · values 5 · tables 3 · errors 2 · command 1 |

## The decisions, area by area

### 1. Lexing

How the decoder reads the characters of a document before it reads any key or value: the encoding, whitespace and newlines, comments, and what may follow a pair or a header on its line.

- **[SW-261126-ZSR5](lexing.md#sw-261126-zsr5-a-document-is-utf-8-and-a-byte-order-mark-at-its-start-is-skipped)** A document is UTF-8, and a byte order mark at its start is skipped *(bound)*  
  Why: TOML 1.0 requires UTF-8, and a decoder that passed invalid bytes through would hand the program strings no other TOML reader would agree with; editors on some systems write the byte order mark unasked.
- **[SW-261126-M55Q](lexing.md#sw-261126-m55q-whitespace-is-space-and-tab-and-a-line-ends-in-lf-or-crlf)** Whitespace is space and tab, and a line ends in LF or CRLF *(bound)*  
  Why: TOML 1.0 defines exactly these, so files written on Windows and on Unix decode alike, and a stray carriage return, which no editor shows, is reported instead of silently splitting or joining lines.
- **[SW-261126-VW84](lexing.md#sw-261126-vw84-a-comment-runs-from-a-hash-to-the-end-of-its-line-and-holds-no-control-character-but-tab)** A comment runs from a hash to the end of its line and holds no control character but tab *(bound)*  
  Why: comments are how a configuration explains itself, so they must never change a value; the control characters TOML 1.0 forbids in them are refused so a corrupted file is noticed.
- **[SW-261126-WYRZ](lexing.md#sw-261126-wyrz-a-keyvalue-pair-or-a-header-ends-its-line)** A key/value pair or a header ends its line *(bound)*  
  Why: TOML 1.0 puts one pair on a line; a decoder that read on would turn a missing newline into a value nobody meant.

### 2. Keys

How the decoder reads a key and decides which table it names: bare and quoted keys, dotted keys, and the rule that a key is defined once.

- **[SW-261126-Y5J5](keys.md#sw-261126-y5j5-a-key-is-bare-or-quoted-and-a-quoted-key-may-hold-anything-a-string-may)** A key is bare or quoted, and a quoted key may hold anything a string may *(bound)*  
  Why: these are TOML 1.0's keys; a decoder that took more bare characters would accept files others refuse, and one that took fewer would refuse valid ones.
- **[SW-261126-F9XD](keys.md#sw-261126-f9xd-a-dotted-key-names-a-path-of-tables-with-whitespace-allowed-around-each-dot)** A dotted key names a path of tables, with whitespace allowed around each dot  
  Why: dotted keys are how TOML 1.0 writes nested settings without a header for each table.
- **[SW-261126-KKT9](keys.md#sw-261126-kkt9-a-key-is-defined-once-in-its-table-however-it-is-spelt)** A key is defined once in its table, however it is spelt *(invariant)*  
  Why: TOML 1.0 forbids defining a key twice, and a file that does is almost always a mistake (two people setting the same option); a decoder that let one win would decide silently which setting the program runs with.

### 3. Values

How the decoder reads the value of a pair: the forms it has built, which decode to exactly what TOML 1.0 says, and the forms it has not, which it refuses by name.

- **[SW-261126-QZWX](values.md#sw-261126-qzwx-a-basic-string-is-one-line-in-double-quotes-with-toml-10s-escapes-decoded)** A basic string is one line in double quotes, with TOML 1.0's escapes decoded *(bound)*  
  Why: these are exactly TOML 1.0's basic strings; an escape the specification does not define has no meaning a program could rely on.
- **[SW-261126-WZ4K](values.md#sw-261126-wz4k-a-literal-string-is-one-line-in-single-quotes-exactly-as-written)** A literal string is one line in single quotes, exactly as written *(bound)*  
  Why: literal strings are how TOML 1.0 writes Windows paths and regular expressions without doubling every backslash.
- **[SW-261126-JD2N](values.md#sw-261126-jd2n-an-integer-is-decimal-optionally-signed-without-leading-zeros-and-fits-in-64-bits)** An integer is decimal, optionally signed, without leading zeros, and fits in 64 bits *(bound)*  
  Why: TOML 1.0 requires integers to be exact and recommends 64 bits; a port of 0080 or a size that wrapped to a negative number would be a value nobody wrote.
- **[SW-261126-WN2V](values.md#sw-261126-wn2v-a-boolean-is-true-or-false-in-lower-case)** A boolean is true or false in lower case *(bound)*  
  Why: TOML 1.0 has exactly these two; reading True or yes as a boolean would accept files every other reader refuses.
- **[SW-261126-H72C](values.md#sw-261126-h72c-a-value-in-a-form-the-decoder-has-not-built-is-refused-by-name)** A value in a form the decoder has not built is refused by name *(bound)*  
  Why: a refusal that names the missing form tells the person what the decoder cannot do yet; reading the form as something else would hand the program a value of the wrong type.

### 4. Tables

How the decoder builds tables: headers, the tables a dotted key makes, and the rules that keep each table defined in one place.

- **[SW-261126-K75F](tables.md#sw-261126-k75f-a-header-defines-its-table-once-and-the-tables-on-its-path-exist-until-then)** A header defines its table once, and the tables on its path exist until then *(invariant)*  
  Why: TOML 1.0 lets a file define a table's subtables before or after the table itself, and forbids defining a table twice, so a setting belongs to one section a reader can find.
- **[SW-261126-596G](tables.md#sw-261126-596g-a-table-made-by-dotted-keys-and-a-table-made-by-a-header-are-each-defined-one-way-only)** A table made by dotted keys and a table made by a header are each defined one way only *(invariant)*  
  Why: TOML 1.0 forbids defining a table both ways, so each table is written in one place; the subtable under a dotted table is the specification's own example of what stays valid.
- **[SW-261126-DY2F](tables.md#sw-261126-dy2f-a-key-that-holds-a-value-is-not-a-table)** A key that holds a value is not a table *(invariant)*  
  Why: a setting is either a value or a section, never both; letting one replace the other would drop a setting the file wrote.

### 5. Errors

What a refusal carries back: where the problem is, in what form, and the rule that a refused document gives the program nothing.

- **[SW-261126-RVMG](errors.md#sw-261126-rvmg-a-refusal-names-the-line-and-column-where-the-problem-begins)** A refusal names the line and column where the problem begins *(bound)*  
  Why: a person fixes a configuration file in an editor that shows lines and columns; a refusal without them leaves the person searching.
- **[SW-261126-BPDX](errors.md#sw-261126-bpdx-a-refused-document-gives-no-table-at-all)** A refused document gives no table at all *(invariant)*  
  Why: a partial configuration is a configuration nobody wrote, and a program that ran on the first half of a file would fail far from the cause.

### 6. The decoder command

How a program, toml-test above all, drives the decoder: the command toml-decoder, its input and its output.

- **[SW-261126-Y24J](command.md#sw-261126-y24j-the-decoder-command-speaks-toml-tests-decoder-protocol)** The decoder command speaks toml-test's decoder protocol *(structure)*  
  Why: toml-test, the suite every TOML implementation is checked against, drives decoders by exactly this protocol, so the decoder is checked against the specification's own cases.
