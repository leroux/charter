# Charter

Charter keeps what a piece of software must do in the repository, above
the code, and holds the code to it. Needs (what a person can rely on) and
decisions (how the system behaves at a boundary) are written once, as
entries in markdown area documents under `claims/`. Each decision is
anchored to the code that carries it and each of its scenarios is cited by
a test, in Go, Python or TypeScript. A static gate, `charter check`, fails
the tree when an active decision has no code that claims it, a scenario no
test cites, or a surface (an environment variable, or a kind the charter
declares) no entry accounts for.

The gate proves the links exist; it does not prove the code right or a
test meaningful. That is review's job, and the charter lets the reviewer
read the entries, their tests and `charter report` instead of every line,
while `charter diff` names any change to the record nobody wrote down. It
works in any repository, existing or new, in any layout: Go modules nested
anywhere, Python and TypeScript beside them. The tooling specifies itself:
its own decisions are the `CH` entries in `claims/`.

**Start with [docs/](docs/README.md):** [why it exists](docs/why.md),
[the loop](docs/loop.md), [adopting it](docs/adopting.md), and, for
changing it, [the architecture](docs/architecture.md).

## What is here

| Path | What it is |
|---|---|
| `cmd/charter` | The one command: `check`, `gen`, `next`, `anchors`, `questions`, `report`, `diff`, `curate` |
| `internal/` | What it is made of: a gate that names no language, and the readers and generators in `lang/{golang,python,typescript,plain}` |
| `claims/` | The charter's own charter: the tooling specified by the method it implements |
| `docs/` | How it works and why, with references the gate checks |
| `examples/ratelimit/` | **ratelimit**, a Go adopter kept as a real one would be, with a doc |
| `examples/slugs/` | **slugs**, a nested Go module and Python beside it, with no go.mod at its root |
| `plugin/` | The `charter` Claude Code plugin: the skills of the loop, draft, implement, review, audit, probe and narrate |
| `backlog/` | The stories in progress or held, one a file |

## The command

```text
charter check               run the gate; print every problem and exit 1 on any, warnings on stderr
charter gen                 write the generated files and the digest from the record, and update the docs
charter next PREFIX [N]     print N new ids with the prefix
charter anchors ID...       an entry's code, tests, relations and docs, as JSON
charter questions [-all]    the open questions; -all the answered ones too
charter report [-md]        run every language's tests and say which scenarios passed
charter diff REV            name what the record changed since REV with no Changes line
charter curate [-by NAME] FILE...
                            record that a person read each doc, with a hash of what it references
```

It works on the charter whose `claims/charter.toml` is nearest above the
working directory. Exit status 1 means it found something; 2 means it
could not run.

## Adopting it

`examples/ratelimit` and `examples/slugs` are the worked examples, and
[docs/adopting.md](docs/adopting.md) walks through both. Nothing of the
tooling is imported: the command is a binary, and the files it generates
are self-contained.

1. Install the command: `go install github.com/leroux/charter/cmd/charter@latest`,
   or pin it in a tools module. The repository is private for now, so
   this needs git credentials that can read it.
2. Write `claims/charter.toml` at the root of the tree the charter covers:

   ```toml
   name = "ratelimit"
   stakeholder-tests = ["cmd/ratelimit/*_test.go"]
   env-patterns = ["^RATELIMIT_[A-Z0-9_]+$"]
   digest = "DIGEST.md"
   docs = ["docs/*.md"]

   [go]
   module = "example.com/ratelimit/claims"   # the generated claims module
   dir = "."                                 # where the tests run
   test = "go test ./..."
   env = "env"                               # the constants package
   ```

   `[python]` and `[typescript]` tables take `dir`, `test` and `env` the
   same way; `[typescript]` also `compiler` (the typescript package) and
   `import`. `[[surface]]` tables declare a surface kind by a command whose
   output lines are its keys, and `[[files]]` tables name files no
   language reads whose comment lines may anchor an entry.
