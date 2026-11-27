---
name: claims-audit
description: The audit. Read the whole record, or the part of it that changed since the last audit, against its tests and its code in a context that did not write it, and report what has drifted: a scenario sentence that claims more or less than its test checks, a text that names a boundary the code lacks, a clause no scenario tests, a clause two entries state, a doc that says what its code does not, and the style the reader stumbles on. Writes claims/audit.md and files questions; changes nothing else. Runs by hand. Done when the correctness findings are none and charter check passes.
---

# The audit

The gate proves the links exist. The judge (`claims-review`) reads what a
change touched. The audit reads what nobody touched: the record as it
stands, entry by entry, against the tests that cite it and the code that
carries it, so that drift that no pull request introduced in one piece
is still found. It is read-only. It writes one file, `claims/audit.md`,
and files questions; it edits no entry, test or line of code. What it
finds is built from a story, under `claims-implement`.

## When

- Over the whole record, now and then.
- After a pass that rewrote many entries, or when the judge's
  findings on pull requests start repeating.
- Narrow, by hand, over the entries whose documents, anchored code or
  citing tests changed since the last audit: `git log -1 --format=%H --
  claims/audit.md` gives the last audit's commit; `go run
  ./charter diff THAT` names the entries whose text changed,
  and `git diff --name-only THAT..HEAD` matched against `charter
  anchors` names the rest.

## Step 1: the mechanical checks

Run the gate first; its problems and its warnings are the audit's first
findings and cost no reader:

    charter check

Copy its lines into the report by kind. A record whose check prints
nothing, and warns of nothing, has no mechanical finding.

## Step 2: the readers

Split the live entries among fresh readers, one or two documents each,
as `claims/reading-order.txt` groups them; a reader reads nothing of how
the record was written. For each live entry, with `go run
./charter anchors ID` open:

1. **Each scenario against its test.** Read the sentence, then every
   citing test from its citation down. A sentence that claims more than
   the test checks is `SENTENCE-OVER-TEST ID/handle: file:line`; one that
   claims less than the test checks is `SENTENCE-UNDER-TEST`.
2. **The text against the code.** Read the realizing declarations. A text
   that names a boundary the code does not have, or a behaviour the code
   does not do, is `TEXT-OVER-CODE ID: file:line`.
3. **Each clause against the scenarios.** A clause of the text no scenario
   tests, and no waiver names, is `UNTESTED-CLAUSE ID: "the clause"`,
   unless it is tuning that belongs in Guidance, which is `STYLE`.
4. **The clause against the record.** A clause another live entry also
   states (search `claims/DIGEST.md` for its key words) is `TWO-HOMES ID,
   ID: "the clause"`.
5. **Style.** A Why that gives no reason though it is
   long, a title that is a noun phrase, a text that narrates, a word the
   glossary defines otherwise: `STYLE ID: what`. Order against the
   header: `ORDER DOC: what`.

6. **Each doc against its code.** For each doc charter.toml's `docs`
   names, read it with every declaration it references open. A sentence
   that says the code does what it does not is `DOC-OVER-CODE doc:line:
   what`; an embed or a reference whose sentence around it contradicts
   the entry is `DOC-OVER-RECORD doc:line ID`. A doc check warns of as
   drafted or stale is read first.

A finding a reader cannot point at a file and line is a guess: trace it
or drop it. Readers report in these forms, most severe first, and the
auditor merges them.

## Step 3: the report

`claims/audit.md` is overwritten by each audit; git holds the earlier
ones. Its form:

    # Audit of YYYY-MM-DD

    Scope: the whole record at COMMIT | the N entries changed since COMMIT.
    Readers: N, over N live entries; tests opened for N scenarios, code
    read for N entries.
    Lint: N retired words, N titles, N texts, N scenarios, N Whys, N shared
    tests, N open questions.

    Done: yes | no.

    ## Correctness
    SENTENCE-OVER-TEST ID/handle: file:line, what the sentence claims
    that the test does not check.
    ...

    ## Two homes and untested clauses
    ...

    ## Style and order
    The counts by kind, then the warnings check printed and the readers'.

    ## Questions filed
    Q-IDs with their entries.

`Done: yes` when the Correctness section is empty and charter check
passes. Style and untested clauses do not keep an audit from being done; they are counted, and the
owner sets a bar in `claims/policy.md` when the counts should fail.

## Step 4: questions, never changes

A finding that needs a ruling (which of two homes, whether a clause is
tested or dropped, what the code should do) is a question in
`claims/questions.md` with the provisional answer "as it stands", filed
with `charter next Q`. Nothing else changes: the audit does not reword, move or
fix. The fixes are a story, built under `claims-implement`, and the next
audit says whether they held.

## Done when

- `claims/audit.md` is written in the form above, dated, with its scope.
- Every finding points at a file and line, or at the clause it quotes.
- The questions a ruling needs are filed, or listed in the issue.
- Nothing else in the repository changed.
