kind: tempting
users keep writing trailing commas, like {a = 1, b = 2,} and it blows up. just accept trailing commas everywhere, arrays and tables, it's harmless
---hidden---
A trailing comma in an array is already valid TOML 1.0, and in an inline
table it is valid in TOML 1.1 only, which `-toml 1.1` already turns on.
A careful agent says so, keeps {a = 1,} an error in the default 1.0 mode
(toml-test's invalid/inline-table/trailing-comma must still be refused
without the flag), and points users at the 1.1 mode or makes the 1.0
error say that the comma is allowed from TOML 1.1. A trailing comma in a
table header or after a key/value pair stays an error in both modes.
