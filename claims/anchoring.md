# Anchors, citations and symbols

Code: chartertest, the charter package

How code and tests name the entries, and what the compiler and the gate
hold those names to.

## CH-020 Every entry has a generated Go and TypeScript symbol, its id with underscores for hyphens, its title the doc comment

Status: active
Serves: CH-002, CH-003
See: CH-012
Accounts: env:CHARTER_UPDATE_SPECS, env:CHARTER_UPDATE_TYPESCRIPT

The entries package's generated Go declares one variable per entry of any
status, named by the id with each hyphen as an underscore, holding the id,
the title, the status and the kind, with the title and the text as its
doc comment; a planned entry is marked planned, a deprecated one carries a
deprecation notice naming its document. The TypeScript module maps the
same symbols to the ids, grouped by prefix, with the title as a doc
comment. The order both declare the entries in is CH-012's. Each file must
equal its rendering, or the gate names the command to run; a variable of
each name makes the gate write the file instead.

Why: a symbol the compiler checks is the link that cannot rot; a symbol
without the slug is one a title can be refined under.

Scenarios:
1. `bare-symbol-declared`: the generated Go declares each entry as its id
   with underscores for hyphens, a dated id as SW_261001_K7Q2, with no
   lens.
2. `planned-and-deprecated-marked`: a planned entry is declared and marked
   planned; a deprecated one carries a deprecation notice naming its
   document.
3. `typescript-module-rendered`: the TypeScript module equals the
   committed fixture module, escapes a comment closer in a title, and
   maps each symbol to its id.
4. `stale-generated-fails`: a generated Go or TypeScript file that differs
   from its rendering, or is missing, fails naming the command to run.
5. `update-writes`: with the update variable set the gate writes the
   module and a second run is a fixed point.

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

## CH-021 Production code anchors an entry by a directive in a declaration's or package's doc comment, or a TypeScript line comment

Status: active
Serves: CH-002

A Go directive is a comment line "//charter:realizes ID, ..." or
"//charter:concerns ID, ...", with or without a space after the slashes,
in the doc comment of a package clause or of a declaration; the anchor
records the kind, the file, the line and the declaration it documents. A
TypeScript directive is the same text in a line comment of a file that is
not a test; the anchor records the file and the line.

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

Changes:
- 2026-10-09: a directive outside a declaration's doc comment was read as
  an anchor on nothing; it is refused.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

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
citation is a call of verifies inside a test call, with the same shape.
The calls open the test, or a subtest it runs: they are its first
statements, one after another, and a call that anything precedes or
encloses, or in a function go test never runs, such as Testx, fails. A TypeScript call in a test
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

Changes:
- 2026-10-09: a citation counted wherever it sat in a test, after a skip,
  under a false branch, in a function go test never runs, in a skipped
  TypeScript test or through a verifies the test declared; each is refused.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.

## CH-024 Verifies and verifies mark the running test with each citation, "ID" or "ID/handle", in the runner's structured output

Status: active
Serves: CH-001

The Go marker is a test attribute per citation, which go test's JSON
stream reports as an attr event on the test before its result; the
helper a TypeScript test runner calls may annotate the running test the
same way. A nil spec is refused at run time. The report reads these events.

Why: the attribute is what turns coverage from a static claim into a
sensor: one run says which scenarios passed, failed or were skipped.

Scenarios:
1. `attr-per-citation`: a test citing an entry emits an attr event with
   its id, and one naming a handle emits "ID/handle"; a skipped subtest
   reports its attr before its skip.
2. `nil-spec-refused`: Verifies with a nil spec fails the test.

## CH-025 The entries package holds only observations and the generated file, and no entry is assigned outside its declaration

Status: active
Serves: CH-003

A hand-written file of the entries package declares observations as
addressed composite literals with keyed, literal fields and nothing else:
no function, no other declaration, no spec written in Go. No test file
assigns to an entry or a field of one, and no test dot-imports the
package.

Why: the generated file is the registry's Go; anything else written there
is a spec the documents do not know, and a mutated entry is a lie to every
test after it.

Scenarios:
1. `functions-refused`: a function in the entries package fails.
2. `spec-in-go-refused`: a hand-written spec value fails, saying a spec is
   written in a document.
3. `assignment-refused`: an assignment to an entry, or to one of its
   fields, from a test file or from the package's own tests, fails.
4. `dot-import-refused`: a test file dot-importing the entries package
   fails.
