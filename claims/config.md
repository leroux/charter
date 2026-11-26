# The configuration

Code: cmd/charter

How the charter command finds a charter, and what claims/charter.toml
says that the tree cannot.

## CH-261126-33GK charter finds claims/charter.toml by walking up from the working directory, and refuses a key it does not know

Status: planned
Serves: CH-002, CH-261126-NS7T
See: CH-036, CH-037

The command walks up from the working directory to the first directory
holding claims/charter.toml, and that directory is the root; with none it
fails, naming where it started. The file holds the charter's name; each
language it reads, a table with the directory its tests run in, the
command that runs them, its constants module and, for Go, the generated
module's path; the stakeholder tests, the probe path, the ignore and only
patterns and the environment patterns; the digest; and the declared
surface kinds. A key it does not know, a value of the wrong type, or a
language table without its test command fails naming the line. Its paths
are relative to the root.

Why: a charter configured in Go code can be read only by a Go build, and a
command that must be told where the charter is gets run in the wrong
place; a file beside the record, found the way git finds a repository,
works in any repository and from any directory in it.

Scenarios:
1. `found-walking-up` (owner, backlog 003): run from a directory three
   below the root, charter reads the charter whose claims/charter.toml is
   at the root.
2. `nearest-wins` (owner, backlog 003): run inside a nested charter's
   tree, charter reads the nested charter, not the one around it.
3. `none-found-named` (owner, backlog 003): with no claims/charter.toml
   above the working directory, charter fails naming the directory it
   started from.
4. `keys-read` (owner, backlog 003): a file setting every key yields each
   value, its paths taken relative to the root.
5. `unknown-key-refused` (owner, backlog 003): an unknown key, a value of
   the wrong type, and a language table without its test command each fail
   naming the line.
