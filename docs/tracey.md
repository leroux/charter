# Charter and tracey compared

[tracey](https://github.com/bearcove/tracey) is the closest tool to charter
that we know of. Both keep requirements in markdown, link them to code and
tests through comments, and specify themselves by their own method. They
differ in what they hold the code to. Tracey measures traceability and
shows it everywhere a developer looks: a dashboard, an editor, an MCP
server. Charter enforces it: a static gate fails the tree when an active
decision has no anchor or a scenario no test cites, and the record also
says why each decision was made and who ruled it. This doc compares the
two part by part, says what each does that the other does not, and ends
with what charter could take from tracey.

What it says of tracey was read at commit `66008a8` on its main branch
(8 October 2026, version 2.0.0-rc.0): the specification at
`docs/content/spec/tracey.md`, the guides beside it, its own configuration
and its source. What it says of charter is in the record, and where this
doc and the record differ, the record binds.

## At a glance

| | Charter | Tracey |
|---|---|---|
| Unit of the record | An entry: id, title, field lines, text, Why, scenarios, Changes | A requirement: a marker `r[id]` and the paragraph after it |
| Spec formats | Markdown area documents in one fixed form | Markdown, StrictDoc `.sdoc`, Typst, AsciiDoc |
| Ids | Prefix, UTC day and a random code, permanent, from `charter next` | Dotted names chosen by the author, renamable through the LSP |
| Link from code | `charter:realizes` on the declaration that carries the decision | `r[impl id]` in any comment; also `depends` and `related` |
| Link from a test | A `verifies` call opening the test, naming each scenario by handle | `r[verify id]` in any comment of a test file |
| Granularity of a test link | A scenario | A requirement |
| Languages | Go, Python, TypeScript, plus comment-line anchors in named files | About twenty, through tree-sitter |
| What fails | An active decision unanchored, a scenario uncited, a surface unowned, a form problem, a broken reference | A broken or stale reference, a duplicate, an unknown prefix, `impl` in a test file |
| Coverage | Required of every active entry | A percentage per spec and implementation |
| Change over time | Dated `Changes:` lines; `charter diff` names a change without one | A version suffix `+N`; older references become stale |
| Who decided | Provenance per scenario; open questions in `questions.md` | Not recorded |
| Runs the tests | `charter report`, which joins each outcome to its scenario | No |
| Interfaces | One command; JSON from `anchors` and `report` | A daemon with a dashboard, an LSP, an MCP server and a query CLI |
| Agent support | The loop's skills: draft, implement, review, audit, probe, narrate | An MCP server and a skill for adding annotations |
| Implementation | Go | Rust |

## The requirement

A tracey requirement is a marker and the text after it. The marker
`r[auth.login]` stands at the start of a paragraph or of a blockquote, and
the requirement's text runs until the next blank line, requirement or
heading; inside a blockquote it can hold code blocks and examples. The
prefix (`r`, `h2`, any lowercase marker) is inferred from the spec files
and names which spec a reference belongs to, so one repository can
implement its own spec beside an external one, such as an HTTP/2 spec
vendored as a submodule, each with its own prefix. Tracey also reads
StrictDoc, Typst and AsciiDoc, so a team that already keeps requirements
in one of those can adopt it without rewriting them.

A charter entry is a heading `## ID Title` followed by field lines, the
text, and paragraphs opened by `Why:`, `Scenarios:`, `Guidance:` and
`Changes:`, in that order
([CH-010](../claims/intent.md#ch-010-the-record-reader-must-read-an-entry-as-a-heading-with-id-and-title-then-field-lines-text-why-scenarios-guidance-and-changes-in-that-order)). The form is
fixed so that a program can read it and `charter diff` can compare it by
part, and the gate fails any departure from it. Three parts of an entry
have no counterpart in tracey:

- **Why.** Each entry says what problem the rule meets and what it costs.
  A tracey requirement can say this in its prose, but nothing asks for it.
- **Scenarios.** Each clause of a decision is a numbered given / when /
  then line with a handle, and the gate asks a test of each one
  ([CH-030](../claims/gate.md#ch-030-the-gate-must-fail-an-active-decision-without-an-anchor-or-with-an-uncited-scenario-a-waiver-may-absorb-a-missing-test-never-a-missing-anchor)). In tracey
  the requirement is the smallest thing a test can point at, so a
  requirement with three behaviours is covered by one test of any of them.
- **Relations and kinds.** `Serves:` ties a decision to the needs it
  answers and `See:` to its neighbours, and the gate fails one that names
  nothing ([CH-015](../claims/intent.md#ch-015-the-gate-must-fail-a-serves-or-see-that-names-no-existing-entry-or-that-names-the-entry-itself)).
  A tracey requirement has no field for its relations to other
  requirements.

The ids differ in purpose. Tracey's ids are names (`auth.token.validation`)
that tell a reader what the requirement is about and that a rename in the
editor rewrites everywhere. A charter id such as `CH-261127-3NY3` means
nothing on its own and never changes; the title carries the meaning
([CH-012](../claims/intent.md#ch-012-an-id-must-be-declared-once-and-a-new-id-must-be-its-prefix-the-day-it-was-taken-and-a-random-code)). The record
says why: numbered ids made two branches take the same number, and a
random code needs no counter to share. A name chosen by its author has the
same property, since two authors rarely choose the same name, and tracey
reports a duplicate within a spec as an error. What the random code adds
is that the id never has to change when the rule's subject does.

## Linking code to a requirement

Tracey reads comments with tree-sitter in about twenty languages, among
them Rust, Go, Python, TypeScript, Swift, Java, C and C++, Zig, YAML and
JSON5. A reference may stand in any line, block or doc comment, with text
around it, and carries one of four verbs: `impl`, `verify`, `depends` (the
code must be rechecked when the requirement changes) and `related` (shown
in review). Tracey then attaches each reference to the smallest code unit
around it (a function, a struct, a module), which gives it reverse
traceability: `tracey query unmapped` shows the source tree with the share
of its code units that some requirement claims, and lists the units no
requirement claims.

Charter reads an anchor only where a decision lives: the last line of a
Go doc comment, a TypeScript line comment, or the comment lines above a
Python def ([CH-021](../claims/anchoring.md#ch-021-production-code-must-anchor-an-entry-in-a-go-doc-comment-a-typescript-line-comment-or-the-comment-lines-above-a-python-def)), and
a comment line in a file of a kind charter.toml names, such as a skill's
frontmatter ([CH-261127-45KM](../claims/anchoring.md#ch-261127-45km-charter-must-read-an-anchor-from-a-comment-line-in-a-file-no-language-reads-where-chartertoml-names-that-file)).
It parses each language with that language's own parser, Go's in process
and Python and TypeScript through their own toolchains, rather than with
a shared grammar. There are two verbs, `realizes`, which counts toward
coverage, and `concerns`, which is navigation only
([CH-022](../claims/anchoring.md#ch-022-an-anchor-must-name-an-entry-that-exists-and-is-not-planned-and-a-concerns-anchor-must-not-count-toward-coverage)); `concerns`
covers roughly what tracey's `related` does, and nothing in charter
corresponds to `depends`.

The trade-off is breadth against precision. Tracey works in nearly any
repository the day it is installed. Charter supports three languages, and
the record says another one is a change to the charter, but it knows each
one's declarations and tests well enough to refuse a citation that its
test runner would never run.

Charter has no measure of unclaimed code. Its counterpart looks outward
instead: the gate enumerates the tree's surfaces, every environment
variable and the keys of any kind the adopter declares, and fails one that
no entry's `Accounts:` names and the baseline does not list
([CH-033](../claims/gate.md#ch-033-the-gate-must-enumerate-every-surface-and-require-each-to-be-accounted-for-by-an-entry-or-listed-in-a-baseline-of-live-unaccounted-surfaces),
[CH-034](../claims/gate.md#ch-034-code-must-read-an-environment-variable-through-a-named-constant-or-under-a-dynamic-directive-and-a-literal-must-not-name-one)). Tracey asks
how much of the code a spec explains; charter asks whether everything an
operator can set or a client can see has an owner.

## Linking a test to a requirement

In tracey a test links to a requirement the same way code does, with an
`r[verify id]` comment. Files matched by `test_include` may hold only
`verify` references, and an `impl` among them is an error. Nothing checks
that the comment stands in a test, that the test runs, or what it asserts,
and tracey never runs the tests.

In charter a citation is code. The test's first statement is a call such
as `claims.Verifies(t, claims.CH_030, "missing-anchor-fails")`, naming the
entry by a symbol `charter gen` generated and each scenario by its handle
([CH-020](../claims/anchoring.md#ch-020-charter-gen-must-generate-a-symbol-per-entry-in-each-language-read-named-by-the-id-with-underscores-for-hyphens-and-documented-by-the-title),
[CH-023](../claims/anchoring.md#ch-023-a-test-must-cite-an-entry-by-a-verifies-call-that-opens-a-test-function-naming-by-handle-the-scenarios-it-exercises)). The gate reads
the call statically and fails one that anything precedes, one in a
function the runner would not run, a skipped TypeScript test, and a handle
the entry does not list. Because the entry is a symbol, citing an entry
that does not exist fails to compile or to import. At run time the same
call writes a line to the file `CHARTER_EVENTS` names, and `charter
report` runs every language's tests and says per scenario whether its
tests passed, failed, were skipped or did not run
([CH-024](../claims/anchoring.md#ch-024-each-languages-verifies-must-record-each-citation-and-its-tests-outcome-as-a-json-line-in-the-events-file),
[CH-042](../claims/review.md#ch-042-charter-report-must-join-each-tests-markers-to-its-scenarios-and-say-which-passed-failed-were-skipped-or-have-no-test)).

The cost of charter's form is a generated module the tests import and a
call in every test that cites. What it buys is that "this scenario is
tested" means a test that runs names it, and the report can say whether
that test passed.

## What fails the build

Tracey's `query validate` exits non-zero on errors: a reference to a
requirement that does not exist, an unknown prefix, a duplicate id, a
stale reference, an `impl` in a test file, a file the include patterns
match but tracey cannot parse. Warnings fail only with `--deny warnings`.
A requirement with no implementation or no test is not among the errors;
`tracey query uncovered` and `untested` list them, and coverage is a
percentage shown in the dashboard, the editor and the head of every MCP
response. Tracey is built for a project to raise that percentage, and for
several implementations of one spec, in different languages, to be
compared by it.

Charter's gate fails the tree on any gap. An active decision with no
`realizes` anchor fails and no waiver can accept it; an active scenario
that no citation names fails unless a line in `claims/waivers.txt`
accepts it with a reason ([CH-030](../claims/gate.md#ch-030-the-gate-must-fail-an-active-decision-without-an-anchor-or-with-an-uncited-scenario-a-waiver-may-absorb-a-missing-test-never-a-missing-anchor)).
Status is how a charter is built up: a planned entry asks for nothing, a
test may cite it before the code exists, and the gate asks coverage only
of active entries ([CH-013](../claims/intent.md#ch-013-an-entrys-status-must-be-planned-active-or-deprecated-every-status-must-generate-a-symbol-and-the-gate-must-ask-coverage-of-active-entries-only)).
So charter's equivalent of tracey's coverage percentage is the line
between planned and active, and an entry crosses it when its work begins.

The two stances suit different situations. A team bringing traceability
to a large existing codebase can start with tracey at a low percentage and
raise it. A team that wants every rule it has set to hold, and a reviewer
who wants to stop reading the code line by line, needs the gate to fail,
which is what charter's needs ask
([CH-261126-P2SK](../claims/needs.md#ch-261126-p2sk-the-charter-gate-must-fail-a-change-that-leaves-a-gap-between-the-record-and-the-code),
[CH-001](../claims/needs.md#ch-001-the-charter-gate-must-let-a-reviewer-judge-a-change-by-its-entries-and-their-tests-without-reading-the-code)).

## When a requirement changes

This is where tracey does something charter does not. A tracey
requirement may carry a version, `r[auth.login+2]`, and a reference
carries the version it was written against. When the spec's version
moves past a reference's, the reference is stale: the requirement no
longer counts as covered, `query validate` fails, the editor shows a
warning whose message begins "Implementation must be changed to match
updated rule text — and ONLY ONCE THAT'S DONE must the code annotation be
bumped", and the hover shows a word-level diff of the rule's text taken
from git history. `tracey pre-commit` fails a commit whose staged spec
changes a requirement's text without bumping its version, and `tracey
bump` bumps it. Each code site that relied on the old text is therefore
named, and someone must touch each one to say it was reviewed against the
new text.

Charter records the change on the entry instead. Each change is a dated
`Changes:` line saying what changed in the ruling and why, and `charter
diff REV` names every change to the record since a revision that no new
line states, including a scenario changed or a ruling lost
([CH-261009-QQC8](../claims/provenance.md#ch-261009-qqc8-charter-diff-must-name-every-change-to-the-record-that-its-changes-lines-do-not-state)).
For docs, charter has staleness of its own: `charter curate` stores a hash
of each entry and declaration a doc references, and `charter check` warns
when one changes ([CH-261127-GXP7](../claims/anchoring.md#ch-261127-gxp7-charter-curate-must-record-who-curated-a-doc-and-a-hash-of-each-reference-and-charter-check-must-warn-when-one-changes)).
For anchors and citations it has none. When an entry's text changes, the
gate still counts its anchors and citations as they were, and only the
claims-review skill, reading the change, asks whether the anchored code
still carries the decision and the citing test still shows the scenario.

Tracey's version is written by hand at every site and bumped by a tool;
charter's history is prose a reviewer reads. Tracey forces the
acknowledgement and says nothing of why the rule changed; charter says why
and leaves the acknowledgement to review.

## Who decided

Charter records who ruled each scenario: the owner with what the ruling
cites, a named delegate, or an agent whose choice nobody has ruled on yet,
which is what a scenario with no provenance written is read as
([CH-261009-VY5C](../claims/provenance.md#ch-261009-vy5c-the-record-reader-must-read-who-ruled-a-scenario-the-owner-with-what-it-cites-a-named-delegate-or-when-nothing-is-written-an-agent-provisionally)). A
choice no entry rules becomes a question in `claims/questions.md` with the
answer the agent took meanwhile
([CH-261009-4B3B](../claims/provenance.md#ch-261009-4b3b-the-question-file-must-hold-what-nobody-has-ruled-each-question-naming-its-entry-who-asked-and-the-provisional-answer-until-it-is-answered)),
and the report counts active scenarios by who ruled them. The need behind
it is that an agent's default and the owner's ruling look the same once
merged ([CH-261009-3FHW](../claims/needs.md#ch-261009-3fhw-the-record-must-show-who-ruled-each-scenario-and-charter-diff-must-name-any-change-to-an-entry-the-record-does-not-state)).
Tracey has no counterpart. Its requirements say what must hold, and git
blame is the only record of who wrote each one.

## Interfaces

Tracey is built around a daemon, one per workspace, that watches the
files, rebuilds on a change and holds the state; the web dashboard, the
LSP server, the MCP server and the query CLI are bridges to it over a Unix
socket. The LSP server gives hover with the requirement's text,
go-to-definition from a reference to the spec, find-references,
completions of ids and verbs, rename across the tree, diagnostics, code
lenses and inlay hints, and a Zed extension packages it. The dashboard
renders the spec with each requirement's status, filters requirements by
coverage and by MUST, SHOULD or MAY, shows the source tree with coverage
badges, and can edit a requirement in place, refusing a file git does not
track and a save whose file hash no longer matches. The daemon keeps an
overlay of files open in the editor, so coverage follows unsaved changes.

Charter is one command with eight verbs: `check`, `gen`, `next`,
`anchors`, `questions`, `report`, `diff` and `curate`
([CH-040](../claims/anchoring.md#ch-040-one-charter-command-must-generate-check-give-ids-and-list-anchors-so-an-adopter-writes-no-main)). It keeps
no state between runs; `anchors` and `report` print JSON, and every
problem the gate prints names a file, a line and what to write
([CH-037](../claims/gate.md#ch-037-charter-check-must-report-every-problem-sorted-with-its-file-line-rule-and-fix-and-must-pass-a-clean-tree)). There is no
editor integration and no dashboard; the digest, when charter.toml names
one, is the whole record as one document in reading order
([CH-261005-6ZBS](../claims/review.md#ch-261005-6zbs-charter-gen-must-write-a-digest-that-reads-the-charter-in-order-needs-by-area-an-area-map-then-each-areas-decisions-with-their-why)).

## Agents

Both tools are written with agents in mind, and they help them at
different steps. Tracey's MCP server answers questions about the current
state (what is uncovered, untested, stale or unmapped, and what a
requirement says), heads every answer with the coverage of each
spec and implementation and a list of what changed since the last query,
and can add include and exclude patterns to the configuration. Its
bundled skill teaches an agent the annotation syntax and a workflow:
check the status, find an uncovered requirement, read it, annotate the
code.

Charter's skills run the whole loop and ship as a plugin
([CH-261126-5JMD](../claims/adoption.md#ch-261126-5jmd-the-loops-skills-must-ship-as-a-plugin-of-this-repository-and-must-read-the-adopters-layout-from-chartertoml-alone)):
claims-draft turns a ruling into entries, claims-implement builds an entry
to covered, claims-review judges a change in a context that did not write
it, claims-audit reads the record against its tests and code for drift,
claims-probe runs probes against live external systems, and claims-narrate
writes docs like this one. The judge is the part with no counterpart in
tracey. Charter's gate, like tracey's validation, proves that links exist;
claims-review asks of each scenario a change touched whether the citing
test shows what the sentence says and whether the anchored code carries
the decision, which is where an agent's annotation can be present and
wrong.

## Docs and the outside world

Charter checks the prose around the record too. A doc is a markdown file
charter.toml names; its links to entries and to declarations must resolve,
`charter gen` keeps their titles and line numbers current, and a curated
doc warns when anything it references changes
([CH-261127-3NY3](../claims/anchoring.md#ch-261127-3ny3-a-doc-must-be-a-markdown-file-chartertoml-names-must-declare-no-entry-and-its-references-to-entries-must-resolve),
[CH-261127-DW8R](../claims/anchoring.md#ch-261127-dw8r-a-doc-must-reference-code-by-a-link-to-the-file-that-declares-the-named-symbol-and-the-gate-must-fail-one-that-file-lacks)). It
also records observations, facts about external systems that entries rely
on, with probes that check them and a gate rule that fails an active entry
resting on a broken one
([CH-035](../claims/gate.md#ch-035-an-observation-must-carry-its-statuss-fields-and-the-gate-must-fail-a-broken-observation-under-an-active-entry)). Tracey renders the
spec itself as documentation and has neither checked docs nor
observations.

## Self-specification

Both tools are specified by their own method. Tracey's specification is
one document of about 325 requirements, covering the annotation language
and the tool down to its dashboard's keyboard bindings, and its own
configuration scans its crates, its dashboard's TypeScript and a sibling
repository. At the commit read, its sources held about 500 `r[impl ...]`
references and 67 `r[verify ...]` references, so most of its requirements
are claimed by code and fewer by tests. Charter's own record is the `CH`
entries in `claims/`, and its own gate must pass
([CH-006](../claims/adoption.md#ch-006-the-charter-must-specify-itself-and-its-own-gate-must-pass)), so every
active `CH` decision is anchored and each of its scenarios is cited by a
test that runs. The examples, `ratelimit` and `slugs`, are charters of
their own whose gates run among the tooling's tests.

## Choosing between them

Tracey fits a project that has a specification, perhaps one written by
someone else, perhaps one with several implementations, and wants to see
which parts each implementation covers, in whatever language it is
written, from inside the editor. It asks little of the spec's form and
nothing of the tests' code, and it does not fail a build for a gap.

Charter fits a project, often one built largely by agents, that wants its
rules written with their reasons, every active rule held by code and a
test of each clause, a record of who ruled what, and a review that reads
the record instead of the code. It asks for a fixed form, a test per
scenario, and code in Go, Python or TypeScript.

## What charter could take from tracey

None of these is ruled; each would be a story of its own, drafted from
what an adopter shows it needs, like the ideas in
[backlog 006](../backlog/006-later.md).

- **Staleness at the anchor.** Charter warns when a curated doc's
  references change but not when an anchored or citing site's entry does.
  A hash of the entry's text, or of a scenario's sentence, recorded beside
  each anchor or citation, would let `charter check` name every site to
  read again, as tracey's versions do, without asking anyone to bump a
  number by hand; the curation hash already shows how such a hash can be
  stored.
- **Editor feedback.** A language server over `charter anchors` and the
  gate's problems would put an entry's text on hover and each problem
  where it is, which is where tracey spends most of its effort.
- **External specs.** Tracey's prefixes let one repository implement its
  own spec and a vendored one side by side. Charter's composed charters
  in backlog 006 are the nearest idea.
- **Unclaimed code.** Tracey's unmapped view asks how much code no
  requirement explains. Charter's surfaces cover the edges of the tree but
  not its interior.

One of tracey's choices the record turns down outright: a language as a
comment grammar. In charter each language is a reader, a generator and a
marker built into the command, and another language is a change to the
charter, not a plug-in
([CH-261126-PQ7X](../claims/adoption.md#ch-261126-pq7x-charter-must-support-each-language-as-a-built-in-reader-generator-and-marker-and-must-read-a-language-only-when-chartertoml-names-it)). The
record's reason is that facts of one shape let every gate rule hold the
same way in each language, and a fixed set of languages is one the charter
can test; the same three parts are what make a citation compile and the
report know each test's outcome.
