# coreutils

Go clones of GNU coreutils' `cut`, `sort`, `uniq` and `wc`, meant as
drop-in replacements for GNU coreutils 9.4 in the C locale (`LC_ALL=C`):
for every flag a tool supports, its output and exit status are GNU's.

Each tool is its own command, `cmd/<tool>`, built with
`go build ./cmd/<tool>`; what they share lives in `internal/`.

| Tool | Supports |
|---|---|
| `wc` | `-l`, `-w`, `-c` (and `--lines`, `--words`, `--bytes`) |
| `cut` | `-f LIST` with `-d DELIM` (and `--fields`, `--delimiter`) |
| `sort` | plain, `-n`, `-r` (and `--numeric-sort`, `--reverse`) |
| `uniq` | plain, `-c` (and `--count`), one input file |

Every tool reads standard input when no file is named, else the named
files. `go test ./...` runs the tests.
