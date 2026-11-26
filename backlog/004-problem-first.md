# 004. The charter's own record is told problem by problem: one need, one area for each problem the charter solves

Status: drafted · Size: M · See: 003, 001

The charter is a higher level to build software at: what the software must
do is written above the code, the code is held to it, and a change is
reviewed through it, by people and agents alike. Its own record was grouped
by the tool's internals, two needs carried 25 of its 35 decisions, and four
of its needs were properties of the record rather than promises to a
person. This story re-cuts the needs and the areas around the problems.

## Ruling

Ruled by the owner, 2026-11-26.

| # | Question | Ruling |
|---|---|---|
| 1 | What problems does the charter solve? | P1, intent has no durable home; P2, code drifts from intent unnoticed; P3, reviewing code line by line does not scale; P4, who decided what is invisible; P5, a newcomer, person or agent, cannot find their way in; P6, every project is different. |
| 2 | How are the needs cut? | One a problem, each a person's: CH-261126-PFG5 (P1, new), CH-261126-P2SK (P2, new), CH-001 (P3), CH-261009-3FHW (P4), CH-002 (P5), CH-261126-NS7T (P6). |
| 3 | What of the needs that were properties? | CH-003, CH-004, CH-005 and CH-006 become decisions serving the need of their problem, ids kept, each with a Changes line. |
| 4 | How are the areas cut? | One a problem, in that order: intent.md, anchoring.md, gate.md, review.md, provenance.md, adoption.md; form.md, tooling.md, config.md and languages.md are folded into them. Entries move verbatim; each Code line names the packages that exist. |
| 5 | Whose is a need? | The need's title names its person, as "A reviewer can…" does; no field. CH-261126-G5BX, drafted for a For field, is deprecated unbuilt. |

## Charter

Drafted 2026-11-26 for the owner's review.

- New, planned: CH-261126-PFG5, CH-261126-P2SK (needs). CH-261126-G5BX (a For field on needs) was drafted and is deprecated.
- Needs made decisions, each with a Changes line and anchored to the code that already carried it: CH-003 (record), CH-004 (gate's surfaces), CH-005 (gate's observations), CH-006 (the walk).
- Serves re-pointed, each with a Changes line: CH-010, CH-011, CH-030, CH-031.
- Every other entry moved verbatim to its problem's area.

## For the review

- The new needs' scenarios are `(agent, provisional)`: the problems are the
  owner's, their wording is the agent's.
- Lint names six tests shared by two decisions: CH-003, CH-004 and CH-005,
  now decisions, restate clauses their children rule (CH-012, CH-033,
  CH-034, CH-035, CH-040). Either they become summaries with no scenarios
  of their own, or their scenarios move to the children.
- The bundles CH-010, CH-023 and CH-041 are not split yet.

## To do

1. The owner reviews the needs and the areas.
2. Settle CH-003, CH-004 and CH-005's overlap with their children.
3. Verify CH-261126-PFG5 and CH-261126-P2SK from the command's side and set them active.
