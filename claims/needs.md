# What people need of the charter

The charter is a higher level to build software at: what the software must
do is written down above the code, the code is held to it, and a change is
reviewed through it, by people and agents alike. Each need below answers
one problem that level exists to solve, in this order: intent has no
durable home; code drifts from intent unnoticed; reviewing code line by
line does not scale; who decided what is invisible; a newcomer, a person or
an agent, cannot find their way in; and every project is different.

The people it serves are the owner of the record, the reviewer of a
change, the builder, a person or an agent building to a decision, and the
adopter who keeps a charter in a repository of their own. Each need is
verified the way they meet the charter: through its command and its
public entry points, over a fixture.

The key words MUST, MUST NOT, SHOULD, SHOULD NOT and MAY in these
documents are to be interpreted as described in BCP 14 (RFC 2119, RFC 8174)
when, and only when, they appear in all capitals.

## CH-261126-PFG5 A repository's charter MUST write what the software must do once, above the code, where people and agents both read and change it

Status: active

- A repository's needs and decisions MUST be written in its claims
  directory, each once, in the form every charter shares, and they are
  the reference a change is held to.
- Each entry MUST keep its id for life, and what changed in it, and
  when, MUST be written beside it.

Why: intent that lives in people's heads, in chats and in pull requests
is lost to the next person and to every agent. The code says what the
software does, not what it should do, so without a written reference
there is nothing to hold the code to; writing each need and decision
once, in one shared form, gives every reader the same reference.

Scenarios:
1. `one-form-everywhere` (agent, provisional): given the same need and
   decision written in two charters of different languages and layouts,
   when the command reads both, then it reads them alike.
2. `reference-survives` (agent, provisional): given an entry changed with
   a dated line, when the record is read, then the entry keeps its id and
   shows the change and when it was made.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-P2SK The charter gate MUST fail a change that leaves a gap between the record and the code

Status: active

- Every active decision MUST have code that claims it and a test that
  shows each of its scenarios.
- Every surface the tree exposes MUST be owned by an entry.
- The gate MUST fail a tree where any of that is missing, so a change,
  a person's or an agent's, that leaves such a gap does not pass.

Why: code written fast, and by agents most of all, drifts from what was
decided without anyone noticing. A check that fails on the gap is the
only kind of adherence that does not depend on someone remembering to
look.

Scenarios:
1. `gap-stops-the-change` (agent, provisional): given a tree whose
   active decision lost its anchor, or whose scenario lost its citation,
   when the check runs, then it fails naming the gap.
2. `unowned-surface-stops-the-change` (agent, provisional): given a new
   environment variable no entry owns, when the check runs, then it fails
   naming it.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-001 The charter gate MUST let a reviewer judge a change by its entries and their tests, without reading the code

Status: active

- Every clause of an active decision MUST be a scenario, and every
  scenario MUST have a test that names it.
- The gate MUST fail on the first one that does not, so the reviewer
  reads the area document diff and the gate's verdict, and opens the
  code only where the review points.
- The gate MUST be static, run in seconds, and give the same verdict on
  every machine.

Why: review is the scarce resource. A gate that proves the links exist
and names the missing one lets the reviewer spend attention on what no
gate can check: whether the decision is right and whether the test means
what the scenario says.

Scenarios:
1. `clean-tree-passes`: given a tree whose every active entry is realized
   and whose every scenario is cited, when the gate runs, then it passes
   with no problem.
2. `missing-scenario-named`: given an active decision with a scenario no
   test cites, when the gate runs, then it fails naming the document, the
   line, the scenario and the citation to write.
3. `every-rule-has-a-mutation`: given the clean fixture, when each rule
   of the gate is considered, then a mutation of the fixture named after
   the rule trips it.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261009-3FHW The record MUST show who ruled each scenario, and charter diff MUST name any change to an entry the record does not state

Status: active

- Every scenario MUST say who ruled it: the owner, with what the ruling
  cites; a named delegate; or an agent whose choice nobody has ruled on,
  which is what a scenario with no provenance written is.
- The registry and the report MUST carry the provenance, and the report
  MUST count the active scenarios by who ruled them.
