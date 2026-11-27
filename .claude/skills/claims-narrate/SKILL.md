---
name: claims-narrate
description: Draft a doc that tells a story over the record and the code, a walkthrough of a flow, an area or a need, with checked references to entries and declarations and embeds of the entries it rests on; or, on a doc charter check warns is stale, propose an edit of what changed. Never writes over a curated doc's prose and never curates: a person reads the draft and runs charter curate. Use when asked to "write a walkthrough of X", "explain how Y works for a newcomer", or "check says docs/Z.md is stale".
---

# Narrate

A doc is the part of a charter a newcomer reads first: how the system
works, where in the code, and which decisions each step carries. The
record holds what binds; a doc tells the story that joins it, and its
prose is free. What the gate checks is its references: a link whose text
opens with an entry's id, a link whose text is a backticked name and
whose target is the file declaring it, and an embed. What check warns of is
whether anything a doc references has changed since a person last read
it.

The agent drafts; a person curates. This skill never runs `charter
curate` and never writes over a curated doc's prose.

## Draft a doc

1. **Find the subject.** The person names a flow ("a request from the
   socket to the answer"), an area, or a need. Read the entries it
   touches (the area documents, `charter anchors ID` for each) and the
   code their anchors point at. Read the docs `charter.toml`'s `docs`
   names already: a new doc does not retell one.
2. **Follow the data.** Order the doc by what happens to the data, from
   where it enters to where it leaves, not by package or by entry. Each
   step names the declaration that does it, as a code reference
   (``[`Limiter.Allow`](../limiter/limiter.go)``), and the entry it
   carries, as a reference (`[SW-…](x)`; `charter gen` writes the
   target and keeps the title current).
3. **Embed what binds.** Where the story rests on a rule, embed the entry
   rather than paraphrase it: `<!-- charter:embed SW-… -->` and
   `<!-- /charter:embed -->`, with nothing between; gen writes the text.
   Paraphrase only what no entry states, and say plainly that it is how
   the code does it today.
4. **Name the traps.** What a newcomer would get wrong, and why the code
   is shaped against it: a lock held for the whole request, a clock read
   in one place, an observation the code rests on.
5. **Check.** Write the doc under a path `docs` matches, run `charter
   gen` (the links' targets, line anchors and embeds), then `charter
   check`: every reference resolves. Read each sentence against the
   declaration it names; one you cannot point at a line for, cut.
6. **Hand it over.** Present the doc and say it is drafted: `charter
   check` warns of it until a person reads it and runs `charter curate
   FILE -by NAME`.

## Refresh a stale doc

1. `charter check` warns of the doc and each reference that changed since it
   was curated; `charter diff REV` names the docs a change touched.
2. For each changed reference, read what changed (the entry's text and
   its Changes lines; the declaration's diff) and the doc's sentences
   around the reference.
3. Propose an edit of those sentences alone, as a diff for the person;
   leave the rest of the doc as it is. A change the doc does not need is
   said so, with the reason: the person curates again either way.
4. Never run curate; the hash is the person's word that they read it.

## Done when

- The doc's references resolve and `charter check` passes.
- Every sentence about the code names a declaration it can be checked
  against, and every rule it rests on is a reference or an embed.
- The person has the draft, or the proposed edit, and curates it.