3. Write `claims/needs.md` and an area document per area (`# Title`, a
   `Code:` line, a summary, then entries); `claims/reading-order.txt`
   names the area documents in the order a reader meets them. Each entry
   is a `## ID Title` heading, field lines (`Status`, `Kind`, `Serves`,
   `See`, `Accounts`), the text, then `Why:`, `Scenarios:`, `Guidance:`
   and `Changes:`. Ids come from `charter next SW`.
4. `charter gen` writes the generated files into `claims/`: a Go module
   (`go.mod`, `specs.go`) your modules require through a `replace`,
   `claims.py`, `index.ts`, and the digest.
5. Anchor code with `//charter:realizes SW-…` as the last line of a
   declaration's doc comment (Python: `# charter:realizes SW-…` on the
   comment lines above the def or class; TypeScript: a line comment), and
   open each test with its citation: `claims.Verifies(t, claims.SW_…,
   "handle")`, `claims.verifies(self, claims.SW_…, "handle")` in unittest
   or `verifies(None, …)` in pytest, `verifies(claims.SW_…, "handle")` in
   vitest.
6. `charter check` is the gate; `charter report -md` runs the tests and
   says which scenarios passed; `charter diff origin/main` names a change
   to the record with no Changes line. Run all three in CI.
7. Docs tell the story the record does not: a walkthrough, a design, a
   guide, with references to entries and declarations the gate checks.
   An agent drafts them (claims-narrate), a person reads each and runs
   `charter curate -by NAME FILE`, and check warns again when anything it
   references changes ([docs/docs.md](docs/docs.md)).

## The skills

The loop's skills ship as the `charter` plugin for Claude Code. Install it
in the repository that keeps a charter:

```
/plugin install charter --marketplace leroux/charter
```

or, from a shell, `claude plugin marketplace add leroux/charter` and then
`claude plugin install charter@charter --scope project`. The skills run as
`/charter:claims-draft`, `/charter:claims-implement`,
`/charter:claims-review`, `/charter:claims-audit`, `/charter:claims-probe`
and `/charter:claims-narrate`, and read the layout from
`claims/charter.toml`. Installing needs git credentials that can read
this repository. A cloud session does not load plugins a repository's
settings name; add the plugin to it at claude.ai/customize/plugins.

Put this section in the repository's CLAUDE.md, with its prefix and paths:

```markdown
## Charter

What this software must do is recorded in `claims/`, configured by
`claims/charter.toml`: needs in `claims/needs.md`, decisions in the area
documents, each anchored to the code that carries it and verified by the
tests that cite it. `charter check` is the gate; run it before pushing.

- A change in behaviour is a change to an entry: its text, a scenario, and
  a dated `Changes:` line saying what changed and why. `charter diff
  origin/main` names a change the record does not say.
- Ids come from `charter next SW`; never write one by hand.
- After editing an area document run `charter gen`; commit what it writes.
- Code carries `//charter:realizes SW-…` as the last line of the doc
  comment of the declaration that carries the decision; a test's first
  statement is `claims.Verifies(t, claims.SW_…, "handle")`.
- A choice no entry rules is a question in `claims/questions.md`, and the
  scenario it produced is `(agent, provisional)` until the owner answers.
- The `charter` plugin's skills run the loop: claims-draft,
  claims-implement, claims-review, claims-audit, claims-probe and
  claims-narrate for docs.
```

## Working on the charter

```sh
pnpm install          # the TypeScript compiler and vitest the tests use
make check            # the Go tests, the charter's own gate and the examples'
go run ./cmd/charter report -md
```

A change to the tooling is a change to an entry in `claims/`, with a
dated `Changes:` line, and `go run ./cmd/charter diff origin/main` names
one that has none. The examples' gates run among the tooling's tests
(`cmd/charter/adopter_test.go`), so a change that would break an adopter
fails here first. [docs/architecture.md](docs/architecture.md) says where
to start for each kind of change.
