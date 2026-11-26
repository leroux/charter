# 005. The record splits into a contract any implementation must keep and the design of this one, and the contract's scenarios run against any implementation

Status: held · Size: L · See: 004, 002, 006

The record should be enough for someone to rebuild the software from it,
with other code and in another language, and to know the rebuild is right
by running its tests. That needs two things the record lacks: a line
between what binds every implementation (the contract) and what this one
chose (the design), and tests of the contract that are not written in the
implementation's language.

> Held, 2026-11-26: the owner put the foundations first. CH-261126-PZ8Z and
> CH-261126-GFR7 are deprecated unbuilt, and the idea waits in backlog 006.
> The map below stays as the analysis it was.

## Ruling

Ruled by the owner, 2026-11-26.

| # | Question | Ruling |
|---|---|---|
| 1 | What is the record for? | Someone can take a project's charter and reimplement the software from it, without its code, and check the result with its contract's tests. |
| 2 | How is contract told from design? | By area: an area's header says `Level: contract` or `Level: design`, contract by default (CH-261126-PZ8Z, planned). |
| 3 | How is the contract tested portably? | By fixtures: txtar files of a tree, a command and the expected output and exit, run by `charter conform` against whatever implementation charter.toml names (CH-261126-GFR7, planned). |
| 4 | How is the contract known to be complete? | Every surface is accounted for by a contract entry; the surface kinds are how a project declares its edge. |

## The map

Each entry of the charter's own record as it stands, by level. "Split"
means part of it is contract and part design; "tests" says where its
citations are today: at the edge (the command) or inside (a package).

### Contract

| Entry | What it rules | Tests today |
|---|---|---|
| CH-001, CH-002, CH-261009-3FHW, CH-261126-NS7T, CH-261126-PFG5, CH-261126-P2SK | The needs: what the charter guarantees | mixed; PFG5 and P2SK none yet |
| CH-010, CH-011, CH-012, CH-013, CH-015 | The form of an entry, needs, ids, status, relations | inside |
| CH-003 | The record is stable | inside |
| CH-261126-HK6S | observations.md's form | inside |
| CH-261009-VY5C, CH-261009-4B3B | Provenance syntax; questions.md and the questions verb | inside; the verb at the edge |
| CH-261126-33GK | charter.toml: found walking up, its keys | inside |
| CH-021, CH-022, CH-023 | How code anchors and tests cite, per language | inside |
| CH-020 | The generated files: names, symbols, self-contained | inside |
| CH-024 | The events file's lines | inside, with real runs |
| CH-025 | What the claims directory may hold | inside |
| CH-004, CH-030 to CH-035, CH-261126-NTRR, CH-005 | The gate's rules | inside |
| CH-036 | What the walk reads and skips | inside |
| CH-037 | check's output and exit status | edge |
| CH-040 | The verbs | edge |
| CH-041, CH-042, CH-261005-6ZBS, CH-261009-QQC8 | lint, report, digest, diff | inside |
| CH-261015-9ZTY (planned), CH-261126-5JMD (planned) | export; the skills plugin | edge, partly |

### Design

| Entry | What it rules |
|---|---|
| CH-006 | The charter specifies itself, and its own gate walks the tooling alone: a property of this project |

### Split

| Entry | Contract | Design |
|---|---|---|
| CH-261126-PQ7X | A language charter.toml does not name is not read; a missing parser fails naming the table; a rule says the same in every language | Every language is a reader, a generator and a marker; the gate's rules import no language package |
| CH-024 | Each line is whole, with id, handle, test, file and outcome | One write per line, opened to append |
| CH-020 | The files and their symbols | The order they declare entries in (CH-012's), if a reimplementation may choose |

### Retired

CH-014, CH-261009-SNX2, CH-261005-VFJR, CH-261126-G5BX.

## What the map shows

- The charter's record is almost all contract. Its rule since the start, a
  decision is "how the system behaves at a boundary", kept design out; the
  split will matter more for a product with an architecture of its own:
  torii's one shared proxy qube, or its own tsnet program over tailscaled,
  are design.
- The contract is not portable yet: of 135 citations, 12 are tests of the
  command; the rest call internal packages. Each contract scenario needs a
  fixture, about 180 of them, for a rewrite to be checkable.
- The contract is not known complete: the charter declares no surface kind
  for its own edge. Its verbs and their flags, charter.toml's keys, the
  record's files and the JSON of export and report are each a surface kind
  to declare, and the gate then names what no entry describes.

## To do, in order

1. The owner reviews the two entries and the map, the splits most of all.
2. Build CH-261126-PZ8Z: the Level line, design entries, surfaces on
   contract only; split CH-261126-PQ7X, CH-024 and CH-020.
3. Build CH-261126-GFR7: the fixture format and charter conform.
4. Declare the charter's own edge as surface kinds, and write the entries
   the gate then finds missing.
5. Move the contract's scenarios to fixtures, area by area; the Go tests
   stay as the design's tests and as faster checks.
6. Prove it with backlog 002: give an agent only claims/ and the fixtures,
   have it rebuild the charter in another language, and run charter conform.
