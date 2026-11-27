# From an id to its code and its tests

Code: internal/lang/golang, internal/lang/python, internal/lang/typescript, internal/lang/plain, internal/facts, internal/docs, cmd/charter

This area answers the problem that a person or an agent new to the code
cannot find their way from what was decided to where it is done. It
covers how code and tests name the entries in each language, the
generated symbols that make a wrong name fail, the command that lists an
entry's code and tests, and the docs that tell a story over them.

## CH-020 charter gen MUST generate a symbol per entry in each language read, named by the id with underscores for hyphens and documented by the title

Status: active
Serves: CH-002, CH-261126-PFG5
See: CH-012, CH-261126-PQ7X

- For each language charter.toml names, `charter gen` MUST write one file
  in the claims directory declaring a symbol per entry of any status,
  named by the id with each hyphen as an underscore, holding the id, the
  title, the status and the kind, with the title as its doc comment.
- A planned entry MUST be marked planned, and a deprecated one MUST carry
  a deprecation notice naming its document.
- Each file MUST be self-contained: it declares the types, the `verifies`
  and the `observes` it offers, and imports nothing outside its
  language's standard library. For Go the file is specs.go in a module of
  its own whose go.mod requires nothing; for TypeScript, index.ts; for
  Python, claims.py.
- The order the files declare the entries in is CH-012's.
- `charter check` MUST fail a file that differs from its rendering, or is
  missing, naming `charter gen`, and MUST NOT write one.

Why: a symbol that the compiler or the import checks is a link from a
test to an entry that cannot rot. Naming it by the id without a slug of
the title lets the title be refined without renaming the symbol, and a
file that imports nothing of the tooling is one an adopter's build never
has to fetch.

Scenarios:
1. `bare-symbol-declared`: given entries with numbered and dated ids, when
   Go is generated, then each is declared as its id with underscores for
   hyphens, a dated id as SW_261001_K7Q2, with no lens.
2. `planned-and-deprecated-marked`: given a planned and a deprecated
   entry, when the files are generated, then the first is declared and
   marked planned and the second carries a deprecation notice naming its
   document.
3. `typescript-module-rendered`: given the fixture, when TypeScript is
   generated, then the module equals the committed fixture module,
   escapes a comment closer in a title, and maps each symbol to its id.
4. `python-module-rendered` (owner, backlog 003): given the fixture, when
   Python is generated, then the module equals the committed fixture
   module, declares each entry as its symbol, and imports only the
   standard library.
5. `go-module-self-contained` (owner, backlog 003): given the generated Go
   module, when its go.mod and file are read, then the go.mod requires
   nothing and the file imports only the standard library.
6. `title-only-comment` (owner, backlog 003): given the generated files,
   when they are read, then each symbol's doc comment is its entry's
   title and no file holds an entry's text.
7. `stale-generated-fails` (owner, backlog 003): given a generated file of
   any language that differs from its rendering, or is missing, when
   check runs, then it fails naming `charter gen`.
8. `gen-writes` (owner, backlog 003): given a record, when `charter gen`
   runs twice and check once, then the first run writes every generated
   file, the second changes nothing, and check writes none.

Changes:
- 2026-10-01: a dated id holds two hyphens, each an underscore in its
  symbol, and the generated files follow the order of the ids, numbered
  then dated, where they followed the number (CH-012).
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
- 2026-11-26: serves CH-002, CH-261126-PFG5, where it served CH-002,
  CH-003: CH-003 is deprecated and the need above it is served directly
  (backlog 004).

## CH-021 Production code MUST anchor an entry in a Go doc comment, a TypeScript line comment, or the comment lines above a Python def

Status: active
Serves: CH-002

- A Go directive is a comment line `//charter:realizes ID, ...` or
  `//charter:concerns ID, ...`, with or without a space after the
  slashes, in the doc comment of a package clause or of a declaration;
  the anchor MUST record the kind, the file, the line and the declaration
  it documents.
