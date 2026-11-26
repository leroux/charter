---
name: claims-probe
description: Run the probes of one subject or all against their live external systems, read the result per observation, and write the Observed date or propose Broken with the entries now suspect. Use when something external broke, when a fact needs revalidating, or occasionally.
---

# Probe a subject

Input: a subject named in `claims/observations.md` (the `Subject` of an
observation, such as `payments-api`), or `all`. An observation is a heading
`## Name` with the fields Subject, Status (observed, documented, broken or
unsure), Observed, Source and Supports, the ids of the entries that rest on
it; its text says what was seen.

Probes touch real systems and may cost money. Confirm with the user which
subjects to reach and that the credentials in the environment are the ones
meant, before running anything.

## 1. Run

A probe is a test under one of the paths `charter.toml`'s `probes` key
names, opening with an observes call:

- Go: `claims.Observes(t, claims.Name)`
- Python: `claims.observes(self, claims.Name)` (unittest) or
  `claims.observes(None, claims.Name)` (pytest)
- TypeScript: `observes(claims.Name)`

Run the probes of the subject with the language's own test command limited
to the probe paths, with `CHARTER_EVENTS` set so each probe's outcome is
written down:

```
CHARTER_EVENTS=/tmp/charter-probe.jsonl <the language's test command> <the probe paths>
```

The repository's CLAUDE.md or Makefile says how its probes choose their
subjects, often an environment variable. A subject with no probe has
nothing to run; its observations stay `documented` with their Source. Each
subject needs what its probe reaches, credentials or a running service,
which the probe's comment names.

## 2. Read

Each line of `/tmp/charter-probe.jsonl` is one observes call: `id` is the
observation's name, `test` the probe, `outcome` passed, failed or skipped.
Print a table: observation, subject, result, and the reason or the failure
from the test command's output. An observation with no line was probed by
nothing: its probe is missing, or it does not open with its observes call.

## 3. Write

In `claims/observations.md`:

- Pass: set `Observed: <today>`; if its Status was `unsure` or
  `documented`, set it to `observed`.
- Fail: propose `Status: broken`, with the failure quoted in the
  observation's text, and list the entries in its Supports with their
  status and title (`charter anchors ID...`): those are now suspect. Do not
  change an entry. Whether to fix the code, amend the entry (with a Changes
  line) or deprecate it is the user's ruling.
- Skip: nothing; print the reason.

Run `charter gen` and `charter check`. The gate fails while a broken
observation supports an active entry; leave it red and hand the ruling to
the user.

## 4. Commit

Only if the user asks: `claims: probe <subject>` with the table in the body.
