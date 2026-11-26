# The command line and input

Code: internal/input, internal/opts

What every tool shares: how its arguments split into options and files,
and how it reads its input and splits it into lines.

## SW-261126-ZG5V A tool reads standard input or the named files in order, split into lines, a last unterminated one included

Status: active
Kind: structure
Serves: STK-261126-3NDQ

A tool with no file operand reads standard input; with file operands it
reads each file in the order named. A tool splits its input on newline
bytes into lines; a last line without a newline is a line all the same,
and an empty input has no lines. Every line a tool prints ends with a
newline.

Why: GNU treats a missing final newline as a line that ends there, and so
does every script that pipes a file made by printf; the tools share one
reader so they cannot differ on it.

Scenarios:
1. `stdin-when-none` (agent, provisional): with no names, the input is
   standard input's bytes.
2. `files-in-order` (agent, provisional): with two names, the inputs are
   the two files in the order named.
3. `unterminated-last-line` (agent, provisional): "a\nb" splits into the
   lines a and b, as "a\nb\n" does.
4. `empty-input-no-lines` (agent, provisional): an empty input has no
   lines, and "\n" has one empty line.

## SW-261126-2QDT Options are parsed as GNU parses them, and one a tool does not know is refused

Status: active
Serves: STK-261126-WSGT, STK-261126-69ZV

Short options may be clustered (-lw) and may come before, between or
after the file operands; "--" ends the options and "-" alone is an
operand. An option that takes a value takes the rest of its cluster
(-d:) or the next argument (-d :). A long option is --name or
--name=value. An option the tool does not know, or a value missing,
is an error in GNU's words, and the tool stops before reading input:
exit status 1, or 2 for sort, as GNU's.

Why: scripts write options every way getopt accepts, and a parser that
took fewer would turn a working command line into an error or, worse,
read an option as a file.

Scenarios:
1. `clustered` (agent, provisional): -lw is -l and -w.
2. `after-operands` (agent, provisional): an option after a file operand
   is still an option.
3. `double-dash` (agent, provisional): after "--", "-l" is an operand.
4. `value-attached-or-next` (agent, provisional): -d: and -d : give the
   same value.
5. `long-forms` (agent, provisional): --fields=2 and --fields 2 are -f 2.
6. `unknown-refused` (agent, provisional): an unknown short option is
   "invalid option -- 'Z'" and an unknown long one "unrecognized option".
