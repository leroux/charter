---
name: claims-review
description: The judge. Review a change's conformance to the charter, over a branch or a pull request, in a context that did not write it. For each scenario the change touched, answer three questions with evidence: does the citing test show what the scenario's sentence says, does the anchored code carry the decision, and does the Changes line say what changed and why. Each no is a finding; claims/policy.md says whether it blocks or warns. Runs as the last step of claims-implement, by hand over a branch, and in CI on every pull request. Read-only; it changes nothing and its verdicts are for the author and the owner.
---

# The judge

The gate proves the links exist: an anchor names the entry, a test cites
each scenario. The judge says whether they are true. It runs in a context
that did not write the change: in CI, its own job; by hand, a fresh
agent or a person who has not read the implementation. It is read-only.
It changes no anchor, test or entry; it posts its verdicts.

## Scope

The scenarios the change touched, and only those:

- scenarios of entries added to an area document or `claims/needs.md`;
- scenarios added, reworded or given a new provenance;
- scenarios whose citing tests or whose entry's anchored code the change
  edits: `git diff BASE... | grep -n 'charter:\|Verifies(\|verifies('`
  and the files it touches, looked up with `go run ./cmd/claimsgen
  anchors ID...`;
- every entry `go run ./cmd/claimsgen diff BASE` names.

Over a pull request, `gh pr diff N` is the change and the base is the
pull request's base branch; by hand, `main`.

## The three questions

For each scenario in scope, `claimsgen anchors ID` gives the sentence,
the realizing code with its declarations, and the citing tests with the
handles each names. Read, then answer each question yes or no, with the
evidence that decided it: the file and line of the assertion, the
declaration, the Changes line.

1. **Does the test show what the sentence says?** Read the scenario's
   sentence, then the test from its citation down. Yes when an assertion
   depends on the sentence's outcome and would fail if the decision were
   violated. No when nothing asserted depends on it, when the test would
   pass against a stubbed implementation, when a bound is checked against
   the constant it bounds rather than a number, when the test exercises
   another scenario than the one it names, or when a need is verified from
   inside a package rather than from the person's side.
2. **Does the anchored code carry the decision?** Read the decision's text
   and Why, then each realizes declaration. Yes when that declaration
   decides what the text says. No when it is a wrapper, a type or a
   constant beside the code that decides, when it implements a narrower or
   another case than the text, or when the code that decides is anchored
   nowhere.
3. **Does the Changes line say what changed and why?** For an entry the
   change edited. Yes when the dated line names the change in place and
   its reason. No when there is no line, when it points at a review or a
   story instead of saying, or when it says what and not why. An entry
   `claimsgen diff BASE` names has no line to read: that is a no.

Also read the provenance. A scenario written as the owner's or a
delegate's whose Why or Changes cite no ruling is a finding: the record
claims a ruling nobody made.

## Verdicts

`claims/policy.md` says, by class (the decision's Kind and its area),
whether a no blocks or warns. Write one line per finding:

    Verdict: block  SW-236/queue-holds-sixty-four  Q1 no: queue/queue_test.go:45 loops over QueueLimit; with the limit at 16 the test passes (bound, queue)
    Verdict: warn   SW-178/panel-at-640  Q2 no: the anchored Panel type decides nothing; the width is in layout.ts:88, unanchored (web)

then a summary: how many scenarios were read, how many findings, how many
block. A finding you cannot point at is a guess: trace it or drop it.

## In CI

Write the verdicts and the summary to `charter-review.md` in the working
directory, then post the same text as one comment on the pull request
with `gh pr comment`. The job fails when any line opens with
`Verdict: block`; the author fixes, or the owner overrules on the thread
and the finding's class is read again on the next push.

## Not a finding

- Code left unanchored because its entry lands in a later change.
- Code that is wrong while faithfully implementing its entry; that is the
  code review's finding.
- Guidance departed from; nothing verifies it.
- A provisional scenario: that is debt the report counts, not drift.

## Done when

Every scenario in scope had its three questions answered with evidence;
every finding has its verdict and its class; the result is presented, or
written to `charter-review.md` and posted. Nothing re-anchored, nothing
changed, nothing merged.
