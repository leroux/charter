# Charter

The charter tooling and method, specified by its own charter. Start with
`README.md`, then `claims/DIGEST.md` if it exists, else the area documents
in `claims/` (needs.md first).

## Layout

- `charter.go`, `chartertest/`: the library; `cli/`: the command line;
  `cmd/claimsgen`, `cmd/claimsmap`: the commands over the charter's own
  charter.
- `claims/`: the charter's own charter, the `CH` entries; its gate is
  `claims/gate_test.go` with options in `claims/gate/`.
- `example/`: ratelimit, an adopter in a module of its own that requires
  this one through a `replace`. Keep it an adopter: it copies nothing of
  the tooling, and `adopter_test.go` holds it to that and runs its gate.
- `backlog/`: the work toward a standalone method, one story a file.

## The loop, here

- A change to the tooling is a change to an entry in `claims/`: its text,
  a scenario, and a dated, append-only `Changes:` line. `go run
  ./cmd/claimsgen diff origin/main` names a change the record does not say.
- Ids come from `go run ./cmd/claimsgen next CH` (in `example/`: `go run
  ./cmd/claimsgen next SW` or `STK`). Never write one by hand, reuse or
  edit one.
- After editing an area document run `go run ./cmd/claimsgen` (or `make
  claims`), and in `example/` its own `go run ./cmd/claimsgen`; commit
  what they write.
- Code carries `//charter:realizes CH-…` as the last line of the doc
  comment of the declaration that carries the decision; a test's first
  statement is `chartertest.Verifies(t, claims.CH_…, "handle")`.
- A choice no entry rules is a question in `claims/questions.md` (id from
  `claimsgen next Q`), and the scenario it produced is `(agent,
  provisional)` until the owner answers.
- `make check` before pushing: the Go tests with the charter's own gate
  and the example's, the example's lint, the viewer. It needs `pnpm
  install` for the TypeScript compiler the citation extractor parses with.

The skills in `.claude/skills/claims-*` still assume one product's layout
in places; backlog 001 lists what is left to make them general.
