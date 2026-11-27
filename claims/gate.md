# Holding the code to the record

Code: internal/gate, internal/check, internal/walk, internal/surface

This area answers the problem that code drifts from intent unnoticed. It
covers what the gate fails a tree on once the record, the anchors and the
citations are read: coverage, the surfaces, the environment, the facts
the entries rest on, the walk, and the gate's own report.

## CH-030 The gate MUST fail an active decision without an anchor or with an uncited scenario; a waiver MAY absorb a missing test, never a missing anchor

Status: active
Kind: invariant
Serves: CH-261126-P2SK, CH-001
See: CH-022

- An active decision with no realizes anchor MUST fail, and a waiver
  MUST NOT absorb it. That a concerns anchor does not count is CH-022's.
- An active entry with scenarios MUST fail for each scenario no citation
  names, naming the document, the line and the call to write; an active
  entry without scenarios MUST fail when no citation names it at all.
- The waiver file lists, one per line with a reason after a tab, an id or
  a scenario as ID/handle, never a glob. A matching line absorbs the
  missing test.
- A waiver line that absorbs nothing, gives no reason or is a glob MUST
  fail.

Why: coverage is how the gate senses drift. A missing anchor means code
nobody can find from the decision, and no reason makes that acceptable;
a missing test, by contrast, may have a reason worth stating.

Scenarios:
1. `not-realized`: given an active decision with no realizes anchor,
   when the gate runs, then it fails naming the directive to add.
2. `scenario-unverified`: given an active entry's scenario that no test
   cites, when the gate runs, then it fails naming the scenario and the
   citation.
3. `spec-unverified`: given an active entry without scenarios that no
   test cites, when the gate runs, then it fails.
4. `waiver-absorbs-a-scenario`: given a waiver line naming an uncited
   scenario, when the gate runs, then the missing test is no problem.
5. `waiver-never-absorbs-an-anchor`: given an entry whose tests are
   waived and which has no anchor, when the gate runs, then it still
   fails.
6. `idle-waiver-fails`: given a waiver that absorbs nothing, or one
   without a reason, when the gate runs, then it fails naming its line.
7. `glob-refused` (owner): given a waiver written as a glob, when the gate
   runs, then it fails naming its line, so what a waiver accepts is
   always on its line.

Changes:
- 2026-10-09: a glob waiver, which absorbed every scenario of an entry at
  once, is refused.
- 2026-10-09: that a concerns anchor realizes nothing moved to CH-022,
  the concerns directive being its subject: not-realized was tested with
  a decision anchored by concerns alone, as CH-022's
  concerns-is-not-realized is, and is now tested with one anchored by
  nothing; the text names CH-022 for concerns.
- 2026-11-26: serves CH-261126-P2SK, CH-001, where it served CH-001: the
  needs are re-cut one a problem, and this decision answers the one that
  names it (backlog 004).

## CH-031 The gate MUST NOT ask an anchor of a need, and MUST count a need's citation only from the stakeholder's side

Status: active
Serves: CH-261126-P2SK

- A need is realized by the decisions that serve it, so the gate MUST NOT
  ask an anchor of it.
- When charter.toml names the test files on the stakeholder's side, a
  citation of a need from any other file MUST fail naming the file and
  MUST NOT count. When it names none, any test file counts.

Why: a need is a promise to a person, so only a test that uses the
software the way that person does can show it is kept. A unit test of a
package proves the package, not the promise.

Scenarios:
1. `need-without-anchor-passes`: given a need with no anchor and a
   citation, when the gate runs, then the need is covered.
2. `offside-citation-refused`: given a citation of a need from a file
   charter.toml does not name, when the gate runs, then it fails and the
   need stays unverified.
3. `onside-citation-counts`: given a citation of a need from a named
   file, when the gate runs, then the need is verified.

Changes:
- 2026-11-26: the stakeholder tests are named in charter.toml, where they
  were in the gate's options (backlog 003).
- 2026-11-26: serves CH-261126-P2SK, where it served CH-001: the needs are
  re-cut one a problem, and this decision answers the one that names it
  (backlog 004).

## CH-032 A probe MUST cite only observations, and a test MUST NOT cite an observation

Status: active
Serves: CH-261126-P2SK

- A probe is a test file under a path charter.toml names for probes, in
  any language the charter reads.
- A probe that references an entry MUST fail, and its references MUST NOT
  count as citations.
- A test file that is not a probe and references an observation MUST
  fail.
- Without a probe path in charter.toml, observations MUST be one
  configuration problem, and no probe rule runs.

