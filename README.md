# Charter

Charter records what a system does as **needs** (what a person can rely
on) and **decisions** (how the system behaves at a boundary), each written
once in a markdown area document, anchored to the code that realizes it and
verified, scenario by scenario, by the tests that name it, in Go, Python
and TypeScript. A gate, `charter check`, fails the tree when an active
decision has no code that claims it, a scenario no test cites, or a surface
(an environment variable, or any kind the charter declares) no decision
accounts for.

The gate proves the links exist; it does not prove the code right or a test
meaningful. That is what review is for, and the charter gives the reviewer
the entries and their tests to read instead of the code.

It is meant for any repository, existing or new, in any layout: Go modules
nested anywhere, Python and TypeScript beside them.

## What is here

| Path | What it is |
|---|---|
| `cmd/charter` | The one command: `check`, `gen`, `lint`, `next`, `anchors`, `questions`, `report`, `diff`, `export`, `curate` |
| `internal/` | What it is made of: a core that names no language, and `lang/{golang,python,typescript}` |
| `claims/` | The charter's own charter: the tooling specified by the method it implements |
| `example/` | **ratelimit**, a Go adopter kept as a real one would be |
| `example-polyglot/` | **slugs**, a nested Go module and Python beside it, with no go.mod at its root |
| `claims-viewer/` | A web viewer of an exported charter |
| `.claude/skills/claims-*` | The agent skills of the loop: draft, implement, review, audit, probe |
| `backlog/` | The work toward a standalone method |

## Adopting it

`example/` and `example-polyglot/` are the worked examples; every piece
below is a file there. Nothing of the tooling is imported: the command is
a binary, and the files it generates are self-contained.

1. Install the command: `go install github.com/leroux/charter/cmd/charter@latest`,
   or pin it in a tools module as torii pins its build tools.
2. Write `claims/charter.toml` at the root of the tree the charter covers:

   ```toml
   name = "ratelimit"
   stakeholder-tests = ["cmd/ratelimit/*_test.go"]
   env-patterns = ["^RATELIMIT_[A-Z0-9_]+$"]
   digest = "DIGEST.md"

   [go]
   module = "example.com/ratelimit/claims"   # the generated claims module
   dir = "."                                 # where the tests run
   test = "go test ./..."
   env = "env"                               # the constants package

   # [python] and [typescript] tables take dir, test and env the same way;
   # [typescript] also compiler (the typescript package) and import.
   # [[surface]] tables declare a surface kind: a kind and a command whose
   # output lines are its keys. [[files]] tables name files no language
   # reads (patterns, and their comment marker) whose comment lines may
   # anchor an entry: `# charter:realizes SW-…` in a firewall ruleset.
   # docs = ["docs/*.md"] names the docs: prose over the record and the
   # code, whose references the gate checks (step 7).
   ```

3. Write `claims/needs.md` and an area document per area (`# Title`, a
   `Code:` line, a summary, then entries); `reading-order.txt` names the
   area documents in the order a reader meets them.
4. `charter gen` writes the generated files into `claims/`: a Go module
   (`go.mod`, `specs.go`) that your modules require through a `replace`,
   `claims.py`, `index.ts`, and the digest. `charter next SW` gives a new id.
5. Anchor code with `//charter:realizes SW-…` as the last line of a
   declaration's doc comment (Python: `# charter:realizes SW-…` on the
   comment lines above the def or class; TypeScript: a line comment), and
   open each test with its citation: `claims.Verifies(t, claims.SW_…,
   "handle")`, `claims.verifies(self, claims.SW_…, "handle")` in unittest
   or `verifies(None, …)` in pytest, `verifies(claims.SW_…, "handle")` in
   vitest.
6. `charter check` is the gate. `charter lint`, `report -md` and `diff
   origin/main` do the rest of the loop; `report` runs each language's test
   command and reads each citation's outcome from the events file.
7. Docs tell the story the record does not: a walkthrough, a design, a
   guide (`example/docs/walkthrough.md`). A link whose text opens with an
   id references an entry, ``[`Limiter.Allow`](../limiter/limiter.go)``
   references a declaration, and `<!-- charter:embed SW-… -->` with its
   closing comment embeds an entry; the gate fails one that does not
   resolve, and `charter gen` keeps titles, line anchors and embeds
   current. An agent drafts a doc (`claims-narrate`), a person reads it
   and runs `charter curate FILE`, and lint names it again when anything
   it references has changed since.

## Working on the charter

```sh
pnpm install          # the TypeScript compiler and vitest the tests use
make check            # the Go tests and the charter's own gate, the examples, the viewer
go run ./cmd/charter lint
go run ./cmd/charter report -md
```

The charter specifies itself: a change to the tooling is a change to an
entry in `claims/`, with a dated `Changes:` line, and `go run
./cmd/charter diff origin/main` names one that has none. The examples'
gates run among the tooling's tests (`cmd/charter/adopter_test.go`), so a
change that would break an adopter fails here first.
