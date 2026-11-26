# Benchmarks

Does a charter keep an agent-written codebase in spec, and does it make
the codebase easier to understand? Each benchmark here answers that with
numbers an oracle computes, never a judgment.

| Benchmark | Start | Oracle | What it tests |
|---|---|---|---|
| `coreutils` | greenfield clone of `cut`, `sort`, `uniq`, `wc` | GNU coreutils, differential | drift as flags pile up |
| `toml` | greenfield TOML parser | the official toml-test suite | staying in a published spec |
| `gjson` | an existing library at an old commit | its own later commits and their tests | adopting charter on real code |

## Conditions

Each condition adds one thing to the one before, so the difference
between two adjacent conditions is the value of that thing.

| | The agent gets | Isolates |
|---|---|---|
| A | the code, a plain CLAUDE.md | baseline |
| B | A plus the charter's markdown (`claims/*.md`), and CLAUDE.md asks it to keep the charter true | a written spec alone |
| C | B plus the tooling: generated symbols, anchors, `Verifies` citations and the gate as a test | the gate |
| D | C plus `claimsgen report` and `diff`, and the claims-implement and claims-review skills | the extra checks and the loop |
| E | D plus `claimsgen mutate` | mutation testing |

## A benchmark's contract

`bench/<name>/` holds:

- `bench.json`: `{"name", "description", "language": "go", "conditions": ["A",…]}`.
- `seed/plain/`: the project as condition A starts it, a directory the
  harness copies. No charter in it.
- `seed/chartered/`: the same project with its charter applied in full
  (condition C's tree): `claims/` with needs and decisions, `claims/gate/`,
  the gate test, generated `specs.go`, anchors in code, `Verifies` in
  tests, and its own `cmd/claimsgen`. Its `go.mod` requires
  `github.com/leroux/charter` with a `replace` to the path the harness
  substitutes for `@CHARTER@`.
- `tasks/NN-slug.md`: the change requests, in order, each the text a
  person would type. Front matter line 1: `kind: neutral` or
  `kind: tempting` (a request that, done naively, breaks a decision the
  charter records). A tempting task's file also says, below a line
  `---hidden---`, what a careful agent would do; the harness strips
  everything from that line on before the agent sees it.
- `oracle/score`: an executable, run by the harness outside the agent's
  workspace as `oracle/score WORKSPACE TASK_NUMBER`. It prints one JSON
  object: `{"spec": {"passed": n, "total": m}, "tasks": {"passed": n,
  "total": m}, "details": [...]}`. `spec` is the hidden checks of the
  behaviour that must hold after every task (drift); `tasks` is the
  hidden checks of the features asked for so far (correctness). It must
  not need network access and should finish in under two minutes.
- `questions.jsonl`: understandability questions, one JSON object a
  line: `{"id", "question"}`. Each asks the agent to predict what the
  code does for a given input, answered in one line.
- `oracle/grade`: run as `oracle/grade WORKSPACE ANSWERS.jsonl`, where
  each answer line is `{"id", "answer"}`. Prints
  `{"correct": n, "total": m, "details": [...]}`, graded against the
  oracle (the reference implementation's real output), never by judgment.

The oracle and the hidden parts of tasks never enter the agent's
workspace.

## Running

```sh
bench/harness/run.py --bench coreutils --condition C --run 1          # one run, every task
bench/harness/run.py --bench coreutils --condition A --tasks 1-2 --dry-run   # plumbing only, no agent
bench/harness/matrix.py --bench coreutils --conditions A,B,C --runs 3 -j 3
```

Each run writes `bench/results/<bench>/<condition>/<run>/` with one JSON
per task (the agent's cost, turns and time, the oracle's scores, the
gate's result, the diff's size), `qa.json` for the questions, the agents'
session transcripts (`transcript.jsonl`, `qa-transcript.jsonl`) and the
workspace's history as `workspace.bundle`, a commit per task, which
`git clone` restores. Results are committed: the trail is kept, smoke runs
included. The
workspace is `/tmp/charter-bench/<bench>-<condition>-<run>/`, outside the
repository, so the agent cannot read the oracle by walking up.

Agents run with `claude -p`, one session a run resumed task after task,
in the workspace, with `--permission-mode bypassPermissions` and
`--setting-sources project,local` so that only the workspace's own
CLAUDE.md and skills apply. The questions run in a fresh session.
