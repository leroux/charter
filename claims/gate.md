# The gate

Code: chartertest

What the gate holds the tree to once the entries, the anchors and the
citations are read.

## CH-030 An active decision is anchored and each scenario cited; a waiver absorbs a missing test, never a missing anchor

Status: active
Kind: invariant
Serves: CH-001
See: CH-022

An active decision with no realizes anchor fails, and nothing waives it;
that a concerns anchor is not one is CH-022's. An active entry with
scenarios fails for each scenario no citation names, naming the document,
the line and the call to write; one without scenarios fails when no
citation names it at all. The waiver file lists, one per line with a
reason after a tab, an id or a scenario as ID/handle, never a glob; a
matching line absorbs the missing test; a line that absorbs nothing, gives
no reason or is a glob fails.

Why: coverage is the sensor; a missing anchor is code nobody can find,
and no reason makes that acceptable.

Scenarios:
1. `not-realized`: an active decision with no realizes anchor fails naming
   the directive to add.
2. `scenario-unverified`: an active entry's scenario that no test cites
   fails naming the scenario and the citation.
3. `spec-unverified`: an active entry without scenarios that no test cites
   fails.
4. `waiver-absorbs-a-scenario`: a waiver line naming the scenario makes
   the missing test no problem.
5. `waiver-never-absorbs-an-anchor`: with the entry's tests waived, a
   missing anchor still fails.
6. `idle-waiver-fails`: a waiver that absorbs nothing fails naming its
   line, and so does one without a reason.
7. `glob-refused` (owner): a waiver written as a glob fails
   naming its line; what a waiver accepts is on its line.

Changes:
- 2026-10-09: a glob waiver, which absorbed every scenario of an entry at
  once, is refused.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-10-09: that a concerns anchor realizes nothing moved to CH-022,
  the concerns directive being its subject: not-realized was tested with
  a decision anchored by concerns alone, as CH-022's
  concerns-is-not-realized is, and is now tested with one anchored by
  nothing; the text names CH-022 for concerns.

## CH-031 A need needs no anchor and is verified only from the stakeholder's side

Status: active
Serves: CH-001

A need is realized by the decisions that serve it, so the gate asks no
anchor of it. When the options name the test files on the stakeholder's
side, a citation of a need from any other file fails naming the file, and
does not count; when they name none, any test file counts.

Why: a need is a promise to a person; a unit test of a package proves the
package, not the promise.

Scenarios:
1. `need-without-anchor-passes`: a need with no anchor and a citation is
   covered.
2. `offside-citation-refused`: a citation of a need from a file the
   options do not name fails and leaves the need unverified.
3. `onside-citation-counts`: a citation from a named file verifies the
   need.

## CH-032 A probe cites observations only, and a test never cites an observation

Status: active
Serves: CH-005

A probe is a test file whose build constraint names the probe tag on a
positive path. A probe referencing an entry fails, and its references are
no citations; a test file that is not a probe referencing an observation
fails. Without a probe tag in the options, observations are one
configuration problem and no probe rule runs.

Why: a probe reaches a live system and a test does not; a fact and a
promise must not be checked by the same run.

Scenarios:
1. `probe-reference-is-no-citation`: an entry referenced only from a probe
   is unverified, and the probe is reported.
2. `test-citing-an-observation-fails`: a test file referencing an
   observation fails naming the file.
3. `negated-tag-is-no-probe`: a file built without the probe tag is an
   ordinary test.
4. `tag-required`: observations without a configured probe tag are one
   problem, and nothing is blamed on a probe.

## CH-033 Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces

Status: active
Kind: invariant
Serves: CH-004
Accounts: env:CHARTER_UPDATE_BASELINE

The surfaces are the exported string constants of the env packages, every
build tag but the probe tag and version tags, every operation of the
OpenAPI files, and every kind of the events file; an adopter may add
kinds. A surface no entry accounts for and the baseline does not list
fails; an entry accounting for a surface nothing enumerates fails; a
baseline line whose surface is gone, or which an entry now accounts for,
or which appears twice, or which lacks a reason, fails. An enumerator
configured by a path that finds nothing is a misconfiguration; zero build
tags is a valid tree. With the update variable set the gate rewrites the
baseline with every unaccounted surface, keeping the reasons it had. A
deprecated entry's Accounts own nothing: what it accounted for is a live
entry's to claim, or the baseline's.

