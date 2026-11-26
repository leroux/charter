# What people need of the charter

The charter tooling serves five people: the reviewer of a change, the agent
or person building to a decision, the owner of the record, the operator
whose surfaces the record owns, and the adopter who keeps a charter in a
repository of their own. Each need here is verified through the
gate's public entry points over a fixture, the way those people meet it.

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

## CH-002 A builder reaches the code and the tests from an id, and the gate names the fix for every problem

Status: active

From an entry's id the tooling lists the code that realizes it, by file,
line and declaration, and the tests that verify it, by file, test and
scenario. Every problem the gate reports names a file, a line and what to
do. Symbols are generated, so a citation of an entry that does not exist
fails to compile.

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

## CH-003 The record is stable: ids never change or get reused, retired entries keep their history, and changes are dated

Status: active

An entry keeps its id for life; a retired one stays in its document with
its status and its changes, and its symbol still builds, so the strip can
follow at its own pace. Two entries with one id, anywhere in the charter,
fail. A change to an entry is a dated line in its Changes.

Why: everything hangs on the ids: anchors, citations, waivers, the
backlog. A record that renumbers or forgets is not a record.

Scenarios:
1. `id-declared-once`: two entries with one id, in one document, in two
   documents, or in a document and a spec file, fail naming both places.
2. `retired-entry-still-builds`: a deprecated entry is generated with a
   deprecation notice and the gate asks no coverage of it.
3. `next-free-id`: a new id with a prefix is one the charter does not
   hold.

Changes:
- 2026-10-01: a new id is one the charter does not hold, where it was one
  above the highest, which two branches took at once and one of them then
  changed (CH-012).
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-004 No surface exists that an entry does not own, and the baseline only shrinks

Status: active

Every environment variable, build tag, contract operation and event kind
the tree exposes is named by an entry's Accounts, or listed with a reason
in a baseline that predates the charter. The gate fails a surface neither
names, a baseline line whose surface is gone or now accounted for, and a
read of an environment variable that goes through anything but a named
constant.

Why: the surfaces are the operator's contract with the system. A variable
nobody claims is a promise nobody keeps.

Scenarios:
1. `unaccounted-surface-fails`: a build tag no entry accounts for and the
   baseline does not list fails naming it.
2. `baseline-holds-only-live-gaps`: a baseline line for a surface that is
   gone, or that an entry now accounts for, fails naming the entry.
3. `env-read-through-a-constant`: an environment read whose argument is a
   string literal, or a variable, fails naming the constant to use.

## CH-005 An external system's fact is recorded with its source and date, and no active entry rests on one found false

Status: active

An observation records a fact about a system the code does not own, with
its subject, where it was read or when a probe last saw it true, and the
entries that rest on it. Probes check observations and nothing else; tests
cite entries and never observations. A broken observation under an active
entry fails the gate.

Why: a decision that rests on a vendor's behaviour is only as good as the
fact; recording the fact apart from the decision says what to re-check
when the vendor changes.

Scenarios:
1. `observed-fact-has-a-date-and-a-probe`: an observed observation without
   a calendar date, or cited by no probe, fails.
2. `broken-fact-under-active-fails`: a broken observation supporting an
   active entry fails, and one supporting a planned or deprecated entry
   does not.
3. `probes-and-tests-stay-apart`: a test file referencing an observation,
   and a probe referencing an entry, each fail.

Changes:
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-006 The charter specifies itself, and its own gate is green

Status: active

The tooling's decisions are entries in this directory, in the form the
charter asks of any adopter, realized by the tooling and verified by its
tests, and checked by a gate of their own over the tooling alone.

Why: a form that cannot carry the charter's own thirty decisions is the
wrong form, and a gate nobody runs on its own maker is a claim, not a
proof.

Scenarios:
1. `self-gate-walks-the-tooling-only`: the charter's own gate walks the
   tooling and nothing else, so an adopter's anchors and tests, the
   example adopter's among them, are invisible to it.

Changes:
- 2026-11-26: the charter became a repository of its own (Q-261015-3Q84). The entries are in the form the charter asks of
  any adopter, where they were in the product's, and the gate walks the
  repository's tooling, where it walked the tooling's directories inside
  the product; the example adopter is a module of its own, which every
  walk skips.

## CH-261126-NS7T An adopter keeps a charter in their own repository with the tooling as a dependency, nothing of it copied

Status: active

A repository takes the charter up by requiring this module: it writes its
entries in its own claims directory, its gate as a test that passes its
options to the gate, and a claimsgen of its own whose main passes the same
options to the command line. Nothing of the tooling is copied into it,
and the example adopter in this repository is kept so, its gate run with
the tooling's tests.

Why: a method one repository can use only by copying its tooling is a
fork, not a method; and an adopter in the same tree is what tells a change
to the tooling that it broke the people who use it, before they find out.

Scenarios:
1. `example-gate-passes` (agent, provisional): the example adopter, with
   the tooling as a dependency, passes its own gate, and its own
   claimsgen lints its charter clean.
2. `example-untested-scenario-named` (agent, provisional): a copy of the
   example whose test of one scenario no longer cites it fails its gate,
   naming that scenario.

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
