# Questions

What nobody ruled about the charter's own tooling, in the form of
`claims/questions.md`: a heading `## Q-ID question`, the id from `go run
./cmd/claimsgen next Q`, then `Asked-by:`, `Entry:` (a CH id, or
`ID/handle`), `Provisional:` and, once a person has answered, `Answer:
YYYY-MM-DD, who: ...`. A question about an entry of this charter is
written here, since the gate holds each question file to its own
charter's entries. `go run ./cmd/claimsgen lint` counts the open ones
and `go run ./cmd/claimsgen questions` lists them.

## Q-261015-5QD4 Is the claims viewer, with claimsmap's export of a charter, kept as the tooling's way to read a charter as pages, or deleted, the digest read as markdown being enough (is a viewer wanted at all)?
Asked-by: agent
Entry: CH-261015-9ZTY
Provisional: the viewer and claimsmap stay as they are and nothing is deleted until the owner rules.
Answer: 2026-10-15, owner: kept, as an observer of the claims and the charters; CH-261015-9ZTY rules what a reader can rely on of it.

## Q-261015-N4HJ Should the charter tooling stop reading the earlier one-file-per-spec form (parseSpec, readSpecs, the old lens list), deprecating CH-014 and changing CH-015 `stakeholder-realizes-nothing`, now that neither charter has a spec file left?
Asked-by: agent
Entry: CH-014
Provisional: nothing is cut or changed until the owner rules; recommended: Yes. Rewrite the chartertest fixture as an area document in its own PR first, then delete. The entries it touches besides this one: CH-015.
Answer: 2026-11-26, owner: yes; the earlier form is no longer read, CH-014 is deprecated and CH-015's stakeholder-realizes-nothing dropped (backlog 003, ruling 26).

## Q-261015-CPW5 Should the generated claims/specs.go give each symbol its entry's title as its doc comment, not the title and the text, narrowing CH-020's text to match its title and claims/README.md:301?
Asked-by: agent
Entry: CH-020
Provisional: nothing is cut or changed until the owner rules; recommended: Yes. About 6,360 generated lines go.
Answer: 2026-11-26, owner: yes; a generated symbol's doc comment is its entry's title (backlog 003, ruling 28).

## Q-261015-Y682 Should the gate stop writing generated files when CHARTER_UPDATE_SPECS, CHARTER_UPDATE_TYPESCRIPT or CHARTER_UPDATE_DIGEST is set, leaving `make claims` (Generate) the one writer, and drop those Accounts from CH-020 and CH-261005-6ZBS?
Asked-by: agent
Entry: CH-020
Provisional: nothing is cut or changed until the owner rules; recommended: Yes; make Options.Regenerate required and have the tests call Generate. The entries it touches besides this one: CH-261005-6ZBS.
Answer: 2026-11-26, owner: yes; charter gen is the one writer, and the update variables go from CH-020 and CH-261005-6ZBS (backlog 003, ruling 27).

## Q-261015-X96S Should a Changes line say only what changed in the ruling, why and its provenance, in about three lines, with how a review found it left to the PR, so that old Changes bullets may be shortened and the 'reworded; the ruling is unchanged' bullets deleted, relaxing 'append-only' in CLAUDE.md and claims/README.md?
Asked-by: agent
Entry: CH-261009-QQC8
Provisional: nothing is cut or changed until the owner rules; recommended: Yes, in order: file web.md's eight 'Not closed' defects, tighten hasNewChange (CH-261009-QQC8), then cap the bullets at three lines (1,855) and delete the 237 pure 'reworded, the ruling is unchanged' bullets and the 89 empty headings they leave (532).

## Q-261015-YSKC Should an answered question leave claims/questions.md once its ruling is written into the record (a scenario's provenance or a Changes line naming its Q-id), changing CH-261009-4B3B's 'A question is never deleted', with git keeping the text?
Asked-by: agent
Entry: CH-261009-4B3B
Provisional: nothing is cut or changed until the owner rules; recommended: Yes, the 73 whose Q-id is already written into the record; the other 102 stay until their rulings are built.

## Q-261015-5YV3 Should `claimsgen lint -baseline` and claims/style-baseline.txt be retired, deprecating CH-261009-SNX2, given style is advice (2026-10-11), nothing runs the ratchet, and it exits 1 on main because four counts grew?
Asked-by: agent
Entry: CH-261009-SNX2
Provisional: nothing is cut or changed until the owner rules; recommended: Yes; move lintFlags' and runLint's anchors to CH-041. The procedure stops using lint -baseline either way. The entries it touches besides this one: CH-041.
Answer: 2026-11-26, owner: yes; lint -baseline and the style baseline are retired and CH-261009-SNX2 is deprecated (backlog 003, ruling 29).

## Q-261015-3Q84 Is go/charter meant to be adopted by other repositories? If not, should the surface baseline (claims/baseline.txt, empty and forbidden to grow) and its reader, rebuilder, CHARTER_UPDATE_BASELINE and CI step be retired, changing CH-033 and CH-004?
Asked-by: agent
Entry: CH-033
Provisional: nothing is cut or changed until the owner rules; recommended: Ask whether go/charter is meant for other repositories; if not, yes; if it is, keep the baseline as the cost of adoptability. The entries it touches besides this one: CH-004.
Answer: 2026-11-26, owner: yes, the charter is meant for other repositories and new ones; it is its own repository, github.com/leroux/charter, a project and a methodology of its own. The baseline stays as the cost of adoptability, and CH-033 and CH-004 are unchanged.
