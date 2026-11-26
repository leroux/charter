# The commands

Code: cli, cmd/claimsgen, chartertest

What the charter command does for a person or an agent, on the charter
it finds.

## CH-040 One charter command generates, checks, lints, gives ids and lists anchors, and an adopter writes no main

Status: active
Serves: CH-002
See: CH-261126-33GK

charter gen writes the generated files from the documents alone, so it
works when one is stale or missing. check runs the gate (CH-037), lint
prints the gate's problems and then the advice, report and diff are
CH-042's and CH-261009-QQC8's, and export is CH-261015-9ZTY's. next prints
a new id for a prefix, the day's UTC date and a code drawn at random,
drawn again when the tree already holds that id, or as many as the count
after the prefix asks for, each different. anchors prints, per id, the
document, the title, the status, the scenarios, the realizing and
concerning sites with their declarations, the verifying tests with the
handles each names, and the relations. Every verb works on the charter
CH-261126-33GK finds, and an adopter writes or builds nothing of its own
to run them.

Why: the agent's way in and the person's way to regenerate are one
command; a command an adopter must write a main for is one they write
wrong, or let fall behind the tooling.

Scenarios:
1. `gen-from-documents-alone` (owner, backlog 003): charter gen writes
   every generated file when one of them is missing or stale, reading the
   documents alone.
2. `registry-lists-everything`: the registry gives each entry with
   scenarios, anchors and citations by handle, each area document in
   reading order with title, code and summary, and each waiver with its
   reason.
3. `next-free`: next gives an id of the prefix, today's UTC date and a
   code of the alphabet, for a prefix in use or not, drawing again one the
   tree holds.
4. `next-count`: next with a count gives that many ids, each different,
   and a count that is not a number above zero is refused.
5. `verbs-on-the-found-charter` (owner, backlog 003): run inside an
   adopter's tree, each verb reads the charter found there, with no
   command of the adopter's own.

Changes:
- 2026-10-01: next gives dated ids, where it gave one above the highest,
  and takes a count, so a story's entries are numbered in one call
  (CH-012).
- 2026-10-09: the registry also gives the area documents, in the
  reading order with their headers, and the waivers, so a renderer such
  as the claims viewer walks a reader through any charter in the order
  it declares and shows a scenario as waived without parsing the files
  again.
- 2026-10-09: the title reworded to twenty words or fewer; the ruling is
  unchanged.
- 2026-10-09: the next free id is ruled here alone, the next verb's
  behaviour being this entry's subject: CH-012's next-free restated
  next-free and is dropped, its test citing this one, which now says what
  the dropped one did, the UTC date, a prefix in use or not and the
  redraw, and the text says how next draws.
- 2026-10-09: `registry-lists-everything` reworded under thirty words; the
  ruling is unchanged.
- 2026-11-26: the charter became a repository of its own, meant for other repositories (Q-261015-3Q84). The verbs live in the
  cli package, which a repository's own command calls with the charters
  it keeps; the charter flag selects one of those by name, where it chose
  between the product's charter and the charter's own. The title says any
  charter, and `charter-flag-selects` is added. A command reads a charter
  from the module above the working directory, since the calling file is
  the tooling's, in another module; `root-from-working-directory` is
  added.
- 2026-11-26: the verbs are one charter command, the gate's check and
  claimsmap's export among them, working on the charter it finds above the
  working directory. An adopter writes no main, so the charter flag and a
  root given in options are gone, `charter-flag-selects` and
  `root-from-working-directory` with them; `generate-writes-both` becomes
  `gen-from-documents-alone`, and `verbs-on-the-found-charter` is added
  (backlog 003).

## CH-041 Lint advises on what the gate does not fail

Status: active
Serves: CH-001

Lint runs the whole gate first and prints every problem it fails on,
then one line counting the open questions of the question file, then
its advice: a decision that serves nothing, a planned entry each of
whose scenarios is verified, a deprecated entry
still anchored or cited, an active decision realized only by a package
clause unless its kind is structure, and a citation of an entry with
scenarios that names none. Then it checks the style a reader would
stumble on, by the bars of the README's Style, in the live entries of
the area documents and the needs: a word retired.txt in the entries
directory retires; a title over twenty words; a text over six sentences;
a scenario sentence over thirty words; a Why under twenty-five words that
names a story, a ruling or a review; and a test whose citations name
scenarios of two or more live decisions and no need. A retired word is
read in a title, a text, a Why, a scenario's sentence and a document's
header, never in a code span or a double-quoted string, and is kept
where the line or the sentence holding it holds one of the word's
phrases. The file holds a word a line: a pattern matched at word
boundaries whatever the case, then -> and the word to write, then | and
the phrases, separated by semicolons. Each style check that finds
anything prints one line, after the rest, counting what it found and
naming the command that lists it; with -style, lint lists every finding
at its line in their place, a retired word with the word to write and a
test once, at the test.

