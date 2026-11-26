# The form of an entry

Code: chartertest

An entry is a heading in an area document; this document holds what the
parser holds an entry and a document to.

## CH-010 An entry is a heading with id and title, then field lines, text, Why, Scenarios, Guidance and Changes in order

Status: active
Serves: CH-001, CH-002
See: CH-013

An area document is a markdown file whose header is a title and an
optional Code line, and whose entries are headings of the form
"## ID Title". The lines after the heading, up to a blank line, are
fields: Status, Kind, Serves, See and Accounts, each at most once; the
values a Status takes are CH-013's. The body is paragraphs: the text
first, then paragraphs opened by Why, Scenarios, Guidance and Changes, in
that order and each at most once; the text and Why are required, and an
active entry lists at least one scenario. A scenario is a numbered item
with a handle in backticks, lower-case kebab, unique in the entry, and a
sentence; a change is a dated item; a continuation line is indented. Every
form problem names the document, the line and the form to write. The
header, before the first entry, may hold prose, a vocabulary list, tables
and diagrams that orient the reader; it is not parsed beyond its title and
Code line and has no section headings, and a line inside a fenced block,
in the header or in an entry, is text and never a heading.

Why: the parts are what the reviewer reads, in the order they read them;
a fixed order and a fixed form make a diff say what changed.

Scenarios:
1. `every-part-read`: a well-formed entry yields its id, symbol, title,
   status, kind, relations, accounts, text, why, scenarios with their
   handles and lines, guidance and dated changes.
2. `heading-form`: a heading without an id and a title is a form problem.
3. `fields-checked`: an unknown field, a field twice, a kind outside the
   three and a malformed surface key are form problems.
4. `required-parts`: an entry without a Status, without text or without a
   Why is a form problem, and so is an active entry without Scenarios.
5. `scenario-form`: a scenario item that is not "N. `handle`: sentence", a
   handle that is not lower-case kebab and a handle used twice are form
   problems.
6. `parts-in-order`: a part out of order, a part twice, prose after a
   labelled part and a heading inside an entry are form problems.
7. `change-form`: a change that is not "- YYYY-MM-DD: text", or whose date
   is not on the calendar, is a form problem.
8. `problems-locate`: every form problem names the document, the line and
   the rule.
9. `blank-lines-before-fields`: blank lines between the heading and the
   first field are allowed; the first blank line after a field ends them.
10. `readme-is-not-a-document`: README.md in the directory is documentation
    and is not read as an area document.
11. `header-is-prose`: a header of prose, a list, a table and a fenced
    block with a "#" line parses clean and yields no entry; a section
    heading there is a form problem.
12. `heading-in-fence-is-text`: a fenced "#" or "## " line in an entry's
    body is text, no form problem and no entry; a fence among the field
    lines is no field line and starts no entry.

Guidance: the header carries no rulings, and a hand-written index of the
entries is not part of it, since the headings are the index; neither is
something the parser can tell, so review holds the header to them.

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
- 2026-10-09: the title, the text and two scenarios reworded to the
  style's bars; the ruling is unchanged.
- 2026-10-09: a status outside the three, which fields-checked named
  beside CH-013's status-form, is CH-013's alone, the status form being
  its subject; fields-checked no longer names it, and the text names
  CH-013 for a Status's values, keeping the fields' presence and order
  here.

## CH-011 A need lives in needs.md and has no Kind, Serves or Guidance

Status: active
Serves: CH-001

An entry of the document named needs.md is a need. A need serves nothing,
carries no kind and no guidance; each of those is a form problem.

Why: a need is what a person can rely on; it is realized by the decisions
that serve it, and it is not the place for how.

Scenarios:
1. `need-by-document`: an entry in needs.md is a need, and one in any other
   document is a decision.
2. `need-fields-refused`: a Kind, a Serves or a Guidance on a need is a
   form problem.

## CH-012 An id is declared once; a new one is its prefix, the day it was taken and a random code

Status: active
Serves: CH-003
See: CH-040

An id is a prefix and then either a number, as SW-118, or the day it was
taken and a code, as SW-261001-K7Q2; a heading whose id is neither is a
form problem. The day is the UTC date as six digits, year, month and
day; the code is four characters of BCDFGHJKMNPQRSTVWXYZ23456789, the
consonants and digits that spell no word and are not read as one
another. The numbered ids stand as they are, and next gives every new
id, of the second kind and one the tree does not hold, as CH-040 rules.
Ids are unique across the area documents; a second declaration fails naming both places and saying that the entry
that has not landed takes a new id. The entries of the area documents are
ordered by prefix, and within one the numbered ids first, by number, then
the dated ones by date and code, and every generated file declares
them in that order.

