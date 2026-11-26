---
name: claims-implement
description: Take charter entries to covered. Set them active, run make claims, build to their text and Why with the repository's implementation and testing skills, write a test per scenario and cite it, anchor the declaration that carries the decision, run the gate, lint and the report, then run claims-review. Use when a story's entries are ready to build, e.g. "implement SW-118", "build backlog 038". Not drafting (claims-draft) or reviewing (claims-review).
---

# Implement entries

The target is an entry, or a story and the ids its `## Charter` section
lists. Covered means: an anchored declaration realizes the entry, a test
cites each of its scenarios, the gate is green and the report says the
tests passed. The form is `claims/README.md`.

## Procedure

1. **Set the entries active.** In the area document, `Status: active`.
   Setting active is not a change and needs no Changes line; add one
   (`- YYYY-MM-DD: what and why`) only if you change the text or the
   scenarios. Set the story's status line to `in progress`. `make claims`,
   so `claims.SW_118` exists in Go and TypeScript. A `deprecated` target:
   stop and report.

2. **Read the intent.** For each id, read the entry's text, Why and
   scenarios, then `go run ./cmd/claimsgen anchors SW-118`: the code
   that realizes or concerns it, the tests that verify it and the handles
   each names. Work from the text, not from what the code does now; an
   anchor is a claim, not a proof, and where the code has drifted the text
   is what shows it.

3. **Build and test** with the repository's own skills:
   `superpowers:test-driven-development` for the loop,
   `superpowers:systematic-debugging` when something fails,
   `superpowers:verification-before-completion` before presenting. Write
   one test per scenario where the boundary allows; one test may cover
   several scenarios, of one entry or several. Each test must fail if its
   scenario were violated; one that cannot is decorative, so reconsider it
   or the scenario. A need is verified from the person's side only: the test
   files the gate's `Options.StakeholderTests` names, such as the
   command driven as a person would drive it, the server's routes or a
   browser suite.

4. **Cite the scenarios** as the first statement of each test, one call
   per spec, naming the handles the test exercises:
   - Go: `chartertest.Verifies(t, claims.SW_118, "branch-keeps-outcome")`,
     importing the repository's entries package (its `claims`) and
     `github.com/leroux/charter/chartertest`.
   - TypeScript: `verifies(claims.SW_118, "handle")`, `claims` from the
     generated module by the import the gate's `Options.TypeScript.Import`
     names, `verifies` from the repository's test helper.
   A handle the entry does not list, a non-literal handle, a call
   outside a test function or anywhere but the test's (or its subtest's)
   leading statements, a `verifies` in a skipped TypeScript test or one
   the test file declares, fails the gate.

5. **Anchor the code.** `//charter:realizes SW-118` as the last line of
   the doc comment of the declaration that carries the decision
   (TypeScript: `// charter:realizes SW-118` above it); several ids with
   commas. `//charter:concerns SW-118` marks code coupled to it without
   realizing it; navigation, not coverage. A package clause carries a
   realizes only for a `structure` decision; lint names any other
   decision realized only by a package clause. A need gets no anchor.
   Never place an anchor you are unsure is true: ask and leave the gap.

6. **Surfaces.** A new environment variable is a constant in an env
   package the gate's `Options.EnvPackages` names; a new operation or event
   kind is in the contract `Options.OpenAPI` or `Options.Events` names; each is named
   in the `Accounts` of the decision that owns it, at the altitude of what
   the operator can rely on. A surface now accounted for comes off
   `claims/baseline.txt` in the same change; the baseline only shrinks.

7. **Waive** only a scenario whose test genuinely cannot exist yet, with
   the user's agreement: a line `SW-118/handle<TAB>reason` in
   `claims/waivers.txt`. Never an anchor, never to go green.

8. **Gate, lint, report.**
   - `go test -count=1 ./claims`: fails an active decision with no
     realizes anchor, a scenario no test cites, a stale generated file.
   - `go run ./cmd/claimsgen lint`: nothing but the line counting the
     open questions.
   - `go run ./cmd/claimsgen report -md`: runs the whole Go suite and
     says, per scenario, which tests cited it and whether they passed;
     `-go FILE` reads a saved `go test -json` stream instead. Every
     targeted scenario is passed, not skipped or none.
   - The TypeScript tests you wrote, if any, pass and typecheck with the
     repository's own commands.
   - `claimsgen mutate` is not a step here. It is opt-in: run it only when
     the story asks for it or a reader doubts that a test would notice
     its decision broken (`claims/README.md`, "The gate and the report").

9. **Review, then present.** Run `claims-review` over the entries this
   work touched and disposition every finding. Present the area document
   diff, the code, the tests and the report. Do not commit. When the user
   accepts, set the story `done` in its file and in `backlog/README.md`.

## Retiring an entry

Set `Status: deprecated` with a Changes line naming what replaces it, if
anything; `make claims`, and its symbol carries a deprecation notice so
everything still builds. `claimsgen lint` lists the anchors and citations
that remain: remove them with the code and tests they served, dropping a
directive whose last id is gone. Keep the entry; the record is the point.

## Done when

- Each target is `active`, regenerated, realized by an anchored
  declaration, with every scenario cited by a test that can fail.
- The gate, lint and the report are green; a gap outside the target was
  reported, not anchored or waived.
- `claims-review` ran and its findings were dispositioned; nothing is
  committed.