- A TypeScript directive is the same text in a line comment of a file
  that is not a test; the anchor MUST record the file and the line.
- A Python directive is `# charter:realizes ID, ...` or
  `# charter:concerns ID, ...` on the comment lines directly above a def
  or a class, or above its decorators, in a file that is not a test; the
  anchor MUST record the file, the line and the declaration.
- A directive anywhere else MUST fail, naming where to move it, and MUST
  count for nothing.

Why: the anchor is the pointer from a decision to the code that carries
it, and the declaration it sits on is the place an agent is sent to.
Tying the directive to a declaration is what makes that destination
precise.

Scenarios:
1. `go-anchor-recorded`: given a directive on a package clause, when the
   tree is read, then it is recorded as the package's anchor with its
   declaration.
2. `spaced-directive-counts`: given a directive written with a space after
   the slashes, when the tree is read, then it is an anchor too.
3. `typescript-anchor-resolves`: given a TypeScript directive naming an
   unknown id, when the gate runs, then it fails naming the file and the
   line.
4. `comment-elsewhere-refused` (owner): given a directive in a comment that
   documents no declaration, inside a function or stray in a file, when
   the gate runs, then it fails naming the doc comment to move it to, and
   counts for nothing.
5. `python-anchor-recorded` (owner, backlog 003): given a Python directive
   above a def, or above its decorators, when the tree is read, then it is
   recorded as that def's anchor with its file and line.
6. `python-comment-elsewhere-refused` (owner, backlog 003): given a Python
   directive in a body, parted from its def by a blank line, or stray in a
   file, when the gate runs, then it fails and counts for nothing.

Changes:
- 2026-10-09: a directive outside a declaration's doc comment was read as
  an anchor on nothing; it is refused.
- 2026-11-26: Python is read: a directive on the comment lines directly
  above a def or a class, or its decorators, anchors it, and one anywhere
  else is refused (backlog 003).

## CH-261127-45KM charter MUST read an anchor from a comment line in a file no language reads, where charter.toml names that file

Status: active
Serves: CH-002
See: CH-021

- A `[[files]]` table in charter.toml names path patterns, matched as the
  walk's ignore patterns are, and the comment marker of the files they
  match, such as `#` for nftables, systemd units and shell scripts.
- A file the walk lists, whose extension is none of Go, Python or
  TypeScript, MUST be read by the first table whose patterns match it.
- In such a file, a line whose text after the marker and any spaces is
  `charter:realizes ID, ...` or `charter:concerns ID, ...` MUST be an
  anchor on the file, recording the kind, the file and the line.
- A directive naming no id MUST fail and count for nothing. A table with
  no pattern or no marker MUST fail naming its line.
- A file no table names MUST NOT be read for directives.

Why: some decisions are carried by a file that is not code, such as a
firewall ruleset or a service unit. The pointer from the decision belongs
on the file that decides, not on code beside it that merely loads it.

Scenarios:
1. `file-anchor-recorded` (agent, provisional): given a directive line in
   a file a table names, when the tree is read, then it is recorded as an
   anchor on that file, with its line.
2. `unnamed-file-unread` (agent, provisional): given the same line in a
   file no table names, when the tree is read, then it is not an anchor.
3. `empty-directive-refused` (agent, provisional): given a directive
   naming no id, when the gate runs, then it fails naming the file and the
   line.
4. `files-table-checked` (agent, provisional): given a `[[files]]` table
   with no pattern, or no comment marker, when the configuration is read,
   then it fails naming its line.

Guidance: name the narrowest patterns that hold the files, such as
`packaging/*/*.nft`; a broad pattern reads every matching file for
directives. The read is by line, so a directive inside a shell heredoc
or string counts as one. The walk skips directories whose names begin
with a dot or an underscore, so no pattern reaches them.

Changes:
- 2026-11-27: a file is read by the first matching table only, and only
  when its extension is no language's, which is what the walk does; a
  review found a file two tables match was read twice.

