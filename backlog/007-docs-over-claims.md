# 007. Docs over claims: free prose that references and embeds claims and code, drafted by an agent, curated by a person, and warned on when what it rests on changes

Status: drafted · Size: M · See: 005, 006

The record holds what binds: claims, each declared once, anchored to code
and verified by tests. It has no place for the rest of what a reader needs
to understand a system: how it works (the flow of data through the code),
why it is shaped the way it is, and a story that connects many claims into
one walkthrough. Today that lives in prose the gate cannot see, like
torii's `docs/design.md`, and it goes stale without anyone being told.

This story adds docs: free markdown that references claims and code by
links the gate checks, can embed a claim so it reads in place, is drafted
by an agent and curated by a person, and is named by lint when anything it
references has changed since it was curated.

## The kinds of knowledge, and where each lives

| Kind | Example | Truth lives in | Checked by |
|---|---|---|---|
| Promise | Each key's bucket is its own | A claim | Its tests |
| Internal rule | The limiter reads time through its Clock alone | A claim | Its tests |
| Fact about the world | Qubes' default policy denies `qubes.ConnectTCP` | A claim (an observation) | A probe |
| Mechanism | How the proxy decides a connection | The code; a doc points at it | References resolve; staleness |
| Rationale | Our own tsnet program rather than tailscaled | A doc, and each claim's Why and Changes | Staleness |

## Ruling

Ruled by the owner, 2026-11-27, in design talk.

| # | Question | Ruling |
|---|---|---|
| 1 | Where are claims declared? | In the claim files alone, in the form they have now, so they stay easy to read by program. A doc never declares one. |
| 2 | What is a doc? | Any markdown that tells a story over claims and code: a walkthrough, a design doc, a guide. Its prose is free; only its references are checked. |
| 3 | Who writes docs? | An agent drafts them; a person curates them. |
| 4 | Are backlinks written by hand? | No. Which docs mention a claim, or a code symbol, is computed from the docs. |
| 5 | Does a stale doc fail the gate? | No, it warns. |

## Proposed, for the owner

Drafted by the agent, 2026-11-27. Each is the default until ruled.

| # | Question | Proposal |
|---|---|---|
| 6 | Which files are docs? | Those charter.toml's `docs` patterns name, e.g. `docs = ["docs/**/*.md", "README.md"]`. Nothing is read as a doc unless named. |
| 7 | How does a doc reference a claim? | A plain markdown link whose text opens with the claim's id: `[SW-261127-M36P](../claims/egress.md#sw-261127-m36p)`. `charter gen` keeps the target and the title after the id current. It renders on GitHub as written. |
| 8 | How does a doc reference code? | A markdown link whose text is a symbol in backticks and whose target is the file declaring it: ``[`engine.Accept`](../proxy/internal/engine/engine.go)``. The gate fails if the file does not declare the symbol; gen keeps a `#L` line anchor current. |
| 9 | How does a doc embed a claim? | A marker pair, `<!-- charter:embed SW-… -->` and `<!-- /charter:embed -->`, between which gen writes the claim's title and text from the claim file. The doc never holds a hand-written copy, and GitHub renders what gen wrote. |
| 10 | What does the gate check? | That every reference resolves: the claim exists and is not deprecated, the file declares the symbol, an embed names a claim. A broken reference fails, as a broken anchor does. |
| 11 | How is a doc's curation recorded? | A comment block at the doc's end, written by `charter curate FILE`: who curated it, the date, and a hash of each claim and each code symbol it references, as they were when curated. |
| 12 | What is a claim's hash? | Of its title, status and text, normalized for whitespace. A new Changes line alone does not change it. |
| 13 | What is a code symbol's hash? | Of its declaration's source, normalized for whitespace, from the language reader that already locates declarations for anchors. A symbol of a plain file (the `[[files]]` tables) is hashed whole-file. |
| 14 | When is a doc stale? | When a referenced claim's or symbol's hash differs from the one recorded at curation. Lint names the doc and each changed reference; a doc with no curation block is named as drafted. |
| 15 | Where do backlinks show? | `charter anchors ID` lists the docs that mention the claim; `export` carries them and the viewer shows them; `charter diff REV` lists the docs whose references changed, as docs to look at. |
| 16 | How is a doc drafted and refreshed? | A skill, `claims-narrate`: given a need, an area or a flow, it drafts a doc from the claims, the code and the docs already there. On a stale doc it proposes an edit to the changed parts only; it never regenerates over a curated doc. |
| 17 | Who checks a drafted doc against the code? | `claims-audit`, in a context that did not write it, adds: does the doc describe what the code it references does? |
| 18 | Can a stale doc ever block? | Only where `claims/policy.md` names it, as for a walkthrough of a security-critical path. The default is a warning. |

## Why hashes, not dates

A curation date or commit compared with git history needs that history
in full: a shallow clone in CI, or a squashed merge, loses it. A hash
recorded in the doc travels with the doc and needs nothing but the tree
as it is. It also changes only when the referenced text does, so a
reformatted file or a moved declaration does not make a doc stale.

## For the review

- Rulings 6 to 18 are the agent's; the reference syntax (7 to 9) most of
  all shapes how docs read and is the first to rule.
- Ruling 13 needs each reader to return a declaration's full span. Go's
  does through go/ast; Python's and TypeScript's readers find where a
  declaration starts and may need to find where it ends.
- Area documents keep their summary prose. Whether it moves to docs, leaving
  area files as claims alone, is left until docs exist.
- The charter specifies itself, so each ruling becomes a CH entry through
  claims-draft before it is built.

## To do, in order

1. The owner rules 6 to 18, the reference syntax first.
2. Draft the CH entries (claims-draft): docs in charter.toml, references
   and their check, embedding, curation and hashes, staleness in lint,
   backlinks in anchors, export and diff.
3. Build them (claims-implement), with the examples: ratelimit gains a
   short walkthrough of how a request spends a token.
4. Write the `claims-narrate` skill and add docs to `claims-audit`.
5. Prove it on torii: `docs/design.md` gains checked references and a
   curation block, and `claims-narrate` drafts a walkthrough of the
   proxy's connection path for the owner to curate.
