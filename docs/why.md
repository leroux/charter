# Why charter exists

Charter is a level to build software at above the code. What the software
must do is written down once, in markdown under `claims/`. The code is held
to that record by a static gate, and a change is reviewed through the
record, by people and by agents. The record states its problems, one per
need in [needs.md](../claims/needs.md), in this order: intent has no
durable home; code drifts from intent unnoticed; reviewing code line by
line does not scale; who decided what is invisible; a newcomer cannot find
the way in; and every project is different. This doc takes each problem
in turn: the need that answers it, the mechanism that meets the need, what
the record says it chose against, and what it costs.

A need is what a person can rely on; a decision says how the system
behaves at a boundary and serves one or more needs. Both are entries,
written in the same form, as [the loop](loop.md) and
[architecture](architecture.md) describe.

## Intent has no durable home

Intent that lives in heads, chats and pull requests is lost to the next
person and to every agent. The code says what the software does, so
without a written reference there is nothing to hold it to.

<!-- charter:embed CH-261126-PFG5 -->
> **[CH-261126-PFG5](../claims/needs.md#ch-261126-pfg5-a-repositorys-charter-must-write-what-the-software-must-do-once-above-the-code-where-people-and-agents-both-read-and-change-it) A repository's charter MUST write what the software must do once, above the code, where people and agents both read and change it**
>
> - A repository's needs and decisions MUST be written in its claims
>   directory, each once, in the form every charter shares, and they are
>   the reference a change is held to.
> - Each entry MUST keep its id for life, and what changed in it, and
>   when, MUST be written beside it.
<!-- /charter:embed -->

The form is fixed so that it can be read by program and diffed by part.
An area document is a header (a title, an optional `Code:` line and free
prose) and then one entry per `## ID Title` heading: field lines (`Status`,
`Kind`, `Serves`, `See`, `Accounts`), the text, then paragraphs opened by
`Why:`, `Scenarios:`, `Guidance:` and `Changes:`, in that order
([CH-010](../claims/intent.md#ch-010-the-record-reader-must-read-an-entry-as-a-heading-with-id-and-title-then-field-lines-text-why-scenarios-guidance-and-changes-in-that-order)). [`parseArea`](../internal/record/areas.go#L141)
reads it and reports every departure as a `form:` problem with its file
and line, which fails the gate.

An id is permanent. A new one is the prefix, the UTC day it was taken and
four characters drawn from `BCDFGHJKMNPQRSTVWXYZ23456789`, such as
`CH-261127-3NY3`, given by `charter next`
([CH-012](../claims/intent.md#ch-012-an-id-must-be-declared-once-and-a-new-id-must-be-its-prefix-the-day-it-was-taken-and-a-random-code), [`Record.NewIDs`](../internal/record/ids.go#L99)).
The record says what it chose against: ids were numbered one above the
highest, so branches cut from one tree took the same number and the later
one rewrote it in every document, anchor and citation at merge. A random
code from about 600,000 needs no counter to share; the tree is the state,
and the gate names both places on the day two branches do collide.

What changed in an entry is a dated `Changes:` line on the entry itself,
`- YYYY-MM-DD: what changed and why`. The record of a decision and its
history therefore travel together, and [charter diff](diff.md) can check
that every change carries one.

The cost is the form itself: an entry that is not in it fails the gate,
and a person or an agent writes the text, the Why and a scenario per
clause before the gate will call it covered.

## Code drifts from intent unnoticed

Code written fast, and by agents most of all, drifts from what was
decided without anyone noticing. The need asks for the only kind of
adherence that does not depend on someone remembering to look: a check
that fails on the gap.

<!-- charter:embed CH-261126-P2SK -->
> **[CH-261126-P2SK](../claims/needs.md#ch-261126-p2sk-the-charter-gate-must-fail-a-change-that-leaves-a-gap-between-the-record-and-the-code) The charter gate MUST fail a change that leaves a gap between the record and the code**
>
> - Every active decision MUST have code that claims it and a test that
>   shows each of its scenarios.
> - Every surface the tree exposes MUST be owned by an entry.
> - The gate MUST fail a tree where any of that is missing, so a change,
>   a person's or an agent's, that leaves such a gap does not pass.
<!-- /charter:embed -->

The link from a decision to code is an anchor, a `charter:realizes`
directive on the declaration that carries the decision; the link from a
scenario to a test is a citation, a `verifies` call that opens the test
and names the scenario's handle. The gate fails an active decision with
no realizes anchor and an active scenario no citation names
([CH-030](../claims/gate.md#ch-030-the-gate-must-fail-an-active-decision-without-an-anchor-or-with-an-uncited-scenario-a-waiver-may-absorb-a-missing-test-never-a-missing-anchor)). A waiver may accept a missing test, with a
reason on its line in `claims/waivers.txt`; nothing accepts a missing
anchor, because code nobody can find from its decision has no reason that
makes it acceptable.

The gate also keeps books of the tree's surfaces: every environment
variable, and the keys of any kind the adopter declares, must be named in
some entry's `Accounts:` or listed with a reason in `claims/baseline.txt`
([CH-033](../claims/gate.md#ch-033-the-gate-must-enumerate-every-surface-and-require-each-to-be-accounted-for-by-an-entry-or-listed-in-a-baseline-of-live-unaccounted-surfaces)). Code reads the environment only through a
named constant of a constants module, so no variable escapes the books by
being read through a literal ([CH-034](../claims/gate.md#ch-034-code-must-read-an-environment-variable-through-a-named-constant-or-under-a-dynamic-directive-and-a-literal-must-not-name-one)).

All of it is static: [`Check`](../internal/gate/gate.go#L56) reads facts the
readers parsed out of the source and never runs a test. That is what lets
it run in seconds and give the same verdict on every machine, and it is
also its limit: the gate proves the links exist, not that the code is right
or that a test means what its scenario says. [How check works](check.md)
follows it rule by rule.

## Reviewing code line by line does not scale

Review is the scarce resource, and an agent can write more code in an hour
than a person can read in a day. The need moves the reviewer's attention
to what no gate can check.

<!-- charter:embed CH-001 -->
> **[CH-001](../claims/needs.md#ch-001-the-charter-gate-must-let-a-reviewer-judge-a-change-by-its-entries-and-their-tests-without-reading-the-code) The charter gate MUST let a reviewer judge a change by its entries and their tests, without reading the code**
>
> - Every clause of an active decision MUST be a scenario, and every
>   scenario MUST have a test that names it.
> - The gate MUST fail on the first one that does not, so the reviewer
>   reads the area document diff and the gate's verdict, and opens the
>   code only where the review points.
> - The gate MUST be static, run in seconds, and give the same verdict on
>   every machine.
<!-- /charter:embed -->

The reviewer reads the area document diff, which says in sentences what
changed in the rules, and the gate's verdict, which says every rule is
still anchored and cited. Three things then cover what the gate cannot see:

- `charter report` runs every language's tests and says, per scenario,
  whether its citing tests passed, failed, were skipped or did not run
  ([CH-042](../claims/review.md#ch-042-charter-report-must-join-each-tests-markers-to-its-scenarios-and-say-which-passed-failed-were-skipped-or-have-no-test); [the report](report.md)).
- The claims-review skill, run in a context that did not write the
  change, asks three questions of each scenario the change touched: does
  the citing test show what the sentence says, does the anchored code
  carry the decision, and does the Changes line say what changed and why.
- The digest, when charter.toml names one, is the whole record as one
  document in reading order ([CH-261005-6ZBS](../claims/review.md#ch-261005-6zbs-charter-gen-must-write-a-digest-that-reads-the-charter-in-order-needs-by-area-an-area-map-then-each-areas-decisions-with-their-why)).

The record says what the report chose against: it read `go test`'s JSON
stream, which only Go writes. Each language's generated `verifies` now
appends one JSON line per citation to the file `CHARTER_EVENTS` names, so
the report makes one join for every test runner
([CH-024](../claims/anchoring.md#ch-024-each-languages-verifies-must-record-each-citation-and-its-tests-outcome-as-a-json-line-in-the-events-file)).

## Who decided what is invisible

An agent's default and the owner's ruling look the same once written.
Without a mark of who ruled, the record cannot say what was actually
decided, and a change nobody wrote down is the drift the form exists to
refuse.

<!-- charter:embed CH-261009-3FHW -->
> **[CH-261009-3FHW](../claims/needs.md#ch-261009-3fhw-the-record-must-show-who-ruled-each-scenario-and-charter-diff-must-name-any-change-to-an-entry-the-record-does-not-state) The record MUST show who ruled each scenario, and charter diff MUST name any change to an entry the record does not state**
>
> - Every scenario MUST say who ruled it: the owner, with what the ruling
>   cites; a named delegate; or an agent whose choice nobody has ruled on,
>   which is what a scenario with no provenance written is.
> - The registry and the report MUST carry the provenance, and the report
>   MUST count the active scenarios by who ruled them.
> - A change to the record that its Changes lines do not state (an entry
>   removed, unbuilt or retired without a line, its text or a scenario
>   changed without one, or a ruling lost) MUST be named against the
>   revision it was made from.
<!-- /charter:embed -->

Each scenario carries its provenance after its handle: `(owner, backlog
003)`, `(delegate name, #12)`, or `(agent, provisional)`, which is also
what a scenario with nothing written is read as
([CH-261009-VY5C](../claims/provenance.md#ch-261009-vy5c-the-record-reader-must-read-who-ruled-a-scenario-the-owner-with-what-it-cites-a-named-delegate-or-when-nothing-is-written-an-agent-provisionally),
[`parseProvenance`](../internal/record/areas.go#L103)). A choice no entry rules
becomes a question in `claims/questions.md` with the answer the agent took
meanwhile ([CH-261009-4B3B](../claims/provenance.md#ch-261009-4b3b-the-question-file-must-hold-what-nobody-has-ruled-each-question-naming-its-entry-who-asked-and-the-provisional-answer-until-it-is-answered)). `charter diff REV`
names every change to the record since a revision that no new Changes
line states, and every ruling lost
([CH-261009-QQC8](../claims/provenance.md#ch-261009-qqc8-charter-diff-must-name-every-change-to-the-record-that-its-changes-lines-do-not-state); [charter diff](diff.md)).

## A newcomer cannot find the way in

An agent building to a decision needs a way into the code and a way to
know when it is done, without reading the whole tree or guessing what a
failure means. So does a person new to the code.

<!-- charter:embed CH-002 -->
> **[CH-002](../claims/needs.md#ch-002-the-charter-command-must-lead-a-builder-from-an-id-to-its-code-and-tests-and-every-gate-problem-must-name-its-fix) The charter command MUST lead a builder from an id to its code and tests, and every gate problem MUST name its fix**
>
> - From an entry's id, the tooling MUST list the code that realizes it, by
>   file, line and declaration, and the tests that verify it, by file, test
>   and scenario.
> - Every problem the gate reports MUST name a file, a line and what to do.
> - Symbols MUST be generated in each language the charter reads, so a
>   citation of an entry that does not exist fails to compile or to import.
<!-- /charter:embed -->

`charter anchors ID` prints, as JSON, an entry's scenarios, the
declarations that realize it, the tests that cite it with their handles,
its relations, the docs that mention it and its Changes
([CH-040](../claims/anchoring.md#ch-040-one-charter-command-must-generate-check-give-ids-and-list-anchors-so-an-adopter-writes-no-main)). The generated symbol for each entry,
`claims.CH_001` in every language, means a citation of an entry that does
not exist fails to compile or to import, where a string id would rot
silently ([CH-020](../claims/anchoring.md#ch-020-charter-gen-must-generate-a-symbol-per-entry-in-each-language-read-named-by-the-id-with-underscores-for-hyphens-and-documented-by-the-title)). Every problem the gate prints
says what to write, such as the `verifies(claims.SW_118, "handle")` call a
missing scenario needs.

Docs, such as this one, are the story over the record and the code. Their
prose is free, but their references to entries and to declarations are
checked, and a doc a person has read carries a hash of everything it
references so that check warns when any of it changes
([CH-261127-3NY3](../claims/anchoring.md#ch-261127-3ny3-a-doc-must-be-a-markdown-file-chartertoml-names-must-declare-no-entry-and-its-references-to-entries-must-resolve),
[CH-261127-GXP7](../claims/anchoring.md#ch-261127-gxp7-charter-curate-must-record-who-curated-a-doc-and-a-hash-of-each-reference-and-charter-check-must-warn-when-one-changes); [docs](docs.md)). The record
says what this chose against: backlinks written by hand in both directions
fall out of step, so the docs that mention an entry are computed from the
docs ([CH-261127-8MYR](../claims/anchoring.md#ch-261127-8myr-charter-must-compute-which-docs-mention-an-entry-from-the-docs-and-charter-anchors-and-charter-diff-must-show-it)); and a curation date
compared against git history needs that history intact, which a shallow
clone or a squashed merge loses, so the hash is stored in the doc.

## Every project is different

A method that a repository can use only by copying its tooling is a fork.

<!-- charter:embed CH-261126-NS7T -->
> **[CH-261126-NS7T](../claims/needs.md#ch-261126-ns7t-an-adopter-must-be-able-to-keep-a-charter-in-any-repository-layout-and-language-with-code-that-imports-nothing-of-the-tooling) An adopter MUST be able to keep a charter in any repository layout and language, with code that imports nothing of the tooling**
>
> - A repository MUST be able to take the charter up by writing its
>   entries and a charter.toml in its own claims directory and running the
>   charter command over it.
> - Its code MUST import nothing of the tooling: its tests cite through the
>   files the command generates.
> - The layout is the repository's own: Go modules nested anywhere, Python
>   and TypeScript beside them, and no go.mod at the root.
> - The examples in this repository MUST be kept so, each checked with the
>   tooling's tests.
<!-- /charter:embed -->

The command finds `claims/charter.toml` by walking up from the working
directory, as git finds a repository, and refuses a key it does not know
([CH-261126-33GK](../claims/adoption.md#ch-261126-33gk-charter-must-find-claimschartertoml-by-walking-up-from-the-working-directory-and-must-refuse-a-key-it-does-not-know),
[`Find`](../internal/config/config.go#L128)). Its record says why the
configuration is a file: a charter configured in Go code could be read only
by a Go build. Go, Python and TypeScript are each a reader, a generator and
a marker built into the command, and the gate's rules read facts of one
shape that name no language ([CH-261126-PQ7X](../claims/adoption.md#ch-261126-pq7x-charter-must-support-each-language-as-a-built-in-reader-generator-and-marker-and-must-read-a-language-only-when-chartertoml-names-it)).
The record says what that replaced: Go woven through the gate, which made
a second language a parallel path beside it.

The adopter's code imports nothing of the tooling. Its tests import the
files `charter gen` writes into its own `claims/`, which import only their
language's standard library. Both examples, `examples/ratelimit` and
`examples/slugs`, are charters of their own whose gates run among the
tooling's tests, so a change that would break an adopter fails here first.
The tooling also specifies itself: its decisions are the `CH` entries in
`claims/`, held by its own gate ([CH-006](../claims/adoption.md#ch-006-the-charter-must-specify-itself-and-its-own-gate-must-pass)).
[Adopting charter](adopting.md) walks through both examples.

## What it costs and where it stops

- Every decision is written before it is built, with a scenario per clause
  and a test per scenario; that is the price of a record a reviewer can
  read instead of the code.
- The gate proves that links exist. Whether a test means what its
  scenario says, and whether anchored code carries the decision, are the
  judge's questions (claims-review) and the audit's (claims-audit), and a
  person's answer.
- A doc's prose is held to nothing; only its references are checked, and
  a stale doc warns, never fails.
- The languages are a fixed set the charter can test. Another language is
  a change to the charter, not a plug-in.
