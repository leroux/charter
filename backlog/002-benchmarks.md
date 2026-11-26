# 002. Benchmarks: measure whether a charter keeps agent-written code in spec and makes it easier to understand

Status: in progress · Size: L · See: 001

The framework's two goals are to keep vibe-coded codebases in spec and to
make them understandable. This story measures both with numbers an oracle
computes, so the results decide what to keep, trim or make an extension.
The contract and the conditions are in `bench/README.md`.

## Done

- **Harness** (`bench/harness/`): `run.py` runs one condition and run,
  driving a `claude -p` session resumed task after task in a workspace
  under `/tmp/charter-bench/`, scoring each task with the benchmark's
  hidden oracle, then asking the questions in a fresh session. `matrix.py`
  runs many in parallel and prints a table.
- **Five conditions**, each adding one thing: A no charter, B the charter's
  markdown, C the gate, D report, diff and the review skills, E mutate.
- **Trail kept**: per-task JSON, both session transcripts and the
  workspace's history as a git bundle, committed under `bench/results/`.
- **Invalid runs**: a task whose agent hit a usage limit or an API error
  is marked invalid and stops the run.
- **Condition B** gets the hand-written documents only, not the generated
  digest.
- **coreutils**: Go clone of cut, sort, uniq, wc against GNU coreutils
  9.4. 275 spec cases, 180 task cases, 14 tasks (4 tempting), 20
  questions. Charter: 3 needs, 12 decisions, 43 scenarios.
- **toml**: TOML decoder against toml-test (submodule). 298 spec checks,
  task cases per feature over 14 tasks (4 tempting), 20 questions.
  Charter: 2 needs, 18 decisions, 54 scenarios.
- **gjson**: gjson at 40764d5 and its next 15 upstream commits as tasks,
  plus 3 tempting (submodule). 53 spec tests, 17 task tests, 20
  questions. Charter written for existing code in about 35 minutes: 4
  needs, 24 decisions, 71 scenarios.
- **Smoke runs**, task 1 of each benchmark, conditions A and C: every run
  scored full marks; C cost 15–60% more per task ($0.27–0.94 against
  $0.23–0.59). Total spend building and testing: about $4.

## To do, in order

1. **Done: the oracles are locked away.** The harness denies the agent's file tools the bench/ directory through --settings and audits every transcript for the hidden paths, marking a run that reached for them invalid. Running the agents as another user was refused, since it needs the session's credential. Was: The hidden tests are in the repository and
   the agents run with permission checks off. Copy each benchmark's
   `oracle/` and the hidden parts of its tasks to a directory the agents'
   user cannot read, and have the harness score from there.
2. **Pilot on coreutils**: A, B and C, one full run each, with questions.
   Estimated $25–50. It answers the first real question (does B or C
   drift less than A?) and gives the real cost of a full run.
3. **Read the pilot's transcripts at the tempting tasks**: did the agent
   notice the conflict, push back, update the charter, or break the spec?
   This is where the conditions should separate.
4. **Owner reviews the tempting tasks**, informed by the pilot (see
   Decisions).
5. **Pin gjson's random tests' seed** in its oracle, so drift scores are
   not noisy.
6. **Full matrix**: three benchmarks × A–E × 3 runs, once the pilot shows
   the conditions separate. Budget from the pilot's real cost.
7. **Condition F**, charting as you go: a seed with `needs.md` only and a
   CLAUDE.md that asks the agent to chart the decisions each change
   touches. Tests whether gradual adoption gives most of the value; added
   if the matrix shows C helps.
8. **Results page**: the tables, the cost per point of drift avoided, and
   what they mean for each feature: keep, trim, or move to an extension.
9. **Decide the core and the extensions** from the results (the critique's
   proposal: a core of record, symbols, links and gate; report, diff,
   digest, style, questions, observations, surfaces, TypeScript and mutate
   as removable extensions).

## Decisions for the owner

| Decision | Default until ruled |
|---|---|
| The tempting tasks: coreutils 04, 07, 10, 13; toml 03, 08, 10, 12; gjson 04, 09, 14 | the drafts as written |
| Budget for the pilot, then the matrix | the pilot waits for a go |
| coreutils questions graded on exact spacing or whitespace-insensitive | exact; the loose score is recorded too |
| toml: a question on a feature the decoder never built counts only "invalid" as correct | yes |
| gjson: SW-261126-2TYJ is left stale against task 06 on purpose, to test whether the agent updates the charter | keep |
| gjson: pin the random tests' seed | pin |
| Whether condition F is added | after the matrix |

## Known limits

- An agent could read the oracle until item 1 is done.
- A smoke run is one task; it shows the loop works, not that a condition
  helps.
- Each benchmark's charter was written by an agent, the tempting tasks
  too; both are part of what the owner reviews.

## Status at handoff (2026-11-26)

- Pilot coreutils A run 1 finished: spec 275/275 after every task
  including all four tempting ones, tasks 170/180 at the end (the misses
  start at task 12, missing-file error messages), questions
  20/20, about $6.58 (the costs recorded are the session's running
  total). Results under `bench/results/coreutils/A/1/`.
- Reading so far: the no-charter agent did not drift on coreutils, so
  coreutils may not separate the conditions (GNU is a well-known spec).
  Read the transcripts at tasks 04, 07, 10 and 13 before running C and B
  on it; toml or gjson, or harder tempting tasks, may be the better next
  pilot.
- Proposed next, not started: questions answered from the charter alone
  (no code) beside code-only and both, to measure the charter as an
  understanding layer; a rendered overview and entry pages from
  claimsmap's JSON; a charter diff on each pull request.
- The owner intends to build a project from scratch with the charter in
  a new session, as product manager with the agent as engineer.
