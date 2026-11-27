# Charter

The charter tooling and method, specified by its own charter. Start with
`README.md`, then the area documents in `claims/` (needs.md first).

Write everything here (the record, skills, docs, commit messages and
replies) as @WRITING.md describes.

## Layout

- `cmd/charter`: the one command; `internal/`: everything it is made of,
  importable by nothing outside. The core (`record`, `facts`, `gate`,
  `check`, `report`, `digest`, `diff`, `docs`) names no language;
  `internal/lang/{golang,python,typescript}` are the readers, generators and
  markers.
- `claims/`: the charter's own charter, the `CH` entries, with
  `charter.toml` and the generated claims module (`go.mod`, `specs.go`).
- `examples/`: `ratelimit`, Go alone, with a doc; `slugs`, a nested Go
  module and Python with no go.mod at its root. Each is a charter of its
  own that imports nothing of the tooling; `cmd/charter/adopter_test.go`
  runs their gates and holds them to that.
- `backlog/`: the stories in progress or held, one a file.

## How claims are written

The owner chose this form on 2026-11-27. Write every entry in `claims/` to it.

- **Title:** the component that carries the obligation, with a BCP 14
  keyword, as one claim a reviewer could disagree with: "The proxy MUST
  check an allowlisted connection's TLS server name before relaying it".
- **Body:** RFC 2119 bullets, one requirement per bullet, by default.
  Use connected prose instead when the entry is a single requirement, or
  when its clauses only make sense together (a sequence, a cause).
- **Key words:** MUST, MUST NOT, SHOULD, SHOULD NOT and MAY, in capitals
  only, with their BCP 14 meanings (RFC 2119, RFC 8174), as `needs.md`
  states.
- **Why:** explain the problem or threat, how the rule meets it, and what
  the rule costs or avoids.
- **Scenarios:** one numbered line each, in given / when / then form, with
  a concrete input and an observable outcome.

## The loop, here

- A change to the tooling is a change to an entry in `claims/`: its text,
  a scenario, and a dated `Changes:` line saying what changed in the
  ruling and why, in about three lines; how a review found it belongs in
  the pull request. `go run
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
  and the examples', and the examples' reports. It needs `pnpm install`
  for the TypeScript compiler and vitest.

The skills that run the loop live in `plugin/skills/claims-*`, the
`charter` plugin other repositories install (CH-261126-5JMD);
`.claude/skills/` links to them so this repository's sessions have them
too. Edit them under `plugin/`.