## CH-022 An anchor MUST name an entry that exists and is not planned, and a concerns anchor MUST NOT count toward coverage

Status: active
Serves: CH-002

- A directive naming an id that is not an entry MUST fail.
- A directive naming a planned entry MUST fail, saying to set the entry
  active as its work begins.
- A directive with no id MUST fail.
- A concerns anchor is navigation only: an active decision with only a
  concerns anchor MUST be treated as unrealized.

Why: an anchor to nothing points the reader at a decision that does not
exist, and no compiler catches it. A realizes anchor on a planned entry
is work that began without the record saying so.

Scenarios:
1. `unknown-id-refused`: given a directive naming an id that is not an
   entry, when the gate runs, then it fails naming the id.
2. `planned-refused`: given a realizes naming a planned entry, when the
   gate runs, then it fails.
3. `concerns-is-not-realized`: given an active decision anchored by
   concerns alone, when the gate runs, then it fails as realized by no
   code.

Changes:
- 2026-10-09: that a concerns anchor realizes nothing is ruled here
  alone, the concerns directive being this entry's subject: CH-030's
  not-realized was tested with a decision anchored by concerns alone, as
  concerns-is-not-realized is, and is now tested with one anchored by
  nothing, CH-030 pointing here.

## CH-023 A test MUST cite an entry by a verifies call that opens a test function, naming by handle the scenarios it exercises

Status: active
Serves: CH-001, CH-002

- A Go citation is a call of `Verifies` inside a function whose name
  starts with `Test`, whose second argument is the entry's symbol and
  whose further arguments are string literals naming the entry's
  scenarios.
- A TypeScript citation is a call of `verifies` inside a test call; a
  Python one is a call of `verifies` in a `test_` method of a
  `unittest.TestCase`, its first argument `self`, or in a pytest `test_`
  function, its first argument `None`; each has the same shape.
- The calls MUST open the test, or a subtest it runs: they are its first
  statements, one after another. A call that anything precedes or
  encloses, or in a function its runner never runs, such as Go's
  `Testx`, MUST fail.
- A TypeScript call in a test marked skip, fixme, fail or todo, or to a
  `verifies` the test file declares itself, MUST fail.
- A call anywhere else, an entry named by anything but its symbol, a
  handle that is not a literal, or a handle the entry does not list, MUST
  each fail naming the file and the line.
- Any other reference to an entry from a test file is a reference, not a
  citation.

Why: the gate reads the call statically, without running the test. A
call in a fixed place and shape is one it can read reliably and hold to
the scenarios it names, and requiring the call to open a test that
actually runs keeps a citation from claiming coverage the test never
exercises.

Scenarios:
1. `go-citation-read`: given the fixture, when it is read, then each
   `Verifies` call is a citation with its file, test function, line and
   handles.
2. `outside-a-test-refused`: given a `Verifies` in a helper function, when
   the gate runs, then it fails and its handles count for nothing.
3. `unknown-handle-refused`: given a handle the entry does not list, when
   the gate runs, then it fails naming the entry and the handle.
4. `literal-arguments`: given a handle that is not a string literal, or an
   entry named by a variable, when the gate runs, then it fails.
5. `typescript-citation-read`: given a `verifies` call in a TypeScript
   test, and one outside a test or with a non-literal handle, when the
   gate runs, then the first verifies the handles it names and the others
   fail.
6. `alias-followed`: given a citation through an import alias, when it is
   read, then it counts.
7. `reference-is-not-citation`: given a bare reference to an entry from a
   test file, when it is read, then it is recorded as a reference and
   counts for no scenario.
8. `later-statement-refused` (owner): given a `Verifies` that another
   statement precedes, or that a block encloses, and one opening a
   subtest's body, when the gate runs, then the first two fail and count
   for nothing and the third counts.
9. `non-test-name-refused` (owner): given a `Verifies` in a function named
   `Test` followed by a lower-case letter, which go test never runs, when
   the gate runs, then it fails naming the function.
