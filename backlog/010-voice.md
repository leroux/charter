# 010. Replace charter's house style with conventional technical prose

Status: held — discuss with the owner before any work · Size: L · See: 008, 009

The owner dislikes the voice charter's record, skills and docs are written
in, and wants it changed. Nothing in this story is decided yet: the first
step is a conversation about what the replacement voice is, and only then
a plan.

## What the current style is

It is not written down anywhere today. It used to be enforced in part by
lint (CH-261126-9K5J: titles at most 20 words, an entry's text at most six
sentences, a scenario sentence at most 30 words; CH-261126-WTMZ: a list of
retired words), both deleted in backlog 008. What remains is convention,
copied from one entry to the next and by agents imitating what they read.
Its traits, with examples from the repository:

- **Telegraphic constructions that drop articles and verbs.** "One doc a
  purpose and a reader", "one doc a flow a newcomer must understand", "a
  word a line", "the stories in progress or held, one a file".
- **Plain words preferred over technical terms**, even where the term is
  more precise. This is what produced "the name" for a DNS query and "the
  way in" for inbound in torii's first doc drafts (fixed in
  CH-261126-5JMD's 2026-11-27 change, for docs only).
- **Coined nouns for the method's own concepts**, used without definition
  at the point of use: the record, the gate, the judge, the loop, a need,
  an area, a decision, curated, drafted, carried.
- **Long comma chains in place of lists, colons or conjunctions.** "the
  syntax, planning the set, drafting, converting an existing doc".
- **Every title is a full declarative sentence**, often a long one, and
  becomes the heading's anchor (`#sw-261127-6s2f-under-an-allowlist-a-tls-…`).
- **Terse, unhedged imperatives** in skills and CLAUDE.md, with little
  explanation of why.

Effects seen so far: entries and skills are slow to read and easy to
misread; precision is lost where a plain word replaces a term; and agents
reproduce the style everywhere they write: docs, commit messages, pull
requests, chat.

## Scope to settle

- charter's record: `claims/*.md` (titles, texts, Whys, scenario sentences,
  area headers, Changes lines).
- The skills under `plugin/skills/`, which teach agents the voice.
- `README.md`, `CLAUDE.md`, the backlog, and the examples' records and docs.
- Commit messages and pull request descriptions.
- Adopters: torii's record was written in the same style; whether it
  follows, and when.

## Questions for the discussion

1. What voice replaces it? For example: conventional technical English as
   in Go or Rust documentation; RFC style with MUST and SHOULD for the
   binding text; something else.
2. Titles: a full sentence, a short noun phrase, or a short imperative?
   Note that a title is also a heading anchor, and docs link to those
   anchors.
3. Do length bars come back, as guidance or as a check?
4. Terminology: keep the method's coined nouns with a glossary, or rename
   them to established terms (spec, requirement, acceptance test)?
5. Rewording an entry without changing its ruling: does it take a dated
   Changes line? Backlog 008 deleted Changes bullets that only recorded
   rewording.
6. Migration: one pass over the whole record, or each entry when it is
   next touched? Who reviews the rewritten text before it lands?
7. How the skills teach the new voice: a short style guide that every
   skill links to, examples, or an audit check.
8. Commit messages and pull requests: does the style apply to them, and
   in what form?

## To do

1. Discuss the questions with the owner and record the answers as a
   ruling table here.
2. Plan the change from the rulings and propose it before editing.
