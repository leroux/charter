# 003. The charter is a binary over a whole repository, any module and any language, configured by a file beside the record

Status: ruled · Size: L · See: 001

The gate today is a Go test inside one Go module: its options are a Go
variable, its root is the nearest go.mod, its walk skips nested modules,
and citations are read from Go and TypeScript only. torii (leroux/torii-qubes)
is the adopter that shows the limit: its Go is in a nested module
(`proxy/`), its root has no go.mod, and half its behaviour is a Python
tool in dom0. This story makes the charter standalone: one `charter`
binary, a `claims/charter.toml` beside the record, Python as a third
language, and no compatibility with the go-test-hosted form.

The record (needs, area documents, questions, the digest) is markdown
already and does not change form. What changes is where the gate runs,
what it reads its configuration from, and which languages it reads.

## Ruling

Ruled by the owner, 2026-11-26.

| # | Question | Ruling |
|---|---|---|
| 1 | Where does the gate run? | One `charter` binary is the gate: `charter check`. No go.mod is needed at the root. |
| 2 | What does the binary do? | `charter check \| gen \| next PREFIX \| lint \| report \| diff REF \| mutate ID \| export`, folding in claimsgen and claimsmap. An adopter writes no main. |
| 3 | How does it find the charter? | It walks up from the working directory to `claims/charter.toml`; the directory holding `claims/` is the root. |
| 4 | Is the go-test-hosted gate kept? | No. `chartertest.Check` as a gate, `gate.Options`, `claims/gate/`, `gate_test.go` and the adopter's `cmd/claimsgen` are removed; the charter's own `claims/` and `example/` move to the new form. |
| 5 | What form is the configuration? | TOML, `claims/charter.toml`: the prefixes, the languages and their test commands, stakeholder tests, ignores, the digest, surface kinds. |
| 6 | Does the walk cross modules? | Yes: the walk is of the repository, nested Go modules included; it skips only what the config ignores and the always-ignored names. |
| 7 | Where do generated Go symbols live? | `claims/` is a Go module of its own holding `specs.go`; a Go module in the repository that cites requires it, through a `replace` within the repository. |
| 8 | Where does a Python anchor go? | `# charter:realizes ID, ...` (and `concerns`) on the comment lines directly above a `def` or `class`, or above its decorators. One anywhere else is refused, naming where to move it. |
| 9 | What is a Python citation? | `verifies(self, claims.SW_..., "handle", ...)` as the first statement of a `test_*` method of a `unittest.TestCase`, or `verifies(None, ...)` as the first statement of a pytest `test_*` function; the same refusals as CH-023's. |
| 10 | Where does the Python helper come from? | Generated: `claims/claims.py` holds the symbols and `verifies`, with no dependency beyond the standard library. Nothing is installed. |
| 11 | How does a Python run report per scenario? | With `CHARTER_EVENTS=path` set, `verifies` appends one JSON line per citation and the test's outcome to that file; unset, it records nothing. Any runner works. |
| 12 | How is Python parsed? | `python3` running an embedded script over the standard `ast`, only when charter.toml has a `[python]` section; a missing `python3` fails naming that section. |
| 13 | How does report run tests? | Each language section in charter.toml names its directory and test command; `charter report` runs them and joins the events. |
| 14 | How does an adopter add a surface kind? | Declared in charter.toml as a command whose output lines are the surface keys of that kind. The built-in kinds are reconsidered against it when drafted. |
| 15 | How does an adopter get the skills? | Backlog 001 items 3 and 4 are done here: the `claims-*` skills read the layout from charter.toml and ship as a Claude Code plugin from this repository; the method gives the CLAUDE.md section as a template. |
| 16 | What proves it? | A second example adopter, Go and Python with a nested Go module and no root go.mod, run by `adopter_test.go` beside ratelimit (moved to the new form), each passing `charter check` and failing it on a copy missing a citation. |
| 17 | Refactor or rewrite? | Rewrite, built against the CH record: `cmd/charter`, a core that knows no language (the record, the registry, the gate's rules, report, digest, lint, diff, export), and a package per language. Code of the old tooling is ported where it is not bound to Go (the area parser, ids, the digest, style) and dropped where it is. Done when every active CH scenario is cited by a passing test again. |
| 18 | How is a language supported? | Built in, and only built in: a language is a reader (source files to facts: anchors, citations, references, refused directives), a generator (the symbol file its tests import) and a marker (its runtime record of each citation's outcome). The gate's rules see facts only. A new language is a change to the charter. |
| 19 | Which languages? | Go, Python and TypeScript, TypeScript ported onto the same three parts. |
| 20 | What becomes of mutate? | Dropped: CH-261005-VFJR is deprecated, and the command and its code go. It may return, per language, if the benchmarks (002) show it earns its cost. |

## To do, in order

1. Draft the CH entries from the ruling (`claims-draft`): changes to
   CH-020 to CH-025 (symbols, anchors, citations, markers, now per
   language), CH-031, CH-033, CH-036, CH-037, CH-040 and the rest of
   `tooling.md` that name the Go options or claimsgen; CH-261005-VFJR
   deprecated; new entries for charter.toml, the repository walk, the
   three parts of a language, Python and the command-based surface kind.
   The owner reviews them before any code.
2. The rewrite (`claims-implement`): the core and the three languages,
   built entry by entry against the record, the old packages removed as
   their entries move; then `claims/` and `example/` in the new form.
3. The Go and Python example.
4. The skills as a plugin and the CLAUDE.md template (001, 3 and 4).
5. Onboard torii: its own story there, from the plan ruled with the owner
   on 2026-11-26 (`claims/` at its root, STK and SW, six needs, ten areas
   with verdict active first, design.md's Decisions reduced to a map,
   machine-file keys as a surface when that area goes active, no CI yet).