10. `disabled-ts-refused` (owner): given a `verifies` in a TypeScript test
    marked skip, fixme or fail, when the gate runs, then it fails and
    counts for nothing.
11. `local-verifies-refused` (owner): given a `verifies` declared in the
    test file itself, when the gate runs, then it fails and counts for
    nothing.
12. `python-citation-read` (owner, backlog 003): given a `verifies`
    opening a unittest `test_` method, or a pytest `test_` function, when
    it is read, then it is a citation with its file, test, line and
    handles.
13. `python-later-statement-refused` (owner, backlog 003): given a Python
    `verifies` that another statement precedes, or in a method or
    function the runner does not collect, when the gate runs, then it
    fails and counts for nothing.

Changes:
- 2026-10-09: a citation counted wherever it sat in a test, after a skip,
  under a false branch, in a function go test never runs, in a skipped
  TypeScript test or through a verifies the test declared; each is refused.
- 2026-11-26: Python is read: a verifies opening a unittest test_ method
  or a pytest test_ function is a citation, held to the same refusals as
  Go's and TypeScript's (backlog 003).

## CH-024 Each language's verifies MUST record each citation and its test's outcome as a JSON line in the events file

Status: active
Serves: CH-001
Accounts: env:CHARTER_EVENTS

- When `CHARTER_EVENTS` names a file, the generated `verifies` of every
  language MUST append to it, once the test ends, one JSON line per
  citation: the id, the handle when one is named, the test, its file and
  its outcome (passed, failed or skipped).
- With the variable unset, it MUST record nothing, and the test runs as
  it would otherwise.
- Lines appended by tests running at the same time MUST each be whole.
- A nil or missing symbol MUST be refused at run time.
- The report reads these lines.

Why: the gate can only show that a test cites a scenario; the line shows
whether that test passed. One run then says which scenarios passed,
failed or were skipped, and one format written by every language means
the report needs one join, not one per test runner.

Scenarios:
1. `event-per-citation` (owner, backlog 003): given a Go, a Python and a
   TypeScript test each citing a scenario, when they run with
   `CHARTER_EVENTS` set, then each appends a line with the id, the handle,
   the test, the file and the outcome.
2. `outcome-recorded` (owner, backlog 003): given a passing, a failing and
   a skipped test, a skip after the citation included, when they run, then
   they record passed, failed and skipped.
3. `unset-records-nothing` (owner, backlog 003): given `CHARTER_EVENTS`
   unset, when a citing test runs, then it writes no file and passes as it
   would.
4. `concurrent-lines-whole` (owner, backlog 003): given tests running in
   parallel and appending at once, when they finish, then every line is a
   whole JSON object.
5. `nil-spec-refused`: given `Verifies` called with a nil entry, when the
   test runs, then it fails.

Changes:
- 2026-11-26: the marker is a JSON line per citation in the file
  CHARTER_EVENTS names, written by every language's generated verifies
  with the test's outcome, where Go's was a test attribute in go test's
  JSON stream; `attr-per-citation` is dropped for `event-per-citation`,
  and `outcome-recorded`, `unset-records-nothing` and
  `concurrent-lines-whole` are added (backlog 003).

## CH-025 The claims directory MUST hold only the record, charter.toml and generated files, and a test MUST NOT assign to an entry

Status: active
Serves: CH-261126-PFG5

- Beside the markdown of the record and charter.toml, the claims
  directory MUST hold only the files `charter gen` writes; a source file
  there that gen did not write MUST fail, saying entries and observations
  are written in markdown.
- A test file MUST NOT assign to an entry or a field of one, and a Go test
  MUST NOT dot-import the generated package.

Why: the generated files are the registry in each language, so anything
else written there would be an entry or a fact the area documents do not
know. An entry mutated by one test would mislead every test that runs
after it.

