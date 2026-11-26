kind: neutral
A document with a key that starts with `@` breaks lookups. On `{"@key":"value"}`, `Get(json, "@key")` should return `value` (there is no modifier named `key`, so it is an ordinary key), and the escaped form `\@key` should return `value` too. At the moment this misbehaves (I have seen a panic on such keys). Please fix and add tests.
