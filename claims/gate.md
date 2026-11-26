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
anchor of it. When charter.toml names the test files on the stakeholder's
side, a citation of a need from any other file fails naming the file, and
does not count; when they name none, any test file counts.

Why: a need is a promise to a person; a unit test of a package proves the
package, not the promise.

Scenarios:
1. `need-without-anchor-passes`: a need with no anchor and a citation is
   covered.
2. `offside-citation-refused`: a citation of a need from a file
   charter.toml does not name fails and leaves the need unverified.
3. `onside-citation-counts`: a citation from a named file verifies the
   need.

Changes:
- 2026-11-26: the stakeholder tests are named in charter.toml, where they
  were in the gate's options (backlog 003).

## CH-032 A probe cites observations only, and a test never cites an observation

Status: active
Serves: CH-005

A probe is a test file under a path charter.toml names for probes, in any
language the charter reads. A probe referencing an entry fails, and its
references are no citations; a test file that is not a probe referencing
an observation fails. Without a probe path in charter.toml, observations
are one configuration problem and no probe rule runs.

Why: a probe reaches a live system and a test does not; a fact and a
promise must not be checked by the same run.

Scenarios:
1. `probe-reference-is-no-citation`: an entry referenced only from a probe
   is unverified, and the probe is reported.
2. `test-citing-an-observation-fails`: a test file referencing an
   observation fails naming the file.
3. `probe-by-path` (owner, backlog 003): a test file under the probe path,
   in Go, Python or TypeScript, is a probe, and one outside it is an
   ordinary test.
4. `probe-path-required` (owner, backlog 003): observations without a
   probe path in charter.toml are one problem, and nothing is blamed on a
   probe.

Changes:
- 2026-11-26: a probe is a test under the path charter.toml names, in any
  language, where it was a Go file carrying the probe tag;
  `negated-tag-is-no-probe` and `tag-required` are dropped for
  `probe-by-path` and `probe-path-required` (backlog 003).

## CH-033 Every surface is enumerated, and accounted for by an entry or listed in a baseline of only live, unaccounted surfaces

Status: active
Kind: invariant
Serves: CH-004
Accounts: env:CHARTER_UPDATE_BASELINE

The surfaces are the environment variables, the values of the string
constants of each language's constants module that charter.toml names,
and the keys of each kind charter.toml declares (CH-261126-NTRR). A surface no entry accounts for and the baseline does not list
fails; an entry accounting for a surface nothing enumerates fails; a
baseline line whose surface is gone, or which an entry now accounts for,
or which appears twice, or which lacks a reason, fails. A constants
module charter.toml names that declares nothing is a misconfiguration. With the update variable set the gate rewrites the
baseline with every unaccounted surface, keeping the reasons it had. A
deprecated entry's Accounts own nothing: what it accounted for is a live
entry's to claim, or the baseline's.

Why: the baseline is a debt; a gate that lets it grow silently, or hold
what is paid off, is not keeping the books.

Scenarios:
1. `surfaces-enumerated` (owner, backlog 003): the fixture's environment
   constants in Go, Python and TypeScript, and the keys its declared kinds
   print, are the surfaces.
2. `unaccounted-fails` (owner, backlog 003): a new key a declared kind
   prints fails naming it.
3. `accounts-must-exist`: an entry accounting for a surface nothing
   enumerates fails.
4. `baseline-holds-live-gaps-once`: a baseline line for a surface that is
   gone, or now accounted for, or listed twice, or without a reason,
   fails.
5. `empty-enumerator` (owner, backlog 003): a constants module that
   declares nothing is a configuration problem.
6. `update-rebuilds`: the update variable rewrites the baseline with every
   unaccounted surface, keeps the reasons it had, and a second run changes
   nothing.
7. `deprecated-accounts-nothing` (owner): a surface accounted
   for only by a deprecated entry is unaccounted.

Changes:
- 2026-10-09: a deprecated entry's Accounts counted; they own nothing now.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-11-26: environment variables are the one built-in kind, in every
  language; build tags, OpenAPI operations and event kinds are kinds an
  adopter declares (CH-261126-NTRR), and the scenarios name the declared
  kinds (backlog 003).

## CH-261126-NTRR A declared surface kind is a command run from the root whose output lines are its keys

Status: active
Serves: CH-004
See: CH-033

A [[surface]] table in charter.toml names a kind, a lower-case word, and a
command, which charter check runs from the root through the shell. Each
non-empty line the command prints is a key of that kind, written kind:key
in Accounts and in the baseline. A command that exits non-zero, or prints
a line twice, fails naming the kind and what it printed. A kind declared
twice, named env, or not a lower-case word is a configuration problem
naming its line.

Why: a surface is whatever the adopter's users can reach, and only the
adopter knows where those are listed: a schema, an OpenAPI file, a list of
build tags. A command lets them say where in whatever language they have,
and the gate keeps those keys in the same books as every other surface.