Scenarios:
1. `hand-written-source-refused` (owner, backlog 003): given a source file
   in the claims directory that `charter gen` did not write, when the gate
   runs, then it fails saying entries and observations are written in
   markdown.
2. `assignment-refused`: given an assignment to an entry, or to one of its
   fields, from a test file or from the package's own tests, when the gate
   runs, then it fails.
3. `dot-import-refused`: given a test file dot-importing the entries
   package, when the gate runs, then it fails.

Changes:
- 2026-11-26: observations are markdown (CH-261126-HK6S), so the claims
  directory holds no hand-written code at all: `functions-refused` and
  `spec-in-go-refused` are dropped for `hand-written-source-refused`
  (backlog 003).
- 2026-11-26: serves CH-261126-PFG5, where it served CH-003: CH-003 is
  deprecated and the need above it is served directly (backlog 004).

## CH-040 One charter command MUST generate, check, give ids and list anchors, so an adopter writes no main

Status: active
Serves: CH-002
See: CH-261126-33GK

- `charter gen` MUST write the generated files from the documents alone,
  so it works when one of them is stale or missing.
- `check` runs the gate (CH-037) and prints its warnings (CH-261127-QRQS);
  `report` and `diff` are CH-042's and CH-261009-QQC8's, and `curate` is
  CH-261127-GXP7's.
- `next` MUST print a new id for a prefix: the day's UTC date and a code
  drawn at random, drawn again when the tree already holds that id; or as
  many as the count after the prefix asks for, each different.
- `anchors` MUST print, per id: the document, the title, the status, the
  scenarios, the realizing and concerning sites with their declarations,
  the verifying tests with the handles each names, the relations, the
  docs that mention it (CH-261127-8MYR) and its Changes.
- Every verb MUST work on the charter CH-261126-33GK finds, and an adopter
  MUST NOT have to write or build anything of its own to run them.

Why: the agent's way into the code and the person's way to regenerate
files are the same command. A command an adopter has to write a main for
is one they will write wrong, or let fall behind the tooling.

Scenarios:
1. `gen-from-documents-alone` (owner, backlog 003): given one generated
   file missing or stale, when `charter gen` runs, then it writes every
   generated file, reading the documents alone.
2. `next-free`: given a prefix in use or not, when `next` runs, then it
   gives an id of the prefix, today's UTC date and a code of the alphabet,
   drawing again one the tree holds.
3. `next-count`: given a count, when `next` runs, then it gives that many
   ids, each different; given a count that is not a number above zero,
   then it refuses.
4. `verbs-on-the-found-charter` (owner, backlog 003): given an adopter's
   tree with no command of its own, when each verb runs inside it, then it
   reads the charter found there.

Changes:
- 2026-10-01: next gives dated ids, where it gave one above the highest,
  and takes a count, so a story's entries are numbered in one call
  (CH-012).
- 2026-10-09: the registry also gives the area documents, in the
  reading order with their headers, and the waivers, so a renderer such
  as the claims viewer walks a reader through any charter in the order
  it declares and shows a scenario as waived without parsing the files
  again.
- 2026-10-09: the next free id is ruled here alone, the next verb's
  behaviour being this entry's subject: CH-012's next-free restated
  next-free and is dropped, its test citing this one, which now says what
  the dropped one did, the UTC date, a prefix in use or not and the
  redraw, and the text says how next draws.
- 2026-11-26: the charter became a repository of its own, meant for other repositories (Q-261015-3Q84). The verbs live in the
  cli package, which a repository's own command calls with the charters
  it keeps; the charter flag selects one of those by name, where it chose
  between the product's charter and the charter's own. The title says any
  charter, and `charter-flag-selects` is added. A command reads a charter
  from the module above the working directory, since the calling file is
  the tooling's, in another module; `root-from-working-directory` is
  added.
