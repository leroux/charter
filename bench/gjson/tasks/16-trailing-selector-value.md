kind: neutral
Issue #141: on

    {"data":[{"q":11,"w":12},{"q":21,"w":22},{"q":31,"w":32}],"sql":"some stuff here"}

`data.#` is `3` and `data.#.q|@ugly` is `[11,21,31]`, but `data.#.{q}|@ugly` should be `[{"q":11},{"q":21},{"q":31}]` and is not: a selector at the end of a `#.` path picks up the wrong value. Please fix and add a test.
