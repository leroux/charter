# 008. A clean slate: what is retired is deleted, and the tooling keeps only what it uses

Status: done · Size: M · See: 006, 007

The repository held what no longer ran or ruled: a viewer and the export
only it read, benchmarks seeded in a form the rewrite removed, nine
deprecated entries, answered questions, a held story about entries since
deprecated, and a lint whose advice was mostly style. Each could be read
as if it were current. Git keeps all of it.

## Ruling

Ruled by the owner, 2026-11-27.

| # | Question | Ruling |
|---|---|---|
| 1 | The viewer? | Deleted, claims-viewer/ and its build; it may be remade later. |
| 2 | The export? | Deleted with the viewer, its one reader (CH-261015-9ZTY). `charter anchors` gives an entry's Changes, which only export gave. |
| 3 | The deprecated entries? | Deleted now. From here on an entry is retired in two changes: deprecated, its code and tests stripped, then deleted; `charter diff` accepts the deletion of an entry deprecated at its base. |
| 4 | Lint? | Deleted, verb and package, with CH-041 and the style entries CH-261126-9K5J and CH-261126-WTMZ. Three of its checks stay, as warnings `charter check` writes to standard error and never fails on: a planned entry fully tested, a deprecated entry still anchored or cited, and a doc drafted or stale (CH-261127-QRQS). |
| 5 | The benchmarks? | bench/ deleted whole, its seeds being in the removed go-test form. |
| 6 | Answered questions? | Deleted once their answer is in the record (Q-261015-YSKC, yes); CH-261009-4B3B changed to match. |
| 7 | Old Changes bullets? | The bullets that only said a title or scenario was reworded to the style's bars are deleted (Q-261015-X96S, yes); a Changes line says what changed in the ruling and why, in about three lines. |
| 8 | Backlog 005? | Deleted; its ideas live in 006. |
| 9 | Q-261127-BVQX, Q-261127-TW9H | Yes, both, as the agent proposed. |
| 10 | Q-261015-CPW5 | Already answered and built (backlog 003, ruling 28); nothing to do. |

## Done

All of the above, 2026-11-27, with the charter's gate green and every
active scenario passing. The long Changes bullets of earlier entries were
not shortened to three lines; ruling 7 applies to lines written from now
on.
