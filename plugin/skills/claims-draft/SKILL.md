---
name: claims-draft
description: Turn a backlog story's rulings into charter entries, a decision or a Changes line in its area document or a need in claims/needs.md, at the right altitude, with a scenario per clause and a new id from charter. Use when a story's Ruling table is settled and its ids are not yet written, e.g. "draft the entries for backlog 038", "write the decisions for this ruling". Not building (claims-implement) or reviewing a change (claims-review).
# charter:realizes CH-261126-5JMD
---

# Draft entries from a story

The form is the one every charter shares (the charter repository's
README, "Adopting it"); `claims/charter.toml` says where the record and
the code are, and `charter anchors ID` reads any entry. Nothing is
minted, drafted or promoted: the ruling table is the draft, the area document is the
record, and an id is one `charter next` gives.

## Voice

Write in the voice the project's agent instructions give: its CLAUDE.md
or AGENTS.md, and any file they import. Where they say nothing about
voice, write clear technical prose for the reader at hand. Lead with the
point, and say how each statement follows from the one before. Use the
field's established terms, and define a term the reader may not know
where it first appears. Give one thing one name. State the mechanism, the
values and the limits rather than a summary of them. Use a list for
separate items and prose where the relationships between them matter.

## Entry form

Write each entry in the form the project's agent instructions (CLAUDE.md
or AGENTS.md) give for claims. Where they give none, use this default:

- **Title:** the component that carries the obligation, with a BCP 14
  keyword, as one claim a reviewer could disagree with: "The proxy MUST
  check an allowlisted connection's TLS server name before relaying it".
- **Text:** RFC 2119 bullets, one requirement per bullet. Use connected
  prose instead when the entry is a single requirement, or when its
  clauses only make sense together (a sequence, a cause).
- **Key words:** MUST, MUST NOT, SHOULD, SHOULD NOT and MAY, in capitals
  only, with their BCP 14 meanings (RFC 2119, RFC 8174). Put the
  negation in the keyword ("MUST NOT send"), never in the subject
  ("nothing MUST be sent"). If `needs.md` does not yet say so, add to its
  header that the key words are interpreted as BCP 14 describes.
- **Why:** explain the problem or threat, how the rule meets it, and what
  the rule costs or avoids. Add no reason the ruling does not support.
- **Scenarios:** one numbered item each, in given / when / then form, with
  a concrete input and an observable outcome.

## Procedure

1. **Read the ruling.** Open the story, `backlog/NNN-slug.md`, and its
   `## Ruling` table: one row per question and its ruling. Note the ids
   the story cites and the surfaces it introduces.

2. **Find the area.** Decisions go in `claims/<area>.md`, needs in
   `claims/needs.md`. `ls claims/*.md` (not the digest, which `charter gen` generates from
   the others, nor `README.md`) and read the document whose
   `Code:` line names the code the ruling touches. A new area gets a
   document with a title line, a `Code:` line naming its packages and a
   sentence on what it holds, and its name goes in
   `claims/reading-order.txt` where a reader should meet it.

3. **New entry or a change?** For each ruling, grep `claims/*.md` for the ids the story cites and the words of the
   ruling. A ruling that alters an existing decision is a `Changes:` line
   on it (`- YYYY-MM-DD: what changed and why`), with its text and
   scenarios edited to match; `charter anchors ID` shows
   the code and tests the change will reach. A ruling nobody has ruled
   before is a new entry. Not everything is a decision: a helper's
   contract, a log line, a refactor are code and tests without an entry.

4. **Take the ids.** `charter next SW 3` gives three
   new ids, one a line, each the day's date and a random code such as
   `SW-261001-K7Q2`; ask for as many as the story has new decisions
   (`next STK` for a need). Never write an id by hand or reuse one.

5. **Write the entry** in the area document, in this order: the heading
   `## SW-261001-K7Q2 Title`, the title as the entry form says; `Status: planned` (`active` when the build lands in the same
   change); `Kind:` only for an invariant, a bound or a structure rule;
   `Serves:` the need or decision it serves; `See:` coupled decisions;
   `Accounts:` each new environment variable, operation, event kind or
   build tag it owns (`env:NAME`, `op:GET /path`, `event:kind`,
   `tag:name`); a blank line; the text, as long as the ruling needs; `Why:`;
   `Scenarios:`, numbered, each "N. `kebab-handle` (who ruled): given an
   input, when it happens, then an observable outcome", the provenance `(owner, #444)` when the
   story's ruling is the owner's, `(delegate name, #12)` for a delegate
   `claims/policy.md` names for the area, else `(agent, provisional)` with
   a question in `claims/questions.md` (id from `charter next Q`) saying
   what was chosen and why; `Guidance:` only when the ruling said how. A need
   has no Kind, Serves or Guidance. An entry in the default form, abridged:

   ```markdown
   ## CH-012 An id MUST be declared once, and a new id MUST be its prefix, the day it was taken and a random code

   Status: active
   Serves: CH-261126-PFG5

   - `charter next` MUST give every new id: its prefix, the day's UTC date
     and a code drawn at random, drawn again when the tree already holds
     that id.
   - Ids MUST be unique across the area documents. A second declaration
     MUST fail, naming both places.

   Why: ids need no state file and no counter to share, because the tree
   is the state. A random code is unlikely to repeat on two branches
   written at the same time, so branches no longer collide at merge.

   Scenarios:
   1. `declared-once`: given an id declared in two documents, when the
      gate runs, then it fails naming both.
   2. `next-free`: given a prefix, when `charter next` runs, then it gives
      the prefix, the UTC date and a code of the alphabet, drawing again a
      code the tree already holds.
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

8. **Check.** `charter gen`, then `charter check`. A planned entry's
   untested scenarios are not problems; a broken relation, form or id is.
   Every decision names what it serves.

9. **Record and stop.** List the ids in the story's `## Charter` section,
   one line each with its title. Present the area document diff and stop
   for the owner's review; the build is `claims-implement`. Do not commit.

## Done when

- Every ruling is a new entry or a dated Changes line, in its area
  document, at the altitude above, with a scenario per clause.
- The ids came from `charter next`; `charter gen` and `charter check`
  are clean; the story's `## Charter` lists the ids.
- The owner has the diff; nothing is built or committed.
