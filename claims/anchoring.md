# Anchors, citations and symbols

Code: chartertest, the charter package

How code and tests name the entries, and what the compiler and the gate
hold those names to.

## CH-020 Every entry has a generated symbol per language read, its id with underscores for hyphens, its title the doc comment

Status: active
Serves: CH-002, CH-003
See: CH-012, CH-261126-PQ7X

For each language charter.toml names, charter gen writes one file in the
claims directory declaring a symbol per entry of any status, named by the
id with each hyphen as an underscore, holding the id, the title, the
status and the kind, with the title as its doc comment; a planned entry is
marked planned, a deprecated one carries a deprecation notice naming its
document. Each file is self-contained: it declares the types, the verifies
and the observes it offers and imports nothing outside its language's
standard library. For Go the file is specs.go in a module of its own whose
go.mod requires nothing, for TypeScript index.ts, for Python claims.py.
The order the files declare the entries in is CH-012's. charter check
fails a file that differs from its rendering, or is missing, naming
charter gen, and never writes one.

Why: a symbol the compiler or the import checks is the link that cannot
rot, and one without the slug is one a title can be refined under; a file
that imports nothing of the tooling is one an adopter's build never
fetches.

Scenarios:
1. `bare-symbol-declared`: the generated Go declares each entry as its id
   with underscores for hyphens, a dated id as SW_261001_K7Q2, with no
   lens.
2. `planned-and-deprecated-marked`: a planned entry is declared and marked
   planned; a deprecated one carries a deprecation notice naming its
   document.
3. `typescript-module-rendered`: the TypeScript module equals the
   committed fixture module, escapes a comment closer in a title, and maps
   each symbol to its id.
4. `python-module-rendered` (owner, backlog 003): the Python module equals
   the committed fixture module, declares each entry as its symbol, and
   imports only the standard library.
5. `go-module-self-contained` (owner, backlog 003): the generated Go
   module's go.mod requires nothing, and its file imports only the
   standard library.
6. `title-only-comment` (owner, backlog 003): each symbol's doc comment is
   its entry's title, and no generated file holds an entry's text.
7. `stale-generated-fails` (owner, backlog 003): a generated file of any
   language that differs from its rendering, or is missing, fails naming
   charter gen.
8. `gen-writes` (owner, backlog 003): charter gen writes every generated
   file, a second run changes nothing, and charter check writes none.

Changes:
- 2026-10-01: a dated id holds two hyphens, each an underscore in its
  symbol, and the generated files follow the order of the ids, numbered
  then dated, where they followed the number (CH-012).
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-10-09: the order the generated files declare the entries in moved
  to CH-012, whose ordered-by-date states it, the order ids sort in being
  its subject: bare-symbol-declared no longer names it and keeps that
  every entry is declared, and the text names CH-012 for the order.
- 2026-10-09: a deprecated entry's marking is ruled here alone, the
  generated files being this entry's subject: CH-013's deprecated-marked
  restated planned-and-deprecated-marked and is dropped, its test citing
  this one, which now names the notice's document as the dropped one
  did, and so does the text.
- 2026-11-26: the symbols are generated per language charter.toml names,
  Python's claims.py beside Go's and TypeScript's, each file
  self-contained and Go's in a module of its own (backlog 003). A symbol's
  doc comment is its title (Q-261015-CPW5). charter gen is the one writer,
  so the update variables and their Accounts are gone and `update-writes`
  is dropped for `gen-writes` (Q-261015-Y682); `stale-generated-fails`
  names any language and the command.

## CH-021 Production code anchors an entry in a Go doc comment, a TypeScript line comment, or above a Python def

Status: active
Serves: CH-002

A Go directive is a comment line "//charter:realizes ID, ..." or
"//charter:concerns ID, ...", with or without a space after the slashes,
in the doc comment of a package clause or of a declaration; the anchor
records the kind, the file, the line and the declaration it documents. A
TypeScript directive is the same text in a line comment of a file that is
not a test; the anchor records the file and the line. A Python directive
is "# charter:realizes ID, ..." or "# charter:concerns ID, ..." on the
comment lines directly above a def or a class, or above its decorators,
in a file that is not a test; the anchor records the file, the line and
the declaration, and a Python directive anywhere else fails naming the
declaration to move it above.

