---
name: claims-audit
description: The audit. Read the whole record, or the part of it that changed since the last audit, against its tests and its code in a context that did not write it, and report what has drifted: a scenario sentence that claims more or less than its test checks, a text that names a boundary the code lacks, a clause no scenario tests, a clause two entries state, and the style the reader stumbles on. Writes claims/audit.md and files questions; changes nothing else. Runs by hand, and monthly in CI. Done when the correctness findings are none and lint prints only the line counting the open questions.
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

- Monthly, in CI (`.github/workflows/charter.yml`, the audit job), over
  the whole record, its report opened as an issue.
- By hand, after a pass that rewrote many entries, or when the judge's
  findings on pull requests start repeating.
- Narrow, by hand, over the entries whose documents, anchored code or
  citing tests changed since the last audit: `git log -1 --format=%H --
  claims/audit.md` gives the last audit's commit; `go run
  ./cmd/claimsgen diff THAT` names the entries whose text changed,
  and `git diff --name-only THAT..HEAD` matched against `claimsgen
  anchors` names the rest.

## Step 1: the mechanical checks

Run the gate and lint first; their lines are the audit's first findings
and cost no reader:

    go test -count=1 ./claims ./claims
    go run ./cmd/claimsgen lint

Lint names the retired words (`claims/retired.txt`), the titles over
twenty words, the texts over six sentences, the scenario sentences over
thirty words, the Whys that only point at a ruling, the tests two
decisions share, and counts the open questions. Copy its lines into the
report by kind. A record whose lint prints nothing but the questions
count has no mechanical finding.

## Step 2: the readers

Split the live entries among fresh readers, one or two documents each,
as `claims/reading-order.txt` groups them; a reader reads nothing of how
the record was written. For each live entry, with `go run
./cmd/claimsgen anchors ID` open:

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
5. **Style lint cannot see.** A Why that gives no reason though it is
   long, a title that is a noun phrase, a text that narrates, a word the
   glossary defines otherwise: `STYLE ID: what`. Order against the
   header: `ORDER DOC: what`.

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
    The counts by kind, then the lines lint printed and the readers'.

    ## Questions filed
    Q-IDs with their entries, or, in CI, the questions as text.

`Done: yes` when the Correctness section is empty and lint printed
nothing but the line counting the open questions. Style and untested
clauses do not keep an audit from being done; they are counted, and the
owner sets a bar in `claims/policy.md` when the counts should fail.

## Step 4: questions, never changes

A finding that needs a ruling (which of two homes, whether a clause is
tested or dropped, what the code should do) is a question in
`claims/questions.md` with the provisional answer "as it stands", filed
by hand with `claimsgen next Q`; in CI it is listed in the issue for a
person to file. Nothing else changes: the audit does not reword, move or
fix. The fixes are a story, built under `claims-implement`, and the next
audit says whether they held.

## Done when

- `claims/audit.md` is written in the form above, dated, with its scope.
- Every finding points at a file and line, or at the clause it quotes.
- The questions a ruling needs are filed, or listed in the issue.
- Nothing else in the repository changed.
