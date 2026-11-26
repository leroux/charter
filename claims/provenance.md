# Knowing who decided

Code: internal/record, internal/diff

The problem this area answers: an agent's default and an owner's ruling
look alike once written. It holds who ruled each scenario, the questions
nobody has ruled yet, and the check that every change to the record is
written down.

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

## CH-261009-QQC8 charter diff names what a change did to the record that its Changes lines do not say

Status: active
Serves: CH-261009-3FHW

Given a git revision, charter diff reads the area documents as they
were there and as they are, and names, by id: an entry removed; an active
entry made planned; an entry deprecated, or its text, Kind or a scenario
changed or removed, with no Changes line the revision lacked; and a
scenario whose ruling was the owner's or a delegate's and is now an
agent's. An entry only set active, or changed with a new dated line, is
not named. It exits 1 when it names anything.

Why: the gate reads one tree; the drift it cannot see is between two, a
test made to assert the opposite of the scenario it cites and nothing
written down. Against the base of a pull request, this is the check that
a change to a decision is a change recorded.

Scenarios:
1. `removed-named` (owner): an id present at the revision and
   absent now is named: an entry is deprecated, never deleted.
2. `unbuilt-named` (owner): an entry active at the revision
   and planned now is named.
3. `unrecorded-change-named` (owner): an entry whose text,
   Kind or scenarios differ from the revision's with no Changes line the
   revision lacked is named, and so is one deprecated without a line.
4. `ruling-lost-named` (owner): a scenario ruled by the owner
   or a delegate at the revision and provisional now is named.
5. `recorded-change-passes` (owner): a change with a new dated
   Changes line, or an entry only set active, names nothing.

Changes:
- 2026-11-26: the verb is charter diff, of the one charter command
  (CH-040); the ruling is unchanged (backlog 003).