Why: these are the gaps a reviewer would want to know about and that no
mechanical rule should block a merge on. The style the README sets can
be measured, so lint measures it and an audit spends its reader only on
what needs one; a test two decisions share is where one may rule what
the other states.

Scenarios:
1. `serves-nothing`: a decision with no Serves is advised.
2. `planned-verified`: a planned entry with every scenario cited is
   advised to become active.
3. `deprecated-still-anchored`: a deprecated entry with an anchor or a
   citation is advised to strip.
4. `package-clause-only`: a decision realized only by a package clause is
   advised to anchor the declaration.
5. `no-scenario-named`: a citation of a scenario'd entry that names no
   handle is advised.
6. `gate-problems-first` (owner): a problem the gate fails on,
   a scenario no test cites, is in lint's output.
7. `retired-word-named` (owner): a retired word outside code,
   quotes and a phrase's line, in a live entry's title, text, Why or
   scenario or a header, is counted, and listed under -style.
8. `long-title-named` (owner): a live entry's title of
   twenty-one words is counted, and listed under -style with its length,
   and one of twenty is not.
9. `long-text-named` (owner): a text of seven sentences is
   counted, and listed under -style, and one of six is not, no period in
   code or of e.g. or i.e. ending one.
10. `long-scenario-named` (owner): a scenario sentence of
    thirty-one words after its handle and provenance is counted, and
    listed at its line under -style, and one of thirty is not.
11. `pointing-why-named` (owner): a Why under twenty-five
    words naming a story, a ruling or a review is counted, and listed
    under -style; a longer one, or one naming none, is not.
12. `shared-test-named` (owner): a test in any language
    citing scenarios of two decisions and no need is counted, and listed
    once under -style; one also citing a need is not.

Changes:
- 2026-10-09: lint ran the load alone and missed the rules that need the
  whole tree; it runs the whole gate first, and lists the open questions.
- 2026-10-09: lint prints one line counting the open questions, where it
  printed one per question and buried its advice; claimsgen questions
  lists them (CH-261009-4B3B).
- 2026-10-09: lint names the style the record's reader would stumble on:
  retired words from claims/retired.txt, titles over twenty words, texts
  over six sentences, scenario sentences over thirty words, a Why that
  only points, and a test two decisions share; by default the style checks
  print their counts, and `lint -style` lists the lines.
- 2026-11-26: `shared-test-named` holds for a test in any language the
  charter reads, Python's among them (backlog 003).

## CH-261009-SNX2 lint -baseline fails when a style check's count exceeds its count in the style baseline, naming both counts

Status: deprecated
Kind: bound
Serves: CH-001
See: CH-041

lint -baseline reads style-baseline.txt in the entries directory, a line
per style check holding its name and a whole count, blank lines and lines
opening with # skipped. The checks are retired-words,
titles-over-twenty-words, texts-over-six-sentences,
scenario-sentences-over-thirty-words, pointing-whys and shared-tests, and
a check the file leaves out has a baseline of zero. It runs the six style
checks and, for each whose count exceeds its baseline, prints one line at
the check's line in the file, or at its first line for a check it leaves
out, naming both counts and what to do, and exits non-zero; when none grew
it prints nothing. A line naming an unknown check or one counted already,
or whose count is not a whole number, is a lint line of its own, and a
missing file is an error saying there is no baseline. The gate's problems
are not part of it, and -baseline with -style is refused.

Why: the style counts are the record's debt, and paying it is slow: the
lines a reader stumbles on are fixed a few at a time, and no change in
another area can wait for them all. A bar that only tightens keeps each
change from adding to the debt while the fixing goes on, and lowering a
count as its lines are fixed keeps the bar where the record stands.

Scenarios:
1. `within-baseline-silent` (owner): with each style count at
   or under its baseline, lint -baseline prints nothing, though the gate
   has problems and questions are open.
2. `grown-count-named` (owner): a count over its baseline is
   named at the check's line with both counts and what to do, and lint
   -baseline exits non-zero.
3. `missing-check-read-as-zero` (owner): a check the file
   leaves out counts from zero, so one finding is named as grown from 0
   to 1; a missing file is an error, not zeros.
4. `bad-line-named` (owner): a line naming an unknown check
   or one counted already, or whose count is not a whole number, is a
   lint line at that line.
5. `flags-exclusive`: lint with -baseline and -style
   together is refused, and either alone is taken.

