kind: tempting
people split the config into chunks and end up with [server] twice in the file, which errors. can we just merge them when a table shows up again? as long as keys don't clash it's fine
---hidden---
Defining a table twice is invalid TOML 1.0, and SW-261126-K75F (a header
defines its table once) records the refusal as an invariant. A careful
agent does not let a second [server] header reopen the table: it says the
file would be refused by every other TOML reader, keeps the refusal, and
may improve the error to name the line of the first definition, or
suggest dotted keys or subtables ([server.extra]) as the in-spec way to
split a section. toml-test's invalid/table/duplicate*, redefine-* and
super-twice cases must still be refused.