Why: the anchor is the pointer from the decision to the code, and the
declaration it sits on is the place an agent is sent to.

Scenarios:
1. `go-anchor-recorded`: a directive on a package clause is recorded as
   the package's anchor with its declaration.
2. `spaced-directive-counts`: a directive written with a space after the
   slashes is an anchor too.
3. `typescript-anchor-resolves`: a TypeScript directive naming an unknown
   id fails naming the file and the line.
4. `comment-elsewhere-refused` (owner): a directive in a
   comment that documents no declaration, inside a function or stray in a
   file, fails naming the doc comment to move it to, and counts for
   nothing.
5. `python-anchor-recorded` (owner, backlog 003): a Python directive above
   a def, or above its decorators, is recorded as that def's anchor with
   its file and line.
6. `python-comment-elsewhere-refused` (owner, backlog 003): a Python
   directive in a body, parted from its def by a blank line, or stray in a
   file fails and counts for nothing.

Changes:
- 2026-10-09: a directive outside a declaration's doc comment was read as
  an anchor on nothing; it is refused.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-11-26: Python is read: a directive on the comment lines directly
  above a def or a class, or its decorators, anchors it, and one anywhere
  else is refused (backlog 003).

## CH-022 An anchor names an entry that exists and is not planned, and concerns never counts toward coverage

Status: active
Serves: CH-002

A directive naming an id that is not an entry fails; one naming a planned
entry fails, saying to set it active as its work begins; a directive with
no id fails. A concerns anchor is navigation: an active decision with only
a concerns anchor is unrealized.

Why: an anchor to nothing is a lie the compiler cannot catch; a realizes
on a planned entry is work that began without saying so.

Scenarios:
1. `unknown-id-refused`: a directive naming an id that is not an entry
   fails naming it.
2. `planned-refused`: a realizes naming a planned entry fails.
3. `concerns-is-not-realized`: an active decision anchored by concerns
   alone fails as realized by no code.

Changes:
- 2026-10-09: that a concerns anchor realizes nothing is ruled here
  alone, the concerns directive being this entry's subject: CH-030's
  not-realized was tested with a decision anchored by concerns alone, as
  concerns-is-not-realized is, and is now tested with one anchored by
  nothing, CH-030 pointing here.

## CH-023 A test cites an entry by a Verifies statement in a test function, naming by handle the scenarios it exercises

Status: active
Serves: CH-001, CH-002

A Go citation is a call of Verifies inside a function whose name starts
with Test, whose second argument is the entry's symbol and whose further
arguments are string literals naming the entry's scenarios. A TypeScript
citation is a call of verifies inside a test call, and a Python one a
call of verifies in a test_ method of a unittest.TestCase, its first
argument self, or in a pytest test_ function, its first argument None,
each with the same shape. The calls open the test, or a subtest it runs:
they are its first statements, one after another, and a call that
anything precedes or encloses, or in a function its runner never runs,
such as Go's Testx, fails. A TypeScript call in a test
marked skip, fixme, fail or todo, or to a verifies the test file declares
itself, fails. A call anywhere else, a spec named by anything but its
symbol, a handle that is not a literal, or a handle the entry does not
list, each fail naming the file and the line. Any other reference to an entry from a test
file is a reference, not a citation.

Why: the gate reads the call, not the run; a call it can read is one it
can hold to the scenarios it names.

Scenarios:
1. `go-citation-read`: each Verifies call in the fixture is a citation
   with its file, test function, line and handles.
2. `outside-a-test-refused`: a Verifies in a helper function fails, and
   its handles count for nothing.
3. `unknown-handle-refused`: a handle the entry does not list fails naming
   the entry and the handle.
4. `literal-arguments`: a handle that is not a string literal, or a spec
   named by a variable, fails.
5. `typescript-citation-read`: a verifies call in a TypeScript test
   verifies the handles it names, and one outside a test, or with a
   non-literal handle, fails.
6. `alias-followed`: a citation through an import alias is read.
7. `reference-is-not-citation`: a bare reference to an entry from a test
   file is recorded as a reference and counts for no scenario.
