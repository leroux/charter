# Recording intent

Code: internal/record

This area answers the problem that what the software must do has no
durable home. It covers the form intent is written in: the entry, its id,
its status and relations, the needs, the observations, and what keeps the
record stable.

## CH-010 The record reader MUST read an entry as a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes, in that order

Status: active
Serves: CH-261126-PFG5, CH-001
See: CH-013

- An area document is a markdown file whose header is a title and an
  optional Code line, and whose entries are headings of the form
  `## ID Title`. What a header may hold is CH-261126-JC9B's.
- The lines after the heading, up to a blank line, MUST be fields:
  Status, Kind, Serves, See and Accounts, each at most once. The values a
  Status takes are CH-013's.
- The body MUST be paragraphs: the text first, then paragraphs opened by
  Why, Scenarios, Guidance and Changes, in that order and each at most
  once. The text and Why MUST be present, and an active entry MUST list
  at least one scenario.
- A scenario MUST be a numbered item with a handle in backticks,
  lower-case kebab and unique in the entry, followed by a sentence. A
  change MUST be a dated item. A continuation line is indented.
- Every form problem MUST name the document, the line and the form to
  write.

Why: the parts are what a reviewer reads, in the order they read them.
A fixed order and a fixed form make a diff of the record say exactly
which part changed, which is what lets review work from the record
instead of the code.

Scenarios:
1. `every-part-read`: given a well-formed entry, when it is read, then it
   yields its id, symbol, title, status, kind, relations, accounts, text,
   why, scenarios with their handles and lines, guidance and dated
   changes.
2. `heading-form`: given a heading without an id and a title, when it is
   read, then it is a form problem.
3. `fields-checked`: given an unknown field, a field twice, a kind outside
   the three or a malformed surface key, when it is read, then each is a
   form problem.
4. `required-parts`: given an entry without a Status, without text or
   without a Why, or an active entry without Scenarios, when it is read,
   then it is a form problem.
5. `scenario-form`: given a scenario item that is not "N. `handle`:
   sentence", a handle that is not lower-case kebab, or a handle used
   twice, when it is read, then each is a form problem.
6. `parts-in-order`: given a part out of order, a part twice, prose after
   a labelled part or a heading inside an entry, when it is read, then
   each is a form problem.
7. `change-form`: given a change that is not "- YYYY-MM-DD: text", or
   whose date is not on the calendar, when it is read, then it is a form
   problem.
8. `problems-locate`: given any form problem, when it is reported, then
   it names the document, the line and the rule.
9. `blank-lines-before-fields`: given blank lines between the heading and
   the first field, when the entry is read, then they are allowed, and
   the first blank line after a field ends the fields.

Changes:
- 2026-10-01: the header may hold prose, lists, tables and diagrams, is
  not parsed beyond its title and Code line and has no section headings;
  a line inside a fenced block is never a heading, which the parser had
  refused inside an entry and split entries on.
- 2026-10-01: a fence among the field lines no longer lets a "## " line
  in it start an entry; that the header carries no rulings and no index
  moved to Guidance, being review's to hold (conformance review).
- 2026-10-09: the reader skips audit.md, the audit's report, as it skips
  README.md, the digest, questions.md, policy.md and glossary.md: none is
  an area document.
- 2026-10-09: a status outside the three, which fields-checked named
  beside CH-013's status-form, is CH-013's alone, the status form being
  its subject; fields-checked no longer names it, and the text names
  CH-013 for a Status's values, keeping the fields' presence and order
  here.
- 2026-11-26: serves CH-261126-PFG5, CH-001, where it served CH-001 and
  CH-002: the needs are re-cut one a problem, and this decision answers
  the one that names it (backlog 004).
- 2026-11-26: split: the header of an area document, what it may hold and
  that a fenced line is never a heading, is CH-261126-JC9B's, a ruling of
  its own; readme-is-not-a-document, header-is-prose and
  heading-in-fence-is-text and the Guidance moved there with it (backlog
  004).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-JC9B The record reader MUST treat an area document's header as a title, a Code line and prose, and MUST NOT read a fenced line as a heading

Status: active
Serves: CH-261126-PFG5
See: CH-010

- The header of an area document, before its first entry, is a title, an
  optional Code line, and prose that orients the reader: paragraphs, a
  vocabulary list, tables and diagrams.
- The reader MUST NOT parse the header beyond its title and Code line,
  and the header MUST NOT have section headings.
- A line inside a fenced block, in the header or in an entry, MUST be
  read as text and never as a heading.
- README.md in the directory is documentation, not an area document.

