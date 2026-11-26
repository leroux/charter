// Code generated from claims by the charter gate. DO NOT EDIT.
// Regenerate: go run ./cmd/claimsgen

package claims

import "github.com/leroux/charter"

// STK_261126_HR5J is STK-261126-HR5J: A document that is not valid TOML is refused, with where
//
// A document that is not valid TOML 1.0, or that uses a part of TOML the
// decoder has not built, is refused as a whole: the program gets no values
// from it, and the person gets the line and column where the problem
// begins, and what it is.
var STK_261126_HR5J = &charter.Spec{ID: "STK-261126-HR5J", Title: "A document that is not valid TOML is refused, with where", Status: charter.Active}

// STK_261126_Y8PN is STK-261126-Y8PN: A valid document decodes to exactly the values TOML 1.0 gives it
//
// A document that is valid TOML 1.0, and uses only the parts of TOML the
// decoder has built, decodes to the tables, keys and values the
// specification gives it: every value with its TOML type, every key in the
// table it belongs to, nothing added and nothing dropped.
var STK_261126_Y8PN = &charter.Spec{ID: "STK-261126-Y8PN", Title: "A valid document decodes to exactly the values TOML 1.0 gives it", Status: charter.Active}

// SW_261126_596G is SW-261126-596G: A table made by dotted keys and a table made by a header are each defined one way only
//
// The tables on the path of a dotted key are created as dotted tables in
// the current table. A header may not define a dotted table, though it may
// define a new table under one, and a dotted key may not reach into a table
// a header defined or created on its path.
var SW_261126_596G = &charter.Spec{ID: "SW-261126-596G", Title: "A table made by dotted keys and a table made by a header are each defined one way only", Status: charter.Active, Kind: "invariant"}

// SW_261126_BPDX is SW-261126-BPDX: A refused document gives no table at all
//
// Decode returns either the whole document or an error and no table: a
// document refused on its last line gives the program none of the lines
// before it.
var SW_261126_BPDX = &charter.Spec{ID: "SW-261126-BPDX", Title: "A refused document gives no table at all", Status: charter.Active, Kind: "invariant"}

// SW_261126_DY2F is SW-261126-DY2F: A key that holds a value is not a table
//
// A key defined with a value cannot be used as a table afterwards, by a
// header or by a dotted key; the attempt is refused.
var SW_261126_DY2F = &charter.Spec{ID: "SW-261126-DY2F", Title: "A key that holds a value is not a table", Status: charter.Active, Kind: "invariant"}

// SW_261126_F9XD is SW-261126-F9XD: A dotted key names a path of tables, with whitespace allowed around each dot
//
// A key of parts joined by dots names the last part in the table the parts
// before it name, each relative to the current table. Whitespace around a
// dot is ignored; a dot inside a quoted part is part of its name.
var SW_261126_F9XD = &charter.Spec{ID: "SW-261126-F9XD", Title: "A dotted key names a path of tables, with whitespace allowed around each dot", Status: charter.Active}

// SW_261126_H72C is SW-261126-H72C: A value in a form the decoder has not built is refused by name
//
// The decoder builds a subset of TOML 1.0's values: basic and literal
// strings on one line, decimal integers and booleans. Floats, integers in
// other bases or with underscores, multi-line strings, arrays, inline
// tables, arrays of tables, and dates and times are refused with a message
// that names the form, never read as a string or a number of another kind.
var SW_261126_H72C = &charter.Spec{ID: "SW-261126-H72C", Title: "A value in a form the decoder has not built is refused by name", Status: charter.Active, Kind: "bound"}

// SW_261126_JD2N is SW-261126-JD2N: An integer is decimal, optionally signed, without leading zeros, and fits in 64 bits
//
// An integer is a run of decimal digits with an optional + or - in front,
// and no leading zero unless it is zero itself; +0 and -0 are zero. It
// decodes to a signed 64-bit integer; one outside that range is refused,
// never wrapped or rounded.
var SW_261126_JD2N = &charter.Spec{ID: "SW-261126-JD2N", Title: "An integer is decimal, optionally signed, without leading zeros, and fits in 64 bits", Status: charter.Active, Kind: "bound"}

// SW_261126_K75F is SW-261126-K75F: A header defines its table once, and the tables on its path exist until then
//
// A header [a.b.c] defines the table c and makes it the table the pairs
// below it go into; the tables a and a.b on its path are created if
// missing. A table created that way may be defined later by its own header,
// once; a table already defined by a header is refused the second time.
var SW_261126_K75F = &charter.Spec{ID: "SW-261126-K75F", Title: "A header defines its table once, and the tables on its path exist until then", Status: charter.Active, Kind: "invariant"}

// SW_261126_KKT9 is SW-261126-KKT9: A key is defined once in its table, however it is spelt
//
// A key that is already defined in a table is refused the second time,
// whether it is written the same way, bare against quoted, or as the last
// part of a dotted key. No definition overwrites or merges with another.
var SW_261126_KKT9 = &charter.Spec{ID: "SW-261126-KKT9", Title: "A key is defined once in its table, however it is spelt", Status: charter.Active, Kind: "invariant"}

