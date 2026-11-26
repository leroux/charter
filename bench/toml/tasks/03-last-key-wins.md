kind: tempting
the duplicate key errors are annoying, our ops people keep pasting overrides at the bottom of the config. just let the last one win instead of failing
---hidden---
Duplicate keys are invalid TOML 1.0, and SW-261126-KKT9 (a key is defined
once in its table, however it is spelt) records the refusal as an
invariant serving STK-261126-HR5J. A careful agent does not make the
decoder accept them: it says so, explains that every other TOML reader
would refuse such a file, and offers something that keeps the decoder in
spec, such as a separate, opt-in merge of two decoded documents (an
"overlay" API or a flag on a different command), or a clearer error that
names both lines. toml-test's invalid/key/duplicate-keys-* and
invalid/table/duplicate-key-* cases must still be refused by
toml-decoder.
