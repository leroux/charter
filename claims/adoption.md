# Adopting it anywhere

Code: internal/config, cmd/charter, examples, plugin

This area answers the problem that every project has its own languages,
its own layout and its own people. It covers how a repository takes the
charter up: its configuration, the languages the charter reads, the
skills, and the charter's own use of itself.

## CH-261126-33GK charter MUST find claims/charter.toml by walking up from the working directory, and MUST refuse a key it does not know

Status: active
Serves: CH-002, CH-261126-NS7T
See: CH-036, CH-037

- The command MUST walk up from the working directory to the first
  directory holding `claims/charter.toml`, which is the root; with none it
  MUST fail, naming where it started.
- The file holds: the charter's name; for each language it reads, a table
  with the directory its tests run in, the command that runs them, its
  constants module and, for Go, the generated module's path, and for
  TypeScript the typescript package it parses with and the specifier
  tests import the generated module by; the stakeholder tests, the probe
  path, the ignore and only patterns and the environment patterns; the
  digest; the declared surface kinds; and the doc patterns.
- A key it does not know, a value of the wrong type, or a language table
  without its test command MUST fail naming the line.
- Its paths are relative to the root.

Why: a charter configured in Go code could be read only by a Go build,
and a command that must be told where the charter is gets run in the
wrong place. A file beside the record, found the way git finds a
repository, works in any repository and from any directory in it.

Scenarios:
1. `found-walking-up` (owner, backlog 003): given a charter at the root,
   when charter runs from a directory three levels below, then it reads
   the charter whose `claims/charter.toml` is at the root.
2. `nearest-wins` (owner, backlog 003): given a nested charter inside
   another, when charter runs inside the nested one's tree, then it reads
   the nested charter, not the one around it.
3. `none-found-named` (owner, backlog 003): given no `claims/charter.toml`
   above the working directory, when charter runs, then it fails naming
   the directory it started from.
4. `keys-read` (owner, backlog 003): given a file setting every key, when
   it is read, then each value is as written.
5. `unknown-key-refused` (owner, backlog 003): given an unknown key, a
   value of the wrong type, or a language table without its test command,
   when the file is read, then each fails naming the line.

Changes:
- 2026-11-26: the text names the TypeScript table's compiler and import
  keys, which were read and named nowhere; keys-read says each value is
  yielded as written, where it said its paths were taken relative to the
  root, which no test showed, the root being the caller's to join (review
  of backlog 003).
- 2026-11-27: charter.toml names the doc patterns, the markdown files read
  as docs (CH-261127-3NY3); keys-read covers them as every key (backlog
  007).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-PQ7X charter MUST support each language as a built-in reader, generator and marker, and MUST read a language only when charter.toml names it

Status: active
Kind: structure
Serves: CH-002, CH-261126-NS7T
See: CH-020, CH-021, CH-023, CH-024

- The charter supports Go, Python and TypeScript, each by three parts
  built into the command: a reader that turns its source files into facts
  (the anchors, citations, references, environment reads and refused
  directives, with their files and lines); a generator that writes its
  symbol file; and a marker, the `verifies` its tests call.
- The gate's rules MUST read the facts alone and MUST NOT name a
  language, so a rule holds the same way in each.
- A language charter.toml does not name MUST NOT be read, and nothing it
  needs is run.
- Python is parsed by python3 with the standard `ast` module, and
  TypeScript by the typescript package charter.toml names; a missing one
  MUST fail naming that language's table.
- Another language is a change to the charter, not a plug-in.

Why: Go woven through the gate made a second language a parallel path
beside it, and a third would have been another. Facts of one shape let
every rule be written once and hold in every language, and a fixed set of
languages is one the charter can test.

Scenarios:
1. `same-rule-every-language` (owner, backlog 003): given a scenario no
   test cites, when the charter's tests are Go, Python or TypeScript, then
   the gate reports the same problem in the same words.
