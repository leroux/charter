# Languages

Code: the language packages

What the charter reads of each language it supports, and the three parts
each of them is made of.

## CH-261126-PQ7X A language is a reader, a generator and a marker, built in, and read only when charter.toml names it

Status: active
Kind: structure
Serves: CH-002, CH-261126-NS7T
See: CH-020, CH-021, CH-023, CH-024

The charter supports Go, Python and TypeScript, each by three parts built
into the command: a reader that turns its source files into facts, the
anchors, citations, references, environment reads and refused directives
with their files and lines; a generator that writes its symbol file; and a
marker, the verifies its tests call. The gate's rules read the facts alone
and name no language, so a rule holds the same way in each. A language
charter.toml does not name is not read, and nothing it needs is run.
Python is parsed by python3 over the standard ast module and TypeScript by
the typescript package charter.toml names; a missing one fails naming that
language's table. Another language is a change to the charter, not a
plug-in.

Why: Go woven through the gate made a second language a parallel path
beside it, and a third would have been another; facts of one shape let
every rule be written once and hold in every language, and a fixed set of
languages is one the charter can test.

Scenarios:
1. `same-rule-every-language` (owner, backlog 003): a scenario no test
   cites is the same problem, in the same words, whether the charter's
   tests are Go, Python or TypeScript.
2. `unnamed-language-unread` (owner, backlog 003): a Python file in a
   charter whose charter.toml has no Python table is not read, and python3
   is never run.
3. `parser-missing-named` (owner, backlog 003): a charter naming Python
   where no python3 runs, or TypeScript without its compiler, fails naming
   that language's table.
4. `rules-name-no-language` (owner, backlog 003): the package holding the
   gate's rules imports no language's reader, parser or generator.