Why: a probe reaches a live outside system and an ordinary test does
not, so a fact about an outside system and a promise of the code must
not be checked by the same run.

Scenarios:
1. `probe-reference-is-no-citation`: given an entry referenced only from
   a probe, when the gate runs, then the entry is unverified and the probe
   is reported.
2. `test-citing-an-observation-fails`: given a test file referencing an
   observation, when the gate runs, then it fails naming the file.
3. `probe-by-path` (owner, backlog 003): given a test file under the probe
   path and one outside it, in Go, Python or TypeScript, when the tree is
   read, then the first is a probe and the second an ordinary test.
4. `probe-path-required` (owner, backlog 003): given observations and no
   probe path in charter.toml, when the gate runs, then it reports one
   problem and blames nothing on a probe.

Changes:
- 2026-11-26: a probe is a test under the path charter.toml names, in any
  language, where it was a Go file carrying the probe tag;
  `negated-tag-is-no-probe` and `tag-required` are dropped for
  `probe-by-path` and `probe-path-required` (backlog 003).
- 2026-11-26: serves CH-261126-P2SK, where it served CH-005: CH-005 is
  deprecated and the need above it is served directly (backlog 004).

## CH-033 The gate MUST enumerate every surface and require each to be accounted for by an entry or listed in a baseline of live, unaccounted surfaces

Status: active
Kind: invariant
Serves: CH-261126-P2SK
Accounts: env:CHARTER_UPDATE_BASELINE

- The surfaces are the environment variables, the values of the string
  constants of each language's constants module that charter.toml names,
  and the keys of each kind charter.toml declares (CH-261126-NTRR).
- A surface that no entry accounts for and the baseline does not list
  MUST fail.
- An entry accounting for a surface nothing enumerates MUST fail.
- A baseline line whose surface is gone, which an entry now accounts for,
  which appears twice, or which lacks a reason MUST fail.
- A constants module charter.toml names that declares nothing MUST be a
  configuration problem.
- With the update variable set, the gate MUST rewrite the baseline with
  every unaccounted surface, keeping the reasons it had.
- A deprecated entry's Accounts MUST own nothing: what it accounted for is
  for a live entry to claim, or for the baseline.

Why: the baseline is a debt of surfaces nobody has ruled on yet. A gate
that lets it grow silently, or keeps listing what has been paid off, is
not keeping the books, so every line must be live, reasoned and unique.

Scenarios:
1. `surfaces-enumerated` (owner, backlog 003): given the fixture's
   environment constants in Go, Python and TypeScript and the keys its
   declared kinds print, when the gate runs, then those are its surfaces.
2. `unaccounted-fails` (owner, backlog 003): given a new key printed by a
   declared kind, when the gate runs, then it fails naming the key.
3. `accounts-must-exist`: given an entry accounting for a surface nothing
   enumerates, when the gate runs, then it fails.
4. `baseline-holds-live-gaps-once`: given a baseline line for a surface
   that is gone, now accounted for, listed twice or without a reason,
   when the gate runs, then it fails.
5. `empty-enumerator` (owner, backlog 003): given a constants module that
   declares nothing, when the gate runs, then it reports a configuration
   problem.
6. `update-rebuilds`: given the update variable set, when the gate runs,
   then it rewrites the baseline with every unaccounted surface, keeping
   the reasons it had, and a second run changes nothing.
7. `deprecated-accounts-nothing` (owner): given a surface accounted for
   only by a deprecated entry, when the gate runs, then the surface is
   unaccounted.

Changes:
- 2026-10-09: a deprecated entry's Accounts counted; they own nothing now.
- 2026-11-26: environment variables are the one built-in kind, in every
  language; build tags, OpenAPI operations and event kinds are kinds an
  adopter declares (CH-261126-NTRR), and the scenarios name the declared
  kinds (backlog 003).
- 2026-11-26: serves CH-261126-P2SK, where it served CH-004: CH-004 is
  deprecated and the need above it is served directly (backlog 004).

## CH-261126-NTRR charter check MUST read a declared surface kind's keys as the output lines of a command run from the root

Status: active
Serves: CH-261126-P2SK
See: CH-033

- A `[[surface]]` table in charter.toml names a kind, a lower-case word,
  and a command, which `charter check` MUST run from the root through the
  shell.
- Each non-empty line the command prints MUST be a key of that kind,
  written `kind:key` in Accounts and in the baseline.
- A command that exits non-zero, or prints a line twice, MUST fail naming
  the kind and what it printed.
- A kind declared twice, named `env`, or not a lower-case word MUST be a
  configuration problem naming its line.

