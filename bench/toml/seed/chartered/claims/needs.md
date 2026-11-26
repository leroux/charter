# Needs

What a person who writes a configuration file in TOML, and the program
that reads it through this decoder, can rely on. A need is verified from
their side: the toml-decoder command, fed a document on stdin and read
back, as toml-test or a script would drive it.

## STK-261126-Y8PN A valid document decodes to exactly the values TOML 1.0 gives it

Status: active

A document that is valid TOML 1.0, and uses only the parts of TOML the
decoder has built, decodes to the tables, keys and values the
specification gives it: every value with its TOML type, every key in the
table it belongs to, nothing added and nothing dropped.

Why: a configuration file is read by a program its author never sees; a
decoder that reads a value as something the specification does not say
changes the program's behaviour without anyone writing a different file.

Scenarios:
1. `config-decodes` (agent, provisional): a configuration file with
   comments, bare, quoted and dotted keys, strings, integers, booleans and
   tables decodes through the command to the tagged values TOML 1.0 gives
   it.
2. `types-kept` (agent, provisional): a quoted number decodes as a string
   and a bare one as an integer.

## STK-261126-HR5J A document that is not valid TOML is refused, with where

Status: active

A document that is not valid TOML 1.0, or that uses a part of TOML the
decoder has not built, is refused as a whole: the program gets no values
from it, and the person gets the line and column where the problem
begins, and what it is.

Why: a decoder that guesses at a broken file hands the program a
configuration nobody wrote; one that refuses without saying where leaves
the person searching the file by hand.

Scenarios:
1. `refused-with-position` (agent, provisional): a document that defines
   a key twice is refused with a failing exit, the line and column of the
   second definition on stderr, and nothing on stdout.
2. `unbuilt-refused` (agent, provisional): a document with a float is
   refused, the message naming floats, rather than read as something
   else.
