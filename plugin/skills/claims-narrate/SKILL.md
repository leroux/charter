---
name: claims-narrate
description: Write the docs of a charter, the prose over the record and the code that a newcomer reads first, with checked references to entries and declarations and embeds of the entries they rest on. Plans a doc set (an index, an architecture, walkthroughs of the flows, design docs, guides), drafts each doc, converts an existing prose doc to checked references, and on a doc charter check warns is stale proposes an edit of what changed. Never curates and never writes over a curated doc's prose; a person reads the draft and runs charter curate. Use when asked to "write the docs for X", "document the architecture", "write a walkthrough of Y", "turn docs/design.md into a charter doc", or "check says docs/Z.md is stale".
# charter:realizes CH-261126-5JMD
---

# Narrate

A doc is the part of a charter a newcomer reads first: how the system
works, where in the code, why it is shaped so, and which decisions each
part carries. The record holds what binds; a doc tells the story that
joins it, and its prose is free. The gate checks a doc's references, and
`charter check` warns when anything a doc references has changed since a
person last read it. `claims/charter.toml` says where everything is: the
`docs` patterns name the docs, and the language tables the code.

The agent drafts; a person curates. This skill never runs `charter
curate` and never writes over a curated doc's prose.

## The syntax

- **An entry:** a link whose text opens with its id, `[SW-261126-VX8Z](x)`.
  `charter gen` writes the target and, when text follows the id, the
  current title; a bare id stays bare, for a link in running prose.
- **An embed:** `<!-- charter:embed SW-… -->` and `<!-- /charter:embed -->`
  on lines of their own, nothing between; gen writes the entry's id, title
  and text as a quote. Use it where a doc rests on the rule itself.