Changes:
- 2026-11-26: deprecated (Q-261015-5YV3): style is advice, nothing runs
  the ratchet and the style baseline is not kept, so lint -baseline goes
  with the rewrite (backlog 003).

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

## CH-261005-VFJR mutate breaks an entry's realizing code and says, per scenario, whether the tests that cite it noticed

Status: deprecated
Serves: CH-001
See: CH-042, CH-030

mutate takes active entries by id, or the active entries whose realizing
Go functions a diff since a revision touched, and makes small changes
inside those functions alone: a comparison or a logical operator flipped,
a condition negated or its branch emptied, a call dropped, an error
returned as nil, a boolean flipped, an integer moved by one, a break
swapped with a continue; a diff-scoped run changes only the lines the
diff changed. A mutant a caller wrote, an exact replacement in a file,
runs beside them. Each mutant runs only the tests that cite its entry,
through an overlay so no file of the tree is written: it is killed when
one of them fails or times out, survives when all pass, and is counted
apart when it does not build. The citing tests first run unmutated, and
an entry whose tests fail then is reported unsteady and none of its
mutants runs. The result gives, per entry, the mutants killed, survived
and unbuilt, the survivors with their place and change, and the
realizing anchors it could not mutate; per scenario, the mutants its
tests killed, a scenario whose tests killed none being hollow. A run
after which the tree holds a file it did not hold before fails, naming
the file; a survivor never fails it.

Why: the gate proves each scenario is cited and the report that its
tests pass; neither proves a test would notice the decision broken. A
survivor points at a test that would not, a hollow scenario at a test or
an anchor that misses the code it names. Advice, not a gate: some
mutants change nothing a decision promises, and only a reader can tell.

Scenarios:
1. `mutants-in-anchored-code`: mutants are made only inside the named
   entries' realizing functions, and in a run scoped to a diff only on its
   changed lines, for the entries whose functions it touched.
2. `killed-by-citing-test`: a mutant a citing test fails on is killed and
   counted for each scenario the test names; one they all pass survives,
   listed with its file, line, function and change.
3. `hollow-scenario`: a scenario whose citing tests kill no mutant of its
   entry is reported hollow.
4. `unbuilt-apart`: a mutant that does not build is counted as unbuilt,
   neither killed nor survived.
5. `unsteady-not-mutated`: an entry whose citing tests fail unmutated is
   reported unsteady, naming the tests, and none of its mutants runs.
6. `given-mutants-run`: a mutant given as a replacement runs like the
   others, and one whose text does not occur exactly once in its file is
   refused, naming the file.
7. `tree-untouched`: no file of the tree changes during a run, and a run
   whose tests leave a new file in the tree fails, naming it.
8. `unmutable-anchors-listed`: a realizing anchor on a type, a variable,
   a constant or a package is listed as not mutated.

Guidance: go test -overlay carries each mutant; the run keeps its own
build cache and empties it past a size, since every mutant links a fresh
test binary. A cap per entry with a seed samples a large entry the same
way twice. When it runs is the loop's to say (claims/README.md, "The gate
and the report"): on request, when a story asks or a reader doubts a
test, never as a step of every change.

Changes:
- 2026-10-09: two scenarios reworded to thirty words or fewer; the ruling
  is unchanged.
- 2026-10-11: Guidance says mutate is opt-in, run on request and never as
  a step of every change, the owner having ruled so when a diff-scoped
  run over one change named 858 mutants and would have run for hours;
  the tool is unchanged.
- 2026-11-26: deprecated: mutate is dropped with the rewrite, its command
  and code with it, and may return per language if the benchmarks show it
  earns its cost (backlog 003).

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
- 2026-10-09: the title and three scenarios reworded to the style's bars;
  the ruling is unchanged.
- 2026-11-26: charter gen is the one writer of the digest, so
  CHARTER_UPDATE_DIGEST and its Accounts are gone (Q-261015-Y682).

## CH-261015-9ZTY A charter exports as one file, which the viewer shows as linked pages in its reading order

Status: planned
Serves: CH-001, CH-002
See: CH-040, CH-261005-6ZBS, CH-261009-3FHW

charter export writes the charter it finds to one file: each area document in the
reading order with its header, each entry with its fields, text, Why,
scenarios and who ruled each, Guidance and Changes, its anchors and its
citations by handle, the waivers, the observations and the gate's
problems. A charter the gate fails on is still exported, its problems
printed apart. The viewer reads that file, or one dropped on the page,
and names no charter's areas or ids, so the charter's own charter reads
as an adopter's does. Its pages walk the areas in the reading order,
place each need under its area as the digest does, show each area
document with coverage beside each entry and each entry with what it
serves and what serves it, and search the record, an exact id opening its
entry; each page is a link that opens it again. A scenario shows tested
when a test cites it, waived when a waiver accepts it and untested
otherwise, and passed or failed only from a report beside the file.

