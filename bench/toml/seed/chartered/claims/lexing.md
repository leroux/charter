# Lexing

Code: lex.go

How the decoder reads the characters of a document before it reads any
key or value: the encoding, whitespace and newlines, comments, and what
may follow a pair or a header on its line.

## SW-261126-ZSR5 A document is UTF-8, and a byte order mark at its start is skipped

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

A document must be valid UTF-8 from its first byte to its last; one that
is not is refused at the first byte that is not. A UTF-8 byte order mark
at the very start is skipped and is not part of the document.

Why: TOML 1.0 requires UTF-8, and a decoder that passed invalid bytes
through would hand the program strings no other TOML reader would agree
with; editors on some systems write the byte order mark unasked.

Scenarios:
1. `invalid-utf8-refused` (agent, provisional): a document with a byte
   that is not UTF-8 inside a string is refused.
2. `bom-skipped` (agent, provisional): a document that starts with a
   byte order mark decodes as the same document without it.

## SW-261126-M55Q Whitespace is space and tab, and a line ends in LF or CRLF

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

Between the parts of a line, whitespace is any run of spaces and tabs and
nothing else. A line ends in a line feed or a carriage return followed by
a line feed; a carriage return alone is refused.

Why: TOML 1.0 defines exactly these, so files written on Windows and on
Unix decode alike, and a stray carriage return, which no editor shows, is
reported instead of silently splitting or joining lines.

Scenarios:
1. `tab-separates` (agent, provisional): a pair written with tabs around
   its equals sign decodes as with spaces.
2. `crlf-ends-line` (agent, provisional): a document with CRLF line ends
   decodes as the same document with LF.
3. `lone-cr-refused` (agent, provisional): a carriage return not followed
   by a line feed is refused.

## SW-261126-VW84 A comment runs from a hash to the end of its line and holds no control character but tab

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

A hash outside a string starts a comment, on a line of its own or after a
pair or a header, and the comment ends at the end of the line. A control
character other than tab inside a comment is refused. A hash inside a
string is part of the string.

Why: comments are how a configuration explains itself, so they must
never change a value; the control characters TOML 1.0 forbids in them
are refused so a corrupted file is noticed.

Scenarios:
1. `comment-ignored` (agent, provisional): comments on their own lines
   and after a pair and a header add nothing to the document.
2. `control-in-comment-refused` (agent, provisional): a comment holding a
   NUL or a DEL is refused.
3. `hash-in-string-is-text` (agent, provisional): a hash inside a quoted
   string is kept in the string.

## SW-261126-WYRZ A key/value pair or a header ends its line

Status: active
Kind: bound
Serves: STK-261126-HR5J

After a pair's value or a header's closing bracket only whitespace and a
comment may follow before the end of the line. Anything else on the line
is refused.

Why: TOML 1.0 puts one pair on a line; a decoder that read on would turn
a missing newline into a value nobody meant.

Scenarios:
1. `two-pairs-one-line-refused` (agent, provisional): two pairs on one
   line are refused.
2. `text-after-header-refused` (agent, provisional): a word after a
   header's closing bracket is refused.
