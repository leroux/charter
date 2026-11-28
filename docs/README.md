# Charter's docs

Software built fast, by people and agents together, loses what it was
meant to do. The intent lives in heads, chats and pull requests; the code
drifts from it and nobody notices; a reviewer cannot read every line an
agent writes; an agent's default and an owner's ruling look the same once
merged; and a newcomer, person or agent, has no way from a decision to the
code that carries it. Charter writes what the software must do once, as
needs and decisions in markdown under `claims/`. A static gate, `charter
check`, fails the tree when an active decision has no code anchored to it,
a scenario no test cites, or a surface no entry accounts for, and a change
is reviewed through the record instead of the code. It works in any
repository layout, over Go, Python and TypeScript, and specifies itself:
the tooling's own decisions are the `CH` entries in `claims/`.

## Reading paths

**Deciding whether to adopt it:**

1. [Why charter exists](why.md): the problems, the need that answers each,
   and what it costs.
2. [The loop](loop.md): one change from a story's ruling to the merge, and
   who does each step.
3. [Adopting charter](adopting.md): from no charter to a green gate, with
   the two examples.
4. [What charter check does](check.md), when you want the gate's rules.

**About to change its code:**

1. [Architecture](architecture.md): the packages, the facts model, why the
   gate names no language, and where to start for each change.
2. [The loop](loop.md), which every change here follows too: an entry's
   text, a scenario and a dated Changes line first.
3. The walkthrough of the part you are changing:
   [gen](gen.md), [check](check.md), [report](report.md),
   [diff](diff.md) or [docs](docs.md).

## The docs

| Doc | Kind | What it covers |
|---|---|---|
| [why.md](why.md) | design | the six problems, the needs, the mechanisms and the alternatives the record rejected |
| [loop.md](loop.md) | walkthrough | story, claims-draft, `charter next`, claims-implement, gate, claims-review, report, diff |
| [architecture.md](architecture.md) | architecture | the packages and their graph, readers and facts, the processes the command starts, the charter's own charter |
| [gen.md](gen.md) | walkthrough | the record parser, the generated `specs.go`, `claims.py` and `index.ts`, the marker, the digest, the docs' rewrite |
| [check.md](check.md) | walkthrough | the walk, the readers, each rule in order, the rule labels, the printed problems and warnings |
| [report.md](report.md) | walkthrough | `CHARTER_EVENTS`, the event line, how outcomes fold, the profile, the exit status |
| [diff.md](diff.md) | walkthrough | the record at a git revision, each change diff names, the docs to read again |
| [docs.md](docs.md) | walkthrough | the reference syntax, `rule 26`, curation, hashes and staleness |
| [adopting.md](adopting.md) | guide | charter.toml, the first entries, wiring the claims module, anchors, citations, CI |
| [tracey.md](tracey.md) | design | charter beside tracey, the nearest tool: the requirement, links from code and tests, what fails, change over time, interfaces, and what charter could take from it |

## The needs

Each need is one problem the method answers; [why.md](why.md) takes them
in this order, and the area documents in
[claims/reading-order.txt](../claims/reading-order.txt) hold the
decisions that serve them.

- [CH-261126-PFG5: A repository's charter MUST write what the software must do once, above the code, where people and agents both read and change it](../claims/needs.md#ch-261126-pfg5-a-repositorys-charter-must-write-what-the-software-must-do-once-above-the-code-where-people-and-agents-both-read-and-change-it)
- [CH-261126-P2SK: The charter gate MUST fail a change that leaves a gap between the record and the code](../claims/needs.md#ch-261126-p2sk-the-charter-gate-must-fail-a-change-that-leaves-a-gap-between-the-record-and-the-code)
- [CH-001: The charter gate MUST let a reviewer judge a change by its entries and their tests, without reading the code](../claims/needs.md#ch-001-the-charter-gate-must-let-a-reviewer-judge-a-change-by-its-entries-and-their-tests-without-reading-the-code)
- [CH-261009-3FHW: The record MUST show who ruled each scenario, and charter diff MUST name any change to an entry the record does not state](../claims/needs.md#ch-261009-3fhw-the-record-must-show-who-ruled-each-scenario-and-charter-diff-must-name-any-change-to-an-entry-the-record-does-not-state)
- [CH-002: The charter command MUST lead a builder from an id to its code and tests, and every gate problem MUST name its fix](../claims/needs.md#ch-002-the-charter-command-must-lead-a-builder-from-an-id-to-its-code-and-tests-and-every-gate-problem-must-name-its-fix)
- [CH-261126-NS7T: An adopter MUST be able to keep a charter in any repository layout and language, with code that imports nothing of the tooling](../claims/needs.md#ch-261126-ns7t-an-adopter-must-be-able-to-keep-a-charter-in-any-repository-layout-and-language-with-code-that-imports-nothing-of-the-tooling)

The record itself is the reference: `claims/needs.md` first, then the
areas in reading order (intent, anchoring, gate, review, provenance,
adoption). These docs explain it and point into it; where a doc and the
record differ, the record binds.
