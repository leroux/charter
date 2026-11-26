kind: neutral
Chaining does not work after an array that collects many results. On a document with `info.friends` holding two objects that both have `"kind":"Person"`:

- `info.friends.#[kind="Person"]#.kind|0` and `info.friends.#.kind|0` should both return `Person` (the pipe applies to the array of kinds),
- `info.friends.#[kind="Person"]#.kind` and `info.friends.#.kind` should return `["Person","Person"]`,
- `info.friends.#[kind="Person"]#|kind` and `info.friends.#|kind` should return nothing (the pipe reads `kind` from an array, which has no such key).

Please fix and add tests.