8. `later-statement-refused` (owner): a Verifies that another
   statement precedes, or that a block encloses, fails and counts for
   nothing; one opening a subtest's body counts.
9. `non-test-name-refused` (owner): a Verifies in a function
   named Test followed by a lower-case letter, which go test never runs,
   fails naming the function.
10. `disabled-ts-refused` (owner): a verifies in a TypeScript
    test marked skip, fixme or fail fails and counts for nothing.
11. `local-verifies-refused` (owner): a verifies declared in
    the test file itself fails and counts for nothing.
12. `python-citation-read` (owner, backlog 003): a verifies opening a
    unittest test_ method, or a pytest test_ function, is a citation with
    its file, test, line and handles.
13. `python-later-statement-refused` (owner, backlog 003): a Python
    verifies that another statement precedes, or in a method or function
    the runner does not collect, fails and counts for nothing.

Changes:
- 2026-10-09: a citation counted wherever it sat in a test, after a skip,
  under a false branch, in a function go test never runs, in a skipped
  TypeScript test or through a verifies the test declared; each is refused.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-11-26: Python is read: a verifies opening a unittest test_ method
  or a pytest test_ function is a citation, held to the same refusals as
  Go's and TypeScript's (backlog 003).

## CH-024 Each language's verifies records each citation and its test's outcome as a JSON line in the events file

Status: active
Serves: CH-001
Accounts: env:CHARTER_EVENTS

When CHARTER_EVENTS names a file, the generated verifies of every language
appends to it, once the test ends, one JSON line per citation: the id, the
handle when one is named, the test, its file and its outcome, passed,
failed or skipped. With the variable unset it records nothing and the test
runs as it would. Lines appended by tests running at once are each whole.
A nil or missing symbol is refused at run time. The report reads these
lines.

Why: the line turns coverage from a static claim into a sensor: one run
says which scenarios passed, failed or were skipped, and one format
written by every language is one join for the report, not one per runner.

Scenarios:
1. `event-per-citation` (owner, backlog 003): a Go, a Python and a
   TypeScript test citing a scenario each append a line with the id, the
   handle, the test, the file and the outcome.
2. `outcome-recorded` (owner, backlog 003): a passing, a failing and a
   skipped test record passed, failed and skipped, a skip after the
   citation included.
3. `unset-records-nothing` (owner, backlog 003): with CHARTER_EVENTS
   unset, a citing test writes no file and passes as it would.
4. `concurrent-lines-whole` (owner, backlog 003): tests running in
   parallel and appending at once leave every line a whole JSON object.
5. `nil-spec-refused`: Verifies with a nil spec fails the test.

Changes:
- 2026-11-26: the marker is a JSON line per citation in the file
  CHARTER_EVENTS names, written by every language's generated verifies
  with the test's outcome, where Go's was a test attribute in go test's
  JSON stream; `attr-per-citation` is dropped for `event-per-citation`,
  and `outcome-recorded`, `unset-records-nothing` and
  `concurrent-lines-whole` are added (backlog 003).

## CH-025 The claims directory holds the record, charter.toml and generated files, and no test assigns to an entry

Status: active
Serves: CH-003

Beside the markdown of the record and charter.toml, the claims directory
holds only the files charter gen writes; a source file there that gen did
not write fails, saying entries and observations are written in markdown.
No test file assigns to an entry or a field of one, and no Go test
dot-imports the generated package.

Why: the generated files are the registry in each language, and anything
else written there is a spec or a fact the documents do not know; a
mutated entry is a lie to every test after it.

Scenarios:
1. `hand-written-source-refused` (owner, backlog 003): a source file in
   the claims directory that charter gen did not write fails, saying
   entries and observations are written in markdown.
2. `assignment-refused`: an assignment to an entry, or to one of its
   fields, from a test file or from the package's own tests, fails.
3. `dot-import-refused`: a test file dot-importing the entries package
   fails.

Changes:
- 2026-11-26: observations are markdown (CH-261126-HK6S), so the claims
  directory holds no hand-written code at all: `functions-refused` and
  `spec-in-go-refused` are dropped for `hand-written-source-refused`
  (backlog 003).
