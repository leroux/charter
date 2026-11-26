# Charter

The charter tooling and method, specified by its own charter. Start with
`README.md`, then the area documents in `claims/` (needs.md first).

## Layout

- `cmd/charter`: the one command; `internal/`: everything it is made of,
  importable by nothing outside. The core (`record`, `facts`, `gate`,
  `check`, `report`, `digest`, `lint`, `diff`, `export`) names no language;
  `internal/lang/{golang,python,typescript}` are the readers, generators and
  markers.
- `claims/`: the charter's own charter, the `CH` entries, with
  `charter.toml` and the generated claims module (`go.mod`, `specs.go`).
- `example/`: ratelimit, Go alone; `example-polyglot/`: slugs, a nested Go
  module and Python with no go.mod at its root. Each is a charter of its
  own that imports nothing of the tooling; `cmd/charter/adopter_test.go`
  runs their gates and holds them to that.
- `backlog/`: the work toward a standalone method, one story a file.

## The loop, here

- A change to the tooling is a change to an entry in `claims/`: its text,
  a scenario, and a dated, append-only `Changes:` line. `go run
  ./cmd/charter diff origin/main` names a change the record does not say.
- Ids come from `go run ./cmd/charter next CH` (in an example, the built
  `build/charter next SW` from its directory). Never write one by hand,
  reuse or edit one.
- After editing an area document run `make claims`; commit what it writes.
- Code carries `//charter:realizes CH-…` as the last line of the doc
  comment of the declaration that carries the decision (Python: `#
  charter:realizes` on the comment lines above the def; TypeScript: a `//
  charter:realizes` line comment); a test's first statement is
  `claims.Verifies(t, claims.CH_…, "handle")`.
- A choice no entry rules is a question in `claims/questions.md` (id from
  `charter next Q`), and the scenario it produced is `(agent,
  provisional)` until the owner answers.
- `make check` before pushing: the Go tests with the charter's own gate
  and the examples', the examples' lint and report, the viewer. It needs
  `pnpm install` for the TypeScript compiler and vitest.

The skills in `.claude/skills/claims-*` still name `claimsgen` and one
product's layout in places; CH-261126-5JMD (planned) makes them a plugin
that reads charter.toml.