Why: no state file and no counter to share: the tree is the state. One
above the highest is the same number on every branch cut from one tree,
so branches written at the same time met at merge, and the later one
rewrote its id in every document, anchor and citation that named it. A
code drawn from six hundred thousand is not the same on two branches, so
they do not meet; the gate still names both places on the day they do.
The date adds nothing to that. It is there so that ids read in the order
they were taken.

Scenarios:
1. `declared-once`: an id declared in two documents fails naming both and
   saying the entry not yet landed takes a new id.
2. `id-form`: numbered and dated ids parse; an off-calendar date or a code
   not four characters of the alphabet is a form problem.
3. `ordered-by-date`: area entries come by prefix, and within one numbered
   ids by number before dated ones by date and code, in every generated
   file alike.

Changes:
- 2026-10-01: a new id is the day it was taken and a random code, where
  it was one above the highest: branches written at the same time took
  the same number, and the later one rewrote it everywhere it was
  written. The numbered ids stand, and the declared-twice problem says
  which entry takes a new id.
- 2026-10-09: the title and four scenarios reworded to the style's bars;
  the ruling is unchanged.
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

## CH-013 Status is planned, active or deprecated, every status generates a symbol, and coverage is asked of active entries only

Status: active
Serves: CH-002, CH-003
See: CH-020

An entry's Status is planned, active or deprecated, and any other is a
form problem. A planned entry is ruled and not built; its symbol is
generated, a test may cite it, and the gate asks nothing else of it. An
active entry must be covered. A deprecated entry is retired; its symbol
is generated, and the gate asks nothing of it. How the generated files
mark a planned or a deprecated entry is CH-020's.

Why: a test written before the code is the point of a scenario; a symbol
that appears only once the entry is active forbids it.

Scenarios:
1. `planned-cited-freely`: a test citing a planned entry is no problem.
2. `status-form`: a status outside the three is a form problem.

Changes:
- 2026-10-09: a deprecated entry's marking moved to CH-020, the generated
  files being its subject: deprecated-marked restated CH-020's
  planned-and-deprecated-marked and is dropped, its test citing CH-020's,
  and the text names CH-020 for how a planned or a deprecated entry is
  marked; the status form stays here.
- 2026-10-09: a status outside the three is ruled here alone, the status
  form being this entry's subject: CH-010's fields-checked named it too
  and now points here, and the text states it.

## CH-014 The earlier form, one file per spec, is read beside the area documents until it is migrated

Status: deprecated
Serves: CH-003

A spec file in the spec directory, named for its numbered id, never a
dated one, with a lens and a slug, is parsed into the same registry as an
area entry: its fields in their fixed order, its Realizes counting as
Serves, and its links pointing at the target's current file. The directory
may be absent.

Why: two hundred and sixty specs move one area at a time; the gate holds
both forms while they do.

Scenarios:
1. `file-form-read`: a spec file yields its id, symbol, lens, status,
   accounts, text and relations, with each link's target.
2. `file-form-checked`: a field out of order, an unknown field, a link
   whose text is not the id, a malformed change row and text outside a
   field are problems.
3. `file-name-checked`: a file not named LENS-NNN-slug.md with a lens of
   the seven, one named by a dated id among them, is a problem.
4. `link-into-a-document`: a spec file's link to an entry that has moved
   into an area document is checked by id alone, whatever file it names.

Changes:
- 2026-10-09: that a dated id names no spec file is ruled here alone, a
  spec file's name being this entry's subject: CH-012's id-form stated it
  too and now points here, and the text and file-name-checked say it.
- 2026-11-26: deprecated (Q-261015-N4HJ): neither charter holds a spec
  file, so the earlier form is no longer read, and its code goes with the
  rewrite (backlog 003).

## CH-015 Serves and See name existing entries, a need serves nothing, and nothing names itself

Status: active
Serves: CH-002

Every id under Serves or See names an entry the charter holds, and an
entry naming itself fails. An observation's Supports names existing
entries too.

Why: a relation to nothing is a broken link, and the reader following
it deserves to be told at the gate, not at the end.

Scenarios:
1. `unknown-target`: a Serves or See naming no entry fails naming the
   relation and the id.
2. `self-relation`: an entry that serves itself fails.
3. `observation-supports-existing`: an observation supporting an id that
   is not a spec fails.

Changes:
- 2026-11-26: the earlier form's relations are no longer read
  (Q-261015-N4HJ): `stakeholder-realizes-nothing` is dropped with them,
  and `unknown-target` names Serves and See alone.