- 2026-11-26: the verbs are one charter command, the gate's check and
  claimsmap's export among them, working on the charter it finds above the
  working directory. An adopter writes no main, so the charter flag and a
  root given in options are gone, `charter-flag-selects` and
  `root-from-working-directory` with them; `generate-writes-both` becomes
  `gen-from-documents-alone`, and `verbs-on-the-found-charter` is added
  (backlog 003).
- 2026-11-27: curate is a verb, CH-261127-GXP7's, and anchors lists the
  docs that mention an entry, CH-261127-8MYR's (backlog 007).
- 2026-11-27: lint and export are gone, the owner ruling the viewer, its
  export and lint's baggage deleted; check prints the warnings lint kept.
  `registry-lists-everything` goes with export, which was its one reader,
  and anchors gives an entry's Changes, which export alone gave.

## CH-261127-3NY3 A doc MUST be a markdown file charter.toml names, MUST declare no entry, and its references to entries MUST resolve

Status: active
Serves: CH-002
See: CH-261126-33GK, CH-261127-DW8R, CH-261127-H4B7

- A doc is a markdown file that charter.toml's doc patterns match,
  matched as the ignore patterns are; charter MUST NOT read any other file
  as a doc.
- A doc tells a story over the record and the code, and its prose is
  held to nothing.
- A doc references an entry by a markdown link whose text opens with the
  entry's id, and embeds one between a `charter:embed ID` comment and its
  closing `/charter:embed` comment.
- A reference or an embed naming no entry, or a deprecated one, MUST fail
  naming the doc and the line.
- A heading in a doc that opens with an id MUST fail, naming the area
  documents as where entries are declared.

Why: the record holds what binds, but a reader also needs to know how the
system works, why it is shaped as it is, and a story that joins many
entries. That prose used to live outside the gate and could keep
pointing at things that were gone without anyone being told. A doc whose
pointers are checked can still be wrong, but it cannot point at nothing,
and keeping entries declared only in area documents keeps the record
easy to read by program.

Scenarios:
1. `doc-by-pattern` (owner, backlog 007): given a markdown file the doc
   patterns match and one they do not, when the tree is read, then only
   the first is read as a doc.
2. `reference-resolves` (owner, backlog 007): given a link whose text
   opens with an active entry's id and one naming no entry, when the gate
   runs, then the first passes and the second fails naming the doc and the
   line.
3. `deprecated-reference-fails` (owner, backlog 007): given a reference to
   a deprecated entry, when the gate runs, then it fails naming it.
4. `embed-resolves` (owner, backlog 007): given an embed naming no entry,
   when the gate runs, then it fails naming the doc and the line.
5. `doc-declares-nothing` (owner, backlog 007): given a heading in a doc
   that opens with an id, when the gate runs, then it fails and no entry is
   read from it.

Guidance: the doc patterns are their own scope: charter.toml's only and
ignore patterns do not narrow them, and a matched file that is not
markdown is not a doc. As in every walk, a directory whose name begins
with a dot or an underscore is skipped, so no doc lives there.

Changes:
- 2026-11-27: the guidance says the doc patterns are not narrowed by only
  or ignore, and read markdown alone; the build read docs through the
  walk's only patterns, which hid a doc the patterns named (review of
  backlog 007).

## CH-261127-DW8R A doc MUST reference code by a link to the file that declares the named symbol, and the gate MUST fail one that file lacks

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261126-PQ7X

- A code reference in a doc is a markdown link whose text is a name in
  backticks and whose target is a file under the root.
- In a file a language reads, the name MUST be a declaration in that
  file: a function, type, class, method or constant, written bare or
  qualified by its package or its type, as `Accept`, `engine.Accept` or
  `Engine.Accept`. In any other file, the file's existing is enough.
- A target that is missing, or a file that declares no such name, MUST
  fail naming the doc, the line and the name.
- A link whose text is not in backticks, or whose target is outside the
  root, is prose and MUST NOT be read.

Why: how a system works lives in its code, so a walkthrough is a guide to
declarations, and one that names a function since renamed sends the
reader nowhere. Resolving the name in the file that each language's
reader already parses makes a rename fail the gate, where the prose would
otherwise go wrong silently.