Why: a surface is whatever the adopter's users can reach, and only the
adopter knows where those surfaces are listed: a schema, an OpenAPI file,
a list of build tags. A command lets them point at that list in whatever
language they use, and the gate keeps those keys in the same books as
every other surface.

Scenarios:
1. `lines-are-keys` (owner, backlog 003): given a declared kind whose
   command prints two lines, when the gate runs, then it has two
   surfaces, `kind:first` and `kind:second`, each to be accounted for or
   listed.
2. `failing-command-named` (owner, backlog 003): given a command that
   exits non-zero, when the gate runs, then it fails naming the kind, the
   command and what it printed.
3. `duplicate-key-named` (owner, backlog 003): given a command printing
   one line twice, when the gate runs, then it fails naming the kind and
   the key.
4. `kind-form-checked` (owner, backlog 003): given a kind declared twice,
   named `env`, or not a lower-case word, when the gate runs, then it
   reports a configuration problem naming its line in charter.toml.

Changes:
- 2026-11-26: serves CH-261126-P2SK, where it served CH-004: CH-004 is
  deprecated and the need above it is served directly (backlog 004).

## CH-034 Code MUST read an environment variable through a named constant, or under a dynamic directive, and a literal MUST NOT name one

Status: active
Kind: invariant
Serves: CH-261126-P2SK

- Outside the constants modules and outside test files, a read of the
  environment MUST take a constant of a constants module as its argument,
  or its line MUST carry the dynamic directive.
- A read is, in Go, a call of a configured reader, `os.Getenv` and
  `os.LookupEnv` by default; in Python `os.getenv`, `os.environ.get` or
  `os.environ[...]`; in TypeScript `process.env`.
- A string literal equal to a declared variable's value, or matching a
  pattern charter.toml gives, MUST fail.

Why: an environment variable read by a literal name is a surface the
enumerators cannot see, so it would escape the books CH-033 keeps.
Reading only through named constants keeps every variable enumerable.

Scenarios:
1. `literal-argument`: given a reader called with a string literal, when
   the gate runs, then it fails naming the constant to use.
2. `variable-argument`: given a reader called with a variable, when the
   gate runs, then it fails.
3. `duplicated-name`: given a literal equal to a declared variable's
   value, when the gate runs, then it fails naming the constant.
4. `undeclared-pattern`: given a literal matching a configured pattern,
   when the gate runs, then it fails.
5. `dynamic-directive`: given a read whose line carries the directive,
   when the gate runs, then it is exempt.
6. `python-read-checked` (owner, backlog 003): given a Python
   `os.environ.get` of a string literal and one of a constant of the
   constants module, when the gate runs, then the first fails naming the
   constant and the second passes.
7. `typescript-read-checked` (owner, backlog 003): given a TypeScript
   `process.env` read by a literal name and one by a constant of the
   constants module, when the gate runs, then the first fails naming the
   constant and the second passes.

Changes:
- 2026-11-26: environment reads are checked in Python and TypeScript as in
  Go, each language's constants module named in charter.toml;
  `python-read-checked` and `typescript-read-checked` are added (backlog
  003).
- 2026-11-26: serves CH-261126-P2SK, where it served CH-004: CH-004 is
  deprecated and the need above it is served directly (backlog 004).

## CH-035 An observation MUST carry its status's fields, and the gate MUST fail a broken observation under an active entry

Status: active
Serves: CH-261126-P2SK

- Every observation MUST name a subject.
- An observed observation MUST have a calendar date and at least one
  citing probe; a documented one MUST have a source; any date given MUST
  be a calendar date.
- A broken observation supporting an active entry MUST fail; under a
  planned or a deprecated entry it MUST NOT.

Why: a fact about an outside system is only worth recording with where it
came from and when it was last seen true. A broken fact under an active
decision means the decision rests on something false, which the gate
must surface.

Scenarios:
1. `observed-needs-date-and-probe`: given an observed observation without
   a date, with a date not on the calendar, or cited by no probe, when the
   gate runs, then it fails.
2. `documented-needs-source`: given a documented observation without a
   source, or with a stray date, when the gate runs, then it fails.
3. `subject-required`: given an observation without a subject, when the
   gate runs, then it fails.
4. `broken-under-active`: given a broken observation under an active
   entry and one under a planned or deprecated entry, when the gate runs,
   then the first fails and the second does not.

Changes:
- 2026-11-26: serves CH-261126-P2SK, where it served CH-005: CH-005 is
  deprecated and the need above it is served directly (backlog 004).

