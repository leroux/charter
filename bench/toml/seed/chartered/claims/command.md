# The decoder command

Code: cmd/toml-decoder

How a program, toml-test above all, drives the decoder: the command
toml-decoder, its input and its output.

## SW-261126-Y24J The decoder command speaks toml-test's decoder protocol

Status: active
Kind: structure
Serves: STK-261126-Y8PN, STK-261126-HR5J

toml-decoder reads one document on stdin. When it decodes, the command
writes it on stdout as toml-test's tagged JSON, tables as objects and
every value as {"type": ..., "value": ...} with the value a string, and
exits 0. When the document is refused, it writes the error on stderr,
nothing on stdout, and exits 1.

Why: toml-test, the suite every TOML implementation is checked against,
drives decoders by exactly this protocol, so the decoder is checked
against the specification's own cases.

Scenarios:
1. `tagged-json` (agent, provisional): a string, an integer and a boolean
   come out tagged string, integer and bool, each value a string.
2. `refusal-exit` (agent, provisional): a refused document exits 1 with
   the error on stderr and nothing on stdout.
