# Adopting charter in a repository

This guide takes a repository from no charter to a green gate. Each step
names the file or command, and points at the two worked examples in this
repository, which are charters of their own:

- `examples/ratelimit`: one Go module at its root, a need verified from
  the command's side, an environment constants package, a digest and a
  curated doc.
- `examples/slugs`: no `go.mod` at its root; a Go module in `svc/` and a
  Python package in `tool/` that must agree on one decision.

Neither imports anything of the tooling. The command is a binary, and the
files it generates into the repository's own `claims/` import only their
language's standard library ([CH-261126-NS7T](../claims/needs.md#ch-261126-ns7t-an-adopter-must-be-able-to-keep-a-charter-in-any-repository-layout-and-language-with-code-that-imports-nothing-of-the-tooling)).

## 1. Install the command

```sh
go install github.com/leroux/charter/cmd/charter@latest
```

or pin it in a tools module. The repository is private for now, so `go
install` needs git credentials that can read it. The readers need, besides
the binary, `python3` on `PATH` when the charter reads Python, and `node`
with the repository's own `typescript` package when it reads TypeScript
([CH-261126-PQ7X](../claims/adoption.md#ch-261126-pq7x-charter-must-support-each-language-as-a-built-in-reader-generator-and-marker-and-must-read-a-language-only-when-chartertoml-names-it)).

## 2. Write claims/charter.toml

The command finds the charter by walking up from the working directory to
the first directory holding `claims/charter.toml`; that directory is the
root and every path in the file is relative to it
([CH-261126-33GK](../claims/adoption.md#ch-261126-33gk-charter-must-find-claimschartertoml-by-walking-up-from-the-working-directory-and-must-refuse-a-key-it-does-not-know)). A key the command does not know
is refused, with its line. ratelimit's:

```toml
name = "ratelimit"
stakeholder-tests = ["cmd/ratelimit/*_test.go"]
env-patterns = ["^RATELIMIT_[A-Z0-9_]+$"]
digest = "DIGEST.md"
docs = ["docs/*.md"]

[go]
module = "example.com/ratelimit/claims"
dir = "."
test = "go test -count=1 ./..."
env = "env"
```

slugs', with Go and Python in their own directories:

```toml
name = "slugs"
env-patterns = ["^SLUGS_[A-Z0-9_]+$"]

[go]
module = "example.com/slugs/claims"
dir = "svc"
test = "go test -count=1 ./..."

[python]
dir = "tool"
test = "PYTHONPATH=../claims python3 -m unittest -q"
env = "tool/slugs/env.py"
```

| Key | What it says |
|---|---|
| `name` | the charter's name |
| `[go]` | `module`, the import path of the generated claims module; `dir`, where the tests run; `test`, the command that runs them; `env`, the constants package's directory |
| `[python]` | `dir`, `test`, and `env`, the constants module's file |
| `[typescript]` | `dir`, `test`, `env`; `compiler`, the typescript package directory the reader parses with; `import`, the specifier tests import the generated module by |
| `stakeholder-tests` | patterns of the test files a need may be cited from; unset, any test |
| `probes` | the paths that hold probes, tests of observations about outside systems |
| `ignore`, `only` | patterns the walk skips, and the paths it is limited to |
| `env-patterns` | regular expressions a string literal naming an undeclared variable matches |
| `digest` | the generated digest's file name under `claims/` |
| `[[surface]]` | a declared surface kind: `kind`, a lower-case word, and `command`, whose output lines are its keys |
| `[[files]]` | `patterns` and a `comment` marker for files no language reads whose comment lines may anchor an entry |
| `docs` | patterns of the markdown files read as docs |

A language table is read only when it is there: a charter with no
`[python]` never walks a `.py` file and never starts `python3`.

## 3. Write the needs and the first area

`claims/needs.md` holds the needs, what a person can rely on. Each other
`.md` file in `claims/` is an area document of decisions: a `# Title`, a
`Code:` line naming the code it covers, a sentence on what it holds, then
the entries. `claims/reading-order.txt` names the area documents, one a
line, in the order a reader meets them. Take every id from the command:

```sh
charter next STK      # a need
charter next SW 2     # two decisions
```

An entry from slugs, as the record reader wants it
([CH-010](../claims/intent.md#ch-010-the-record-reader-must-read-an-entry-as-a-heading-with-id-and-title-then-field-lines-text-why-scenarios-guidance-and-changes-in-that-order)):

```markdown
## SW-261126-X6HB The tool takes another separator from SLUGS_SEPARATOR

Status: active
Serves: STK-261126-78RN
Accounts: env:SLUGS_SEPARATOR

The tool joins words with the value of SLUGS_SEPARATOR when it is set, and
with a hyphen otherwise.

Why: a site that links with underscores keeps the tool's slugs its own
without a fork of the tool.

Scenarios:
1. `separator-from-env` (owner, 2026-11-26): with SLUGS_SEPARATOR set to
   an underscore, "Hello World" gives hello_world.
```

Start new entries `Status: planned`: a planned entry has its symbol, so
tests can cite it, and the gate asks nothing else of it until it is set
active ([CH-013](../claims/intent.md#ch-013-an-entrys-status-must-be-planned-active-or-deprecated-every-status-must-generate-a-symbol-and-the-gate-must-ask-coverage-of-active-entries-only)). The form an agent writes entries
in comes from the repository's CLAUDE.md or AGENTS.md; the claims-draft
skill has a default.

## 4. Generate and wire the claims module

```sh
charter gen
```

writes, for each language table, the symbol file into `claims/`:
`go.mod` and `specs.go`, `claims.py`, `index.ts`, and the digest when
`digest` is set ([gen](gen.md)). Commit them; check fails a tree where
they differ from what the record renders. Then make the tests see them:

- **Go.** The module that holds the tests requires the claims module and
  replaces it with its directory. ratelimit's root `go.mod`:

  ```text
  require example.com/ratelimit/claims v0.0.0

  replace example.com/ratelimit/claims => ./claims
  ```

  slugs' `svc/go.mod` points one level up, `=> ../claims`.
- **Python.** The generated module is `claims`; put `claims/` on the
  path the tests import from, as slugs' `test` command does with
  `PYTHONPATH=../claims`. Under pytest, also list `claims` in
  `pytest_plugins`.
- **TypeScript.** Map the specifier `[typescript].import` names to
  `claims/index.ts` in the module resolution your tests use, and run
  vitest with `globals`, which the marker records through.

## 5. Anchor the code

Put the directive on the declaration that carries the decision
([CH-021](../claims/anchoring.md#ch-021-production-code-must-anchor-an-entry-in-a-go-doc-comment-a-typescript-line-comment-or-the-comment-lines-above-a-python-def)): the last line of a Go doc comment,
the comment lines directly above a Python `def` or `class`, or a
TypeScript line comment. From slugs' [slug.go](../examples/slugs/svc/slug/slug.go):

```go
// Of is the slug of title: its letters and digits in lower case, each run
// of them joined by one hyphen.
//
//charter:realizes SW-261126-4NQQ
func Of(title string) string {
```

and from its [slug.py](../examples/slugs/tool/slugs/slug.py):

```python
# charter:realizes SW-261126-X6HB
def separator():
```

A need takes no anchor. `charter:concerns` marks code coupled to a
decision without realizing it, and counts for nothing.

## 6. Cite the scenarios

A test opens with its citation, naming the handles it exercises
([CH-023](../claims/anchoring.md#ch-023-a-test-must-cite-an-entry-by-a-verifies-call-that-opens-a-test-function-naming-by-handle-the-scenarios-it-exercises)). slugs' Go test and Python test of the
same decision:

```go
func TestASlugIsLowerCaseWordsJoined(t *testing.T) {
	claims.Verifies(t, claims.SW_261126_4NQQ, "lower-and-joined", "no-edge-separator")
```

```python
class SlugTest(unittest.TestCase):
    def test_lower_and_joined(self):
        claims.verifies(self, claims.SW_261126_4NQQ, "lower-and-joined", "no-edge-separator")
```

A pytest function passes `None` for `self`; vitest calls
`verifies(claims.SW_…, "handle")`. With `stakeholder-tests` set, a need is
cited only from those files: ratelimit cites its needs from
`cmd/ratelimit/main_test.go`, which drives the command as a caller would
([CH-031](../claims/gate.md#ch-031-the-gate-must-not-ask-an-anchor-of-a-need-and-must-count-a-needs-citation-only-from-the-stakeholders-side)).

## 7. Name the environment

Every environment variable is a string constant of the constants module
the language table's `env` names, such as ratelimit's
[env.go](../examples/ratelimit/env/env.go) and slugs'
[env.py](../examples/slugs/tool/slugs/env.py), and code reads it through
the constant: `os.environ.get(env.SEPARATOR)`. Each variable is a surface
some entry names in `Accounts: env:NAME`; slugs' `SW-261126-X6HB` owns
`env:SLUGS_SEPARATOR` ([CH-033](../claims/gate.md#ch-033-the-gate-must-enumerate-every-surface-and-require-each-to-be-accounted-for-by-an-entry-or-listed-in-a-baseline-of-live-unaccounted-surfaces),
[CH-034](../claims/gate.md#ch-034-code-must-read-an-environment-variable-through-a-named-constant-or-under-a-dynamic-directive-and-a-literal-must-not-name-one)). Another kind of surface, such as the routes
of an OpenAPI file, is a `[[surface]]` table whose command prints one key
a line.

## 8. Run the gate, the report and diff

```sh
charter check               # the gate: problems on stdout, exit 1; warnings on stderr
charter report -md          # runs every test command, says which scenarios passed
charter diff origin/main    # names a change to the record with no Changes line
```

Run all three before pushing and in CI ([check](check.md),
[report](report.md), [diff](diff.md)). `make check` in this repository
does the same for both examples.

## 9. An existing repository

A repository with code and no record starts with surfaces nobody has
ruled on. Run the gate once with

```sh
CHARTER_UPDATE_BASELINE=1 charter check
```

and it writes `claims/baseline.txt`, every unaccounted surface with the
reason `unclaimed at baseline init`; the gate then fails a line once an
entry accounts for its surface or the surface is gone, so the baseline
only shrinks. A scenario whose test cannot exist yet may be waived, one
line `SW-…/handle`, a tab and the reason, in `claims/waivers.txt`; a
missing anchor never can be ([CH-030](../claims/gate.md#ch-030-the-gate-must-fail-an-active-decision-without-an-anchor-or-with-an-uncited-scenario-a-waiver-may-absorb-a-missing-test-never-a-missing-anchor)). `only` limits
the walk to the paths the charter covers so far.

## 10. The skills and the agent instructions

The loop's skills ship as the `charter` plugin of this repository
([CH-261126-5JMD](../claims/adoption.md#ch-261126-5jmd-the-loops-skills-must-ship-as-a-plugin-of-this-repository-and-must-read-the-adopters-layout-from-chartertoml-alone)):

```text
/plugin install charter --marketplace leroux/charter
```

They run as `/charter:claims-draft`, `/charter:claims-implement`,
`/charter:claims-review`, `/charter:claims-audit`, `/charter:claims-probe`
and `/charter:claims-narrate`, and read the layout from
`claims/charter.toml`. The repository's README holds the section to put in
the adopter's CLAUDE.md, with its prefix and paths; [the loop](loop.md)
says what each skill does in a change.

## 11. Docs

Add `docs = ["docs/*.md"]` and write the prose over the record and the
code there, with references the gate checks; the claims-narrate skill
drafts them and a person curates each with `charter curate -by NAME FILE`
([docs](docs.md)). ratelimit's
[walkthrough.md](../examples/ratelimit/docs/walkthrough.md) is one.

## Traps

- Each example holds its own `claims/charter.toml`, so a charter around
  it never walks it: a repository can keep several charters, each the
  nearest above its code.
- A Python test is a file named `test_*.py`, `*_test.py` or
  `conftest.py`; a TypeScript test ends `.test.ts(x)` or `.spec.ts(x)`.
  Anything else is production code. A `verifies` there counts for
  nothing, and in Go and Python the gate prints nothing at the call
  either: a test in a misnamed file shows only as its scenario's `rule 25`,
  with no word of where the call went.
- Passing `os.Getenv` around as a function value, as ratelimit's `main`
  does, hides the reads from the environment check; the literal check by
  `env-patterns` still applies.
