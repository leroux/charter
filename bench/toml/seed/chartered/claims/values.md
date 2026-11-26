# Values

Code: values.go

How the decoder reads the value of a pair: the forms it has built, which
decode to exactly what TOML 1.0 says, and the forms it has not, which it
refuses by name.

## SW-261126-QZWX A basic string is one line in double quotes, with TOML 1.0's escapes decoded

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

A basic string runs from a double quote to the next unescaped double
quote on the same line. The escapes \b \t \n \f \r \" \\ decode to their
characters, and \uXXXX and \UXXXXXXXX to the Unicode scalar value they
name; any other escape, a \u or \U naming a surrogate or a value beyond
U+10FFFF, a control character other than tab, and a newline before the
closing quote are refused.

Why: these are exactly TOML 1.0's basic strings; an escape the
specification does not define has no meaning a program could rely on.

Scenarios:
1. `escapes-decoded` (agent, provisional): each of the seven short
   escapes decodes to its character.
2. `unicode-escapes` (agent, provisional): \u00E9 and \U0001F600 decode
   to their characters.
3. `unknown-escape-refused` (agent, provisional): \x41, \e and \q are
   refused.
4. `surrogate-refused` (agent, provisional): \uD800 is refused.
5. `one-line` (agent, provisional): a basic string with no closing quote
   before the end of the line is refused.
6. `control-refused` (agent, provisional): a basic string holding a raw
   NUL is refused.

## SW-261126-WZ4K A literal string is one line in single quotes, exactly as written

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

A literal string runs from a single quote to the next single quote on the
same line and is its text exactly: a backslash is a backslash. A control
character other than tab inside it is refused.

Why: literal strings are how TOML 1.0 writes Windows paths and regular
expressions without doubling every backslash.

Scenarios:
1. `backslash-kept` (agent, provisional): 'C:\Users\n' decodes with its
   backslashes and the letter n.
2. `control-refused` (agent, provisional): a literal string holding a
   raw NUL is refused.

## SW-261126-JD2N An integer is decimal, optionally signed, without leading zeros, and fits in 64 bits

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

An integer is a run of decimal digits with an optional + or - in front,
and no leading zero unless it is zero itself; +0 and -0 are zero. It
decodes to a signed 64-bit integer; one outside that range is refused,
never wrapped or rounded.

Why: TOML 1.0 requires integers to be exact and recommends 64 bits; a
port of 0080 or a size that wrapped to a negative number would be a value
nobody wrote.

Scenarios:
1. `signed-decimal` (agent, provisional): 42, +17, -5 and -0 decode to
   their integers.
2. `leading-zero-refused` (agent, provisional): 007 and +01 are refused.
3. `out-of-range-refused` (agent, provisional): 9223372036854775808 is
   refused and -9223372036854775808 decodes.

## SW-261126-WN2V A boolean is true or false in lower case

Status: active
Kind: bound
Serves: STK-261126-Y8PN, STK-261126-HR5J

The values true and false, in lower case, are booleans. Any other
spelling is refused.

Why: TOML 1.0 has exactly these two; reading True or yes as a boolean
would accept files every other reader refuses.

Scenarios:
1. `lowercase-only` (agent, provisional): true and false decode, and
   True, FALSE and truthy are refused.

## SW-261126-H72C A value in a form the decoder has not built is refused by name

Status: active
Kind: bound
Serves: STK-261126-HR5J

The decoder builds a subset of TOML 1.0's values: basic and literal
strings on one line, decimal integers and booleans. Floats, integers in
other bases or with underscores, multi-line strings, arrays, inline
tables, arrays of tables, and dates and times are refused with a message
that names the form, never read as a string or a number of another kind.

Why: a refusal that names the missing form tells the person what the
decoder cannot do yet; reading the form as something else would hand the
program a value of the wrong type.

Scenarios:
1. `float-refused` (agent, provisional): 3.14, 1e6 and inf are refused as
   floats.
2. `array-refused` (agent, provisional): [1, 2] and [[t]] are refused as
   arrays and arrays of tables.
3. `date-refused` (agent, provisional): 1979-05-27 and 07:32:00 are
   refused as dates and times.
4. `other-forms-refused` (agent, provisional): 0xff, 1_000, a multi-line
   string and an inline table are each refused by name.