Why: the baseline is a debt; a gate that lets it grow silently, or hold
what is paid off, is not keeping the books.

Scenarios:
1. `surfaces-enumerated`: the fixture's env constants, tags, operations
   and event kinds are the surfaces, the probe tag and version tags left
   out.
2. `unaccounted-fails`: a new build tag fails naming it.
3. `accounts-must-exist`: an entry accounting for a surface nothing
   enumerates fails.
4. `baseline-holds-live-gaps-once`: a baseline line for a surface that is
   gone, or now accounted for, or listed twice, or without a reason,
   fails.
5. `empty-enumerator`: an events file with no kinds is a configuration
   problem, and a tree with no build tags is not.
6. `update-rebuilds`: the update variable rewrites the baseline with every
   unaccounted surface, keeps the reasons it had, and a second run changes
   nothing.
7. `deprecated-accounts-nothing` (owner): a surface accounted
   for only by a deprecated entry is unaccounted.

Changes:
- 2026-10-09: a deprecated entry's Accounts counted; they own nothing now.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-034 An environment variable is read through a named constant, or under a dynamic directive, and no literal names one

Status: active
Kind: invariant
Serves: CH-004

Outside the env packages and outside test files, a call of a configured
reader takes a constant of an env package as its first argument, or the
line carries the dynamic directive. A string literal equal to a declared
variable's value, or matching a configured pattern, fails.

Why: a variable read by a literal is a surface the enumerators cannot
see.

Scenarios:
1. `literal-argument`: a reader called with a string literal fails naming
   the constant to use.
2. `variable-argument`: a reader called with a variable fails.
3. `duplicated-name`: a literal equal to a declared variable's value fails
   naming the constant.
4. `undeclared-pattern`: a literal matching a configured pattern fails.
5. `dynamic-directive`: a read whose line carries the directive is exempt.

## CH-035 Observations carry their status's fields, observed ones a date and a probe; one broken under an active entry fails

Status: active
Serves: CH-005

Every observation names a subject. An observed one has a calendar date and
at least one citing probe; a documented one has a source; any date given
is a calendar date. A broken observation supporting an active entry fails;
under a planned or a deprecated entry it does not.

Why: a fact is only worth recording with where it came from and when it
was last true.

Scenarios:
1. `observed-needs-date-and-probe`: an observed observation without a
   date, with a date that is not on the calendar, or cited by no probe,
   fails.
2. `documented-needs-source`: a documented observation without a source
   fails, and one with a stray date fails.
3. `subject-required`: an observation without a subject fails.
4. `broken-under-active`: a broken observation under an active entry
   fails, and under a planned or deprecated one does not.

Changes:
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-036 The walk covers every Go and TypeScript file under the root within the scope, skipping what the options ignore

Status: active
Serves: CH-006

The walk parses every Go file under the root and, when TypeScript is
configured, lists every .ts and .tsx file that is not a declaration file.
It skips .git, node_modules, vendor and testdata, hidden and underscore
directories, nested modules, every Ignore pattern matched against the
path and the base name, and everything outside Only when Only names
paths.

Why: two charters in one module must not see each other's anchors; the
scope is what keeps them apart.

Scenarios:
1. `hidden-and-nested-skipped`: a hidden directory and a nested module are
   not walked.
2. `ignore-patterns`: a pattern matches a path or a base name.
3. `only-restricts`: with Only naming paths, a directory outside them is
   not walked, so its anchors and citations are unseen.
4. `typescript-listed`: with TypeScript configured, the tree lists the
   TypeScript files and no declaration file.

## CH-037 The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree

Status: active
Serves: CH-001, CH-002

Check reports each problem as one test error; Problems returns the same
list, sorted. The options resolve the module root from the calling file or
the working directory, require the entries package inside the module, and
validate the TypeScript settings.

Why: a gate whose output a person cannot act on is noise.

Scenarios:
1. `clean-fixture-passes`: the clean fixture reports nothing.
2. `root-resolved`: the root is the module above the calling file, else
   above the working directory; a package outside the module is refused.
3. `typescript-options-validated`: TypeScript needs a module, an import
   and a compiler, and resolves a relative compiler under the root.

Changes:
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