- **A declaration:** a link whose text is a name in backticks and whose
  target is the file declaring it, ``[`Limiter.Allow`](../limiter/limiter.go)``.
  The name is bare, or qualified by its package (in Python and TypeScript,
  the file's name) or by its type for a method. gen keeps a `#L` line
  anchor on it. A link to a file no language reads only has to exist.
- **Not read:** a link whose text is neither, anything in a code block or
  a code span, and anything outside the root. A heading that opens with an
  id fails: a doc declares no entry.
- **Curation:** `charter curate -by NAME FILE` writes, at the doc's end, who
  read it, the date and a hash of each reference. A doc with no block is
  drafted; one whose references changed since is stale. Both are warnings.

## Set up the docs

If `charter.toml` has no `docs` key, propose one, such as `docs =
["docs/*.md", "docs/**/*.md"]`, and say that the patterns are their own
scope: `ignore` and `only` do not narrow them, and only markdown is read.
A directory whose name begins with a dot or an underscore is never walked.

## Plan a doc set

Before writing any doc, propose the set to the person and wait.

1. **Take stock.** Read `claims/needs.md`, every area document's header
   and titles (`claims/reading-order.txt` gives their order), the
   observations, the README, and every prose doc the repository already
   has. List what each existing doc covers and whether it is current.
2. **Choose the docs.** One doc a purpose and a reader, from these kinds:
   - **Index** (`docs/README.md`): what the system is in a paragraph, who
     each doc is for, the order to read them in. Every other doc is
     reached from it.
   - **Architecture:** the parts (processes, packages, machines, files on
     disk), what runs where, what crosses each boundary, and the trust
     boundaries. One picture in text or a diagram block, then a section a
     part, each naming its declarations and the areas whose decisions it
     carries.
   - **Walkthrough:** one flow followed through the code, from where the
     data enters to where it leaves: a request, a packet, a command run
     end to end. One doc a flow a newcomer must understand.
   - **Design:** why the system is shaped so: the problem, the choice,
     the alternatives rejected and why, and what the choice led to, as
     references to the entries it produced. A design choice no entry
     rules stays prose; say it is how the code does it today.
   - **Guide:** how a person does one task: install, configure, operate,
     debug, add a part. Steps, each naming the command or file, and the
     entries that rule what the person sees.
   Not a kind: a reference of the rules. The record is the reference; a
   doc that lists every entry is a second copy.
3. **Give each fact one home.** The rule's home is its entry; a doc
   references or embeds it. A part is described once, in the
   architecture; a walkthrough or a guide links there. Docs link to each
   other with plain links.
4. **Cover the record.** Every need is reached from the index, and every
   area by at least one doc; an area no doc reaches is named in the plan
   as left out, with the reason.
5. **Propose the plan** as a table: file, kind, reader, one line of
   purpose, the areas and needs it covers, and for each existing doc
   whether it is kept, converted or replaced. Write nothing until the
   person agrees.

## Draft a doc

1. **Read what it covers.** The entries (`charter anchors ID` gives each
   one's code, tests and the docs that mention it already) and the code
   their anchors point at. Read the declarations, not only the entries.
2. **Follow the data.** Order a walkthrough by what happens to the data,
   an architecture by the parts, a design by the choices, a guide by the
   person's steps: never by package or by entry id. Each claim about the
   code names the declaration as a code reference, and the entry it
   carries as a reference.
3. **Write for the reader the plan names.** The record's register, short
   plain words, is the record's; a doc speaks its reader's language. For
   engineers and security reviewers that is the field's own vocabulary
   and the exact mechanism:
   - The term, never a plainer word standing in for it: a DNS query and
     its A/AAAA answer, not "the name"; inbound, not "the way in";
     `connect(2)` or the dial, named as the code names it; SNI, the
     ClientHello, TPROXY, conntrack, a hook and a chain.
   - The mechanism and its values: the port, the socket, the file, the
     feature key, the timeout and the limit, the error code returned, as
     the code has them; what happens at the boundary of each.
   - The security vocabulary where it applies: the adversary, the trust
     boundary and what authenticates across it, fail-closed, spoofing,
     rebinding, exfiltration, the residual risk.
   - The system's own coined terms only as the record defines them, each
     defined once with a link; where the code names a thing otherwise,
     give both once.
   - Headings that say what the section is about ("DNS resolution",
     "Inbound"), and no sentence that only announces the next one.
   A doc for another reader (an operator new to the field, a manager) is
   another doc, planned as such.
4. **Embed what binds.** Where the story rests on a rule, embed the entry
   rather than paraphrase it. Paraphrase only what no entry states.
5. **Name the traps.** What a newcomer would get wrong, and why the code
   is shaped against it: a lock held for a whole request, a clock read in
   one place, an observation of an outside system the code rests on.
6. **Check.** Run `charter gen` (targets, line anchors, embeds), then
   `charter check`: every reference resolves. Read each sentence against
   the declaration it names; one you cannot point at a line for, cut. A
   sentence that contradicts an entry is a question for the person, not a
   quiet fix in either place.
7. **Hand it over.** Present the doc, say it is drafted, and list what to
   read it against. `charter check` warns of it until a person reads it
   and runs `charter curate -by NAME FILE`.

## Convert an existing doc

A prose doc the repository already has (a design doc, an architecture
note) keeps its author's prose and gains checked references.

1. Read it whole, then the entries and code it talks about.
2. Where it states a rule an entry rules, link the entry, or embed it
   where the doc rests on it; where it names code, make the name a code
   reference.
3. Where it says what the record or the code does not, do not change
   either: list each such sentence for the person, with the entry or the
   declaration it disagrees with. Their answer is a change to the doc, a
   change to the record through a story, or a question.
4. Keep its rationale, its history and its voice; cut nothing the record
   does not already say better. Then check and hand it over as above.

## Refresh a stale doc

1. `charter check` warns of the doc and each reference that changed since
   it was curated; `charter diff REV` lists the docs a change touched.
2. For each changed reference, read what changed (the entry's text and
   its Changes lines; the declaration's diff) and the doc's sentences
   around the reference.
3. Propose an edit of those sentences alone, as a diff for the person;
   leave the rest of the doc as it is. A change the doc does not need is
   said so, with the reason: the person curates again either way.
4. Never run curate; the hash is the person's word that they read it.

## Done when

- The plan was agreed before the docs were written.
- Every doc's references resolve and `charter check` passes.
- Every sentence about the code names a declaration it can be checked
  against, and every rule a doc rests on is a reference or an embed.
- Every disagreement found between prose, record and code is listed for
  the person.
- The person has the drafts, or the proposed edits, and curates them.
