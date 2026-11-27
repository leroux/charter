# From an id to its code and its tests

Code: internal/lang/golang, internal/lang/python, internal/lang/typescript, internal/lang/plain, internal/facts, internal/docs, cmd/charter

The problem this area answers: a person or an agent new to the code cannot
find their way from what was decided to where it is done. It holds how code
and tests name the entries in each language, the generated symbols that
make a wrong name fail, and the command that lists an entry's code and
tests, and the docs that tell a story over them.

## CH-020 Every entry has a generated symbol per language read, its id with underscores for hyphens, its title the doc comment

Status: active
Serves: CH-002, CH-261126-PFG5
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
- 2026-11-26: serves CH-002, CH-261126-PFG5, where it served CH-002,
  CH-003: CH-003 is deprecated and the need above it is served directly
  (backlog 004).

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

## CH-261127-45KM A file no language reads anchors an entry by a comment line, where charter.toml names it

Status: active
Serves: CH-002
See: CH-021

A `[[files]]` table in charter.toml names path patterns, matched as the
walk's ignore patterns are, and the comment marker of the files they
match, such as `#` for nftables, systemd units and shell scripts. A file
the walk lists, whose extension is none of Go, Python or TypeScript, is
read by the first table whose patterns match it. In such a file, a line
whose text after the marker and any spaces is "charter:realizes ID, ..."
or "charter:concerns ID, ..." is an anchor on the file; the anchor
records the kind, the file and the line. A directive naming no id fails
and counts for nothing. A table with no pattern or no marker fails
naming its line, and a file no table names is not read for directives.

Why: some decisions are carried by a file that is not code, a firewall
ruleset or a service unit, and the pointer from the decision belongs on
the file that decides, not on code beside it that merely loads it.

Scenarios:
1. `file-anchor-recorded` (agent, provisional): a directive line in a file
   a table names is recorded as an anchor on that file, with its line.
2. `unnamed-file-unread` (agent, provisional): the same line in a file no
   table names is not an anchor.
3. `empty-directive-refused` (agent, provisional): a directive naming no
   id fails naming the file and the line.
4. `files-table-checked` (agent, provisional): a `[[files]]` table with no
   pattern, or no comment marker, fails naming its line.

Guidance: name the narrowest patterns that hold the files, such as
`packaging/*/*.nft`; a broad pattern reads every matching file for
directives. The read is by line, so a directive inside a shell heredoc
or string counts as one. The walk skips directories whose names begin
with a dot or an underscore, so no pattern reaches them.

Changes:
- 2026-11-27: a file is read by the first matching table only, and only
  when its extension is no language's, which is what the walk does; a
  review found a file two tables match was read twice.

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
Serves: CH-261126-PFG5

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
- 2026-11-26: serves CH-261126-PFG5, where it served CH-003: CH-003 is
  deprecated and the need above it is served directly (backlog 004).

## CH-040 One charter command generates, checks, lints, gives ids and lists anchors, and an adopter writes no main

Status: active
Serves: CH-002
See: CH-261126-33GK

charter gen writes the generated files from the documents alone, so it
works when one is stale or missing. check runs the gate (CH-037), lint
prints the gate's problems and then the advice, report and diff are
CH-042's and CH-261009-QQC8's, and export is CH-261015-9ZTY's, curate CH-261127-GXP7's. next prints
a new id for a prefix, the day's UTC date and a code drawn at random,
drawn again when the tree already holds that id, or as many as the count
after the prefix asks for, each different. anchors prints, per id, the
document, the title, the status, the scenarios, the realizing and
concerning sites with their declarations, the verifying tests with the
handles each names, the relations and the docs that mention it
(CH-261127-8MYR). Every verb works on the charter
CH-261126-33GK finds, and an adopter writes or builds nothing of its own
to run them.

Why: the agent's way in and the person's way to regenerate are one
command; a command an adopter must write a main for is one they write
wrong, or let fall behind the tooling.

Scenarios:
1. `gen-from-documents-alone` (owner, backlog 003): charter gen writes
   every generated file when one of them is missing or stale, reading the
   documents alone.
2. `registry-lists-everything`: the registry gives each entry with
   scenarios, anchors and citations by handle, each area document in
   reading order with title, code and summary, and each waiver with its
   reason.
3. `next-free`: next gives an id of the prefix, today's UTC date and a
   code of the alphabet, for a prefix in use or not, drawing again one the
   tree holds.
4. `next-count`: next with a count gives that many ids, each different,
   and a count that is not a number above zero is refused.
5. `verbs-on-the-found-charter` (owner, backlog 003): run inside an
   adopter's tree, each verb reads the charter found there, with no
   command of the adopter's own.

Changes:
- 2026-10-01: next gives dated ids, where it gave one above the highest,
  and takes a count, so a story's entries are numbered in one call
  (CH-012).