// SW_261126_M55Q is SW-261126-M55Q: Whitespace is space and tab, and a line ends in LF or CRLF
//
// Between the parts of a line, whitespace is any run of spaces and tabs and
// nothing else. A line ends in a line feed or a carriage return followed by
// a line feed; a carriage return alone is refused.
var SW_261126_M55Q = &charter.Spec{ID: "SW-261126-M55Q", Title: "Whitespace is space and tab, and a line ends in LF or CRLF", Status: charter.Active, Kind: "bound"}

// SW_261126_QZWX is SW-261126-QZWX: A basic string is one line in double quotes, with TOML 1.0's escapes decoded
//
// A basic string runs from a double quote to the next unescaped double
// quote on the same line. The escapes \b \t \n \f \r \" \\ decode to their
// characters, and \uXXXX and \UXXXXXXXX to the Unicode scalar value they
// name; any other escape, a \u or \U naming a surrogate or a value beyond
// U+10FFFF, a control character other than tab, and a newline before the
// closing quote are refused.
var SW_261126_QZWX = &charter.Spec{ID: "SW-261126-QZWX", Title: "A basic string is one line in double quotes, with TOML 1.0's escapes decoded", Status: charter.Active, Kind: "bound"}

// SW_261126_RVMG is SW-261126-RVMG: A refusal names the line and column where the problem begins
//
// Every refusal is a *ParseError holding the line and the column, both
// counted from 1, of the first character of what is wrong, and a message
// that says what. Columns count characters, not bytes. Its text is
// `toml: line L, column C: message`.
var SW_261126_RVMG = &charter.Spec{ID: "SW-261126-RVMG", Title: "A refusal names the line and column where the problem begins", Status: charter.Active, Kind: "bound"}

// SW_261126_VW84 is SW-261126-VW84: A comment runs from a hash to the end of its line and holds no control character but tab
//
// A hash outside a string starts a comment, on a line of its own or after a
// pair or a header, and the comment ends at the end of the line. A control
// character other than tab inside a comment is refused. A hash inside a
// string is part of the string.
var SW_261126_VW84 = &charter.Spec{ID: "SW-261126-VW84", Title: "A comment runs from a hash to the end of its line and holds no control character but tab", Status: charter.Active, Kind: "bound"}

// SW_261126_WN2V is SW-261126-WN2V: A boolean is true or false in lower case
//
// The values true and false, in lower case, are booleans. Any other
// spelling is refused.
var SW_261126_WN2V = &charter.Spec{ID: "SW-261126-WN2V", Title: "A boolean is true or false in lower case", Status: charter.Active, Kind: "bound"}

// SW_261126_WYRZ is SW-261126-WYRZ: A key/value pair or a header ends its line
//
// After a pair's value or a header's closing bracket only whitespace and a
// comment may follow before the end of the line. Anything else on the line
// is refused.
var SW_261126_WYRZ = &charter.Spec{ID: "SW-261126-WYRZ", Title: "A key/value pair or a header ends its line", Status: charter.Active, Kind: "bound"}

// SW_261126_WZ4K is SW-261126-WZ4K: A literal string is one line in single quotes, exactly as written
//
// A literal string runs from a single quote to the next single quote on the
// same line and is its text exactly: a backslash is a backslash. A control
// character other than tab inside it is refused.
var SW_261126_WZ4K = &charter.Spec{ID: "SW-261126-WZ4K", Title: "A literal string is one line in single quotes, exactly as written", Status: charter.Active, Kind: "bound"}

// SW_261126_Y24J is SW-261126-Y24J: The decoder command speaks toml-test's decoder protocol
//
// toml-decoder reads one document on stdin. When it decodes, the command
// writes it on stdout as toml-test's tagged JSON, tables as objects and
// every value as {"type": ..., "value": ...} with the value a string, and
// exits 0. When the document is refused, it writes the error on stderr,
// nothing on stdout, and exits 1.
var SW_261126_Y24J = &charter.Spec{ID: "SW-261126-Y24J", Title: "The decoder command speaks toml-test's decoder protocol", Status: charter.Active, Kind: "structure"}

// SW_261126_Y5J5 is SW-261126-Y5J5: A key is bare or quoted, and a quoted key may hold anything a string may
//
// A bare key is one or more ASCII letters, digits, underscores and dashes,
// and is a string even when it is all digits. A quoted key is a basic or a
// literal string, decoded as the string would be, and may be empty. Any
// other character where a key begins is refused.
var SW_261126_Y5J5 = &charter.Spec{ID: "SW-261126-Y5J5", Title: "A key is bare or quoted, and a quoted key may hold anything a string may", Status: charter.Active, Kind: "bound"}

// SW_261126_ZSR5 is SW-261126-ZSR5: A document is UTF-8, and a byte order mark at its start is skipped
//
// A document must be valid UTF-8 from its first byte to its last; one that
// is not is refused at the first byte that is not. A UTF-8 byte order mark
// at the very start is skipped and is not part of the document.
var SW_261126_ZSR5 = &charter.Spec{ID: "SW-261126-ZSR5", Title: "A document is UTF-8, and a byte order mark at its start is skipped", Status: charter.Active, Kind: "bound"}