Why: the digest reads the record as prose in one order, but a reviewer
following an entry to the needs it reaches, or the owner checking which
scenarios a test cites and who ruled them, needs the record as linked
pages with what the gate knows beside each part. A reader that names no
charter keeps working as the record grows and reads the charter's own as
well, and one that writes nothing cannot drift from what it shows.

Scenarios:
1. `export-carries-the-record` (agent, provisional): exporting a
   charter writes one file with each area in the reading order and each
   entry's text, scenarios, anchors and citations, the waivers and the
   observations.
2. `provenance-exported` (agent, provisional): each scenario in the
   export says who ruled it, as the registry does, and one with nothing
   written is an agent's, provisional.
3. `exported-despite-problems` (agent, provisional): a charter the gate
   fails on is exported, its problems printed apart and carried in the
   file.
4. `any-charter-reads` (agent, provisional): the charter's own
   charter, whose areas and ids the viewer names nowhere, shows its
   walkthrough, areas and entries as an adopter's does.
5. `pages-in-reading-order` (agent, provisional): the walkthrough lists
   the areas in the reading order, each need sits under the area the
   digest gives it, and an area's previous and next follow the order.
6. `entry-shows-its-relations` (agent, provisional): an entry's page
   shows the needs it reaches and the tree of decisions that serve it.
7. `page-opens-from-its-link` (agent, provisional): every page opens
   again from its link, and an exact id searched opens its entry.
8. `coverage-from-citations` (agent, provisional): a scenario a test
   cites shows tested, one a waiver accepts waived, one neither untested,
   and passed or failed only with a report beside the file.

Changes:
- 2026-11-26: the export is charter export, of the charter the command
  finds, where claimsmap chose between the product's charter and the
  charter's own; the text and two scenarios name an adopter's charter for
  the product's (backlog 003).

## CH-261009-QQC8 charter diff names what a change did to the record that its Changes lines do not say

Status: active
Serves: CH-261009-3FHW

Given a git revision, charter diff reads the area documents as they
were there and as they are, and names, by id: an entry removed; an active
entry made planned; an entry deprecated, or its text, Kind or a scenario
changed or removed, with no Changes line the revision lacked; and a
scenario whose ruling was the owner's or a delegate's and is now an
agent's. An entry only set active, or changed with a new dated line, is
not named. It exits 1 when it names anything.

Why: the gate reads one tree; the drift it cannot see is between two, a
test made to assert the opposite of the scenario it cites and nothing
written down. Against the base of a pull request, this is the check that
a change to a decision is a change recorded.

Scenarios:
1. `removed-named` (owner): an id present at the revision and
   absent now is named: an entry is deprecated, never deleted.
2. `unbuilt-named` (owner): an entry active at the revision
   and planned now is named.
3. `unrecorded-change-named` (owner): an entry whose text,
   Kind or scenarios differ from the revision's with no Changes line the
   revision lacked is named, and so is one deprecated without a line.
4. `ruling-lost-named` (owner): a scenario ruled by the owner
   or a delegate at the revision and provisional now is named.
5. `recorded-change-passes` (owner): a change with a new dated
   Changes line, or an entry only set active, names nothing.

Changes:
- 2026-11-26: the verb is charter diff, of the one charter command
  (CH-040); the ruling is unchanged (backlog 003).

## CH-261126-5JMD The loop's skills ship as a plugin of this repository and read the layout from charter.toml alone

Status: planned
Serves: CH-261126-NS7T
See: CH-261126-33GK

The claims skills, draft, implement, review, audit and probe, ship as a
Claude Code plugin of this repository, which an adopter installs. A skill
names the commands of the method and finds the adopter's layout, the
claims directory, the languages and their test commands, in charter.toml;
it names no file of one repository's layout. The method's document gives
the section an adopter puts in its CLAUDE.md as a template.

Why: skills written for one product send an agent in another repository to
files that are not there; skills that read the layout from the file the
command reads go wherever the charter does, and a plugin is how they get
there.

Scenarios:
1. `plugin-lists-the-skills` (owner, backlog 003): the plugin's manifest
   lists the five skills, and each skill's file is in the plugin.
2. `no-fixed-layout` (owner, backlog 003): no skill names cmd/claimsgen,
   claims/README.md or a path of one product, and each one needing the
   layout reads charter.toml.
3. `claude-md-template` (owner, backlog 003): the method's document holds
   the CLAUDE.md section as a template naming the charter command and
   charter.toml.