## CH-036 The walk MUST cover every Go, Python and TypeScript file under the root, skipping another charter's tree and what charter.toml ignores

Status: active
Serves: CH-006

- The walk MUST read every file under the root of each language
  charter.toml names: `.go` files, `.py` files, and `.ts` and `.tsx` files
  that are not declaration files.
- It MUST cross into nested Go modules.
- It MUST skip `.git`, `node_modules`, `vendor` and `testdata`; hidden and
  underscore directories; a directory below the root that holds its own
  `claims/charter.toml`; every ignore pattern, matched against the path and
  the base name; and, when `only` names paths, everything outside them.

Why: two charters in one repository must not see each other's anchors,
and one charter's code may sit in several modules. The boundary between
charters is therefore the other charter's configuration, not the build
tool's module boundary.

Scenarios:
1. `hidden-skipped` (owner, backlog 003): given a hidden and an
   underscore directory, when the tree is walked, then neither is read.
2. `nested-module-walked` (owner, backlog 003): given a Go module nested
   below the root, when the tree is walked, then its anchors and
   citations are seen.
3. `other-charter-skipped` (owner, backlog 003): given a directory below
   the root with its own `claims/charter.toml`, when the tree is walked,
   then its anchors and citations are unseen.
4. `ignore-patterns`: given an ignore pattern, when the tree is walked,
   then it matches a path or a base name.
5. `only-restricts`: given `only` naming paths, when the tree is walked,
   then a directory outside them is not read and its anchors and
   citations are unseen.
6. `typescript-listed`: given TypeScript configured, when the tree is
   walked, then it lists the TypeScript files and no declaration file.
7. `python-listed` (owner, backlog 003): given Python configured, when the
   tree is walked, then it lists every `.py` file; given Python not
   configured, then none is read.

Changes:
- 2026-11-26: the walk is of the repository: it reads Python, crosses into
  nested Go modules, and stops at a directory holding another charter's
  claims/charter.toml, where it stopped at a nested go.mod;
  `hidden-and-nested-skipped` is dropped for `hidden-skipped`,
  `nested-module-walked` and `other-charter-skipped`, and `python-listed`
  is added (backlog 003).

## CH-037 charter check MUST report every problem, sorted, with its file, line, rule and fix, and MUST pass a clean tree

Status: active
Serves: CH-001, CH-002
See: CH-261126-33GK

`charter check` MUST print each problem on a line of its own, sorted, and
exit non-zero when it prints any; on a clean tree it MUST print nothing
and exit zero. How it finds the root and checks the configuration is
CH-261126-33GK's.

Why: a gate whose output a person cannot act on is noise. Each problem
on its own sorted line, with what to fix, makes the output a worklist,
and a silent zero exit makes a clean tree unambiguous to a script.

Scenarios:
1. `clean-fixture-passes`: given the clean fixture, when the gate runs,
   then it reports nothing.
2. `exit-status` (owner, backlog 003): given a tree with a problem and the
   clean fixture, when `charter check` runs on each, then the first exits
   non-zero and the second exits zero with no output.

Changes:
- 2026-11-26: charter check prints the problems and sets its exit status,
  where Check reported them as test errors; how the root is found and the
  configuration checked moved to CH-261126-33GK, `root-resolved` and
  `typescript-options-validated` with them, and `exit-status` is added
  (backlog 003).

## CH-261127-QRQS charter check MUST warn, without failing, of a planned entry fully tested, a deprecated one still in the code, and a doc to read again

Status: active
Serves: CH-001
See: CH-037, CH-261127-GXP7

- After its problems, `charter check` MUST write each warning on a line of
  its own to standard error, with its file and line, for:
  - a planned entry each of whose scenarios a test cites, to set active;
  - a deprecated entry still anchored or cited, with each site, to strip
    before it is deleted;
  - each doc CH-261127-GXP7 names drafted or stale.
- Warnings MUST NOT change the exit status, and with none, check MUST
  write nothing to standard error.

Why: each warning is a step of the loop not yet taken, worth seeing on
every run but wrong to block a merge on. Keeping warnings apart from the
problems, on standard error, leaves the gate's own output and exit status
exactly as they were.

Scenarios:
1. `planned-fully-tested` (owner, backlog 008): given a planned entry with
   every scenario cited, when check runs, then it warns to set it active.
2. `deprecated-still-in-code` (owner, backlog 008): given a deprecated
   entry with an anchor and a citation, when check runs, then it warns at
   each.
3. `warnings-never-fail` (owner, backlog 008): given a tree whose only
   findings are warnings, when check runs, then it passes and writes
   nothing to standard output.
