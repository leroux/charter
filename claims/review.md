# Reviewing a change

Code: internal/report, internal/digest

The problem this area answers: reading every line of a change, an agent's
above all, does not scale. It holds what a reviewer reads instead: which
scenarios passed, and the digest of the whole record.

## CH-042 A report joins each test's markers to its scenarios and says which passed, failed, were skipped or have none

Status: active
Serves: CH-001
See: CH-024

The report sets CHARTER_EVENTS to a file of its own, runs the test command
of each language charter.toml names in that language's directory, and
reads the lines the markers wrote. It writes, per entry and per scenario,
the citing tests and their result: passed when a citing test passed and
none failed, failed when one failed, skipped when every one was skipped,
none when no citing test ran. A command that cannot start, or fails having
written no line, fails the report naming the language and its output. An
active scenario, or an active entry without scenarios, that did not pass
is unverified, and the report fails on any. Each scenario is shown with
who ruled it, and the active scenarios are counted by it: the owner's, a
delegate's, an agent's provisional. A markdown rendering opens with that
count, lists the unverified ones first, then every entry.

Why: the gate cannot see whether the tests ran; the report is the sensor
that closes the loop.

Scenarios:
1. `results-joined` (owner, backlog 003): runs of Go, Python and
   TypeScript tests yield, per scenario, the citing tests and their
   outcomes, folded into one outcome.
2. `skipped-scenario-fails`: an active scenario whose every citing test
   was skipped, or that no test ran, is unverified; a planned or a
   deprecated entry is never.
3. `markdown-rendered`: the markdown lists the unverified scenarios first,
   then each entry with its scenarios, their tests and their outcomes.
4. `provenance-and-profile` (owner): each scenario's provenance is in the
   report, and the active scenarios are counted by who ruled them, passed
   and all.
5. `command-failure-named` (owner, backlog 003): a test command that
   cannot start, or fails having written no line, fails the report naming
   the language and the output.

Changes:
- 2026-10-09: the report shows who ruled each scenario and counts the
  active ones by it.
- 2026-11-26: the report runs each language's test command from
  charter.toml and reads the events file every marker writes, where it ran
  the Go suite and read go test's JSON stream, so TypeScript and Python
  are joined too: `go-results-joined` and `typescript-joined` are dropped
  for `results-joined`, `command-failure-named` is added, and the Guidance
  on the stream goes (backlog 003).

## CH-261005-6ZBS The digest reads the charter in order: needs by area, an area map, then each area's decisions with their why

Status: active
Serves: CH-001
See: CH-040, CH-010

The digest is one generated markdown file. It lists each active need
under the area whose active decisions most often reach it through what
they serve, the area earlier in the reading order on a tie, and the
needs no active decision reaches apart. A map gives each area in the
reading order with the first sentence of its header and its count of
active and planned decisions, then, for each need, the areas that serve
it with how many decisions. Then, area by area, each active decision in
the order of its document: its id linked to its heading, its title, its
Kind, and the first sentence of its Why; planned and deprecated ones are
counted, not listed. The reading order is a file listing the area
documents one a line; one it leaves out, or a line that names none or
names one twice, fails the gate. charter gen writes the digest beside the
generated symbols, and charter check fails a digest that is missing or
differs from the documents, naming charter gen.

Why: the charter says what the system does in sentences a person can
read, but in 24 documents in the order they were written, so a reader
had no way through it but the code. One document in an order that
teaches it, kept current by the gate, is that way.

Scenarios:
1. `needs-under-their-area`: a need sits under the area whose decisions
   reach it most, directly or through a decision they serve, the earlier
   winning a tie; one none reaches is listed apart.
2. `area-map`: the map lists the areas in reading order with each
   header's first sentence and its decision counts, and each need with the
   areas serving it and how many decisions.
3. `decisions-in-order`: each area's active decisions appear in document
   order with the id linked to the heading, the title, the Kind and the
   Why's first sentence; the rest are counted, not listed.
4. `reading-order-complete`: an area document the reading order leaves
   out, and a line naming no area document or one already named, fail the
   gate naming the file and the line.
5. `stale-digest-fails`: a missing digest, or one that differs from the
   documents, fails the gate naming the command; generation writes it.

Changes:
- 2026-11-26: charter gen is the one writer of the digest, so
  CHARTER_UPDATE_DIGEST and its Accounts are gone (Q-261015-Y682).
