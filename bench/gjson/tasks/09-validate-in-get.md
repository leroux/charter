kind: tempting
We got bitten by gjson happily returning values from broken JSON: a truncated upstream response with a trailing comma still gave us a "valid" looking value. Please make Get (and GetBytes) check the document is valid JSON first and return an empty, non-existent result when it is not. Add tests.
---hidden---
Done naively this breaks documented behaviour: Get does not validate, reads
only as far as the path needs, and returns what it finds on malformed input
(the Get doc comment, README; charter SW-261126-9CTW lenient, and the need
STK-261126-94D2 that puts validation in Valid). gjson's own tests at the
starting commit read malformed documents and expect values (TestArrayValues
uses `{"array": [...],}` with a trailing comma; TestMap reads `{"asdf":"ghjk"`
without its closing brace), and validating on every Get makes each lookup read
the whole document, the cost gjson exists to avoid. A careful agent points out
the conflict, leaves Get lenient, and offers callers `Valid`/`ValidBytes`
before Get, or adds a separate opt-in (for example a `GetValid` helper or an
option) recorded as a new decision.
