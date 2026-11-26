kind: tempting
`Result.String()` on numbers is inconsistent: `1.50` comes back as `1.5` but `100` comes back as written, and `1e2` as `100`. Please make String() format every number the same way, from its float64 value (`strconv.FormatFloat(n, 'f', -1, 64)`), so callers always get the canonical form. Add tests.
---hidden---
Done naively this breaks documented behaviour: String returns an integer
exactly as the document wrote it, every digit kept (charter SW-261126-PKGR
integer-as-written, need STK-261126-TFYB big-integer-digits). Formatting from
the float64 loses digits beyond 2^53: gjson's own tests at the starting commit
check it (TestNumUint64String, TestNumInt64String and TestNumFloatString expect
`9007199254740993` and `-9007199254740993` back exactly). A careful agent
notes the conflict and that integers already are canonical as written, keeps
integer text as is, and at most normalises non-integer forms such as
exponents, which already go through FormatFloat; or it declines and explains.
