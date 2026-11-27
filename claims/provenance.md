# Knowing who decided

Code: internal/record, internal/diff

This area answers the problem that an agent's default and an owner's
ruling look alike once written. It covers who ruled each scenario, the
questions nobody has ruled yet, and the check that every change to the
record is written down.

## CH-261009-VY5C The record reader MUST read who ruled a scenario: the owner with what it cites, a named delegate, or, when nothing is written, an agent, provisionally

Status: active
Serves: CH-261009-3FHW

- A scenario's handle MAY be followed, in parentheses, by who ruled it:
  "owner", with what the ruling cites after a comma; "delegate" with the
  delegate's name and, after a comma, what it cites; or "agent,
  provisional".
- A scenario with nothing written MUST be read as an agent's,
  provisional, and the registry MUST say that it was not written.
- Anything else MUST be a form problem naming the line and the three
  forms.
- The provenance is part of the scenario as the registry gives it.

Why: who decided is a fact a person writes, never one the tooling infers.
Written beside the scenario it rules, it lets a reader, a reviewer and a
policy tell an owner's ruling from an agent's default.

Scenarios:
1. `owner-read` (owner): given "(owner, #444)" and "(owner)", when they
   are read, then the first gives the owner with the reference and the
   second the owner alone.
2. `delegate-read` (owner): given "(delegate name, #12)", when it is
   read, then it gives the delegate's name and the reference.
3. `default-is-provisional` (owner): given a scenario with nothing in
   parentheses, when it is read, then it is an agent's, provisional, and
   says it was not written.
4. `malformed-refused` (owner): given a provenance naming anyone but
   owner, delegate or agent, a named owner, an unnamed delegate, or an
   agent with a reference, when it is read, then it is a form problem
   naming the forms to write.

Changes:
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261009-4B3B The question file MUST hold what nobody has ruled, each question naming its entry, who asked and the provisional answer until it is answered

Status: active
Serves: CH-261009-3FHW

- questions.md in the claims directory holds one question per heading,
  `## Q-ID question`, with field lines Asked-by, Entry (an id, or an id
  and a scenario as ID/handle), Provisional (the answer taken meanwhile)
  and, once a person has answered, Answer, opening with the date and who
  answered.
- A question without Asked-by, Entry or Provisional, naming an entry or
  scenario the record does not hold, with an Answer that does not open
  with a date, or with an id declared twice, MUST be a form problem naming
  the line.
- Once its answer is written into the record, as a scenario's provenance
  or a Changes line naming its id, an answered question MUST be deleted;
  git keeps it.
- `charter questions` MUST list each open question with its id, the entry
  it is about, its text, the provisional answer and who asked, then how
  many are open; with `-all` it MUST list the answered ones too, each with
  its answer.

Why: an agent that must choose something nobody has ruled needs a place
to say so and go on, and the owner needs the list to clear; without it
the agent's choice lands looking like a decision. An answer already
written where it rules would be a second copy that goes stale, so an
answered question is deleted once its answer is in the record.

Scenarios:
1. `question-read` (owner): given a well-formed question, when it is
   read, then it yields its id, text, asker, entry, provisional answer
   and, when present, its answer.
2. `unknown-entry-refused` (owner): given a question about an entry or a
   scenario the record does not hold, when it is read, then it is a form
   problem.
3. `missing-parts-refused` (owner): given a question without Asked-by,
   Entry or Provisional, or whose Answer does not open with a date, when
   it is read, then it is a form problem naming what to write.
4. `questions-verb-lists` (owner): given open and answered questions,
   when `charter questions` runs, then it prints each open question with
   its id, entry, text, provisional answer and who asked, then how many
   are open; with `-all`, then it adds the answered ones with their
   answers.

Changes:
- 2026-10-09: lint summarises the open questions in one line and the
  questions verb lists them.
- 2026-11-27: an answered question is deleted once its answer is written
  into the record (Q-261015-YSKC, answered yes by the owner), and lint,
  which counted the open ones, is gone; `open-listed` goes with it.
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.

## CH-261009-QQC8 charter diff MUST name every change to the record that its Changes lines do not state

Status: active
Serves: CH-261009-3FHW

Given a git revision, `charter diff` reads the area documents as they
were there and as they are now, and MUST name, by id:

- an entry removed that was not deprecated there, since an entry is
  deprecated first and deleted in a later change;
- an active entry made planned;
- an entry deprecated, or its text, Kind or a scenario changed or
  removed, with no Changes line the revision lacked;
- a scenario whose ruling was the owner's or a delegate's and is now an
  agent's.

An entry only set active, or changed with a new dated line, MUST NOT be
named. diff MUST exit 1 when it names anything.

Why: the gate reads one tree, so the drift it cannot see is between two
trees: a test changed to assert the opposite of the scenario it cites,
with nothing written down. Run against the base of a pull request, diff
is the check that every change to a decision is also a change recorded
in the record.

Scenarios:
1. `removed-named` (owner): given an id present at the revision and
   absent now, when diff runs, then it is named, unless it was deprecated
   there.
2. `unbuilt-named` (owner): given an entry active at the revision and
   planned now, when diff runs, then it is named.
3. `unrecorded-change-named` (owner): given an entry whose text, Kind or
   scenarios differ from the revision's with no Changes line the revision
   lacked, or one deprecated without a line, when diff runs, then it is
   named.
4. `ruling-lost-named` (owner): given a scenario ruled by the owner or a
   delegate at the revision and provisional now, when diff runs, then it
   is named.
5. `recorded-change-passes` (owner): given a change with a new dated
   Changes line, or an entry only set active, when diff runs, then it
   names nothing.

Changes:
- 2026-11-26: the verb is charter diff, of the one charter command
  (CH-040); the ruling is unchanged (backlog 003).
- 2026-11-27: an entry deprecated at the revision may be gone now: the
  owner ruled that retiring is deprecating an entry and stripping its
  code and tests, then deleting it in a later change, git keeping the
  text; until now an entry was never deleted.
- 2026-11-27: reworded in the owner's claim form; the ruling is unchanged.