- A change to the record that its Changes lines do not state (an entry
  removed, unbuilt or retired without a line, its text or a scenario
  changed without one, or a ruling lost) MUST be named against the
  revision it was made from.

Why: an agent's default and the owner's decision look the same on the
page, so without provenance the record could not say what was actually
decided. A change that nobody wrote down is exactly the drift the whole
form exists to refuse.

Scenarios:
1. `provenance-in-the-registry` (owner): given the registry, when it is
   read, then each scenario has its provenance, and a scenario with none
   written is an agent's, provisional.
2. `report-counts-who-ruled` (owner): given a record, when the report is
   made, then it shows each scenario's provenance and counts the active
   scenarios by it.
3. `unrecorded-change-named` (owner): given an entry's text changed with
   no dated Changes line, when diff runs against the earlier revision,
   then it names the entry; given the line, then it does not.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-002 The charter command MUST lead a builder from an id to its code and tests, and every gate problem MUST name its fix

Status: active

- From an entry's id, the tooling MUST list the code that realizes it, by
  file, line and declaration, and the tests that verify it, by file, test
  and scenario.
- Every problem the gate reports MUST name a file, a line and what to do.
- Symbols MUST be generated in each language the charter reads, so a
  citation of an entry that does not exist fails to compile or to import.

Why: an agent building to a decision needs a way into the code and a way
to know when it is done, without reading the whole tree or guessing what
a failure means.

Scenarios:
1. `anchors-and-citations-listed`: given a record with anchors and
   citations, when the registry is read, then each entry lists its
   anchors with their declarations and its citations with their handles.
2. `problems-name-file-line-and-fix`: given a form problem in an area
   document, when the gate runs, then it names the document, the line and
   the form to write.
3. `planned-entry-cited-first`: given a test citing a planned entry, when
   the gate runs, then it accepts the citation and asks nothing else of
   the entry, so a test can be written before the code.

Changes:
- 2026-11-26: symbols are generated in each language the charter reads,
  Go, Python and TypeScript, so a citation of an entry that does not exist
  fails to compile or to import, where it named Go's compiler alone
  (backlog 003).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-NS7T An adopter MUST be able to keep a charter in any repository layout and language, with code that imports nothing of the tooling

Status: active

- A repository MUST be able to take the charter up by writing its
  entries and a charter.toml in its own claims directory and running the
  charter command over it.
- Its code MUST import nothing of the tooling: its tests cite through the
  files the command generates.
- The layout is the repository's own: Go modules nested anywhere, Python
  and TypeScript beside them, and no go.mod at the root.
- The examples in this repository MUST be kept so, each checked with the
  tooling's tests.

Why: a method that a repository can use only by copying its tooling is a
fork, not a method. Keeping an adopter inside the same tree also means a
change to the tooling learns that it broke the people who use it before
they find out.

Scenarios:
1. `example-gate-passes` (owner, backlog 003): given the ratelimit
   example, when charter check runs over its own charter, then it passes,
   and charter lint finds nothing it fails on.
2. `example-untested-scenario-named` (agent, provisional): given a copy of
   the example whose test of one scenario no longer cites it, when its
   gate runs, then it fails naming that scenario.
3. `polyglot-example-passes` (owner, backlog 003): given an example of Go
   in a nested module and Python beside it, with no go.mod at its root,
   when charter check runs over its own charter, then it passes.
4. `polyglot-untested-scenario-named` (owner, backlog 003): given a copy
   of that example whose Python test no longer cites a scenario, when
   charter check runs, then it fails naming the scenario.
5. `nothing-of-the-tooling-imported` (owner, backlog 003): given either
   example, when its files and go.mod files are read, then none imports a
   package of the tooling or requires it.

Changes:
- 2026-11-26: an adopter keeps claims/charter.toml and runs the charter
  command, where it required this module and wrote a gate test and a
  claimsgen of its own, and its code imports nothing of the tooling, the
  generated files carrying what its tests call. Any layout and the three
  languages are kept, a second example of Go and Python proving it;
  `polyglot-example-passes`, `polyglot-untested-scenario-named` and
  `nothing-of-the-tooling-imported` are added (backlog 003).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.
