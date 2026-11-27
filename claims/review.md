# Reviewing a change

Code: internal/report, internal/digest

This area answers the problem that reading every line of a change, an
agent's above all, does not scale. It covers what a reviewer reads
instead: which scenarios passed, and the digest of the whole record.

## CH-042 charter report MUST join each test's markers to its scenarios and say which passed, failed, were skipped or have no test

Status: active
Serves: CH-001
See: CH-024

- The report MUST set CHARTER_EVENTS to a file of its own, run the test
  command of each language charter.toml names in that language's
  directory, and read the lines the markers wrote.
- It MUST write, per entry and per scenario, the citing tests and their
  result: passed when a citing test passed and none failed; failed when
  one failed; skipped when every one was skipped; none when no citing
  test ran.
- A command that cannot start, or fails having written no line, MUST
  fail the report naming the language and its output.
- An active scenario, or an active entry without scenarios, that did not
  pass is unverified, and the report MUST fail on any.
- Each scenario MUST be shown with who ruled it, and the active scenarios
  MUST be counted by it: the owner's, a delegate's, an agent's
  provisional. A markdown rendering opens with that count, lists the
  unverified ones first, then every entry.

Why: the gate proves that every scenario is cited, but it cannot see
whether the citing tests ran or passed. The report runs them and joins
their outcomes to the scenarios, which closes the loop between the record
and the tests.

Scenarios:
1. `results-joined` (owner, backlog 003): given runs of Go, Python and
   TypeScript tests, when the report is made, then each scenario has its
   citing tests and their outcomes, folded into one outcome.
2. `skipped-scenario-fails`: given an active scenario whose every citing
   test was skipped, or that no test ran, when the report is made, then
   it is unverified; given a planned or a deprecated entry, then it never
   is.
3. `markdown-rendered`: given a report, when it is rendered as markdown,
   then it lists the unverified scenarios first, then each entry with its
   scenarios, their tests and their outcomes.
4. `provenance-and-profile` (owner): given scenarios ruled by different
   people, when the report is made, then each scenario's provenance is in
   it and the active scenarios are counted by who ruled them, passed and
   all.
5. `command-failure-named` (owner, backlog 003): given a test command that
   cannot start, or fails having written no line, when the report runs,
   then it fails naming the language and the output.

Changes:
- 2026-10-09: the report shows who ruled each scenario and counts the
  active ones by it.
- 2026-11-26: the report runs each language's test command from
  charter.toml and reads the events file every marker writes, where it ran
  the Go suite and read go test's JSON stream, so TypeScript and Python
  are joined too: `go-results-joined` and `typescript-joined` are dropped
  for `results-joined`, `command-failure-named` is added, and the Guidance
  on the stream goes (backlog 003).

## CH-261005-6ZBS charter gen MUST write a digest that reads the charter in order: needs by area, an area map, then each area's decisions with their Why

Status: active
Serves: CH-001
See: CH-040, CH-010

- The digest MUST be one generated markdown file.
- It MUST list each active need under the area whose active decisions
  most often reach it through what they serve, the area earlier in the
  reading order on a tie, and list apart the needs no active decision
  reaches.
- A map MUST give each area in the reading order with the first sentence
  of its header and its count of active and planned decisions, then, for
  each need, the areas that serve it with how many decisions.
- Then, area by area, it MUST give each active decision in the order of
  its document: its id linked to its heading, its title, its Kind, and the
  first sentence of its Why. Planned and deprecated ones are counted, not
  listed.
- The reading order is a file listing the area documents one per line;
  one it leaves out, or a line that names none or names one twice, MUST
  fail the gate.
- `charter gen` writes the digest beside the generated symbols, and
  `charter check` MUST fail a digest that is missing or differs from the
  documents, naming `charter gen`.

Why: the charter says what the system does in sentences a person can
read, but spread over many documents in the order they were written, so
a reader had no way through it except the code. One generated document,
in an order that teaches the system and kept current by the gate, gives
the reader that way through.

Scenarios:
1. `needs-under-their-area`: given needs reached by several areas'
   decisions, directly or through a decision they serve, when the digest
   is written, then each sits under the area that reaches it most, the
   earlier winning a tie, and one none reaches is listed apart.
2. `area-map`: given a record, when the digest is written, then the map
   lists the areas in reading order with each header's first sentence and
   its decision counts, and each need with the areas serving it and how
   many decisions.
3. `decisions-in-order`: given an area, when the digest is written, then
   its active decisions appear in document order with the id linked to
   the heading, the title, the Kind and the Why's first sentence, and the
   rest are counted, not listed.
4. `reading-order-complete`: given an area document the reading order
   leaves out, or a line naming no area document or one already named,
   when the gate runs, then it fails naming the file and the line.
5. `stale-digest-fails`: given a missing digest, or one that differs from
   the documents, when the gate runs, then it fails naming the command;
   when gen runs, then it writes the digest.

Changes:
- 2026-11-26: charter gen is the one writer of the digest, so
  CHARTER_UPDATE_DIGEST and its Accounts are gone (Q-261015-Y682).
