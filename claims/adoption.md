# Adopting it anywhere

Code: internal/config, cmd/charter, example, example-polyglot

The problem this area answers: every project has its own languages, its
own layout and its own people. It holds how a repository takes the charter
up: its configuration, the languages the charter reads, the skills, and
the charter's own use of itself.

## CH-261126-33GK charter finds claims/charter.toml by walking up from the working directory, and refuses a key it does not know

Status: active
Serves: CH-002, CH-261126-NS7T
See: CH-036, CH-037

The command walks up from the working directory to the first directory
holding claims/charter.toml, and that directory is the root; with none it fails,
naming where it started. The file holds the charter's name; each language
it reads, a table with the directory its tests run in, the command that
runs them, its constants module and, for Go, the generated module's path,
and for TypeScript the typescript package it parses with and the specifier
tests import the generated module by; the stakeholder tests, the probe
path, the ignore and only patterns and the environment patterns; the
digest; and the declared surface kinds. A key it does not know, a value of
the wrong type, or a language table without its test command fails naming
the line. Its paths are relative to the root.

Why: a charter configured in Go code can be read only by a Go build, and a
command that must be told where the charter is gets run in the wrong
place; a file beside the record, found the way git finds a repository,
works in any repository and from any directory in it.

Scenarios:
1. `found-walking-up` (owner, backlog 003): run from a directory three
   below the root, charter reads the charter whose claims/charter.toml is
   at the root.
2. `nearest-wins` (owner, backlog 003): run inside a nested charter's
   tree, charter reads the nested charter, not the one around it.
3. `none-found-named` (owner, backlog 003): with no claims/charter.toml
   above the working directory, charter fails naming the directory it
   started from.
4. `keys-read` (owner, backlog 003): a file setting every key yields
   each value as written.
5. `unknown-key-refused` (owner, backlog 003): an unknown key, a value of
   the wrong type, and a language table without its test command each fail
   naming the line.

Changes:
- 2026-11-26: the text names the TypeScript table's compiler and import
  keys, which were read and named nowhere; keys-read says each value is
  yielded as written, where it said its paths were taken relative to the
  root, which no test showed, the root being the caller's to join (review
  of backlog 003).

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

## CH-006 The charter specifies itself, and its own gate is green

Status: active
Serves: CH-261126-NS7T

The tooling's decisions are entries in this directory, in the form the
charter asks of any adopter, realized by the tooling and verified by its
tests, and checked by a gate of their own over the tooling alone.

Why: a form that cannot carry the charter's own thirty decisions is the
wrong form, and a gate nobody runs on its own maker is a claim, not a
proof.

Scenarios:
1. `self-gate-walks-the-tooling-only`: the charter's own gate walks the
   tooling and nothing else, so an adopter's anchors and tests, the
   example adopter's among them, are invisible to it.

Changes:
- 2026-11-26: the charter became a repository of its own (Q-261015-3Q84). The entries are in the form the charter asks of
  any adopter, where they were in the product's, and the gate walks the
  repository's tooling, where it walked the tooling's directories inside
  the product; the example adopter is a module of its own, which every
  walk skips.
- 2026-11-26: a decision serving CH-261126-NS7T, where it was a need of
  its own, the needs being re-cut one a problem: the charter using its own
  method is how an adopter knows the method carries a real tool, a
  property of the project and not a promise to one person (backlog 004).

## CH-261126-5JMD The loop's skills ship as a plugin of this repository and read the layout from charter.toml alone

Status: planned
Serves: CH-261126-NS7T
See: CH-261126-33GK

The claims skills, draft, implement, review, audit and probe, ship as a
Claude Code plugin of this repository, which an adopter installs. A skill
names the commands of the method and finds the adopter's layout, the
claims directory, the languages and their test commands, in charter.toml;
it names no file of one repository's layout. The method's document gives
the section an adopter puts in its CLAUDE.md as a template.

Why: skills written for one product send an agent in another repository to
files that are not there; skills that read the layout from the file the
command reads go wherever the charter does, and a plugin is how they get
there.

Scenarios:
1. `plugin-lists-the-skills` (owner, backlog 003): the plugin's manifest
   lists the five skills, and each skill's file is in the plugin.
2. `no-fixed-layout` (owner, backlog 003): no skill names cmd/claimsgen,
   claims/README.md or a path of one product, and each one needing the
   layout reads charter.toml.
3. `claude-md-template` (owner, backlog 003): the method's document holds
   the CLAUDE.md section as a template naming the charter command and
   charter.toml.
