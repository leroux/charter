# 003. The charter is a binary over a whole repository, any module and any language, configured by a file beside the record

Status: in progress · Size: L · See: 001, 002

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
| 21 | Does an adopter's code import the charter? | No. Each language's generated file is self-contained: `claims/specs.go` declares its own Spec and Verifies and `claims/go.mod` requires nothing; `claims/index.ts` exports verifies; `claims/claims.py` as ruling 10. Only the binary depends on this module. |
| 22 | How does every language report a citation's outcome? | One sink: with `CHARTER_EVENTS` set, every generated verifies, Go's included, appends a JSON line per citation with the test and its outcome. The report sets it and runs each configured command; go test -json is not read. |
| 23 | How does one charter keep out of another's tree? | A directory below the root that holds its own `claims/charter.toml` is another charter's and the walk skips it; charter.toml's ignore covers the rest. |
| 24 | What becomes of observations? | They move to `claims/observations.md`, parsed like entries; each language's generated file has observes; a probe is a test under a path charter.toml names. CH-025's Go-literal rules go. |
| 25 | Which surface kinds are built in? | Environment variables, in all three languages: each reader reports the reads, a read by literal is refused, and charter.toml names each language's constants module. Build tags, OpenAPI operations and event kinds become declared kinds. |
| 26 | Q-261015-N4HJ | Yes: the one-file-per-spec form is no longer read; CH-014 is deprecated. |
| 27 | Q-261015-Y682 | Yes: `charter gen` is the one writer of generated files; the update variables go. |
| 28 | Q-261015-CPW5 | Yes: a generated symbol's doc comment is its entry's title. |
| 29 | Q-261015-5YV3 | Yes: lint -baseline and the style baseline are retired; CH-261009-SNX2 is deprecated. |
| 30 | How is the rewrite laid out? | `cmd/charter` and `internal/`: config, record, facts, lang/{golang,python,typescript}, walk, surface, gate (importing no lang package), report, digest, lint, diff, export. Nothing is importable by an adopter. `charter.go`, `chartertest/`, `cli/`, `cmd/claimsgen` and `cmd/claimsmap` are removed. |
| 31 | Which TOML parser? | github.com/BurntSushi/toml, the binary's one dependency. |
| 32 | In what order does it land? | Bottom-up, one commit a step with its tests: config, record and facts; the Go language; gate, walk and surfaces; Python and TypeScript; gen, report, lint, diff, digest and export; then claims/ and example/ move to charter.toml, the old packages go, and the gate is green. |

## Done

- The rewrite, built against the record (rulings 1 to 32): `cmd/charter`
  and `internal/`, Go, Python and TypeScript, `claims/charter.toml`, the
  repository walk, the events file, observations in markdown, declared
  surface kinds; the old packages removed. The charter's own gate is green
  and its report passes every active scenario.
- Both examples in the new form: ratelimit, and slugs (`example-polyglot/`),
  a nested Go module and Python with no go.mod at its root, whose need is
  verified across the two languages.

## To do, in order

1. The skills as a plugin and the CLAUDE.md template (001, 3 and 4):
   CH-261126-5JMD, planned. The skills still name `claimsgen`.
2. The benchmarks' seeds (002) are adopters of the earlier, go-test form:
   their `@CHARTER@` replace now reaches a tree without `chartertest`.
   Move them to charter.toml, or pin the harness to the commit before the
   rewrite.
3. Onboard torii: its own story there, from the plan ruled with the owner
   on 2026-11-26 (`claims/` at its root, STK and SW, six needs, ten areas
   with verdict active first, design.md's Decisions reduced to a map,
   machine-file keys as a surface when that area goes active, no CI yet).

## Charter

Drafted 2026-11-26 and approved by the owner the same day in design review, before any build.

New, active since the build, but CH-261126-5JMD, still planned:

- CH-261126-33GK charter finds claims/charter.toml by walking up from the working directory, and refuses a key it does not know (`config.md`)
- CH-261126-PQ7X A language is a reader, a generator and a marker, built in, and read only when charter.toml names it (`languages.md`)
- CH-261126-NTRR A declared surface kind is a command run from the root whose output lines are its keys (`gate.md`)
- CH-261126-HK6S An observation is a heading in observations.md, with its subject, status, date or source and the entries it supports (`form.md`)
- CH-261126-5JMD The loop's skills ship as a plugin of this repository and read the layout from charter.toml alone (`tooling.md`)

Changed, each with a dated line: CH-002, CH-003, CH-004, CH-261126-NS7T
(`needs.md`); CH-012, CH-015 (`form.md`); CH-020, CH-021, CH-023, CH-024,
CH-025 (`anchoring.md`); CH-031, CH-032, CH-033, CH-034, CH-036, CH-037
(`gate.md`); CH-040, CH-041, CH-042, CH-261005-6ZBS, CH-261015-9ZTY,
CH-261009-QQC8 (`tooling.md`).

Deprecated: CH-014, CH-261009-SNX2, CH-261005-VFJR.

Answered: Q-261015-N4HJ, Q-261015-Y682, Q-261015-CPW5, Q-261015-5YV3.

The changed active entries state the rewrite's behaviour ahead of its
code, so until the build the gate names their new scenarios as untested,
the old tests' dropped handles, and the update variables as unaccounted.
