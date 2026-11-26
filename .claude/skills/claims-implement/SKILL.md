---
name: claims-implement
description: Take charter entries to covered. Set them active, run charter gen, build to their text and Why with the repository's implementation and testing skills, write a test per scenario and cite it, anchor the declaration that carries the decision, run the gate, lint and the report, then run claims-review. Use when a story's entries are ready to build, e.g. "implement SW-118", "build backlog 038". Not drafting (claims-draft) or reviewing (claims-review).
---

# Implement entries

The target is an entry, or a story and the ids its `## Charter` section
lists. Covered means: an anchored declaration realizes the entry, a test
cites each of its scenarios, the gate is green and the report says the
tests passed. The layout is the repository's `claims/charter.toml`.

## Procedure

1. **Set the entries active.** In the area document, `Status: active`.
   Setting active is not a change and needs no Changes line; add one
   (`- YYYY-MM-DD: what and why`) only if you change the text or the
   scenarios. Set the story's status line to `in progress`. `charter gen`,
   so `claims.SW_118` exists in each language `charter.toml` names. A `deprecated` target:
   stop and report.

2. **Read the intent.** For each id, read the entry's text, Why and
   scenarios, then `charter anchors SW-118`: the code
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
   or the scenario. A need is verified from the outside only when
   `charter.toml`'s `stakeholder-tests` names test files: the command driven
   as a person would drive it, the server's routes or a browser suite.

4. **Cite the scenarios** as the first statement of each test, one call
   per spec, naming the handles the test exercises:
   - Go: `claims.Verifies(t, claims.SW_118, "branch-keeps-outcome")`,
     importing the generated claims module, the `[go]` table's `module`.
   - Python: `claims.verifies(self, claims.SW_118, "handle")` opening a
     unittest `test_` method, or `claims.verifies(None, ...)` opening a
     pytest `test_` function (pytest also lists `claims` in
     `pytest_plugins`).
   - TypeScript: `verifies(claims.SW_118, "handle")`, both imported from
     the generated module by the `[typescript]` table's `import`.
   A handle the entry does not list, a non-literal handle, a call
   outside a test function or anywhere but the test's (or its subtest's)
   leading statements, a `verifies` in a skipped TypeScript test or one
   the test file declares, fails the gate.

5. **Anchor the code.** `//charter:realizes SW-118` as the last line of
   the doc comment of the declaration that carries the decision
   (Python: `# charter:realizes SW-118` on the comment lines directly
   above the `def` or `class`; TypeScript: a `// charter:realizes SW-118`
   line comment); several ids with
   commas. `//charter:concerns SW-118` marks code coupled to it without
   realizing it; navigation, not coverage. A package clause carries a
   realizes only for a `structure` decision; lint names any other
   decision realized only by a package clause. A need gets no anchor.
   Never place an anchor you are unsure is true: ask and leave the gap.

6. **Surfaces.** A new environment variable is a constant in the
   constants module a language table's `env` names; a key of a declared
   kind (`[[surface]]` in `charter.toml`) comes from that kind's command;
   each is named in the `Accounts` of the decision that owns it, at the altitude of what
   the operator can rely on. A surface now accounted for comes off
   `claims/baseline.txt` in the same change; the baseline only shrinks.

7. **Waive** only a scenario whose test genuinely cannot exist yet, with
   the user's agreement: a line `SW-118/handle<TAB>reason` in
   `claims/waivers.txt`. Never an anchor, never to go green.

8. **Gate, lint, report.**
   - `charter check`: fails an active decision with no realizes anchor, a
     scenario no test cites, a stale generated file.
   - `charter lint`: nothing but the line counting the
     open questions.
   - `charter report -md`: runs each language's test command and says, per
     scenario, which tests cited it and whether they passed. Every targeted
     scenario is passed, not skipped or none.
   - The TypeScript tests you wrote, if any, typecheck with the
     repository's own commands.

9. **Review, then present.** Run `claims-review` over the entries this
   work touched and disposition every finding. Present the area document
   diff, the code, the tests and the report. Do not commit. When the user
   accepts, set the story `done` in its file.

## Retiring an entry

Set `Status: deprecated` with a Changes line naming what replaces it, if
anything; `charter gen`, and its symbol carries a deprecation notice so
everything still builds. `charter lint` lists the anchors and citations
that remain: remove them with the code and tests they served, dropping a
directive whose last id is gone. Keep the entry; the record is the point.

## Done when

- Each target is `active`, regenerated, realized by an anchored
  declaration, with every scenario cited by a test that can fail.
- The gate, lint and the report are green; a gap outside the target was
  reported, not anchored or waived.
- `claims-review` ran and its findings were dispositioned; nothing is
  committed.
