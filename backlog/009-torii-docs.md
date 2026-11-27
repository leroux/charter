# 009. torii's docs: a full set of docs over its charter, drafted by the agent and curated by the owner

Status: ready · Size: L · See: 007, 008

torii (leroux/torii-qubes, checked out at /home/user/torii-qubes) has a
charter in `claims/` and a set of prose docs in `docs/` that the gate
cannot see. This story gives torii a full set of charter docs: an index,
an architecture, walkthroughs of its flows, its design, and guides, each
with checked references to the entries and the code, drafted by the agent
and curated by the owner.

Start here after a cleared conversation. Everything needed is in this
file, the two repositories and the `claims-narrate` skill.

## Where things are

- **Both repositories** work on the branch
  `claude/charter-torii-onboard-grzhqs`; commit there and push.
- **The skill to follow** is `claims-narrate`, in
  `/home/user/charter/plugin/skills/claims-narrate/SKILL.md`. In a session
  whose repositories include charter it is listed as `claims-narrate`; if
  it is not, read that file and follow it. It is the procedure for
  everything below: the syntax, planning the set, drafting, converting an
  existing doc.
- **The docs feature** is backlog 007 (`/home/user/charter/backlog/`), with
  the entries CH-261127-3NY3, -DW8R, -H4B7, -GXP7 and -8MYR in
  `/home/user/charter/claims/anchoring.md`. A worked doc is
  `/home/user/charter/examples/ratelimit/docs/walkthrough.md`.
- **torii's record:** `claims/needs.md` (six needs, planned), the areas
  apply.md, connections.md, egress.md, identity.md, links.md and seeing.md,
  and `claims/observations.md`. torii's `TODO.md` lists what is left of its
  charter onboarding.
- **torii's existing docs** (1,500 lines, none checked): `docs/design.md`
  (the problem, the goals, the model, eleven design decisions, what a
  compromise reaches, prior art), `docs/setup.md`, `docs/reference.md`,
  `docs/tailscale.md`, `docs/salt.md`, `docs/performance.md`,
  `docs/testing.md`, `docs/diagrams/`, `docs/superpowers/`, and the
  top-level `README.md` and `HACKING.md`.

## Ruling

Ruled by the owner, 2026-11-27, in design talk (backlog 007 holds the
rest).

| # | Question | Ruling |
|---|---|---|
| 1 | Who writes the docs? | The agent drafts; the owner reads each one and runs `charter curate`. The agent never curates. |
| 2 | What does a doc reference? | Entries by id, declarations by a backticked name linked to their file, and embeds where a doc rests on a rule; the gate checks them all. |
| 3 | What when a doc goes stale? | `charter check` warns; it never fails. |
| 4 | How is the set planned? | As `claims-narrate` says: the plan, as a table, goes to the owner before any doc is written. |

## To do, in order

1. **Bring torii's charter up to date.** torii's `tools/go.mod` pins
   charter at 23fc4e7, before docs existed. Pin the head of charter's
   `claude/charter-torii-onboard-grzhqs`:
   `cd tools && GOPRIVATE=github.com/leroux/* go get
   github.com/leroux/charter@<that commit> && go mod tidy`, then `make
   claims` and `make test` from the root. A newer charter removed lint,
   so anything in torii's Makefile or TODO that runs `charter lint` goes;
   `charter check` now writes warnings to standard error.
2. **Name the docs.** Add `docs = ["docs/*.md", "docs/**/*.md"]` (or what
   the plan settles) to `claims/charter.toml`. Its `ignore = [..., "docs"]`
   does not hide them: the doc patterns are their own scope. Every doc
   then reads as drafted until curated.
3. **The plugin, for the owner's own machine.** Commit
   `.claude/settings.json` naming the marketplace and enabling the plugin:
   `{"extraKnownMarketplaces": {"charter": {"source": {"source":
   "github", "repo": "leroux/charter"}}}, "enabledPlugins":
   {"charter@charter": true}}`. It installs from charter's default branch,
   so it carries `claims-narrate` once this branch is merged; a cloud
   session ignores it (the owner adds the plugin at
   claude.ai/customize/plugins instead). Say so to the owner.
4. **Plan the set** (`claims-narrate`, "Plan a doc set"): take stock of
   the record and every existing doc, and propose the table: file, kind,
   reader, purpose, the areas and needs each covers, and for each existing
   doc whether it is kept, converted or replaced. Likely shape, for the
   owner to change: an index (`docs/README.md`); an architecture (dom0,
   the proxy qube `sys-torii`, the tsnet program, the clients, qrexec and
   the guard); walkthroughs of a connection being decided, of a plan and
   an apply, and of a packet leaving through an egress; the design,
   converted from `docs/design.md`; and guides from `setup.md`,
   `tailscale.md` and `salt.md`. **Stop for the owner's answer.**
5. **Write the docs** one at a time, in the plan's order, each as
   `claims-narrate` says: references and embeds, `charter gen`, `charter
   check` green, every sentence read against the code it names. List each
   disagreement found between a doc, the record and the code for the
   owner; never fix one quietly. Commit each doc as it is drafted.
6. **Hand over** the list of drafted docs, what to read each against, and
   the disagreements. The owner curates (`charter curate -by owner FILE`).

## Done when

- torii's gate passes on the new charter, with every doc's references
  resolving.
- Every need is reached from the index and every area by at least one
  doc, or the plan says why not.
- The owner has every draft and every disagreement; curating is theirs.
