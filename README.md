# Charter

Charter records what a system does as **needs** (what a person can rely
on) and **decisions** (how the system behaves at a boundary), each written
once in a markdown area document, anchored to the code that realizes it and
verified, scenario by scenario, by the tests that name it. A gate, run as
an ordinary Go test, fails the tree when an active decision has no code
that claims it, a scenario no test cites, or a surface (an environment
variable, an operation, an event) no decision accounts for.

The gate proves the links exist; it does not prove the code right or a test
meaningful. That is what review is for, and the charter gives the reviewer
the entries and their tests to read instead of the code.

It is meant for any repository, existing or new.

## What is here

| Path | What it is |
|---|---|
| `charter.go` | The types a generated entries package declares: `Spec`, `Observation` |
| `chartertest/` | The gate (`Check`), the registry, generation, the report, the digest, mutation testing |
| `cli/` | The command line, for an adopter's own `claimsgen` and `claimsmap` |
| `cmd/claimsgen`, `cmd/claimsmap` | Those commands over the charter's own charter |
| `claims/` | The charter's own charter: the tooling specified by the method it implements |
| `example/` | **ratelimit**, a small adopter kept as a real one would be: its own module, requiring this one |
| `claims-viewer/` | A web viewer of an exported charter |
| `.claude/skills/claims-*` | The agent skills of the loop: draft, implement, review, audit, probe |
| `backlog/` | The work toward a standalone method |

## Adopting it

`example/` is the worked example; every piece below is a file there.

1. Require the module: `go get github.com/leroux/charter`.
2. Make an entries package, `claims/`, with a `doc.go`, `needs.md`, an area
   document per area of the product (`# Title`, a `Code:` line, a summary,
   then entries), and `reading-order.txt` naming the area documents.
3. Give the gate its options in a package of its own, `claims/gate`:

   ```go
   var Options = chartertest.Options{
       Entries:     "example.com/ratelimit/claims",
       EnvPackages: []string{"example.com/ratelimit/env"},
       EnvLiterals: []string{`^RATELIMIT_[A-Z0-9_]+$`},
       Regenerate:  "go run ./cmd/claimsgen",
       Digest:      "DIGEST.md",
   }
   ```

4. The gate is a test, `claims/gate_test.go`:
   `chartertest.Check(t, gate.Options)`.
5. Your own `claimsgen` is five lines, `cmd/claimsgen/main.go`:

   ```go
   cli.Claimsgen(os.Args[1:], cli.Charter{Name: "product", Dir: "claims", Options: gate.Options})
   ```

6. `go run ./cmd/claimsgen` writes `claims/specs.go` (a symbol per entry)
   and the digest; `go run ./cmd/claimsgen next SW` gives a new id.
7. Anchor code with `//charter:realizes SW-…` as the last line of a
   declaration's doc comment, and open each test with
   `chartertest.Verifies(t, claims.SW_…, "handle")`.
8. `go test ./...` runs the gate. `claimsgen lint`, `report -md`,
   `mutate -md ID` and `diff origin/main` do the rest of the loop.

## Working on the charter

```sh
pnpm install          # the TypeScript compiler the citation extractor uses
make check            # Go tests and the charter's own gate, the example, the viewer
go run ./cmd/claimsgen lint
go run ./cmd/claimsgen report -md
```

The charter specifies itself: a change to the tooling is a change to an
entry in `claims/`, with a dated `Changes:` line, and `go run
./cmd/claimsgen diff origin/main` names one that has none. The example
adopter's gate runs among the tooling's tests (`adopter_test.go`), so a
change that would break an adopter fails here first.