Scenarios:
1. `lines-are-keys` (owner, backlog 003): a declared kind whose command
   prints two lines has two surfaces, kind:first and kind:second, each
   accounted for or listed.
2. `failing-command-named` (owner, backlog 003): a command that exits
   non-zero fails naming the kind, the command and what it printed.
3. `duplicate-key-named` (owner, backlog 003): a command printing one line
   twice fails naming the kind and the key.
4. `kind-form-checked` (owner, backlog 003): a kind declared twice, named
   env, or not a lower-case word is a configuration problem naming its
   line in charter.toml.

## CH-034 An environment variable is read through a named constant, or under a dynamic directive, and no literal names one

Status: active
Kind: invariant
Serves: CH-004

Outside the constants modules and outside test files, a read of the
environment takes a constant of a constants module as its argument, or its
line carries the dynamic directive. In Go a read is a call of a configured
reader, os.Getenv and os.LookupEnv by default; in Python os.getenv,
os.environ.get or os.environ[...]; in TypeScript process.env. A string
literal equal to a declared variable's value, or matching a pattern
charter.toml gives, fails.

Why: a variable read by a literal is a surface the enumerators cannot see.

Scenarios:
1. `literal-argument`: a reader called with a string literal fails naming
   the constant to use.
2. `variable-argument`: a reader called with a variable fails.
3. `duplicated-name`: a literal equal to a declared variable's value fails
   naming the constant.
4. `undeclared-pattern`: a literal matching a configured pattern fails.
5. `dynamic-directive`: a read whose line carries the directive is exempt.
6. `python-read-checked` (owner, backlog 003): a Python os.environ.get of
   a string literal fails naming the constant, and one of a constant of
   the constants module passes.
7. `typescript-read-checked` (owner, backlog 003): a TypeScript
   process.env read by a literal name fails naming the constant, and one
   by a constant of the constants module passes.

Changes:
- 2026-11-26: environment reads are checked in Python and TypeScript as in
  Go, each language's constants module named in charter.toml;
  `python-read-checked` and `typescript-read-checked` are added (backlog
  003).

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

## CH-036 The walk covers every Go, Python and TypeScript file under the root, skipping another charter's tree and what charter.toml ignores

Status: active
Serves: CH-006

The walk reads every file under the root of each language charter.toml
names: .go files, .py files, and .ts and .tsx files that are not
declaration files. It crosses into nested Go modules. It skips .git,
node_modules, vendor and testdata, hidden and underscore directories, a
directory below the root that holds its own claims/charter.toml, every
ignore pattern matched against the path and the base name, and everything
outside only when only names paths.

Why: two charters in one repository must not see each other's anchors, and
one charter's code may sit in several modules; the boundary is the other
charter's configuration, not the build tool's.

Scenarios:
1. `hidden-skipped` (owner, backlog 003): a hidden and an underscore
   directory are not walked.
2. `nested-module-walked` (owner, backlog 003): a Go module nested below
   the root is walked, and its anchors and citations are seen.
3. `other-charter-skipped` (owner, backlog 003): a directory below the
   root holding its own claims/charter.toml is not walked, so its anchors
   and citations are unseen.
4. `ignore-patterns`: a pattern matches a path or a base name.
5. `only-restricts`: with Only naming paths, a directory outside them is
   not walked, so its anchors and citations are unseen.
6. `typescript-listed`: with TypeScript configured, the tree lists the
   TypeScript files and no declaration file.
7. `python-listed` (owner, backlog 003): with Python configured the tree
   lists every .py file, and without it none is read.

Changes:
- 2026-11-26: the walk is of the repository: it reads Python, crosses into
  nested Go modules, and stops at a directory holding another charter's
  claims/charter.toml, where it stopped at a nested go.mod;
  `hidden-and-nested-skipped` is dropped for `hidden-skipped`,
  `nested-module-walked` and `other-charter-skipped`, and `python-listed`
  is added (backlog 003).

## CH-037 The gate reports every problem, sorted, with its file, line, rule and fix, and passes on a clean tree

Status: active
Serves: CH-001, CH-002
See: CH-261126-33GK

charter check prints each problem on a line of its own, sorted, and exits
non-zero when it prints any; on a clean tree it prints nothing and exits
zero. How it finds the root and checks the configuration is
CH-261126-33GK's.

Why: a gate whose output a person cannot act on is noise.

Scenarios:
1. `clean-fixture-passes`: the clean fixture reports nothing.
2. `exit-status` (owner, backlog 003): charter check exits non-zero on a
   tree with a problem, and zero with no output on the clean fixture.

Changes:
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-11-26: charter check prints the problems and sets its exit status,
  where Check reported them as test errors; how the root is found and the
  configuration checked moved to CH-261126-33GK, `root-resolved` and
  `typescript-options-validated` with them, and `exit-status` is added
  (backlog 003).
