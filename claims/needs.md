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

## CH-261126-PFG5 What the software must do is written once, above the code, where people and agents both read and change it

Status: planned

The needs and decisions of a repository are written in its claims
directory, each once, in the form every charter shares, and they are the
reference a change is held to. Each keeps its id for life, and what changed
in it, and when, is written beside it.

Why: intent that lives in heads, chats and pull requests is lost to the
next person and to every agent; code says what the software does, not what
it should, so without a written reference there is nothing to hold it to.

Scenarios:
1. `one-form-everywhere` (agent, provisional): the same need and decision,
   written in two charters of different languages and layouts, are read
   alike by the command.
2. `reference-survives` (agent, provisional): an entry changed with a dated
   line keeps its id, and the record shows the change and when it was made.

## CH-261126-P2SK The code does what the record says, and a gap between them stops the change

Status: planned

Every active decision has code that claims it and a test that shows each of
its scenarios, every surface the tree exposes is owned by an entry, and the
gate fails a tree where any of that is missing. A change, a person's or an
agent's, that leaves such a gap does not pass.

Why: code written fast, by agents most of all, drifts from what was decided
without anyone noticing; a check that fails on the gap is the only
adherence that does not depend on someone remembering to look.

Scenarios:
1. `gap-stops-the-change` (agent, provisional): a tree whose active decision
   lost its anchor, or whose scenario lost its citation, fails the check
   naming the gap.
2. `unowned-surface-stops-the-change` (agent, provisional): a new
   environment variable no entry owns fails the check naming it.

## CH-001 A reviewer can judge a change by its entries and their tests, not by reading the code

Status: active

Every clause of an active decision is a scenario, every scenario has a
test that names it, and the gate fails on the first one that does not, so
the reviewer reads the area document diff and the gate's verdict and
opens the code only where the review points. The gate is static and runs
in seconds, and its verdict is the same on every machine.

Why: review is the scarce resource. A gate that proves the links exist and
names the missing one lets the reviewer spend attention on whether the
decision is right and the test means what the scenario says, which no gate
can do.

Scenarios:
1. `clean-tree-passes`: a tree whose every active entry is realized and
   whose every scenario is cited passes with no problem.
2. `missing-scenario-named`: an active decision with a scenario no test
   cites fails, naming the document, the line, the scenario and the
   citation to write.
3. `every-rule-has-a-mutation`: every rule of the gate has a mutation of
   the clean fixture that trips it, named after the rule.

## CH-261009-3FHW A reader tells what was ruled from what an agent chose, and an entry's change is recorded or named

Status: active

Every scenario says who ruled it: the owner, with what the ruling cites;
a named delegate; or an agent whose choice nobody has ruled on, which a
scenario with nothing written is. The registry and the report carry it,
and the report counts the active scenarios by who ruled them. A change to
the record that its Changes lines do not say, an entry removed, unbuilt
or retired without a line, text or a scenario changed without one, or a
ruling lost, is named against the revision it was made from.

Why: an agent's default and the owner's decision look the same on the
page, so the record could not say what was decided; a change nobody wrote
down is the drift the whole form exists to refuse.

Scenarios:
1. `provenance-in-the-registry` (owner): reading the registry
   gives each scenario's provenance, and a scenario with none written is
   an agent's, provisional.
2. `report-counts-who-ruled` (owner): the report shows each
   scenario's provenance and counts the active scenarios by it.
3. `unrecorded-change-named` (owner): a change to an entry's
   text with no dated Changes line is named against the earlier revision,
   and one with the line is not.

Changes:
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-002 A builder reaches the code and the tests from an id, and the gate names the fix for every problem

Status: active

From an entry's id the tooling lists the code that realizes it, by file,
line and declaration, and the tests that verify it, by file, test and
scenario. Every problem the gate reports names a file, a line and what to
do. Symbols are generated in each language the charter reads, so a citation
of an entry that does not exist fails to compile or to import.

Why: an agent building to a decision needs the way in and the way to know
it is done, without reading the whole tree or guessing what a failure
means.

Scenarios:
1. `anchors-and-citations-listed`: reading the registry gives, per entry,
   its anchors with their declarations and its citations with their
   handles.
2. `problems-name-file-line-and-fix`: a form problem in an area document
   names the document, the line and the form to write.
3. `planned-entry-cited-first`: a test may cite a planned entry, so a test
   can be written before the code, and the gate asks nothing else of it.

Changes:
- 2026-11-26: symbols are generated in each language the charter reads,
  Go, Python and TypeScript, so a citation of an entry that does not exist
  fails to compile or to import, where it named Go's compiler alone
  (backlog 003).

## CH-261126-NS7T An adopter keeps a charter in any repository layout and language, and its code imports nothing of the tooling

Status: active

A repository takes the charter up by writing its entries and a
charter.toml in its own claims directory and running the charter command
over it. Its code imports nothing of the tooling: its tests cite through
the files the command generates. The layout is the repository's own: Go
modules nested anywhere, Python and TypeScript beside them, and no go.mod
at the root. The examples in this repository are kept so, each checked
with the tooling's tests.

Why: a method one repository can use only by copying its tooling is a
fork, not a method; and an adopter in the same tree is what tells a change
to the tooling that it broke the people who use it, before they find out.

Scenarios:
1. `example-gate-passes` (owner, backlog 003): the ratelimit example
   passes charter check over its own charter, and charter lint finds
   nothing it fails on.
2. `example-untested-scenario-named` (agent, provisional): a copy of the
   example whose test of one scenario no longer cites it fails its gate,
   naming that scenario.
3. `polyglot-example-passes` (owner, backlog 003): an example of Go in a
   nested module and Python beside it, with no go.mod at its root, passes
   charter check over its own charter.
4. `polyglot-untested-scenario-named` (owner, backlog 003): a copy of that
   example whose Python test no longer cites a scenario fails charter
   check, naming the scenario.
5. `nothing-of-the-tooling-imported` (owner, backlog 003): no file of
   either example imports a package of the tooling, and no go.mod of
   theirs requires it.

Changes:
- 2026-11-26: an adopter keeps claims/charter.toml and runs the charter
  command, where it required this module and wrote a gate test and a
  claimsgen of its own, and its code imports nothing of the tooling, the
  generated files carrying what its tests call. Any layout and the three
  languages are kept, a second example of Go and Python proving it;
  `polyglot-example-passes`, `polyglot-untested-scenario-named` and
  `nothing-of-the-tooling-imported` are added (backlog 003).
