---
name: claims-probe
description: Run the probes of one subject or all against their live external systems, read the result per observation, and write the Observed date or propose Broken with the entries now suspect. Use when something external broke, when a fact needs revalidating, or occasionally.
---

# Probe a subject

Input: a subject named in the observations in `claims/<subject>.go`
(the `Subject` of an observation, such as `payments-api`), or `all`. The form is
`claims/README.md`, "Observations": an observation's `Supports` lists the
ids of the entries that rest on it (`"CON-002"`).

Probes touch real systems and may cost money. Confirm with the user which
subjects to reach and that the credentials in the environment are the ones
meant, before running anything.

## 1. Run

```
<PROBE_SWITCH>=<subject[,subject]> go test -tags <probe tag> -count=1 -run '^TestProbe' -json ./... > /tmp/charter-probe.json
```

The probe tag is the gate's `Options.ProbeTag`, and the switch is the
environment variable the repository's probes read to choose their
subjects; its CLAUDE.md or Makefile names both. Use the command above for
a stream you can read. The switch has no `all` keyword: for `all`, pass
every subject the observations name. A subject with no probe file has
nothing to run; its observations stay `charter.Documented` with their
`Source`.

Each subject needs what its probe reaches: credentials or a running
service, which the probe file's comment names.

## 2. Read

From the JSON stream, for every `attr` event with key `observation`, take the
test name and the observation id; then take that test's final `pass`, `fail`
or `skip`. Print a table: observation, subject, result, the skip reason or
the failure output. A probe that never appears in the stream is misnamed:
probe functions start with `TestProbe` and cite with
`chartertest.Observes(t, claims.X)`, never `Verifies`.

## 3. Write

- Pass: set `Observed: "<today>"` on the observation; if its status was
  `charter.ObservationUnsure` or `charter.Documented`, set `charter.Observed`.
- Fail: propose `charter.Broken`, with the failure quoted in a comment above
  the observation, and list the entries in its `Supports` with their status
  and title (`go run ./cmd/claimsgen anchors ID...`): those are now
  suspect. Do not change an entry. Whether to fix the code, amend the entry
  (with a `Changes` line) or deprecate it is the user's ruling.
- Skip: nothing; print the reason.

Run `gofmt -w claims && go test -count=1 ./claims`. The gate fails while a
Broken observation supports an active entry; leave it red and hand the
ruling to the user.

## 4. Commit

Only if the user asks: `claims: probe <subject>` with the table in the body.