- 2026-10-09: the registry also gives the area documents, in the
  reading order with their headers, and the waivers, so a renderer such
  as the claims viewer walks a reader through any charter in the order
  it declares and shows a scenario as waived without parsing the files
  again.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-10-09: the next free id is ruled here alone, the next verb's
  behaviour being this entry's subject: CH-012's next-free restated
  next-free and is dropped, its test citing this one, which now says what
  the dropped one did, the UTC date, a prefix in use or not and the
  redraw, and the text says how next draws.
- 2026-10-09: `registry-lists-everything` reworded under thirty words; the
  ruling is unchanged.
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

## CH-261127-3NY3 A doc is a markdown file charter.toml names; it declares no entry, and its references to entries must resolve

Status: active
Serves: CH-002
See: CH-261126-33GK, CH-261127-DW8R, CH-261127-H4B7

A doc is a markdown file that charter.toml's doc patterns match, matched
as the ignore patterns are; no other file is read as one. A doc tells a
story over the record and the code, and its prose is held to nothing. It
references an entry by a markdown link whose text opens with the entry's
id, and embeds one between a `charter:embed ID` comment and its closing
`/charter:embed` comment. A reference or an embed naming no entry, or a
deprecated one, fails naming the doc and the line. A heading in a doc
that opens with an id fails naming the area documents, where entries are
declared.

Why: the record holds what binds, but a reader also needs how the system
works, why it is shaped so, and a story joining many entries; that prose
lives outside the gate today and points at things that are gone without
anyone being told. A doc whose pointers are checked can still be wrong,
but it cannot point at nothing, and an entry declared in one place alone
stays easy to read by program.

Scenarios:
1. `doc-by-pattern` (owner, backlog 007): a markdown file the doc patterns
   match is read as a doc, and one they do not match is not.
2. `reference-resolves` (owner, backlog 007): a link whose text opens
   with an active entry's id passes, and one naming no entry fails naming
   the doc and the line.
3. `deprecated-reference-fails` (owner, backlog 007): a reference to a
   deprecated entry fails naming it.
4. `embed-resolves` (owner, backlog 007): an embed naming no entry fails
   naming the doc and the line.
5. `doc-declares-nothing` (owner, backlog 007): a heading in a doc that
   opens with an id fails, and no entry is read from it.

Guidance: the doc patterns are their own scope: charter.toml's only and
ignore patterns do not narrow them, and a matched file that is not
markdown is not a doc. As in every walk, a directory whose name begins
with a dot or an underscore is skipped, so no doc lives there.

Changes:
- 2026-11-27: the guidance says the doc patterns are not narrowed by only
  or ignore, and read markdown alone; the build read docs through the
  walk's only patterns, which hid a doc the patterns named (review of
  backlog 007).

## CH-261127-DW8R A doc references code by a link to the file declaring the named symbol, and one that file lacks fails

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261126-PQ7X

A code reference in a doc is a markdown link whose text is a name in
backticks and whose target is a file under the root. In a file a
language reads, the name is a declaration in that file, a function, type,
class, method or constant, written bare or qualified by its package or its
type, as `Accept`, `engine.Accept` or `Engine.Accept`; in any other file,
the file's existing is enough. A target that is missing, or a file that
declares no such name, fails naming the doc, the line and the name. A
link whose text is not in backticks, or whose target is outside the root,
is prose and not read.

Why: how a system works lives in its code, so a walkthrough is a guide to
declarations, and one that names a function since renamed sends the
reader nowhere; reading the declaration from the file each language's
reader already parses makes a rename fail the gate where the prose would
silently go wrong.

Scenarios:
1. `go-symbol-resolves` (owner, backlog 007): a reference to a Go file
   naming a function it declares passes, bare or qualified by its
   package; one naming an undeclared function fails naming the doc, line
   and name.
2. `python-symbol-resolves` (owner, backlog 007): a reference to a Python
   file naming a class it declares passes, and a name it lacks fails.
3. `typescript-symbol-resolves` (owner, backlog 007): a reference to a
   TypeScript file naming a function it exports passes, and a name it
   lacks fails.
4. `method-by-type` (agent, provisional): a name qualified by a type
   resolves to that type's method, and to nothing when the type has no
   such method.
5. `plain-file-exists` (owner, backlog 007): a reference to a file no
   language reads passes when the file exists, and a missing target fails.
6. `prose-link-unread` (owner, backlog 007): a link whose text is not in
   backticks is not read as a code reference.

## CH-261127-H4B7 charter gen keeps each doc's references and embeds current with the record and the code

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261127-DW8R, CH-040

charter gen rewrites, in every doc, the text after the id of each entry
reference to the entry's current title, leaving a bare id bare, and its
target to the entry's heading in its area document; the target of each
code reference to the file with a line anchor at the declaration; and
what lies between each embed's two comments to the entry's id, title and
text, written from the area document. It changes nothing else in a doc, and a second run with
nothing changed writes nothing.

Why: a doc must read well where people read it, on a code host that
renders markdown as written and transcludes nothing; written out by gen,
a reference keeps its title and an embed its text without a person ever
holding a copy to keep up to date by hand.