Why: a reader needs a way into an area before its entries, and a header
that the parser held to a form would leave no room for one. A fenced
example of an entry must stay an example, so it must not be read as an
entry.

Scenarios:
1. `readme-is-not-a-document`: given README.md in the claims directory,
   when the record is read, then it is not read as an area document.
2. `header-is-prose`: given a header of prose, a list, a table and a
   fenced block with a "#" line, when it is read, then it parses clean and
   yields no entry; given a section heading there, then it is a form
   problem.
3. `heading-in-fence-is-text`: given a fenced "#" or "## " line in an
   entry's body, when it is read, then it is text, no form problem and no
   entry; given a fence among the field lines, then it is no field line
   and starts no entry.

Guidance: the header carries no rulings, and a hand-written index of the
entries is not part of it, since the headings are the index; neither is
something the parser can tell, so review holds the header to them.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-011 The record reader MUST treat an entry of needs.md as a need, with no Kind, Serves or Guidance

Status: active
Serves: CH-261126-PFG5

An entry of the document named needs.md is a need. A need MUST NOT
serve anything, carry a kind or carry guidance, and the reader MUST
report each of those as a form problem.

Why: a need is what a person can rely on. It is realized by the
decisions that serve it, so it has nothing to serve itself, and how it
is met belongs in those decisions rather than in the need.

Scenarios:
1. `need-by-document`: given an entry in needs.md and one in another
   document, when they are read, then the first is a need and the second
   a decision.
2. `need-fields-refused`: given a need with a Kind, a Serves or a
   Guidance, when it is read, then each is a form problem.

Changes:
- 2026-11-26: serves CH-261126-PFG5, where it served CH-001: the needs are
  re-cut one a problem, and this decision answers the one that names it
  (backlog 004).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-012 An id MUST be declared once, and a new id MUST be its prefix, the day it was taken and a random code

Status: active
Serves: CH-261126-PFG5
See: CH-040

- An id MUST be a prefix followed by either a number, as SW-118, or the
  day it was taken and a code, as SW-261001-K7Q2; a heading whose id is
  neither is a form problem.
- The day is the UTC date as six digits, year, month and day. The code
  is four characters of BCDFGHJKMNPQRSTVWXYZ23456789, the consonants and
  digits that spell no word and are not read as one another.
- The numbered ids stand as they are. `charter next` gives every new id,
  of the second kind and not already in the tree, as CH-040 rules.
- Ids MUST be unique across the area documents. A second declaration
  MUST fail, naming both places and saying that the entry that has not
  landed takes a new id.
- The entries of the area documents are ordered by prefix, and within one
  prefix the numbered ids first, by number, then the dated ones by date
  and code; every generated file MUST declare them in that order.

Why: ids need no state file and no counter to share, because the tree is
the state. Numbering one above the highest gave the same number on every
branch cut from one tree, so branches written at the same time collided
at merge, and the later one had to rewrite its id in every document,
anchor and citation that named it. A code drawn from about six hundred
thousand is unlikely to repeat on two branches, and the gate still names
both places on the day it does. The date adds nothing to uniqueness; it
is there so that ids read in the order they were taken.

Scenarios:
1. `declared-once`: given an id declared in two documents, when the gate
   runs, then it fails naming both and saying the entry not yet landed
   takes a new id.
2. `id-form`: given numbered and dated ids, when they are read, then they
   parse; given an off-calendar date or a code that is not four
   characters of the alphabet, then it is a form problem.
3. `ordered-by-date`: given area entries of several prefixes and both id
   kinds, when the generated files are written, then every one declares
   them by prefix, and within one numbered ids by number before dated ones
   by date and code.

Changes:
- 2026-10-01: a new id is the day it was taken and a random code, where
  it was one above the highest: branches written at the same time took
  the same number, and the later one rewrote it everywhere it was
  written. The numbered ids stand, and the declared-twice problem says
  which entry takes a new id.
- 2026-10-09: the next free id moved to CH-040, the next verb's behaviour
  being its subject: this entry's next-free restated CH-040's and is
  dropped, its test citing CH-040's, and the text names CH-040 for how
  next gives a new id; the id's form stays here.
- 2026-10-09: that a dated id names no spec file moved to CH-014, a spec
  file's name being its subject: id-form stated it beside CH-014's
  file-name-checked and no longer does, and the text names CH-014 for
  which ids may name a spec file.
- 2026-10-09: the order ids sort in is ruled here alone, the order the
  generated files declare them in with it: CH-020's bare-symbol-declared
  stated it too and now points here, and the text says the generated Go
  and TypeScript module follow it.