## CH-261126-HK6S An observation is a heading in observations.md, with its subject, status, date or source and the entries it supports

Status: active
Serves: CH-005
See: CH-035, CH-032

observations.md in the claims directory holds one observation per heading,
"## Name", the name an upper-camel word unique among the observations,
then its text, what was seen or read. Its field lines are Subject, Status
(observed, documented, broken or unsure), Observed (a date), Source and
Supports, the entries that rest on it. An unknown field, a field twice, a
status outside the four or a name declared twice is a form problem naming
the line. Each language's generated file declares a symbol per
observation, by its name, and an observes its probes call.

Why: a fact written as Go could be checked only by a Go probe and read
only by someone reading Go; written as markdown beside the entries it
supports, a probe in any language checks it and every reader of the record
sees it.

Scenarios:
1. `observation-read` (owner, backlog 003): a well-formed observation
   yields its name, subject, status, date, source, the entries it supports
   and its text.
2. `observation-form-checked` (owner, backlog 003): an unknown field, a
   field twice, a status outside the four and a name declared twice are
   form problems naming the line.
3. `observation-symbol-generated` (owner, backlog 003): each generated
   file declares a symbol per observation, and a probe's observes naming
   one is read in Go, Python and TypeScript.

## CH-261009-VY5C A scenario says who ruled it: the owner with what it cites, a named delegate, or, unwritten, an agent, provisional

Status: active
Serves: CH-261009-3FHW

A scenario's handle may be followed, in parentheses, by who ruled it:
"owner", with what the ruling cites after a comma; "delegate" with the
delegate's name and, after a comma, what it cites; or "agent,
provisional". Nothing written is an agent's, provisional, and the
registry says it was not written. Anything else is a form problem naming
the line and the three forms. The provenance is part of the scenario as
the registry gives it.

Why: who decided is a fact a person writes, never one the tooling infers;
written beside the promise it lets a reader, a reviewer and a policy tell
a ruling from a default.

Scenarios:
1. `owner-read` (owner): "(owner, #444)" gives the owner with
   the reference, and "(owner)" the owner alone.
2. `delegate-read` (owner): "(delegate name, #12)"
   gives the delegate's name and the reference.
3. `default-is-provisional` (owner): a scenario with nothing
   in parentheses is an agent's, provisional, and says it was not written.
4. `malformed-refused` (owner): a provenance naming anyone but
   owner, delegate or agent, a named owner, an unnamed delegate, or an
   agent with a reference is a form problem naming the forms to write.

Changes:
- 2026-10-09: the title and the malformed-refused scenario reworded to the
  style's bars; the ruling is unchanged.

## CH-261009-4B3B The question file holds what nobody ruled, each question naming its entry, who asked and the provisional answer until answered

Status: active
Serves: CH-261009-3FHW

questions.md in the entries directory holds one question per heading,
"## Q-ID question", with field lines Asked-by, Entry (an id, or an id and
a scenario as ID/handle), Provisional (the answer taken meanwhile) and,
once a person has answered, Answer, opening with the date and who
answered. A question without Asked-by, Entry or Provisional, naming an
entry or scenario the record does not hold, with an Answer that does not
open with a date, or with an id declared twice, is a form problem naming
the line. A question is never deleted. Lint prints one line, at the
question file's first line, saying how many questions are open and naming
the command that lists them; with none open it prints none. claimsgen
questions lists each open question with its id, the entry it is about,
its text, the provisional answer and who asked, then how many are open;
with -all it lists the answered ones too, each with its answer.

Why: an agent that must choose something unruled needs a place to say so
and go on, and the owner needs the list to clear; without it the choice
lands as if it were a decision. Lint is read for what to look at, and a
line per open question buried its advice; a count keeps it readable, and
the list is one command away.

Scenarios:
1. `question-read` (owner): a well-formed question yields its
   id, text, asker, entry, provisional answer and, when present, its
   answer.
2. `unknown-entry-refused` (owner): a question about an entry
   or a scenario the record does not hold is a form problem.
3. `missing-parts-refused` (owner): a question without
   Asked-by, Entry or Provisional, or whose Answer does not open with a
   date, is a form problem naming what to write.
4. `open-listed` (owner): lint prints one line, at the
   question file's first line, counting the open questions and naming the
   command that lists them, and none when no question is open.
5. `questions-verb-lists` (owner): claimsgen questions prints
   each open question with its id, entry, text, provisional answer and who
   asked, then how many are open; -all adds the answered ones with their
   answers.

Changes:
- 2026-10-09: lint summarises the open questions in one line and the
  questions verb lists them.
- 2026-10-09: the title reworded to twenty words; the ruling is unchanged.
