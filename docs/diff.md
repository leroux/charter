# charter diff against a revision

The gate reads one tree, so it cannot see what changed between two. A
test can be edited to assert the opposite of its scenario, and the
scenario's sentence rewritten to match, and both trees pass. `charter diff
REV` is the check between trees: it reads the area documents as they were
at a git revision and as they are now, and names by id every change to the
record that no new dated `Changes:` line states
([CH-261009-QQC8](../claims/provenance.md#ch-261009-qqc8-charter-diff-must-name-every-change-to-the-record-that-its-changes-lines-do-not-state)). Run against the base of a pull
request, it makes every change to a decision a change recorded in the
record. After its findings it lists the docs whose references changed
since the revision, to read again ([CH-261127-8MYR](../claims/anchoring.md#ch-261127-8myr-charter-must-compute-which-docs-mention-an-entry-from-the-docs-and-charter-anchors-and-charter-diff-must-show-it)).

<!-- charter:embed CH-261009-QQC8 -->
> **[CH-261009-QQC8](../claims/provenance.md#ch-261009-qqc8-charter-diff-must-name-every-change-to-the-record-that-its-changes-lines-do-not-state) charter diff MUST name every change to the record that its Changes lines do not state**
>
> Given a git revision, `charter diff` reads the area documents as they
> were there and as they are now, and MUST name, by id:
>
> - an entry removed that was not deprecated there, since an entry is
>   deprecated first and deleted in a later change;
> - an active entry made planned;
> - an entry deprecated, or its text, Kind or a scenario changed or
>   removed, with no Changes line the revision lacked;
> - a scenario whose ruling was the owner's or a delegate's and is now an
>   agent's.
>
> An entry only set active, or changed with a new dated line, MUST NOT be
> named. diff MUST exit 1 when it names anything.
<!-- /charter:embed -->

## Reading the record at the revision

[`Diff`](../internal/diff/diff.go#L27) reads the current record with
[`Read`](../internal/record/record.go#L154), then the record at the revision
with `previous`:

1. `git ls-tree --name-only REV -- claims/` lists the files of the claims
   directory at the revision.
2. Every `.md` file except `README.md`, `questions.md`,
   `observations.md`, `policy.md`, `glossary.md`, `audit.md` and the
   digest is an area document, `needs.md` among them.
3. `git show REV:claims/FILE` gives each one's text, and
   [`ParseArea`](../internal/record/areas.go#L519) parses it as `charter gen`
   parses the current one. Form problems at the revision are ignored; the
   first declaration of an id wins.

Nothing is checked out, so diff works in a dirty tree and needs only the
objects of the revision. A revision git cannot resolve is an error,
`charter diff: reading the record at REV: git ls-tree …`, and exit 2.

## What it names

[`compare`](../internal/diff/diff.go#L164) walks the entries of the revision in
the order it met them, and for each one writes a finding
`claims/FILE:LINE: history: …`, at the entry's heading (the revision's,
for an entry that is gone):

| At the revision | Now | Named when | Finding |
|---|---|---|---|
| any status but deprecated | gone | always | `ID, in FILE at REV, is gone; an entry is deprecated first and deleted in a later change` |
| deprecated | gone | never | (retiring's second step) |
| active | planned | always | `ID went from active to planned; a built decision is deprecated with a Changes line, not unbuilt` |
| not deprecated | deprecated | no new Changes line | `ID was deprecated with no Changes line saying what replaces it` |
| title, text, Kind or a scenario | differs | no new Changes line | `ID changed since REV (text, scenario X reworded, …) with no Changes line saying what and why` |
| a scenario `(owner…)` or `(delegate…)` | the same handle, provisional | always | `ID/handle lost its ruling: it was ruled (owner, …) and is now an agent's, provisional` |

A new Changes line is a dated item whose date and text together the
entry did not have at the revision ([`hasNewChange`](../internal/diff/diff.go#L202)).
[`whatChanged`](../internal/diff/diff.go#L215) compares the title, the text and
the Kind, and each scenario by handle: added, removed, or its sentence
reworded. An entry added since the revision, an entry only set active,
and a change with a new line are not named. diff exits 1 when it names
anything, 0 otherwise.

From a copy of `examples/ratelimit` committed as `HEAD`, with one sentence
of `SW-261126-VX8Z`'s text changed and no Changes line:

```text
$ charter diff HEAD
claims/limiter.md:8: history: SW-261126-VX8Z changed since HEAD (text) with no Changes line saying what and why
docs/walkthrough.md: docs: SW-261126-VX8Z changed since HEAD; read the doc again
$ echo $?
1
```

The second line is a doc to read again, not a finding; the exit status is
the first line's.

## The docs to read again

When charter.toml names docs, [`run`](../cmd/charter/main.go#L47) loads the
tree with [`Load`](../internal/check/check.go#L56) and calls
[`Docs`](../internal/diff/diff.go#L82), which compares each doc's references
by the same hashes curation records ([docs](docs.md)):

- an entry reference or embed, by the hash of the entry's title, status and
  text at the revision and now, so a new Changes line alone does not list
  the doc;
- a code reference, by the hash of the declaration's source at the
  revision and now. Each referenced file is written from `git show` into a
  temporary directory, read there by its language's reader to find the
  declaration's lines at the revision, and hashed; a file that was not
  there is a change.

Each changed reference is a line `DOC: docs: KEY changed since REV; read
the doc again`, the key an id or `file#Name`. The list never changes diff's
exit status, and if it cannot be made (a reader that fails, say) diff
prints `charter diff: the docs to read again are not listed: …` to
standard error and its status stands.

## Where it runs

The loop runs `charter diff origin/main` before every push; in this
repository that is `go run ./cmd/charter diff origin/main`. The
claims-review skill reads every entry diff names, since an entry with no
Changes line has no line to judge ([the loop](loop.md)).

## Traps

- Only the area documents are compared. A change to `observations.md`,
  `questions.md`, the waivers or the baseline is never named.
- An entry's `Why:`, `Guidance:`, fields other than `Kind` and `Status`,
  and a provenance going from provisional to ruled are not compared, so a
  change to them alone is not named.
- Retiring is two changes: deprecate with a Changes line, strip the code
  and tests, then delete the entry in a later change. Deleting an active
  entry in one change is named, whatever its Changes said.
- The text is compared line by line as written, so rewrapping an entry's
  text is a change diff names, though a doc's hash, which normalizes
  whitespace, does not change. A scenario's continuation lines are joined
  with single spaces, so rewrapping a scenario is no change.
- A Changes line is matched by its date and text together, so rewording
  an old line counts as a new one.
