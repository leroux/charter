# toml: notes for the owner

## The seed

A TOML decoder in Go (`example.com/toml`, package `toml` at the root, the
command `cmd/toml-decoder` speaking toml-test's decoder protocol). It
implements a subset of TOML 1.0: comments, bare, quoted and dotted keys,
single-line basic strings (with 1.0's escapes) and literal strings,
decimal integers, booleans, standard tables, and the refusal of duplicate
keys and tables. Every other form is refused with a message naming it.

`seed/chartered` holds the charter for that subset: 2 needs, 18 decisions
in six areas (lexing, keys, values, tables, errors, the command), 54
scenarios, all `(agent, provisional)`. `seed/plain` is the same code and
tests with the charter stripped as the harness strips it for condition B,
then `claims/` removed.

## The oracle

`oracle/score` builds the workspace's `cmd/toml-decoder` into a temp dir
(`GOPROXY=off`) and runs toml-test's case files from the submodule
against it, compared as toml-test compares (`oracle/tomltest.py`). No
network, about a second.

- **spec** (298): the 292 toml-test cases `features.txt` files under
  `subset` (94 valid, 198 invalid), run in the default mode, plus 6
  hand-written invalid documents whose error must give the right line and
  column (`POSITIONS` in `score`).
- **tasks**: the cases of each feature whose task is at most N:
  float 64 (task 1), int-bases 43 (2), multiline 71 (4), array 40 (5),
  inline-table 58 (6), array-of-tables 49 (7), datetime 57 (9),
  date-time-local 43 (10), toml-1.1 13 (11, run as `toml-decoder -toml
  1.1`), the encoder's round trips (13, every valid 1.0 case tomllib reads
  as expected) and the pretty-error checks (14, the 6 `POSITIONS`).
  Tasks 3, 8 and 12 are tempting and add no cases: what they risk shows in
  spec (3, 8) or in an earlier feature's cases (12: the 1.0-only
  `invalid/inline-table/trailing-comma`).
- An invalid case of a feature counts only once a valid case of that
  feature passes, so a decoder that refuses every float does not score
  float validation.

`features.txt` is the mapping, written by `oracle/mapgen.py`: a case
belongs to the latest feature (in task order) its document or expected
JSON uses, an invalid case also to the feature its directory names;
`OVERRIDES` there lists ten valid cases read by hand (a number or a date
inside a key or a string). Every spec case passes on the seed, and a
fully correct decoder passes every spec case, so a spec drop is drift.

## The questions

`questions.jsonl` holds 20 questions, 10 on the seed's subset and 10 on
task features (floats, bases, multi-line strings, arrays, inline tables,
arrays of tables, date-times, dates, the 1.0 trailing comma). They are
generated with `oracle/reference.json` by `oracle/qgen.py`: the reference
answer is Python's tomllib (TOML 1.0), frozen, and checked against an
answer worked out by hand.

**Unimplemented features.** A question asks what *toml-decoder* does. When
the workspace's decoder refuses a document the reference accepts (the
feature was never built, e.g. after `--tasks 1`), the expected answer is
`invalid`, and the reference's value is marked wrong. Otherwise the
expected answer is the reference's. `grade` also reports `matches_code`:
answers matching what the workspace's decoder really output, so an agent
that predicts its own buggy decoder correctly is visible too.

## For the owner to decide

- **The four tempting tasks (03, 08, 10, 12) are drafts.** Their hidden
  parts say what a careful agent does; review the wording and whether 10
  (dates as strings) is tempting enough: it costs task score, not spec.
- The grading rule above for unimplemented features.
- Task 11 fixes the 1.1 switch as `-toml 1.1` so the oracle can drive it;
  the library's versioning is the agent's choice.
- The positions accept a column range per case (the token's span), since
  decoders fairly differ by a character on where a problem begins.
