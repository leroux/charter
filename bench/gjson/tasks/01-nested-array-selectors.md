kind: neutral
An array path inside an array path drops its results. On

    [{"c":[{"a":10.11}]},{"c":[{"a":11.11}]}]

`@ugly|#.c.#[a=10.11]` should return `[{"a":10.11}]` (the matches, from the elements that have any), and `@ugly|#.c.#` should return the count of each inner array, `[1,1]`. Today both come back wrong. Please fix and add tests.
