---
name: claims-draft
description: Turn a backlog story's rulings into charter entries, a decision or a Changes line in its area document or a need in claims/needs.md, at the right altitude, with a scenario per clause and a new id from claimsgen. Use when a story's Ruling table is settled and its ids are not yet written, e.g. "draft the entries for backlog 038", "write the decisions for this ruling". Not building (claims-implement) or reviewing a change (claims-review).
---

# Draft entries from a story

The form is `claims/README.md` ("Writing a decision", "Altitude"); the
best worked example is `claims/`. Nothing is minted, drafted
or promoted: the ruling table is the draft, the area document is the
record, and an id is one `claimsgen next` gives.

## Procedure

1. **Read the ruling.** Open the story, `backlog/NNN-slug.md`, and its
   `## Ruling` table: one row per question and its ruling. Note the ids
   the story cites and the surfaces it introduces.

2. **Find the area.** Decisions go in `claims/<area>.md`, needs in
   `claims/needs.md`. `ls claims/*.md` (not `DIGEST.md`, which `make claims` generates from
   the others, nor `README.md`) and read the document whose
   `Code:` line names the code the ruling touches. A new area gets a
   document with a title line, a `Code:` line naming its packages and a
   sentence on what it holds. Until an area has moved, its older specs
   are files in `claims/specs/`; grep both.

3. **New entry or a change?** For each ruling, grep `claims/*.md` and
   `claims/specs/` for the ids the story cites and the words of the
   ruling. A ruling that alters an existing decision is a `Changes:` line
   on it (`- YYYY-MM-DD: what changed and why`), with its text and
   scenarios edited to match; `go run ./cmd/claimsgen anchors ID` shows
   the code and tests the change will reach. A ruling nobody has ruled
   before is a new entry. Not everything is a decision: a helper's
   contract, a log line, a refactor are code and tests without an entry.

4. **Take the ids.** `go run ./cmd/claimsgen next SW 3` gives three
   new ids, one a line, each the day's date and a random code such as
   `SW-261001-K7Q2`; ask for as many as the story has new decisions
   (`next STK` for a need). Never write an id by hand or reuse one.

5. **Write the entry** in the area document, in this order: the heading
   `## SW-261001-K7Q2 Title`, the title a sentence a reviewer could disagree
   with; `Status: planned` (`active` when the build lands in the same
   change); `Kind:` only for an invariant, a bound or a structure rule;
   `Serves:` the need or decision it serves; `See:` coupled decisions;
   `Accounts:` each new environment variable, operation, event kind or
   build tag it owns (`env:NAME`, `op:GET /path`, `event:kind`,
   `tag:name`); a blank line; the text, one to six sentences; `Why:`;
   `Scenarios:`, numbered, each "N. `kebab-handle` (who ruled): an input
   and an observable outcome", the provenance `(owner, #444)` when the
   story's ruling is the owner's, `(delegate name, #12)` for a delegate
   `claims/policy.md` names for the area, else `(agent, provisional)` with
   a question in `claims/questions.md` (id from `claimsgen next Q`) saying
   what was chosen and why; `Guidance:` only when the ruling said how. A need
   has no Kind, Serves or Guidance. Abridged from `claims/form.md`:

   ```markdown
   ## CH-012 An id is a prefix with a number or with a date and a code, declared once, and a new id is the day it was taken and a code drawn at random

   Status: active
   Serves: CH-003

   A new id is its prefix, the day's date and a code drawn at random,
   drawn again when the tree already holds that id. Ids are unique across
   the area documents; a second declaration fails naming both places.

   Why: no state file and no counter to share: the tree is the state, and
   a random code is not the same on two branches written at the same time.

   Scenarios:
   1. `declared-once`: an id declared in two documents fails naming both.
   2. `next-free`: a new id is its prefix, the UTC date and a code of the
      alphabet, and a draw the tree already holds is drawn again.
   ```

6. **Hold the altitude.** Name the boundary (a request, an event on the
   log, a file, what the model is sent, what the person sees), never an
   identifier; a surface goes in `Accounts`, not the text. One ruling is
   one decision: a behaviour that shows at several surfaces is one entry
   with a scenario per surface. Read the text clause by clause: every
   clause gets a scenario, and a clause not worth a scenario is not worth
   writing. A ruled number is part of the text; tuning is Guidance or
   nothing.

7. **Ask the sufficiency questions** of a decision that validates input,
   gates a call, bounds a resource, kills a process or persists state:
   where is the trust boundary; what happens on the failure path; what is
   the limit; what happens in the concurrent case; what is never done.
   Each answer is a clause with a scenario, or an exclusion stated in the
   text.

8. **Check.** `make claims`, then `go run ./cmd/claimsgen lint`; both
   clean but for lint's line counting the open questions. Lint names a
   decision that serves nothing.

9. **Record and stop.** List the ids in the story's `## Charter` section,
   one line each with its title. Present the area document diff and stop
   for the owner's review; the build is `claims-implement`. Do not commit.

## Done when

- Every ruling is a new entry or a dated Changes line, in its area
  document, at the altitude above, with a scenario per clause.
- The ids came from `claimsgen next`; `make claims` and `claimsgen lint`
  are clean but for the open-questions line; the story's `## Charter`
  lists the ids.
- The owner has the diff; nothing is built or committed.