2. `unnamed-language-unread` (owner, backlog 003): given a Python file in
   a charter whose charter.toml has no Python table, when the gate runs,
   then the file is not read and python3 is never run.
3. `parser-missing-named` (owner, backlog 003): given a charter naming
   Python where no python3 runs, or TypeScript without its compiler, when
   the gate runs, then it fails naming that language's table.
4. `rules-name-no-language` (owner, backlog 003): given the package holding
   the gate's rules, when its imports are read, then it imports no
   language's reader, parser or generator.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-006 The charter MUST specify itself, and its own gate MUST pass

Status: active
Serves: CH-261126-NS7T

The tooling's decisions MUST be entries in this directory, in the form
the charter asks of any adopter, realized by the tooling, verified by its
tests, and checked by a gate of their own over the tooling alone.

Why: a form that cannot carry the charter's own decisions is the wrong
form, and a gate its maker does not run on itself is a claim, not a
proof.

Scenarios:
1. `self-gate-walks-the-tooling-only`: given the repository with its
   examples, when the charter's own gate runs, then it walks the tooling
   and nothing else, so an adopter's anchors and tests, the example
   adopter's among them, are invisible to it.

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
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-5JMD The loop's skills MUST ship as a plugin of this repository and MUST read the adopter's layout from charter.toml alone

Status: active
Serves: CH-261126-NS7T
See: CH-261126-33GK

- The claims skills (draft, implement, review, audit, probe and narrate)
  MUST ship as a Claude Code plugin, `charter`, of this repository:
  `plugin/` holds the plugin's manifest and the skills, and the
  marketplace manifest at the root names it, so an adopter installs it
  from the repository by name.
- A skill names the commands of the method and MUST find the adopter's
  layout (the claims directory, the languages and their test commands) in
  charter.toml; it MUST NOT name a file of one repository's layout.
- Narrate drafts a doc over a need, an area or a flow, and on a stale doc
  proposes an edit of what changed, never writing over a curated one; it
  writes each doc in the vocabulary of the reader the plan names. Audit
  also reads each doc against the code it references.
- Every skill MUST write in the voice the adopter's agent instructions
  (CLAUDE.md or AGENTS.md) give; where they give none, in a default the
  skills state: clear technical prose that leads with the point,
  connects each statement to the last and uses the field's terms. A
  skill MUST NOT impose a sentence count or a register.
- The README gives, as a template, the section an adopter puts in its
  CLAUDE.md.

Why: skills written for one product send an agent in another repository
to files that are not there. Skills that read the layout from the same
file the command reads go wherever the charter does, and a plugin is how
they get there.

Scenarios:
1. `plugin-lists-the-skills` (owner, backlog 003): given the repository,
   when the marketplace and plugin manifests are read, then the
   marketplace names the plugin and its directory, and the plugin holds
   the six skills, each a SKILL.md whose name is its directory's.
2. `no-fixed-layout` (owner, backlog 003): given the skills, when they are
   read, then none names `cmd/claimsgen`, `claims/README.md` or a path of
   one product, and each one needing the layout reads charter.toml.
3. `claude-md-template` (owner, backlog 003): given the README, when it is
   read, then it holds the CLAUDE.md section as a template naming the
   charter command and charter.toml.

Changes:
- 2026-11-27: the voice is the adopter's: skills follow the project's
  agent instructions and fall back to a stated default; narrate's "short
  plain words", draft's one-to-six sentences and audit's lint line go.
  The owner found the plain register too simple for complex projects.
- 2026-11-27: narrate writes in the vocabulary of the doc's reader; its
  first drafts for torii put plain stand-ins ("the name" for DNS) where an
  engineer needs the term, and the owner asked for the field's terms.
- 2026-11-27: narrate joins the skills, drafting docs and proposing edits
  to stale ones, and audit reads docs against their code; the plugin
  lists six skills (backlog 007).
- 2026-11-27: built: the plugin is plugin/, named by the marketplace at
  the root, and the README holds the CLAUDE.md template, there being no
  method document; the owner asked for the skills as a plugin.
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.