Scenarios:
1. `title-kept-current` (owner, backlog 007): after an entry's title
   changes, gen rewrites each reference's text after the id to the new
   title.
2. `target-kept-current` (owner, backlog 007): after an entry moves to
   another area document, gen rewrites each reference's target to its
   heading there.
3. `line-kept-current` (owner, backlog 007): after a declaration moves
   down its file, gen rewrites the code reference's line anchor to the
   declaration's new line.
4. `embed-written` (owner, backlog 007): between an embed's comments gen
   writes the entry's id, title and text, replacing what was there.
5. `prose-untouched` (owner, backlog 007): gen leaves a doc's prose
   byte for byte as it was, and a second run writes nothing.
6. `bare-id-stays` (agent, provisional): a reference whose text is the id
   alone keeps it alone, its target made current.

Guidance: an embed renders as a quote: the id, linked to the entry, and
the title in bold, then the text.

Changes:
- 2026-11-27: a reference whose text is the id alone stays so, a link in
  running prose; the guidance says the embed's link is on its id, where it
  said after the text (build of backlog 007).
- 2026-11-27: the bare id is the agent's choice, not a ruling, so its
  clause has a scenario of its own, provisional, and Q-261127-TW9H asks
  the owner (review of backlog 007).

## CH-261127-GXP7 charter curate records who curated a doc and a hash of what it references; lint names one whose references changed

Status: active
Serves: CH-002
See: CH-261127-3NY3, CH-261127-DW8R, CH-041

charter curate FILE writes, at the end of the doc, a comment block naming
who curated it, the date, and a hash of each entry and each code
reference the doc holds as they are now, replacing the block that was
there. An entry's hash is of its title, status and text, whitespace
normalized, so a new Changes line alone leaves it unchanged; a code
reference's hash is of its declaration's source, whitespace normalized,
or of the whole file for a file no language reads. Curate refuses a doc
whose references do not resolve, and changes nothing of it outside the
block. A doc is stale when a reference's hash now differs from the one
recorded, and drafted when it has no block; lint names each such doc,
with each changed reference, and the gate never fails on either.

Why: an agent drafts a doc and a person curates it, and what the person
checked is what the doc rested on when they read it. A date compared with
history needs that history whole, which a shallow clone or a squashed
merge loses; a hash in the doc travels with it and changes only when the
referenced text does, so lint can say which prose may no longer be true
without claiming to know whether it is.

Scenarios:
1. `curate-records` (owner, backlog 007): curating a doc writes a block
   with the curator, the date and a hash per reference, and curating again
   replaces it.
2. `stale-on-entry-change` (owner, backlog 007): after a referenced
   entry's text changes, lint names the doc as stale with that entry.
3. `changes-line-not-stale` (owner, backlog 007): after a referenced
   entry gains only a Changes line, lint does not name the doc.
4. `stale-on-go-change` (owner, backlog 007): after a referenced Go
   function's body changes, lint names the doc with that reference, and
   a change of whitespace alone does not.
5. `stale-on-python-change` (owner, backlog 007): after a referenced
   Python function's body changes, lint names the doc with that reference.
6. `stale-on-typescript-change` (owner, backlog 007): after a referenced
   TypeScript function's body changes, lint names the doc with that
   reference.
7. `drafted-named` (owner, backlog 007): a doc with no curation block is
   named by lint as drafted.
8. `stale-never-fails` (owner, backlog 007): a stale doc and a drafted
   one leave charter check passing.
9. `curate-refuses-broken` (owner, backlog 007): curating a doc with a
   reference that does not resolve fails naming it, and writes nothing.
10. `curate-touches-only-its-block` (owner, backlog 007): curating leaves
    every byte of the doc before its block as it was.

## CH-261127-8MYR Which docs mention an entry is computed from the docs, and anchors, export and diff show it

Status: active
Serves: CH-002, CH-001
See: CH-040, CH-261015-9ZTY, CH-261009-QQC8, CH-261127-GXP7

Nobody writes a backlink: the docs that reference or embed an entry are
read from the docs. charter anchors lists them for each id, with the line
of each mention; the export carries them with each entry, and the viewer
shows them on its page; and charter diff, given a revision, lists after
what it names each doc holding a reference whose hash differs from the
revision's, as docs to read again, without changing its exit status.

Why: a person changing an entry needs to know which prose rests on it,
and a reviewer of a change which docs to read again; links written by
hand in both directions fall out of step, and the docs already say
everything the backlinks need.

Scenarios:
1. `anchors-lists-docs` (owner, backlog 007): anchors for an id that two
   docs mention lists both, with the line of each mention.
2. `export-carries-docs` (owner, backlog 007): the export carries, with
   each entry, the docs that mention it.
3. `diff-lists-docs` (owner, backlog 007): diff against a revision lists
   a doc whose referenced entry's text changed since, and exits as it
   would without the doc.
