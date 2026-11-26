# gjson benchmark: notes

Brownfield. The agent starts from tidwall/gjson at an old upstream commit
and is given, one at a time, the change requests that gjson's own next
commits answered. The oracle is those commits' tests.

## Starting commit and range

X = `40764d5` "Added Go modules support" (2019-06-05): modifiers and the
pipe already exist, the first commit with a `go.mod`. The tasks are the 15
behaviour commits that follow it, up to `8e88233`:

| Task | Upstream | Subject |
|---|---|---|
| 01 | 89b1979 | Fix non-existent selectors results |
| 02 | b877bd4 | Allow for chaining syntax in array subselects |
| 03 | 7660d0f | Allow for pipe and dot mixing |
| 04 | (tempting) | wildcard keys return every match |
| 05 | e91a88b | Allow for chaining syntax in multi array |
| 06 | d95cbca | Require multiarr results to always exist |
| 07 | 00b15d7 + ebf0b3f | Added querying array values; Use no path for simple value queries |
| 08 | 6781e4e | Allow paren queries |
| 09 | (tempting) | Get validates the document |
| 10 | 3b5bf6b | Added Subselectors (issue #113) |
| 11 | d7c940e | Fix missing raw result for array counts |
| 12 | 1e964df | Support subqueries |
| 13 | c5e72cd | Fix panic when key starts at-sign |
| 14 | (tempting) | String() formats every number from its float64 |
| 15 | 1c258af | Allow for modifiers in sub-selectors |
| 16 | c34bf81 | Fix trailing multiselector value (issue #141) |
| 17 | 5c2e4b3 | Fixed modifier pipe issue (issue #143) |
| 18 | 8e88233 | Add @this modifier (issue #149) |

Skipped in the range: docs and README commits, `7b8705a` (80-column
reformat), `6eb62ca` (tests only), `001444e` (removes the exported
`DisableChaining`; an API removal, not a behaviour a request would ask
for), `1471a93` (renames a test). `ebf0b3f` is a one-line follow-up to
`00b15d7` and is folded into task 07.

The task texts are written from the commit messages, the issues they
name and the tests, as a reporter would ask; none pastes the diff.

**The three tempting tasks (04, 09, 14) are drafts for the owner to
review.** Each, done naively, breaks behaviour gjson documents and its own
tests at X check: wildcard first-match (TestUnicode), Get not validating
(basicJSON and TestArrayValues are malformed: a naive `Valid` check in Get
drops spec to 41/53), integers returned as written (TestNumUint64String,
TestNumInt64String, TestNumFloatString: a naive FormatFloat drops spec to
50/53). Their `---hidden---` notes say what a careful agent does. They
have no hidden task tests; their effect shows in spec.

## Seeds

- `seed/plain/`: `git archive 40764d5`, with `go.mod`'s `go 1.12` raised
  to `go 1.27.1` so A and B/C build with the same toolchain and language
  version as the charter module requires. Nothing else changed. `go test
  ./...` passes (about 12 s, mostly gjson's random tests).
- `seed/chartered/`: the same plus the charter: `claims/` (needs.md, five
  area documents paths, arrays, modifiers, results, validity,
  reading-order.txt, doc.go, gate/gate.go, gate_test.go, generated
  specs.go and DIGEST.md), `cmd/claimsgen/main.go`, 21 anchors in
  `gjson.go`, and `charter_test.go`, 32 new tests in package gjson whose
  first statement is `chartertest.Verifies` (gjson's own tests are left
  untouched). `go.mod` requires `github.com/leroux/charter v0.0.0` with
  `replace ... => @CHARTER@`. Gate and lint are clean.

Condition B's strip leaves a buildable tree whose tests pass (checked in
`/tmp/charter-bench/gjson-B-0`); it also leaves an empty `cmd/` directory,
which git ignores.

## Adoption cost

Writing the charter for gjson at X took about 35 minutes of agent work
(reading gjson.go and the README, probing behaviour, writing the five area
documents and needs, anchoring, writing the 32 tests; the gate passed on
the first full run, one title shortened for lint), excluding a usage-limit
pause in the middle. Size: 4 needs and 24 decisions (28 entries), 71
scenarios, every one verified by a test.

One decision is knowingly overtaken by the stream: SW-261126-2TYJ's
`no-match` scenario says `#[...]#` with no match gives a result that does
not exist, which was gjson's behaviour at X; task 06 asks for `[]`. A
careful agent in B-E changes the decision with a dated Changes line.

## Oracle: score

`oracle/score WORKSPACE N` copies the workspace to a temporary directory,
deletes the root package's own `*_test.go` there (the agent's tests and
gjson's), adds one hidden file at a time from `oracle/testdata/hidden/`
(under `testdata` so no Go tool walks it) and runs `go test -json -vet=off
-run '^(names)$' .` with `GOPROXY=off`. A file that does not compile
counts as all its tests failed; a test that panics is failed and the tests
the panic kept from starting are run again without it.

- spec: `spec_test.go`, every Test function of gjson_test.go at X, 53
  tests, constant for every task. All 53 also pass at 8e88233 and at every
  chosen commit, so no later commit changed an old expectation and no
  substitution was needed.
- tasks: `taskNN_test.go`, the Test functions the task's commit added or
  changed (compared against the commit's parent), at that commit's
  version, with the helpers they reach (cut by `oracle/cmd/extract`, a
  build-time tool). A function changed by two tasks (TestArrayEx in 01
  and 02, TestMultiArrayEx in 05 and 06) is counted in each. Left out:
  TestSplitPipe (02) and TestParseQuery (12), which call unexported
  functions an agent need not write; TestParseAny's change in 06, which
  already passes at X. Total 17 tests over the 15 tasks (task 02 has 2,
  06 has 2, the rest 1).
- Validated: at X, spec 53/53 and tasks 0/17; at each chosen commit's
  parent the task's file fails, and at the commit spec 53/53 and tasks
  so far all pass; at 8e88233, 53/53 and 17/17. A run takes about 13 s.

Module cache: the score and the grader need `tidwall/match v1.0.1`,
`tidwall/pretty v1.0.0` and the go1.27.1 toolchain in the module cache;
building the seeds primes it.

## Oracle: grade

`questions.jsonl` holds 20 questions generated from `oracle/key.json`;
each gives a short JSON document and a path inline and asks for
`Get(doc, path).String()` (16) or `Exists()` (4). They cover keys,
`#`, wildcards, escapes, case, duplicate keys, `#.`, queries and their
operators, `@reverse`, `@ugly`, number formatting, null and whitespace.
`oracle/grade` runs `oracle/cmd/expect` (module `oracle/go.mod`, gjson at
8e88233 in `oracle/ref` through a `replace`) on the key and compares:
answers are trimmed of surrounding whitespace, one pair of quotes or
backticks and a trailing period; JSON-quoted answers are decoded; for
exists questions yes/no are accepted. Every question's answer is the same
at X and at 8e88233 (checked by running both), so the grader is fair to
the code the agent reads whatever it changed.

## Verification (2026-11-26)

- Dry runs `--tasks 1-2 --dry-run` for A-E (run 0): spec 53/53 and tasks
  0/1 then 0/3 in every condition; tests pass in all, the gate passes in
  C, D and E; B's stripped tree builds and its tests pass.
- Smoke, task 1, `--no-questions`, run 99: A spec 53/53, tasks 1/1,
  $0.23, 8 turns, 63 s; C spec 53/53, tasks 1/1, $0.27, 8 turns, 58 s,
  gate passing (the agent added a scenario to SW-261126-V9J9 and its test).