Scenarios:
1. `go-symbol-resolves` (owner, backlog 007): given a reference to a Go
   file naming a function it declares, bare or qualified by its package,
   and one naming an undeclared function, when the gate runs, then the
   first passes and the second fails naming the doc, line and name.
2. `python-symbol-resolves` (owner, backlog 007): given a reference to a
   Python file naming a class it declares and one naming a name it lacks,
   when the gate runs, then the first passes and the second fails.
3. `typescript-symbol-resolves` (owner, backlog 007): given a reference to
   a TypeScript file naming a function it exports and one naming a name it
   lacks, when the gate runs, then the first passes and the second fails.
4. `method-by-type` (owner, Q-261127-BVQX): given a name qualified by a
   type, when it is resolved, then it resolves to that type's method, and
   to nothing when the type has no such method.
5. `plain-file-exists` (owner, backlog 007): given a reference to a file
   no language reads, when the gate runs, then it passes when the file
   exists and fails when the target is missing.
6. `prose-link-unread` (owner, backlog 007): given a link whose text is
   not in backticks, when the doc is read, then it is not read as a code
   reference.

Changes:
- 2026-11-27: `method-by-type` is the owner's, answering Q-261127-BVQX
  yes: a type qualifies its method, and in Python and TypeScript the
  file's name stands for the package.

## CH-261127-H4B7 charter gen MUST keep each doc's references and embeds current with the record and the code

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261127-DW8R, CH-040

- In every doc, `charter gen` MUST rewrite: the text after the id of each
  entry reference to the entry's current title, leaving a bare id bare;
  its target to the entry's heading in its area document; the target of
  each code reference to the file, with a line anchor at the declaration;
  and what lies between each embed's two comments to the entry's id,
  title and text, taken from the area document.
- gen MUST change nothing else in a doc, and a second run with nothing
  changed MUST write nothing.

Why: a doc must read well where people read it, on a code host that
renders markdown as written and transcludes nothing. Because gen writes
the titles, targets and embedded text, a reference keeps its title and an
embed its text without anyone keeping a copy up to date by hand.

Scenarios:
1. `title-kept-current` (owner, backlog 007): given an entry whose title
   changed, when gen runs, then each reference's text after the id is the
   new title.
2. `target-kept-current` (owner, backlog 007): given an entry moved to
   another area document, when gen runs, then each reference's target is
   its heading there.
3. `line-kept-current` (owner, backlog 007): given a declaration moved down
   its file, when gen runs, then the code reference's line anchor is the
   declaration's new line.
4. `embed-written` (owner, backlog 007): given an embed, when gen runs,
   then the entry's id, title and text are written between its comments,
   replacing what was there.
5. `prose-untouched` (owner, backlog 007): given a doc, when gen runs
   twice, then its prose is byte for byte as it was and the second run
   writes nothing.
6. `bare-id-stays` (owner, Q-261127-TW9H): given a reference whose text is
   the id alone, when gen runs, then the text stays the id alone and its
   target is made current.

Guidance: an embed renders as a quote: the id, linked to the entry, and
the title in bold, then the text.

Changes:
- 2026-11-27: a reference whose text is the id alone stays so, a link in
  running prose; the guidance says the embed's link is on its id, where it
  said after the text (build of backlog 007).
- 2026-11-27: the bare id is the agent's choice, not a ruling, so its
  clause has a scenario of its own, provisional, and Q-261127-TW9H asks
  the owner (review of backlog 007).
- 2026-11-27: `bare-id-stays` is the owner's, answering Q-261127-TW9H
  yes.

## CH-261127-GXP7 charter curate MUST record who curated a doc and a hash of each reference, and charter check MUST warn when one changes

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261127-DW8R, CH-261127-QRQS

- `charter curate FILE` MUST replace the comment block at the end of the
  doc with one naming the curator, the date, and a hash of each entry and
  code reference the doc holds.