- 2026-11-26: spec files are no longer read (Q-261015-N4HJ), so ids are
  unique across the area documents alone and `declared-once` names two
  documents; the order holds in every generated file, Python's among them
  (backlog 003).
- 2026-11-26: serves CH-261126-PFG5, where it served CH-003: CH-003 is
  deprecated and the need above it is served directly (backlog 004).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-013 An entry's Status MUST be planned, active or deprecated; every status MUST generate a symbol, and the gate MUST ask coverage of active entries only

Status: active
Serves: CH-002, CH-261126-PFG5
See: CH-020

- An entry's Status MUST be planned, active or deprecated; any other is a
  form problem.
- A planned entry is ruled but not built: its symbol MUST be generated, a
  test MAY cite it, and the gate MUST ask nothing else of it.
- An active entry MUST be covered.
- A deprecated entry is retired: its symbol MUST be generated, and the
  gate MUST ask nothing of it.
- How the generated files mark a planned or a deprecated entry is
  CH-020's.

Why: a test written before the code is the point of a scenario. A symbol
that appeared only once the entry was active would make such a test fail
to compile, so every status generates one and only active entries are
held to coverage.

Scenarios:
1. `planned-cited-freely`: given a test citing a planned entry, when the
   gate runs, then it reports no problem.
2. `status-form`: given a status outside the three, when it is read, then
   it is a form problem.

Changes:
- 2026-10-09: a deprecated entry's marking moved to CH-020, the generated
  files being its subject: deprecated-marked restated CH-020's
  planned-and-deprecated-marked and is dropped, its test citing CH-020's,
  and the text names CH-020 for how a planned or a deprecated entry is
  marked; the status form stays here.
- 2026-10-09: a status outside the three is ruled here alone, the status
  form being this entry's subject: CH-010's fields-checked named it too
  and now points here, and the text states it.
- 2026-11-26: serves CH-002, CH-261126-PFG5, where it served CH-002,
  CH-003: CH-003 is deprecated and the need above it is served directly
  (backlog 004).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-015 The gate MUST fail a Serves or See that names no existing entry, or that names the entry itself

Status: active
Serves: CH-002

- Every id under Serves or See MUST name an entry the charter holds.
- An entry naming itself MUST fail.
- An observation's Supports MUST name existing entries too.

Why: a relation to nothing is a broken link, and a reader following it
should be told at the gate, not when they reach the end of the link.

Scenarios:
1. `unknown-target`: given a Serves or See naming no entry, when the gate
   runs, then it fails naming the relation and the id.
2. `self-relation`: given an entry that serves itself, when the gate
   runs, then it fails.
3. `observation-supports-existing`: given an observation supporting an id
   that is not an entry, when the gate runs, then it fails.

Changes:
- 2026-11-26: the earlier form's relations are no longer read
  (Q-261015-N4HJ): `stakeholder-realizes-nothing` is dropped with them,
  and `unknown-target` names Serves and See alone.
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261126-HK6S The record reader MUST read an observation as a heading in observations.md, with its subject, status, date or source and the entries it supports

Status: active
Serves: CH-261126-P2SK
See: CH-035, CH-032

- observations.md in the claims directory holds one observation per
  heading, `## Name`, the name an upper-camel word unique among the
  observations, then its text: what was seen or read.
- Its field lines are Subject, Status (observed, documented, broken or
  unsure), Observed (a date), Source and Supports, the entries that rest
  on it.
- An unknown field, a field twice, a status outside the four or a name
  declared twice MUST be a form problem naming the line.
- Each language's generated file MUST declare a symbol per observation,
  by its name, which its probes' observes call names.

Why: a fact written as Go code could be checked only by a Go probe and
read only by someone reading Go. Written as markdown beside the entries
it supports, it can be checked by a probe in any language and seen by
every reader of the record.

Scenarios:
1. `observation-read` (owner, backlog 003): given a well-formed
   observation, when it is read, then it yields its name, subject,
   status, date, source, the entries it supports and its text.
2. `observation-form-checked` (owner, backlog 003): given an unknown
   field, a field twice, a status outside the four or a name declared
   twice, when it is read, then each is a form problem naming the line.
3. `observation-symbol-generated` (owner, backlog 003): given an
   observation, when the files are generated, then each declares a symbol
   for it, and a probe's observes naming it is read in Go, Python and
   TypeScript.

Changes:
- 2026-11-26: serves CH-261126-P2SK, where it served CH-005: CH-005 is
  deprecated and the need above it is served directly (backlog 004).
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.