- An entry's hash MUST cover its title, status and text, whitespace
  normalized; a new Changes line alone MUST NOT change it.
- A code reference's hash MUST cover its declaration's source, whitespace
  normalized, or the whole file for a file no language reads.
- curate MUST refuse a doc with an unresolved reference, and MUST NOT
  change the doc outside the block.
- A doc with no block is drafted; one whose recorded hash differs from the
  current one is stale. `charter check` MUST warn of each, naming the
  changed references, and MUST NOT fail on either.

Why: an agent drafts a doc and a person curates it, and what the person
checked is what the doc rested on when they read it. A curation date
compared against git history would need that history intact, which a
shallow clone or a squashed merge loses. A hash stored in the doc travels
with it and changes only when the referenced text does, so check can say
which prose may no longer be true without claiming to know whether it is.

Scenarios:
1. `curate-records` (owner, backlog 007): given a doc, when it is curated
   and curated again, then it holds one block with the curator, the date
   and a hash per reference, the second replacing the first.
2. `stale-on-entry-change` (owner, backlog 007): given a curated doc that
   references an entry, when the entry's text changes, then check warns
   that the doc is stale and names that entry.
3. `changes-line-not-stale` (owner, backlog 007): given a curated doc that
   references an entry, when the entry gains only a Changes line, then
   check does not warn of the doc.
4. `stale-on-go-change` (owner, backlog 007): given a curated doc that
   references a Go function, when the function's body changes, then check
   warns of the doc with that reference; when only its whitespace changes,
   then it does not.
5. `stale-on-python-change` (owner, backlog 007): given a curated doc that
   references a Python function, when the function's body changes, then
   check warns of the doc with that reference.
6. `stale-on-typescript-change` (owner, backlog 007): given a curated doc
   that references a TypeScript function, when the function's body
   changes, then check warns of the doc with that reference.
7. `drafted-named` (owner, backlog 007): given a doc with no curation
   block, when check runs, then it warns of the doc as drafted.
8. `stale-never-fails` (owner, backlog 007): given a stale doc and a
   drafted one, when check runs, then it passes.
9. `curate-refuses-broken` (owner, backlog 007): given a doc with a
   reference that does not resolve, when it is curated, then curate fails
   naming the reference and writes nothing.
10. `curate-touches-only-its-block` (owner, backlog 007): given a doc,
    when it is curated, then every byte before its block is as it was.

Changes:
- 2026-11-27: check warns of a drafted or stale doc, where lint named
  it: the owner ruled lint deleted and its three useful checks kept as
  check's warnings (CH-261127-QRQS).

## CH-261127-8MYR charter MUST compute which docs mention an entry from the docs, and charter anchors and charter diff MUST show it

Status: active
Serves: CH-002, CH-001
See: CH-040, CH-261009-QQC8, CH-261127-GXP7

- Nobody writes a backlink: the docs that reference or embed an entry
  MUST be read from the docs themselves.
- `charter anchors` MUST list them for each id, with the line of each
  mention.
- `charter diff`, given a revision, MUST list after what it names each doc
  holding a reference whose hash differs from the revision's, as docs to
  read again, and MUST NOT change its exit status for them.

Why: a person changing an entry needs to know which prose rests on it,
and a reviewer of a change needs to know which docs to read again. Links
written by hand in both directions fall out of step, and the docs already
say everything the backlinks need.

Scenarios:
1. `anchors-lists-docs` (owner, backlog 007): given an id two docs
   mention, when `charter anchors` runs for it, then it lists both docs
   with the line of each mention.
2. `diff-lists-docs` (owner, backlog 007): given a doc whose referenced
   entry's text changed since a revision, when `charter diff` runs against
   it, then it lists the doc and exits as it would without the doc.

Changes:
- 2026-11-27: the export and the viewer are deleted, by the owner's
  ruling, so neither shows the docs; `export-carries-docs` goes with them.
